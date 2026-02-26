// Copyright © 2025 Attestant Limited.
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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/spec/capella"
)

// AttestationRewards are the rewards for a number of attesting validators.
type AttestationRewards struct {
	IdealRewards []IdealAttestationRewards     `json:"ideal_rewards"`
	TotalRewards []ValidatorAttestationRewards `json:"total_rewards"`
}

// IdealAttestationRewards are the ideal attestation rewards for an attestation.
type IdealAttestationRewards struct {
	EffectiveBalance capella.Shor
	Head             capella.Shor
	Target           capella.Shor
	Source           capella.Shor
	InclusionDelay   *capella.Shor
	Inactivity       capella.Shor
}

// idealAttestationRewardsJSON is the spec representation of the struct.
type idealAttestationRewardsJSON struct {
	EffectiveBalance string `json:"effective_balance"`
	Head             string `json:"head"`
	Target           string `json:"target"`
	Source           string `json:"source"`
	InclusionDelay   string `json:"inclusion_delay,omitempty"`
	Inactivity       string `json:"inactivity"`
}

// MarshalJSON implements json.Marshaler.
func (i *IdealAttestationRewards) MarshalJSON() ([]byte, error) {
	inclusionDelay := ""
	if i.InclusionDelay != nil {
		inclusionDelay = fmt.Sprintf("%d", *i.InclusionDelay)
	}

	return json.Marshal(&idealAttestationRewardsJSON{
		EffectiveBalance: fmt.Sprintf("%d", i.EffectiveBalance),
		Head:             fmt.Sprintf("%d", i.Head),
		Target:           fmt.Sprintf("%d", i.Target),
		Source:           fmt.Sprintf("%d", i.Source),
		InclusionDelay:   inclusionDelay,
		Inactivity:       fmt.Sprintf("%d", i.Inactivity),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *IdealAttestationRewards) UnmarshalJSON(input []byte) error {
	var err error

	var data idealAttestationRewardsJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.EffectiveBalance == "" {
		return errors.New("effective balance missing")
	}

	effectiveBalance, err := strconv.ParseUint(data.EffectiveBalance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for effective balance")
	}

	i.EffectiveBalance = capella.Shor(effectiveBalance)

	if data.Head == "" {
		return errors.New("head missing")
	}

	head, err := strconv.ParseUint(data.Head, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for head")
	}

	i.Head = capella.Shor(head)

	if data.Target == "" {
		return errors.New("target missing")
	}

	target, err := strconv.ParseUint(data.Target, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for target")
	}

	i.Target = capella.Shor(target)

	if data.Source == "" {
		return errors.New("source missing")
	}

	source, err := strconv.ParseInt(data.Source, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for source")
	}

	i.Source = capella.Shor(source)

	if data.InclusionDelay != "" {
		inclusionDelay, err := strconv.ParseUint(data.InclusionDelay, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for inclusion delay")
		}

		tmp := capella.Shor(inclusionDelay)
		i.InclusionDelay = &tmp
	}

	if data.Inactivity == "" {
		return errors.New("inactivity missing")
	}

	inactivity, err := strconv.ParseUint(data.Inactivity, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for inactivity")
	}

	i.Inactivity = capella.Shor(inactivity)

	return nil
}

// String returns a string version of the structure.
func (i *IdealAttestationRewards) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// ValidatorAttestationRewards are the ideal attestation rewards for a validator.
type ValidatorAttestationRewards struct {
	ValidatorIndex capella.ValidatorIndex
	Head           capella.Shor
	// Target can be negative, so it is an int64 (but still a Shor value).
	Target int64
	// Source can be negative, so it is an int64 (but still a Shor value).
	Source         int64
	InclusionDelay *capella.Shor
	Inactivity     capella.Shor
}

// validatorAttestationRewardsJSON is the spec representation of the struct.
type validatorAttestationRewardsJSON struct {
	ValidatorIndex string `json:"validator_index"`
	Head           string `json:"head"`
	Target         string `json:"target"`
	Source         string `json:"source"`
	InclusionDelay string `json:"inclusion_delay,omitempty"`
	Inactivity     string `json:"inactivity"`
}

// MarshalJSON implements json.Marshaler.
func (v *ValidatorAttestationRewards) MarshalJSON() ([]byte, error) {
	inclusionDelay := ""
	if v.InclusionDelay != nil {
		inclusionDelay = fmt.Sprintf("%d", *v.InclusionDelay)
	}

	return json.Marshal(&validatorAttestationRewardsJSON{
		ValidatorIndex: fmt.Sprintf("%d", v.ValidatorIndex),
		Head:           fmt.Sprintf("%d", v.Head),
		Target:         fmt.Sprintf("%d", v.Target),
		Source:         fmt.Sprintf("%d", v.Source),
		InclusionDelay: inclusionDelay,
		Inactivity:     fmt.Sprintf("%d", v.Inactivity),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *ValidatorAttestationRewards) UnmarshalJSON(input []byte) error {
	var err error

	var data validatorAttestationRewardsJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}

	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}

	v.ValidatorIndex = capella.ValidatorIndex(validatorIndex)

	if data.Head == "" {
		return errors.New("head missing")
	}

	head, err := strconv.ParseUint(data.Head, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for head")
	}

	v.Head = capella.Shor(head)

	if data.Target == "" {
		return errors.New("target missing")
	}

	v.Target, err = strconv.ParseInt(data.Target, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for target")
	}

	if data.Source == "" {
		return errors.New("source missing")
	}

	v.Source, err = strconv.ParseInt(data.Source, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for source")
	}

	if data.InclusionDelay != "" {
		inclusionDelay, err := strconv.ParseUint(data.InclusionDelay, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for inclusion delay")
		}

		tmp := capella.Shor(inclusionDelay)
		v.InclusionDelay = &tmp
	}

	if data.Inactivity == "" {
		return errors.New("inactivity missing")
	}

	inactivity, err := strconv.ParseUint(data.Inactivity, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for inactivity")
	}

	v.Inactivity = capella.Shor(inactivity)

	return nil
}

// String returns a string version of the structure.
func (v *ValidatorAttestationRewards) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
