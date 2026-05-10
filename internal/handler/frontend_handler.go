package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/service"
)

type FrontendHandler interface {
	HomeHandler(c *gin.Context)
	UserDetailHandler(c *gin.Context)
}

type frontendHandler struct {
	userService service.UserService
}

func NewFrontendHandler(userService service.UserService) FrontendHandler {
	return &frontendHandler{userService: userService}
}

func (h *frontendHandler) HomeHandler(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error",
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home",
		"Users": users,
	})
}

func (h *frontendHandler) UserDetailHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)

	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title": "Error",
			"Error": err.Error(),
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title": "Error",
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "user-details.html", gin.H{
		"Title": user.Username,
		"User":  user,
	})
}
