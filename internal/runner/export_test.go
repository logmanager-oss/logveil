package runner

import (
	"bytes"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/logmanager-oss/logveil/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestLmExport(t *testing.T) {
	tests := []struct {
		name            string
		inputFilename   string
		outputFilename  string
		anonymizingData string
		expectedOutput  string
		expectedProof   []map[string]interface{}
	}{
		{
			name:            "Test LM Export Anonymizer",
			inputFilename:   "../../examples/logs/example_logs.csv",
			anonymizingData: "../../examples/anon_data",
			expectedOutput:  "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.20.0.53\", \"username\":\"ladislav.dosek\", \"organization\":\"Apple\"}\n",
			expectedProof: []map[string]interface{}{
				{"original": "89.239.31.49", "new": "10.20.0.53"},
				{"original": "test.user@test.cz", "new": "ladislav.dosek"},
				{"original": "TESTuser.test.com", "new": "Apple"},
			},
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

			anonymizingData, err := loader.Load(tt.anonymizingData)
			if err != nil {
				t.Fatal(err)
			}
			proofWriter := proof.New(true)
			anonymizer := anonymizer.New(anonymizingData, proofWriter)
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })

			err = AnonymizeLmExport(input, &output, anonymizer)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedOutput, output.String())

			proofWriter.Close()

			actualProof, err := utils.UnpackProofOutput()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expectedProof, actualProof)

			os.Remove("proof.json")
		})
	}
}
