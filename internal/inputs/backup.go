package inputs

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
)

type LmBackup struct {
	Source LmLog `json:"_source"`
}

type LmLog struct {
	Raw string                 `json:"raw"`
	Msg map[string]interface{} `json:"msg"`
}

func AnonymizeLmBackup(input *os.File, output io.Writer, anonymizer *anonymizer.Anonymizer) error {
	gzReader, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)

	for scanner.Scan() {
		line := scanner.Bytes()

		lmBackup := &LmBackup{}
		err = json.Unmarshal(line, &lmBackup)
		if err != nil {
			return fmt.Errorf("unmarshaling log line: %w", err)
		}

		// Convert map[string]interface{} to map[string]string as requred by anonymizer
		logLine := make(map[string]string)
		for key, value := range lmBackup.Source.Msg {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			logLine[strKey] = strValue
		}
		logLine["raw"] = lmBackup.Source.Raw

		anonymizedLogLine := anonymizer.Anonymize(logLine)

		_, err = io.WriteString(output, fmt.Sprintln(anonymizedLogLine))
		if err != nil {
			return fmt.Errorf("writing anonymized data: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	return nil
}
