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

// Web3 exposes metrics defined by the Web3 module.
type Web3 struct {
	client          *qrlclient.Client
	api             api.ExecutionClient
	qrlRPCClient    *qrlrpc.QRLRPC
	log             logrus.FieldLogger
	ClientVersion   prometheus.GaugeVec
	previousVersion string
}

const (
	NameWeb3 = "web3"
)

func (w *Web3) Name() string {
	return NameWeb3
}

func (w *Web3) RequiredModules() []string {
	return []string{"web3"}
}

// NewWeb3 returns a new Web3 instance.
func NewWeb3(client *qrlclient.Client, internalAPI api.ExecutionClient, qrlRPCClient *qrlrpc.QRLRPC, log logrus.FieldLogger, namespace string, constLabels map[string]string) Web3 {
	namespace += "_web3"

	constLabels["module"] = NameWeb3

	return Web3{
		client:       client,
		api:          internalAPI,
		qrlRPCClient: qrlRPCClient,
		log:          log.WithField("module", NameWeb3),
		ClientVersion: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   namespace,
				Name:        "client_version",
				Help:        "Client version.",
				ConstLabels: constLabels,
			},
			[]string{
				"version",
			},
		),
	}
}

func (w *Web3) Start(ctx context.Context) {
	w.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
			w.tick(ctx)
		}
	}
}

//nolint:unparam // context will be used in the future
func (w *Web3) tick(context.Context) {
	clientVersion, err := w.qrlRPCClient.Web3ClientVersion()
	if err != nil {
		w.log.WithError(err).Error("Failed to get node info")
	} else {
		if w.previousVersion != clientVersion {
			w.ClientVersion.Reset()

			w.ClientVersion.WithLabelValues(clientVersion).Set(1)
		}

		w.previousVersion = clientVersion
	}
}
