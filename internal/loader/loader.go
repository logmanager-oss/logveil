package loader

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

// Load() loads anonymization data from given directory and returns it in a map format of: [filename][]values. Anonymization data is needed for the purposes of masking original values.
func Load(anonDataDir string) (map[string][]string, error) {
	var anonData = make(map[string][]string)

	files, err := os.ReadDir(anonDataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := loadAnonymizingData(filepath.Join(anonDataDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("loading anonymizing data from file %s: %v", file.Name(), err)
		}

		anonData[file.Name()] = data
		slog.Debug(fmt.Sprintf("Loaded anonymizing data for field: %s; values loaded: %d\n", file.Name(), len(data)))
	}

	return anonData, nil
}

func loadAnonymizingData(filepath string) ([]string, error) {
	anonDataFile, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var anonData []string
	scanner := bufio.NewScanner(anonDataFile)
	for scanner.Scan() {
		anonData = append(anonData, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading anon data: %w", err)
	}

	return anonData, anonDataFile.Close()
}
