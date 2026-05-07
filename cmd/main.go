package main

import (
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/magistraapta/golang-devops/api"
	"github.com/magistraapta/golang-devops/internal/config"
	"github.com/magistraapta/golang-devops/internal/metrics"
	"github.com/magistraapta/golang-devops/internal/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.LoadEnv()
	router := gin.Default()
	metrics.RunMetrics(10 * time.Second)
	router.Use(middleware.HTTPMiddleware())

	// root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	// health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server is running properly"})
	})

	api.UserRoutes(router)

	reg := prometheus.NewRegistry()
	metrics.RegisterAllMetrics(reg)
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))

	slog.Info("Starting server on port " + os.Getenv("PORT"))
	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		slog.Error("Error starting server:", "error", err)
		os.Exit(1)
	}
}
