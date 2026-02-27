// Copyright © 2020 Attestant Limited.
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

package v1_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	api "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

func TestBeaconBlockHeaderJSON(t *testing.T) {
	validSignature := "0x" + strings.Repeat("61", capella.SignatureLength)

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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.beaconBlockHeaderJSON",
		},
		{
			name:  "RootMissing",
			input: []byte(fmt.Sprintf(`{"canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "root missing",
		},
		{
			name:  "RootWrongType",
			input: []byte(fmt.Sprintf(`{"root":true,"canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.root of type string",
		},
		{
			name:  "RootInvalid",
			input: []byte(fmt.Sprintf(`{"root":"invalid","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "invalid value for root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "RootShort",
			input: []byte(fmt.Sprintf(`{"root":"0x354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "incorrect length 31 for root",
		},
		{
			name:  "RootLong",
			input: []byte(fmt.Sprintf(`{"root":"0xbcbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "incorrect length 33 for root",
		},
		{
			name:  "CanonicalWrongType",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":"true","header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal string into Go struct field beaconBlockHeaderJSON.canonical of type bool",
		},
		{
			name:  "CanonicalInvalid",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":maybe,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
			err:   "invalid character 'm' looking for beginning of value",
		},
		{
			name:  "HeaderMissing",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"signature":"%s"}`, validSignature)),
			err:   "header missing",
		},
		{
			name:  "HeaderWrongType",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":true,"signature":"%s"}`, validSignature)),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type capella.signedBeaconBlockHeaderJSON",
		},
		{
			name:  "HeaderInvalid",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"signature":"%s"}}`, validSignature)),
			err:   "invalid JSON: message missing",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"%s"}}`, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BeaconBlockHeader
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
