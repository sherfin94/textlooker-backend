package elastic

import (
	"time"
)

type QueryResult struct {
	Took             int              `json:"took"`
	TimedOut         bool             `json:"timed_out"`
	Shards           shardsPart       `json:"_shards"`
	Hits             hitsPart         `json:"hits"`
	AggregationsPart aggregationsPart `json:"aggregations,omitempty"`
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
	Analyzed bool      `json:"analyzed,omitempty"`
	People   []string  `json:"people,omitempty"`
	GPE      []string  `json:"gpe,omitempty"`
	Tokens   []string  `json:"tokens,omitempty"`
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

type aggregationsPart struct {
	AuthorAggregation aggregationResultPart `json:"authors,omitempty"`
	PeopleAggregation aggregationResultPart `json:"people,omitempty"`
	GPEAggregation    aggregationResultPart `json:"gpe,omitempty"`
	TokenAggregation  aggregationResultPart `json:"tokens,omitempty"`
	DateAggregation   aggregationResultPart `json:"date,omitempty"`
}

type aggregationResultPart struct {
	Buckets []count `json:"buckets,omitempty"`
}

type count struct {
	Key   interface{} `json:"key"`
	Value int         `json:"doc_count"`
	Date  float64     `json:"date,omitempty"`
}
