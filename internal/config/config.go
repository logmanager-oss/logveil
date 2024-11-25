package config

import (
	"flag"
	"fmt"
	"os"
)

// Config represents user supplied program input
type Config struct {
	AnonymizationDataPath    string
	InputPath                string
	OutputPath               string
	CustomReplacementMapPath string
	IsVerbose                bool
	IsLmExport               bool
	IsProofWriter            bool
	IsPersistReplacementMap  bool
}

// LoadAndValidate loads values from user supplied input into Config struct and validates them
func (c *Config) LoadAndValidate() {
	flag.Func("d", "Path to directory with anonymizing data", validateDir(c.AnonymizationDataPath))

	flag.Func("i", "Path to input file containing logs to be anonymized", validateInput(c.InputPath))

	flag.Func("c", "Path to input file containing custom anonymization mappings", validateInput(c.CustomReplacementMapPath))

	flag.Func("o", "Path to output file (default: Stdout)", validateOutput(c.OutputPath))

	flag.BoolVar(&c.IsVerbose, "v", false, "Enable verbose logging (default: Disabled)")
	flag.BoolVar(&c.IsLmExport, "e", false, "Change input file type to LM export (default: LM Backup)")
	flag.BoolVar(&c.IsProofWriter, "p", true, "Disable proof writer (default: Enabled)")
	flag.BoolVar(&c.IsPersistReplacementMap, "r", true, "Disable persistent (per session) replacement map (default: Enabled)")

	flag.Parse()

	// Check if mandatory flags are set
	if c.InputPath == "" {
		fmt.Println("Error: -i flag is mandatory")
		flag.Usage()
		os.Exit(1)
	}
}
