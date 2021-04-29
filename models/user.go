package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string           `gorm:"not null;unique" validate:"required,email"`
	EncryptedPassword  string           `gorm:"not null" validate:"required,min=8,max=20"`
	UserRegistrationID int              `gorm:"not null;unique"`
	UserRegistration   UserRegistration `gorm:"not null;foreignkey:UserRegistrationID;unique"`
}

func (user *User) BeforeSave(database *gorm.DB) (err error) {

	userValidator := validator.New()
	err = userValidator.Struct(user)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	hashedPassword, hashingError := bcrypt.GenerateFromPassword([]byte(user.EncryptedPassword), 10)
	if hashingError != nil {
		return hashingError
	} else {
		user.EncryptedPassword = string(hashedPassword)
		err = nil
	}

	return err
}

func NewUser(email string, password string, userRegistration UserRegistration) (*User, error) {
	user := &User{
		Email:             email,
		EncryptedPassword: password,
		UserRegistration:  userRegistration,
	}

	result := Database.Create(user)
	return user, result.Error
}
