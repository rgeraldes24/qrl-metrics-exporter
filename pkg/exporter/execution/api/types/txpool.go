package types

import (
	"github.com/theQRL/go-zond/common/hexutil"
)

// TXPoolStatus is the information about the transaction pool.
type TXPoolStatus struct {
	Pending hexutil.Uint64 `json:"pending"`
	Queued  hexutil.Uint64 `json:"queued"`
}
