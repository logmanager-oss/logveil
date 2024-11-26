package reader

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

var syntaxError *json.SyntaxError

// LmBackup represents raw and msg fields contained in LM Backup format
type LmLog struct {
	Raw string                 `json:"raw"`
	Msg map[string]interface{} `json:"msg"`
}

// LmExportReader represents a reader for LM Backup filetype, which should be a gzip
type LmBackupReader struct {
	backupReader *bufio.Scanner
	file         *os.File
}

func NewLmBackupReader(input *os.File, maxCapacity int) (*LmBackupReader, error) {
	gzReader, err := gzip.NewReader(input)
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %w", err)
	}

	backupReader := bufio.NewScanner(gzReader)
	buf := make([]byte, maxCapacity)
	backupReader.Buffer(buf, maxCapacity)

	return &LmBackupReader{
		backupReader: backupReader,
	}, nil
}

// ReadLine returns a single log line from LM Backup file. Log line is formatted into map[string]string as expected by Anonymizer.
func (r *LmBackupReader) ReadLine() (map[string]string, error) {
	if !r.backupReader.Scan() {
		err := r.backupReader.Err()
		if err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	line := r.backupReader.Bytes()

	lmLog := &LmLog{}
	err := json.Unmarshal(line, &lmLog)
	if err != nil {
		if errors.As(err, &syntaxError) {
			return nil, fmt.Errorf("Malformed lm backup file: %v", err)
		}
		return nil, err
	}

	if lmLog.Raw == "" {
		return nil, fmt.Errorf("Malformed lm backup file - raw field cannot be empty")
	}

	// Convert map[string]interface{} to map[string]string as requred by anonymizer
	logLine := make(map[string]string)
	for key, value := range lmLog.Msg {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		logLine[strKey] = strValue
	}
	logLine["raw"] = lmLog.Raw

	return logLine, nil
}

func (r *LmBackupReader) Close() error {
	return r.file.Close()
}
