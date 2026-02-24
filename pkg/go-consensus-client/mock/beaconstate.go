// Copyright © 2020 Attestant Limited.
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

package mock

import (
	"context"

	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// BeaconState fetches a beacon state given a state ID.
func (s *Service) BeaconState(ctx context.Context,
	opts *api.BeaconStateOpts,
) (
	*api.Response[*spec.VersionedBeaconState],
	error,
) {
	if s.BeaconStateFunc != nil {
		return s.BeaconStateFunc(ctx, opts)
	}

	data := &spec.VersionedBeaconState{
		Version: spec.DataVersionCapella,
		Capella: &capella.BeaconState{
			LatestBlockHeader:           &capella.BeaconBlockHeader{},
			ExecutionData:               &capella.ExecutionData{},
			PreviousJustifiedCheckpoint: &capella.Checkpoint{},
			CurrentJustifiedCheckpoint:  &capella.Checkpoint{},
			FinalizedCheckpoint:         &capella.Checkpoint{},
		},
	}

	return &api.Response[*spec.VersionedBeaconState]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
