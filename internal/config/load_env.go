package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	paths := []string{
		".env.local",
		"../.env.local",
		"../../.env.local",
		".env",
		"../.env",
		"../../.env",
	}

	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			slog.Info("Environment variables loaded", "path", path)
			return nil
		}
	}

	slog.Error("No environment file found in any search path")
	return nil
}
