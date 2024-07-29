package writer

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_Outputwriter(t *testing.T) {
	tests := []struct {
		name           string
		outputFile     string
		expectedOutput string
	}{
		{
			name:           "Test Output Writer",
			outputFile:     "output.txt",
			expectedOutput: "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.10.10.1\", \"username\":\"miloslav.illes\", \"organization\":\"Microsoft\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputwriter := &Output{
				Output: []string{tt.expectedOutput},
			}

			defer os.Remove(tt.outputFile)

			err := outputwriter.Write(tt.outputFile)
			if err != nil {
				t.Fatal(err)
			}

			data, err := os.ReadFile(tt.outputFile)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedOutput, strings.TrimRight(string(data), "\n"))
		})
	}
}
