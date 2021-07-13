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

func GetTexts(content string, filterItems []elastic.FilterItem, dateStart time.Time, dateEnd time.Time, sourceID int, dateAvailableForSource bool) (texts []Text, err error) {
	textQuery := elastic.NewTextQuery(content, filterItems, dateStart, dateEnd, sourceID, dateAvailableForSource)
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

func SendToProcessQueue(textSet []Text) {
	kafkaTextSet := kafka.TextSet{Set: []kafka.Text{}}

	for _, text := range textSet {
		kafkaText := kafka.Text{
			ID:        text.ID,
			Content:   text.Content,
			Author:    text.Author,
			SourceID:  text.SourceID,
			CreatedAt: text.CreatedAt,
			UpdatedAt: text.UpdatedAt,
			DeletedAt: text.DeletedAt,
			Date:      text.Date.Format("2006-01-02T15:04:05-0700"),
		}
		if text.Date.IsZero() {
			kafkaText.Date = time.Now().Format("2006-01-02T15:04:05-0700")
		}

		kafkaTextSet.Set = append(kafkaTextSet.Set, kafkaText)
	}

	*kafka.TextProcessChannel <- kafkaTextSet
}
