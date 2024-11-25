package anonymizer

import (
	"net"
	"net/mail"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/logmanager-oss/logveil/internal/config"
	"github.com/logmanager-oss/logveil/internal/proof"
)

func TestAnonimizer_AnonymizeData(t *testing.T) {
	tests := []struct {
		name                           string
		anonymizationDataDir           string
		customAnonymizationMappingPath string
		input                          map[string]string
		expectedOutput                 string
	}{
		{
			name:                           "Test AnonymizeData",
			anonymizationDataDir:           "../../tests/data/anonymization_data",
			customAnonymizationMappingPath: "../../tests/data/custom_mappings.txt",
			input: map[string]string{
				"@timestamp":   "2024-06-05T14:59:27.000+00:00",
				"src_ip":       "10.10.10.1",
				"src_ipv6":     "7f1d:64ed:536a:1fd7:fe8e:cc29:9df4:7911",
				"mac":          "71:e5:41:18:cb:3e",
				"email":        "atest@atest.com",
				"url":          "https://www.testurl.com",
				"username":     "miloslav.illes",
				"organization": "Microsoft",
				"custom:":      "replacement_test",
				"raw":          "2024-06-05T14:59:27.000+00:00, 10.10.10.1, 7f1d:64ed:536a:1fd7:fe8e:cc29:9df4:7911, miloslav.illes, Microsoft, 71:e5:41:18:cb:3e, test@test.com, https://www.testurl.com, replace_this",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anonymizer, err := CreateAnonymizer(
				&config.Config{
					AnonymizationDataPath:    tt.anonymizationDataDir,
					CustomReplacementMapPath: tt.customAnonymizationMappingPath,
				},
				&proof.ProofWriter{IsEnabled: false},
			)
			if err != nil {
				t.Fatal(err)
			}
			output := anonymizer.Anonymize(tt.input)

			// Verify each part of the output individually - Is generated value valid in terms of its type and not the same as input?
			parts := strings.Split(output, ", ")

			ipv4 := net.ParseIP(parts[1])
			if ipv4 == nil || ipv4.String() == tt.input["src_ip"] {
				t.Fatalf("invalid IPv4 generated or it didn't got replaced at all: %s", parts[1])
			}

			ipv6 := net.ParseIP(parts[2])
			if ipv6 == nil || ipv6.String() == tt.input["src_ip"] {
				t.Fatalf("invalid IPv6 generated or it didn't got replaced at all: %s", parts[2])
			}

			if !slices.Contains(anonymizer.anonymizationData["username"], parts[3]) || parts[3] == tt.input["username"] {
				t.Fatalf("invalid username or it didn't got replaced at all: %s", parts[3])
			}

			if !slices.Contains(anonymizer.anonymizationData["organization"], parts[4]) || parts[4] == tt.input["organization"] {
				t.Fatalf("invalid organization or it didn't got replaced at all: %s", parts[4])
			}

			mac, err := net.ParseMAC(parts[5])
			if err != nil {
				t.Fatalf("invalid MAC generated: %s", parts[5])
			}
			if mac.String() == tt.input["mac"] {
				t.Fatalf("MAC not replaced at all")
			}

			email, err := mail.ParseAddress(parts[6])
			if err != nil {
				t.Fatalf("invalid email generated: %s", parts[6])
			}
			if email.Address == tt.input["email"] {
				t.Fatalf("email not replaced at all")
			}

			url, err := url.ParseRequestURI(parts[7])
			if err != nil {
				t.Fatalf("invalid url generated: %s", parts[7])
			}
			if url.String() == tt.input["url"] {
				t.Fatalf("url not replaced at all")
			}

			if parts[8] != anonymizer.replacementMap["replace_this"] {
				t.Fatalf("custom replacement didn't work")
			}
		})
	}
}
