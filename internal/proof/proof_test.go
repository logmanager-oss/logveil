package proof

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/files"
	"github.com/stretchr/testify/assert"
)

func TestProof_Write(t *testing.T) {
	tests := []struct {
		name           string
		isProofWriter  bool
		replacementMap map[string]string
		expectedOutput string
	}{
		{
			name:          "Test case 1: write proof",
			isProofWriter: true,
			replacementMap: map[string]string{
				"test": "masked",
			},
			expectedOutput: "{\"original\":\"test\",\"new\":\"masked\"}\n",
		},
		{
			name:          "Test case 2: proof writer disabled",
			isProofWriter: false,
			replacementMap: map[string]string{
				"test": "masked",
			},
			expectedOutput: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filesHandler := &files.FilesHandler{}
			defer filesHandler.Close()

			p, err := CreateProofWriter(&config.Config{IsProofWriter: tt.isProofWriter}, filesHandler)
			if err != nil {
				t.Fatal(err)
			}

			p.Write(tt.replacementMap)
			p.Flush()

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
