package models

import (
	"textlooker-backend/database"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type DashboardInsight struct {
	gorm.Model
	DashboardID int `gorm:"not null" validate:"required"`
	InsightID   int `gorm:"not null" validate:"required"`
}

func (dashboardInsight *DashboardInsight) BeforeSave(database *gorm.DB) (err error) {
	dashboardInsightValidator := validator.New()
	err = dashboardInsightValidator.Struct(dashboardInsight)

	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}

func NewDashboardInsight(dashboardID int, insightID int) (*DashboardInsight, error) {
	dashboardInsight := &DashboardInsight{
		DashboardID: dashboardID,
		InsightID:   insightID,
	}

	result := database.Database.Create(dashboardInsight)
	return dashboardInsight, result.Error
}
