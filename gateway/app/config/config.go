package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

func init() {
	// Load the config first, which will initialize all the configuration fields
	LoadConfig("config.yaml")
}

type Config struct {
	Service struct {
		Name           string
		Port           int
		Environment    string
		DeploymentMode string `mapstructure:"deployment_mode"`
	}
	Logging struct {
		Output string
	}
	Api struct {
		Timeout    int
		MaxRetries int `mapstructure:"max_retries"`
		RetryDelay int `mapstructure:"retry_delay"`
	}
}

func LoadConfig(path string) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}
