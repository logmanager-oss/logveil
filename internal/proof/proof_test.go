package proof

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestProof_Write(t *testing.T) {
	tests := []struct {
		name           string
		isProofWriter  bool
		originalValue  string
		newValue       string
		expectedOutput string
	}{
		{
			name:           "Test case 1: write proof",
			isProofWriter:  true,
			originalValue:  "test",
			newValue:       "masked",
			expectedOutput: "{\"original\":\"test\",\"new\":\"masked\"}\n",
		},
		{
			name:           "Test case 2: proof writer disabled",
			isProofWriter:  false,
			originalValue:  "test",
			newValue:       "masked",
			expectedOutput: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filesHandler := &handlers.Files{}
			defer filesHandler.Close()

			buffersHandler := &handlers.Buffers{}
			defer buffersHandler.Flush()

			p, err := CreateProofWriter(&config.Config{IsProofWriter: tt.isProofWriter}, filesHandler, buffersHandler)
			if err != nil {
				t.Fatal(err)
			}

			p.Write(tt.originalValue, tt.newValue)
			p.Flush()

			file, err := os.OpenFile(ProofFilename, os.O_RDWR|os.O_CREATE, 0644)
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

			os.Remove(ProofFilename)
		})
	}
}
