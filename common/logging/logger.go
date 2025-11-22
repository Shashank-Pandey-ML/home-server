package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/shashank/home-server/common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

const (
	DEFAULT_LOG_DIR       = "/app/logs"
	DEFAULT_LOG_FILE_NAME = "app.log"
)

func InitLogger(cfg config.LoggingConfig, serviceName string) error {
	var encoder zapcore.Encoder
	var zapLevel zapcore.Level

	// Set log level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Set encoder (json or console)
	var encoderConfig zapcore.EncoderConfig
	switch strings.ToLower(cfg.Format) {
	case "json":
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.CallerKey = "caller"
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Set output (stdout, stderr, or file)
	var output zapcore.WriteSyncer
	switch strings.ToLower(cfg.Output) {
	case "stdout":
		output = zapcore.Lock(os.Stdout)
	case "stderr":
		output = zapcore.Lock(os.Stderr)
	case "file":
		LOG_DIR := DEFAULT_LOG_DIR + "/" + serviceName
		if err := os.MkdirAll(LOG_DIR, 0755); err != nil {
			return fmt.Errorf("could not create log directory: %w", err)
		}
		logPath := fmt.Sprintf("%s/app.log", LOG_DIR)
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not open log file: %w", err)
		}
		output = zapcore.AddSync(file)
	default:
		return fmt.Errorf("invalid log output specified: %s", cfg.Output)
	}

	core := zapcore.NewCore(encoder, output, zapLevel)

	// Create logger with caller information
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0), zap.AddStacktrace(zapcore.ErrorLevel))
	Log.Info("Zap logger initialized",
		zap.String("level", zapLevel.String()),
		zap.String("format", cfg.Format),
		zap.String("service", serviceName))
	return nil
}
