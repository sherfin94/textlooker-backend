package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email             string `gorm:"not null;unique" validate:"required,email"`
	EncryptedPassword string `gorm:"not null" validate:"required,min=8,max=20"`
}

func NewUser(email string, password string) (*User, error) {
	userValidator := validator.New()
	user := &User{Email: email, EncryptedPassword: password}
	err := userValidator.Struct(user)

	if err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic("Password hashing failed. Please check.")
	} else {
		user.EncryptedPassword = string(hashedPassword)
	}

	db := ConnectDatabase()
	result := db.Create(user)

	return user, result.Error
}
