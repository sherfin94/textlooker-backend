package elastic

import (
	"time"
)

func NewAnalyzedTextQuery(content string, author string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	conditions := generateBasicConditions(
		makeDateRange(startDate, endDate),
		sourceID, content, author,
	)

	for _, person := range people {
		conditions = append(conditions, matchPart{Match: peoplePart{Person: person}})
	}

	for _, gpeItem := range gpe {
		conditions = append(conditions, matchPart{Match: gpePart{GPE: gpeItem}})
	}

	aggregations := generateAllAggregationQueryParts()

	textQuery := generateTextQuery(conditions, &aggregations)

	return textQuery
}
