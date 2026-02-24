// Copyright © 2021 Attestant Limited.
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

const syncCommitteeContributionBeaconBlockRoot = "0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c"

func TestSyncCommitteeContributionJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.syncCommitteeContributionJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(fmt.Sprintf(`{"beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(fmt.Sprintf(`{"slot":true,"beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(fmt.Sprintf(`{"slot":"-1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BeaconBlockRootMissing",
			input: []byte(fmt.Sprintf(`{"slot":"1","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, validSignature)),
			err:   "beacon block root missing",
		},
		{
			name:  "BeaconBlockRootWrongType",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":true,"subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.beacon_block_root of type string",
		},
		{
			name:  "BeaconBlockRootInvalid",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"invalid","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, validSignature)),
			err:   "invalid value for beacon block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BeaconBlockRootShort",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d06448","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, validSignature)),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "SubcommitteeIndexMissing",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "subcommittee index missing",
		},
		{
			name:  "SubcommitteeIndexWrongType",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":true,"aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.subcommittee_index of type string",
		},
		{
			name:  "SubcommitteeIndexInvalid",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"-3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid value for subcommittee index: strconv.ParseUint: parsing \"-3\": invalid syntax",
		},
		{
			name:  "AggregationBitsMissing",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "aggregation bits missing",
		},
		{
			name:  "AggregationBitsWrongType",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":true,"signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.aggregation_bits of type string",
		},
		{
			name:  "AggregationBitsInvalid",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"invalid","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "invalid value for aggregation bits: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "AggregationBitsWrongLength",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x000400","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
			err:   "incorrect length for aggregation bits",
		},
		{
			name:  "SignatureMissing",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004"}`, syncCommitteeContributionBeaconBlockRoot)),
			err:   "signatures missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":true}`, syncCommitteeContributionBeaconBlockRoot)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeContributionJSON.signatures of type []string",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["invalid"]}`, syncCommitteeContributionBeaconBlockRoot)),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, shortSignature)),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, longSignature)),
			err:   "incorrect length for signature",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{"slot":"1","beacon_block_root":"%s","subcommittee_index":"3","aggregation_bits":"0x0004","signatures":["%s"]}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncCommitteeContribution
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

func TestSyncCommitteeContributionYAML(t *testing.T) {
	validSignature := "0x" + strings.Repeat("61", capella.SignatureLength)

	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{slot: 1, beacon_block_root: '%s', subcommittee_index: 3, aggregation_bits: '0x0004', signatures: ['%s']}`, syncCommitteeContributionBeaconBlockRoot, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncCommitteeContribution
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
