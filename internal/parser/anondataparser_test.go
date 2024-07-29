package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_Anondataloader(t *testing.T) {
	tests := []struct {
		name             string
		anonDataDir      string
		fieldNames       []string
		expectedAnonData map[string][]string
	}{
		{
			name:        "Test Anondataloader",
			fieldNames:  []string{"msg.organization"},
			anonDataDir: "../../examples/anon_data",
			expectedAnonData: map[string][]string{
				"msg.organization": {
					"Microsoft", "Apple", "H&P", "IBM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonData, err := ParseAnonData(tt.anonDataDir, tt.fieldNames)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.expectedAnonData, anonData)
		})
	}
}
