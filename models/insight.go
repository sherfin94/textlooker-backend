package models

import (
	"encoding/json"
	"textlooker-backend/database"
	"textlooker-backend/elastic"
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

type filterItem struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

type filterObject struct {
	FilterItems []filterItem `json:"filter"`
}

type visualizeTextSet struct {
	Texts []string `json:"visualizeTexts"`
}

func (insight *Insight) Aggregation() (aggregation map[string]interface{}, err error) {
	var filter filterObject
	var visualizeTexts visualizeTextSet

	err = json.Unmarshal([]byte(insight.Filter), &filter)
	if err != nil {
		return aggregation, err
	}

	err = json.Unmarshal([]byte(insight.VisualizeTexts), &visualizeTexts)
	if err != nil {
		return aggregation, err
	}

	filterItems := []elastic.FilterItem{}
	for _, item := range filter.FilterItems {
		filterItems = append(filterItems, elastic.FilterItem{Label: item.Label, Text: item.Text})
	}

	if insight.DateRangeAvailable {
		aggregation, err = GetAggregation(
			"*", filterItems,
			insight.StartDate, insight.EndDate,
			insight.SourceID,
		)
	} else {
		aggregation, err = GetDatelessAggregation(
			"*", filterItems,
			insight.SourceID,
		)
	}

	aggregation["visualizeTexts"] = visualizeTexts.Texts

	return aggregation, err
}
