package models

import (
	"textlooker-backend/database"
	"textlooker-backend/util"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistration struct {
	gorm.Model
	Email             string `gorm:"not null;unique" validate:"required,email"`
	VerificationToken string `gorm:"not null" validate:"required"`
	EncryptedPassword string `gorm:"not null" validate:"required"`
}

func (userRegistration *UserRegistration) BeforeSave(database *gorm.DB) (err error) {
	userRegistrationValidator := validator.New()
	err = userRegistrationValidator.Struct(userRegistration)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	hashedPassword, hashingError := bcrypt.GenerateFromPassword([]byte(userRegistration.EncryptedPassword), 10)
	if hashingError != nil {
		return hashingError
	} else {
		userRegistration.EncryptedPassword = string(hashedPassword)
		err = nil
	}

	return nil
}

func NewUserRegistration(email string, password string) (*UserRegistration, error) {
	var userRegistration *UserRegistration

	otp, err := util.GenerateToken(5)
	if err != nil {
		return userRegistration, err
	}

	userRegistration = &UserRegistration{
		Email:             email,
		VerificationToken: otp,
		EncryptedPassword: password,
	}

	result := database.Database.Create(userRegistration)

	return userRegistration, result.Error
}
