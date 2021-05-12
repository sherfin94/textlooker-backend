package elastic

import (
	"bytes"
	"textlooker-backend/util"
	"time"
)

type peoplePart struct {
	Person string `json:"people,omitempty"`
}

type gpePart struct {
	GPE string `json:"gpe,omitempty"`
}

type date struct {
	GTE      string `json:"gte"`
	LTE      string `json:"lte"`
	TimeZone string `json:"time_zone,omitempty"`
}

type datePart struct {
	Date date `json:"date"`
}

type sourcePart struct {
	SourceID int `json:"source_id"`
}

type contentPart struct {
	Content string `json:"content"`
}

type authorPart struct {
	Author string `json:"author"`
}

type wildcardPart struct {
	WildCard interface{} `json:"wildcard"`
}

type matchPart struct {
	Match interface{} `json:"match"`
}

type rangePart struct {
	Range interface{} `json:"range"`
}

type mustPart struct {
	Must []interface{} `json:"must"`
}

type boolPart struct {
	Bool mustPart `json:"bool"`
}

type TextQuery struct {
	Query boolPart `json:"query"`
}

func (textQuery *TextQuery) Buffer() (bytesBuffer bytes.Buffer, err error) {
	return util.StructToBytesBuffer(textQuery)
}

func NewTextQuery(content string, author string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	date := date{
		GTE: util.MakeTimestamp(startDate),
		LTE: util.MakeTimestamp(endDate),
	}

	textQuery := TextQuery{
		Query: boolPart{
			Bool: mustPart{
				Must: []interface{}{
					rangePart{Range: datePart{Date: date}},
					matchPart{Match: sourcePart{SourceID: sourceID}},
					wildcardPart{WildCard: contentPart{Content: content}},
					wildcardPart{WildCard: authorPart{Author: author}},
				},
			},
		},
	}

	// jsonBytes, _ := json.Marshal(textQuery)
	// fmt.Printf("%v\n", string(jsonBytes))

	return textQuery
}

func NewAnalyzedTextQuery(content string, author string, people []string, gpe []string, startDate time.Time, endDate time.Time, sourceID int) TextQuery {
	date := date{
		GTE: util.MakeTimestamp(startDate),
		LTE: util.MakeTimestamp(endDate),
	}

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
	}

	return textQuery
}
