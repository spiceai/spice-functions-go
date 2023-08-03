# spice-functions-go

Spice.xyz Go Function handler

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
