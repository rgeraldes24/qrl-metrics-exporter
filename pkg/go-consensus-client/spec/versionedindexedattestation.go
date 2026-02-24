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

	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedIndexedAttestation contains a versioned indexed attestation.
type VersionedIndexedAttestation struct {
	Version DataVersion
	Capella *capella.IndexedAttestation
}

// IsEmpty returns true if there is no block.
func (v *VersionedIndexedAttestation) IsEmpty() bool {
	return v.Capella == nil
}

// AttestingIndices returns the attesting indices of the indexed attestation.
func (v *VersionedIndexedAttestation) AttestingIndices() ([]uint64, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		return v.Capella.AttestingIndices, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Data returns the data of the indexed attestation.
func (v *VersionedIndexedAttestation) Data() (*capella.AttestationData, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		return v.Capella.Data, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Signature returns the signature of the indexed attestation.
func (v *VersionedIndexedAttestation) Signature() (capella.BLSSignature, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return capella.BLSSignature{}, errors.New("no Capella indexed attestation")
		}

		return v.Capella.Signature, nil
	default:
		return capella.BLSSignature{}, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedIndexedAttestation) String() string {
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
