package anonymizer

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/exp/rand"
)

type Anonymizer struct {
	anonData map[string][]string
	randFunc func(int) int
}

func New(anonData map[string][]string) *Anonymizer {
	return &Anonymizer{
		anonData: anonData,
		randFunc: rand.Intn,
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

			slog.Debug(fmt.Sprintf("Replacing the values for field %s. From %s to %s.\n", field, value, newAnonValue))

			logLine["raw"] = strings.Replace(logLine["raw"], value, newAnonValue, -1)
		}
	}

	return logLine["raw"]
}

func (an *Anonymizer) SetRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}
