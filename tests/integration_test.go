package testing

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/cmd/logveil"
	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/reader"
	"github.com/stretchr/testify/assert"
)

func TestLogVeil_IntegrationTest(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		isLmExport         bool
		anonymizingDataDir string
		expectedOutput     string
	}{
		{
			name:               "Test Test LM Backup Anonymizer",
			inputFilename:      "data/lm_backup_test_input.gz",
			isLmExport:         false,
			anonymizingDataDir: "data/anonymization_data",
			expectedOutput:     "<189>date=2024-11-06 time=12:29:25 devname=\"LM-FW-70F-Praha\" devid=\"FGT70FTK22012016\" eventtime=1730892565525108329 tz=\"+0100\" logid=\"0000000013\" type=\"traffic\" subtype=\"forward\" level=\"notice\" vd=\"root\" srcip=10.20.0.53 srcport=57158 srcintf=\"lan1\" srcintfrole=\"wan\" dstip=227.51.221.89 dstport=80 dstintf=\"lan1\" dstintfrole=\"lan\" srccountry=\"China\" dstcountry=\"Czech Republic\" sessionid=179455916 proto=6 action=\"client-rst\" policyid=9 policytype=\"policy\" poluuid=\"d8ccb3e4-74d4-51ef-69a3-73b41f46df74\" policyname=\"Gitlab web from all\" service=\"HTTP\" trandisp=\"noop\" duration=6 sentbyte=80 rcvdbyte=44 sentpkt=2 rcvdpkt=1 appcat=\"unscanned\" srchwvendor=\"H3C\" devtype=\"Router\" mastersrcmac=\"00:23:89:39:a4:ef\" srcmac=\"00:23:89:39:a4:ef\" srcserver=0 dsthwvendor=\"H3C\" dstdevtype=\"Router\" masterdstmac=\"00:23:89:39:a4:fa\" dstmac=\"00:23:89:39:a4:fa\" dstserver=0\n",
		},
		{
			name:               "Test LM Export Anonymizer",
			inputFilename:      "data/lm_export_test_input.csv",
			isLmExport:         true,
			anonymizingDataDir: "data/anonymization_data",
			expectedOutput:     "{\"@timestamp\": \"2024-06-05T14:59:27.000+00:00\", \"msg.src_ip\":\"10.20.0.53\", \"username\":\"ladislav.dosek\", \"organization\":\"Apple\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile, err := os.Open(tt.inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			var inputReader reader.InputReader
			if tt.isLmExport {
				inputReader, err = reader.NewLmExportReader(inputFile)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				inputReader, err = reader.NewLmBackupReader(inputFile)
				if err != nil {
					t.Fatal(err)
				}
			}

			var output bytes.Buffer
			outputWriter := bufio.NewWriter(&output)

			anonymizer, err := anonymizer.CreateAnonymizer(&config.Config{AnonymizationDataPath: tt.anonymizingDataDir})
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
		})
	}
}
