package beacon

import (
	"context"
	"errors"
	"fmt"
	"time"

	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	v1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

func (n *node) ensureBeaconSubscription(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second * 2):
			if len(n.options.BeaconSubscription.Topics) == 0 {
				continue
			}

			if n.client == nil {
				continue
			}

			if !n.options.BeaconSubscription.Enabled {
				continue
			}

			if err := n.subscribeToBeaconEvents(ctx); err != nil {
				n.log.WithError(err).Error("Failed to subscribe to beacon")

				continue
			}

			return nil
		}
	}
}

func (n *node) subscribeToBeaconEvents(ctx context.Context) error {
	provider, isProvider := n.client.(consensusclient.EventsProvider)
	if !isProvider {
		return errors.New("client does not implement consensusclient.Subscriptions")
	}

	topics := n.options.BeaconSubscription.Topics
	n.log.WithField("topics", topics).Info("Subscribing to events upstream")

	// Open a new subscription for each topic.
	for _, topic := range topics {
		n.log.WithField("topic", topic).Info("Subscribing to event")

		if err := provider.Events(ctx, &api.EventsOpts{
			Topics: []string{topic},
			Handler: func(event *v1.Event) {
				n.lastEventTimeMu.Lock()
				n.lastEventTime = time.Now()
				n.lastEventTimeMu.Unlock()

				if err := n.handleEvent(ctx, event); err != nil {
					n.log.Errorf("Failed to handle event: %v", err)
				}
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (n *node) handleEvent(ctx context.Context, event *v1.Event) error {
	n.publishEvent(ctx, event)

	switch event.Topic {
	case topicAttestation:
		return n.handleAttestation(ctx, event)
	case topicBlock:
		return n.handleBlock(ctx, event)
	case topicBlockGossip:
		return n.handleBlockGossip(ctx, event)
	case topicChainReorg:
		return n.handleChainReorg(ctx, event)
	case topicFinalizedCheckpoint:
		return n.handleFinalizedCheckpoint(ctx, event)
	case topicHead:
		return n.handleHead(ctx, event)
	case topicVoluntaryExit:
		return n.handleVoluntaryExit(ctx, event)
	case topicContributionAndProof:
		return n.handleContributionAndProof(ctx, event)

	default:
		return fmt.Errorf("unknown event topic %s", event.Topic)
	}
}

func (n *node) handleAttestation(ctx context.Context, event *v1.Event) error {
	attestation, valid := event.Data.(*spec.VersionedAttestation)
	if !valid {
		return errors.New("invalid attestation event")
	}

	n.publishAttestation(ctx, attestation)

	return nil
}

func (n *node) handleBlock(ctx context.Context, event *v1.Event) error {
	block, valid := event.Data.(*v1.BlockEvent)
	if !valid {
		return errors.New("invalid block event")
	}

	n.publishBlock(ctx, block)

	return nil
}

func (n *node) handleBlockGossip(ctx context.Context, event *v1.Event) error {
	blockGossip, valid := event.Data.(*v1.BlockGossipEvent)
	if !valid {
		return errors.New("invalid block gossip event")
	}

	n.publishBlockGossip(ctx, blockGossip)

	return nil
}

func (n *node) handleChainReorg(ctx context.Context, event *v1.Event) error {
	chainReorg, valid := event.Data.(*v1.ChainReorgEvent)
	if !valid {
		return errors.New("invalid chain reorg event")
	}

	n.publishChainReOrg(ctx, chainReorg)

	return nil
}

func (n *node) handleFinalizedCheckpoint(ctx context.Context, event *v1.Event) error {
	checkpoint, valid := event.Data.(*v1.FinalizedCheckpointEvent)
	if !valid {
		return errors.New("invalid checkpoint event")
	}

	n.publishFinalizedCheckpoint(ctx, checkpoint)

	return nil
}

func (n *node) handleHead(ctx context.Context, event *v1.Event) error {
	head, valid := event.Data.(*v1.HeadEvent)
	if !valid {
		return errors.New("invalid head event")
	}

	n.publishHead(ctx, head)

	return nil
}

func (n *node) handleVoluntaryExit(ctx context.Context, event *v1.Event) error {
	exit, valid := event.Data.(*capella.SignedVoluntaryExit)
	if !valid {
		return errors.New("invalid voluntary exit event")
	}

	n.publishVoluntaryExit(ctx, exit)

	return nil
}

func (n *node) handleContributionAndProof(ctx context.Context, event *v1.Event) error {
	contributionAndProof, valid := event.Data.(*capella.SignedContributionAndProof)
	if !valid {
		return errors.New("invalid contribution and proof event")
	}

	n.publishContributionAndProof(ctx, contributionAndProof)

	return nil
}
