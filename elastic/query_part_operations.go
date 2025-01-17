package elastic

import (
	"bytes"
	"textlooker-backend/util"
	"time"
)

const DateFormat = "2006-01-02T15:04:05-0700"

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func makeDateRange(startDate time.Time, endDate time.Time) dateRange {
	return dateRange{
		GTE: startDate.Format(DateFormat),
		LTE: endDate.Format(DateFormat),
	}
}

func generateBasicConditions(requiredDateRange *dateRange, sourceID int, content string, filterItems []FilterItem, dateAvailableForSource bool) []interface{} {
	var parts []interface{}
	for _, filterItem := range filterItems {

		found := false
		for _, field := range AggregatableFields {
			if field == filterItem.Label {
				found = true
				break
			}
		}
		if found {
			parts = append(parts, matchPart{Match: map[string]interface{}{filterItem.Label: filterItem.Text}})
		}
	}
	if requiredDateRange != nil {
		if dateAvailableForSource {
			parts = append(parts, rangePart{Range: datePart{Date: *requiredDateRange}})
		} else {
			parts = append(parts, rangePart{Range: datePart{Date: *requiredDateRange}})
		}
	}
	parts = append(parts, matchPart{Match: sourcePart{SourceID: sourceID}})
	if content != "" {
		parts = append(parts, matchPart{Match: contentPart{Content: matchQueryPart{Query: content}}})
	}

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
				"size":  50,
			},
		}
	}

	aggregateQuery["date"] = map[string]interface{}{
		"terms": map[string]interface{}{
			"field": "date",
			"size":  50,
		},
	}

	query.AggregateQuery = aggregateQuery
	return query
}

func AddSingleFieldCompositeAggregationPart(query TextQuery, fieldToAggregate string) TextQuery {
	aggregationPart := map[string]interface{}{
		fieldToAggregate: map[string]interface{}{
			"terms": map[string]interface{}{
				"field": fieldToAggregate,
			},
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
