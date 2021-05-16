package elastic

import (
	"time"
)

func NewTextQuery(content string, author string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	conditions := generateBasicConditions(
		makeDateRange(startDate, endDate),
		sourceID, content, author,
	)

	textQuery := TextQuery{
		Query: boolPart{
			Bool: mustPart{
				Must: conditions,
			},
		},
	}

	return textQuery
}
