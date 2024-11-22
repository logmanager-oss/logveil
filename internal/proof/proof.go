package proof

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/files"
)

type ProofWriter struct {
	IsEnabled bool
	writer    *bufio.Writer
	file      *os.File
}

func CreateProofWriter(config *config.Config, openFiles *files.FilesHandler) (*ProofWriter, error) {
	if config.IsProofWriter {
		file, err := os.Create("proof.json")
		if err != nil {
			return nil, fmt.Errorf("creating/opening proof file: %v", err)
		}
		openFiles.Add(file)

		return &ProofWriter{
			IsEnabled: true,
			writer:    bufio.NewWriter(file),
			file:      file,
		}, nil
	}

	return &ProofWriter{IsEnabled: false}, nil
}

func (p *ProofWriter) Write(originalValue string, newValue string) {
	if !p.IsEnabled {
		return
	}

	proof := struct {
		OriginalValue string `json:"original"`
		NewValue      string `json:"new"`
	}{
		OriginalValue: originalValue,
		NewValue:      newValue,
	}

	bytes, err := json.Marshal(proof)
	if err != nil {
		slog.Error("marshalling anonymisation proof", "error", err)
	}

	_, err = fmt.Fprintf(p.writer, "%s\n", bytes)
	if err != nil {
		slog.Error("writing anonymisation proof", "error", err)
	}
}

func (p *ProofWriter) Flush() {
	if !p.IsEnabled {
		return
	}

	p.writer.Flush()
}
