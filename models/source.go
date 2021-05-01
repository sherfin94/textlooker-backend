package models

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Source struct {
	gorm.Model
	Name   string `gorm:"not null;unique" validate:"required"`
	UserID int    `gorm:"not null;unique"`
	User   *User  `validate:"structonly"`
}

func (source *Source) BeforeSave(database *gorm.DB) (err error) {
	sourceValidator := validator.New()
	err = sourceValidator.Struct(source)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func NewSource(name string, user *User) (*Source, error) {
	source := &Source{
		Name: name,
		User: user,
	}

	result := Database.Create(source)
	return source, result.Error
}
