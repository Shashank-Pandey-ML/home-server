package config

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleYAML = `
service:
  name: "auth"
  port: 8080
  environment: "prod"

logging:
  level: "info"
  format: "json"
  output: "stdout"

database:
  type: "postgresql"
  host: "localhost"
  port: 5432
  name: "authDB"
  user: "postgres"

api:
  base_url: "/api/v1"
  timeout: 30s
  max_retries: 3

security:
  enable_tls: true
  cert_file: "certs/server.crt"
  key_file: "certs/server.key"
  allowed_origins:
    - "https://example.com"
    - "https://another.com"

feature_flags:
  enable_feature_x: true
  enable_feature_y: false
`

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory and file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_config.yaml")

	if err := os.WriteFile(tmpFile, []byte(sampleYAML), 0644); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}

	// Set required environment variable
	os.Setenv("DB_PASSWORD", "test_secret")
	defer os.Unsetenv("DB_PASSWORD")

	// Load the config
	err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Validate config values
	if AppConfig == nil {
		t.Fatal("AppConfig is nil")
	}

	if AppConfig.Service.Name != "auth" {
		t.Errorf("Expected service.name 'auth', got '%s'", AppConfig.Service.Name)
	}

	if AppConfig.Database.Password != "test_secret" {
		t.Errorf("Expected DB_PASSWORD to be 'test_secret', got '%s'", AppConfig.Database.Password)
	}

	if !AppConfig.Security.EnableTLS {
		t.Errorf("Expected security.enable_tls to be true")
	}
}
