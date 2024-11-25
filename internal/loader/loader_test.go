package loader

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_LoadCustomReplacementMap(t *testing.T) {
	tests := []struct {
		name                     string
		customReplacementMapPath string
		expectedMapping          map[string]string
	}{
		{
			name:                     "Test Loading Custom Anonymization Mapping",
			customReplacementMapPath: "../../tests/data/custom_mappings.txt",
			expectedMapping:          map[string]string{"replace_this": "with_that", "test123": "test1234", "test_custom_replacement": "test_custom_replacement123"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping, err := LoadCustomReplacementMap(tt.customReplacementMapPath)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedMapping, mapping)
		})
	}
}

func TestAnonimizer_LoadAnonymizationData(t *testing.T) {
	tests := []struct {
		name           string
		anonDataDir    string
		expectedFields []string
	}{
		{
			name:           "Test Anonymization Data Loading",
			anonDataDir:    "../../tests/data/anonymization_data",
			expectedFields: []string{"dst_iface", "dst_ip", "ip", "name", "organization", "src_iface", "src_ip", "username"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonData, err := LoadAnonymizationData(tt.anonDataDir)
			if err != nil {
				t.Fatal(err)
			}

			for field, value := range anonData {
				assert.Contains(t, tt.expectedFields, field)
				assert.Equal(t, readLines(t, filepath.Join(tt.anonDataDir, field)), value)
			}
		})
	}
}

func readLines(t *testing.T, path string) []string {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		t.Fatal(err)
	}

	return lines
}
