package anonymizer

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/logmanager-oss/logveil/internal/generator"
	"github.com/logmanager-oss/logveil/internal/lookup"
	"github.com/logmanager-oss/logveil/internal/proof"
	"golang.org/x/exp/rand"
)

type Anonymizer struct {
	anonData    map[string][]string
	randFunc    func(int) int
	proofWriter *proof.Proof
	lookup      *lookup.Lookup
	generator   *generator.Generator
}

func New(anonData map[string][]string, proofWriter *proof.Proof, lookup *lookup.Lookup, generator *generator.Generator) *Anonymizer {
	return &Anonymizer{
		anonData:    anonData,
		randFunc:    rand.Intn,
		proofWriter: proofWriter,
		lookup:      lookup,
		generator:   generator,
	}
}

func (an *Anonymizer) Anonymize(logLine map[string]string) string {
	logLine["raw"] = an.dynamicReplacements(logLine["raw"])
	logLine = an.staticReplacements(logLine)

	return logLine["raw"]
}

func (an *Anonymizer) SetRandFunc(randFunc func(int) int) {
	an.randFunc = randFunc
}

func (an *Anonymizer) staticReplacements(logLine map[string]string) map[string]string {
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

	return logLine
}

func (an *Anonymizer) dynamicReplacements(rawLog string) string {
	return an.lookup.ValidIp.ReplaceAllStringFunc(rawLog, func(original string) string {
		randIp := an.generator.GenerateRandomIPv4()
		an.proofWriter.Write(original, randIp)

		return randIp
	})
}
