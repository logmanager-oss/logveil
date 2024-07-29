package anonymizer

import (
	"fmt"
	"log/slog"

	"github.com/logmanager-oss/logveil/internal/flags"
	"github.com/logmanager-oss/logveil/internal/parser"
	"github.com/logmanager-oss/logveil/internal/writer"
)

func Run() {
	slog.Info("Anonymization process started...")

	anonDataDir, inputFile, outputFile := flags.Load()

	fieldNames, csvData, err := parser.ParseCSV(inputFile)
	if err != nil {
		slog.Error("reading input file %s: %v", inputFile, err)
		return
	}

	anonData, err := parser.ParseAnonData(anonDataDir, fieldNames)
	if err != nil {
		slog.Error("loading anonymizing data from dir %s: %v", anonDataDir, err)
		return
	}

	anonymizer := New(csvData, anonData)
	anonymizedData := anonymizer.anonymize()
	if outputFile != "" {
		outputwriter := &writer.Output{
			Output: anonymizedData,
		}
		err := outputwriter.Write(outputFile)
		if err != nil {
			slog.Error("writing anonymized data to output file %s: %v", outputFile, err)
		}
	} else {
		fmt.Println(anonymizedData)
	}

	slog.Info("All done. Exiting...")
}
