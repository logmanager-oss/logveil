package anonymizer

import (
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/logmanager-oss/logveil/internal/generator"
	"github.com/logmanager-oss/logveil/internal/lookup"
	"github.com/logmanager-oss/logveil/internal/parser"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_AnonymizeData(t *testing.T) {
	tests := []struct {
		name                 string
		isProofWriterEnabled bool
		anonDataDir          string
		input                map[string]string
		expectedOutput       string
	}{
		{
			name:                 "Test AnonymizeData",
			isProofWriterEnabled: false,
			anonDataDir:          "../../examples/anon_data",
			input:                map[string]string{"@timestamp": "2024-06-05T14:59:27.000+00:00", "src_ip": "10.10.10.1", "username": "miloslav.illes", "organization": "Microsoft", "raw": "2024-06-05T14:59:27.000+00:00, 10.10.10.1, miloslav.illes, Microsoft"},
			expectedOutput:       "2024-06-05T14:59:27.000+00:00, 33.15.199.187, ladislav.dosek, Apple",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonData, err := parser.LoadAnonData(tt.anonDataDir)
			if err != nil {
				t.Fatalf("loading anonymizing data from dir %s: %v", tt.anonDataDir, err)
			}

			proofWriter := proof.New(tt.isProofWriterEnabled)
			lookup := lookup.New()
			generator := &generator.Generator{}
			anonymizer := New(anonData, proofWriter, lookup, generator)
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })
			faker.SetRandomSource(rand.NewSource(1))

			output := anonymizer.Anonymize(tt.input)

			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
