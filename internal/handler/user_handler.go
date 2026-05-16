package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/service"
)

type UserHandler interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserByEmail(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *userHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUserByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uuid.MustParse(id) // ← force the correct ID after binding

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userService.DeleteUser(uuid.MustParse(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
