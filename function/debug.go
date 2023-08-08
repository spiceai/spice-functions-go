package function

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/gofrs/flock"
	"github.com/spiceai/gospice/v2"
)

type debugOptions struct {
	inputsDir        string
	dataDir          string
	outputsDir       string
	flightAddress    string
	firecacheAddress string
	apiKey           string

	chainName   string
	blockNumber int64
	blockHash   string
}

func defaultDebugOptions() *debugOptions {
	return &debugOptions{
		inputsDir:        "./inputs",
		dataDir:          "./data",
		outputsDir:       "./outputs",
		flightAddress:    "flight.spiceai.io:443",
		firecacheAddress: "firecache.spiceai.io:443",
		apiKey:           os.Getenv("SPICE_API_KEY"),
		chainName:        "eth",
	}
}

func WithInputsDir(inputsDir string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.inputsDir = inputsDir
	}
}

func WithDataDir(dataDir string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.dataDir = dataDir
	}
}

func WithOutputsDir(outputsDir string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.outputsDir = outputsDir
	}
}

func WithFlightAddress(flightAddress string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.flightAddress = flightAddress
	}
}

func WithFirecacheAddress(firecacheAddress string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.firecacheAddress = firecacheAddress
	}
}

func WithApiKey(apiKey string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.apiKey = apiKey
	}
}

func WithChain(chain string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.chainName = chain
	}
}

func WithBlockNumber(blockNumber int64) func(*debugOptions) {
	return func(o *debugOptions) {
		o.blockNumber = blockNumber
	}
}

func WithBlockHash(blockHash string) func(*debugOptions) {
	return func(o *debugOptions) {
		o.blockHash = blockHash
	}
}

func Debug(handler func(ctx *FunctionCtx, duckDb *sql.DB, spiceClient *gospice.SpiceClient) error, options ...func(*debugOptions)) (*sql.DB, error) {
	opts := defaultDebugOptions()
	for _, opt := range options {
		opt(opts)
	}

	err := mkdir(opts.inputsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create inputs dir: %w", err)
	}
	err = mkdir(opts.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}
	err = mkdir(opts.outputsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create outputs dir: %w", err)
	}

	spiceClient := gospice.NewSpiceClientWithAddress(opts.flightAddress, opts.firecacheAddress)
	err = spiceClient.Init(opts.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Spice client: %w", err)
	}

	lockFile := flock.New(filepath.Join(opts.dataDir, "persistent_data.lock"))
	err = lockFile.Lock()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire persistent data lock: %w", err)
	}
	defer lockFile.Unlock() // nolint: errcheck

	blockNumber := opts.blockNumber
	blockHash := opts.blockHash

	if blockNumber == 0 {
		rr, err := spiceClient.Query(context.Background(), fmt.Sprintf("SELECT MAX(number) from %s.recent_blocks", opts.chainName))
		if err != nil {
			return nil, fmt.Errorf("failed to query recent blocks: %w", err)
		}
		rr.Next()
		blockNumber = rr.Record().Column(0).(*array.Int64).Value(0)
		rr.Release()
	}

	if blockHash == "" {
		rr, err := spiceClient.Query(context.Background(), fmt.Sprintf("SELECT hash from %s.blocks WHERE number = %d", opts.chainName, blockNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to query recent blocks: %w", err)
		}
		rr.Next()
		blockHash = rr.Record().Column(0).(*array.String).Value(0)
		rr.Release()
	}

	functionCtx := NewFunctionCtx(context.Background(), blockNumber, blockHash)

	persistentDataDb := filepath.Join(opts.dataDir, "persistent_data.duckdb")
	duckDb, err := openDuckDb(functionCtx, persistentDataDb)
	if err != nil {
		return nil, fmt.Errorf("failed to open persistent data duckdb: %w", err)
	}

	inputDb := filepath.Join(opts.inputsDir, "input.duckdb")
	_, err = duckDb.ExecContext(functionCtx, fmt.Sprintf("ATTACH '%s' AS input", inputDb))
	if err != nil {
		return nil, fmt.Errorf("failed to attach input duckdb: %w", err)
	}
	outputDb := filepath.Join(opts.outputsDir, "output.duckdb")
	err = rmIfExists(outputDb)
	if err != nil {
		return nil, fmt.Errorf("failed to remove output duckdb: %w", err)
	}
	_, err = duckDb.ExecContext(functionCtx, fmt.Sprintf("ATTACH '%s' AS output", outputDb))
	if err != nil {
		return nil, fmt.Errorf("failed to attach output duckdb: %w", err)
	}

	err = handler(functionCtx, duckDb, spiceClient)
	if err != nil {
		return nil, err
	}

	return duckDb, nil
}

func mkdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}
	}
	return nil
}

func rmIfExists(path string) error {
	if _, err := os.Stat(path); err == nil {
		err = os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("failed to remove path: %w", err)
		}
	}
	return nil
}
