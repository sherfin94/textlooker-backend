package util

import (
	"fmt"
	"strconv"
	"time"
)

func MakeTimestamp(givenTime time.Time) string {
	timestamp := givenTime.UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", timestamp)
}

func ParseTimestamp(timestamp string) *time.Time {
	timeAsInteger, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}

	result := time.Unix(timeAsInteger, 0)
	return &result
}
