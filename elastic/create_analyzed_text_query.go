package elastic

import (
	"time"
)

func NewAnalyzedTextQuery(content string, author string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	date := makeDate(startDate, endDate)

	conditions := []interface{}{
		rangePart{Range: datePart{Date: date}},
		matchPart{Match: sourcePart{SourceID: sourceID}},
		wildcardPart{WildCard: contentPart{Content: content}},
		wildcardPart{WildCard: authorPart{Author: author}},
	}

	for _, person := range people {
		conditions = append(conditions, matchPart{Match: peoplePart{Person: person}})
	}

	for _, gpeItem := range gpe {
		conditions = append(conditions, matchPart{Match: gpePart{GPE: gpeItem}})
	}

	textQuery := TextQuery{
		Query: boolPart{
			Bool: mustPart{
				Must: conditions,
			},
		},
		AggregateQuery: aggregations{
			AuthorAggregation: aggregation{Terms: field{Field: "author"}},
			PeopleAggregation: aggregation{Terms: field{Field: "people"}},
			GPEAggregation:    aggregation{Terms: field{Field: "gpe"}},
			TokenAggregation:  aggregation{Terms: field{Field: "tokens"}},
			DateAggregation:   aggregation{Terms: field{Field: "date"}},
		},
	}

	return textQuery
}
