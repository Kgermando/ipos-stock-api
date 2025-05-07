package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	UUID string `gorm:"type:text;not null;unique" json:"uuid"` // Explicitly set type:text

	Fullname        string `gorm:"not null" json:"fullname"`
	Email           string `json:"email" gorm:"unique;not null"`
	Title           string `json:"title"`
	Phone           string `json:"phone" gorm:"not null;unique"` // Added unique constraint
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" gorm:"-"`

	Role       string `json:"role"`
	Permission string `json:"permission"`
	Image      string `json:"image"`
	Status     bool   `json:"status"`
	Signature  string `json:"signature"`
}

type UserResponse struct {
	ID         uint      `json:"id"`
	UUID       string    `json:"uuid"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Title      string    `json:"title"`
	Role       string    `json:"role"`
	Permission string    `json:"permission"`
	Status     bool      `json:"status"`
	Signature  string    `json:"signature"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Login struct {
	// Email    string `json:"email" validate:"required,email"`
	// Phone    string `json:"phone" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (u *User) SetPassword(p string) {
	hp, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
	u.Password = string(hp)
}

func (u *User) ComparePassword(p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	return err
}
