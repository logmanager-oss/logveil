package reader

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLmExport(t *testing.T) {
	tests := []struct {
		name           string
		inputFilename  string
		outputFilename string
		expectedOutput map[string]string
	}{
		{
			name:           "Test LM Export Anonymizer",
			inputFilename:  "../../tests/data/lm_export_test_input.csv",
			expectedOutput: map[string]string{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip": "89.239.31.49", "username": "test.user@test.cz", "organization": "TESTuser.test.com", "raw": "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"89.239.31.49\", \"username\":\"test.user@test.cz\", \"organization\":\"TESTuser.test.com\"}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile, err := os.Open(tt.inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			inputReader, err := NewLmExportReader(inputFile)
			if err != nil {
				t.Fatal(err)
			}

			for {
				logLine, err := inputReader.ReadLine()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedOutput, logLine)
			}
		})
	}
}
