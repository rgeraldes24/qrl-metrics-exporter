// Copyright © 2024 Attestant Limited.
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
	"fmt"

	"github.com/theQRL/go-bitfield"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedAttestation contains a versioned attestation.
type VersionedAttestation struct {
	Version        DataVersion
	ValidatorIndex *capella.ValidatorIndex
	Capella        *capella.Attestation
}

// IsEmpty returns true if there is no block.
func (v *VersionedAttestation) IsEmpty() bool {
	return v.Capella == nil
}

// AggregationBits returns the aggregation bits of the attestation.
func (v *VersionedAttestation) AggregationBits() (bitfield.Bitlist, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella attestation")
		}

		return v.Capella.AggregationBits, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Data returns the data of the attestation.
func (v *VersionedAttestation) Data() (*capella.AttestationData, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella attestation")
		}

		return v.Capella.Data, nil
	default:
		return nil, fmt.Errorf("unknown version: %d", v.Version)
	}
}

// CommitteeBits returns the committee bits of the attestation.
func (v *VersionedAttestation) CommitteeBits() (bitfield.Bitvector64, error) {
	switch v.Version {
	case DataVersionCapella:
		return nil, errors.New("attestation does not provide committee bits")
	default:
		return nil, errors.New("unknown version")
	}
}

// CommitteeIndex returns the index if only one bit is set, otherwise error.
func (v *VersionedAttestation) CommitteeIndex() (capella.CommitteeIndex, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no Capella attestation")
		}

		return v.Capella.Data.Index, nil
	default:
		return 0, errors.New("unknown version")
	}
}

func (v *VersionedAttestation) HashTreeRoot() ([32]byte, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no Capella attestation")
		}

		return v.Capella.HashTreeRoot()
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// Signature returns the signature of the attestation.
func (v *VersionedAttestation) Signatures() ([]capella.MLDSA87Signature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return []capella.MLDSA87Signature{}, errors.New("no Capella attestation")
		}

		return v.Capella.Signatures, nil
	default:
		return []capella.MLDSA87Signature{}, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedAttestation) String() string {
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
