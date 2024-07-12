package config

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	LogLevel    string `default:"info" envconfig:"LOG_LEVEL"`
	LogFilePath string `default:"" envconfig:"LOG_FILE_PATH"`
}

func NewConfig(ctx context.Context) Config {
	zap.L().Info("Starting qgames")
	envPath := getEnvOrDefault("ENV_PATH", ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		zap.L().Error("failed to load env file", zap.Error(err))
	}

	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		zap.L().Fatal("failed to process env vars", zap.Error(err))
	}
	return *config
}

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
