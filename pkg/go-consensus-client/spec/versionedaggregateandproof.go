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

// VersionedAggregateAndProof contains a versioned aggregate and proof.
type VersionedAggregateAndProof struct {
	Version DataVersion
	Capella *capella.AggregateAndProof
}

// AggregatorIndex returns the aggregator index of the aggregate.
func (v *VersionedAggregateAndProof) AggregatorIndex() (capella.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella aggregate and proof")
		}

		return v.Capella.AggregatorIndex, nil
	default:
		return 0, errors.New("unknown version for aggregate and proof")
	}
}

// HashTreeRoot returns the hash tree root of the aggregate and proof.
func (v *VersionedAggregateAndProof) HashTreeRoot() ([32]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella aggregate and proof")
		}

		return v.Capella.HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// IsEmpty returns true if there is no aggregate and proof.
func (v *VersionedAggregateAndProof) IsEmpty() bool {
	return v.Capella == nil
}

// String returns a string version of the structure.
func (v *VersionedAggregateAndProof) String() string {
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

// SelectionProof returns the selection proof of the aggregate.
func (v *VersionedAggregateAndProof) SelectionProof() (capella.MLDSA87Signature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.MLDSA87Signature{}, errors.New("no capella aggregate and proof")
		}

		return v.Capella.SelectionProof, nil
	default:
		return capella.MLDSA87Signature{}, errors.New("unknown version")
	}
}
