package bytesUtils

import (
	"bytes"
)

func Contains(scannerBytes []byte, words []string) bool {
	for _, word := range words {
		if bytes.Contains(scannerBytes, []byte(word)) {
			return true
		}
	}
	return false
}
