package elastic

import (
	"time"
)

func NewTextQuery(content string, author string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	dateRange := makeDateRange(startDate, endDate)

	textQuery := TextQuery{
		Query: boolPart{
			Bool: mustPart{
				Must: []interface{}{
					rangePart{Range: datePart{Date: dateRange}},
					matchPart{Match: sourcePart{SourceID: sourceID}},
					wildcardPart{WildCard: contentPart{Content: content}},
					wildcardPart{WildCard: authorPart{Author: author}},
				},
			},
		},
	}

	return textQuery
}
