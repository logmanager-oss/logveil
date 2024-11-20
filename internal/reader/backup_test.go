package reader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLmBackup(t *testing.T) {
	tests := []struct {
		name           string
		inputFilename  string
		expectedOutput map[string]string
		wantErr        bool
		expectedErr    error
	}{
		{
			name:           "Test LM Backup Anonymizer",
			inputFilename:  "../../tests/data/lm_backup_test_input.gz",
			expectedOutput: map[string]string{"appcat": "unscanned", "device_id": "FGT70FTK22012016", "device_name": "LM-FW-70F-Praha", "dst_iface": "dev-uplink", "dst_ip": "95.80.197.108", "dst_ip@ip": "map[as_number:29208 as_organization:Quantcom, a.s. city:Unknown country_code:CZ country_name:Czechia is_link_local:false is_multicast:false is_reserved:false ptr:95.80.197.108 value:95.80.197.108 version:4]", "dst_port": "80", "dst_port@int": "map[value:80]", "duration": "6.0", "duration@float": "map[value:6]", "policy_id": "9", "protocol": "TCP", "raw": "<189>date=2024-11-06 time=12:29:25 devname=\"LM-FW-70F-Praha\" devid=\"FGT70FTK22012016\" eventtime=1730892565525108329 tz=\"+0100\" logid=\"0000000013\" type=\"traffic\" subtype=\"forward\" level=\"notice\" vd=\"root\" srcip=27.221.126.209 srcport=57158 srcintf=\"wan1-lm\" srcintfrole=\"wan\" dstip=95.80.197.108 dstport=80 dstintf=\"dev-uplink\" dstintfrole=\"lan\" srccountry=\"China\" dstcountry=\"Czech Republic\" sessionid=179455916 proto=6 action=\"client-rst\" policyid=9 policytype=\"policy\" poluuid=\"d8ccb3e4-74d4-51ef-69a3-73b41f46df74\" policyname=\"Gitlab web from all\" service=\"HTTP\" trandisp=\"noop\" duration=6 sentbyte=80 rcvdbyte=44 sentpkt=2 rcvdpkt=1 appcat=\"unscanned\" srchwvendor=\"H3C\" devtype=\"Router\" mastersrcmac=\"00:23:89:39:a4:ef\" srcmac=\"00:23:89:39:a4:ef\" srcserver=0 dsthwvendor=\"H3C\" dstdevtype=\"Router\" masterdstmac=\"00:23:89:39:a4:fa\" dstmac=\"00:23:89:39:a4:fa\" dstserver=0", "rcvd_byte": "44", "rcvd_byte@int": "map[value:44]", "rcvd_pkt": "1", "rcvd_pkt@int": "map[value:1]", "sent_byte": "80", "sent_byte@int": "map[value:80]", "sent_pkt": "2", "sent_pkt@int": "map[value:2]", "service": "HTTP", "src_iface": "wan1-lm", "src_ip": "27.221.126.209", "src_ip@ip": "map[as_number:4837 as_organization:CHINA UNICOM China169 Backbone city:Unknown country_code:CN country_name:China is_link_local:false is_multicast:false is_reserved:false ptr:27.221.126.209 value:27.221.126.209 version:4]", "src_port": "57158", "src_port@int": "map[value:57158]", "status": "client-rst", "subtype": "forward", "type": "traffic", "vd": "root"},
		},
		{
			name:           "Test LM Backup Anonymizer - RAW missing",
			inputFilename:  "../../tests/data/lm_backup_test_input_raw_missing.gz",
			expectedOutput: map[string]string{},
			wantErr:        true,
			expectedErr:    fmt.Errorf("Malformed lm backup file: unexpected end of JSON input"),
		},
		{
			name:           "Test LM Backup Anonymizer - RAW empty",
			inputFilename:  "../../tests/data/lm_backup_test_input_raw_empty.gz",
			expectedOutput: map[string]string{},
			wantErr:        true,
			expectedErr:    fmt.Errorf("Malformed lm backup file - raw field cannot be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile, err := os.Open(tt.inputFilename)
			if err != nil {
				t.Fatal(err)
			}
			defer inputFile.Close()

			inputReader, err := NewLmBackupReader(inputFile)
			if err != nil {
				t.Fatal(err)
			}

			for {
				logLine, err := inputReader.ReadLine()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					if tt.wantErr {
						assert.Equal(t, tt.expectedErr, err)
						return
					}
					t.Fatal(err)
				}

				assert.Equal(t, tt.expectedOutput, logLine)
			}
		})
	}
}
