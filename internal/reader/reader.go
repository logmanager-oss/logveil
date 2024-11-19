package reader

import (
	"fmt"
	"os"

	"github.com/logmanager-oss/logveil/internal/config"
	file "github.com/logmanager-oss/logveil/internal/files"
)

type InputReader interface {
	ReadLine() (map[string]string, error)
	Close() error
}

func CreateInputReader(config *config.Config, openFiles *file.FilesHandler) (InputReader, error) {
	inputFile, err := os.Open(config.InputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file for reading: %v", err)
	}
	openFiles.Add(inputFile)

	if *config.IsLmExport {
		inputReader, err := NewLmExportReader(inputFile)
		if err != nil {
			return nil, fmt.Errorf("initializin LM Export reader: %v", err)
		}
		return inputReader, nil
	}

	inputReader, err := NewLmBackupReader(inputFile)
	if err != nil {
		return nil, fmt.Errorf("initializin LM Backup reader: %v", err)
	}

	return inputReader, nil
}
