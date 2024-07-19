package flags

import (
	"flag"
	"log/slog"
)

func Load() (string, string, string) {
	var anonDataDir input
	flag.Var(&anonDataDir, "d", "Path to directory with anonymizing data")

	var inputFile input
	flag.Var(&inputFile, "i", "Path to input file containing logs to be anonymized")

	var outputFile output
	flag.Var(&outputFile, "o", "Path to output file containing anonymized logs")

	var verbose = flag.Bool("v", false, "Enable verbose logging")
	flag.Parse()

	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return anonDataDir.String(), inputFile.String(), outputFile.String()
}
