# spice-functions-go

Spice.ai Go Function handler

# Getting Started

1. Get the spice-functions-go package.

```bash
go get github.com/spiceai/spice-functions-go
```

2. Wrap the function handler with the spice-functions-go handler.

```go
// main.go
package main

import (
	"fmt"
	"github.com/spiceai/spice-functions-go/function"
	"github.com/spiceai/gospice/v2"
)

func hello(ctx *function.FunctionCtx, duckDb *sql.DB, spiceClient *gospice.SpiceClient) error {
	fmt.Println("Hello Spice Functions!")
	return nil
}

func main() {
	function.Run(hello)
}
```

# Local Development

You can debug your function locally by calling `function.Debug()` instead of `function.Run()` from within a test.

```go
// main_test.go
package main

import (
	"testing"
	"github.com/spiceai/spice-functions-go/function"
)

func TestHello(t *testing.T) {
	// Debug returns the DuckDB database that was used in the function, as well as any errors.
	// Close the database when you're done with it.
	outputDb, err := function.Debug(hello, function.WithApiKey("<your-api-key>"))
	if err != nil {
		t.Error(err)
	}
	defer outputDb.Close()

	// Use the outputDb to query the results of your function.
}
```

You can pass a number of options to `function.Debug()` to configure the function's execution.

- `function.WithApiKey("<spice-api-key>")`: **Required**. Sets the Spice API key to use when running the function. Defaults to the `SPICE_API_KEY` environment variable.
- `function.WithOutputDatasetMigration("CREATE TABLE ...")`: Adds a migration SQL to use to create the output dataset. Can be specified multiple times to add multiple migrations.
- `function.WithInputsDir("./inputs")`: Sets the path to the inputs directory (currently unused). Defaults to `./inputs`.
- `function.WithDataDir("./data")`: Sets the path to the persistent data directory. Defaults to `./data`.
- `function.WithOutputsDir("./outputs")`: Sets the path to the outputs directory. Defaults to `./outputs`.
- `function.WithPathTrigger("eth")`: Sets the path this function is triggered on. Defaults to `eth`.
- `function.WithBlockNumber(17400000)`: Sets the block number to use. Defaults to the latest block from the chain.
- `function.WithBlockHash("0x...")`: Sets the block hash to use. Defaults to the latest block from the chain.
