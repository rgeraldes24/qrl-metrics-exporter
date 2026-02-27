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
	require "github.com/stretchr/testify/require"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

func buildTestPubKeys(count int) []string {
	pubkeys := make([]string, count)
	for i := range count {
		pubkeys[i] = "0x" + strings.Repeat(fmt.Sprintf("%02x", i+1), capella.PublicKeyLength)
	}

	return pubkeys
}

func syncCommitteeJSONInput(pubkeys []string) []byte {
	items := make([]string, len(pubkeys))
	for i := range pubkeys {
		items[i] = fmt.Sprintf("\"%s\"", pubkeys[i])
	}

	return []byte(fmt.Sprintf("{\"pubkeys\":[%s]}", strings.Join(items, ",")))
}

func syncCommitteeYAMLInput(pubkeys []string) []byte {
	items := make([]string, len(pubkeys))
	for i := range pubkeys {
		items[i] = fmt.Sprintf("'%s'", pubkeys[i])
	}

	return []byte(fmt.Sprintf("{pubkeys: [%s]}", strings.Join(items, ", ")))
}

func TestSyncCommitteeJSON(t *testing.T) {
	validPubKeys := buildTestPubKeys(16)
	shortPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength-1)
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.syncCommitteeJSON",
		},
		{
			name:  "PubKeysMissing",
			input: []byte(`{"pubkeys":[]}`),
			err:   "public keys missing",
		},
		{
			name:  "InvalidPubKey",
			input: []byte(`{"pubkeys":["invalid"]}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ShortPubKey",
			input: syncCommitteeJSONInput([]string{shortPubKey}),
			err:   "incorrect length for public key",
		},
		{
			name:  "Good",
			input: syncCommitteeJSONInput(validPubKeys),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncCommittee
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

func TestSyncCommitteeYAML(t *testing.T) {
	validPubKeys := buildTestPubKeys(16)
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name:  "Good",
			input: syncCommitteeYAMLInput(validPubKeys),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.SyncCommittee
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
