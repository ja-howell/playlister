package videoclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatLength(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PT56M10S", "56:10"},
		{"PT1H35M41S", "1:35:41"},
		{"PT56M", "56:00"},
		{"PT1H10S", "1:00:10"},
		{"PT1H", "1:00:00"},
		{"PT1H5M", "1:05:00"},
		{"PT5M", "5:00"},
		{"PT1S", "0:01"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test-%d", i), func(t *testing.T) {
			assert := assert.New(t)
			got := formatLength(test.input)
			assert.Equal(test.expected, got)
		})
	}
}
