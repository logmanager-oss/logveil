package config

import (
	"flag"
)

// Config represents user supplied program input
type Config struct {
	AnonymizationDataPath string
	InputPath             string
	OutputPath            string
	IsVerbose             *bool
	IsLmExport            *bool
}

// LoadAndValidate loads values from user supplied input into Config struct and validates them
func (c *Config) LoadAndValidate() {
	flag.Func("d", "Path to directory with anonymizing data", validateDir(c.AnonymizationDataPath))

	flag.Func("i", "Path to input file containing logs to be anonymized", validateInput(c.InputPath))

	flag.Func("o", "Path to output file (default: Stdout)", validateOutput(c.OutputPath))

	c.IsVerbose = flag.Bool("v", false, "Enable verbose logging (default: Disabled)")
	c.IsLmExport = flag.Bool("e", false, "Change input file type to LM export (default: LM Backup)")

	flag.Parse()
}
