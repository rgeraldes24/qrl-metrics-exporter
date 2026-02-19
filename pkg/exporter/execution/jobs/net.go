package jobs

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/theQRL/go-zond/qrlclient"
	"github.com/theQRL/qrl-metrics-exporter/pkg/exporter/execution/api"
	"github.com/theQRL/qrl-metrics-exporter/pkg/qrlrpc"
)

// Net exposes metrics defined by the net module.
type Net struct {
	client       *qrlclient.Client
	api          api.ExecutionClient
	qrlRPCClient *qrlrpc.QRLRPC
	log          logrus.FieldLogger
	PeerCount    prometheus.Gauge
}

const (
	NameNet = "net"
)

func (n *Net) Name() string {
	return NameNet
}

func (n *Net) RequiredModules() []string {
	return []string{"net"}
}

// NewNet returns a new Net instance.
func NewNet(client *qrlclient.Client, internalAPI api.ExecutionClient, qrlRPCClient *qrlrpc.QRLRPC, log logrus.FieldLogger, namespace string, constLabels map[string]string) Net {
	namespace += "_net"

	constLabels["module"] = NameWeb3

	return Net{
		client:       client,
		api:          internalAPI,
		qrlRPCClient: qrlRPCClient,
		log:          log.WithField("module", NameNet),
		PeerCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "peer_count",
				Help:        "The amount of peers connected to the node.",
				ConstLabels: constLabels,
			},
		),
	}
}

func (n *Net) Start(ctx context.Context) {
	n.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			n.tick(ctx)
		}
	}
}

//nolint:unparam // context will be used in the future
func (n *Net) tick(context.Context) {
	count, err := n.qrlRPCClient.NetPeerCount()
	if err != nil {
		n.log.WithError(err).Error("Failed to get peer count")
	} else {
		n.PeerCount.Set(float64(count))
	}
}
