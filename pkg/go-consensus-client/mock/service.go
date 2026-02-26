// Copyright © 2020 - 2025 Attestant Limited.
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

package mock

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	apiv1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// Service is a mock QRL client service, providing data locally.
//
//nolint:revive
type Service struct {
	name    string
	timeout time.Duration

	genesisTime time.Time

	// Various information from the node that does not change during the
	// lifetime of a beacon node.
	// genesis         *api.Genesis
	// spec            map[string]interface{}
	// depositContract *api.DepositContract
	// forkSchedule    []*capella.Fork
	nodeVersion string

	// Values that can be altered if required.
	HeadSlot     capella.Slot
	SyncDistance capella.Slot

	// Functions that can be provided to mock specific responses from this client.
	AggregateAttestationFunc      func(context.Context, *api.AggregateAttestationOpts) (*api.Response[*spec.VersionedAttestation], error)
	AttesterDutiesFunc            func(context.Context, *api.AttesterDutiesOpts) (*api.Response[[]*apiv1.AttesterDuty], error)
	AttestationDataFunc           func(context.Context, *api.AttestationDataOpts) (*api.Response[*capella.AttestationData], error)
	AttestationRewardsFunc        func(context.Context, *api.AttestationRewardsOpts) (*api.Response[*apiv1.AttestationRewards], error)
	BeaconBlockHeaderFunc         func(context.Context, *api.BeaconBlockHeaderOpts) (*api.Response[*apiv1.BeaconBlockHeader], error)
	BeaconBlockRootFunc           func(context.Context, *api.BeaconBlockRootOpts) (*api.Response[*capella.Root], error)
	BeaconStateFunc               func(context.Context, *api.BeaconStateOpts) (*api.Response[*spec.VersionedBeaconState], error)
	BeaconStateRandaoFunc         func(context.Context, *api.BeaconStateRandaoOpts) (*api.Response[*capella.Root], error)
	BeaconStateRootFunc           func(context.Context, *api.BeaconStateRootOpts) (*api.Response[*capella.Root], error)
	BlockRewardsFunc              func(context.Context, *api.BlockRewardsOpts) (*api.Response[*apiv1.BlockRewards], error)
	DepositContractFunc           func(context.Context, *api.DepositContractOpts) (*api.Response[*apiv1.DepositContract], error)
	EventsFunc                    func(context.Context, *api.EventsOpts) error
	FinalityFunc                  func(context.Context, *api.FinalityOpts) (*api.Response[*apiv1.Finality], error)
	ForkChoiceFunc                func(context.Context, *api.ForkChoiceOpts) (*api.Response[*apiv1.ForkChoice], error)
	ForkFunc                      func(context.Context, *api.ForkOpts) (*api.Response[*capella.Fork], error)
	ForkScheduleFunc              func(context.Context, *api.ForkScheduleOpts) (*api.Response[[]*capella.Fork], error)
	GenesisFunc                   func(context.Context, *api.GenesisOpts) (*api.Response[*apiv1.Genesis], error)
	NodePeersFunc                 func(context.Context, *api.NodePeersOpts) (*api.Response[[]*apiv1.Peer], error)
	NodeSyncingFunc               func(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error)
	NodeVersionFunc               func(context.Context, *api.NodeVersionOpts) (*api.Response[string], error)
	BlindedProposalFunc           func(context.Context, *api.BlindedProposalOpts) (*api.Response[*api.VersionedBlindedProposal], error)
	ProposalFunc                  func(context.Context, *api.ProposalOpts) (*api.Response[*api.VersionedProposal], error)
	ProposerDutiesFunc            func(context.Context, *api.ProposerDutiesOpts) (*api.Response[[]*apiv1.ProposerDuty], error)
	SignedBeaconBlockFunc         func(context.Context, *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error)
	SpecFunc                      func(context.Context, *api.SpecOpts) (*api.Response[map[string]any], error)
	SyncCommitteeContributionFunc func(context.Context, *api.SyncCommitteeContributionOpts) (*api.Response[*capella.SyncCommitteeContribution], error)
	SyncCommitteeDutiesFunc       func(context.Context, *api.SyncCommitteeDutiesOpts) (*api.Response[[]*apiv1.SyncCommitteeDuty], error)
	SyncCommitteeRewardsFunc      func(context.Context, *api.SyncCommitteeRewardsOpts) (*api.Response[[]*apiv1.SyncCommitteeReward], error)
	ValidatorBalancesFunc         func(context.Context, *api.ValidatorBalancesOpts) (*api.Response[map[capella.ValidatorIndex]capella.Shor], error)
	ValidatorLivenessFunc         func(context.Context, *api.ValidatorLivenessOpts) (*api.Response[[]*apiv1.ValidatorLiveness], error)
	ValidatorsFunc                func(context.Context, *api.ValidatorsOpts) (*api.Response[map[capella.ValidatorIndex]*apiv1.Validator], error)
	VoluntaryExitPoolFunc         func(context.Context, *api.VoluntaryExitPoolOpts) (*api.Response[[]*capella.SignedVoluntaryExit], error)
}

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new QRL client service, mocking connections.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "mock").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	s := &Service{
		name:        parameters.name,
		genesisTime: parameters.genesisTime,
		timeout:     parameters.timeout,
		nodeVersion: "mock",

		HeadSlot:     12345,
		SyncDistance: 0,
	}

	// Fetch static values to confirm the connection is good.
	if err := s.fetchStaticValues(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to confirm node connection")
	}

	// Close the service on context done.
	go func(*Service) {
		<-ctx.Done()
		log.Trace().Msg("Context done; closing connection")
		s.close()
	}(s)

	return s, nil
}

// fetchStaticValues fetches values that never change.
// This caches the values, avoiding future API calls.

// Name provides the name of the service.
func (*Service) Name() string {
	return "Mock"
}

// Address provides the address of the service.
func (s *Service) Address() string {
	return s.name
}

// IsActive returns true if the client is active.
func (*Service) IsActive() bool {
	return true
}

// IsSynced returns true if the client is synced.
func (s *Service) IsSynced() bool {
	return s.SyncDistance == 0
}

// fetchStaticValues fetches values that never change.
// This caches the values, avoiding future API calls.
func (s *Service) fetchStaticValues(ctx context.Context) error {
	if _, err := s.Genesis(ctx, &api.GenesisOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch genesis")
	}
	//	if _, err := s.Spec(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch spec")
	//	}
	//	if _, err := s.DepositContract(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch deposit contract")
	//	}
	//	if _, err := s.ForkSchedule(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch fork schedule")
	//	}

	return nil
}

// close closes the service, freeing up resources.
func (*Service) close() {
}
