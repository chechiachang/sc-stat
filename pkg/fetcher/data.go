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
			d.Time.Format("2006-01-02 15:04:05"),
			d.Name,
			d.Number,
			d.Capacity,
		})
	}
	return csv
}
