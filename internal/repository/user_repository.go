package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/magistraapta/golang-devops/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uuid.UUID) error
	GetAllUsers() ([]model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username":   user.Username,
		"email":      user.Email,
		"updated_at": time.Now(),
	}).Error
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, id).Where("id = ?", id).Error
}
