package inputs

import (
	"bytes"
	"math/rand"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/logmanager-oss/logveil/internal/anonymizer"
	"github.com/logmanager-oss/logveil/internal/generator"
	"github.com/logmanager-oss/logveil/internal/lookup"
	"github.com/logmanager-oss/logveil/internal/parser"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/logmanager-oss/logveil/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestLmBackup(t *testing.T) {
	tests := []struct {
		name                 string
		isProofWriterEnabled bool
		inputFilename        string
		anonDataDir          string
		expectedOutput       string
		expectedProof        []map[string]interface{}
	}{
		{
			name:                 "Test Test LM Backup Anonymizer",
			isProofWriterEnabled: true,
			inputFilename:        "../../examples/logs/lm-2024-06-09_0000.gz",
			anonDataDir:          "../../examples/anon_data",
			expectedOutput:       "<189>date=2024-11-06 time=12:29:25 devname=\"LM-FW-70F-Praha\" devid=\"FGT70FTK22012016\" eventtime=1730892565525108329 tz=\"+0100\" logid=\"0000000013\" type=\"traffic\" subtype=\"forward\" level=\"notice\" vd=\"root\" srcip=33.15.199.187 srcport=57158 srcintf=\"lan1\" srcintfrole=\"wan\" dstip=129.134.57.172 dstport=80 dstintf=\"lan1\" dstintfrole=\"lan\" srccountry=\"China\" dstcountry=\"Czech Republic\" sessionid=179455916 proto=6 action=\"client-rst\" policyid=9 policytype=\"policy\" poluuid=\"d8ccb3e4-74d4-51ef-69a3-73b41f46df74\" policyname=\"Gitlab web from all\" service=\"HTTP\" trandisp=\"noop\" duration=6 sentbyte=80 rcvdbyte=44 sentpkt=2 rcvdpkt=1 appcat=\"unscanned\" srchwvendor=\"H3C\" devtype=\"Router\" mastersrcmac=\"00:23:89:39:a4:ef\" srcmac=\"00:23:89:39:a4:ef\" srcserver=0 dsthwvendor=\"H3C\" dstdevtype=\"Router\" masterdstmac=\"00:23:89:39:a4:fa\" dstmac=\"00:23:89:39:a4:fa\" dstserver=0\n",
			expectedProof: []map[string]interface{}{
				{"original": "dev-uplink", "new": "lan1"},
				{"original": "95.80.197.108", "new": "129.134.57.172"},
				{"original": "27.221.126.209", "new": "33.15.199.187"},
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

			anonData, err := parser.LoadAnonData(tt.anonDataDir)
			if err != nil {
				t.Fatal(err)
			}

			proofWriter := proof.New(tt.isProofWriterEnabled)
			lookup := lookup.New()
			generator := &generator.Generator{}
			anonymizer := anonymizer.New(anonData, proofWriter, lookup, generator)
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })
			faker.SetRandomSource(rand.NewSource(1))

			var output bytes.Buffer
			err = AnonymizeLmBackup(input, &output, anonymizer)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedOutput, output.String())

			proofWriter.Close()

			actualProof, err := utils.UnpackProofOutput()
			if err != nil {
				t.Fatal(err)
			}

			assert.ElementsMatch(t, tt.expectedProof, actualProof)

			os.Remove("proof.json")
		})
	}
}
