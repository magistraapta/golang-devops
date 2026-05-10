package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magistraapta/golang-devops/internal/config"
	"github.com/magistraapta/golang-devops/internal/handler"
	"github.com/magistraapta/golang-devops/internal/repository"
	"github.com/magistraapta/golang-devops/internal/service"
)

func FrontendRoutes(router *gin.Engine) {
	db := config.DatabaseConnection()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	frontendHandler := handler.NewFrontendHandler(userService)

	// frontend routes
	frontendGroup := router.Group("/")
	{
		frontendGroup.GET("/", frontendHandler.HomeHandler)
		frontendGroup.GET("/users/:id", frontendHandler.UserDetailHandler)

	}
}
