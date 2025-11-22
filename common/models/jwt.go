package models

import "github.com/golang-jwt/jwt/v5"

// JWTClaims represents the custom claims for JWT tokens
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	Type    string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenType constants
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// PublicKeyResponse represents the response structure for public key endpoint
type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
	Algorithm string `json:"algorithm"`
	KeyType   string `json:"key_type"`
}
