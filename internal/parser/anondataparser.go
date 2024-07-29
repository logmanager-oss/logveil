package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// ParseAnonData reads text files from provided directory based on provided field names.
// In other words if file name matches one of the provided field names, it is loaded into the map[fieldName][]anonymizationValues.
// Returned map will be used in anonymization process to match original values with corresponding anonymization values.
func ParseAnonData(anonDataDir string, fieldNames []string) (map[string][]string, error) {
	var anonData = make(map[string][]string)

	for i := range fieldNames {
		if fieldNames[i] == "raw" {
			continue
		}

		filename := filepath.Join(anonDataDir, fieldNames[i])
		_, err := os.Stat(filename)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				slog.Debug(fmt.Sprintf("Anonymizing data not found for field %s. Skipping.\n", fieldNames[i]))
				continue
			}
			return nil, err
		}

		data, err := loadAnonymizingData(filename)
		if err != nil {
			return nil, fmt.Errorf("loading anonymizing data from file %s: %v", filename, err)
		}

		anonData[fieldNames[i]] = data
		slog.Debug(fmt.Sprintf("Loaded anonymizing data for field: %s; values loaded: %d\n", fieldNames[i], len(data)))
	}

	return anonData, nil
}

func loadAnonymizingData(filepath string) ([]string, error) {
	anonDataFile, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var anonData []string
	anonDataFileScanner := bufio.NewScanner(anonDataFile)
	for anonDataFileScanner.Scan() {
		anonData = append(anonData, anonDataFileScanner.Text())
	}

	return anonData, anonDataFile.Close()
}
