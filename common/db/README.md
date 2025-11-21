# Database Package

A comprehensive PostgreSQL database package for Go applications that provides connection management, query building, migrations, and common database operations.

## Features

- **Connection Management**: PostgreSQL connection pooling with pgx/v5
- **Query Builders**: Fluent API for building SELECT, INSERT, UPDATE, DELETE queries
- **Migrations**: Database schema migration management
- **Repository Pattern**: Base repository with common CRUD operations
- **Transactions**: Safe transaction handling with automatic rollback
- **Health Checks**: Database health monitoring
- **Utility Functions**: Common database helpers and utilities
- **Type Safety**: Strong typing with custom types for JSONB, nullable fields
- **Pagination**: Built-in pagination support
- **Soft Deletes**: Support for soft delete patterns
- **Audit Logging**: Built-in audit trail capabilities

## Installation

```bash
go get github.com/jackc/pgx/v5
go get go.uber.org/zap
```

## Quick Start

### 1. Database Configuration

```go
import "github.com/shashank/home-server/common/db"

// Create configuration
config := db.DefaultConfig()
config.Host = "localhost"
config.Port = 5432
config.Username = "your_user"
config.Password = "your_password"
config.Database = "your_database"

// Create logger
logger, _ := zap.NewDevelopment()

// Connect to database
database, err := db.NewDB(config, logger)
if err != nil {
    log.Fatal("Failed to connect:", err)
}
defer database.Close()
```

### 2. Running Migrations

```go
// Create migration manager
migrationManager := db.NewMigrationManager(database)

// Add migrations
common := &db.CommonMigrations{}
migrationManager.AddMigration("001", "Create UUID extension", common.AddUUIDExtension())
migrationManager.AddMigration("002", "Create users table", common.CreateUsersTable())
migrationManager.AddMigration("003", "Create sessions table", common.CreateSessionsTable())

// Run migrations
ctx := context.Background()
if err := migrationManager.MigrateUp(ctx); err != nil {
    log.Fatal("Migration failed:", err)
}
```

### 3. Using Query Builders

```go
// SELECT query
selectBuilder := db.NewSelectBuilder("users").
    Select("id", "email", "first_name", "last_name").
    Where("is_active = $1", true).
    Where("created_at > $2", time.Now().AddDate(0, 0, -30)).
    OrderBy("created_at DESC").
    Limit(10)

query, args := selectBuilder.Build()
rows, err := database.Pool.Query(ctx, query, args...)

// INSERT query
insertBuilder := db.NewInsertBuilder("users").
    Columns("email", "first_name", "last_name").
    Values("user@example.com", "John", "Doe")

query, args = insertBuilder.Build()
_, err = database.Pool.Exec(ctx, query, args...)

// UPDATE query
updateBuilder := db.NewUpdateBuilder("users").
    Set("first_name", "Jane").
    Set("updated_at", time.Now()).
    Where("id = $1", userID)

query, args = updateBuilder.Build()
_, err = database.Pool.Exec(ctx, query, args...)
```

### 4. Using Transactions

```go
err := database.Transaction(ctx, func(tx pgx.Tx) error {
    // Insert user
    var userID string
    err := tx.QueryRow(ctx, 
        "INSERT INTO users (email, first_name) VALUES ($1, $2) RETURNING id",
        "user@example.com", "John",
    ).Scan(&userID)
    if err != nil {
        return err
    }
    
    // Insert related data
    _, err = tx.Exec(ctx,
        "INSERT INTO user_profiles (user_id, bio) VALUES ($1, $2)",
        userID, "User bio",
    )
    return err
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### 5. Repository Pattern

```go
// Define your entity
type User struct {
    ID        string     `json:"id" db:"id"`
    Email     string     `json:"email" db:"email"`
    FirstName string     `json:"first_name" db:"first_name"`
    LastName  string     `json:"last_name" db:"last_name"`
    CreatedAt time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Implement entity interfaces
func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(string) }
func (u *User) SetCreatedAt(t time.Time) { u.CreatedAt = t }
func (u *User) SetUpdatedAt(t time.Time) { u.UpdatedAt = t }
func (u *User) GetCreatedAt() time.Time { return u.CreatedAt }
func (u *User) GetUpdatedAt() time.Time { return u.UpdatedAt }
func (u *User) SetDeletedAt(t time.Time) { u.DeletedAt = &t }
func (u *User) GetDeletedAt() *time.Time { return u.DeletedAt }
func (u *User) IsDeleted() bool { return u.DeletedAt != nil }

// Create repository
type UserRepository struct {
    *db.BaseRepository
}

func NewUserRepository(database *db.DB) *UserRepository {
    return &UserRepository{
        BaseRepository: db.NewBaseRepository(database, "users"),
    }
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (email, first_name, last_name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    return r.DB().Pool.QueryRow(ctx, query,
        user.Email, user.FirstName, user.LastName,
        time.Now(), time.Now(),
    ).Scan(&user.ID)
}
```

### 6. Pagination and Filtering

```go
// Define query options
opts := db.QueryOptions{
    Page:     1,
    PageSize: 20,
    Filters: []db.Filter{
        {Field: "is_active", Operator: "=", Value: true},
        {Field: "created_at", Operator: ">", Value: time.Now().AddDate(0, 0, -30)},
    },
    Sorts: []db.Sort{
        {Field: "created_at", Order: "DESC"},
    },
}

// Use in repository
result, err := userRepo.List(ctx, opts)
if err != nil {
    log.Fatal(err)
}

log.Printf("Total users: %d", result.Total)
log.Printf("Current page: %d of %d", result.Page, result.TotalPages)
for _, user := range result.Data {
    log.Printf("User: %+v", user)
}
```

### 7. Health Checks

```go
// Simple ping
if err := database.Ping(ctx); err != nil {
    log.Printf("Database is down: %v", err)
}

// Detailed health check
health := database.HealthCheck(ctx)
log.Printf("Database status: %s", health["status"])
log.Printf("Total connections: %v", health["total_connections"])
log.Printf("Idle connections: %v", health["idle_connections"])
```

### 8. Upsert Operations

```go
upsertBuilder := db.NewUpsertBuilder("users").
    Columns("email", "first_name", "last_name").
    OnConflict("email").
    DoUpdate("first_name", "last_name")

err := database.Upsert(ctx, upsertBuilder, 
    "user@example.com", "Jane", "Smith",
)
```

## Configuration Options

```go
type Config struct {
    Host            string        // Database host
    Port            int           // Database port
    Username        string        // Database username
    Password        string        // Database password
    Database        string        // Database name
    SSLMode         string        // SSL mode (disable, require, etc.)
    MaxConnections  int32         // Maximum connections in pool
    MinConnections  int32         // Minimum connections in pool
    MaxConnLifetime time.Duration // Maximum connection lifetime
    MaxConnIdleTime time.Duration // Maximum connection idle time
    ConnectTimeout  time.Duration // Connection timeout
}
```

## Query Builder Examples

### Complex SELECT with Joins

```go
query, args := db.NewSelectBuilder("users u").
    Select("u.id", "u.email", "p.bio", "r.name as role_name").
    LeftJoin("user_profiles p", "p.user_id = u.id").
    InnerJoin("user_roles ur", "ur.user_id = u.id").
    InnerJoin("roles r", "r.id = ur.role_id").
    Where("u.is_active = $1", true).
    Where("u.created_at > $2", time.Now().AddDate(0, 0, -30)).
    GroupBy("u.id", "u.email", "p.bio", "r.name").
    Having("COUNT(ur.role_id) > $3", 0).
    OrderBy("u.created_at DESC").
    Limit(50).
    Build()
```

### Pagination Builder

```go
paginationBuilder := db.NewPaginationBuilder("users", 2, 20). // page 2, 20 items per page
    Select("id", "email", "first_name").
    Where("is_active = $1", true).
    OrderBy("created_at DESC")

// Get paginated results
query, args := paginationBuilder.BuildWithPagination()
rows, err := database.Pool.Query(ctx, query, args...)

// Get total count for pagination
countQuery, countArgs := paginationBuilder.BuildCountQuery()
var total int64
err = database.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
```

## Migration Patterns

### Custom Migration

```go
// Add custom migration
migrationManager.AddMigration("005", "Add user preferences", `
    CREATE TABLE user_preferences (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
        key VARCHAR(100) NOT NULL,
        value JSONB,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW(),
        UNIQUE(user_id, key)
    );
    
    CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
    CREATE INDEX idx_user_preferences_key ON user_preferences(key);
`)
```

### Common Migration Helpers

```go
common := &db.CommonMigrations{}

// Create extensions
migrationManager.AddMigration("001", "Extensions", common.AddUUIDExtension())

// Create standard tables
migrationManager.AddMigration("002", "Users table", common.CreateUsersTable())
migrationManager.AddMigration("003", "Sessions table", common.CreateSessionsTable())
migrationManager.AddMigration("004", "Audit logs", common.CreateAuditLogsTable())

// Add triggers
migrationManager.AddMigration("005", "Timestamp triggers", common.AddTimestampTriggers())

// Create indexes (production-safe)
migrationManager.AddMigration("006", "User indexes", 
    common.CreateIndexesConcurrently("users", "email,created_at"))
```

## Utilities

### JSONB Support

```go
// Define JSONB field
type UserPreferences struct {
    Theme    string `json:"theme"`
    Language string `json:"language"`
    Settings map[string]interface{} `json:"settings"`
}

// Use in struct
type User struct {
    ID          string                 `db:"id"`
    Email       string                 `db:"email"`
    Preferences db.JSONB              `db:"preferences"`
}

// Insert with JSONB
preferences := db.JSONB{
    "theme": "dark",
    "language": "en",
    "notifications": map[string]interface{}{
        "email": true,
        "sms": false,
    },
}

_, err = database.Pool.Exec(ctx,
    "INSERT INTO users (email, preferences) VALUES ($1, $2)",
    "user@example.com", preferences,
)
```

### Struct Mapping

```go
// Convert struct to map for dynamic queries
user := &User{Email: "test@example.com", FirstName: "John"}
data := db.StructToMap(user)
// data = {"email": "test@example.com", "first_name": "John", ...}

// Convert map to struct
data := map[string]interface{}{
    "email": "test@example.com",
    "first_name": "John",
}
user := &User{}
err := db.MapToStruct(data, user)
```

## Error Handling

```go
// Check for no rows error
user, err := userRepo.GetByEmail(ctx, email)
if db.IsNoRowsError(err) {
    // User not found
    return nil, nil
} else if err != nil {
    // Other error
    return nil, err
}

// Handle no rows gracefully
err = db.HandleNoRowsError(err) // Returns nil if no rows, otherwise returns original error
```

## Best Practices

1. **Always use context**: Pass context for timeout and cancellation support
2. **Use transactions**: For operations that modify multiple tables
3. **Handle errors properly**: Check for specific error types
4. **Use connection pooling**: Configure appropriate pool sizes
5. **Run migrations safely**: Test migrations in staging first
6. **Monitor health**: Implement health checks in your services
7. **Use prepared statements**: Query builders generate parameterized queries
8. **Index properly**: Use migration helpers to create indexes
9. **Soft deletes**: Implement soft deletes for audit trails
10. **Log queries**: Use the debug logging for development

## Testing

The package includes example tests that demonstrate how to set up and test database operations. For actual testing, you'll need to:

1. Set up a test database
2. Run migrations
3. Use repositories for testing
4. Clean up after tests

```go
func TestUserRepository(t *testing.T) {
    // Set up test database connection
    database := setupTestDB(t)
    defer database.Close()
    
    // Run migrations
    migrationManager := db.NewMigrationManager(database)
    // ... add migrations
    err := migrationManager.MigrateUp(context.Background())
    require.NoError(t, err)
    
    // Test repository operations
    userRepo := NewUserRepository(database)
    ctx := context.Background()
    
    // Test create
    user := &User{Email: "test@example.com", FirstName: "Test"}
    err = userRepo.Create(ctx, user)
    require.NoError(t, err)
    require.NotEmpty(t, user.ID)
    
    // Test get
    found, err := userRepo.GetByEmail(ctx, "test@example.com")
    require.NoError(t, err)
    require.Equal(t, user.Email, found.Email)
    
    // Clean up
    cleanupTestDB(t, database)
}
```

This database package provides a solid foundation for PostgreSQL operations in your Go applications with proper error handling, connection management, and common patterns built-in.
