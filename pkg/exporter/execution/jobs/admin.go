package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/theQRL/go-zond/qrlclient"
	"github.com/theQRL/qrl-metrics-exporter/pkg/exporter/execution/api"
	"github.com/theQRL/qrl-metrics-exporter/pkg/exporter/execution/api/types"
	"github.com/theQRL/qrl-metrics-exporter/pkg/qrlrpc"
)

// Admin exposes metrics defined by the admin module.
type Admin struct {
	client       *qrlclient.Client
	api          api.ExecutionClient
	qrlRPCClient *qrlrpc.QRLRPC
	log          logrus.FieldLogger
	NodeInfo     prometheus.GaugeVec
	Port         prometheus.GaugeVec
	Peers        prometheus.Gauge
}

const (
	NameAdmin = "admin"
)

func (a *Admin) Name() string {
	return NameAdmin
}

func (a *Admin) RequiredModules() []string {
	return []string{"admin"}
}

// NewAdmin returns a new Admin instance.
func NewAdmin(client *qrlclient.Client, internalAPI api.ExecutionClient, qrlRPCClient *qrlrpc.QRLRPC, log logrus.FieldLogger, namespace string, constLabels map[string]string) Admin {
	namespace += "_admin"

	constLabels["module"] = NameAdmin

	return Admin{
		client:       client,
		api:          internalAPI,
		qrlRPCClient: qrlRPCClient,
		log:          log.WithField("module", NameAdmin),
		NodeInfo: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "node_info",
				Help:        "Node info.",
				ConstLabels: constLabels,
			},
			[]string{
				"ip",
				"listenAddr",
				"name",
				"discovery_port",
				"listener_port",
				"network",
			},
		),
		Port: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "node_port",
				Help:        "The ports for the node.",
				ConstLabels: constLabels,
			},
			[]string{
				"name",
				"port_name",
			},
		),
		Peers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "peers",
				Help:        "The number of peers connected with the node.",
				ConstLabels: constLabels,
			},
		),
	}
}

func (a *Admin) Start(ctx context.Context) {
	a.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			a.tick(ctx)
		}
	}
}

func (a *Admin) tick(ctx context.Context) {
	nodeInfo, err := a.api.AdminNodeInfo(ctx)
	if err != nil {
		a.log.WithError(err).Error("Failed to get node info")
	} else {
		a.ObserveNodeInfo(nodeInfo)
	}

	peers, err := a.api.AdminPeers(ctx)
	if err != nil {
		a.log.WithError(err).Error("Failed to get peers")
	} else {
		a.ObservePeers(len(peers))
	}
}

func (a *Admin) ObserveNodeInfo(nodeInfo *types.NodeInfo) {
	// Info
	a.NodeInfo.WithLabelValues(nodeInfo.IP,
		nodeInfo.ListenAddr,
		nodeInfo.Name,
		fmt.Sprint(nodeInfo.Ports.Discovery),
		fmt.Sprint(nodeInfo.Ports.Listener),
		fmt.Sprint(nodeInfo.Protocols.QRL.NetworkID),
	).Set(1)

	// Ports
	a.Port.WithLabelValues("discovery", "discovery").Set(float64(nodeInfo.Ports.Discovery))
	a.Port.WithLabelValues("listener", "listener").Set(float64(nodeInfo.Ports.Listener))
}

func (a *Admin) ObservePeers(peers int) {
	a.Peers.Set(float64(peers))
}
