package fetcher

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func getPathFromTime(t time.Time) string {
	return fmt.Sprintf("%d-%d-%d.csv", t.Year(), t.Month(), t.Day())
}

func writeCSV(path string, data [][]string) error {
	// Write CSV file

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(data); err != nil {
		return err
	}
	writer.Flush()

	return nil
}
