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

package capella

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ExecutionData provides information about the state of Execution as viewed by the
// Consensus chain.
type ExecutionData struct {
	DepositRoot  Root `ssz-size:"32"`
	DepositCount uint64
	BlockHash    []byte `ssz-size:"32"`
}

// executionDataJSON is the spec representation of the struct.
type executionDataJSON struct {
	DepositRoot  string `json:"deposit_root"`
	DepositCount string `json:"deposit_count"`
	BlockHash    string `json:"block_hash"`
}

// executionDataYAML is the spec representation of the struct.
type executionDataYAML struct {
	DepositRoot  string `yaml:"deposit_root"`
	DepositCount uint64 `yaml:"deposit_count"`
	BlockHash    string `yaml:"block_hash"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionDataJSON{
		DepositRoot:  fmt.Sprintf("%#x", e.DepositRoot),
		DepositCount: strconv.FormatUint(e.DepositCount, 10),
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionData) UnmarshalJSON(input []byte) error {
	var executionDataJSON executionDataJSON
	if err := json.Unmarshal(input, &executionDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&executionDataJSON)
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionData) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&executionDataYAML{
		DepositRoot:  fmt.Sprintf("%#x", e.DepositRoot),
		DepositCount: e.DepositCount,
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionData) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var executionDataJSON executionDataJSON
	if err := yaml.Unmarshal(input, &executionDataJSON); err != nil {
		return err
	}

	return e.unpack(&executionDataJSON)
}

// String returns a string version of the structure.
func (e *ExecutionData) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (e *ExecutionData) unpack(executionDataJSON *executionDataJSON) error {
	if executionDataJSON.DepositRoot == "" {
		return errors.New("deposit root missing")
	}

	depositRoot, err := hex.DecodeString(strings.TrimPrefix(executionDataJSON.DepositRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for deposit root")
	}

	if len(depositRoot) != RootLength {
		return errors.New("incorrect length for deposit root")
	}

	copy(e.DepositRoot[:], depositRoot)

	if executionDataJSON.DepositCount == "" {
		return errors.New("deposit count missing")
	}

	if e.DepositCount, err = strconv.ParseUint(executionDataJSON.DepositCount, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for deposit count")
	}

	if executionDataJSON.BlockHash == "" {
		return errors.New("block hash missing")
	}

	if e.BlockHash, err = hex.DecodeString(strings.TrimPrefix(executionDataJSON.BlockHash, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for block hash")
	}

	if len(e.BlockHash) != HashLength {
		return errors.New("incorrect length for block hash")
	}

	return nil
}
