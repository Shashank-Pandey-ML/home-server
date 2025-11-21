package db

import (
	"context"
	"fmt"

	"github.com/shashank/home-server/common/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserRepository provides user-specific database operations
type UserRepository struct {
	*GormRepository[models.User]
	logger *zap.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		GormRepository: NewGormRepository[models.User](db),
		logger:         db.logger,
	}
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	return &user, nil
}

// EmailExists checks if a user with the given email exists
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		r.logger.Error("Failed to check if email exists", zap.Error(err), zap.String("email", email))
		return false, err
	}
	return count > 0, nil
}

// GetUsersByRole retrieves users by role with pagination
func (r *UserRepository) GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]models.User, int64, error) {
	conditions := map[string]interface{}{
		"role": role,
	}
	return r.Paginate(ctx, page, pageSize, conditions)
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		r.logger.Error("Failed to update password", zap.Error(result.Error), zap.Uint("user_id", userID))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}
	return nil
}

// SearchUsers searches users by name or email
func (r *UserRepository) SearchUsers(ctx context.Context, searchTerm string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}).Where(
		"name ILIKE ? OR email ILIKE ?",
		"%"+searchTerm+"%",
		"%"+searchTerm+"%",
	)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count users for search", zap.Error(err))
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		r.logger.Error("Failed to search users", zap.Error(err))
		return nil, 0, err
	}

	return users, total, nil
}

// GetActiveUsers retrieves users that are not soft deleted
func (r *UserRepository) GetActiveUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		r.logger.Error("Failed to get active users", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// SoftDeleteUser soft deletes a user
func (r *UserRepository) SoftDeleteUser(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, userID)
	if result.Error != nil {
		r.logger.Error("Failed to soft delete user", zap.Error(result.Error), zap.Uint("user_id", userID))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}
	return nil
}

// RestoreUser restores a soft deleted user
func (r *UserRepository) RestoreUser(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("id = ?", userID).Update("deleted_at", nil)
	if result.Error != nil {
		r.logger.Error("Failed to restore user", zap.Error(result.Error), zap.Uint("user_id", userID))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}
	return nil
}

// CreateUserIfNotExists creates a user only if email doesn't exist
func (r *UserRepository) CreateUserIfNotExists(ctx context.Context, user *models.User) (*models.User, error) {
	exists, err := r.EmailExists(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return r.GetByEmail(ctx, user.Email)
	}

	if err := r.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
