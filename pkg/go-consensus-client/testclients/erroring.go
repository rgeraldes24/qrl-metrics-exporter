// Copyright © 2021 - 2025 Attestant Limited.
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

package testclients

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	apiv1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// Erroring is a QRL client that errors at a given rate.
type Erroring struct {
	errorRate float64
	next      consensusclient.Service
}

// NewErroring creates a new QRL client that errors at a given rate.
func NewErroring(_ context.Context,
	errorRate float64,
	next consensusclient.Service,
) (consensusclient.Service, error) {
	if next == nil {
		return nil, errors.New("no next service supplied")
	}

	if errorRate < 0 {
		return nil, errors.New("error rate cannot be less than 0")
	}

	if errorRate > 1 {
		return nil, errors.New("error rate cannot be more than 1")
	}

	return &Erroring{
		errorRate: errorRate,
		next:      next,
	}, nil
}

// Name returns the name of the client implementation.
func (s *Erroring) Name() string {
	nextName := s.next.Name()

	return fmt.Sprintf("erroring(%v,%s)", s.errorRate, nextName)
}

// Address returns the address of the client.
func (s *Erroring) Address() string {
	nextAddress := s.next.Address()

	return fmt.Sprintf("erroring:%v,%s", s.errorRate, nextAddress)
}

// IsActive returns true if the client is active.
func (*Erroring) IsActive() bool {
	return true
}

// IsSynced returns true if the client is synced.
func (*Erroring) IsSynced() bool {
	return true
}

// maybeError may return an error depending on the error arte.
func (s *Erroring) maybeError(_ context.Context) error {
	// #nosec G404
	roll := rand.Float64()
	if roll < s.errorRate {
		return errors.New("error")
	}

	return nil
}

// EpochFromStateID converts a state ID to its epoch.
//
// Deprecated: use chaintime.
func (s *Erroring) EpochFromStateID(ctx context.Context, stateID string) (capella.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.EpochFromStateIDProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
//
// Deprecated: use chaintime.
func (s *Erroring) SlotFromStateID(ctx context.Context, stateID string) (capella.Slot, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.SlotFromStateIDProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SlotFromStateID(ctx, stateID)
}

// NodeVersion returns a free-text string with the node version.
func (s *Erroring) NodeVersion(ctx context.Context,
	opts *api.NodeVersionOpts,
) (
	*api.Response[string],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.NodeVersionProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.NodeVersion(ctx, opts)
}

// SlotDuration provides the duration of a slot of the chain.
//
// Deprecated: use Spec().
func (s *Erroring) SlotDuration(ctx context.Context) (time.Duration, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.SlotDurationProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SlotDuration(ctx)
}

// SlotsPerEpoch provides the slots per epoch of the chain.
//
// Deprecated: use Spec().
func (s *Erroring) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.SlotsPerEpochProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Erroring) FarFutureEpoch(ctx context.Context) (capella.Epoch, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.FarFutureEpochProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.FarFutureEpoch(ctx)
}

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
//
// Deprecated: use Spec().
func (s *Erroring) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	if err := s.maybeError(ctx); err != nil {
		return 0, err
	}

	next, isNext := s.next.(consensusclient.TargetAggregatorsPerCommitteeProvider)
	if !isNext {
		return 0, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.TargetAggregatorsPerCommittee(ctx)
}

// AggregateAttestation fetches the aggregate attestation for the given options.
func (s *Erroring) AggregateAttestation(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*spec.VersionedAttestation],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.AggregateAttestationProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AggregateAttestation(ctx, opts)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Erroring) AttestationData(ctx context.Context,
	opts *api.AttestationDataOpts,
) (
	*api.Response[*capella.AttestationData],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.AttestationDataProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AttestationData(ctx, opts)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Erroring) AttestationPool(ctx context.Context,
	opts *api.AttestationPoolOpts,
) (
	*api.Response[[]*spec.VersionedAttestation],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.AttestationPoolProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AttestationPool(ctx, opts)
}

// AttesterDuties obtains attester duties.
// If validatorIndices is nil it will return all duties for the given epoch.
func (s *Erroring) AttesterDuties(ctx context.Context,
	opts *api.AttesterDutiesOpts,
) (
	*api.Response[[]*apiv1.AttesterDuty],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.AttesterDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AttesterDuties(ctx, opts)
}

// BeaconBlockHeader provides the block header of a given block ID.
func (s *Erroring) BeaconBlockHeader(ctx context.Context,
	opts *api.BeaconBlockHeaderOpts,
) (
	*api.Response[*apiv1.BeaconBlockHeader],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BeaconBlockHeadersProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BeaconBlockHeader(ctx, opts)
}

// BeaconBlockRoot fetches a block's root given a block ID.
func (s *Erroring) BeaconBlockRoot(ctx context.Context,
	opts *api.BeaconBlockRootOpts,
) (
	*api.Response[*capella.Root],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BeaconBlockRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BeaconBlockRoot(ctx, opts)
}

// BeaconCommittees fetches all beacon committees for the epoch at the given state.
func (s *Erroring) BeaconCommittees(ctx context.Context,
	opts *api.BeaconCommitteesOpts,
) (
	*api.Response[[]*apiv1.BeaconCommittee],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BeaconCommitteesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BeaconCommittees(ctx, opts)
}

// Proposal fetches a proposal for signing.
func (s *Erroring) Proposal(ctx context.Context,
	opts *api.ProposalOpts,
) (
	*api.Response[*api.VersionedProposal],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ProposalProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Proposal(ctx, opts)
}

// BeaconState fetches a beacon state.
func (s *Erroring) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BeaconStateProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BeaconState(ctx, opts)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Erroring) Events(ctx context.Context, opts *api.EventsOpts) error {
	if err := s.maybeError(ctx); err != nil {
		return err
	}

	next, isNext := s.next.(consensusclient.EventsProvider)
	if !isNext {
		return fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Events(ctx, opts)
}

// Finality provides the finality given a state ID.
func (s *Erroring) Finality(ctx context.Context,
	opts *api.FinalityOpts,
) (
	*api.Response[*apiv1.Finality],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.FinalityProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Finality(ctx, opts)
}

// Fork fetches fork information for the given state.
func (s *Erroring) Fork(ctx context.Context,
	opts *api.ForkOpts,
) (
	*api.Response[*capella.Fork],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ForkProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Fork(ctx, opts)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Erroring) ForkSchedule(ctx context.Context,
	opts *api.ForkScheduleOpts,
) (
	*api.Response[[]*capella.Fork],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ForkScheduleProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.ForkSchedule(ctx, opts)
}

// Genesis fetches genesis information for the chain.
func (s *Erroring) Genesis(ctx context.Context,
	opts *api.GenesisOpts,
) (
	*api.Response[*apiv1.Genesis],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.GenesisProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Genesis(ctx, opts)
}

// NodeSyncing provides the state of the node's synchronization with the chain.
func (s *Erroring) NodeSyncing(ctx context.Context,
	opts *api.NodeSyncingOpts,
) (
	*api.Response[*apiv1.SyncState],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.NodeSyncingProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.NodeSyncing(ctx, opts)
}

// NodePeers provides the peers of the node.
func (s *Erroring) NodePeers(ctx context.Context,
	opts *api.NodePeersOpts,
) (
	*api.Response[[]*apiv1.Peer],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.NodePeersProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.NodePeers(ctx, opts)
}

// ProposerDuties obtains proposer duties for the given epoch.
func (s *Erroring) ProposerDuties(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*apiv1.ProposerDuty],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ProposerDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.ProposerDuties(ctx, opts)
}

// SyncCommittee fetches the sync committee for the given state.
func (s *Erroring) SyncCommittee(ctx context.Context,
	opts *api.SyncCommitteeOpts,
) (
	*api.Response[*apiv1.SyncCommittee],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SyncCommitteesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SyncCommittee(ctx, opts)
}

// SyncCommitteeContribution provides a sync committee contribution.
func (s *Erroring) SyncCommitteeContribution(ctx context.Context,
	opts *api.SyncCommitteeContributionOpts,
) (
	*api.Response[*capella.SyncCommitteeContribution],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SyncCommitteeContributionProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SyncCommitteeContribution(ctx, opts)
}

// SyncCommitteeDuties obtains sync committee duties.
// If validatorIndices is nil it will return all duties for the given epoch.
func (s *Erroring) SyncCommitteeDuties(ctx context.Context,
	opts *api.SyncCommitteeDutiesOpts,
) (
	*api.Response[[]*apiv1.SyncCommitteeDuty],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SyncCommitteeDutiesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SyncCommitteeDuties(ctx, opts)
}

// Spec provides the spec information of the chain.
func (s *Erroring) Spec(ctx context.Context,
	opts *api.SpecOpts,
) (
	*api.Response[map[string]any],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SpecProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Spec(ctx, opts)
}

// ValidatorBalances provides the validator balances for a given state.
func (s *Erroring) ValidatorBalances(ctx context.Context,
	opts *api.ValidatorBalancesOpts,
) (
	*api.Response[map[capella.ValidatorIndex]capella.Gwei],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ValidatorBalancesProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.ValidatorBalances(ctx, opts)
}

// Validators provides the validators, with their balance and status, for a given state.
func (s *Erroring) Validators(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[capella.ValidatorIndex]*apiv1.Validator],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Validators(ctx, opts)
}

// VoluntaryExitPool fetches the voluntary exit pool.
func (s *Erroring) VoluntaryExitPool(ctx context.Context,
	opts *api.VoluntaryExitPoolOpts,
) (
	*api.Response[[]*capella.SignedVoluntaryExit],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.VoluntaryExitPoolProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.VoluntaryExitPool(ctx, opts)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Erroring) Domain(ctx context.Context, domainType capella.DomainType, epoch capella.Epoch) (capella.Domain, error) {
	if err := s.maybeError(ctx); err != nil {
		return capella.Domain{}, err
	}

	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return capella.Domain{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.Domain(ctx, domainType, epoch)
}

// GenesisDomain provides a domain for a given domain type.
func (s *Erroring) GenesisDomain(ctx context.Context, domainType capella.DomainType) (capella.Domain, error) {
	if err := s.maybeError(ctx); err != nil {
		return capella.Domain{}, err
	}

	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return capella.Domain{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.GenesisDomain(ctx, domainType)
}

// GenesisTime provides the genesis time of the chain.
//
// Deprecated: use Genesis().
func (s *Erroring) GenesisTime(ctx context.Context) (time.Time, error) {
	if err := s.maybeError(ctx); err != nil {
		return time.Time{}, err
	}

	next, isNext := s.next.(consensusclient.GenesisTimeProvider)
	if !isNext {
		return time.Time{}, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.GenesisTime(ctx)
}

// DepositContract provides details of the QRL deposit contract for the chain.
func (s *Erroring) DepositContract(ctx context.Context,
	opts *api.DepositContractOpts,
) (
	*api.Response[*apiv1.DepositContract],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.DepositContractProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.DepositContract(ctx, opts)
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Erroring) SignedBeaconBlock(ctx context.Context,
	opts *api.SignedBeaconBlockOpts,
) (
	*api.Response[*spec.VersionedSignedBeaconBlock],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SignedBeaconBlockProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SignedBeaconBlock(ctx, opts)
}

// BeaconStateRoot fetches a beacon state root given a state ID.
func (s *Erroring) BeaconStateRoot(ctx context.Context,
	opts *api.BeaconStateRootOpts,
) (
	*api.Response[*capella.Root],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BeaconStateRootProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BeaconStateRoot(ctx, opts)
}

// ForkChoice fetches the node's current fork choice context.
func (s *Erroring) ForkChoice(ctx context.Context,
	opts *api.ForkChoiceOpts,
) (
	*api.Response[*apiv1.ForkChoice],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ForkChoiceProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.ForkChoice(ctx, opts)
}

// AttestationRewards provides rewards to the given validators for attesting.
func (s *Erroring) AttestationRewards(ctx context.Context,
	opts *api.AttestationRewardsOpts,
) (
	*api.Response[*apiv1.AttestationRewards],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.AttestationRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AttestationRewards(ctx, opts)
}

// BlockRewards provides rewards for proposing a block.
func (s *Erroring) BlockRewards(ctx context.Context,
	opts *api.BlockRewardsOpts,
) (
	*api.Response[*apiv1.BlockRewards],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.BlockRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BlockRewards(ctx, opts)
}

// SyncCommitteeRewards provides rewards to the given validators for being members of a sync committee.
func (s *Erroring) SyncCommitteeRewards(ctx context.Context,
	opts *api.SyncCommitteeRewardsOpts,
) (
	*api.Response[[]*apiv1.SyncCommitteeReward],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.SyncCommitteeRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SyncCommitteeRewards(ctx, opts)
}

// ValidatorLiveness provides the liveness data to the given validators.
func (s *Erroring) ValidatorLiveness(ctx context.Context,
	opts *api.ValidatorLivenessOpts,
) (
	*api.Response[[]*apiv1.ValidatorLiveness],
	error,
) {
	if err := s.maybeError(ctx); err != nil {
		return nil, err
	}

	next, isNext := s.next.(consensusclient.ValidatorLivenessProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.ValidatorLiveness(ctx, opts)
}
