package handlers

import (
	"bufio"
	"log/slog"
)

type Buffers struct {
	buffers []*bufio.Writer
}

func (b *Buffers) Add(buffer *bufio.Writer) {
	b.buffers = append(b.buffers, buffer)
}

func (b *Buffers) Flush() {
	slog.Debug("flush buffers on exit...")
	for _, buffer := range b.buffers {
		err := buffer.Flush()
		if err != nil {
			slog.Error("buffer flush failed", "error", err)
		}
	}
}
