package app

import (
	"go2web/internal/cli"
	"go2web/internal/ui"
)

func Run() {
	cfg := cli.GetFlags()
	
	ui.Print("URL: %s\n", cfg.URL)
	ui.Print("Search: %s\n", cfg.Search)
	ui.Print("Help: %t\n", cfg.Help)

	if cfg.Help || (cfg.URL == "" && cfg.Search == "") {
		cli.PrintHelp()
		return
	} 
	
	if cfg.URL != "" {
		url, err := cli.ValidateURL(cfg.URL)
		if err != nil {
			ui.Print("error: %v\n", err)
			return
		} else {
			ui.Print("\nURL mode selected\n")
			ui.Print("URL: %s\n", url)
			// call url mode functions
			return
		}
	}
	
	if cfg.Search != "" {
		// call search mode functions
	}
}