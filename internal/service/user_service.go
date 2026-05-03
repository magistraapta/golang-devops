package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"github.com/magistraapta/golang-devops/internal/repository"
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
	return s.userRepository.GetAllUsers()
}

func (s *userService) CreateUser(user *model.User) error {
	// increment UUID
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepository.CreateUser(user)
}

func (s *userService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *userService) UpdateUser(user *model.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.userRepository.DeleteUser(id)
}
