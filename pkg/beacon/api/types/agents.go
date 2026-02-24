package types

import (
	"strings"
)

// Agent is a peer's agent.
type Agent string

const (
	// AgentUnknown is an unknown agent.
	AgentUnknown Agent = "unknown"
	// AgentQrysm is a Qrysm agent.
	AgentQrysm Agent = "qrysm"
)

// AllAgents is a list of all agents.
var AllAgents = []Agent{
	AgentUnknown,
	AgentQrysm,
}

// AgentCount represents the number of peers with each agent.
type AgentCount struct {
	Unknown    int `json:"unknown"`
	Qrysm      int `json:"qrysm"`
}

// AgentFromString returns the agent from the given string.
func AgentFromString(agent string) Agent {
	asLower := strings.ToLower(agent)

	if strings.Contains(asLower, "qrysm") {
		return AgentQrysm
	}

	return AgentUnknown
}
