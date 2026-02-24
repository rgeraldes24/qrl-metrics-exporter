// Copyright © 2020 - 2023 Attestant Limited.
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

package client

import (
	"context"
	"time"

	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	apiv1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// Service is the service providing a connection to a QRL client.
type Service interface {
	// Name returns the name of the client implementation.
	Name() string

	// Address returns the address of the client.
	Address() string

	// IsActive returns true if the client is active.
	IsActive() bool

	// IsSynced returns true if the client is synced.
	IsSynced() bool
}

// EpochFromStateIDProvider is the interface for providing epochs from state IDs.
type EpochFromStateIDProvider interface {
	// EpochFromStateID converts a state ID to its epoch.
	//
	// Deprecated: will be removed in a future release.
	EpochFromStateID(ctx context.Context, stateID string) (capella.Epoch, error)
}

// SlotFromStateIDProvider is the interface for providing slots from state IDs.
type SlotFromStateIDProvider interface {
	// SlotFromStateID converts a state ID to its slot.
	//
	// Deprecated: will be removed in a future release.
	SlotFromStateID(ctx context.Context, stateID string) (capella.Slot, error)
}

// SlotDurationProvider is the interface for providing the duration of each slot of a chain.
type SlotDurationProvider interface {
	// SlotDuration provides the duration of a slot of the chain.
	//
	// Deprecated: use Spec()
	SlotDuration(ctx context.Context) (time.Duration, error)
}

// SlotsPerEpochProvider is the interface for providing the number of slots in each epoch of a chain.
type SlotsPerEpochProvider interface {
	// SlotsPerEpoch provides the slots per epoch of the chain.
	//
	// Deprecated: use Spec()
	SlotsPerEpoch(ctx context.Context) (uint64, error)
}

// FarFutureEpochProvider is the interface for providing the far future epoch of a chain.
type FarFutureEpochProvider interface {
	// FarFutureEpoch provides the far future epoch of the chain.
	FarFutureEpoch(ctx context.Context) (capella.Epoch, error)
}

// TargetAggregatorsPerCommitteeProvider is the interface for providing the target number of
// aggregators in each attestation committee.
type TargetAggregatorsPerCommitteeProvider interface {
	// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
	//
	// Deprecated: use Spec()
	TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error)
}

// ValidatorIndexProvider is the interface for entities that can provide the index of a validator.
type ValidatorIndexProvider interface {
	// Index provides the index of the validator.
	Index(ctx context.Context) (capella.ValidatorIndex, error)
}

// ValidatorPubKeyProvider is the interface for entities that can provide the public key of a validator.
type ValidatorPubKeyProvider interface {
	// PubKey provides the public key of the validator.
	PubKey(ctx context.Context) (capella.BLSPubKey, error)
}

// ValidatorIDProvider is the interface that provides the identifiers (pubkey, index) of a validator.
type ValidatorIDProvider interface {
	ValidatorIndexProvider
	ValidatorPubKeyProvider
}

// SignedBeaconBlockProvider is the interface for providing beacon blocks.
type SignedBeaconBlockProvider interface {
	// SignedBeaconBlock fetches a signed beacon block given a block ID.
	SignedBeaconBlock(ctx context.Context,
		opts *api.SignedBeaconBlockOpts,
	) (
		*api.Response[*spec.VersionedSignedBeaconBlock],
		error,
	)
}

// BeaconCommitteesProvider is the interface for providing beacon committees.
type BeaconCommitteesProvider interface {
	// BeaconCommittees fetches all beacon committees for the given options.
	BeaconCommittees(ctx context.Context,
		opts *api.BeaconCommitteesOpts,
	) (*api.Response[[]*apiv1.BeaconCommittee],
		error,
	)
}

// SyncCommitteesProvider is the interface for providing sync committees.
type SyncCommitteesProvider interface {
	// SyncCommittee fetches the sync committee for the given state.
	SyncCommittee(ctx context.Context,
		opts *api.SyncCommitteeOpts,
	) (
		*api.Response[*apiv1.SyncCommittee],
		error,
	)
}

//
// Standard API
//

// AggregateAttestationProvider is the interface for providing aggregate attestations.
type AggregateAttestationProvider interface {
	// AggregateAttestation fetches the aggregate attestation for the given options.
	AggregateAttestation(ctx context.Context,
		opts *api.AggregateAttestationOpts,
	) (
		*api.Response[*spec.VersionedAttestation],
		error,
	)
}

// AttestationDataProvider is the interface for providing attestation data.
type AttestationDataProvider interface {
	// AttestationData fetches the attestation data for the given options.
	AttestationData(ctx context.Context,
		opts *api.AttestationDataOpts,
	) (
		*api.Response[*capella.AttestationData],
		error,
	)
}

// AttestationPoolProvider is the interface for providing attestation pools.
type AttestationPoolProvider interface {
	// AttestationPool fetches the attestation pool for the given options.
	AttestationPool(ctx context.Context,
		opts *api.AttestationPoolOpts,
	) (
		*api.Response[[]*spec.VersionedAttestation],
		error,
	)
}

// AttestationRewardsProvider is the interface for providing attestation rewards.
type AttestationRewardsProvider interface {
	// AttestationRewards provides rewards to the given validators for attesting.
	AttestationRewards(ctx context.Context,
		opts *api.AttestationRewardsOpts,
	) (
		*api.Response[*apiv1.AttestationRewards],
		error,
	)
}

// AttesterDutiesProvider is the interface for providing attester duties.
type AttesterDutiesProvider interface {
	// AttesterDuties obtains attester duties.
	AttesterDuties(ctx context.Context,
		opts *api.AttesterDutiesOpts,
	) (
		*api.Response[[]*apiv1.AttesterDuty],
		error,
	)
}

// BlockRewardsProvider is the interface for providing block rewards.
type BlockRewardsProvider interface {
	// BlockRewards provides rewards for proposing a block.
	BlockRewards(ctx context.Context,
		opts *api.BlockRewardsOpts,
	) (
		*api.Response[*apiv1.BlockRewards],
		error,
	)
}

// DepositContractProvider is the interface for providing details about the deposit contract.
type DepositContractProvider interface {
	// DepositContract provides details of the execution deposit contract for the chain.
	DepositContract(ctx context.Context,
		opts *api.DepositContractOpts,
	) (
		*api.Response[*apiv1.DepositContract],
		error,
	)
}

// SyncCommitteeDutiesProvider is the interface for providing sync committee duties.
type SyncCommitteeDutiesProvider interface {
	// SyncCommitteeDuties obtains sync committee duties.
	// If validatorIndices is nil it will return all duties for the given epoch.
	SyncCommitteeDuties(ctx context.Context,
		opts *api.SyncCommitteeDutiesOpts,
	) (
		*api.Response[[]*apiv1.SyncCommitteeDuty],
		error,
	)
}

// SyncCommitteeContributionProvider is the interface for providing sync committee contributions.
type SyncCommitteeContributionProvider interface {
	// SyncCommitteeContribution provides a sync committee contribution.
	SyncCommitteeContribution(ctx context.Context,
		opts *api.SyncCommitteeContributionOpts,
	) (
		*api.Response[*capella.SyncCommitteeContribution],
		error,
	)
}

// SyncCommitteeRewardsProvider is the interface for providing sync committee rewards.
type SyncCommitteeRewardsProvider interface {
	// SyncCommitteeRewards provides rewards to the given validators for being members of a sync committee.
	SyncCommitteeRewards(ctx context.Context,
		opts *api.SyncCommitteeRewardsOpts,
	) (
		*api.Response[[]*apiv1.SyncCommitteeReward],
		error,
	)
}

// BeaconBlockHeadersProvider is the interface for providing beacon block headers.
type BeaconBlockHeadersProvider interface {
	// BeaconBlockHeader provides the block header of a given block ID.
	BeaconBlockHeader(ctx context.Context,
		opts *api.BeaconBlockHeaderOpts,
	) (
		*api.Response[*apiv1.BeaconBlockHeader],
		error,
	)
}

// ProposalProvider is the interface for providing proposals.
type ProposalProvider interface {
	// Proposal fetches a proposal for signing.
	Proposal(ctx context.Context,
		opts *api.ProposalOpts,
	) (
		*api.Response[*api.VersionedProposal],
		error,
	)
}

// BeaconBlockRootProvider is the interface for providing beacon block roots.
type BeaconBlockRootProvider interface {
	// BeaconBlockRoot fetches a block's root given a set of options.
	BeaconBlockRoot(ctx context.Context,
		opts *api.BeaconBlockRootOpts,
	) (
		*api.Response[*capella.Root],
		error,
	)
}

// BeaconStateProvider is the interface for providing beacon state.
type BeaconStateProvider interface {
	// BeaconState fetches a beacon state given a state ID.
	BeaconState(ctx context.Context,
		opts *api.BeaconStateOpts,
	) (*api.Response[*spec.VersionedBeaconState],
		error,
	)
}

// BeaconStateRandaoProvider is the interface for providing beacon state RANDAOs.
type BeaconStateRandaoProvider interface {
	// BeaconStateRandao fetches a beacon state RANDAO given a state ID.
	BeaconStateRandao(ctx context.Context,
		opts *api.BeaconStateRandaoOpts,
	) (
		*api.Response[*capella.Root],
		error,
	)
}

// BeaconStateRootProvider is the interface for providing beacon state roots.
type BeaconStateRootProvider interface {
	// BeaconStateRoot fetches a beacon state root given a state ID.
	BeaconStateRoot(ctx context.Context,
		opts *api.BeaconStateRootOpts,
	) (
		*api.Response[*capella.Root],
		error,
	)
}

// EventsProvider is the interface for providing events.
type EventsProvider interface {
	// Events feeds requested events with the given topics to the supplied handler.
	Events(ctx context.Context, opts *api.EventsOpts) error
}

// FinalityProvider is the interface for providing finality information.
type FinalityProvider interface {
	// Finality provides the finality given a state ID.
	Finality(ctx context.Context,
		opts *api.FinalityOpts,
	) (
		*api.Response[*apiv1.Finality],
		error,
	)
}

// ForkChoiceProvider is the interface for providing fork choice information.
type ForkChoiceProvider interface {
	// Fork fetches all current fork choice context.
	ForkChoice(ctx context.Context,
		opts *api.ForkChoiceOpts,
	) (
		*api.Response[*apiv1.ForkChoice],
		error,
	)
}

// ForkProvider is the interface for providing fork information.
type ForkProvider interface {
	// Fork fetches fork information for the given state.
	Fork(ctx context.Context,
		opts *api.ForkOpts,
	) (
		*api.Response[*capella.Fork],
		error,
	)
}

// ForkScheduleProvider is the interface for providing fork schedule data.
type ForkScheduleProvider interface {
	// ForkSchedule provides details of past and future changes in the chain's fork version.
	ForkSchedule(ctx context.Context,
		opts *api.ForkScheduleOpts,
	) (
		*api.Response[[]*capella.Fork],
		error,
	)
}

// GenesisProvider is the interface for providing genesis information.
type GenesisProvider interface {
	// Genesis fetches genesis information for the chain.
	Genesis(ctx context.Context,
		opts *api.GenesisOpts,
	) (
		*api.Response[*apiv1.Genesis],
		error,
	)
}

// NodePeersProvider is the interface for providing peer information.
type NodePeersProvider interface {
	// NodePeers provides the peers of the node.
	NodePeers(ctx context.Context,
		opts *api.NodePeersOpts,
	) (
		*api.Response[[]*apiv1.Peer],
		error,
	)
}

// NodeSyncingProvider is the interface for providing synchronization state.
type NodeSyncingProvider interface {
	// NodeSyncing provides the state of the node's synchronization with the chain.
	NodeSyncing(ctx context.Context,
		opts *api.NodeSyncingOpts,
	) (
		*api.Response[*apiv1.SyncState],
		error,
	)
}

// ValidatorLivenessProvider is the interface for providing validator liveness data.
type ValidatorLivenessProvider interface {
	// ValidatorLiveness provides the liveness data to the given validators.
	ValidatorLiveness(ctx context.Context,
		opts *api.ValidatorLivenessOpts,
	) (
		*api.Response[[]*apiv1.ValidatorLiveness],
		error,
	)
}

// NodeVersionProvider is the interface for providing the node version.
type NodeVersionProvider interface {
	// NodeVersion returns a free-text string with the node version.
	NodeVersion(ctx context.Context,
		opts *api.NodeVersionOpts,
	) (
		*api.Response[string],
		error,
	)
}

// ProposerDutiesProvider is the interface for providing proposer duties.
type ProposerDutiesProvider interface {
	// ProposerDuties obtains proposer duties for the given options.
	ProposerDuties(ctx context.Context,
		opts *api.ProposerDutiesOpts,
	) (
		*api.Response[[]*apiv1.ProposerDuty],
		error,
	)
}

// SpecProvider is the interface for providing spec data.
type SpecProvider interface {
	// Spec provides the spec information of the chain.
	Spec(ctx context.Context,
		opts *api.SpecOpts,
	) (
		*api.Response[map[string]any],
		error,
	)
}

// SyncStateProvider is the interface for providing synchronization state.
type SyncStateProvider interface {
	// SyncState provides the state of the node's synchronization with the chain.
	//
	// Deprecated: use NodeSyncing()
	SyncState(ctx context.Context) (*apiv1.SyncState, error)
}

// ValidatorBalancesProvider is the interface for providing validator balances.
type ValidatorBalancesProvider interface {
	// ValidatorBalances provides the validator balances for the given options.
	ValidatorBalances(ctx context.Context,
		opts *api.ValidatorBalancesOpts,
	) (
		*api.Response[map[capella.ValidatorIndex]capella.Gwei],
		error,
	)
}

// ValidatorsProvider is the interface for providing validator information.
type ValidatorsProvider interface {
	// Validators provides the validators, with their balance and status, for the given options.
	Validators(ctx context.Context,
		opts *api.ValidatorsOpts,
	) (
		*api.Response[map[capella.ValidatorIndex]*apiv1.Validator],
		error,
	)
}

// VoluntaryExitPoolProvider is the interface for providing voluntary exit pools.
type VoluntaryExitPoolProvider interface {
	// VoluntaryExitPool fetches the voluntary exit pool.
	VoluntaryExitPool(ctx context.Context,
		opts *api.VoluntaryExitPoolOpts,
	) (
		*api.Response[[]*capella.SignedVoluntaryExit],
		error,
	)
}

//
// Local extensions
//

// DomainProvider provides a domain for a given domain type at an epoch.
type DomainProvider interface {
	// Domain provides a domain for a given domain type at a given epoch.
	Domain(ctx context.Context, domainType capella.DomainType, epoch capella.Epoch) (capella.Domain, error)

	// GenesisDomain returns the domain for the given domain type at genesis.
	// N.B. this is not always the same as the domain at epoch 0.  It is possible
	// for a chain's fork schedule to have multiple forks at genesis.  In this situation,
	// GenesisDomain() will return the first, and Domain() will return the last.
	GenesisDomain(ctx context.Context, domainType capella.DomainType) (capella.Domain, error)
}

// GenesisTimeProvider is the interface for providing the genesis time of a chain.
//
// Deprecated: use Genesis().
type GenesisTimeProvider interface {
	// GenesisTime provides the genesis time of the chain.
	GenesisTime(ctx context.Context) (time.Time, error)
}

// NodeClientProvider provides the client for the node.
type NodeClientProvider interface {
	// NodeClient provides the client for the node.
	NodeClient(ctx context.Context) (*api.Response[string], error)
}
