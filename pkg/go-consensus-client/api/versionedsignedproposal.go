// Copyright © 2023, 2024 Attestant Limited.
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

package api

import (
	"errors"
	"math/big"

	apiv1capella "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1/capella"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedSignedProposal contains a versioned signed beacon node proposal.
type VersionedSignedProposal struct {
	Version        spec.DataVersion
	Blinded        bool
	ConsensusValue *big.Int
	ExecutionValue *big.Int
	Capella        *capella.SignedBeaconBlock
	CapellaBlinded *apiv1capella.SignedBlindedBeaconBlock
}

// AssertPresent throws an error if the expected proposal
// given the version and blinded fields is not present.
func (v *VersionedSignedProposal) AssertPresent() error {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil && !v.Blinded {
			return errors.New("capella proposal not present")
		}

		if v.CapellaBlinded == nil && v.Blinded {
			return errors.New("blinded capella proposal not present")
		}
	default:
		return errors.New("unsupported version")
	}

	return nil
}

// Slot returns the slot of the signed proposal.
func (v *VersionedSignedProposal) Slot() (capella.Slot, error) {
	err := v.assertMessagePresent()
	if err != nil {
		return 0, err
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.Slot, nil
		}

		return v.Capella.Message.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the signed proposal.
func (v *VersionedSignedProposal) ProposerIndex() (capella.ValidatorIndex, error) {
	if err := v.assertMessagePresent(); err != nil {
		return 0, err
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.ProposerIndex, nil
		}

		return v.Capella.Message.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ExecutionBlockHash returns the hash of the execution payload.
func (v *VersionedSignedProposal) ExecutionBlockHash() (capella.Hash32, error) {
	if err := v.assertExecutionPayloadPresent(); err != nil {
		return capella.Hash32{}, err
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Message.Body.ExecutionPayloadHeader.BlockHash, nil
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockHash, nil
	default:
		return capella.Hash32{}, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedProposal) String() string {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil {
				return ""
			}

			return v.CapellaBlinded.String()
		}

		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	default:
		return "unsupported version"
	}
}

// assertMessagePresent throws an error if the expected message
// given the version and blinded fields is not present.
//
//nolint:gocyclo // ignore
func (v *VersionedSignedProposal) assertMessagePresent() error {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil ||
				v.CapellaBlinded.Message == nil {
				return ErrDataMissing
			}
		} else {
			if v.Capella == nil ||
				v.Capella.Message == nil {
				return ErrDataMissing
			}
		}
	default:
		return ErrUnsupportedVersion
	}

	return nil
}

// assertExecutionPayloadPresent throws an error if the expected execution payload or payload header
// given the version and blinded fields is not present.
//
//nolint:gocyclo
func (v *VersionedSignedProposal) assertExecutionPayloadPresent() error {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			if v.CapellaBlinded == nil ||
				v.CapellaBlinded.Message == nil ||
				v.CapellaBlinded.Message.Body == nil ||
				v.CapellaBlinded.Message.Body.ExecutionPayloadHeader == nil {
				return ErrDataMissing
			}
		} else {
			if v.Capella == nil ||
				v.Capella.Message == nil ||
				v.Capella.Message.Body == nil ||
				v.Capella.Message.Body.ExecutionPayload == nil {
				return ErrDataMissing
			}
		}
	default:
		return ErrUnsupportedVersion
	}

	return nil
}
