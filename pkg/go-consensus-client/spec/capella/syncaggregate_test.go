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

func TestSyncAggregateJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.syncAggregateJSON",
		},
		{
			name:  "SyncCommitteeBitsMissing",
			input: []byte(fmt.Sprintf(`{"sync_committee_signatures":["%s"]}`, validSignature)),
			err:   "sync committee bits missing",
		},
		{
			name:  "SyncCommitteeBitsInvalid",
			input: []byte(fmt.Sprintf(`{"sync_committee_bits":"invalid","sync_committee_signatures":["%s"]}`, validSignature)),
			err:   "invalid value for sync committee bits: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SyncCommitteeSignatureMissing",
			input: []byte(`{"sync_committee_bits":"0xe7fc"}`),
			err:   "sync committee signatures missing",
		},
		{
			name:  "SyncCommitteeSignaturesWrongType",
			input: []byte(`{"sync_committee_bits":"0xe7fc","sync_committee_signatures":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncAggregateJSON.sync_committee_signatures of type []string",
		},
		{
			name:  "SyncCommitteeSignatureInvalid",
			input: []byte(`{"sync_committee_bits":"0xe7fc","sync_committee_signatures":["invalid"]}`),
			err:   "invalid value for sync committee signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(fmt.Sprintf(`{"sync_committee_bits":"0xe7fc","sync_committee_signatures":["%s"]}`, shortSignature)),
			err:   "incorrect length for sync committee signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(fmt.Sprintf(`{"sync_committee_bits":"0xe7fc","sync_committee_signatures":["%s"]}`, longSignature)),
			err:   "incorrect length for sync committee signature",
		},
		{
			name:  "SyncCommitteeBitsWrongLength",
			input: []byte(fmt.Sprintf(`{"sync_committee_bits":"0xe7fcbc","sync_committee_signatures":["%s"]}`, validSignature)),
			err:   "incorrect length for sync committee bits",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{"sync_committee_bits":"0xe7fc","sync_committee_signatures":["%s"]}`, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncAggregate
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

func TestSyncAggregateYAML(t *testing.T) {
	validSignature := "0x" + strings.Repeat("61", capella.SignatureLength)

	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{sync_committee_bits: '0xe7fc', sync_committee_signatures: ['%s']}`, validSignature)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncAggregate
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
