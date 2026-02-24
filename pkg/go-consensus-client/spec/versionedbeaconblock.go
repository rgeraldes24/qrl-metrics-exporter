// Copyright © 2021 - 2024 Attestant Limited.
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

package spec

import (
	"errors"

	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedBeaconBlock contains a versioned beacon block.
type VersionedBeaconBlock struct {
	Version DataVersion
	Capella *capella.BeaconBlock
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconBlock) IsEmpty() bool {
	return v.Capella == nil
}

// Slot returns the slot of the beacon block.
func (v *VersionedBeaconBlock) Slot() (capella.Slot, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// RandaoReveal returns the RANDAO reveal of the beacon block.
func (v *VersionedBeaconBlock) RandaoReveal() (capella.BLSSignature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.BLSSignature{}, errors.New("no capella block")
		}

		if v.Capella.Body == nil {
			return capella.BLSSignature{}, errors.New("no capella block body")
		}

		return v.Capella.Body.RANDAOReveal, nil
	default:
		return capella.BLSSignature{}, errors.New("unknown version")
	}
}

// Graffiti returns the graffiti of the beacon block.
func (v *VersionedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella block")
		}

		if v.Capella.Body == nil {
			return [32]byte{}, errors.New("no capella block body")
		}

		return v.Capella.Body.Graffiti, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedBeaconBlock) ProposerIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedBeaconBlock) Root() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.HashTreeRoot()
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedBeaconBlock) BodyRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		if v.Capella.Body == nil {
			return capella.Root{}, errors.New("no capella block body")
		}

		return v.Capella.Body.HashTreeRoot()
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedBeaconBlock) ParentRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.ParentRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedBeaconBlock) StateRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.StateRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedBeaconBlock) Attestations() ([]VersionedAttestation, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedAttestations := make([]VersionedAttestation, len(v.Capella.Body.Attestations))
		for i, attestation := range v.Capella.Body.Attestations {
			versionedAttestations[i] = VersionedAttestation{
				Version: DataVersionCapella,
				Capella: attestation,
			}
		}

		return versionedAttestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedBeaconBlock) AttesterSlashings() ([]VersionedAttesterSlashing, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedAttesterSlashings := make([]VersionedAttesterSlashing, len(v.Capella.Body.AttesterSlashings))
		for i, attesterSlashing := range v.Capella.Body.AttesterSlashings {
			versionedAttesterSlashings[i] = VersionedAttesterSlashing{
				Version: DataVersionCapella,
				Capella: attesterSlashing,
			}
		}

		return versionedAttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedBeaconBlock) ProposerSlashings() ([]*capella.ProposerSlashing, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ExecutionPayload returns the execution payload of the beacon block.
func (v *VersionedBeaconBlock) ExecutionPayload() (*VersionedExecutionPayload, error) {
	versionedExecutionPayload := &VersionedExecutionPayload{
		Version: v.Version,
	}

	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedExecutionPayload.Capella = v.Capella.Body.ExecutionPayload
	default:
		return nil, errors.New("unknown version")
	}

	return versionedExecutionPayload, nil
}

// String returns a string version of the structure.
func (v *VersionedBeaconBlock) String() string {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	default:
		return "unknown version"
	}
}
