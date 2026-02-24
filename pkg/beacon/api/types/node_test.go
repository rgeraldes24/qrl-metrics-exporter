package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theQRL/qrl-metrics-exporter/pkg/beacon/api/types"
)

func TestIdentity_GetEnode(t *testing.T) {
	identity := &types.Identity{
		QNR: "qnr:-IS4QHCYrYZbAKWCBRlAy5zzaDZXJBGkcnh4MHcBFZntXNFrdvJjX04jRzjzCBOonrkTfj499SZuOh8R33Ls8RRcy5wBgmlkgnY0gmlwhH8AAAGJc2VjcDI1NmsxoQPKY0yuDUmstAHYpMa2_oxVtw0RW_QAdpzBQA8yWM0xOIN1ZHCCdl8",
	}

	qnode, err := identity.GetQnode()
	require.NoError(t, err)
	require.NotNil(t, qnode)

	// Verify enode details
	require.Equal(t, "127.0.0.1", qnode.IP().String())
	require.Equal(t, 30303, qnode.UDP())
	require.Equal(t, 0, qnode.TCP())
}
