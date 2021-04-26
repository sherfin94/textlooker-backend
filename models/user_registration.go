package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistration struct {
	gorm.Model
	Email             string `gorm:"not null;unique" validate:"required,email"`
	EncryptedPassword string `gorm:"not null" validate:"required,min=8,max=20"`
}

func NewUserRegistration(email string, password string) (*UserRegistration, error) {
	userRegistrationValidator := validator.New()
	userRegistration := &UserRegistration{Email: email, EncryptedPassword: password}
	err := userRegistrationValidator.Struct(userRegistration)

	if err != nil {
		return nil, err.(validator.ValidationErrors)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic("Password hashing failed. Please check.")
	} else {
		userRegistration.EncryptedPassword = string(hashedPassword)
	}

	result := Database.Create(userRegistration)

	return userRegistration, result.Error
}
