package main

import (
	"fmt"

	"go2web/internal/cli"
)

func main() {
	cfg := cli.GetFlags()
	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Search Results: %s\n", cfg.Search)
	fmt.Printf("Help: %t\n", cfg.Help)
}
