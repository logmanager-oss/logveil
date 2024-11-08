package proof

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/logmanager-oss/logveil/internal/utils"
)

type Proof struct {
	isEnabled bool
	writer    *bufio.Writer
	file      *os.File
}

func New(isEnabled bool) *Proof {
	var writer *bufio.Writer
	var file *os.File

	if isEnabled {
		var err error
		proofFile, err := os.OpenFile("proof.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			slog.Error("opening/creating proof file", "error", err)
			return nil
		}

		writer = bufio.NewWriter(proofFile)
		file = proofFile
	}

	return &Proof{
		isEnabled: isEnabled,
		writer:    writer,
		file:      file,
	}
}

func (p *Proof) Write(originalValue string, maskedValue string) {
	if !p.isEnabled {
		return
	}

	proof := struct {
		OriginalValue string `json:"original"`
		MaskedValue   string `json:"new"`
	}{
		OriginalValue: originalValue,
		MaskedValue:   maskedValue,
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

func (p *Proof) Close() {
	if !p.isEnabled {
		return
	}

	err := p.writer.Flush()
	if err != nil {
		slog.Error("flushing buffer", "error", err)
	}

	utils.CloseFile(p.file)
}
