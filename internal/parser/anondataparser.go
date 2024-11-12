package parser

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func LoadAnonData(anonDataDir string) (map[string][]string, error) {
	var anonData = make(map[string][]string)

	files, err := os.ReadDir(anonDataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.Contains(file.Name(), "ip") {
			slog.Debug("skipping file with IP addresses - generator handles IP replacements")
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
