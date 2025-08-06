package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID(data string) string {
	// This function generates a UUID based on the input data.
	prefix := data
	if len(data) > 3 {
		prefix = data[:3]
	}
	rand.Seed(time.Now().UnixNano())
	suffix := fmt.Sprintf("%06d", rand.Intn(1000000))
	uuid := prefix + suffix
	return uuid
}
