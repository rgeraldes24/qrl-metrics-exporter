package types

import "testing"

func TestAgentParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input  string
		expect Agent
	}{
		{"Qrysm/v2.0.2/4a4a7e97dfd2285a5e48a178f693d870e9a4ff60", AgentQrysm},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			te := test
			t.Parallel()
			if actual := AgentFromString(te.input); actual != te.expect {
				t.Errorf("Expected %s, got %s", te.expect, actual)
			}
		})
	}
}
