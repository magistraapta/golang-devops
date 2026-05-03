package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	createdUser *model.User
	users       []model.User
}

func (r *fakeUserRepository) CreateUser(user *model.User) error {
	copiedUser := *user
	r.createdUser = &copiedUser
	return nil
}

func (r *fakeUserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (r *fakeUserRepository) GetUserByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (r *fakeUserRepository) UpdateUser(user *model.User) error {
	return nil
}

func (r *fakeUserRepository) DeleteUser(id uuid.UUID) error {
	return nil
}

func (r *fakeUserRepository) GetAllUsers() ([]model.User, error) {
	return r.users, nil
}

func TestUserServiceCreateUserSetsIDTimestampsAndHashesPassword(t *testing.T) {
	const originalPassword = "secret123"

	userRepository := &fakeUserRepository{}
	userService := service.NewUserService(userRepository)
	user := &model.User{
		Username: "Alice",
		Email:    "alice@example.com",
		Password: originalPassword,
	}

	if err := userService.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	createdUser := userRepository.createdUser
	if createdUser == nil {
		t.Fatal("expected repository CreateUser to be called")
	}
	if createdUser.ID == uuid.Nil {
		t.Fatal("expected user ID to be generated")
	}
	if createdUser.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
	if createdUser.UpdatedAt.IsZero() {
		t.Fatal("expected UpdatedAt to be set")
	}
	if createdUser.Password == originalPassword {
		t.Fatal("expected password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.Password), []byte(originalPassword)); err != nil {
		t.Fatalf("expected hashed password to match original password: %v", err)
	}
}

func TestUserServiceGetAllUsersDelegatesToRepository(t *testing.T) {
	expectedUsers := []model.User{
		{ID: uuid.New(), Username: "Alice", Email: "alice@example.com"},
	}
	userRepository := &fakeUserRepository{users: expectedUsers}
	userService := service.NewUserService(userRepository)

	users, err := userService.GetAllUsers()
	if err != nil {
		t.Fatalf("get all users: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Email != expectedUsers[0].Email {
		t.Fatalf("expected email %q, got %q", expectedUsers[0].Email, users[0].Email)
	}
}
