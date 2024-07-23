package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/chechiachang/sc-stat/pkg/utils"
)

func main() {
	date := time.Now()
	files := utils.Glob("data", func(s string) bool {
		return filepath.Ext(s) == ".csv" &&
			(strings.Contains(s, date.Format("2006-1-2")) ||
				strings.Contains(s, date.AddDate(0, 0, -1).Format("2006-1-2"))) // today and yesterday
	})
	fmt.Println(strings.Join(files, ","))
}
