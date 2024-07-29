package anonymizer

import (
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/exp/rand"
)

type Anonymizer struct {
	csvData  []map[string]string
	anonData map[string][]string
	randFunc func(int) int
}

func New(csvData []map[string]string, anonData map[string][]string) *Anonymizer {
	return &Anonymizer{
		csvData:  csvData,
		anonData: anonData,
		randFunc: rand.Intn,
	}
}

func (an *Anonymizer) anonymize() []string {
	var output []string
	for _, logLine := range an.csvData {
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

		output = append(output, fmt.Sprint(logLine["raw"]))
	}

	return output
}

func (an *Anonymizer) setRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}
