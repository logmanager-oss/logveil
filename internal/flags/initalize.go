package flags

import (
	"flag"
)

func LoadAndValidate() (string, string, string, bool, bool) {
	var anonymizationDataPath string
	flag.Func("d", "Path to directory with anonymizing data", validateDir(anonymizationDataPath))

	var inputPath string
	flag.Func("i", "Path to input file containing logs to be anonymized", validateInput(inputPath))

	var outputPath string
	flag.Func("o", "Path to output file (default: Stdout)", validateOutput(outputPath))

	var isVerbose = flag.Bool("v", false, "Enable verbose logging")
	var isLmExport = flag.Bool("e", false, "Change input file type to LM export (default input file type is LM Backup)")

	flag.Parse()

	return anonymizationDataPath, inputPath, outputPath, *isVerbose, *isLmExport
}
