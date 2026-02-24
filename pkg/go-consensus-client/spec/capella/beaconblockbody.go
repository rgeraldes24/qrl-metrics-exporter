// Copyright © 2022, 2024 Attestant Limited.
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
)

// BeaconBlockBody represents the body of a beacon block.
type BeaconBlockBody struct {
	RANDAOReveal      MLDSA87Signature `ssz-size:"4627"`
	ExecutionData     *ExecutionData
	Graffiti          [32]byte               `ssz-size:"32"`
	ProposerSlashings []*ProposerSlashing    `dynssz-max:"MAX_PROPOSER_SLASHINGS" ssz-max:"16"`
	AttesterSlashings []*AttesterSlashing    `dynssz-max:"MAX_ATTESTER_SLASHINGS" ssz-max:"2"`
	Attestations      []*Attestation         `dynssz-max:"MAX_ATTESTATIONS"       ssz-max:"128"`
	Deposits          []*Deposit             `dynssz-max:"MAX_DEPOSITS"           ssz-max:"16"`
	VoluntaryExits    []*SignedVoluntaryExit `dynssz-max:"MAX_VOLUNTARY_EXITS"    ssz-max:"16"`
	SyncAggregate     *SyncAggregate
	ExecutionPayload  *ExecutionPayload
}

// beaconBlockBodyJSON is the spec representation of the struct.
type beaconBlockBodyJSON struct {
	RANDAOReveal      string                 `json:"randao_reveal"`
	ExecutionData     *ExecutionData         `json:"execution_data"`
	Graffiti          string                 `json:"graffiti"`
	ProposerSlashings []*ProposerSlashing    `json:"proposer_slashings"`
	AttesterSlashings []*AttesterSlashing    `json:"attester_slashings"`
	Attestations      []*Attestation         `json:"attestations"`
	Deposits          []*Deposit             `json:"deposits"`
	VoluntaryExits    []*SignedVoluntaryExit `json:"voluntary_exits"`
	SyncAggregate     *SyncAggregate         `json:"sync_aggregate"`
	ExecutionPayload  *ExecutionPayload      `json:"execution_payload"`
}

// beaconBlockBodyYAML is the spec representation of the struct.
type beaconBlockBodyYAML struct {
	RANDAOReveal      string                 `yaml:"randao_reveal"`
	ExecutionData     *ExecutionData         `yaml:"execution_data"`
	Graffiti          string                 `yaml:"graffiti"`
	ProposerSlashings []*ProposerSlashing    `yaml:"proposer_slashings"`
	AttesterSlashings []*AttesterSlashing    `yaml:"attester_slashings"`
	Attestations      []*Attestation         `yaml:"attestations"`
	Deposits          []*Deposit             `yaml:"deposits"`
	VoluntaryExits    []*SignedVoluntaryExit `yaml:"voluntary_exits"`
	SyncAggregate     *SyncAggregate         `yaml:"sync_aggregate"`
	ExecutionPayload  *ExecutionPayload      `yaml:"execution_payload"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockBodyJSON{
		RANDAOReveal:      fmt.Sprintf("%#x", b.RANDAOReveal),
		ExecutionData:     b.ExecutionData,
		Graffiti:          fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings: b.ProposerSlashings,
		AttesterSlashings: b.AttesterSlashings,
		Attestations:      b.Attestations,
		Deposits:          b.Deposits,
		VoluntaryExits:    b.VoluntaryExits,
		SyncAggregate:     b.SyncAggregate,
		ExecutionPayload:  b.ExecutionPayload,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalJSON(input []byte) error {
	var data beaconBlockBodyJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconBlockBody) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconBlockBodyYAML{
		RANDAOReveal:      fmt.Sprintf("%#x", b.RANDAOReveal),
		ExecutionData:     b.ExecutionData,
		Graffiti:          fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings: b.ProposerSlashings,
		AttesterSlashings: b.AttesterSlashings,
		Attestations:      b.Attestations,
		Deposits:          b.Deposits,
		VoluntaryExits:    b.VoluntaryExits,
		SyncAggregate:     b.SyncAggregate,
		ExecutionPayload:  b.ExecutionPayload,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlockBody) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data beaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return b.unpack(&data)
}

// String returns a string version of the structure.
func (b *BeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (b *BeaconBlockBody) unpack(data *beaconBlockBodyJSON) error {
	if data.RANDAOReveal == "" {
		return errors.New("RANDAO reveal missing")
	}

	randaoReveal, err := hex.DecodeString(strings.TrimPrefix(data.RANDAOReveal, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for RANDAO reveal")
	}

	if len(randaoReveal) != SignatureLength {
		return errors.New("incorrect length for RANDAO reveal")
	}

	copy(b.RANDAOReveal[:], randaoReveal)

	if data.ExecutionData == nil {
		return errors.New("Execution data missing")
	}

	b.ExecutionData = data.ExecutionData
	if data.Graffiti == "" {
		return errors.New("graffiti missing")
	}

	graffiti, err := hex.DecodeString(strings.TrimPrefix(data.Graffiti, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for graffiti")
	}

	if len(graffiti) != GraffitiLength {
		return errors.New("incorrect length for graffiti")
	}

	copy(b.Graffiti[:], graffiti)

	if data.ProposerSlashings == nil {
		return errors.New("proposer slashings missing")
	}

	for i := range data.ProposerSlashings {
		if data.ProposerSlashings[i] == nil {
			return fmt.Errorf("proposer slashings entry %d missing", i)
		}
	}

	b.ProposerSlashings = data.ProposerSlashings
	if data.AttesterSlashings == nil {
		return errors.New("attester slashings missing")
	}

	for i := range data.AttesterSlashings {
		if data.AttesterSlashings[i] == nil {
			return fmt.Errorf("attester slashings entry %d missing", i)
		}
	}

	b.AttesterSlashings = data.AttesterSlashings
	if data.Attestations == nil {
		return errors.New("attestations missing")
	}

	for i := range data.Attestations {
		if data.Attestations[i] == nil {
			return fmt.Errorf("attestations entry %d missing", i)
		}
	}

	b.Attestations = data.Attestations
	if data.Deposits == nil {
		return errors.New("deposits missing")
	}

	for i := range data.Deposits {
		if data.Deposits[i] == nil {
			return fmt.Errorf("deposits entry %d missing", i)
		}
	}

	b.Deposits = data.Deposits
	if data.VoluntaryExits == nil {
		return errors.New("voluntary exits missing")
	}

	for i := range data.VoluntaryExits {
		if data.VoluntaryExits[i] == nil {
			return fmt.Errorf("voluntary exits entry %d missing", i)
		}
	}

	b.VoluntaryExits = data.VoluntaryExits
	if data.SyncAggregate == nil {
		return errors.New("sync aggregate missing")
	}

	b.SyncAggregate = data.SyncAggregate
	if data.ExecutionPayload == nil {
		return errors.New("execution payload missing")
	}

	b.ExecutionPayload = data.ExecutionPayload

	return nil
}
