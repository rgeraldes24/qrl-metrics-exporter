// Copyright © 2023 Attestant Limited.
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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// PayloadAttributesEvent represents the data of a payload_attributes event.
type PayloadAttributesEvent struct {
	// Version is the fork version of the beacon chain.
	Version spec.DataVersion
	// Data is the data of the event.
	Data *PayloadAttributesData
}

// PayloadAttributesData represents the data of a payload_attributes event.
type PayloadAttributesData struct {
	// ProposerIndex is the index of the proposer.
	ProposerIndex capella.ValidatorIndex
	// ProposalSlot is the slot of the proposal.
	ProposalSlot capella.Slot
	// ParentBlockNumber is the number of the parent block.
	ParentBlockNumber uint64
	// ParentBlockRoot is the root of the parent block.
	ParentBlockRoot capella.Root
	// ParentBlockHash is the hash of the parent block.
	ParentBlockHash capella.Hash32
	// V2 is the v2 payload attributes.
	V2 *PayloadAttributesV2
}

// PayloadAttributesV2 represents the payload attributes v2.
type PayloadAttributesV2 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient capella.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
}

// payloadAttributesEventJSON is the spec representation of the event.
type payloadAttributesEventJSON struct {
	Version spec.DataVersion           `json:"version"`
	Data    *payloadAttributesDataJSON `json:"data"`
}

// payloadAttributesDataJSON is the spec representation of the payload attributes data.
type payloadAttributesDataJSON struct {
	ProposerIndex     string          `json:"proposer_index"`
	ProposalSlot      string          `json:"proposal_slot"`
	ParentBlockNumber string          `json:"parent_block_number"`
	ParentBlockRoot   string          `json:"parent_block_root"`
	ParentBlockHash   string          `json:"parent_block_hash"`
	PayloadAttributes json.RawMessage `json:"payload_attributes"`
}

// payloadAttributesV2JSON is the spec representation of the payload attributes v2.
type payloadAttributesV2JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV2) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV2JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV2) unpack(data *payloadAttributesV2JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}

	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}

	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}

	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}

	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}

	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}

	if len(feeRecipient) != capella.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}

	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}

	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}

	p.Withdrawals = data.Withdrawals

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e *PayloadAttributesEvent) MarshalJSON() ([]byte, error) {
	var (
		payloadAttributes []byte
		err               error
	)

	switch e.Version {
	case spec.DataVersionCapella:
		if e.Data.V2 == nil {
			return nil, errors.New("no payload attributes v2 data")
		}

		payloadAttributes, err = json.Marshal(&payloadAttributesV2JSON{
			Timestamp:             strconv.FormatUint(e.Data.V2.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V2.PrevRandao),
			SuggestedFeeRecipient: e.Data.V2.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V2.Withdrawals,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v2")
		}
	default:
		return nil, fmt.Errorf("unsupported payload attributes version: %s", e.Version)
	}

	data := payloadAttributesDataJSON{
		ProposerIndex:     fmt.Sprintf("%d", e.Data.ProposerIndex),
		ProposalSlot:      fmt.Sprintf("%d", e.Data.ProposalSlot),
		ParentBlockNumber: strconv.FormatUint(e.Data.ParentBlockNumber, 10),
		ParentBlockRoot:   fmt.Sprintf("%#x", e.Data.ParentBlockRoot),
		ParentBlockHash:   fmt.Sprintf("%#x", e.Data.ParentBlockHash),
		PayloadAttributes: payloadAttributes,
	}

	return json.Marshal(&payloadAttributesEventJSON{
		Version: e.Version,
		Data:    &data,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *PayloadAttributesEvent) UnmarshalJSON(input []byte) error {
	var event payloadAttributesEventJSON
	if err := json.Unmarshal(input, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&event)
}

// String returns a string version of the structure.
func (e *PayloadAttributesEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (e *PayloadAttributesEvent) unpack(data *payloadAttributesEventJSON) error {
	var err error

	if data.Data == nil {
		return errors.New("payload attributes data missing")
	}

	e.Data = &PayloadAttributesData{}

	if data.Data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}

	proposerIndex, err := strconv.ParseUint(data.Data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}

	e.Data.ProposerIndex = capella.ValidatorIndex(proposerIndex)

	if data.Data.ProposalSlot == "" {
		return errors.New("proposal slot missing")
	}

	proposalSlot, err := strconv.ParseUint(data.Data.ProposalSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposal slot")
	}

	e.Data.ProposalSlot = capella.Slot(proposalSlot)

	if data.Data.ParentBlockNumber == "" {
		return errors.New("parent block number missing")
	}

	parentBlockNumber, err := strconv.ParseUint(data.Data.ParentBlockNumber, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block number")
	}

	e.Data.ParentBlockNumber = parentBlockNumber

	if data.Data.ParentBlockRoot == "" {
		return errors.New("parent block root missing")
	}

	parentBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block root")
	}

	if len(parentBlockRoot) != capella.RootLength {
		return errors.New("incorrect length for parent block root")
	}

	copy(e.Data.ParentBlockRoot[:], parentBlockRoot)

	if data.Data.ParentBlockHash == "" {
		return errors.New("parent block hash missing")
	}

	parentBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block hash")
	}

	if len(parentBlockHash) != capella.Hash32Length {
		return errors.New("incorrect length for parent block hash")
	}

	copy(e.Data.ParentBlockHash[:], parentBlockHash)

	if data.Data.PayloadAttributes == nil {
		return errors.New("payload attributes missing")
	}

	switch data.Version {
	case spec.DataVersionCapella:
		var payloadAttributes PayloadAttributesV2

		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}

		e.Data.V2 = &payloadAttributes
	default:
		return errors.New("unsupported data version")
	}

	e.Version = data.Version

	return nil
}
