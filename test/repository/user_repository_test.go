package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestUserRepository(t *testing.T) repository.UserRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate user model: %v", err)
	}

	return repository.NewUserRepository(db)
}

func TestUserRepositoryCreateAndGetAllUsers(t *testing.T) {
	userRepository := newTestUserRepository(t)
	user := &model.User{
		ID:       uuid.New(),
		Username: "Alice",
		Email:    "alice@example.com",
		Password: "hashed-password",
	}

	if err := userRepository.CreateUser(user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	users, err := userRepository.GetAllUsers()
	if err != nil {
		t.Fatalf("get all users: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].Email != user.Email {
		t.Fatalf("expected email %q, got %q", user.Email, users[0].Email)
	}
}

func TestUserRepositoryGetUserByEmailReturnsNilWhenMissing(t *testing.T) {
	userRepository := newTestUserRepository(t)

	user, err := userRepository.GetUserByEmail("missing@example.com")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}
