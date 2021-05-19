package models

import (
	"textlooker-backend/database"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string           `gorm:"not null;unique" validate:"required,email"`
	EncryptedPassword  string           `gorm:"not null" validate:"required"`
	UserRegistrationID int              `gorm:"not null"`
	UserRegistration   UserRegistration `gorm:"not null" validate:"structonly"`
}

func (user *User) BeforeSave(database *gorm.DB) (err error) {

	userValidator := validator.New()
	err = userValidator.Struct(user)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return err
}

func NewUser(email string, userRegistration UserRegistration) (*User, error) {
	user := &User{
		Email:             email,
		EncryptedPassword: userRegistration.EncryptedPassword,
		UserRegistration:  userRegistration,
	}

	result := database.Database.Create(user)
	return user, result.Error
}
