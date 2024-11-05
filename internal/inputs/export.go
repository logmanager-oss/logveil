package inputs

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
)

func AnonymizeLmExport(input *os.File, output io.Writer, anonymizer *anonymizer.Anonymizer) error {
	csvReader := csv.NewReader(input)

	// First element of the csvReader contains field names
	fieldNames, err := csvReader.Read()
	if err != nil {
		return err
	}

	// Trimming prefix from field names
	for i, fieldName := range fieldNames {
		fieldNames[i] = strings.TrimPrefix(fieldName, "msg.")
	}

	for {
		row, err := csvReader.Read()
		if err != nil {
			break
		}

		logLine := make(map[string]string)
		for i, val := range row {
			logLine[fieldNames[i]] = val
		}

		anonymizedLogLine := anonymizer.Anonymize(logLine)

		_, err = io.WriteString(output, fmt.Sprintln(anonymizedLogLine))
		if err != nil {
			return fmt.Errorf("writing anonymized data: %v", err)
		}
	}

	return nil
}
