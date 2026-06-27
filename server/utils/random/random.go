package random

import (
	"crypto/rand"
	"fmt"
)

func RandBytes(length uint) ([]byte, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return b, fmt.Errorf("failed generate: %v", err)
	}
	return b, nil
}
