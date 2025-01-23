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
	ReaderMaxCapacity        int
	WriterMaxCapacity        int
}

// LoadAndValidate loads values from user supplied input into Config struct and validates them
func (c *Config) LoadAndValidate() {
	flag.Func("d", "Path to directory with anonymizing data", c.validateDirPath())

	flag.Func("i", "Path to input file containing logs to be anonymized", c.validateInputPath())

	flag.Func("c", "Path to input file containing custom anonymization mappings", c.validateCustomMappingPath())

	flag.Func("o", "Path to output file (default: Stdout)", c.validateOutput())

	flag.BoolVar(&c.IsVerbose, "v", false, "Enable verbose logging (default: Disabled)")
	flag.BoolVar(&c.IsLmExport, "e", false, "Change input file type to LM export (default: LM Backup)")
	flag.BoolVar(&c.IsProofWriter, "p", false, "Enable proof writer (default: Disabled)")
	flag.BoolVar(&c.IsPersistReplacementMap, "r", false, "Enable persistent (per session) replacement map (default: Disabled)")

	c.ReaderMaxCapacity = *flag.Int("rs", 4000000, "Size of the reader buffer in Bytes")
	c.WriterMaxCapacity = *flag.Int("ws", 2000000, "Size of the writer buffer in Bytes")

	flag.Parse()

	// Check if mandatory flags are set
	if c.InputPath == "" {
		fmt.Println("Error: -i flag is mandatory")
		flag.Usage()
		os.Exit(1)
	}
}
