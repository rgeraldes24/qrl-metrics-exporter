// Copyright © 2022, 2023 Attestant Limited.
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

//nolint:revive
// Need to `go install github.com/ferranbt/fastssz/sszgen@latest` for this to work.
//go:generate rm -f aggregateandproof_ssz.go attestation_ssz.go attestationdata_ssz.go attesterslashing_ssz.go beaconblock_ssz.go beaconblockbody_ssz.go beaconblockheader_ssz.go beaconstate_ssz.go checkpoint_ssz.go contributionandproof_ssz.go deposit_ssz.go depositdata_ssz.go depositmessage_ssz.go executiondata_ssz.go executionpayload_ssz.go executionpayloadheader_ssz.go fork_ssz.go forkdata_ssz.go historicalsummary_ssz.go indexedattestation_ssz.go pendingattestation_ssz.go proposerslashing_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go signedbeaconblockheader_ssz.go signedcontributionandproof_ssz.go signedvoluntaryexit_ssz.go signingdata_ssz.go syncaggregate_ssz.go syncaggregatorselectiondata_ssz.go synccommittee_ssz.go synccommitteecontribution_ssz.go synccommitteemessage_ssz.go validator_ssz.go voluntaryexit_ssz.go withdrawal_ssz.go
//go:generate sszgen -suffix ssz -path . -objs AggregateAndProof,Attestation,AttestationData,AttesterSlashing,BeaconBlock,BeaconBlockBody,BeaconBlockHeader,BeaconState,Checkpoint,ContributionAndProof,Deposit,DepositData,DepositMessage,ExecutionData,ExecutionPayload,ExecutionPayloadHeader,Fork,ForkData,HistoricalSummary,IndexedAttestation,PendingAttestation,ProposerSlashing,SignedAggregateAndProof,SignedBeaconBlock,SignedBeaconBlockHeader,SignedContributionAndProof,SignedVoluntaryExit,SigningData,SyncAggregate,SyncAggregatorSelectionData,SyncCommittee,SyncCommitteeContribution,SyncCommitteeMessage,Validator,VoluntaryExit,Withdrawal
//go:generate goimports -w aggregateandproof_ssz.go attestation_ssz.go attestationdata_ssz.go attesterslashing_ssz.go beaconblock_ssz.go beaconblockbody_ssz.go beaconblockheader_ssz.go beaconstate_ssz.go checkpoint_ssz.go contributionandproof_ssz.go deposit_ssz.go depositdata_ssz.go depositmessage_ssz.go executiondata_ssz.go executionpayload_ssz.go executionpayloadheader_ssz.go fork_ssz.go forkdata_ssz.go historicalsummary_ssz.go indexedattestation_ssz.go pendingattestation_ssz.go proposerslashing_ssz.go signedaggregateandproof_ssz.go signedbeaconblock_ssz.go signedbeaconblockheader_ssz.go signedcontributionandproof_ssz.go signedvoluntaryexit_ssz.go signingdata_ssz.go syncaggregate_ssz.go syncaggregatorselectiondata_ssz.go synccommittee_ssz.go synccommitteecontribution_ssz.go synccommitteemessage_ssz.go validator_ssz.go voluntaryexit_ssz.go withdrawal_ssz.go
