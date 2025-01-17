package model

import (
	"errors"

	"time"

	"github.com/geedotrar/mygram/pkg/helper"
	"gorm.io/gorm"
)

type User struct {
	ID           uint64         `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Password     string         `json:"-"`
	Dob          time.Time      `json:"dob" gorm:"column:dob"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
	Photos       []Photo        `json:"photos,omitempty"`
	Comments     []Comment      `json:"comments,omitempty"`
	SocialMedias []SocialMedia  `json:"social_medias,omitempty"`
}

type UserSignUp struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Dob      string `json:"dob" binding:"required"`
}
type UserUpdate struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	Dob       string    `json:"dob" binding:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}
type UserView struct {
	ID       uint64    `json:"id"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"-" binding:"required"`
	Email    string    `json:"email" binding:"required"`
	Dob      time.Time `json:"dob" binding:"required"`
}

func (u UserSignUp) Validate() error {
	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if !helper.IsValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
