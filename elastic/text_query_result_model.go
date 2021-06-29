package elastic

import (
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02T15:04:05-0700"
	s := string(b)
	s = s[1 : len(s)-1]
	result, err := time.Parse(layout, s)
	*t = Time{result}
	return err
}

type Config struct {
	T Time
}

type QueryResult struct {
	Took             int         `json:"took"`
	TimedOut         bool        `json:"timed_out"`
	Shards           shardsPart  `json:"_shards"`
	Hits             hitsPart    `json:"hits"`
	AggregationsPart interface{} `json:"aggregations,omitempty"`
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
	ID       string   `json:"-"`
	Content  string   `json:"content" validate:"required"`
	Author   []string `json:"author" validate:"required"`
	Date     Time     `json:"date,omitempty" validate:"required"`
	SourceID int      `json:"source_id" validate:"required"`
	Analyzed bool     `json:"analyzed,omitempty"`
	// People   []string `json:"people,omitempty"`
	// GPE      []string `json:"gpe,omitempty"`
	// Tokens   []string `json:"tokens,omitempty"`
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

type AggregationResultPart struct {
	Buckets []count `json:"buckets,omitempty"`
}

type count struct {
	Key   interface{} `json:"key"`
	Value int         `json:"doc_count"`
	Date  float64     `json:"date,omitempty"`
}

type perDateAggregationResultPart struct {
	Buckets []perDateBucket `json:"buckets,omitempty"`
}

type perDateBucket struct {
	Key   perDateCount `json:"key"`
	Count int          `json:"doc_count"`
}

type perDateCount struct {
	Date       float64     `json:"date,omitempty"`
	FieldValue interface{} `json:"field_value"`
}
