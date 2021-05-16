package elastic

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"textlooker-backend/util"
	"time"
)

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func ParseResult(body io.ReadCloser) (queryResult QueryResult, err error) {
	err = json.NewDecoder(body).Decode(&queryResult)
	if err != nil {
		log.Fatal(err)
		return queryResult, err
	}
	return queryResult, err
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

func generateAllAggregationQueryParts() aggregations {
	return aggregations{
		AuthorAggregation: aggregation{Terms: field{Field: "author"}},
		PeopleAggregation: aggregation{Terms: field{Field: "people"}},
		GPEAggregation:    aggregation{Terms: field{Field: "gpe"}},
		TokenAggregation:  aggregation{Terms: field{Field: "tokens"}},
		DateAggregation:   aggregation{Terms: field{Field: "date"}},
	}
}
