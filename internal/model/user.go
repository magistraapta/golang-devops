package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" validate:"required,min=3,max=20" json:"username"`
	Email     string    `gorm:"unique" validate:"required,email" json:"email"`
	Password  string    `gorm:"not null" validate:"required,min=8" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" validate:"required" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" validate:"required" json:"updated_at"`
}
