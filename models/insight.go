package models

import (
	"textlooker-backend/database"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Insight struct {
	gorm.Model
	Title              string    `gorm:"not null,unique" validate:"required"`
	SourceID           int       `gorm:"not null"`
	Filter             string    `gorm:"not null" validate:"required"`
	LookForHandle      string    `gorm:"not null" validate:"required"`
	VisualizeTexts     string    `gorm:"not null" validate:"required"`
	VisualizationType  string    `gorm:"not null" validate:"required"`
	StartDate          time.Time `gorm:"not null" validate:"required"`
	EndDate            time.Time `gorm:"not null" validate:"required"`
	DateRangeAvailable bool      `gorm:"not null"`
}

func (insight *Insight) BeforeSave(database *gorm.DB) (err error) {
	sourceValidator := validator.New()
	err = sourceValidator.Struct(insight)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func NewInsight(title string, filter string, lookForHandle string, visualizeTexts string, startDate time.Time, endDate time.Time, visualizationType string, dateRangeAvailable bool, sourceID int) (*Insight, error) {
	insight := &Insight{
		Title:              title,
		LookForHandle:      lookForHandle,
		VisualizeTexts:     visualizeTexts,
		Filter:             filter,
		SourceID:           sourceID,
		StartDate:          startDate,
		EndDate:            endDate,
		DateRangeAvailable: dateRangeAvailable,
		VisualizationType:  visualizationType,
	}

	result := database.Database.Create(insight)
	return insight, result.Error
}
