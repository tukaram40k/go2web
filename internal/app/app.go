package app

import (
	"fmt"
	"go2web/internal/cli"
)

func Run() {
	cfg := cli.GetFlags()
	
	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("Search: %s\n", cfg.Search)
	fmt.Printf("Help: %t\n", cfg.Help)

	if cfg.Help {
		cli.PrintHelp()
		return
	} 
	
	if cfg.URL != "" {
		url, err := cli.ValidateURL(cfg.URL)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		} else {
			fmt.Printf("URL: %s\n", url)
			// call url mode functions
			return
		}
	}
	
	if cfg.Search != "" {
		// call search mode functions
	}
}