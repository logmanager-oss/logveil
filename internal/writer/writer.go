package writer

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/config"
	file "github.com/logmanager-oss/logveil/internal/files"
)

func CreateOutputWriter(configFile *config.Config, openFiles *file.FilesHandler) (*bufio.Writer, error) {
	var outputFile *os.File
	if configFile.OutputPath != "" {
		outputFile, err := os.Create(configFile.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("opening output file for writing: %v", err)
		}
		openFiles.Add(outputFile)

	} else {
		slog.Debug("output path empty - defaulting to stdout")
		outputFile = os.Stdout
	}

	outputWriter := bufio.NewWriter(outputFile)

	return outputWriter, nil
}
