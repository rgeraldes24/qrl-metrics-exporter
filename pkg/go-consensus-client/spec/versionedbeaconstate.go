// Copyright © 2021 - 2025 Attestant Limited.
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

	ssz "github.com/ferranbt/fastssz"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
	proofutil "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/util/proof"
)

// VersionedBeaconState contains a versioned beacon state.
type VersionedBeaconState struct {
	Version DataVersion
	Capella *capella.BeaconState
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconState) IsEmpty() bool {
	return v.Capella == nil
}

// Slot returns the slot of the state.
func (v *VersionedBeaconState) Slot() (capella.Slot, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return v.Capella.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// NextWithdrawalValidatorIndex returns the next withdrawal validator index of the state.
func (v *VersionedBeaconState) NextWithdrawalValidatorIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return v.Capella.NextWithdrawalValidatorIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Validators returns the validators of the state.
func (v *VersionedBeaconState) Validators() ([]*capella.Validator, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.Validators, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ValidatorBalances returns the validator balances of the state.
func (v *VersionedBeaconState) ValidatorBalances() ([]capella.Shor, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.Balances, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// DepositRequestsStartIndex returns the deposit requests start index of the state.
func (v *VersionedBeaconState) DepositRequestsStartIndex() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide deposit requests start index")
	default:
		return 0, errors.New("unknown version")
	}
}

// DepositBalanceToConsume returns the deposit balance to consume of the state.
func (v *VersionedBeaconState) DepositBalanceToConsume() (capella.Shor, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide deposit balance to consume")
	default:
		return 0, errors.New("unknown version")
	}
}

// ExitBalanceToConsume returns the deposit balance to consume of the state.
func (v *VersionedBeaconState) ExitBalanceToConsume() (capella.Shor, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide exit balance to consume")
	default:
		return 0, errors.New("unknown version")
	}
}

// EarliestExitEpoch returns the earliest exit epoch of the state.
func (v *VersionedBeaconState) EarliestExitEpoch() (capella.Epoch, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide earliest exit epoch")
	default:
		return 0, errors.New("unknown version")
	}
}

// ConsolidationBalanceToConsume returns the consolidation balance to consume of the state.
func (v *VersionedBeaconState) ConsolidationBalanceToConsume() (capella.Shor, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide consolidation balance to consume")
	default:
		return 0, errors.New("unknown version")
	}
}

// EarliestConsolidationEpoch returns the earliest consolidation epoch of the state.
func (v *VersionedBeaconState) EarliestConsolidationEpoch() (capella.Epoch, error) {
	switch v.Version {
	case DataVersionCapella:
		return 0, errors.New("state does not provide earliest consolidation epoch")
	default:
		return 0, errors.New("unknown version")
	}
}

// ValidatorAtIndex returns the validator at the given index.
// This is a convenience method that handles accessing the validators array.
// Parameters:
//   - index: The index of the validator to retrieve
//
// Returns:
//   - *capella.Validator: The validator at the given index
//   - error: If the index is invalid or there's an error accessing the validators
func (v *VersionedBeaconState) ValidatorAtIndex(index capella.ValidatorIndex) (*capella.Validator, error) {
	validators, err := v.Validators()
	if err != nil {
		return nil, err
	}

	if index >= capella.ValidatorIndex(len(validators)) {
		return nil, errors.New("validator index out of bounds")
	}

	return validators[index], nil
}

// ValidatorBalance returns the balance of the validator at the given index.
// This is a convenience method that handles accessing the balances array.
// Parameters:
//   - index: The index of the validator whose balance to retrieve
//
// Returns:
//   - capella.Shor: The balance in Shor
//   - error: If the index is invalid or there's an error accessing the balances
func (v *VersionedBeaconState) ValidatorBalance(index capella.ValidatorIndex) (capella.Shor, error) {
	balances, err := v.ValidatorBalances()
	if err != nil {
		return 0, err
	}

	if index >= capella.ValidatorIndex(len(balances)) {
		return 0, errors.New("validator index out of bounds")
	}

	return balances[index], nil
}

// GetTree returns the GetTree of the specific beacon state version.
func (v *VersionedBeaconState) GetTree() (*ssz.Node, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella state")
		}

		return v.Capella.GetTree()
	default:
		return nil, errors.New("unknown version")
	}
}

// HashTreeRoot returns the HashTreeRoot of the specific beacon state version.
func (v *VersionedBeaconState) HashTreeRoot() (capella.Hash32, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Hash32{}, errors.New("no Capella state")
		}

		return v.Capella.HashTreeRoot()
	default:
		return capella.Hash32{}, errors.New("unknown version")
	}
}

// FieldIndex returns the struct field index for a given field name.
// The index represents the field's position in the struct's memory layout.
// Returns an error if the field doesn't exist or the state is empty.
func (v *VersionedBeaconState) FieldIndex(name string) (int, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return proofutil.FieldIndex(v.Capella, name)
	default:
		return 0, errors.New("unknown version")
	}
}

// FieldGeneralizedIndex returns the generalized index for a given field name.
// The generalized index represents the field's absolute position in the Merkle tree.
// This is used for generating and verifying Merkle proofs.
// Returns an error if the field doesn't exist or the state is empty.
func (v *VersionedBeaconState) FieldGeneralizedIndex(name string) (int, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella state")
		}

		return proofutil.FieldGeneralizedIndex(v.Capella, name)
	default:
		return 0, errors.New("unknown version")
	}
}

// FieldRoot returns the SSZ hash root of a specific field in the beacon state.
// Parameters:
//   - name: The name of the field to get the root for
//
// Returns:
//   - capella.Hash32: The SSZ hash root of the field
//   - error: If the field doesn't exist, the state is empty, or the field is not hash tree rootable
func (v *VersionedBeaconState) FieldRoot(name string) (capella.Hash32, error) {
	fieldTree, err := v.FieldTree(name)
	if err != nil {
		return capella.Hash32{}, err
	}

	var root capella.Hash32
	copy(root[:], fieldTree.Hash())

	return root, nil
}

// FieldTree returns the Merkle subtree for a specific field in the beacon state.
func (v *VersionedBeaconState) FieldTree(name string) (*ssz.Node, error) {
	stateTree, err := v.GetTree()
	if err != nil {
		return nil, err
	}

	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return nil, err
	}

	return stateTree.Get(fieldGeneralizedIndex)
}

// ProveField generates a Merkle proof for a specific field against the beacon state root.
// Parameters:
//   - name: The name of the field to generate a proof for
//
// Returns:
//   - []capella.Hash32: The Merkle proof as a sequence of 32-byte hashes
//   - error: If the field doesn't exist or there's an error generating the proof
func (v *VersionedBeaconState) ProveField(name string) ([]capella.Hash32, error) {
	stateTree, err := v.GetTree()
	if err != nil {
		return nil, err
	}

	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return nil, err
	}

	proof, err := stateTree.Prove(fieldGeneralizedIndex)
	if err != nil {
		return nil, err
	}

	proofBytes := make([]capella.Hash32, len(proof.Hashes))
	for i, hash := range proof.Hashes {
		copy(proofBytes[i][:], hash)
	}

	return proofBytes, nil
}

// VerifyFieldProof verifies a Merkle proof for a field against the beacon state root.
// Parameters:
//   - proof: The Merkle proof as a sequence of 32-byte hashes
//   - name: The name of the field the proof is for
//
// Returns:
//   - bool: True if the proof is valid, false otherwise
//   - error: If there's an error during verification
func (v *VersionedBeaconState) VerifyFieldProof(proof []capella.Hash32, name string) (bool, error) {
	// Get the state root
	stateRoot, err := v.HashTreeRoot()
	if err != nil {
		return false, err
	}

	// Get the field's generalized index
	fieldGeneralizedIndex, err := v.FieldGeneralizedIndex(name)
	if err != nil {
		return false, err
	}

	// Get the field's root
	fieldRoot, err := v.FieldRoot(name)
	if err != nil {
		return false, err
	}

	// Convert proof to ssz.Proof format
	proofHashes := make([][]byte, len(proof))
	for i, hash := range proof {
		proofHashes[i] = make([]byte, 32)
		copy(proofHashes[i], hash[:])
	}

	// Create and verify the proof
	sszProof := &ssz.Proof{
		Index:  fieldGeneralizedIndex,
		Leaf:   fieldRoot[:],
		Hashes: proofHashes,
	}

	return ssz.VerifyProof(stateRoot[:], sszProof)
}

// String returns a string version of the structure.
func (v *VersionedBeaconState) String() string {
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
