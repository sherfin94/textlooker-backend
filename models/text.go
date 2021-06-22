package models

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync/atomic"
	"textlooker-backend/deployment"
	"textlooker-backend/elastic"
	"textlooker-backend/kafka"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/go-playground/validator/v10"
)

type Text struct {
	ID        string    `json:"-"`
	Content   string    `json:"content" validate:"required"`
	Author    []string  `json:"author"`
	Date      time.Time `json:"date,omitempty"`
	SourceID  int       `json:"source_id" validate:"required"`
	Analyzed  bool      `json:"analyzed"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	DeletedAt time.Time `json:"deleted_at"`
}

func NewText(content string, author []string, date time.Time, sourceID int) (text Text, err error) {
	text = Text{
		Content:   content,
		Author:    author,
		Date:      date,
		SourceID:  sourceID,
		Analyzed:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validator := validator.New()
	if err = validator.Struct(text); err != nil {
		return text, err
	}

	if text.ID, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, ""); err != nil {
		return text, err
	}

	return text, nil
}

func BulkSaveText(textSet []Text) (int, error) {
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"),
		Client:        elastic.Client,
		FlushInterval: 30 * time.Second,
	})

	if err != nil {
		return 0, err
	}

	var countSuccessful uint64 = 0

	for _, text := range textSet {
		data, err := json.Marshal(text)
		if err != nil {
			log.Printf("Cannot encode article : %s", err)
			return 0, err
		}

		err = bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Printf("Unexpected error: %s", err)
			return int(countSuccessful), err
		}
	}

	if err := bulkIndexer.Close(context.Background()); err != nil {
		log.Printf("Unexpected error: %s", err)
		return int(countSuccessful), err
	}

	return int(countSuccessful), err
}

func NewTextWithoutDate(content string, author []string, sourceID int) (text Text, err error) {
	text = Text{
		Content:   content,
		Author:    author,
		SourceID:  sourceID,
		Analyzed:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validator := validator.New()
	if err = validator.Struct(text); err != nil {
		return text, err
	}

	if text.ID, err = elastic.Save(deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT"), text, ""); err != nil {
		return text, err
	}

	return text, nil
}

func GetTexts(content string, author []string, dateStart time.Time, dateEnd time.Time, sourceID int) (texts []Text, err error) {
	textQuery := elastic.NewTextQuery(content, author, dateStart, dateEnd, sourceID)
	texts = []Text{}

	if err != nil {
		log.Println(err)
		return texts, err
	}
	if queryResult, err := elastic.Query(textQuery, deployment.GetEnv("ELASTIC_INDEX_FOR_TEXT")); err != nil {
		log.Println(err)
		return texts, err
	} else {
		for _, hit := range queryResult.Hits.Hits {
			texts = append(texts, Text{
				ID:       hit.ID,
				Content:  hit.Source.Content,
				Author:   hit.Source.Author,
				SourceID: hit.Source.SourceID,
				Analyzed: hit.Source.Analyzed,
				// Date:     hit.Source.Date,
			})
		}
	}
	return texts, err
}

func (text *Text) SendToProcessQueue() {
	kafkaText := kafka.Text{
		ID:        text.ID,
		Content:   text.Content,
		Author:    text.Author,
		SourceID:  text.SourceID,
		CreatedAt: text.CreatedAt,
		UpdatedAt: text.UpdatedAt,
		DeletedAt: text.DeletedAt,
	}

	if !text.Date.IsZero() {
		kafkaText.Date = text.Date.Format("2006-01-02T15:04:05-0700")
	}

	*kafka.TextProcessChannel <- kafkaText
}
