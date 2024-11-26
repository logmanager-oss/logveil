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
			expectedOutput: map[string]string{"collisions": "0", "collisions@int": "map[value:0]", "device_model": "XGS126", "device_name": "SFW", "device_serial_id": "X121073DFG3TY7C", "display_interface": "Port6", "interface": "Port6", "log_component": "Interface", "log_id": "123526618031", "log_subtype": "Usage", "log_type": "System Health", "log_version": "1", "meta_host": "", "meta_pid": "", "meta_program": "", "raw": "<30>device_name=\"SFW\" timestamp=\"2024-10-02T17:05:55+0200\" device_model=\"XGS126\" device_serial_id=\"X121073DFG3TY7C\" log_id=123526618031 log_type=\"System Health\" log_component=\"Interface\" log_subtype=\"Usage\" log_version=1 severity=\"Information\" display_interface=\"Port6\" interface=Port6 receivedkbits=1.16 transmittedkbits=1.53 receivederrors=0.00 transmitteddrops=0.00 collisions=0.00 transmittederrors=0.00 receiveddrops=0.00", "rcvd_drop": "0", "rcvd_drop@int": "map[value:0]", "rcvd_error": "0", "rcvd_error@int": "map[value:0]", "rcvd_kbit": "1", "rcvd_kbit@int": "map[value:1]", "sent_drop": "0", "sent_drop@int": "map[value:0]", "sent_error": "0", "sent_error@int": "map[value:0]", "sent_kbit": "1", "sent_kbit@int": "map[value:1]", "severity": "Information", "timestamp": "2024-10-02T17:05:55+0200"},
		},
		{
			name:           "Test LM Backup Anonymizer - RAW missing/empty",
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

			maxCapacity := 4000000
			inputReader, err := NewLmBackupReader(inputFile, maxCapacity)
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
