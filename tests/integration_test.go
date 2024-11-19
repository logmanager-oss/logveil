package testing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/cmd/logveil"
	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/files"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/logmanager-oss/logveil/internal/reader"
	"github.com/stretchr/testify/assert"
)

func TestLogVeil_IntegrationTest(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectedOutput string
		expectedProof  []map[string]interface{}
	}{
		{
			name: "Test Test LM Backup Anonymizer",
			config: &config.Config{
				AnonymizationDataPath: "data/anonymization_data",
				InputPath:             "data/lm_backup_test_input.gz",
				IsLmExport:            false,
				IsProofWriter:         true,
			},
			expectedOutput: "<189>date=2024-11-06 time=12:29:25 devname=\"LM-FW-70F-Praha\" devid=\"FGT70FTK22012016\" eventtime=1730892565525108329 tz=\"+0100\" logid=\"0000000013\" type=\"traffic\" subtype=\"forward\" level=\"notice\" vd=\"root\" srcip=10.20.0.53 srcport=57158 srcintf=\"lan1\" srcintfrole=\"wan\" dstip=227.51.221.89 dstport=80 dstintf=\"lan1\" dstintfrole=\"lan\" srccountry=\"China\" dstcountry=\"Czech Republic\" sessionid=179455916 proto=6 action=\"client-rst\" policyid=9 policytype=\"policy\" poluuid=\"d8ccb3e4-74d4-51ef-69a3-73b41f46df74\" policyname=\"Gitlab web from all\" service=\"HTTP\" trandisp=\"noop\" duration=6 sentbyte=80 rcvdbyte=44 sentpkt=2 rcvdpkt=1 appcat=\"unscanned\" srchwvendor=\"H3C\" devtype=\"Router\" mastersrcmac=\"00:23:89:39:a4:ef\" srcmac=\"00:23:89:39:a4:ef\" srcserver=0 dsthwvendor=\"H3C\" dstdevtype=\"Router\" masterdstmac=\"00:23:89:39:a4:fa\" dstmac=\"00:23:89:39:a4:fa\" dstserver=0\n",
			expectedProof: []map[string]interface{}{
				{"original": "dev-uplink", "new": "lan1"},
				{"original": "95.80.197.108", "new": "227.51.221.89"},
				{"original": "27.221.126.209", "new": "10.20.0.53"},
				{"original": "wan1-lm", "new": "lan1"},
			},
		},
		{
			name: "Test LM Export Anonymizer",
			config: &config.Config{
				AnonymizationDataPath: "data/anonymization_data",
				InputPath:             "data/lm_export_test_input.csv",
				IsLmExport:            true,
				IsProofWriter:         true,
			},
			expectedOutput: "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.20.0.53\", \"username\":\"ladislav.dosek\", \"organization\":\"Apple\"}\n",
			expectedProof: []map[string]interface{}{
				{"original": "89.239.31.49", "new": "10.20.0.53"},
				{"original": "test.user@test.cz", "new": "ladislav.dosek"},
				{"original": "TESTuser.test.com", "new": "Apple"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filesHandler := &files.FilesHandler{}
			defer filesHandler.Close()

			inputReader, err := reader.CreateInputReader(tt.config, filesHandler)
			if err != nil {
				t.Fatal(err)
			}

			var output bytes.Buffer
			outputWriter := bufio.NewWriter(&output)

			proofWriter, err := proof.CreateProofWriter(tt.config, filesHandler)
			if err != nil {
				t.Fatal(err)
			}

			anonymizer, err := anonymizer.CreateAnonymizer(tt.config, proofWriter)
			if err != nil {
				t.Fatal(err)
			}
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })

			err = logveil.RunAnonymizationLoop(inputReader, outputWriter, anonymizer)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedOutput, output.String())

			actualProof, err := unpackProofOutput()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expectedProof, actualProof)

			err = os.Remove("proof.json")
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func unpackProofOutput() ([]map[string]interface{}, error) {
	outputFile, err := os.OpenFile("proof.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	scanner := bufio.NewScanner(outputFile)
	for scanner.Scan() {
		var unpackedLine map[string]interface{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &unpackedLine)
		if err != nil {
			return nil, err
		}
		output = append(output, unpackedLine)
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return output, nil
}
