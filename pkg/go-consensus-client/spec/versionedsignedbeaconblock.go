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

// VersionedSignedBeaconBlock contains a versioned signed beacon block.
type VersionedSignedBeaconBlock struct {
	Version DataVersion
	Capella *capella.SignedBeaconBlock
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedSignedBeaconBlock) Slot() (capella.Slot, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedSignedBeaconBlock) ProposerIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionBlockHash returns the block hash of the beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionBlockHash() (capella.Hash32, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayload == nil {
			return capella.Hash32{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockHash, nil
	default:
		return capella.Hash32{}, errors.New("unknown version")
	}
}

// ExecutionBlockNumber returns the block number of the beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionBlockNumber() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayload == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionTransactions returns the execution payload transactions for the block.
func (v *VersionedSignedBeaconBlock) ExecutionTransactions() ([]capella.Transaction, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Graffiti returns the graffiti for the block.
func (v *VersionedSignedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return [32]byte{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.Graffiti, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
//
//nolint:gocyclo
func (v *VersionedSignedBeaconBlock) Attestations() ([]*VersionedAttestation, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedAttestations := make([]*VersionedAttestation, len(v.Capella.Message.Body.Attestations))
		for i, attestation := range v.Capella.Message.Body.Attestations {
			versionedAttestations[i] = &VersionedAttestation{
				Version: DataVersionCapella,
				Capella: attestation,
			}
		}

		return versionedAttestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedSignedBeaconBlock) Root() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.HashTreeRoot()
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedSignedBeaconBlock) BodyRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.HashTreeRoot()
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedSignedBeaconBlock) ParentRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.ParentRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedSignedBeaconBlock) StateRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return capella.Root{}, errors.New("no capella block")
		}

		return v.Capella.Message.StateRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// RandaoReveal returns the randao reveal of the beacon block.
func (v *VersionedSignedBeaconBlock) RandaoReveal() (capella.BLSSignature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return capella.BLSSignature{}, errors.New("no capella block")
		}

		return v.Capella.Message.Body.RANDAOReveal, nil
	default:
		return capella.BLSSignature{}, errors.New("unknown version")
	}
}

// ExecutionData returns the execution data of the beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionData() (*capella.ExecutionData, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionData, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Deposits returns the deposits of the beacon block.
func (v *VersionedSignedBeaconBlock) Deposits() ([]*capella.Deposit, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.Deposits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// VoluntaryExits returns the voluntary exits of the beacon block.
func (v *VersionedSignedBeaconBlock) VoluntaryExits() ([]*capella.SignedVoluntaryExit, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.VoluntaryExits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
//
//nolint:gocyclo
func (v *VersionedSignedBeaconBlock) AttesterSlashings() ([]VersionedAttesterSlashing, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedAttesterSlashings := make([]VersionedAttesterSlashing, len(v.Capella.Message.Body.AttesterSlashings))
		for i, attesterSlashing := range v.Capella.Message.Body.AttesterSlashings {
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
func (v *VersionedSignedBeaconBlock) ProposerSlashings() ([]*capella.ProposerSlashing, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// SyncAggregate returns the sync aggregate of the beacon block.
func (v *VersionedSignedBeaconBlock) SyncAggregate() (*capella.SyncAggregate, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.SyncAggregate, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals of the beacon block.
func (v *VersionedSignedBeaconBlock) Withdrawals() ([]*capella.Withdrawal, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Message == nil ||
			v.Capella.Message.Body == nil ||
			v.Capella.Message.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Message.Body.ExecutionPayload.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ExecutionPayload returns the execution payload of the signed beacon block.
func (v *VersionedSignedBeaconBlock) ExecutionPayload() (*VersionedExecutionPayload, error) {
	versionedExecutionPayload := &VersionedExecutionPayload{
		Version: v.Version,
	}

	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}

		versionedExecutionPayload.Capella = v.Capella.Message.Body.ExecutionPayload
	default:
		return nil, errors.New("unknown version")
	}

	return versionedExecutionPayload, nil
}

// String returns a string version of the structure.
func (v *VersionedSignedBeaconBlock) String() string {
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
