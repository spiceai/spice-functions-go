package examples

import (
	"database/sql"
	"fmt"

	"github.com/spiceai/gospice/v2"
	"github.com/spiceai/spice-functions-go/function"
)

func HelloWorldGo(ctx *function.FunctionCtx, duckDb *sql.DB, client *gospice.SpiceClient) error {
	fmt.Println("Hello from Spice Go runtime!")

	_, err := duckDb.ExecContext(ctx, "INSERT INTO output.hello_world_golang (block_number, greeting) VALUES ($1, $2);", ctx.BlockNumber(), "Hello from Spice Go runtime!")
	if err != nil {
		return err
	}

	return nil
}
