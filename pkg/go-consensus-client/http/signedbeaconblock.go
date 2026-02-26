// Copyright © 2020 - 2024 Attestant Limited.
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
	"errors"
	"fmt"

	dynssz "github.com/pk910/dynamic-ssz"
	client "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
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
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/qrl/v1/beacon/blocks/%s", opts.Block)

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	var response *api.Response[*spec.VersionedSignedBeaconBlock]

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.signedBeaconBlockFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		response, err = s.signedBeaconBlockFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) signedBeaconBlockFromSSZ(ctx context.Context,
	res *httpResponse,
) (
	*api.Response[*spec.VersionedSignedBeaconBlock],
	error,
) {
	response := &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	var dynSSZ *dynssz.DynSsz

	if s.customSpecSupport {
		specs, err := s.Spec(ctx, &api.SpecOpts{})
		if err != nil {
			return nil, errors.Join(errors.New("failed to request specs"), err)
		}

		dynSSZ = dynssz.NewDynSsz(specs.Data)
	}

	var err error

	switch res.consensusVersion {
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.SignedBeaconBlock{}
		if s.customSpecSupport {
			err = dynSSZ.UnmarshalSSZ(response.Data.Capella, res.body)
		} else {
			err = response.Data.Capella.UnmarshalSSZ(res.body)
		}

		if err != nil {
			return nil, errors.Join(errors.New("failed to decode capella signed beacon block"), err)
		}
	default:
		return nil, fmt.Errorf("unhandled block version %s", res.consensusVersion)
	}

	return response, nil
}

func (*Service) signedBeaconBlockFromJSON(res *httpResponse) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	response := &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: res.consensusVersion,
		},
	}

	var err error

	switch res.consensusVersion {
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body),
			&capella.SignedBeaconBlock{},
		)
	default:
		return nil, fmt.Errorf("unhandled version %s", res.consensusVersion)
	}

	if err != nil {
		return nil, err
	}

	return response, nil
}
