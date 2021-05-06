package models

import (
	"textlooker-backend/elastic"
	"time"

	"github.com/go-playground/validator/v10"
)

type Text struct {
	Content  string    `json:"content" validate:"required"`
	Author   string    `json:"author" validate:"required"`
	Time     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
}

func NewText(content string, author string, time time.Time, sourceID int) (err error) {
	text := Text{Content: content, Author: author, Time: time, SourceID: sourceID}
	validator := validator.New()
	if err = validator.Struct(text); err != nil {
		return err
	} else {
		elastic.Save("text", text)
	}
	return nil
}
