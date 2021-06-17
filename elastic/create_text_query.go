package elastic

import (
	"time"
)

func NewTextQuery(content string, author []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	dateRange := makeDateRange(startDate, endDate)
	conditions := generateBasicConditions(
		&dateRange,
		sourceID, content, author,
	)

	textQuery := generateTextQuery(conditions)

	return textQuery
}
