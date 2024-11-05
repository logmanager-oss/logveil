package loader

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_Anondataloader(t *testing.T) {
	tests := []struct {
		name           string
		anonDataDir    string
		expectedFields []string
	}{
		{
			name:           "Test Anondataloader",
			anonDataDir:    "../../examples/anon_data",
			expectedFields: []string{"dst_iface", "dst_ip", "ip", "name", "organization", "src_iface", "src_ip", "username"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonData, err := Load(tt.anonDataDir)
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
