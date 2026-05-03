package config

import (
	"log/slog"
	"os"

	"github.com/magistraapta/golang-devops/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	db.AutoMigrate(&model.User{})
	slog.Info("connected to database")
	return db
}
