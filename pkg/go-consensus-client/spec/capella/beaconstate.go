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
	bitfield "github.com/theQRL/go-bitfield"
)

// BeaconState represents a beacon state.
type BeaconState struct {
	GenesisTime                  uint64
	GenesisValidatorsRoot        Root `ssz-size:"32"`
	Slot                         Slot
	Fork                         *Fork
	LatestBlockHeader            *BeaconBlockHeader
	BlockRoots                   []Root `dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"1024,32"`
	StateRoots                   []Root `dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"1024,32"`
	HistoricalRoots              []Root `dynssz-max:"HISTORICAL_ROOTS_LIMIT"        ssz-max:"16777216" ssz-size:"?,32"`
	ExecutionData                *ExecutionData
	ExecutionDataVotes           []*ExecutionData `dynssz-max:"EPOCHS_PER_EXECUTION_VOTING_PERIOD*SLOTS_PER_EPOCH" ssz-max:"2048"`
	ExecutionDepositIndex        uint64
	Validators                   []*Validator         `dynssz-max:"VALIDATOR_REGISTRY_LIMIT"         ssz-max:"1099511627776"`
	Balances                     []Gwei               `dynssz-max:"VALIDATOR_REGISTRY_LIMIT"         ssz-max:"1099511627776"`
	RANDAOMixes                  []Root               `dynssz-size:"EPOCHS_PER_HISTORICAL_VECTOR,32" ssz-size:"65536,32"`
	Slashings                    []Gwei               `dynssz-size:"EPOCHS_PER_SLASHINGS_VECTOR"     ssz-size:"1024"`
	PreviousEpochParticipation   []ParticipationFlags `dynssz-max:"VALIDATOR_REGISTRY_LIMIT"         ssz-max:"1099511627776"`
	CurrentEpochParticipation    []ParticipationFlags `dynssz-max:"VALIDATOR_REGISTRY_LIMIT"         ssz-max:"1099511627776"`
	JustificationBits            bitfield.Bitvector4  `ssz-size:"1"`
	PreviousJustifiedCheckpoint  *Checkpoint
	CurrentJustifiedCheckpoint   *Checkpoint
	FinalizedCheckpoint          *Checkpoint
	InactivityScores             []uint64 `dynssz-max:"VALIDATOR_REGISTRY_LIMIT" ssz-max:"1099511627776"`
	CurrentSyncCommittee         *SyncCommittee
	NextSyncCommittee            *SyncCommittee
	LatestExecutionPayloadHeader *ExecutionPayloadHeader
	NextWithdrawalIndex          WithdrawalIndex
	NextWithdrawalValidatorIndex ValidatorIndex
	HistoricalSummaries          []*HistoricalSummary `dynssz-max:"HISTORICAL_ROOTS_LIMIT" ssz-max:"16777216"`
}

// String returns a string version of the structure.
func (s *BeaconState) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
