package handlers

import (
	"log/slog"
	"os"
)

type Files struct {
	files []*os.File
}

func (f *Files) Add(file *os.File) {
	f.files = append(f.files, file)
}

func (f *Files) Close() {
	slog.Debug("closing files on exit...")
	for _, file := range f.files {
		file.Close()
	}
}
