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

	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedSignedAggregateAndProof contains a versioned signed aggregate and proof.
type VersionedSignedAggregateAndProof struct {
	Version DataVersion
	Capella *capella.SignedAggregateAndProof
}

// AggregatorIndex returns the aggregator index of the aggregate.
func (v *VersionedSignedAggregateAndProof) AggregatorIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.AggregatorIndex, nil
	default:
		return 0, errors.New("unknown version for signed aggregate and proof")
	}
}

// IsEmpty returns true if there is no aggregate and proof.
func (v *VersionedSignedAggregateAndProof) IsEmpty() bool {
	return v.Capella == nil
}

// SelectionProof returns the selection proof of the signed aggregate.
func (v *VersionedSignedAggregateAndProof) SelectionProof() (capella.MLDSA87Signature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.MLDSA87Signature{}, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.SelectionProof, nil
	default:
		return capella.MLDSA87Signature{}, errors.New("unknown version")
	}
}

// Signature returns the signature of the signed aggregate and proof.
func (v *VersionedSignedAggregateAndProof) Signature() (capella.MLDSA87Signature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.MLDSA87Signature{}, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Signature, nil
	default:
		return capella.MLDSA87Signature{}, errors.New("unknown version")
	}
}

// Slot returns the slot of the signed aggregate and proof.
func (v *VersionedSignedAggregateAndProof) Slot() (capella.Slot, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella signed aggregate and proof")
		}

		return v.Capella.Message.Aggregate.Data.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedSignedAggregateAndProof) String() string {
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
