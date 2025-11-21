package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/shashank/home-server/common/config"
)

// DB represents a database connection wrapper using GORM
type DB struct {
	*gorm.DB
	logger *zap.Logger
}

const (
	DB_MAX_IDLE_CONNECTIONS = 10
	DB_MAX_OPEN_CONNECTIONS = 100
	DB_MAX_CONN_LIFETIME    = time.Hour
)

// GormLogger wraps zap logger for GORM
type GormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
}

// NewGormLogger creates a new GORM logger that uses zap
func NewGormLogger(zapLogger *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger: zapLogger,
		LogLevel:  logger.Info,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLogger{
		ZapLogger: l.ZapLogger,
		LogLevel:  level,
	}
}

// Info logs info messages
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn logs warning messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error logs error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace logs SQL queries
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && l.LogLevel >= logger.Error {
		l.ZapLogger.Error("SQL query failed",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn("Slow SQL query",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	} else if l.LogLevel >= logger.Info {
		l.ZapLogger.Debug("SQL query executed",
			zap.Duration("elapsed", elapsed),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
		)
	}
}

// ConnectionString builds and returns the PostgreSQL connection string for GORM
func ConnectionString(c config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
	)
}

// InitDbConnection creates a new database connection with GORM
func InitDbConnection(config config.DatabaseConfig, logger *zap.Logger) (*DB, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	// Configure GORM with custom logger
	gormConfig := &gorm.Config{
		Logger: NewGormLogger(logger),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(ConnectionString(config)), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(DB_MAX_IDLE_CONNECTIONS)
	sqlDB.SetMaxOpenConns(DB_MAX_OPEN_CONNECTIONS)
	sqlDB.SetConnMaxLifetime(DB_MAX_CONN_LIFETIME)

	dbWrapper := &DB{
		DB:     db,
		logger: logger,
	}

	// Test the connection
	if err := dbWrapper.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established successfully",
		zap.String("database", config.Name),
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
	)

	return dbWrapper, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	db.logger.Info("Closing database connection")
	return sqlDB.Close()
}

// Ping tests the database connection
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Stats returns connection pool statistics
func (db *DB) Stats() map[string]interface{} {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
	}
}

// HealthCheck returns database health information
func (db *DB) HealthCheck(ctx context.Context) map[string]interface{} {
	health := map[string]interface{}{
		"status": "unknown",
	}

	// Add connection stats
	for k, v := range db.Stats() {
		health[k] = v
	}

	// Test connection
	if err := db.Ping(ctx); err != nil {
		health["status"] = "unhealthy"
		health["error"] = err.Error()
	} else {
		health["status"] = "healthy"
	}

	return health
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return db.DB.WithContext(ctx).Transaction(fn)
}

// AutoMigrate runs auto migration for given models
func (db *DB) AutoMigrate(models ...interface{}) error {
	db.logger.Info("Running auto migration")
	if err := db.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	db.logger.Info("Auto migration completed successfully")
	return nil
}

// Utilities for common database operations

// IsRecordNotFound checks if the error is a "record not found" error
func IsRecordNotFound(err error) bool {
	return err != nil && err == gorm.ErrRecordNotFound
}

// HandleRecordNotFound returns nil if the error is "record not found", otherwise returns the original error
func HandleRecordNotFound(err error) error {
	if IsRecordNotFound(err) {
		return nil
	}
	return err
}
