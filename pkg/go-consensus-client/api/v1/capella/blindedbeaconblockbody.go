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

package capella

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// BlindedBeaconBlockBody represents the body of a blinded beacon block.
type BlindedBeaconBlockBody struct {
	RANDAOReveal           capella.BLSSignature `ssz-size:"96"`
	ExecutionData          *capella.ExecutionData
	Graffiti               [32]byte                       `ssz-size:"32"`
	ProposerSlashings      []*capella.ProposerSlashing    `ssz-max:"16"`
	AttesterSlashings      []*capella.AttesterSlashing    `ssz-max:"2"`
	Attestations           []*capella.Attestation         `ssz-max:"128"`
	Deposits               []*capella.Deposit             `ssz-max:"16"`
	VoluntaryExits         []*capella.SignedVoluntaryExit `ssz-max:"16"`
	SyncAggregate          *capella.SyncAggregate
	ExecutionPayloadHeader *capella.ExecutionPayloadHeader
}

// String returns a string version of the structure.
func (b *BlindedBeaconBlockBody) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
