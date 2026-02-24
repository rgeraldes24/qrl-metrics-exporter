// Copyright © 2022, 2023 Attestant Limited.
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
	apiv1capella "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1/capella"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedSignedBlindedBeaconBlock contains a versioned signed blinded beacon block.
type VersionedSignedBlindedBeaconBlock struct {
	Version spec.DataVersion
	Capella *apiv1capella.SignedBlindedBeaconBlock
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Slot() (capella.Slot, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Message.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the signed blinded beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Attestations() ([]spec.VersionedAttestation, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Capella.Message.Body.Attestations))
		for i, attestation := range v.Capella.Message.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionCapella,
				Capella: attestation,
			}
		}

		return versionedAttestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Root returns the root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Root() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Message.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) BodyRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Message.Body.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ParentRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Message.ParentRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) StateRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Message.StateRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) AttesterSlashings() ([]spec.VersionedAttesterSlashing, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttesterSlashings := make([]spec.VersionedAttesterSlashing, len(v.Capella.Message.Body.AttesterSlashings))
		for i, attesterSlashing := range v.Capella.Message.Body.AttesterSlashings {
			versionedAttesterSlashings[i] = spec.VersionedAttesterSlashing{
				Version: spec.DataVersionCapella,
				Capella: attesterSlashing,
			}
		}

		return versionedAttesterSlashings, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ProposerSlashings() ([]*capella.ProposerSlashing, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Capella.Message.Body.ProposerSlashings, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ProposerIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Message.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ExecutionParentHash returns the parent hash of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ExecutionParentHash() (capella.Hash32, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayloadHeader == nil {
			return capella.Hash32{}, ErrDataMissing
		}

		return v.Capella.Message.Body.ExecutionPayloadHeader.ParentHash, nil
	default:
		return capella.Hash32{}, ErrUnsupportedVersion
	}
}

// ExecutionBlockHash returns the hash of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ExecutionBlockHash() (capella.Hash32, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayloadHeader == nil {
			return capella.Hash32{}, ErrDataMissing
		}

		return v.Capella.Message.Body.ExecutionPayloadHeader.BlockHash, nil
	default:
		return capella.Hash32{}, ErrUnsupportedVersion
	}
}

// ExecutionBlockNumber returns the block number of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ExecutionBlockNumber() (uint64, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Message.Body.ExecutionPayloadHeader.BlockNumber, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Signature returns the signature of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Signature() (capella.MLDSA87Signature, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return capella.MLDSA87Signature{}, ErrDataMissing
		}

		return v.Capella.Signature, nil
	default:
		return capella.MLDSA87Signature{}, ErrUnsupportedVersion
	}
}
