package main

import (
	"net/http"
	"os"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/magistraapta/golang-devops/api"
	"github.com/magistraapta/golang-devops/internal/config"
)

func main() {
	config.LoadEnv()
	router := gin.Default()

	// root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	// health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server is running properly"})
	})

	api.UserRoutes(router)

	slog.Info("Starting server on port " + os.Getenv("PORT"))
	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		slog.Error("Error starting server:", "error", err)
		os.Exit(1)
	}
}
