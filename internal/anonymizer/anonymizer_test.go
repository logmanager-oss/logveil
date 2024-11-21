package anonymizer

import (
	"math/rand"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/proof"
	"github.com/stretchr/testify/assert"
)

func TestAnonimizer_AnonymizeData(t *testing.T) {
	tests := []struct {
		name               string
		anonymizingDataDir string
		input              map[string]string
		expectedOutput     string
	}{
		{
			name:               "Test AnonymizeData",
			anonymizingDataDir: "../../tests/data/anonymization_data",
			input: map[string]string{
				"@timestamp":   "2024-06-05T14:59:27.000+00:00",
				"src_ip":       "10.10.10.1",
				"src_ipv6":     "7f1d:64ed:536a:1fd7:fe8e:cc29:9df4:7911",
				"mac":          "71:e5:41:18:cb:3e",
				"email":        "test@test.com",
				"url":          "https://www.testurl.com",
				"username":     "miloslav.illes",
				"organization": "Microsoft",
				"raw":          "2024-06-05T14:59:27.000+00:00, 10.10.10.1, 7f1d:64ed:536a:1fd7:fe8e:cc29:9df4:7911, miloslav.illes, Microsoft, 71:e5:41:18:cb:3e, test@test.com, https://www.testurl.com",
			},
			expectedOutput: "2024-06-05T14:59:27.000+00:00, 10.20.0.53, 8186:39ac:48a4:c6af:a2f1:581a:8b95:25e2, ladislav.dosek, Apple, 0f:da:68:92:7f:2b, QHtPwsw@RJSkoHl.top, http://soqovkq.com/NfkcUjG.php",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonymizer, err := CreateAnonymizer(&config.Config{AnonymizationDataPath: tt.anonymizingDataDir}, &proof.ProofWriter{IsEnabled: false})
			if err != nil {
				t.Fatal(err)
			}
			// Disabling randomization so we know which values to expect
			anonymizer.SetRandFunc(func(int) int { return 1 })
			faker.SetRandomSource(rand.NewSource(1))
			output := anonymizer.Anonymize(tt.input)

			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
