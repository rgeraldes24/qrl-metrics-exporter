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

func TestValidatorJSON(t *testing.T) {
	validPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength)
	shortPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength-1)
	longPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength+1)

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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.validatorJSON",
		},
		{
			name:  "PublicKeyMissing",
			input: []byte(`{"withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"pubkey":true,"withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"pubkey":"invalid","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, shortPubKey)),
			err:   fmt.Sprintf("incorrect length %d for public key", capella.PublicKeyLength-1),
		},
		{
			name:  "PublicKeyLong",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, longPubKey)),
			err:   fmt.Sprintf("incorrect length %d for public key", capella.PublicKeyLength+1),
		},
		{
			name:  "WithdrawalCredentialsMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "withdrawal credentials missing",
		},
		{
			name:  "WithdrawalCredentialsWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":true,"effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.withdrawal_credentials of type string",
		},
		{
			name:  "WithdrawalCredentialsInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"invalid","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid value for withdrawal credentials: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "WithdrawalCredentialsShort",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0xec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "incorrect length 31 for withdrawal credentials",
		},
		{
			name:  "WithdrawalCredentialsLong",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x0000ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "incorrect length 33 for withdrawal credentials",
		},
		{
			name:  "EffectiveBalanceMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "effective balance missing",
		},
		{
			name:  "EffectiveBalanceWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":true,"slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.effective_balance of type string",
		},
		{
			name:  "EffectiveBalanceInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"-1","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid value for effective balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SlashedWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":"false","activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal string into Go struct field validatorJSON.slashed of type bool",
		},
		{
			name:  "ActivationEligibilityEpochMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "activation eligibility epoch missing",
		},
		{
			name:  "ActivationEligibilityEpochWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":true,"activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.activation_eligibility_epoch of type string",
		},
		{
			name:  "ActivationEligibilityInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"-1","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid value for activation eligibility epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ActivationEpochMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "activation epoch missing",
		},
		{
			name:  "ActivationEligibilityEpochWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":true,"exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.activation_epoch of type string",
		},
		{
			name:  "ActivationInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"-1","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid value for activation epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ExitEpochMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "exit epoch missing",
		},
		{
			name:  "ExitEligibilityEpochWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":true,"withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.exit_epoch of type string",
		},
		{
			name:  "ExitInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"-1","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "invalid value for exit epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "WithdrawableEpochMissing",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615"}`, validPubKey)),
			err:   "withdrawable epoch missing",
		},
		{
			name:  "WithdrawableEligibilityEpochWrongType",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":true}`, validPubKey)),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.withdrawable_epoch of type string",
		},
		{
			name:  "WithdrawableInvalid",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"-1"}`, validPubKey)),
			err:   "invalid value for withdrawable epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(fmt.Sprintf(`{"pubkey":"%s","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}`, validPubKey)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.Validator
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

func TestValidatorYAML(t *testing.T) {
	validPubKey := "0x" + strings.Repeat("11", capella.PublicKeyLength)
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(fmt.Sprintf("{pubkey: '%s', withdrawal_credentials: '0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594', effective_balance: 32000000000, slashed: false, activation_eligibility_epoch: 0, activation_epoch: 0, exit_epoch: 18446744073709551615, withdrawable_epoch: 18446744073709551615}", validPubKey)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.Validator
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
