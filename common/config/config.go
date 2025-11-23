package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ServiceConfig holds metadata about the running microservice.
type ServiceConfig struct {
	Name        string `mapstructure:"name"`        // Logical name of the service (e.g., "auth", "catalog").
	Port        int    `mapstructure:"port"`        // Port on which the service will listen (e.g., 8080).
	Environment string `mapstructure:"environment"` // Deployment environment: "dev", "staging", or "prod".
}

// LoggingConfig controls the behavior of the application logger.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // Logging level: "debug", "info", "warn", or "error".
	Format string `mapstructure:"format"` // Log format: "json" for structured logs or "text" for console logs.
	Output string `mapstructure:"output"` // Destination for logs: "stdout", "stderr", or a "file".
}

// DatabaseConfig provides connection parameters for the backing database.
type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // Type of database: "postgresql", "mysql", "sqlite", etc.
	Host     string `mapstructure:"host"` // Hostname or IP address of the database (e.g., "db.example.com").
	Port     int    `mapstructure:"port"` // Port number the DB listens on (default for PostgreSQL is 5432).
	Name     string `mapstructure:"name"` // Name of the database to connect to.
	User     string `mapstructure:"user"` // Database user for authentication.
	Password string // Database password, loaded securely via environment variable.
	SSLMode  string `mapstructure:"ssl_mode"` // SSL mode for the connection: "disable", "require", "verify-ca", etc.
}

// APIConfig sets the behavior of the service's outbound or internal API communication.
type APIConfig struct {
	BaseURL    string        `mapstructure:"base_url"`    // Base URL for exposed APIs (e.g., "/api/v1").
	Timeout    time.Duration `mapstructure:"timeout"`     // Request timeout duration (e.g., "30s", "1m").
	MaxRetries int           `mapstructure:"max_retries"` // Number of retry attempts for failed requests.
}

// SecurityConfig defines security-related settings such as TLS and CORS.
type SecurityConfig struct {
	EnableTLS bool   `mapstructure:"enable_tls"` // If true, TLS is enabled; requires cert and key files.
	CertFile  string `mapstructure:"cert_file"`  // Path to the TLS certificate file.
	KeyFile   string `mapstructure:"key_file"`   // Path to the TLS private key file.
}

// JWTConfig defines JWT token configuration for authentication services.
type JWTConfig struct {
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`  // Duration for access tokens (e.g., "30m", "1h").
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"` // Duration for refresh tokens (e.g., "168h", "7d").
	Issuer               string        `mapstructure:"issuer"`                 // JWT issuer identifier.
	KeySize              int           `mapstructure:"key_size"`               // RSA key size for JWT signing (e.g., 2048, 4096).
	KeyFile              string        `mapstructure:"key_file"`               // Path to the JWT private key file.
	AllowedOrigins       []string      `mapstructure:"allowed_origins"`        // List of allowed origins for CORS (e.g., ["https://example.com"]).
}

// Config aggregates all other configurations into a single structure.
type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`  // Service-related configuration.
	Logging  LoggingConfig  `mapstructure:"logging"`  // Logging configuration.
	Database DatabaseConfig `mapstructure:"database"` // Database connection settings.
	API      APIConfig      `mapstructure:"api"`      // API-related configuration.
	Security SecurityConfig `mapstructure:"security"` // Security/TLS/CORS configuration.
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT authentication configuration.
}

// AppConfig is the globally accessible parsed configuration for the running service.
var AppConfig *Config

// DefaultConfigPath is the default file path where config.yaml is expected to be found.
const DefaultConfigPath = "config/config.yaml"

// LoadConfig reads the configuration from the specified file path and environment variables.
func LoadConfig(configPath string) error {
	setDefaults()

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Load secrets from environment variable
	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	AppConfig = &cfg
	return nil
}

// setDefaults initializes default values for the configuration.
func setDefaults() {
	viper.SetDefault("service.port", 8080)
	viper.SetDefault("service.environment", "prod")

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	viper.SetDefault("api.base_url", "/api/v1")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.max_retries", 3)

	viper.SetDefault("security.enable_tls", true)
	viper.SetDefault("security.cert_file", "cert.pem")
	viper.SetDefault("security.key_file", "key.pem")

	// Database defaults
	viper.SetDefault("database.ssl_mode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.access_token_duration", "30m")
	viper.SetDefault("jwt.refresh_token_duration", "168h") // 7 days
	viper.SetDefault("jwt.issuer", "home-server-auth")
	viper.SetDefault("jwt.key_size", 2048)
	viper.SetDefault("jwt.key_file", "jwt_key.pem")
	// Default allowed origins for CORS, can be overridden in config.yaml
	viper.SetDefault("jwt.allowed_origins", []string{})
}
