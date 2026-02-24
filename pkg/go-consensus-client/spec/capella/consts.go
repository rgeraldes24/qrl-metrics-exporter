// Copyright © 2022 Attestant Limited.
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

import "math/big"

// ForkVersionLength is the number of bytes in a fork version.
const ForkVersionLength = 4

// HashLength is the number of bytes in a hash.
const HashLength = 32

// RootLength is the number of bytes in a root.
const RootLength = 32

// PublicKeyLength is the number of bytes in a public key.
const PublicKeyLength = 48

// GraffitiLength is the number of bytes in a block graffiti.
const GraffitiLength = 32

// DomainTypeLength is the number of bytes in a domain type.
const DomainTypeLength = 4

// DomainLength is the number of bytes in a domain.
const DomainLength = 32

// Hash32Length is the number of bytes in a 32-byte hash.
const Hash32Length = 32

var maxBaseFeePerGas = new(big.Int).SetBytes([]byte{
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
})

// FeeRecipientLength is the number of bytes in an execution fee recipient.
const FeeRecipientLength = 20

// ExecutionAddressLength is the number of bytes in an execution address.
const ExecutionAddressLength = 20

// MaxBytesPerTransaction is the maximum number of bytes in a transaction.
const MaxBytesPerTransaction = 1_073_741_824

// MaxTransactionsPerPayload is the maximum number of transactions in a payload.
const MaxTransactionsPerPayload = 1_048_576
