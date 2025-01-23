package logveil

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/handlers"
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

	filesHandler := &handlers.Files{}
	defer filesHandler.Close()

	buffersHandler := &handlers.Buffers{}
	defer buffersHandler.Flush()

	inputReader, err := reader.CreateInputReader(config, filesHandler)
	if err != nil {
		slog.Error("initializing input reader", "error", err)
		return
	}
	outputWriter, err := writer.CreateOutputWriter(config, filesHandler, buffersHandler)
	if err != nil {
		slog.Error("initializing output writer", "error", err)
		return
	}
	proofWriter, err := proof.CreateProofWriter(config, filesHandler, buffersHandler)
	if err != nil {
		slog.Error("initializing proof writer", "error", err)
		return
	}
	anonymizerDoer, err := anonymizer.CreateAnonymizer(config, proofWriter)
	if err != nil {
		slog.Error("initializing anonymizer", "error", err)
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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
			return fmt.Errorf("writing log line %s: %v", anonymizedLogLine, err)
		}

		select {
		case <-ctx.Done():
			fmt.Println("\nInterrupt received, closing...")
			stop()
			return nil
		default:
			continue
		}
	}
}
