package main

import "github.com/chechiachang/sc-stat/pkg/git"

func main() {
	url := "https://github.com/chechiachang/sc-stat-data"
	directory := "data"
	git.GitInit(url, directory)
}
