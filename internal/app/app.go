package app

import (
	"fmt"
	"go2web/internal/cli"
)

func Run() {
	cfg := cli.GetFlags()
	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Search Results: %s\n", cfg.Search)
	fmt.Printf("Help: %t\n", cfg.Help)
}