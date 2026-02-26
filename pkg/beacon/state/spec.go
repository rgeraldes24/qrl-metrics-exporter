package state

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	sp "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// Spec represents the state of the spec.
type Spec struct {
	PresetBase string `json:"PRESET_BASE"`
	ConfigName string `json:"CONFIG_NAME"`

	DepositChainID         uint64 `json:"DEPOSIT_CHAIN_ID,string"`
	DepositContractAddress string `json:"DEPOSIT_CONTRACT_ADDRESS"`

	SafeSlotsToUpdateJustified capella.Slot `json:"SAFE_SLOTS_TO_UPDATE_JUSTIFIED,string"`
	SlotsPerEpoch              capella.Slot `json:"SLOTS_PER_EPOCH,string"`

	EpochsPerSyncCommitteePeriod capella.Epoch `json:"EPOCHS_PER_SYNC_COMMITTEE_PERIOD,string"`
	MinSyncCommitteeParticipants uint64        `json:"MIN_SYNC_COMMITTEE_PARTICIPANTS,string"`
	TargetCommitteeSize          uint64        `json:"TARGET_COMMITTEE_SIZE,string"`
	SyncCommitteeSize            uint64        `json:"SYNC_COMMITTEE_SIZE,string"`

	MaxValidatorsPerCommittee uint64       `json:"MAX_VALIDATORS_PER_COMMITTEE,string"`
	BaseRewardFactor          uint64       `json:"BASE_REWARD_FACTOR,string"`
	EffectiveBalanceIncrement capella.Shor `json:"EFFECTIVE_BALANCE_INCREMENT,string"`
	MaxEffectiveBalance       capella.Shor `json:"MAX_EFFECTIVE_BALANCE,string"`
	MinDepositAmount          capella.Shor `json:"MIN_DEPOSIT_AMOUNT,string"`
	MaxAttestations           uint64       `json:"MAX_ATTESTATIONS,string"`

	SecondsPerExecutionBlock       StringerDuration `json:"SECONDS_PER_EXECUTION_BLOCK,string"`
	GenesisDelay                   StringerDuration `json:"GENESIS_DELAY,string"`
	SecondsPerSlot                 StringerDuration `json:"SECONDS_PER_SLOT,string"`
	MaxDeposits                    uint64           `json:"MAX_DEPOSITS,string"`
	MinGenesisActiveValidatorCount uint64           `json:"MIN_GENESIS_ACTIVE_VALIDATOR_COUNT,string"`
	ExecutionFollowDistance        uint64           `json:"EXECUTION_FOLLOW_DISTANCE,string"`

	ForkEpochs ForkEpochs     `json:"-"`
	FullSpec   map[string]any `json:"-"`
}

// NewSpec creates a new spec instance.
//
//nolint:gocyclo // existing.
func NewSpec(data map[string]interface{}) Spec {
	spec := Spec{
		ForkEpochs: ForkEpochs{},
		FullSpec:   data,
	}

	if safeSlotsToUpdateJustified, exists := data["SAFE_SLOTS_TO_UPDATE_JUSTIFIED"]; exists {
		spec.SafeSlotsToUpdateJustified = capella.Slot(cast.ToUint64(safeSlotsToUpdateJustified))
	}

	if depositChainID, exists := data["DEPOSIT_CHAIN_ID"]; exists {
		spec.DepositChainID = cast.ToUint64(depositChainID)
	}

	if depositContractAddress, exists := data["DEPOSIT_CONTRACT_ADDRESS"]; exists {
		spec.DepositContractAddress = fmt.Sprintf("%#x", cast.ToString(depositContractAddress))
	}

	if configName, exists := data["CONFIG_NAME"]; exists {
		spec.ConfigName = cast.ToString(configName)
	}

	if maxValidatorsPerCommittee, exists := data["MAX_VALIDATORS_PER_COMMITTEE"]; exists {
		spec.MaxValidatorsPerCommittee = cast.ToUint64(maxValidatorsPerCommittee)
	}

	if secondsPerExecutionBlock, exists := data["SECONDS_PER_EXECUTION_BLOCK"]; exists {
		spec.SecondsPerExecutionBlock = StringerDuration(cast.ToDuration(secondsPerExecutionBlock))
	}

	if baseRewardFactor, exists := data["BASE_REWARD_FACTOR"]; exists {
		spec.BaseRewardFactor = cast.ToUint64(baseRewardFactor)
	}

	if epochsPerSyncComitteePeriod, exists := data["EPOCHS_PER_SYNC_COMMITTEE_PERIOD"]; exists {
		spec.EpochsPerSyncCommitteePeriod = capella.Epoch(cast.ToUint64(epochsPerSyncComitteePeriod))
	}

	if effectiveBalanceIncrement, exists := data["EFFECTIVE_BALANCE_INCREMENT"]; exists {
		spec.EffectiveBalanceIncrement = capella.Shor(cast.ToUint64(effectiveBalanceIncrement))
	}

	if maxAttestations, exists := data["MAX_ATTESTATIONS"]; exists {
		spec.MaxAttestations = cast.ToUint64(maxAttestations)
	}

	if minSyncCommitteeParticipants, exists := data["MIN_SYNC_COMMITTEE_PARTICIPANTS"]; exists {
		spec.MinSyncCommitteeParticipants = cast.ToUint64(minSyncCommitteeParticipants)
	}

	if genesisDelay, exists := data["GENESIS_DELAY"]; exists {
		spec.GenesisDelay = StringerDuration(cast.ToDuration(genesisDelay))
	}

	if secondsPerSlot, exists := data["SECONDS_PER_SLOT"]; exists {
		spec.SecondsPerSlot = StringerDuration(cast.ToDuration(secondsPerSlot))
	}

	if maxEffectiveBalance, exists := data["MAX_EFFECTIVE_BALANCE"]; exists {
		spec.MaxEffectiveBalance = capella.Shor(cast.ToUint64(maxEffectiveBalance))
	}

	if maxDeposits, exists := data["MAX_DEPOSITS"]; exists {
		spec.MaxDeposits = cast.ToUint64(maxDeposits)
	}

	if minGenesisActiveValidatorCount, exists := data["MIN_GENESIS_ACTIVE_VALIDATOR_COUNT"]; exists {
		spec.MinGenesisActiveValidatorCount = cast.ToUint64(minGenesisActiveValidatorCount)
	}

	if targetCommitteeSize, exists := data["TARGET_COMMITTEE_SIZE"]; exists {
		spec.TargetCommitteeSize = cast.ToUint64(targetCommitteeSize)
	}

	if syncCommitteeSize, exists := data["SYNC_COMMITTEE_SIZE"]; exists {
		spec.SyncCommitteeSize = cast.ToUint64(syncCommitteeSize)
	}

	if executionFollowDistance, exists := data["EXECUTION_FOLLOW_DISTANCE"]; exists {
		spec.ExecutionFollowDistance = cast.ToUint64(executionFollowDistance)
	}

	if minDepositAmount, exists := data["MIN_DEPOSIT_AMOUNT"]; exists {
		spec.MinDepositAmount = capella.Shor(cast.ToUint64(minDepositAmount))
	}

	if slotsPerEpoch, exists := data["SLOTS_PER_EPOCH"]; exists {
		spec.SlotsPerEpoch = capella.Slot(cast.ToUint64(slotsPerEpoch))
	}

	if presetBase, exists := data["PRESET_BASE"]; exists {
		spec.PresetBase = cast.ToString(presetBase)
	}

	forkEpochs := make(map[string]capella.Epoch)
	forkVersions := make(map[string]string)

	forkEpochs["GENESIS"] = 0

	for k, v := range data {
		if strings.Contains(k, "_FORK_EPOCH") {
			forkName := strings.ReplaceAll(k, "_FORK_EPOCH", "")

			forkEpochs[forkName] = capella.Epoch(cast.ToUint64(v))
		}

		if strings.Contains(k, "_FORK_VERSION") {
			forkName := strings.ReplaceAll(k, "_FORK_VERSION", "")

			forkVersions[forkName] = fmt.Sprintf("%#x", v)
		}
	}

	for k, v := range forkEpochs {
		version := ""
		if v, exists := forkVersions[k]; exists {
			version = v
		}

		// Convert the name to a DataVersion.
		dataVersion, err := dataVersionFromString(k)
		if err != nil {
			continue
		}

		spec.ForkEpochs = append(spec.ForkEpochs, &ForkEpoch{
			Epoch:   v,
			Name:    dataVersion,
			Version: version,
		})
	}

	return spec
}

// Validate performs basic validation of the spec.
func (s *Spec) Validate() error {
	return nil
}

func dataVersionFromString(name string) (sp.DataVersion, error) {
	var v sp.DataVersion
	if err := json.Unmarshal([]byte(fmt.Sprintf("\"%s\"", name)), &v); err != nil {
		return sp.DataVersionUnknown, err
	}

	return v, nil
}
