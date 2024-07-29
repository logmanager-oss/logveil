package flags

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type input string

func (f *input) String() string {
	return fmt.Sprint(*f)
}

func (f *input) Set(value string) error {
	_, err := os.Stat(value)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("Provided file or dir %s does not exist. Aborting.", value)
		}
	}

	*f = input(value)

	return nil
}

type output string

func (f *output) String() string {
	return fmt.Sprint(*f)
}

func (f *output) Set(value string) error {
	file, err := os.Create(value)
	if err != nil {
		return err
	}
	defer file.Close()

	*f = output(value)

	return nil
}
