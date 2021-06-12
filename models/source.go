package models

import (
	"textlooker-backend/database"

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
	}

	result := database.Database.Create(source)
	return source, result.Error
}
