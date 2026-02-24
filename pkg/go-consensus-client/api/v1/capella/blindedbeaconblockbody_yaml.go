// Copyright © 2022 Attestant Limited.
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
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// blindedBeaconBlockBodyYAML is the spec representation of the struct.
type blindedBeaconBlockBodyYAML struct {
	RANDAOReveal           string                          `yaml:"randao_reveal"`
	ExecutionData          *capella.ExecutionData          `yaml:"execution_data"`
	Graffiti               string                          `yaml:"graffiti"`
	ProposerSlashings      []*capella.ProposerSlashing     `yaml:"proposer_slashings"`
	AttesterSlashings      []*capella.AttesterSlashing     `yaml:"attester_slashings"`
	Attestations           []*capella.Attestation          `yaml:"attestations"`
	Deposits               []*capella.Deposit              `yaml:"deposits"`
	VoluntaryExits         []*capella.SignedVoluntaryExit  `yaml:"voluntary_exits"`
	SyncAggregate          *capella.SyncAggregate          `yaml:"sync_aggregate"`
	ExecutionPayloadHeader *capella.ExecutionPayloadHeader `yaml:"execution_payload_header"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BlindedBeaconBlockBody) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blindedBeaconBlockBodyYAML{
		RANDAOReveal:           fmt.Sprintf("%#x", b.RANDAOReveal),
		ExecutionData:          b.ExecutionData,
		Graffiti:               fmt.Sprintf("%#x", b.Graffiti),
		ProposerSlashings:      b.ProposerSlashings,
		AttesterSlashings:      b.AttesterSlashings,
		Attestations:           b.Attestations,
		Deposits:               b.Deposits,
		VoluntaryExits:         b.VoluntaryExits,
		SyncAggregate:          b.SyncAggregate,
		ExecutionPayloadHeader: b.ExecutionPayloadHeader,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BlindedBeaconBlockBody) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data blindedBeaconBlockBodyJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return b.unpack(&data)
}
