package elastic

type aggregations struct {
	AuthorAggregation aggregation `json:"authors,omitempty"`
	PeopleAggregation aggregation `json:"people,omitempty"`
	GPEAggregation    aggregation `json:"gpe,omitempty"`
	TokenAggregation  aggregation `json:"tokens,omitempty"`
	DateAggregation   aggregation `json:"date,omitempty"`
}

type aggregation struct {
	Terms field `json:"terms"`
}

type field struct {
	Field string `json:"field"`
}

type peoplePart struct {
	Person string `json:"people,omitempty"`
}

type gpePart struct {
	GPE string `json:"gpe,omitempty"`
}

type dateRange struct {
	GTE      string `json:"gte"`
	LTE      string `json:"lte"`
	TimeZone string `json:"time_zone,omitempty"`
}

type datePart struct {
	Date dateRange `json:"date"`
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
	Query          boolPart     `json:"query"`
	AggregateQuery aggregations `json:"aggs,omitempty"`
}
