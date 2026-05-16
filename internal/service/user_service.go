package service

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/repository"
	"github.com/magistraapta/golang-devops/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetAllUsers(context context.Context) ([]model.User, error) {
	users, err := s.userRepository.GetAllUsers(context)
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get all users", Code: http.StatusInternalServerError}
	}
	return users, nil
}

func (s *userService) CreateUser(context context.Context, user *model.User) error {
	// increment UUID
	user.ID = uuid.New()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return &utils.CustomError{Message: "Failed to hash password", Code: http.StatusInternalServerError}
	}
	user.Password = string(hashedPassword)

	if err := s.userRepository.CreateUser(context, user); err != nil {
		return &utils.CustomError{Message: "Failed to create user", Code: http.StatusInternalServerError}
	}
	return nil
}

func (s *userService) GetUserByID(context context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(context, id)
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get user", Code: http.StatusInternalServerError}
	}
	return user, nil
}

func (s *userService) GetUserByEmail(context context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetUserByEmail(context, email)
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get user", Code: http.StatusInternalServerError}
	}
	return user, nil
}

func (s *userService) UpdateUser(context context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	if err := s.userRepository.UpdateUser(context, user); err != nil {
		return &utils.CustomError{Message: "Failed to update user", Code: http.StatusInternalServerError}
	}
	return nil
}

func (s *userService) DeleteUser(context context.Context, id uuid.UUID) error {
	if err := s.userRepository.DeleteUser(context, id); err != nil {
		return &utils.CustomError{Message: "Failed to delete user", Code: http.StatusInternalServerError}
	}
	return nil
}
