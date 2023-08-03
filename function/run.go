package function

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/spiceai/gospice/v2"
	"gopkg.in/yaml.v3"
)

type FunctionContextVariables struct {
	BlockNumber int64  `mapstructure:"block_number,omitempty" json:"block_number,omitempty" yaml:"block_number,omitempty"`
	BlockHash   string `mapstructure:"block_hash,omitempty" json:"block_hash,omitempty" yaml:"block_hash,omitempty"`
}

func Run(handler func(ctx *FunctionCtx, duckDb *sql.DB, spiceClient *gospice.SpiceClient) error) {
	inputsDir := os.Getenv("INPUT_DIR")
	dataDir := os.Getenv("DATA_DIR")
	outputsDir := os.Getenv("OUTPUT_DIR")
	flightAddress := os.Getenv("SPICE_FLIGHT_ADDRESS")
	firecacheAddress := os.Getenv("SPICE_FIRECACHE_ADDRESS")
	apiKey := os.Getenv("SPICE_API_KEY")

	client := gospice.NewSpiceClientWithAddress(flightAddress, firecacheAddress)
	err := client.Init(apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize Spice client: %s", err)
	}

	lockFile := flock.New(filepath.Join(dataDir, "persistent_data.lock"))
	err = lockFile.Lock()
	if err != nil {
		log.Fatalf("Failed to acquire persistent data lock: %s", err)
	}
	defer lockFile.Unlock()

	contextYamlBytes, err := os.ReadFile(filepath.Join(inputsDir, "context.yaml"))
	if err != nil {
		log.Fatalf("Failed to read context.yaml: %s", err)
	}

	var contextVars FunctionContextVariables
	err = yaml.Unmarshal(contextYamlBytes, &contextVars)
	if err != nil {
		log.Fatalf("Failed to unmarshal context.yaml: %s", err)
	}

	functionCtx := NewFunctionCtx(context.Background(), contextVars.BlockNumber, contextVars.BlockHash)

	persistentDataDb := filepath.Join(dataDir, "persistent_data.duckdb")
	duckDb, err := openDuckDb(functionCtx, persistentDataDb)
	if err != nil {
		log.Fatalf("Failed to open persistent data duckdb: %s", err)
	}
	defer duckDb.Close()

	inputDb := filepath.Join(inputsDir, "input.duckdb")
	_, err = duckDb.ExecContext(functionCtx, fmt.Sprintf("ATTACH '%s' AS input", inputDb))
	if err != nil {
		log.Fatalf("Failed to attach input duckdb: %s", err)
	}
	outputDb := filepath.Join(outputsDir, "output.duckdb")
	_, err = duckDb.ExecContext(functionCtx, fmt.Sprintf("ATTACH '%s' AS output", outputDb))
	if err != nil {
		log.Fatalf("Failed to attach output duckdb: %s", err)
	}

	err = handler(functionCtx, duckDb, client)
	if err != nil {
		panic(err)
	}
}
