package config

import (
	"errors"
	"fmt"
	"os"
)

func (c *Config) validateInputPath() func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Input file %s cannot be a directory.\n", flagValue)
		}

		c.InputPath = flagValue

		return nil
	}
}

func (c *Config) validateCustomMappingPath() func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Path to custom mapping file %s cannot be a directory.\n", flagValue)
		}

		c.CustomReplacementMapPath = flagValue

		return nil
	}
}

func (c *Config) validateOutput() func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			// If output path does not exist it's ok - we will create it
			if errors.Is(err, os.ErrNotExist) {
				c.OutputPath = flagValue
				return nil
			}
			return err
		}

		// If output path exists check if it's a directory - which would be wrong
		if fileInfo.IsDir() {
			return fmt.Errorf("Output file %s cannot be a directory.\n", flagValue)
		}

		// If output path exists and is not a dir it's ok - file will be truncated
		c.OutputPath = flagValue

		return nil
	}
}

func (c *Config) validateDirPath() func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			return fmt.Errorf("Path to anonymization data %s needs to be a directory.\n", flagValue)
		}

		c.AnonymizationDataPath = flagValue

		return nil
	}
}
