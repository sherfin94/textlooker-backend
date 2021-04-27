package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateToken() string {
	b := make([]byte, 10)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
