package util

import (
	"fmt"
	"time"
)

func MakeTimestamp(givenTime time.Time) string {
	timestamp := float64(givenTime.UnixNano()) / float64(time.Millisecond)
	return fmt.Sprintf("%f", timestamp)
}

func ParseTimestamp(timestamp float64) *time.Time {
	timeAsInteger := int64(timestamp) / 1000

	result := time.Unix(timeAsInteger, 0)
	return &result
}
