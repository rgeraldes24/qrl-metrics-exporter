package types

import (
	"encoding/json"

	"github.com/theQRL/go-zond/common"
)

// NodeInfo is the information about the node.
type NodeInfo struct {
	Qnode      string `json:"qnode"`
	ID         string `json:"id"`
	IP         string `json:"ip"`
	ListenAddr string `json:"listenAddr"`
	Name       string `json:"name"`
	Ports      struct {
		Discovery int `json:"discovery"`
		Listener  int `json:"listener"`
	} `json:"ports"`
	Protocols struct {
		QRL QRLProtocol `json:"qrl"`
	} `json:"protocols"`
}

// QRLProtocol is the information about the qrl protocol.
type QRLProtocol struct {
	Genesis   common.Hash `json:"genesis"`
	Head      common.Hash `json:"head"`
	NetworkID int         `json:"networkID"`
}

// UnmarshalJSON implements the json.Unmarshaler interface, overriding to handle
// clients returning string prefixed with 0x.
func (e *QRLProtocol) UnmarshalJSON(data []byte) error {
	var v struct {
		Genesis   common.Hash `json:"genesis"`
		Head      common.Hash `json:"head"`
		NetworkID int         `json:"networkID"`
	}

	var objMap map[string]*json.RawMessage

	err := json.Unmarshal(data, &objMap)
	if err != nil {
		return err
	}

	e.Genesis = v.Genesis
	e.Head = v.Head
	e.NetworkID = v.NetworkID

	return nil
}
