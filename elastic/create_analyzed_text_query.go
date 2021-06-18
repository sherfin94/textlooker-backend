package elastic

import (
	"time"
)

func NewAnalyzedTextQuery(content string, author []string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	dateRange := makeDateRange(startDate, endDate)
	conditions := generateBasicConditions(
		&dateRange,
		sourceID, content, author,
	)

	for _, person := range people {
		conditions = append(conditions, matchPart{Match: peoplePart{Person: person}})
	}

	for _, gpeItem := range gpe {
		conditions = append(conditions, matchPart{Match: gpePart{GPE: gpeItem}})
	}

	textQuery := generateTextQuery(conditions)

	return textQuery
}

func NewDatelessAnalyzedTextQuery(content string, author []string, people []string, gpe []string, tokens []string, sourceID int) TextQuery {
	conditions := generateBasicConditions(
		nil, sourceID, content, author,
	)

	for _, person := range people {
		conditions = append(conditions, matchPart{Match: peoplePart{Person: person}})
	}

	for _, gpeItem := range gpe {
		conditions = append(conditions, matchPart{Match: gpePart{GPE: gpeItem}})
	}

	for _, token := range tokens {
		conditions = append(conditions, matchPart{Match: tokenPart{Token: token}})
	}

	textQuery := generateTextQuery(conditions)

	return textQuery
}
