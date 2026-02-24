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
	"math/big"
	"strings"

	dynssz "github.com/pk910/dynamic-ssz"
	client "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
	apiv1capella "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api/v1/capella"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
	"go.opentelemetry.io/otel"
)

// Proposal fetches a potential beacon block for signing.
func (s *Service) Proposal(ctx context.Context,
	opts *api.ProposalOpts,
) (
	*api.Response[*api.VersionedProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-consensus-client.http").Start(ctx, "Proposal")
	defer span.End()

	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.Slot == 0 {
		return nil, errors.Join(errors.New("no slot specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/qrl/v3/validator/blocks/%d", opts.Slot)
	query := fmt.Sprintf("randao_reveal=%#x&graffiti=%#x", opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		query = fmt.Sprintf("%s&skip_randao_verification", query)
	}

	if opts.BuilderBoostFactor == nil {
		query += "&builder_boost_factor=100"
	} else {
		query = fmt.Sprintf("%s&builder_boost_factor=%d", query, *opts.BuilderBoostFactor)
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, true)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request beacon block proposal"), err)
	}

	var response *api.Response[*api.VersionedProposal]

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.beaconBlockProposalFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		response, err = s.beaconBlockProposalFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	if err != nil {
		return nil, err
	}

	// Ensure the data returned to us is as expected given our input.
	blockSlot, err := response.Data.Slot()
	if err != nil {
		return nil, err
	}

	if blockSlot != opts.Slot {
		return nil, errors.Join(
			fmt.Errorf("beacon block proposal for slot %d; expected %d", blockSlot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}

	// Only check the RANDAO reveal if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, errors.Join(
				fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], opts.RandaoReveal[:]),
				client.ErrInconsistentResult,
			)
		}
	}

	return response, nil
}

//nolint:nestif
func (s *Service) beaconBlockProposalFromSSZ(ctx context.Context,
	res *httpResponse,
) (
	*api.Response[*api.VersionedProposal],
	error,
) {
	response := &api.Response[*api.VersionedProposal]{
		Data: &api.VersionedProposal{
			Version:        res.consensusVersion,
			ConsensusValue: big.NewInt(0),
			ExecutionValue: big.NewInt(0),
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	if err := s.populateProposalDataFromHeaders(response, res.headers); err != nil {
		return nil, err
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
		if response.Data.Blinded {
			response.Data.CapellaBlinded = &apiv1capella.BlindedBeaconBlock{}
			if s.customSpecSupport {
				err = dynSSZ.UnmarshalSSZ(response.Data.CapellaBlinded, res.body)
			} else {
				err = response.Data.CapellaBlinded.UnmarshalSSZ(res.body)
			}
		} else {
			response.Data.Capella = &capella.BeaconBlock{}
			if s.customSpecSupport {
				err = dynSSZ.UnmarshalSSZ(response.Data.Capella, res.body)
			} else {
				err = response.Data.Capella.UnmarshalSSZ(res.body)
			}
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to decode %v SSZ beacon block (blinded: %v)", res.consensusVersion, response.Data.Blinded),
			err,
		)
	}

	return response, nil
}

func (s *Service) beaconBlockProposalFromJSON(res *httpResponse) (*api.Response[*api.VersionedProposal], error) {
	response := &api.Response[*api.VersionedProposal]{
		Data: &api.VersionedProposal{
			Version:        res.consensusVersion,
			ConsensusValue: big.NewInt(0),
			ExecutionValue: big.NewInt(0),
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	if err := s.populateProposalDataFromHeaders(response, res.headers); err != nil {
		return nil, err
	}

	var err error

	switch res.consensusVersion {
	case spec.DataVersionCapella:
		if response.Data.Blinded {
			response.Data.CapellaBlinded, response.Metadata, err = decodeJSONResponse(
				bytes.NewReader(res.body),
				&apiv1capella.BlindedBeaconBlock{},
			)
		} else {
			response.Data.Capella, response.Metadata, err = decodeJSONResponse(
				bytes.NewReader(res.body),
				&capella.BeaconBlock{},
			)
		}
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}

	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to decode %v JSON beacon block (blinded: %v)", res.consensusVersion, response.Data.Blinded),
			err,
		)
	}

	return response, nil
}

func (*Service) populateProposalDataFromHeaders(response *api.Response[*api.VersionedProposal],
	headers map[string]string,
) error {
	for k, v := range headers {
		switch {
		case strings.EqualFold(k, "Qrl-Execution-Payload-Blinded"):
			response.Data.Blinded = strings.EqualFold(v, "true")
		case strings.EqualFold(k, "Qrl-Execution-Payload-Value"):
			var success bool

			response.Data.ExecutionValue, success = new(big.Int).SetString(v, 10)
			if !success {
				return fmt.Errorf("proposal header Qrl-Execution-Payload-Value %s not a valid integer", v)
			}
		case strings.EqualFold(k, "Qrl-Consensus-Block-Value"):
			var success bool

			response.Data.ConsensusValue, success = new(big.Int).SetString(v, 10)
			if !success {
				return fmt.Errorf("proposal header Qrl-Consensus-Block-Value %s not a valid integer", v)
			}
		default:
			// Unknown header, ignore
		}
	}

	return nil
}
