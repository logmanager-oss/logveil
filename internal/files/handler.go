package files

import "os"

type FilesHandler struct {
	files []*os.File
}

func (fh *FilesHandler) Add(file *os.File) {
	fh.files = append(fh.files, file)
}

func (fh *FilesHandler) Close() {
	for _, file := range fh.files {
		file.Close()
	}
}
