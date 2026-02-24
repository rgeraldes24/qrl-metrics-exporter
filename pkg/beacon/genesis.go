package beacon

import (
	"context"
	"errors"

	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	v1 "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1"
)

func (n *node) FetchGenesis(ctx context.Context) (*v1.Genesis, error) {
	provider, isProvider := n.client.(consensusclient.GenesisProvider)
	if !isProvider {
		return nil, errors.New("client does not implement consensusclient.GenesisProvider")
	}

	rsp, err := provider.Genesis(ctx, &api.GenesisOpts{})
	if err != nil {
		return nil, err
	}

	n.genesisMu.Lock()
	n.genesis = rsp.Data
	n.genesisMu.Unlock()

	return rsp.Data, nil
}
