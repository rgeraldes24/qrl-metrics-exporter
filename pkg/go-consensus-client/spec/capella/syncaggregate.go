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

// SyncAggregate is the QRL sync aggregate structure.
type SyncAggregate struct {
	SyncCommitteeBits       bitfield.Bitvector16 `dynssz-size:"SYNC_COMMITTEE_SIZE/8" ssz-size:"2"`
	SyncCommitteeSignatures []MLDSA87Signature   `ssz-max:"16" ssz-size:"?,4627"`
}

// syncAggregateJSON is the spec representation of the struct.
type syncAggregateJSON struct {
	SyncCommitteeBits       string   `json:"sync_committee_bits"`
	SyncCommitteeSignatures []string `json:"sync_committee_signatures"`
}

// syncAggregateYAML is the spec representation of the struct.
type syncAggregateYAML struct {
	SyncCommitteeBits       string   `yaml:"sync_committee_bits"`
	SyncCommitteeSignatures []string `yaml:"sync_committee_signatures"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncAggregate) MarshalJSON() ([]byte, error) {
	signatures := make([]string, len(s.SyncCommitteeSignatures))
	for i := range s.SyncCommitteeSignatures {
		signatures[i] = fmt.Sprintf("%#x", s.SyncCommitteeSignatures[i])
	}

	return json.Marshal(&syncAggregateJSON{
		SyncCommitteeBits:       fmt.Sprintf("%#x", s.SyncCommitteeBits.Bytes()),
		SyncCommitteeSignatures: signatures,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncAggregate) UnmarshalJSON(input []byte) error {
	var syncAggregateJSON syncAggregateJSON
	if err := json.Unmarshal(input, &syncAggregateJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&syncAggregateJSON)
}

func (s *SyncAggregate) unpack(syncAggregateJSON *syncAggregateJSON) error {
	if syncAggregateJSON.SyncCommitteeBits == "" {
		return errors.New("sync committee bits missing")
	}

	syncCommitteeBits, err := hex.DecodeString(strings.TrimPrefix(syncAggregateJSON.SyncCommitteeBits, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for sync committee bits")
	}
	if len(syncCommitteeBits) != 2 {
		return errors.New("incorrect length for sync committee bits")
	}

	s.SyncCommitteeBits = syncCommitteeBits

	if len(syncAggregateJSON.SyncCommitteeSignatures) == 0 {
		return errors.New("sync committee signatures missing")
	}

	s.SyncCommitteeSignatures = make([]MLDSA87Signature, len(syncAggregateJSON.SyncCommitteeSignatures))
	for i := range syncAggregateJSON.SyncCommitteeSignatures {
		syncCommitteeSignature, err := hex.DecodeString(strings.TrimPrefix(syncAggregateJSON.SyncCommitteeSignatures[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for sync committee signature")
		}

		if len(syncCommitteeSignature) != SignatureLength {
			return errors.New("incorrect length for sync committee signature")
		}

		copy(s.SyncCommitteeSignatures[i][:], syncCommitteeSignature)
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SyncAggregate) MarshalYAML() ([]byte, error) {
	signatures := make([]string, len(s.SyncCommitteeSignatures))
	for i := range s.SyncCommitteeSignatures {
		signatures[i] = fmt.Sprintf("%#x", s.SyncCommitteeSignatures[i])
	}

	yamlBytes, err := yaml.MarshalWithOptions(&syncAggregateYAML{
		SyncCommitteeBits:       fmt.Sprintf("%#x", s.SyncCommitteeBits.Bytes()),
		SyncCommitteeSignatures: signatures,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SyncAggregate) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncAggregateJSON syncAggregateJSON
	if err := yaml.Unmarshal(input, &syncAggregateJSON); err != nil {
		return err
	}

	return s.unpack(&syncAggregateJSON)
}

// String returns a string version of the structure.
func (s *SyncAggregate) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
