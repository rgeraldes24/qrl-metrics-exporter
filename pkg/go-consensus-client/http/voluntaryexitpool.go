// Copyright © 2021, 2024 Attestant Limited.
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

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	client "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

type voluntaryExitPoolJSON struct {
	Data []*capella.SignedVoluntaryExit `json:"data"`
}

// VoluntaryExitPool obtains the voluntary exit pool.
func (s *Service) VoluntaryExitPool(ctx context.Context,
	opts *api.VoluntaryExitPoolOpts,
) (
	*api.Response[[]*capella.SignedVoluntaryExit],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := "/qrl/v1/beacon/pool/voluntary_exits"

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, false)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request voluntary exit pool"), err)
	}

	var voluntaryExitPoolJSON voluntaryExitPoolJSON
	if err := json.NewDecoder(bytes.NewReader(httpResponse.body)).Decode(&voluntaryExitPoolJSON); err != nil {
		return nil, errors.Join(errors.New("failed to parse voluntary exit pool"), err)
	}

	// Ensure the data returned to us is as expected given our input.
	if voluntaryExitPoolJSON.Data == nil {
		return nil, errors.New("voluntary exit pool not returned")
	}

	return &api.Response[[]*capella.SignedVoluntaryExit]{
		Data:     voluntaryExitPoolJSON.Data,
		Metadata: make(map[string]any),
	}, nil
}
