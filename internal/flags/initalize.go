package flags

import (
	"flag"
)

func LoadAndValidate() (string, string, string, bool, bool, bool) {
	var anonymizationDataPath string
	flag.Func("d", "Path to directory with anonymizing data", validateDir(anonymizationDataPath))

	var inputPath string
	flag.Func("i", "Path to input file containing logs to be anonymized", validateInput(inputPath))

	var outputPath string
	flag.Func("o", "Path to output file (default: Stdout)", validateOutput(outputPath))

	var isVerbose = flag.Bool("v", false, "Enable verbose logging (default: Disabled)")
	var isLmExport = flag.Bool("e", false, "Change input file type to LM export (default: LM Backup)")
	var isProofWriter = flag.Bool("p", true, "Disable proof wrtier (default: Enabled)")

	flag.Parse()

	return anonymizationDataPath, inputPath, outputPath, *isVerbose, *isLmExport, *isProofWriter
}
