package writer

import (
	"fmt"
	"log/slog"
	"os"
)

type Output struct {
	Output []string
}

func (o *Output) Write(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(fs *os.File) {
		if err := fs.Close(); err != nil {
			slog.Error(err.Error())
		}
	}(file)

	for _, line := range o.Output {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("writing anonymized data to output file %s: %v", filename, err)
		}
	}

	return nil
}
