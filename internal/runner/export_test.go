package runner

import (
	"bytes"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/stretchr/testify/assert"
)

func TestLmExport(t *testing.T) {
	tests := []struct {
		name            string
		inputFilename   string
		outputFilename  string
		anonymizingData string
		expectedOutput  string
	}{
		{
			name:            "Test LM Export Anonymizer",
			inputFilename:   "../../examples/logs/example_logs.csv",
			anonymizingData: "../../examples/anon_data",
			expectedOutput:  "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.20.0.53\", \"username\":\"ladislav.dosek\", \"organization\":\"Apple\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.Open(tt.inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			defer input.Close()

			var output bytes.Buffer

			anonData, err := loader.Load(tt.anonymizingData)
			if err != nil {
				t.Fatal(err)
			}
			anonymizer := anonymizer.New(anonData)
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })

			err = AnonymizeLmExport(input, &output, anonymizer)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}
