package logveil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/files"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/logmanager-oss/logveil/internal/reader"
	"github.com/logmanager-oss/logveil/internal/writer"
)

func Start() {
	slog.Info("Anonymization process started...")

	config := &config.Config{}
	config.LoadAndValidate()

	if config.IsVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	filesHandler := &files.FilesHandler{}
	defer filesHandler.Close()

	inputReader, err := reader.CreateInputReader(config, filesHandler)
	if err != nil {
		return
	}
	outputWriter, err := writer.CreateOutputWriter(config, filesHandler)
	if err != nil {
		return
	}
	proofWriter, err := proof.CreateProofWriter(config, filesHandler)
	if err != nil {
		return
	}
	anonymizerDoer, err := anonymizer.CreateAnonymizer(config, proofWriter)
	if err != nil {
		return
	}

	err = RunAnonymizationLoop(inputReader, outputWriter, anonymizerDoer)
	if err != nil {
		slog.Error("running anonymisation loop", "error", err)
		return
	}

	slog.Info("All done. Exiting...")
}

func RunAnonymizationLoop(inputReader reader.InputReader, outputWriter *bufio.Writer, anonymizerDoer *anonymizer.Anonymizer) error {
	defer outputWriter.Flush()

	for {
		logLine, err := inputReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("reading line: %v", err)
		}

		anonymizedLogLine := anonymizerDoer.Anonymize(logLine)

		_, err = fmt.Fprintln(outputWriter, anonymizedLogLine)
		if err != nil {
			return fmt.Errorf("writing log line to buffer: %v", err)
		}
	}
}
