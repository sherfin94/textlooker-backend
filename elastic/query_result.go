package elastic

import (
	"encoding/json"
	"io"
	"log"
	"time"
)

type QueryResult struct {
	Took     int        `json:"took"`
	TimedOut bool       `json:"timed_out"`
	Shards   shardsPart `json:"_shards"`
	Hits     hitsPart   `json:"hits"`
}

type innerHitsPart struct {
	Index  string  `json:"_index,omitempty"`
	Type   string  `json:"_type,omitempty"`
	ID     string  `json:"_id"`
	Score  float32 `json:"_score,omitempty"`
	Source Text    `json:"_source"`
}

type hitsPart struct {
	Total    totalPart       `json:"total"`
	MaxScore float32         `json:"max_score"`
	Hits     []innerHitsPart `json:"hits"`
}

type Text struct {
	ID       string    `json:"-"`
	Content  string    `json:"content" validate:"required"`
	Author   string    `json:"author" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
	SourceID int       `json:"source_id" validate:"required"`
	Analyzed bool      `json:"analyzed"`
}

type shardsPart struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type totalPart struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

func ParseResult(body io.ReadCloser) (queryResult QueryResult, err error) {
	err = json.NewDecoder(body).Decode(&queryResult)
	if err != nil {
		log.Fatal(err)
		return queryResult, err
	}

	return queryResult, err
}

// {
// 	"took": 5,
// 	"timed_out": false,
// 	"_shards": {
// 			"total": 1,
// 			"successful": 1,
// 			"skipped": 0,
// 			"failed": 0
// 	},
// 	"hits": {
// 			"total": {
// 					"value": 1,
// 					"relation": "eq"
// 			},
// 			"max_score": 4.0,
// 			"hits": [
// 					{
// 							"_index": "test",
// 							"_type": "_doc",
// 							"_id": "vMQHW3kBUlFJV57c8hsa",
// 							"_score": 4.0,
// 							"_source": {
// 									"content": "Pinarayi Vijayan is the chief minister of Kerala.",
// 									"author": "Kyle Kulinski",
// 									"date": "2021-05-11T16:15:17.97696539+05:30",
// 									"source_id": 1288,
// 									"analyzed": false
// 							}
// 					}
// 			]
// 	}
// }
