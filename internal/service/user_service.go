package service

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/repository"
	"github.com/magistraapta/golang-devops/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uuid.UUID) error
	GetAllUsers() ([]model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	users, err := s.userRepository.GetAllUsers()
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get all users", Code: http.StatusInternalServerError}
	}
	return users, nil
}

func (s *userService) CreateUser(user *model.User) error {
	// increment UUID
	user.ID = uuid.New()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return &utils.CustomError{Message: "Failed to hash password", Code: http.StatusInternalServerError}
	}
	user.Password = string(hashedPassword)

	if err := s.userRepository.CreateUser(user); err != nil {
		return &utils.CustomError{Message: "Failed to create user", Code: http.StatusInternalServerError}
	}
	return nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get user", Code: http.StatusInternalServerError}
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, &utils.CustomError{Message: "Failed to get user", Code: http.StatusInternalServerError}
	}
	return user, nil
}

func (s *userService) UpdateUser(user *model.User) error {
	user.UpdatedAt = time.Now()
	if err := s.userRepository.UpdateUser(user); err != nil {
		return &utils.CustomError{Message: "Failed to update user", Code: http.StatusInternalServerError}
	}
	return nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	if err := s.userRepository.DeleteUser(id); err != nil {
		return &utils.CustomError{Message: "Failed to delete user", Code: http.StatusInternalServerError}
	}
	return nil
}
