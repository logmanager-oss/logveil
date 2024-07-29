package anonymizer

import (
	"testing"

	"github.com/logmanager-oss/logveil/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_AnonymizeData(t *testing.T) {
	tests := []struct {
		name           string
		anonDataDir    string
		inputFile      string
		expectedOutput []string
	}{
		{
			name:           "Test AnonymizeData",
			anonDataDir:    "../../examples/anon_data",
			inputFile:      "../../examples/logs/example_logs.csv",
			expectedOutput: []string{"{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.10.10.1\", \"username\":\"miloslav.illes\", \"organization\":\"Microsoft\"}"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldNames, csvData, err := parser.ParseCSV(tt.inputFile)
			if err != nil {
				t.Fatalf("reading input file %s: %v", tt.inputFile, err)
			}

			anonData, err := parser.ParseAnonData(tt.anonDataDir, fieldNames)
			if err != nil {
				t.Fatalf("loading anonymizing data from dir %s: %v", tt.anonDataDir, err)
			}

			anonymizer := New(csvData, anonData)
			anonymizer.setRandFunc(func(int) int { return 0 })
			output := anonymizer.anonymize()

			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
