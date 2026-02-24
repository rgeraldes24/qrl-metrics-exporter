// Copyright © 2021 - 2023 Attestant Limited.
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

// Sleepy is a QRL client that sleeps for a random amount of time within a
// set of bounds before continuing.
type Sleepy struct {
	minSleep time.Duration
	maxSleep time.Duration
	next     consensusclient.Service
}

// NewSleepy creates a new QRL client that sleeps for random amount of time
// within a set of bounds between minSleep and maxSleep before continuing.
func NewSleepy(_ context.Context,
	minSleep time.Duration,
	maxSleep time.Duration,
	next consensusclient.Service,
) (consensusclient.Service, error) {
	if next == nil {
		return nil, errors.New("no next service supplied")
	}

	if maxSleep < minSleep {
		return nil, errors.New("max sleep less than min sleep")
	}

	return &Sleepy{
		minSleep: minSleep,
		maxSleep: maxSleep,
		next:     next,
	}, nil
}

// Name returns the name of the client implementation.
func (s *Sleepy) Name() string {
	nextName := s.next.Name()

	return fmt.Sprintf("sleepy(%v,%v,%s)", s.minSleep, s.maxSleep, nextName)
}

// Address returns the address of the client.
func (s *Sleepy) Address() string {
	nextAddress := s.next.Address()

	return fmt.Sprintf("sleepy:%v,%v,%s", s.minSleep, s.maxSleep, nextAddress)
}

// IsActive returns true if the client is active.
func (*Sleepy) IsActive() bool {
	return true
}

// IsSynced returns true if the client is synced.
func (*Sleepy) IsSynced() bool {
	return true
}

// sleep sleeps for a bounded amount of time.
func (s *Sleepy) sleep(_ context.Context) {
	duration := time.Duration(s.minSleep.Milliseconds()+
		// #nosec G404
		rand.Int63n(s.maxSleep.Milliseconds()-s.minSleep.Milliseconds())) * time.Millisecond
	time.Sleep(duration)
}

// EpochFromStateID converts a state ID to its epoch.
//
// Deprecated: use chaintime.
func (s *Sleepy) EpochFromStateID(ctx context.Context, stateID string) (capella.Epoch, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.EpochFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.EpochFromStateID(ctx, stateID)
}

// SlotFromStateID converts a state ID to its slot.
//
// Deprecated: use chaintime.
func (s *Sleepy) SlotFromStateID(ctx context.Context, stateID string) (capella.Slot, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.SlotFromStateIDProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.SlotFromStateID(ctx, stateID)
}

// NodeVersion returns a free-text string with the node version.
func (s *Sleepy) NodeVersion(ctx context.Context,
	opts *api.NodeVersionOpts,
) (
	*api.Response[string],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.NodeVersionProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.NodeVersion(ctx, opts)
}

// SlotDuration provides the duration of a slot of the chain.
//
// Deprecated: use Spec().
func (s *Sleepy) SlotDuration(ctx context.Context) (time.Duration, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.SlotDurationProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.SlotDuration(ctx)
}

// SlotsPerEpoch provides the slots per epoch of the chain.
//
// Deprecated: use Spec().
func (s *Sleepy) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.SlotsPerEpochProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.SlotsPerEpoch(ctx)
}

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Sleepy) FarFutureEpoch(ctx context.Context) (capella.Epoch, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.FarFutureEpochProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.FarFutureEpoch(ctx)
}

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
//
// Deprecated: use Spec().
func (s *Sleepy) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.TargetAggregatorsPerCommitteeProvider)
	if !isNext {
		return 0, errors.New("next does not support this call")
	}

	return next.TargetAggregatorsPerCommittee(ctx)
}

// AggregateAttestation fetches the aggregate attestation for the given options.
func (s *Sleepy) AggregateAttestation(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*spec.VersionedAttestation],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.AggregateAttestationProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.AggregateAttestation(ctx, opts)
}

// AttestationData fetches the attestation data for the given slot and committee index.
func (s *Sleepy) AttestationData(ctx context.Context,
	opts *api.AttestationDataOpts,
) (
	*api.Response[*capella.AttestationData],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.AttestationDataProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.AttestationData(ctx, opts)
}

// AttestationPool fetches the attestation pool for the given slot.
func (s *Sleepy) AttestationPool(ctx context.Context,
	opts *api.AttestationPoolOpts,
) (
	*api.Response[[]*spec.VersionedAttestation],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.AttestationPoolProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.AttestationPool(ctx, opts)
}

// AttesterDuties obtains attester duties.
// If validatorIndices is nil it will return all duties for the given epoch.
func (s *Sleepy) AttesterDuties(ctx context.Context,
	opts *api.AttesterDutiesOpts,
) (
	*api.Response[[]*apiv1.AttesterDuty],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.AttesterDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.AttesterDuties(ctx, opts)
}

// BeaconBlockHeader provides the block header of a given block ID.
func (s *Sleepy) BeaconBlockHeader(ctx context.Context,
	opts *api.BeaconBlockHeaderOpts,
) (
	*api.Response[*apiv1.BeaconBlockHeader],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.BeaconBlockHeadersProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.BeaconBlockHeader(ctx, opts)
}

// Proposal fetches a proposal for signing.
func (s *Sleepy) Proposal(ctx context.Context,
	opts *api.ProposalOpts,
) (
	*api.Response[*api.VersionedProposal],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ProposalProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Proposal(ctx, opts)
}

// BeaconState fetches a beacon state.
func (s *Sleepy) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.BeaconStateProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.BeaconState(ctx, opts)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Sleepy) Events(ctx context.Context, opts *api.EventsOpts) error {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.EventsProvider)
	if !isNext {
		return errors.New("next does not support this call")
	}

	return next.Events(ctx, opts)
}

// Finality provides the finality given a state ID.
func (s *Sleepy) Finality(ctx context.Context,
	opts *api.FinalityOpts,
) (
	*api.Response[*apiv1.Finality],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.FinalityProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Finality(ctx, opts)
}

// Fork fetches fork information for the given state.
func (s *Sleepy) Fork(ctx context.Context,
	opts *api.ForkOpts,
) (
	*api.Response[*capella.Fork],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ForkProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Fork(ctx, opts)
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Sleepy) ForkSchedule(ctx context.Context,
	opts *api.ForkScheduleOpts,
) (
	*api.Response[[]*capella.Fork],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ForkScheduleProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.ForkSchedule(ctx, opts)
}

// Genesis fetches genesis information for the chain.
func (s *Sleepy) Genesis(ctx context.Context,
	opts *api.GenesisOpts,
) (
	*api.Response[*apiv1.Genesis],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.GenesisProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Genesis(ctx, opts)
}

// NodeSyncing provides the state of the node's synchronization with the chain.
func (s *Sleepy) NodeSyncing(ctx context.Context,
	opts *api.NodeSyncingOpts,
) (
	*api.Response[*apiv1.SyncState],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.NodeSyncingProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.NodeSyncing(ctx, opts)
}

// NodePeers provides the peers of the node.
func (s *Sleepy) NodePeers(ctx context.Context,
	opts *api.NodePeersOpts,
) (
	*api.Response[[]*apiv1.Peer],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.NodePeersProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.NodePeers(ctx, opts)
}

// ProposerDuties obtains proposer duties for the given epoch.
func (s *Sleepy) ProposerDuties(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*apiv1.ProposerDuty],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ProposerDutiesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.ProposerDuties(ctx, opts)
}

// Spec provides the spec information of the chain.
func (s *Sleepy) Spec(ctx context.Context,
	opts *api.SpecOpts,
) (
	*api.Response[map[string]any],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.SpecProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Spec(ctx, opts)
}

// ValidatorBalances provides the validator balances for a given state.
func (s *Sleepy) ValidatorBalances(ctx context.Context,
	opts *api.ValidatorBalancesOpts,
) (
	*api.Response[map[capella.ValidatorIndex]capella.Gwei],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ValidatorBalancesProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.ValidatorBalances(ctx, opts)
}

// Validators provides the validators, with their balance and status, for a given state.
func (s *Sleepy) Validators(ctx context.Context,
	opts *api.ValidatorsOpts,
) (
	*api.Response[map[capella.ValidatorIndex]*apiv1.Validator],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ValidatorsProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.Validators(ctx, opts)
}

// VoluntaryExitPool fetches the voluntary exit pool.
func (s *Sleepy) VoluntaryExitPool(ctx context.Context,
	opts *api.VoluntaryExitPoolOpts,
) (
	*api.Response[[]*capella.SignedVoluntaryExit],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.VoluntaryExitPoolProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.VoluntaryExitPool(ctx, opts)
}

// Domain provides a domain for a given domain type at a given epoch.
func (s *Sleepy) Domain(ctx context.Context, domainType capella.DomainType, epoch capella.Epoch) (capella.Domain, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return capella.Domain{}, errors.New("next does not support this call")
	}

	return next.Domain(ctx, domainType, epoch)
}

// GenesisDomain provides a domain for a given domain type at genesis.
func (s *Sleepy) GenesisDomain(ctx context.Context, domainType capella.DomainType) (capella.Domain, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.DomainProvider)
	if !isNext {
		return capella.Domain{}, errors.New("next does not support this call")
	}

	return next.GenesisDomain(ctx, domainType)
}

// GenesisTime provides the genesis time of the chain.
//
// Deprecated: use Genesis().
func (s *Sleepy) GenesisTime(ctx context.Context) (time.Time, error) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.GenesisTimeProvider)
	if !isNext {
		return time.Time{}, errors.New("next does not support this call")
	}

	return next.GenesisTime(ctx)
}

// ForkChoice fetches the node's current fork choice context.
func (s *Sleepy) ForkChoice(ctx context.Context,
	opts *api.ForkChoiceOpts,
) (
	*api.Response[*apiv1.ForkChoice],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ForkChoiceProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.ForkChoice(ctx, opts)
}

// AttestationRewards provides rewards to the given validators for attesting.
func (s *Sleepy) AttestationRewards(ctx context.Context,
	opts *api.AttestationRewardsOpts,
) (
	*api.Response[*apiv1.AttestationRewards],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.AttestationRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.AttestationRewards(ctx, opts)
}

// BlockRewards provides rewards for proposing a block.
func (s *Sleepy) BlockRewards(ctx context.Context,
	opts *api.BlockRewardsOpts,
) (
	*api.Response[*apiv1.BlockRewards],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.BlockRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.BlockRewards(ctx, opts)
}

// SyncCommitteeRewards provides rewards to the given validators for being members of a sync committee.
func (s *Sleepy) SyncCommitteeRewards(ctx context.Context,
	opts *api.SyncCommitteeRewardsOpts,
) (
	*api.Response[[]*apiv1.SyncCommitteeReward],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.SyncCommitteeRewardsProvider)
	if !isNext {
		return nil, fmt.Errorf("%s@%s does not support this call", s.next.Name(), s.next.Address())
	}

	return next.SyncCommitteeRewards(ctx, opts)
}

// ValidatorLiveness provides the liveness data to the given validators.
func (s *Sleepy) ValidatorLiveness(ctx context.Context,
	opts *api.ValidatorLivenessOpts,
) (
	*api.Response[[]*apiv1.ValidatorLiveness],
	error,
) {
	s.sleep(ctx)

	next, isNext := s.next.(consensusclient.ValidatorLivenessProvider)
	if !isNext {
		return nil, errors.New("next does not support this call")
	}

	return next.ValidatorLiveness(ctx, opts)
}
