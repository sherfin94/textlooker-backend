package elastic

import (
	"time"
)

func NewAggregateAllQuery(content string, author []string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	query := NewAnalyzedTextQuery(content, author, people, gpe, startDate, endDate, sourceID)
	query.Size = 0
	query = AddGeneralAggregationPart(query, true)
	return query
}

func NewAggregateByOneFieldQuery(content string, author []string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int, field AggregationField) TextQuery {
	query := NewAnalyzedTextQuery(content, author, people, gpe, startDate, endDate, sourceID)
	query.Size = 1000
	query = AddSingleFieldCompositeAggregationPart(query, field)
	return query
}

func NewDatelessAggregateAllQuery(content string, author []string, people []string, gpe []string, tokens []string, sourceID int) TextQuery {
	query := NewDatelessAnalyzedTextQuery(content, author, people, gpe, tokens, sourceID)
	query.Size = 0
	query = AddGeneralAggregationPart(query, false)
	return query
}
