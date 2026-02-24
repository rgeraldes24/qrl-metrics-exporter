// Copyright © 2025 Attestant Limited.
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

	"github.com/holiman/uint256"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedExecutionPayload contains a versioned execution payload.
type VersionedExecutionPayload struct {
	Version DataVersion
	Capella *capella.ExecutionPayload
}

// IsEmpty returns true if there is no block.
func (v *VersionedExecutionPayload) IsEmpty() bool {
	return v.Capella == nil
}

// ParentHash returns the parent hash of the execution payload.
func (v *VersionedExecutionPayload) ParentHash() (capella.Hash32, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Hash32{}, errors.New("no capella execution payload")
		}

		return v.Capella.ParentHash, nil
	default:
		return capella.Hash32{}, errors.New("unknown version")
	}
}

// FeeRecipient returns the fee recipient of the execution payload.
func (v *VersionedExecutionPayload) FeeRecipient() (capella.ExecutionAddress, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.ExecutionAddress{}, errors.New("no capella execution payload")
		}

		return v.Capella.FeeRecipient, nil
	default:
		return capella.ExecutionAddress{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the execution payload.
func (v *VersionedExecutionPayload) StateRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella execution payload")
		}

		return v.Capella.StateRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// ReceiptsRoot returns the receipts root of the execution payload.
func (v *VersionedExecutionPayload) ReceiptsRoot() (capella.Root, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Root{}, errors.New("no capella execution payload")
		}

		return v.Capella.ReceiptsRoot, nil
	default:
		return capella.Root{}, errors.New("unknown version")
	}
}

// LogsBloom returns the logs bloom of the execution payload.
func (v *VersionedExecutionPayload) LogsBloom() ([256]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return [256]byte{}, errors.New("no capella execution payload")
		}

		return v.Capella.LogsBloom, nil
	default:
		return [256]byte{}, errors.New("unknown version")
	}
}

// PrevRandao returns the prev randao of the execution payload.
func (v *VersionedExecutionPayload) PrevRandao() ([32]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella execution payload")
		}

		return v.Capella.PrevRandao, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// BlockNumber returns the block number of the execution payload.
func (v *VersionedExecutionPayload) BlockNumber() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// GasLimit returns the gas limit of the execution payload.
func (v *VersionedExecutionPayload) GasLimit() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.GasLimit, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// GasUsed returns the gas used of the execution payload.
func (v *VersionedExecutionPayload) GasUsed() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.GasUsed, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Timestamp returns the timestamp of the execution payload.
func (v *VersionedExecutionPayload) Timestamp() (uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella execution payload")
		}

		return v.Capella.Timestamp, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExtraData returns the extra data of the execution payload.
func (v *VersionedExecutionPayload) ExtraData() ([]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.ExtraData, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BaseFeePerGas returns the base fee per gas of the execution payload.
func (v *VersionedExecutionPayload) BaseFeePerGas() (*uint256.Int, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return uint256.NewInt(0).SetBytes(v.Capella.BaseFeePerGas[:]), nil
	default:
		return nil, errors.New("unknown version")
	}
}

// BlockHash returns the block hash of the execution payload.
func (v *VersionedExecutionPayload) BlockHash() (capella.Hash32, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.Hash32{}, errors.New("no capella execution payload")
		}

		return v.Capella.BlockHash, nil
	default:
		return capella.Hash32{}, errors.New("unknown version")
	}
}

// Transactions returns the transactions of the execution payload.
func (v *VersionedExecutionPayload) Transactions() ([]capella.Transaction, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.Transactions, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Withdrawals returns the withdrawals of the execution payload.
func (v *VersionedExecutionPayload) Withdrawals() ([]*capella.Withdrawal, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella execution payload")
		}

		return v.Capella.Withdrawals, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedExecutionPayload) String() string {
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
