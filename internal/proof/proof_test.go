package proof

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProof_Write(t *testing.T) {
	tests := []struct {
		name                 string
		isProofWriterEnabled bool
		originalValue        string
		maskedValue          string
		expectedOutput       string
	}{
		{
			name:                 "Test case 1: write proof",
			isProofWriterEnabled: true,
			originalValue:        "test",
			maskedValue:          "masked",
			expectedOutput:       "{\"original\":\"test\",\"new\":\"masked\"}\n",
		},
		{
			name:                 "Test case 2: proof writer disabled",
			isProofWriterEnabled: false,
			originalValue:        "test",
			maskedValue:          "masked",
			expectedOutput:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.isProofWriterEnabled)

			p.Write(tt.originalValue, tt.maskedValue)

			p.Close()

			file, err := os.OpenFile("proof.json", os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				t.Fatal(err)
			}

			buf := bytes.NewBuffer(nil)
			_, err = io.Copy(buf, file)
			if err != nil {
				t.Fatal(err)
			}

			file.Close()

			assert.Equal(t, tt.expectedOutput, buf.String())

			os.Remove("proof.json")
		})
	}
}
