# go-consensus-client

Go library providing an abstraction to multiple QRL beacon nodes.  Its external API follows the official [QRL beacon APIs](https://github.com/ethereum/beacon-APIs) specification.

This library is under development; expect APIs and data structures to change until it reaches version 1.0.  In addition, clients' implementations of both their own and the standard API are themselves under development so implementation of the full API can be incomplete.

> Between versions 0.18.0 and 0.19.0 the API has undergone a number of changes.  Please see [the detailed documentation](docs/0.19.0-changes.md) regarding these changes.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-consensus-client` is a standard Go module which can be installed with:

```sh
go get github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client
```

## Support

`go-consensus-client` supports beacon nodes that comply with the standard beacon node API.  To date it has been tested against the following beacon nodes:

  - [Qrysm](https://github.com/theQRL/qrysm) minimum version ?

## Usage

Please read the [Go documentation for this library](https://godoc.org/github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client) for interface information.

## Example

Below is a complete annotated example to access a beacon node.

```go
package main

import (
    "context"
    "fmt"
    
    consensusclient "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client"
    "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/api"
    "github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/http"
    "github.com/rs/zerolog"
)

func main() {
    // Provide a cancellable context to the creation function.
    ctx, cancel := context.WithCancel(context.Background())
    client, err := http.New(ctx,
        // WithAddress supplies the address of the beacon node, as a URL.
        http.WithAddress("http://localhost:5052/"),
        // LogLevel supplies the level of logging to carry out.
        http.WithLogLevel(zerolog.WarnLevel),
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Connected to %s\n", client.Name())
    
    // Client functions have their own interfaces.  Not all functions are
    // supported by all clients, so checks should be made for each function when
    // casting the service to the relevant interface.
    if provider, isProvider := client.(consensusclient.GenesisProvider); isProvider {
        genesisResponse, err := provider.Genesis(ctx, &api.GenesisOpts{})
        if err != nil {
            // Errors may be API errors, in which case they will have more detail
            // about the failure.
            var apiErr *api.Error
            if errors.As(err, &apiErr) {
                switch apiErr.StatusCode {
                  case 404:
                    panic("genesis not found")
                  case 503:
                    panic("node is syncing")
                }
            }
            panic(err)
        }
        fmt.Printf("Genesis time is %v\n", genesisResponse.Data.GenesisTime)
    }

    // You can also access the struct directly if required.
    httpClient := client.(*http.Service)
    genesisResponse, err := httpClient.Genesis(ctx, &api.GenesisOpts{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Genesis validators root is %s\n", genesisResponse.Data.GenesisValidatorsRoot)

    // Cancelling the context passed to New() frees up resources held by the
    // client, closes connections, clears handlers, etc.
    cancel()
}
```

## Maintainers

Chris Berry: [@bez625](https://github.com/Bez625).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/theQRL/qrl-metrics-exporter/pkg/go-consensus-client/issues).

## License

[Apache-2.0](LICENSE) © 2020, 2021 Attestant Limited
