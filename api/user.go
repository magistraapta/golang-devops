package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magistraapta/golang-devops/internal/config"
	"github.com/magistraapta/golang-devops/internal/handler"
	"github.com/magistraapta/golang-devops/internal/repository"
	"github.com/magistraapta/golang-devops/internal/service"
)

func UserRoutes(router *gin.Engine) {
	db := config.DatabaseConnection()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	// user routes
	userGroup := router.Group("/api/v1/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/", userHandler.GetAllUsers)
		userGroup.GET("/:id", userHandler.GetUserByID)
		userGroup.GET("/email/:email", userHandler.GetUserByEmail)
		userGroup.PUT("/:id", userHandler.UpdateUser)
		userGroup.DELETE("/:id", userHandler.DeleteUser)
	}

}
