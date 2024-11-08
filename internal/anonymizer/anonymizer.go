package anonymizer

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/logmanager-oss/logveil/internal/proof"
	"golang.org/x/exp/rand"
)

type Anonymizer struct {
	anonData    map[string][]string
	randFunc    func(int) int
	proofWriter *proof.Proof
}

func New(anonData map[string][]string, proofWriter *proof.Proof) *Anonymizer {
	return &Anonymizer{
		anonData:    anonData,
		randFunc:    rand.Intn,
		proofWriter: proofWriter,
	}
}

func (an *Anonymizer) Anonymize(logLine map[string]string) string {
	for field, value := range logLine {
		if field == "raw" {
			continue
		}

		if value == "" {
			continue
		}

		if anonValues, exists := an.anonData[field]; exists {
			newAnonValue := anonValues[an.randFunc(len(anonValues))]

			an.proofWriter.Write(value, newAnonValue)

			slog.Debug(fmt.Sprintf("Replacing the values for field %s. From %s to %s.\n", field, value, newAnonValue))

			logLine["raw"] = strings.Replace(logLine["raw"], value, newAnonValue, -1)
		}
	}

	return logLine["raw"]
}

func (an *Anonymizer) SetRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}
