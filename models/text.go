package models

import (
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"time"

	"github.com/go-playground/validator/v10"
)

type Text struct {
	ID       string    `json:"-"`
	Content  string    `json:"content" validate:"required"`
	Author   string    `json:"author" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
}

func NewText(content string, author string, date time.Time, sourceID int) (text Text, err error) {
	text = Text{Content: content, Author: author, Date: date, SourceID: sourceID}
	validator := validator.New()
	if err = validator.Struct(text); err != nil {
		return text, err
	}

	if text.ID, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, ""); err != nil {
		return text, err
	}

	return text, nil
}
