package types

import (
	"github.com/theQRL/go-zond/p2p/qnode"
)

// Identity represents the node identity.
type Identity struct {
	PeerID             string   `json:"peer_id"`
	QNR                string   `json:"qnr"`
	P2PAddresses       []string `json:"p2p_addresses"`
	DiscoveryAddresses []string `json:"discovery_addresses"`
	Metadata           struct {
		SeqNumber string `json:"seq_number"`
		Attnets   string `json:"attnets"`
		Syncnets  string `json:"syncnets"`
	} `json:"metadata"`
}

func (i *Identity) GetQnode() (*qnode.Node, error) {
	var node qnode.Node

	err := node.UnmarshalText([]byte(i.QNR))
	if err != nil {
		return nil, err
	}

	return &node, nil
}
