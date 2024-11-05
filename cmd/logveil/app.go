package logveil

import (
	"io"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/flags"
	"github.com/logmanager-oss/logveil/internal/inputs"
	"github.com/logmanager-oss/logveil/internal/parser"
)

func Run() {
	slog.Info("Anonymization process started...")

	anonDataDir, inputFilename, outputFilename, enableLMexport := flags.Load()

	input, err := os.Open(inputFilename)
	if err != nil {
		slog.Error("Opening input file", "error", err)
		return
	}
	defer closeFile(input)

	var output io.Writer
	if outputFilename != "" {
		output, err = os.Create(outputFilename)
		if err != nil {
			slog.Error("Opening input file", "error", err)
			return
		}
		defer closeFile(output.(*os.File))
	} else {
		output = os.Stdout
	}

	anonData, err := parser.LoadAnonData(anonDataDir)
	if err != nil {
		slog.Error("loading anonymizing data from dir %s: %v", anonDataDir, err)
		return
	}
	anonymizer := anonymizer.New(anonData)

	if enableLMexport {
		err := inputs.AnonymizeLmExport(input, output, anonymizer)
		if err != nil {
			slog.Error("reading lm export input file %s: %v", inputFilename, err)
			return
		}
	} else {
		err := inputs.AnonymizeLmBackup(input, output, anonymizer)
		if err != nil {
			slog.Error("reading lm backup input file %s: %v", inputFilename, err)
			return
		}
	}

	slog.Info("All done. Exiting...")
}

func closeFile(fs *os.File) {
	err := fs.Close()
	if err != nil {
		slog.Error(err.Error())
	}
}
