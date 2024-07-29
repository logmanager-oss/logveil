package parser

import (
	"encoding/csv"
	"log/slog"
	"os"
)

// ParseCSV takes a CSV file containing logs and transforms it into a list of maps, where each map entry represents a log line.
// Such format is required to be able to modify log data (replace original values with anonymous values).
// It is also returning names of the CSV columns. Names of the columns (field names) are needed to grab corresponding anonymization data.
func ParseCSV(filename string) ([]string, []map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer func(fs *os.File) {
		if err := fs.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(file)

	csvReader := csv.NewReader(file)

	// First element of the csvReader contains field names
	fieldNames, err := csvReader.Read()
	if err != nil {
		return nil, nil, err
	}

	var csvData []map[string]string
	for {
		row, err := csvReader.Read()
		if err != nil {
			break
		}

		m := make(map[string]string)
		for i, val := range row {
			m[fieldNames[i]] = val
		}
		csvData = append(csvData, m)
	}

	return fieldNames, csvData, nil
}
