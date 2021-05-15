package elastic

// import (
// 	"textlooker-backend/util"
// 	"time"
// )

// func AggregateQuery(startDate time.Time, endDate time.Time, sourceID int) {
// 	date := Date{
// 		GTE: util.MakeTimestamp(startDate),
// 		LTE: util.MakeTimestamp(endDate),
// 	}

// 	textQuery := TextQuery{
// 		Query: BoolPart{
// 			Bool: MustPart{
// 				Must: []interface{}{
// 					RangePart{Range: DatePart{Date: date}},
// 					MatchPart{Match: SourcePart{SourceID: sourceID}},
// 					WildcardPart{WildCard: ContentPart{Content: content}},
// 					WildcardPart{WildCard: AuthorPart{Author: author}},
// 				},
// 			},
// 		},
// 	}

// }
