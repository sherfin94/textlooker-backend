package elastic

import (
	"bytes"
	"textlooker-backend/util"
	"time"
)

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func makeDateRange(startDate time.Time, endDate time.Time) dateRange {
	return dateRange{
		GTE: util.MakeTimestamp(startDate),
		LTE: util.MakeTimestamp(endDate),
	}
}

func generateBasicConditions(requiredDateRange *dateRange, sourceID int, content string, filterItems []FilterItem) []interface{} {
	var parts []interface{}
	for _, filterItem := range filterItems {
		parts = append(parts, matchPart{Match: map[string]interface{}{filterItem.Label: filterItem.Text}})
	}
	if requiredDateRange != nil {
		parts = append(parts, rangePart{Range: datePart{Date: *requiredDateRange}})
	}
	parts = append(parts, matchPart{Match: sourcePart{SourceID: sourceID}})
	parts = append(parts, wildcardPart{WildCard: contentPart{Content: content}})

	return parts
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

func AddGeneralAggregationPart(query TextQuery, includeDate bool) TextQuery {
	aggregateQuery := map[string]interface{}{}

	for _, field := range AggregatableFields {
		aggregateQuery[field] = map[string]interface{}{
			"terms": map[string]interface{}{
				"field": field,
			},
		}
	}

	aggregateQuery["date"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "date",
		},
	}

	query.AggregateQuery = aggregateQuery
	return query
}

func AddSingleFieldCompositeAggregationPart(query TextQuery, fieldToAggregate string) TextQuery {
	aggregationPart := map[string]interface{}{
		"terms": map[string]interface{}{
			"field": fieldToAggregate,
		},
	}

	query.AggregateQuery =
		customBucketNamePartForCompositeQuery{
			compositeAggregationQueryPart{
				Sources: aggregationsQuerySourcePart{
					Size: 100,
					Aggregations: []interface{}{
						aggregationDatePart{Date: dateHistogramAggregation{Terms: field{Field: "date", Interval: "1d"}}},
						aggregationPart,
					},
				},
			},
		}

	return query
}
