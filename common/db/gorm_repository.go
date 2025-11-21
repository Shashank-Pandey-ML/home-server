package db

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GormRepository provides common CRUD operations using GORM
type GormRepository[T any] struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewGormRepository creates a new repository for the given model type
func NewGormRepository[T any](db *DB) *GormRepository[T] {
	return &GormRepository[T]{
		db:     db.DB,
		logger: db.logger,
	}
}

// Create creates a new record
func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Failed to create record", zap.Error(err))
		return err
	}
	return nil
}

// CreateBatch creates multiple records in a single query
func (r *GormRepository[T]) CreateBatch(ctx context.Context, entities []T, batchSize int) error {
	if err := r.db.WithContext(ctx).CreateInBatches(entities, batchSize).Error; err != nil {
		r.logger.Error("Failed to create batch records", zap.Error(err))
		return err
	}
	return nil
}

// GetByID retrieves a record by ID
func (r *GormRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to get record by ID", zap.Error(err), zap.Uint("id", id))
		return nil, err
	}
	return &entity, nil
}

// GetAll retrieves all records with optional preloading
func (r *GormRepository[T]) GetAll(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get all records", zap.Error(err))
		return nil, err
	}
	return entities, nil
}

// Update updates a record
func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Failed to update record", zap.Error(err))
		return err
	}
	return nil
}

// Delete soft deletes a record by ID
func (r *GormRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, id).Error; err != nil {
		r.logger.Error("Failed to delete record", zap.Error(err), zap.Uint("id", id))
		return err
	}
	return nil
}

// HardDelete permanently deletes a record by ID
func (r *GormRepository[T]) HardDelete(ctx context.Context, id uint) error {
	var entity T
	if err := r.db.WithContext(ctx).Unscoped().Delete(&entity, id).Error; err != nil {
		r.logger.Error("Failed to hard delete record", zap.Error(err), zap.Uint("id", id))
		return err
	}
	return nil
}

// FindWhere finds records based on conditions
func (r *GormRepository[T]) FindWhere(ctx context.Context, conditions map[string]interface{}) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := query.Find(&entities).Error; err != nil {
		r.logger.Error("Failed to find records with conditions", zap.Error(err))
		return nil, err
	}
	return entities, nil
}

// FindOneWhere finds a single record based on conditions
func (r *GormRepository[T]) FindOneWhere(ctx context.Context, conditions map[string]interface{}) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)

	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := query.First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("Failed to find record with conditions", zap.Error(err))
		return nil, err
	}
	return &entity, nil
}

// Count returns the count of records matching the conditions
func (r *GormRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := query.Count(&count).Error; err != nil {
		r.logger.Error("Failed to count records", zap.Error(err))
		return 0, err
	}
	return count, nil
}

// Exists checks if a record exists with the given conditions
func (r *GormRepository[T]) Exists(ctx context.Context, conditions map[string]interface{}) (bool, error) {
	count, err := r.Count(ctx, conditions)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Paginate returns paginated results
func (r *GormRepository[T]) Paginate(ctx context.Context, page, pageSize int, conditions map[string]interface{}) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.WithContext(ctx).Model(new(T))

	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count records for pagination", zap.Error(err))
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		r.logger.Error("Failed to get paginated records", zap.Error(err))
		return nil, 0, err
	}

	return entities, total, nil
}

// CustomQuery executes a custom query
func (r *GormRepository[T]) CustomQuery(ctx context.Context, query string, args ...interface{}) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Raw(query, args...).Find(&entities).Error; err != nil {
		r.logger.Error("Failed to execute custom query", zap.Error(err), zap.String("query", query))
		return nil, err
	}
	return entities, nil
}
