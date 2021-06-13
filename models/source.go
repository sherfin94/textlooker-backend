package models

import (
	"errors"
	"textlooker-backend/database"
	"textlooker-backend/token"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Source struct {
	gorm.Model
	Name            string `gorm:"not null" validate:"required"`
	UserID          int    `gorm:"not null"`
	User            *User  `validate:"structonly"`
	DateAvailable   bool   `gorm:"not null"`
	AuthorAvailable bool   `gorm:"not null"`
	ApiToken        string `gorm:"not null,index:api_token,unique" validate:"required,min=10"`
}

func (source *Source) BeforeSave(database *gorm.DB) (err error) {
	sourceValidator := validator.New()
	err = sourceValidator.Struct(source)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func NewSource(name string, user *User, dateAvailable bool, authorAvailable bool) (*Source, error) {
	source := &Source{
		Name:            name,
		User:            user,
		DateAvailable:   dateAvailable,
		AuthorAvailable: authorAvailable,
		ApiToken:        token.GenerateSecureToken(20),
	}

	result := database.Database.Create(source)
	return source, result.Error
}

func GetSourceByID(sourceID int) (source *Source, err error) {
	result := database.Database.Where("id = ?", sourceID).Find(&source)

	if result.Error != nil {
		err = errors.New("source not found")
		return source, err
	}

	return source, err
}

func GetSourceByToken(token string) (source *Source, err error) {
	result := database.Database.Where("api_token = ?", token).Find(&source)

	if result.Error != nil {
		err = errors.New("source not found")
		return source, err
	}

	return source, err
}
