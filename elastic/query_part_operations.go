package elastic

import (
	"bytes"
	"textlooker-backend/util"
	"time"
)

type AggregationField uint8

const People, GPE, Tokens, Authors = 1, 2, 3, 4

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func makeDateRange(startDate time.Time, endDate time.Time) dateRange {
	return dateRange{
		GTE: util.MakeTimestamp(startDate),
		LTE: util.MakeTimestamp(endDate),
	}
}

func generateBasicConditions(requiredDateRange dateRange, sourceID int, content string, author string) []interface{} {
	return []interface{}{
		rangePart{Range: datePart{Date: requiredDateRange}},
		matchPart{Match: sourcePart{SourceID: sourceID}},
		wildcardPart{WildCard: contentPart{Content: content}},
		wildcardPart{WildCard: authorPart{Author: author}},
	}
}

func generateTextQuery(conditions []interface{}) TextQuery {
	textQuery := TextQuery{
		Query: boolPart{
			Bool: mustPart{
				Must: conditions,
			},
		},
	}
	return textQuery
}

func AddGeneralAggregationPart(query TextQuery) TextQuery {
	query.AggregateQuery = aggregationsQueryPart{
		AuthorAggregation: aggregation{Terms: field{Field: "author"}},
		PeopleAggregation: aggregation{Terms: field{Field: "people"}},
		GPEAggregation:    aggregation{Terms: field{Field: "gpe"}},
		TokenAggregation:  aggregation{Terms: field{Field: "tokens"}},
		DateAggregation:   aggregation{Terms: field{Field: "date"}},
	}
	return query
}

func AddSingleFieldCompositeAggregationPart(query TextQuery, fieldToAggregate AggregationField) TextQuery {
	var aggregationPart interface{}

	switch fieldToAggregate {
	case People:
		aggregationPart = aggregationPersonPart{Person: aggregation{Terms: field{Field: "people"}}}
	case Tokens:
		aggregationPart = aggregationTokenPart{Token: aggregation{Terms: field{Field: "tokens"}}}
	case GPE:
		aggregationPart = aggregationGPEPart{GPE: aggregation{Terms: field{Field: "gpe"}}}
	case Authors:
		aggregationPart = aggregationAuthorPart{Author: aggregation{Terms: field{Field: "author"}}}
	}

	query.AggregateQuery = compositeAggregationQueryPart{
		Sources: aggregationsQuerySourcePart{
			Aggregations: []interface{}{
				aggregationDatePart{Date: aggregation{Terms: field{Field: "date"}}},
				aggregationPart,
			},
		},
	}

	return query
}
