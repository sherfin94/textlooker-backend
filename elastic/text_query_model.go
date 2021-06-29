package elastic

import "encoding/json"

type aggregationsQueryPart struct {
	AuthorAggregation aggregation `json:"authors,omitempty"`
	PeopleAggregation aggregation `json:"people,omitempty"`
	GPEAggregation    aggregation `json:"gpe,omitempty"`
	TokenAggregation  aggregation `json:"tokens,omitempty"`
	DateAggregation   aggregation `json:"date,omitempty"`
}

type compositeAggregationQueryPart struct {
	Sources aggregationsQuerySourcePart `json:"composite"`
}

type customBucketNamePartForCompositeQuery struct {
	Composite compositeAggregationQueryPart `json:"per_date"`
}

type aggregationsQuerySourcePart struct {
	Aggregations []interface{} `json:"sources"`
	Size         int           `json:"size,omitempty"`
}

type aggregation struct {
	Terms field `json:"terms"`
}

type dateHistogramAggregation struct {
	Terms field `json:"date_histogram"`
}

type field struct {
	Field    string `json:"field"`
	Interval string `json:"interval,omitempty"`
}

type peoplePart struct {
	Person string `json:"people,omitempty"`
}

type gpePart struct {
	GPE string `json:"gpe,omitempty"`
}

type tokenPart struct {
	Token string `json:"tokens,omitempty"`
}

type dateRange struct {
	GTE      string `json:"gte"`
	LTE      string `json:"lte"`
	TimeZone string `json:"time_zone,omitempty"`
}

type datePart struct {
	Date dateRange `json:"date"`
}

type aggregationGenericFieldPart struct {
	Field aggregation `json:"field_value"`
}

type aggregationPersonPart struct {
	Person aggregation `json:"person"`
}

type aggregationGPEPart struct {
	GPE aggregation `json:"gpe"`
}

type aggregationTokenPart struct {
	Token aggregation `json:"token"`
}

type aggregationAuthorPart struct {
	Author aggregation `json:"author"`
}

type aggregationDatePart struct {
	Date dateHistogramAggregation `json:"date"`
}

type sourcePart struct {
	SourceID int `json:"source_id"`
}

type contentPart struct {
	Content interface{} `json:"content"`
}

type matchQueryPart struct {
	Query string `json:"query"`
}

type authorPart struct {
	Author string `json:"author"`
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
	Query          boolPart    `json:"query"`
	Size           int         `json:"size"`
	From           int         `json:"from"`
	AggregateQuery interface{} `json:"aggs,omitempty"`
}

func (query *TextQuery) RequestString() (requestString string) {
	jsonBytes, _ := json.Marshal(query)
	return string(jsonBytes)
}
