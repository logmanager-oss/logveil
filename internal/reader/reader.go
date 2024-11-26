package reader

import (
	"fmt"
	"os"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/handlers"
)

type InputReader interface {
	ReadLine() (map[string]string, error)
	Close() error
}

func CreateInputReader(config *config.Config, filesHandler *handlers.Files) (InputReader, error) {
	inputFile, err := os.Open(config.InputPath)
	if err != nil {
		return nil, fmt.Errorf("opening input file for reading: %v", err)
	}
	filesHandler.Add(inputFile)

	if config.IsLmExport {
		inputReader, err := NewLmExportReader(inputFile)
		if err != nil {
			return nil, fmt.Errorf("initializin LM Export reader: %v", err)
		}
		return inputReader, nil
	}

	inputReader, err := NewLmBackupReader(inputFile, config.ReaderMaxCapacity)
	if err != nil {
		return nil, fmt.Errorf("initializin LM Backup reader: %v", err)
	}

	return inputReader, nil
}
