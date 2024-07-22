package main

import (
	"fmt"
	"time"
)

func main() {
	// Code
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Println(now)
}
