package writer

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/handlers"
)

func CreateOutputWriter(config *config.Config, filesHandler *handlers.Files, buffersHandler *handlers.Buffers) (*bufio.Writer, error) {
	var outputFile *os.File
	var err error
	if config.OutputPath != "" {
		outputFile, err = os.Create(config.OutputPath)
		if err != nil {
			return nil, fmt.Errorf("opening output file for writing: %v", err)
		}
		filesHandler.Add(outputFile)

	} else {
		slog.Debug("output path empty - defaulting to stdout")
		outputFile = os.Stdout
	}

	outputWriter := bufio.NewWriterSize(outputFile, config.WriterMaxCapacity)

	buffersHandler.Add(outputWriter)

	return outputWriter, nil
}
