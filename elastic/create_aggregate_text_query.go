package elastic

import (
	"time"
)

func NewAggregateAllQuery(content string, filterItems []FilterItem, startDate time.Time, endDate time.Time, sourceID int, dateAvailableForSource bool) TextQuery {
	query := NewAnalyzedTextQuery(content, filterItems, startDate, endDate, sourceID, true, dateAvailableForSource)
	query.Size = 0
	query = AddGeneralAggregationPart(query, true)
	return query
}

func NewAggregateByOneFieldQuery(content string, filterItems []FilterItem, field string, startDate time.Time, endDate time.Time, sourceID int, dateAvailableForSource bool) TextQuery {
	query := NewAnalyzedTextQuery(content, filterItems, startDate, endDate, sourceID, true, dateAvailableForSource)
	query.Size = 0
	query = AddSingleFieldCompositeAggregationPart(query, field)
	return query
}

func NewDatelessAggregateAllQuery(content string, filterItems []FilterItem, sourceID int, dateAvailableForSource bool) TextQuery {
	query := NewDatelessAnalyzedTextQuery(content, filterItems, sourceID, dateAvailableForSource)
	query.Size = 0
	query = AddGeneralAggregationPart(query, false)
	return query
}
