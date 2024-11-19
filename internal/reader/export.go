package reader

import (
	"encoding/csv"
	"os"
	"strings"
)

// LmExportReader represents a reader for LM Export filetype, which should be a CSV
type LmExportReader struct {
	exportReader *csv.Reader
	fieldNames   []string
	file         *os.File
}

func NewLmExportReader(input *os.File) (*LmExportReader, error) {
	csvReader := csv.NewReader(input)

	// First element of the csvReader contains field names
	fieldNames, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// Trimming prefix from field names
	for i, fieldName := range fieldNames {
		fieldNames[i] = strings.TrimPrefix(fieldName, "msg.")
	}

	return &LmExportReader{
		exportReader: csvReader,
		fieldNames:   fieldNames,
	}, nil
}

// ReadLine returns a single log line from LM Export file. Log line is formatted into map[string]string as expected by Anonymizer.
func (r *LmExportReader) ReadLine() (map[string]string, error) {
	row, err := r.exportReader.Read()
	if err != nil {
		return nil, err
	}

	logLine := make(map[string]string)
	for i, val := range row {
		logLine[r.fieldNames[i]] = val
	}

	return logLine, nil
}

func (r *LmExportReader) Close() error {
	return r.file.Close()
}
