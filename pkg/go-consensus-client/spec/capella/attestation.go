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

package capella

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	bitfield "github.com/theQRL/go-bitfield"
)

// Attestation is the QRL attestation structure.
type Attestation struct {
	AggregationBits bitfield.Bitlist `dynssz-max:"MAX_VALIDATORS_PER_COMMITTEE" ssz-max:"128"`
	Data            *AttestationData
	Signatures      []MLDSA87Signature `ssz-max:"128" ssz-size:"?,4627"`
}

// attestationJSON is a raw representation of the struct.
type attestationJSON struct {
	AggregationBits string           `json:"aggregation_bits"`
	Data            *AttestationData `json:"data"`
	Signatures      []string         `json:"signatures"`
}

// attestationYAML is a raw representation of the struct.
type attestationYAML struct {
	AggregationBits string           `yaml:"aggregation_bits"`
	Data            *AttestationData `yaml:"data"`
	Signatures      []string         `yaml:"signatures"`
}

// MarshalJSON implements json.Marshaler.
func (a *Attestation) MarshalJSON() ([]byte, error) {
	signatures := make([]string, len(a.Signatures))
	for i := range a.Signatures {
		signatures[i] = fmt.Sprintf("%#x", a.Signatures[i])
	}

	return json.Marshal(&attestationJSON{
		AggregationBits: fmt.Sprintf("%#x", []byte(a.AggregationBits)),
		Data:            a.Data,
		Signatures:      signatures,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attestation) UnmarshalJSON(input []byte) error {
	var attestationJSON attestationJSON

	err := json.Unmarshal(input, &attestationJSON)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&attestationJSON)
}

// MarshalYAML implements yaml.Marshaler.
func (a *Attestation) MarshalYAML() ([]byte, error) {
	signatures := make([]string, len(a.Signatures))
	for i := range a.Signatures {
		signatures[i] = fmt.Sprintf("%#x", a.Signatures[i])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&attestationYAML{
		AggregationBits: fmt.Sprintf("%#x", []byte(a.AggregationBits)),
		Data:            a.Data,
		Signatures:      signatures,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *Attestation) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var attestationJSON attestationJSON
	if err := yaml.Unmarshal(input, &attestationJSON); err != nil {
		return err
	}

	return a.unpack(&attestationJSON)
}

// String returns a string version of the structure.
func (a *Attestation) String() string {
	data, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (a *Attestation) unpack(attestationJSON *attestationJSON) error {
	var err error

	if attestationJSON.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}

	if a.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(attestationJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}

	a.Data = attestationJSON.Data
	if a.Data == nil {
		return errors.New("data missing")
	}

	if len(attestationJSON.Signatures) == 0 {
		return errors.New("signatures missing")
	}

	a.Signatures = make([]MLDSA87Signature, len(attestationJSON.Signatures))
	for i := range attestationJSON.Signatures {
		signature, err := hex.DecodeString(strings.TrimPrefix(attestationJSON.Signatures[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for signature")
		}

		if len(signature) != SignatureLength {
			return errors.New("incorrect length for signature")
		}

		copy(a.Signatures[i][:], signature)
	}

	return nil
}
