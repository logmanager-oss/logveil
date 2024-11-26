package loader

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func LoadCustomReplacementMap(path string) (map[string]string, error) {
	customReplacementMap := make(map[string]string)

	if path != "" {
		file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			values := strings.Split(line, ":")
			if len(values) == 1 {
				slog.Error("wrong custom mapping: %s", "error", line)
			}

			originalValue := values[0]
			newValue := values[1]

			customReplacementMap[originalValue] = newValue
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading custom anonymization mapping: %w", err)
		}
	}

	return customReplacementMap, nil
}

// LoadAnonymizationData() loads anonymization data from given directory and returns it in a map format of: [filename][]values. Anonymization data is needed for the purposes of masking original values.
func LoadAnonymizationData(path string) (map[string][]string, error) {
	anonymizationData := make(map[string][]string)

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := loadFromFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("loading anonymizing data from file %s: %v", file.Name(), err)
		}

		anonymizationData[file.Name()] = data
		slog.Debug(fmt.Sprintf("Loaded anonymizing data for field: %s; values loaded: %d\n", file.Name(), len(data)))
	}

	return anonymizationData, nil
}

func loadFromFile(filepath string) ([]string, error) {
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
