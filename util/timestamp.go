package util

import (
	"fmt"
	"time"
)

func MakeTimestamp(givenTime time.Time) string {
	timestamp := givenTime.UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", timestamp)
}
