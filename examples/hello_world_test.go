package examples

import (
	"testing"

	"github.com/spiceai/spice-functions-go/function"
)

func TestHelloWorld(t *testing.T) {
	outputDb, err := function.Debug(HelloWorldGo)
	if err != nil {
		t.Error(err)
		return
	}
	defer outputDb.Close()

	rows, err := outputDb.Query("SELECT * FROM output.hello_world_golang;")
	if err != nil {
		t.Error(err)
		return
	}

	var blockNumber int64
	var greeting string
	var numRows int
	for rows.Next() {
		err := rows.Scan(&blockNumber, &greeting)
		if err != nil {
			t.Error(err)
			return
		}
		numRows++
	}

	if numRows != 1 {
		t.Errorf("Expected 1 row, got %d", numRows)
	}
	if blockNumber <= 17500000 {
		t.Errorf("Expected a recent eth block number, got %d", blockNumber)
	}
	if greeting != "Hello from Spice Go runtime!" {
		t.Errorf("Expected greeting 'Hello from Spice Go runtime!', got '%s'", greeting)
	}
}
