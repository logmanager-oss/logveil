package logveil

import (
	"bufio"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/flags"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/runner"
)

func Run() {
	slog.Info("Anonymization process started...")

	anonymizingDataDir, inputPath, outputPath, isVerbose, isLmExport := flags.LoadAndValidate()

	if isVerbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	inputReader, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer inputReader.Close()

	var outputFile *os.File
	if outputPath != "" {
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()

	anonymizingData, err := loader.Load(anonymizingDataDir)
	if err != nil {
		slog.Error("loading anonymizing data from dir %s: %v", anonymizingDataDir, err)
		return
	}
	anonymizer := anonymizer.New(anonymizingData)

	if isLmExport {
		err := runner.AnonymizeLmExport(inputReader, outputWriter, anonymizer)
		if err != nil {
			slog.Error("reading lm export input file %s: %v", inputReader.Name(), err)
			return
		}
	} else {
		err := runner.AnonymizeLmBackup(inputReader, outputWriter, anonymizer)
		if err != nil {
			slog.Error("reading lm backup input file %s: %v", inputReader.Name(), err)
			return
		}
	}

	slog.Info("All done. Exiting...")
}
