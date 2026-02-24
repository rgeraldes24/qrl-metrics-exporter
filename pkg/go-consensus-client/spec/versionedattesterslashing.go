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

// VersionedAttesterSlashing contains a versioned attestation.
type VersionedAttesterSlashing struct {
	Version DataVersion
	Capella *capella.AttesterSlashing
}

// IsEmpty returns true if there is no block.
func (v *VersionedAttesterSlashing) IsEmpty() bool {
	return v.Capella == nil
}

// Attestation1 returns the first indexed attestation.
func (v *VersionedAttesterSlashing) Attestation1() (*VersionedIndexedAttestation, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionCapella,
			Capella: v.Capella.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Attestation2 returns the second indexed attestation.
func (v *VersionedAttesterSlashing) Attestation2() (*VersionedIndexedAttestation, error) {
	switch v.Version {
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionCapella,
			Capella: v.Capella.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedAttesterSlashing) String() string {
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
