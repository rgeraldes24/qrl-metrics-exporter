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
	"bytes"
	"context"

	"github.com/pkg/errors"
	consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// emptyDomain is used for comparison purposes.
var emptyDomain capella.Domain

// Domain provides a domain for a given domain type at a given epoch.
func (s *Service) Domain(ctx context.Context,
	domainType capella.DomainType,
	epoch capella.Epoch,
) (
	capella.Domain,
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		domain, err := client.(consensusclient.DomainProvider).Domain(ctx, domainType, epoch)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(domain[:], emptyDomain[:]) {
			return nil, errors.New("empty domain not a valid response")
		}

		return domain, nil
	}, nil)
	if err != nil {
		return capella.Domain{}, err
	}

	response, isResponse := res.(capella.Domain)
	if !isResponse {
		return capella.Domain{}, ErrIncorrectType
	}

	return response, nil
}

// GenesisDomain provides a domain for a given domain type.
func (s *Service) GenesisDomain(ctx context.Context,
	domainType capella.DomainType,
) (
	capella.Domain,
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		domain, err := client.(consensusclient.DomainProvider).GenesisDomain(ctx, domainType)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(domain[:], emptyDomain[:]) {
			return nil, errors.New("empty domain not a valid response")
		}

		return domain, nil
	}, nil)
	if err != nil {
		return capella.Domain{}, err
	}

	response, isResponse := res.(capella.Domain)
	if !isResponse {
		return capella.Domain{}, ErrIncorrectType
	}

	return response, nil
}
