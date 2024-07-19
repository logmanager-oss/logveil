package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_CSVloader(t *testing.T) {
	tests := []struct {
		name               string
		filename           string
		expectedFieldNames []string
		expectedValues     []map[string]string
	}{
		{
			name:               "Test CSVLoader",
			filename:           "../../examples/logs/example_logs.csv",
			expectedFieldNames: []string{"@timestamp", "raw", "msg.src_ip", "msg.username", "msg.organization"},
			expectedValues: []map[string]string{{
				"@timestamp":       "2024-06-05T14:59:27.000+00:00",
				"msg.organization": "TESTuser.test.com",
				"msg.src_ip":       "89.239.31.49", "msg.username": "test.user@test.cz",
				"raw": "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"89.239.31.49\", \"username\":\"test.user@test.cz\", \"organization\":\"TESTuser.test.com\"}",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldNames, csvData, err := ParseCSV(tt.filename)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expectedFieldNames, fieldNames)
			assert.Equal(t, tt.expectedValues, csvData)
		})
	}
}
