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
	"math/big"

	apiv1capella "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1/capella"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedProposal contains a versioned proposal.
type VersionedProposal struct {
	Version        spec.DataVersion
	Blinded        bool
	ConsensusValue *big.Int
	ExecutionValue *big.Int
	Capella        *capella.BeaconBlock
	CapellaBlinded *apiv1capella.BlindedBeaconBlock
}

// IsEmpty returns true if there is no proposal.
func (v *VersionedProposal) IsEmpty() bool {
	return v.Capella == nil &&
		v.CapellaBlinded == nil
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedProposal) BodyRoot() (capella.Root, error) {
	if !v.bodyPresent() {
		return capella.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.HashTreeRoot()
		}

		return v.Capella.Body.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedProposal) ParentRoot() (capella.Root, error) {
	if !v.proposalPresent() {
		return capella.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.ParentRoot, nil
		}

		return v.Capella.ParentRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedProposal) ProposerIndex() (capella.ValidatorIndex, error) {
	if !v.proposalPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.ProposerIndex, nil
		}

		return v.Capella.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Root returns the root of the proposal.
func (v *VersionedProposal) Root() (capella.Root, error) {
	if !v.proposalPresent() {
		return capella.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.HashTreeRoot()
		}

		return v.Capella.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// Slot returns the slot of the proposal.
func (v *VersionedProposal) Slot() (capella.Slot, error) {
	if !v.proposalPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Slot, nil
		}

		return v.Capella.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the proposal.
func (v *VersionedProposal) StateRoot() (capella.Root, error) {
	if !v.proposalPresent() {
		return capella.Root{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.StateRoot, nil
		}

		return v.Capella.StateRoot, nil
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the proposal.
func (v *VersionedProposal) Attestations() ([]spec.VersionedAttestation, error) {
	if !v.bodyPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			versionedAttestations := make([]spec.VersionedAttestation, len(v.CapellaBlinded.Body.Attestations))
			for i, attestation := range v.CapellaBlinded.Body.Attestations {
				versionedAttestations[i] = spec.VersionedAttestation{
					Version: spec.DataVersionCapella,
					Capella: attestation,
				}
			}

			return versionedAttestations, nil
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

// Graffiti returns the graffiti of the proposal.
func (v *VersionedProposal) Graffiti() ([32]byte, error) {
	if !v.bodyPresent() {
		return [32]byte{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.Graffiti, nil
		}

		return v.Capella.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedProposal) RandaoReveal() (capella.BLSSignature, error) {
	if !v.bodyPresent() {
		return capella.BLSSignature{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.RANDAOReveal, nil
		}

		return v.Capella.Body.RANDAOReveal, nil
	default:
		return capella.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Transactions returns the transactions of the proposal.
func (v *VersionedProposal) Transactions() ([]capella.Transaction, error) {
	if !v.payloadPresent() {
		return nil, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedProposal) FeeRecipient() (capella.ExecutionAddress, error) {
	if !v.payloadPresent() {
		return capella.ExecutionAddress{}, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.FeeRecipient, nil
		}

		return v.Capella.Body.ExecutionPayload.FeeRecipient, nil
	default:
		return capella.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedProposal) Timestamp() (uint64, error) {
	if !v.payloadPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.Timestamp, nil
		}

		return v.Capella.Body.ExecutionPayload.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// GasLimit returns the gas limit of the proposal.
func (v *VersionedProposal) GasLimit() (uint64, error) {
	if !v.payloadPresent() {
		return 0, ErrDataMissing
	}

	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded.Body.ExecutionPayloadHeader.GasLimit, nil
		}

		return v.Capella.Body.ExecutionPayload.GasLimit, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Value returns the value of the proposal, in Wei.
func (v *VersionedProposal) Value() *big.Int {
	value := big.NewInt(0)
	if v.ConsensusValue != nil {
		value = value.Add(value, v.ConsensusValue)
	}

	if v.ExecutionValue != nil {
		value = value.Add(value, v.ExecutionValue)
	}

	return value
}

// String returns a string version of the structure.
func (v *VersionedProposal) String() string {
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

func (v *VersionedProposal) proposalPresent() bool {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil
		}

		return v.Capella != nil
	}

	return false
}

func (v *VersionedProposal) bodyPresent() bool {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil && v.CapellaBlinded.Body != nil
		}

		return v.Capella != nil && v.Capella.Body != nil
	}

	return false
}

//nolint:gocyclo // ignore
func (v *VersionedProposal) payloadPresent() bool {
	switch v.Version {
	case spec.DataVersionCapella:
		if v.Blinded {
			return v.CapellaBlinded != nil && v.CapellaBlinded.Body != nil && v.CapellaBlinded.Body.ExecutionPayloadHeader != nil
		}

		return v.Capella != nil && v.Capella.Body != nil && v.Capella.Body.ExecutionPayload != nil
	}

	return false
}
