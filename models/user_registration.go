package models

import (
	"textlooker-backend/database"
	"textlooker-backend/util"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserRegistration struct {
	gorm.Model
	Email             string `gorm:"not null;unique" validate:"required,email"`
	VerificationToken string `gorm:"not null" validate:"required"`
}

func (userRegistration *UserRegistration) BeforeSave(database *gorm.DB) (err error) {
	userRegistrationValidator := validator.New()
	err = userRegistrationValidator.Struct(userRegistration)

	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

func NewUserRegistration(email string) (*UserRegistration, error) {
	userRegistration := &UserRegistration{
		Email:             email,
		VerificationToken: util.GenerateToken(),
	}

	result := database.Database.Create(userRegistration)

	return userRegistration, result.Error
}
