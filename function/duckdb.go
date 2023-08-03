package function

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/marcboeker/go-duckdb"
)

func openDuckDb(ctx context.Context, dbFile string) (*sql.DB, error) {
	connector, err := duckdb.NewConnector(dbFile, func(execer driver.ExecerContext) error {
		bootQueries := []string{
			"INSTALL 'parquet'",
			"LOAD 'parquet'",
			"INSTALL 'json'",
			"LOAD 'json'",
		}

		for _, qry := range bootQueries {
			_, err := execer.ExecContext(ctx, qry, nil)
			if err != nil {
				return fmt.Errorf("error executing boot query: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error creating duckdb connector: %w", err)
	}

	return sql.OpenDB(connector), nil
}
