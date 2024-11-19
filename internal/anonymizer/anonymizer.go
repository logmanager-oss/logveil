package anonymizer

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/proof"
	"golang.org/x/exp/rand"
)

// Anonymizer represents an object responsible for anonymizing indivisual log lines feed to it. It contains anonymization data which will be used to anonymize input and a random number generator funtion used to select values from anonymization data.
type Anonymizer struct {
	anonData    map[string][]string
	randFunc    func(int) int
	proofWriter *proof.ProofWriter
}

func CreateAnonymizer(config *config.Config, proofWriter *proof.ProofWriter) (*Anonymizer, error) {
	anonymizingData, err := loader.Load(config.AnonymizationDataPath)
	if err != nil {
		return nil, fmt.Errorf("loading anonymizing data from dir %s: %v", config.AnonymizationDataPath, err)
	}

	return &Anonymizer{
		anonData:    anonymizingData,
		randFunc:    rand.Intn,
		proofWriter: proofWriter,
	}, nil
}

func (an *Anonymizer) Anonymize(logLine map[string]string) string {
	defer an.proofWriter.Flush()

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

// SetRandFunc sets the function used by Anonymize() to select values from anonymization data at random
func (an *Anonymizer) SetRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}
