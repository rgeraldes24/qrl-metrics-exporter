// Copyright © 2020, 2021 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package capella_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

const (
	indexedAttestationDataJSON = `{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`
	indexedAttestationDataYAML = `{slot: 100, index: 1, beacon_block_root: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f', source: {epoch: 1, root: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f'}, target: {epoch: 2, root: '0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f'}}`
)

func TestIndexedAttestationJSON(t *testing.T) {
	validSignature := "0x" + strings.Repeat("61", capella.SignatureLength)
	shortSignature := "0x" + strings.Repeat("61", capella.SignatureLength-1)
	longSignature := "0x" + strings.Repeat("61", capella.SignatureLength+1)

	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.indexedAttestationJSON",
		},
		// Spec tests contain indexed attestations without attesting indices.
		// {
		// 	name:  "AttestingIndicesMissing",
		// 	input: []byte(`{"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signatures":["0x..."]}`),
		// 	err:   "attesting indices missing",
		// },
		// {
		// 	name:  "AttestingIndicesEmpty",
		// 	input: []byte(`{"attesting_indices":[],"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signatures":["0x..."]}`),
		// 	err:   "attesting indices missing",
		// },
		{
			name:  "AttestingIndicesWrongType",
			input: []byte(fmt.Sprintf(`{"attesting_indices":true,"data":%s,"signatures":["%s"]}`, indexedAttestationDataJSON, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field indexedAttestationJSON.attesting_indices of type []string",
		},
		{
			name:  "AttestingIndicesInvalid",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["-1","2","3"],"data":%s,"signatures":["%s"]}`, indexedAttestationDataJSON, validSignature)),
			err:   "failed to parse attesting index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "DataMissing",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"signatures":["%s"]}`, validSignature)),
			err:   "data missing",
		},
		{
			name:  "DataWrongType",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":true,"signatures":["%s"]}`, validSignature)),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type capella.attestationDataJSON",
		},
		{
			name:  "DataInvalid",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":{},"signatures":["%s"]}`, validSignature)),
			err:   "invalid JSON: slot missing",
		},
		{
			name:  "SignatureMissing",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s}`, indexedAttestationDataJSON)),
			err:   "signatures missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s,"signatures":true}`, indexedAttestationDataJSON)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field indexedAttestationJSON.signatures of type []string",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s,"signatures":["invalid"]}`, indexedAttestationDataJSON)),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s,"signatures":["%s"]}`, indexedAttestationDataJSON, shortSignature)),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s,"signatures":["%s"]}`, indexedAttestationDataJSON, longSignature)),
			err:   "incorrect length for signature",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{"attesting_indices":["1","2","3"],"data":%s,"signatures":["%s"]}`, indexedAttestationDataJSON, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.IndexedAttestation
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestIndexedAttestationYAML(t *testing.T) {
	validSignature := "0x" + strings.Repeat("61", capella.SignatureLength)

	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{attesting_indices: [1, 2, 3], data: %s, signatures: ['%s']}`, indexedAttestationDataYAML, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.IndexedAttestation
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(rt), res.String())
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
