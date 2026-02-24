// Copyright © 2020 - 2023 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	client "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

func TestBeaconCommittees(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	epoch := capella.Epoch(0)

	tests := []struct {
		name     string
		opts     *api.BeaconCommitteesOpts
		expected *capella.Attestation
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NilState",
			opts: &api.BeaconCommitteesOpts{},
			err:  "no state specified",
		},
		{
			name: "Good",
			opts: &api.BeaconCommitteesOpts{
				State: "head",
			},
		},
		{
			name: "Genesis",
			opts: &api.BeaconCommitteesOpts{
				State: "genesis",
				Epoch: &epoch,
			},
		},
	}

	service := testService(ctx, t).(client.Service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BeaconCommitteesProvider).BeaconCommittees(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.expected != nil {
					require.Equal(t, test.expected, response.Data)
				}
			}
		})
	}
}
