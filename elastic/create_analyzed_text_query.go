package elastic

import (
	"time"
)

func NewAnalyzedTextQuery(content string, filterItems []FilterItem, startDate time.Time, endDate time.Time, sourceID int, dateRangeProvided bool, dateAvailableForSource bool) TextQuery {
	var dateRange dateRange
	var conditions []interface{}
	if dateRangeProvided {
		dateRange = makeDateRange(startDate, endDate)
		conditions = generateBasicConditions(
			&dateRange,
			sourceID, content, filterItems, dateAvailableForSource,
		)
	} else {
		conditions = generateBasicConditions(
			nil,
			sourceID, content, filterItems, false,
		)
	}

	textQuery := generateTextQuery(conditions)

	return textQuery
}

func NewDatelessAnalyzedTextQuery(content string, filterItems []FilterItem, sourceID int, dateAvailableForSource bool) TextQuery {
	conditions := generateBasicConditions(
		nil, sourceID, content, filterItems, false,
	)

	textQuery := generateTextQuery(conditions)

	return textQuery
}
