package beacon

import (
	"context"
	"errors"

	"github.com/theQRL/qrl-metrics-exporter/pkg/beacon/api/types"
	"github.com/theQRL/qrl-metrics-exporter/pkg/beacon/state"
	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	v1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

func (n *node) FetchSyncStatus(ctx context.Context) (*v1.SyncState, error) {
	provider, isProvider := n.client.(consensusclient.NodeSyncingProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.NodeSyncingProvider")
	}

	status, err := provider.NodeSyncing(ctx, &api.NodeSyncingOpts{})
	if err != nil {
		return nil, err
	}

	n.stat.UpdateSyncState(status.Data)

	n.publishSyncStatus(ctx, status.Data)

	return status.Data, nil
}

func (n *node) FetchPeers(ctx context.Context) (*types.Peers, error) {
	peers, err := n.api.NodePeers(ctx)
	if err != nil {
		return nil, err
	}

	n.peers = peers

	n.publishPeersUpdated(ctx, peers)

	return &peers, nil
}

func (n *node) FetchNodeVersion(ctx context.Context) (string, error) {
	provider, isProvider := n.client.(consensusclient.NodeVersionProvider)
	if !isProvider {
		return "", errors.New("client does not implement consensusclient.NodeVersionProvider")
	}

	rsp, err := provider.NodeVersion(ctx, &api.NodeVersionOpts{})
	if err != nil {
		return "", err
	}

	n.nodeVersionMu.Lock()
	n.nodeVersion = rsp.Data
	n.nodeVersionMu.Unlock()

	n.publishNodeVersionUpdated(ctx, rsp.Data)

	return rsp.Data, nil
}

func (n *node) FetchBlock(ctx context.Context, stateID string) (*spec.VersionedSignedBeaconBlock, error) {
	return n.getBlock(ctx, stateID)
}

func (n *node) FetchRawBlock(ctx context.Context, stateID string, contentType string) ([]byte, error) {
	return n.api.RawBlock(ctx, stateID, contentType)
}

func (n *node) FetchBlockRoot(ctx context.Context, stateID string) (*capella.Root, error) {
	return n.getBlockRoot(ctx, stateID)
}

func (n *node) FetchBeaconState(ctx context.Context, stateID string) (*spec.VersionedBeaconState, error) {
	provider, isProvider := n.client.(consensusclient.BeaconStateProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.NodeVersionProvider")
	}

	rsp, err := provider.BeaconState(ctx, &api.BeaconStateOpts{
		State: stateID,
	})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchRawBeaconState(ctx context.Context, stateID string, contentType string) ([]byte, error) {
	return n.api.RawDebugBeaconState(ctx, stateID, contentType)
}

func (n *node) FetchFinality(ctx context.Context, stateID string) (*v1.Finality, error) {
	provider, isProvider := n.client.(consensusclient.FinalityProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.FinalityProvider")
	}

	rsp, err := provider.Finality(ctx, &api.FinalityOpts{
		State: stateID,
	})
	if err != nil {
		return nil, err
	}

	finality := rsp.Data

	if stateID == "head" {
		changed := false
		if n.finality == nil ||
			finality.Finalized.Root != n.finality.Finalized.Root ||
			finality.Finalized.Epoch != n.finality.Finalized.Epoch ||
			finality.Justified.Root != n.finality.Justified.Root ||
			finality.Justified.Epoch != n.finality.Justified.Epoch ||
			finality.PreviousJustified.Epoch != n.finality.PreviousJustified.Epoch ||
			finality.PreviousJustified.Root != n.finality.PreviousJustified.Root {
			changed = true
		}

		n.finality = finality

		if changed {
			n.publishFinalityCheckpointUpdated(ctx, finality)
		}
	}

	return finality, nil
}

func (n *node) FetchRawSpec(ctx context.Context) (map[string]any, error) {
	provider, isProvider := n.client.(consensusclient.SpecProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.SpecProvider")
	}

	rsp, err := provider.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchSpec(ctx context.Context) (*state.Spec, error) {
	provider, isProvider := n.client.(consensusclient.SpecProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.SpecProvider")
	}

	rsp, err := provider.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return nil, err
	}

	sp := state.NewSpec(rsp.Data)

	n.specMu.Lock()
	n.spec = &sp
	n.specMu.Unlock()

	n.publishSpecUpdated(ctx, &sp)

	return &sp, nil
}

func (n *node) FetchProposerDuties(ctx context.Context, epoch capella.Epoch) ([]*v1.ProposerDuty, error) {
	n.log.WithField("epoch", epoch).Debug("Fetching proposer duties")

	provider, isProvider := n.client.(consensusclient.ProposerDutiesProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.ProposerDutiesProvider")
	}

	rsp, err := provider.ProposerDuties(ctx, &api.ProposerDutiesOpts{
		Epoch: epoch,
	})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchForkChoice(ctx context.Context) (*v1.ForkChoice, error) {
	provider, isProvider := n.client.(consensusclient.ForkChoiceProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.ForkChoiceProvider")
	}

	rsp, err := provider.ForkChoice(ctx, &api.ForkChoiceOpts{})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchDepositSnapshot(ctx context.Context) (*types.DepositSnapshot, error) {
	return n.api.DepositSnapshot(ctx)
}

func (n *node) FetchNodeIdentity(ctx context.Context) (*types.Identity, error) {
	return n.api.NodeIdentity(ctx)
}

func (n *node) FetchBeaconStateRoot(ctx context.Context, state string) (capella.Root, error) {
	provider, isProvider := n.client.(consensusclient.BeaconStateRootProvider)
	if !isProvider {
		return capella.Root{}, errors.New("client does not implement consensusclient.StateRootProvider")
	}

	rsp, err := provider.BeaconStateRoot(ctx, &api.BeaconStateRootOpts{
		State: state,
	})
	if err != nil {
		return capella.Root{}, err
	}

	return *rsp.Data, nil
}

func (n *node) FetchValidators(ctx context.Context, state string, indices []capella.ValidatorIndex, pubKeys []capella.BLSPubKey) (map[capella.ValidatorIndex]*v1.Validator, error) {
	provider, isProvider := n.client.(consensusclient.ValidatorsProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.ValidatorsProvider")
	}

	rsp, err := provider.Validators(ctx, &api.ValidatorsOpts{
		State:   state,
		Indices: indices,
		PubKeys: pubKeys,
	})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchBeaconCommittees(ctx context.Context, state string, epoch *capella.Epoch) ([]*v1.BeaconCommittee, error) {
	provider, isProvider := n.client.(consensusclient.BeaconCommitteesProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.BeaconCommitteesProvider")
	}

	opts := &api.BeaconCommitteesOpts{
		State: state,
	}

	if epoch != nil {
		opts.Epoch = epoch
	}

	rsp, err := provider.BeaconCommittees(ctx, opts)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchAttestationData(ctx context.Context, slot capella.Slot, committeeIndex capella.CommitteeIndex) (*capella.AttestationData, error) {
	provider, isProvider := n.client.(consensusclient.AttestationDataProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.AttestationDataProvider")
	}

	rsp, err := provider.AttestationData(ctx, &api.AttestationDataOpts{
		Slot:           slot,
		CommitteeIndex: committeeIndex,
	})
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (n *node) FetchBeaconBlockHeader(ctx context.Context, opts *api.BeaconBlockHeaderOpts) (*v1.BeaconBlockHeader, error) {
	provider, isProvider := n.client.(consensusclient.BeaconBlockHeadersProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.BeaconBlockHeadersProvider")
	}

	rsp, err := provider.BeaconBlockHeader(ctx, opts)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}
