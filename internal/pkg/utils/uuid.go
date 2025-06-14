package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a random UUID v4
func GenerateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Sprintf("error-generating-uuid-%d", time.Now().UnixNano())
	}

	return id.String()
}
