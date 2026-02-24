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

// VersionedBlindedBeaconBlock contains a versioned blinded beacon block.
type VersionedBlindedBeaconBlock struct {
	Version spec.DataVersion
	Capella *apiv1capella.BlindedBeaconBlock
}

// IsEmpty returns true if there is no block.
func (v *VersionedBlindedBeaconBlock) IsEmpty() bool {
	return v.Capella == nil
}

// Slot returns the slot of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Slot() (capella.Slot, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedBlindedBeaconBlock) ProposerIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) RandaoReveal() (capella.MLDSA87Signature, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return capella.MLDSA87Signature{}, ErrDataMissing
		}

		return v.Capella.Body.RANDAOReveal, nil
	default:
		return capella.MLDSA87Signature{}, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Capella.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Attestations() ([]spec.VersionedAttestation, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Capella.Body.Attestations))
		for i, attestation := range v.Capella.Body.Attestations {
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

// Root returns the root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Root() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) BodyRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Body.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) ParentRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.ParentRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) StateRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.StateRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// TransactionsRoot returns the transactions root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) TransactionsRoot() (capella.Root, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) FeeRecipient() (capella.ExecutionAddress, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return capella.ExecutionAddress{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.FeeRecipient, nil
	default:
		return capella.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Timestamp() (uint64, error) {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedBlindedBeaconBlock) String() string {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	default:
		return "unknown version"
	}
}
