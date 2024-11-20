package anonymizer

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/generator"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/lookup"
	"github.com/logmanager-oss/logveil/internal/proof"
	"golang.org/x/exp/rand"
)

// Anonymizer represents an object responsible for anonymizing indivisual log lines feed to it. It contains anonymization data which will be used to anonymize input and a random number generator funtion used to select values from anonymization data.
type Anonymizer struct {
	anonData       map[string][]string
	randFunc       func(int) int
	proofWriter    *proof.ProofWriter
	lookup         *lookup.Lookup
	generator      *generator.Generator
	replacementMap map[string]string
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
		lookup:      lookup.New(),
		generator:   &generator.Generator{},
	}, nil
}

func (an *Anonymizer) Anonymize(logLine map[string]string) string {
	an.replacementMap = make(map[string]string)

	an.loadAndReplace(logLine)

	logLineRaw := logLine["raw"]
	an.generateAndReplace(logLineRaw, an.lookup.ValidIpv4, an.generator.GenerateRandomIPv4())
	an.generateAndReplace(logLineRaw, an.lookup.ValidIpv6, an.generator.GenerateRandomIPv6())
	an.generateAndReplace(logLineRaw, an.lookup.ValidMac, an.generator.GenerateRandomMac())
	an.generateAndReplace(logLineRaw, an.lookup.ValidEmail, an.generator.GenerateRandomEmail())
	an.generateAndReplace(logLineRaw, an.lookup.ValidUrl, an.generator.GenerateRandomUrl())

	an.proofWriter.Write(an.replacementMap)
	an.proofWriter.Flush()

	return an.replace(logLineRaw)
}

// SetRandFunc sets the function used by Anonymize() to select values from anonymization data at random
func (an *Anonymizer) SetRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}

func (an *Anonymizer) loadAndReplace(logLine map[string]string) {
	for field, value := range logLine {
		if field == "raw" {
			continue
		}

		if value == "" {
			continue
		}

		if _, ok := an.replacementMap[value]; ok {
			continue
		}

		if anonValues, exists := an.anonData[field]; exists {
			newAnonValue := anonValues[an.randFunc(len(anonValues))]
			an.replacementMap[value] = newAnonValue

			slog.Debug(fmt.Sprintf("Replacing the values for field %s. From %s to %s.\n", field, value, newAnonValue))
		}
	}
}

func (an *Anonymizer) generateAndReplace(rawLog string, regexp *regexp.Regexp, generatedData string) {
	values := regexp.FindAllString(rawLog, -1)

	for _, value := range values {
		if _, ok := an.replacementMap[value]; ok {
			continue
		}

		an.replacementMap[value] = generatedData
	}
}

func (an *Anonymizer) replace(rawLog string) string {
	for oldValue, newValue := range an.replacementMap {
		rawLog = strings.ReplaceAll(rawLog, oldValue, newValue)
	}

	return rawLog
}
