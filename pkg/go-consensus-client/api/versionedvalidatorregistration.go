// Copyright © 2022 Attestant Limited.
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
	"time"

	apiv1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// VersionedValidatorRegistration contains a versioned ValidatorRegistrationV1.
type VersionedValidatorRegistration struct {
	Version spec.BuilderVersion
	V1      *apiv1.ValidatorRegistration
}

// IsEmpty returns true if there is no block.
func (v *VersionedValidatorRegistration) IsEmpty() bool {
	return v.V1 == nil
}

// FeeRecipient returns the fee recipient of the validator registration.
func (v *VersionedValidatorRegistration) FeeRecipient() (capella.ExecutionAddress, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return capella.ExecutionAddress{}, ErrDataMissing
		}

		return v.V1.FeeRecipient, nil
	default:
		return capella.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// GasLimit returns the gas limit of the validator registration.
func (v *VersionedValidatorRegistration) GasLimit() (uint64, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return 0, ErrDataMissing
		}

		return v.V1.GasLimit, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the validator registration.
func (v *VersionedValidatorRegistration) Timestamp() (time.Time, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return time.Time{}, ErrDataMissing
		}

		return v.V1.Timestamp, nil
	default:
		return time.Time{}, ErrUnsupportedVersion
	}
}

// PubKey returns the public key of the validator registration.
func (v *VersionedValidatorRegistration) PubKey() (capella.BLSPubKey, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return capella.BLSPubKey{}, ErrDataMissing
		}

		return v.V1.Pubkey, nil
	default:
		return capella.BLSPubKey{}, ErrUnsupportedVersion
	}
}

// Root returns the root of the validator registration.
func (v *VersionedValidatorRegistration) Root() (capella.Root, error) {
	switch v.Version {
	case spec.BuilderVersionV1:
		if v.V1 == nil {
			return capella.Root{}, ErrDataMissing
		}

		return v.V1.HashTreeRoot()
	default:
		return capella.Root{}, ErrUnsupportedVersion
	}
}
