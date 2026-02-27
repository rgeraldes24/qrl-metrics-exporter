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

func TestValidatorJSON(t *testing.T) {
	validPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength)
	validValidator := fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)

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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.validatorJSON",
		},
		{
			name:  "IndexMissing",
			input: fmt.Appendf(nil, `{"balance":"32000000000","status":"active_ongoing","validator":%s}`, validValidator),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: fmt.Appendf(nil, `{"index":true,"balance":"32000000000","status":"active_ongoing","validator":%s}`, validValidator),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: fmt.Appendf(nil, `{"index":"-1","balance":"32000000000","status":"active_ongoing","validator":%s}`, validValidator),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BalanceMissing",
			input: []byte(fmt.Sprintf(`{"index":"1","status":"active_ongoing","validator":%s}`, validValidator)),
			err:   "balance missing",
		},
		{
			name:  "BalanceWrongType",
			input: fmt.Appendf(nil, `{"index":"1","balance":true,"status":"active_ongoing","validator":%s}`, validValidator),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.balance of type string",
		},
		{
			name:  "BalanceInvalid",
			input: fmt.Appendf(nil, `{"index":"1","balance":"-1","status":"active_ongoing","validator":%s}`, validValidator),
			err:   "invalid value for balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "StateWrongType",
			input: fmt.Appendf(nil, `{"index":"1","balance":"32000000000","status":true,"validator":%s}`, validValidator),
			err:   "invalid JSON: unrecognised validator state true",
		},
		{
			name:  "StateInvalid",
			input: fmt.Appendf(nil, `{"index":"1","balance":"32000000000","status":"invalid","validator":%s}`, validValidator),
			err:   "invalid JSON: unrecognised validator state \"invalid\"",
		},
		{
			name:  "ValidatorMissing",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing"}`),
			err:   "validator missing",
		},
		{
			name:  "ValidatorWrongType",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing","validator":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type capella.validatorJSON",
		},
		{
			name:  "ValidatorInvalid",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing","validator":{}}`),
			err:   "invalid JSON: public key missing",
		},
		{
			name:  "Good",
			input: fmt.Appendf(nil, `{"index":"1","balance":"32000000000","status":"active_ongoing","validator":%s}`, validValidator),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Validator
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
