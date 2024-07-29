package fetcher

import "time"

type Data struct {
	Time     time.Time
	Name     string
	Number   string
	Capacity string
}

func DataToCsv(data []Data) [][]string {
	csv := make([][]string, 0)
	for _, d := range data {
		csv = append(csv, []string{
			d.Time.Format(time.RFC3339),
			d.Name,
			d.Number,
			d.Capacity,
		})
	}
	return csv
}
