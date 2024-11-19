package config

import (
	"errors"
	"fmt"
	"os"
)

func validateInput(inputPath string) func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Input file %s cannot be a directory.\n", flagValue)
		}

		inputPath = flagValue

		return nil
	}
}

func validateOutput(outputPath string) func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}

		if fileInfo.IsDir() {
			return fmt.Errorf("Output file %s cannot be a directory.\n", flagValue)
		}

		outputPath = flagValue

		return nil
	}
}

func validateDir(dir string) func(string) error {
	return func(flagValue string) error {
		fileInfo, err := os.Stat(flagValue)
		if err != nil {
			return err
		}

		if !fileInfo.IsDir() {
			return fmt.Errorf("Path to anonymization data %s needs to be a directory.\n", flagValue)
		}

		dir = flagValue

		return nil
	}
}
