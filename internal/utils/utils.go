package utils

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
)

func CloseFile(fs *os.File) {
	err := fs.Close()
	if err != nil {
		slog.Error(err.Error())
	}
}

func UnpackProofOutput() ([]map[string]interface{}, error) {
	outputFile, err := os.OpenFile("proof.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	scanner := bufio.NewScanner(outputFile)
	for scanner.Scan() {
		var unpackedLine map[string]interface{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &unpackedLine)
		if err != nil {
			return nil, err
		}
		output = append(output, unpackedLine)
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return output, nil
}
