// Copyright © 2020, 2023 Attestant Limited.
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

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Service) SignedBeaconBlock(ctx context.Context,
	opts *api.SignedBeaconBlockOpts,
) (
	*api.Response[*spec.VersionedSignedBeaconBlock],
	error,
) {
	if s.SignedBeaconBlockFunc != nil {
		return s.SignedBeaconBlockFunc(ctx, opts)
	}

	return &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: spec.DataVersionCapella,
			Capella: &capella.SignedBeaconBlock{
				Message: &capella.BeaconBlock{
					Body: &capella.BeaconBlockBody{
						ExecutionData: &capella.ExecutionData{},
					},
				},
			},
		},
		Metadata: make(map[string]any),
	}, nil
}
