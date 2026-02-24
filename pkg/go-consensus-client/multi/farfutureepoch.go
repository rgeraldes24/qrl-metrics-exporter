// Copyright © 2021, 2022 Attestant Limited.
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

package multi

import (
	"context"

	"github.com/pkg/errors"
	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// FarFutureEpoch provides the far future epoch of the chain.
func (s *Service) FarFutureEpoch(ctx context.Context) (capella.Epoch, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		epoch, err := client.(consensusclient.FarFutureEpochProvider).FarFutureEpoch(ctx)
		if err != nil {
			return nil, err
		}

		if epoch == 0 {
			return nil, errors.New("zero epoch not a valid response")
		}

		return epoch, nil
	}, nil)
	if err != nil {
		return 0, err
	}

	response, isResponse := res.(capella.Epoch)
	if !isResponse {
		return 0, ErrIncorrectType
	}

	return response, nil
}
