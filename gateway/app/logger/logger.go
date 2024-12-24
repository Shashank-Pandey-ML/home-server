package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"gateway/app/config"
)

var Logger *zap.Logger

func init() {
	// Initialize the logger
	InitializeLogger()
}

// InitializeLogger initializes the global logger
func InitializeLogger() {
	writer := zapcore.AddSync(os.Stdout) // Output to console

	if config.AppConfig.Logging.Output == "file" {
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/gateway.log",
			MaxSize:    10, // Megabytes
			MaxBackups: 5,
			MaxAge:     30, // Days
		})
	}
	var core zapcore.Core
	if config.AppConfig.Service.Environment == "prod" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			writer,
			zap.InfoLevel,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			writer,
			zap.DebugLevel,
		)
	}
	Logger = zap.New(core)
	Logger = Logger.With(
		zap.String("service", config.AppConfig.Service.Name),
		zap.String("environment", config.AppConfig.Service.Environment),
	)
	defer Logger.Sync()
}
