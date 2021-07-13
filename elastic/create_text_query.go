package elastic

import (
	"time"
)

func NewTextQuery(content string, filterItems []FilterItem, startDate time.Time, endDate time.Time, sourceID int, dateAvailableForSource bool) TextQuery {
	dateRange := makeDateRange(startDate, endDate)
	conditions := generateBasicConditions(
		&dateRange,
		sourceID, content, filterItems, dateAvailableForSource,
	)

	textQuery := generateTextQuery(conditions)

	return textQuery
}
