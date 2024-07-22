package fetcher

import (
	"testing"
)

func TestWriteCSV(t *testing.T) {
	// Write CSV file

	err := writeCSV("test.csv", [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
	})

	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
