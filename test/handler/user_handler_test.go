package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/handler"
	"github.com/magistraapta/golang-devops/internal/model"
)

type fakeUserService struct {
	createErr   error
	createdUser *model.User
	users       []model.User
	deleteErr   error
	updateErr   error // ← add this
}

func (s *fakeUserService) CreateUser(ctx context.Context, user *model.User) error {
	copiedUser := *user
	s.createdUser = &copiedUser
	return s.createErr
}

func (s *fakeUserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *fakeUserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (s *fakeUserService) UpdateUser(ctx context.Context, user *model.User) error {
	return nil
}

func (s *fakeUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.deleteErr
}

func (s *fakeUserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.users, nil
}

func newTestRouter(userService *fakeUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.POST("/users/", userHandler.CreateUser)
	router.GET("/users/", userHandler.GetAllUsers)
	router.GET("/users/:id", userHandler.GetUserByID)

	return router
}

func performRequest(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func TestUserHandlerCreateUserReturnsCreated(t *testing.T) {
	userService := &fakeUserService{}
	router := newTestRouter(userService)
	body := `{"username":"Alice","email":"alice@example.com","password":"secret123"}`

	response := performRequest(router, http.MethodPost, "/users/", body)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, response.Code, response.Body.String())
	}
	createdUser := userService.createdUser
	if createdUser == nil {
		t.Fatal("expected service CreateUser to be called")
	}
	if createdUser.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", createdUser.Email)
	}
	responseBody := response.Body.String()
	if !strings.Contains(responseBody, "User created successfully") {
		t.Fatalf("expected success message, got %s", responseBody)
	}
}

func TestUserHandlerCreateUserReturnsBadRequestForInvalidJSON(t *testing.T) {
	userService := &fakeUserService{}
	router := newTestRouter(userService)

	response := performRequest(router, http.MethodPost, "/users/", `{"email":`)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
	if userService.createdUser != nil {
		t.Fatal("expected service CreateUser not to be called")
	}
}

func TestUserHandlerCreateUserReturnsInternalServerError(t *testing.T) {
	userService := &fakeUserService{createErr: errors.New("create failed")}
	router := newTestRouter(userService)
	body := `{"username":"Alice","email":"alice@example.com","password":"secret123"}`

	response := performRequest(router, http.MethodPost, "/users/", body)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
	responseBody := response.Body.String()
	if !strings.Contains(responseBody, "create failed") {
		t.Fatalf("expected error body, got %s", responseBody)
	}
}

func TestUserHandlerGetAllUsersReturnsData(t *testing.T) {
	userService := &fakeUserService{
		users: []model.User{
			{ID: uuid.New(), Username: "Alice", Email: "alice@example.com"},
		},
	}
	router := newTestRouter(userService)

	response := performRequest(router, http.MethodGet, "/users/", "")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	responseBody := response.Body.String()
	if !strings.Contains(responseBody, "alice@example.com") {
		t.Fatalf("expected user email in response, got %s", responseBody)
	}
}

func TestGetUserDetailByID(t *testing.T) {
	// ✅ Save the UUID so we can use it in the URL
	aliceID := uuid.New()

	userService := &fakeUserService{
		users: []model.User{
			{ID: aliceID, Username: "Alice", Email: "alice@example.com"},
			{ID: uuid.New(), Username: "Bob", Email: "bob@example.com"},
		},
	}
	router := newTestRouter(userService)

	// ✅ Use the real UUID in the URL path
	response := performRequest(router, http.MethodGet, "/users/"+aliceID.String(), "")

	// ✅ Assert correct status
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	body := response.Body.String()

	// ✅ Only assert the requested user appears
	if !strings.Contains(body, "alice@example.com") {
		t.Fatalf("expected alice@example.com in response, got %s", body)
	}

	// ✅ Assert the OTHER user does NOT appear
	if strings.Contains(body, "bob@example.com") {
		t.Fatalf("bob@example.com should not appear in a single user response")
	}
}

func TestUserHandlerGetUserByID_NotFound(t *testing.T) {
	userService := &fakeUserService{users: []model.User{}}
	router := newTestRouter(userService)

	response := performRequest(router, http.MethodGet, "/users/"+uuid.New().String(), "")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
	if !strings.Contains(response.Body.String(), "User not found") {
		t.Fatalf("expected 'User not found' in response, got %s", response.Body.String())
	}
}

func TestUserHandlerUpdateUser_Success(t *testing.T) {
	userID := uuid.New()
	userService := &fakeUserService{
		users: []model.User{
			{ID: userID, Username: "Alice", Email: "alice@example.com"},
		},
	}
	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.PUT("/users/:id", userHandler.UpdateUser)

	body := `{"username":"Alice_Updated","email":"newalice@example.com"}`
	response := performRequest(router, http.MethodPut, "/users/"+userID.String(), body)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if !strings.Contains(response.Body.String(), "User updated successfully") {
		t.Fatalf("expected success message, got %s", response.Body.String())
	}
}

func TestUserHandlerUpdateUser_InvalidJSON(t *testing.T) {
	userID := uuid.New()
	userService := &fakeUserService{
		users: []model.User{
			{ID: userID, Username: "Alice", Email: "alice@example.com"},
		},
	}
	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.PUT("/users/:id", userHandler.UpdateUser)

	response := performRequest(router, http.MethodPut, "/users/"+userID.String(), `{"username":`)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestUserHandlerUpdateUser_NotFound(t *testing.T) {
	userService := &fakeUserService{users: []model.User{}}
	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.PUT("/users/:id", userHandler.UpdateUser)

	body := `{"username":"Updated","email":"updated@example.com"}`
	response := performRequest(router, http.MethodPut, "/users/"+uuid.New().String(), body)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}

func TestUserHandlerDeleteUser_Success(t *testing.T) {
	userService := &fakeUserService{}
	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	response := performRequest(router, http.MethodDelete, "/users/"+uuid.New().String(), "")

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if !strings.Contains(response.Body.String(), "User deleted successfully") {
		t.Fatalf("expected success message, got %s", response.Body.String())
	}
}

func TestUserHandlerDeleteUser_Error(t *testing.T) {
	userService := &fakeUserService{}
	userService.deleteErr = errors.New("delete failed")

	router := gin.New()
	userHandler := handler.NewUserHandler(userService)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	response := performRequest(router, http.MethodDelete, "/users/"+uuid.New().String(), "")

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, response.Code)
	}
}
