package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/shashank/home-server/common/config"
	"github.com/shashank/home-server/common/db"
	"github.com/shashank/home-server/common/logging"
	"github.com/shashank/home-server/common/models"
)

// JWT key pair for signing and validation
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// JWTClaims represents the claims stored in JWT tokens
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	Type    string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo *db.UserRepository
}

func NewAuthService(userRepo *db.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// initializeJWTKeys generates RSA key pair for JWT signing
func InitializeJWTKeys() error {
	// Get key size from config, fallback to 2048 if not set
	keySize := config.AppConfig.JWT.KeySize
	if keySize == 0 {
		keySize = 2048
	}

	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKey = privKey
	publicKey = &privKey.PublicKey

	logging.Log.Info("JWT keys initialized successfully",
		zap.Int("key_size", keySize))
	return nil
}

// Login handles user login and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, int64, error) {
	user, err := s.validateUserCredentials(ctx, email, password)
	if err != nil {
		return "", "", 0, err
	}

	// Generate token pair
	accessToken, refreshToken, expiresIn, err := GenerateTokenPair(user)
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, expiresIn, nil
}

// Logout handles user logout and token invalidation
func (s *AuthService) Logout(ctx context.Context, refreshToken string, userID uint) error {
	// We use Stateless JWT (most common in microservices). Here Backend does not store the token at all.
	// If you go purely stateless with JWTs, then there is no real “logout” on the backend — because the
	// backend never stores or tracks tokens. Once issued, the JWT is valid until it expires.

	// We invalidate the refresh token
	err := s.InvalidateRefreshToken(ctx, refreshToken, userID)
	if err != nil {
		logging.Log.Error("Failed to invalidate refresh token",
			zap.Uint("user_id", userID),
			zap.Error(err))
		return err
	}

	// Invalidate the refresh token
	return nil
}

// GenerateTokenPair creates both access and refresh tokens for a user
func GenerateTokenPair(user *models.User) (accessToken, refreshToken string, expiresIn int64, err error) {
	now := time.Now()

	// Get durations from config
	accessTokenDuration := config.AppConfig.JWT.AccessTokenDuration
	refreshTokenDuration := config.AppConfig.JWT.RefreshTokenDuration
	issuer := config.AppConfig.JWT.Issuer

	// Generate access token
	accessClaims := JWTClaims{
		UserID:  strconv.Itoa(int(user.ID)),
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		Type:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.Itoa(int(user.ID)),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(privateKey)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := JWTClaims{
		UserID:  strconv.Itoa(int(user.ID)),
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		Type:    "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   strconv.Itoa(int(user.ID)),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(privateKey)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	expiresIn = int64(accessTokenDuration.Seconds())

	logging.Log.Debug("Generated token pair",
		zap.Uint("user_id", user.ID),
		zap.Int64("expires_in", expiresIn))

	return accessToken, refreshToken, expiresIn, nil
}

// validateJWTToken validates and parses a JWT token
func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Validate token type (should be "access" for API requests)
	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// verifyPassword checks if the provided password matches the hash
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validateUserCredentials validates user email and password
func (s *AuthService) validateUserCredentials(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if !verifyPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// validateRefreshToken validates a refresh token and returns the user
func (s *AuthService) ValidateRefreshToken(ctx context.Context, tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token claims")
	}

	// Validate token type (should be "refresh")
	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type for refresh")
	}

	// TODO: Check if refresh token exists in database and is not revoked

	// Return user from claims (in production, fetch from database)
	userID, _ := strconv.ParseUint(claims.UserID, 10, 0)
	return &models.User{
		BaseModel: models.BaseModel{
			ID: uint(userID),
		},
		Email:   claims.Email,
		IsAdmin: claims.IsAdmin,
	}, nil
}

// invalidateRefreshToken marks a refresh token as invalid
func (s *AuthService) InvalidateRefreshToken(ctx context.Context, tokenString string, userID uint) error {
	// TODO: Implement database logic to mark token as revoked
	logging.Log.Info("Refresh token invalidated",
		zap.Uint("user_id", userID))
	return nil
}

// getUserByID fetches a user by their ID
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// checkDatabaseHealth verifies database connectivity
func (s *AuthService) CheckDatabaseHealth(ctx context.Context) error {
	// TODO: Implement actual database health check
	// For now, return success
	return nil
}

// GetPublicKeyPEM returns the public key in PEM format for the gateway service
func (s *AuthService) GetPublicKeyPEM(ctx context.Context) (string, error) {
	if publicKey == nil {
		return "", errors.New("public key not initialized")
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	return string(pubKeyPEM), nil
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

// UpdateUserProfile updates the user's profile information
func (s *AuthService) UpdateUserProfile(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}
