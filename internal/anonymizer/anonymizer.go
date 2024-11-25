package anonymizer

import (
	"fmt"
	"log/slog"
	"maps"
	"regexp"

	"math/rand/v2"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/generator"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/lookup"
	"github.com/logmanager-oss/logveil/internal/proof"
)

// Anonymizer represents an object responsible for anonymizing indivisual log lines feed to it. It contains anonymization data which will be used to anonymize input and a random number generator funtion used to select values from anonymization data.
type Anonymizer struct {
	anonymizationData       map[string][]string
	replacementMap          map[string]string
	randFunc                func(int) int
	proofWriter             *proof.ProofWriter
	lookup                  *lookup.Lookup
	generator               *generator.Generator
	isPersistReplacementMap bool
}

func CreateAnonymizer(config *config.Config, proofWriter *proof.ProofWriter) (*Anonymizer, error) {
	customReplacementMap, err := loader.LoadCustomReplacementMap(config.CustomReplacementMapPath)
	if err != nil {
		return nil, fmt.Errorf("loading custom replacement map from path %s: %v", config.CustomReplacementMapPath, err)
	}

	anonymizationData, err := loader.LoadAnonymizationData(config.AnonymizationDataPath)
	if err != nil {
		return nil, fmt.Errorf("loading anonymizing data from dir %s: %v", config.AnonymizationDataPath, err)
	}

	return &Anonymizer{
		anonymizationData:       anonymizationData,
		replacementMap:          customReplacementMap,
		randFunc:                rand.IntN,
		proofWriter:             proofWriter,
		lookup:                  lookup.New(),
		generator:               &generator.Generator{},
		isPersistReplacementMap: config.IsPersistReplacementMap,
	}, nil
}

func (an *Anonymizer) Anonymize(logLine map[string]string) string {
	replacementMap := an.loadAndReplace(logLine, an.replacementMap)

	logLineRaw := logLine["raw"]
	replacementMap = an.generateAndReplace(logLineRaw, replacementMap, an.lookup.ValidIpv4, an.generator.GenerateRandomIPv4())
	replacementMap = an.generateAndReplace(logLineRaw, replacementMap, an.lookup.ValidIpv6, an.generator.GenerateRandomIPv6())
	replacementMap = an.generateAndReplace(logLineRaw, replacementMap, an.lookup.ValidMac, an.generator.GenerateRandomMac())
	replacementMap = an.generateAndReplace(logLineRaw, replacementMap, an.lookup.ValidEmail, an.generator.GenerateRandomEmail())
	replacementMap = an.generateAndReplace(logLineRaw, replacementMap, an.lookup.ValidUrl, an.generator.GenerateRandomUrl())

	if an.isPersistReplacementMap {
		maps.Copy(an.replacementMap, replacementMap)
	}

	return an.replace(logLineRaw, replacementMap)
}

func (an *Anonymizer) loadAndReplace(logLine map[string]string, replacementMap map[string]string) map[string]string {
	for field, value := range logLine {
		if field == "raw" {
			continue
		}

		if value == "" {
			continue
		}

		if _, ok := replacementMap[value]; ok {
			continue
		}

		if anonValues, exists := an.anonymizationData[field]; exists {
			newAnonValue := anonValues[an.randFunc(len(anonValues))]
			replacementMap[value] = newAnonValue

			slog.Debug(fmt.Sprintf("Replacing the values for field %s. From %s to %s.\n", field, value, newAnonValue))
		}
	}

	return replacementMap
}

func (an *Anonymizer) generateAndReplace(rawLog string, replacementMap map[string]string, regexp *regexp.Regexp, generatedData string) map[string]string {
	values := regexp.FindAllString(rawLog, -1)

	for _, value := range values {
		if _, ok := an.replacementMap[value]; ok {
			continue
		}

		replacementMap[value] = generatedData
	}

	return replacementMap
}

func (an *Anonymizer) replace(rawLog string, replacementMap map[string]string) string {
	for originalValue, newValue := range replacementMap {
		// Added word boundary to avoid matching words withing word. For example "test" in "testing".
		r := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, originalValue))

		var found bool
		rawLog = r.ReplaceAllStringFunc(rawLog, func(originalValue string) string {
			found = true
			return newValue
		})

		if found {
			an.proofWriter.Write(originalValue, newValue)
		}
	}

	return rawLog
}
