package logging

import (
	"testing"

	"github.com/shashank/home-server/common/config"
)

func TestInitLogger(t *testing.T) {
	cfg := config.LoggingConfig{
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}

	err := InitLogger(cfg, "test-service")
	if err != nil {
		t.Fatalf("InitLogger failed: %v", err)
	}

	Log.Debug("Debug message")
	Log.Info("Info message")
	Log.Warn("Warning message")
	Log.Error("Error message")
}
