package runner

import (
	"bytes"
	"os"
	"testing"

	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/loader"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/logmanager-oss/logveil/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestLmBackup(t *testing.T) {
	tests := []struct {
		name               string
		inputFilename      string
		anonymizingDataDir string
		expectedOutput     string
		expectedProof      []map[string]interface{}
	}{
		{
			name:               "Test Test LM Backup Anonymizer",
			inputFilename:      "../../examples/logs/lm-2024-06-09_0000.gz",
			anonymizingDataDir: "../../examples/anon_data",
			expectedOutput:     "<189>date=2024-11-06 time=12:29:25 devname=\"LM-FW-70F-Praha\" devid=\"FGT70FTK22012016\" eventtime=1730892565525108329 tz=\"+0100\" logid=\"0000000013\" type=\"traffic\" subtype=\"forward\" level=\"notice\" vd=\"root\" srcip=10.20.0.53 srcport=57158 srcintf=\"lan1\" srcintfrole=\"wan\" dstip=227.51.221.89 dstport=80 dstintf=\"lan1\" dstintfrole=\"lan\" srccountry=\"China\" dstcountry=\"Czech Republic\" sessionid=179455916 proto=6 action=\"client-rst\" policyid=9 policytype=\"policy\" poluuid=\"d8ccb3e4-74d4-51ef-69a3-73b41f46df74\" policyname=\"Gitlab web from all\" service=\"HTTP\" trandisp=\"noop\" duration=6 sentbyte=80 rcvdbyte=44 sentpkt=2 rcvdpkt=1 appcat=\"unscanned\" srchwvendor=\"H3C\" devtype=\"Router\" mastersrcmac=\"00:23:89:39:a4:ef\" srcmac=\"00:23:89:39:a4:ef\" srcserver=0 dsthwvendor=\"H3C\" dstdevtype=\"Router\" masterdstmac=\"00:23:89:39:a4:fa\" dstmac=\"00:23:89:39:a4:fa\" dstserver=0\n",
			expectedProof: []map[string]interface{}{
				{"original": "dev-uplink", "new": "lan1"},
				{"original": "95.80.197.108", "new": "227.51.221.89"},
				{"original": "27.221.126.209", "new": "10.20.0.53"},
				{"original": "wan1-lm", "new": "lan1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := os.Open(tt.inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			defer input.Close()

			var output bytes.Buffer

			anonymizingData, err := loader.Load(tt.anonymizingDataDir)
			if err != nil {
				t.Fatal(err)
			}
			proofWriter := proof.New(true)
			anonymizer := anonymizer.New(anonymizingData, proofWriter)
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })

			err = AnonymizeLmBackup(input, &output, anonymizer)
			if err != nil {
				t.Fatal(err)
			}

			proofWriter.Close()

			actualProof, err := utils.UnpackProofOutput()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expectedProof, actualProof)

			os.Remove("proof.json")

			assert.Equal(t, tt.expectedOutput, output.String())
		})
	}
}
