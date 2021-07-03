package models

import (
	"textlooker-backend/database"
	"textlooker-backend/token"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Dashboard struct {
	gorm.Model
	Title    string `gorm:"not null" validate:"required"`
	SourceID int    `gorm:"not null" validate:"required"`
	Token    string `gorm:"not null" validate:"required"`
}

func (dashboard *Dashboard) BeforeSave(database *gorm.DB) (err error) {
	dashboardValidator := validator.New()
	err = dashboardValidator.Struct(dashboard)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func NewDashboard(title string, sourceID int) (*Dashboard, error) {
	dashboard := &Dashboard{
		Title:    title,
		SourceID: sourceID,
		Token:    token.GenerateSecureToken(20),
	}

	result := database.Database.Create(dashboard)
	return dashboard, result.Error
}
