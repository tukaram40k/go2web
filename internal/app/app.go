package app

import (
	"go2web/internal/cli"
	"go2web/internal/ui"
	"go2web/internal/tcp"
)

func Run() {
	flags := cli.GetFlags()

	ui.Print("URL: %s\n", flags.URL)
	ui.Print("Search: %s\n", flags.Search)
	ui.Print("Help: %t\n", flags.Help)

	if flags.Help || (flags.URL == "" && flags.Search == "") {
		cli.PrintHelp()
		return
	}

	if flags.URL != "" {
		url, err := cli.ValidateURL(flags.URL)
		if err != nil {
			ui.Print("error: %v\n", err)
			return
		} else {
			ui.Print("\nURL mode selected\n")
			ui.Print("URL: %s\n", url)

			// call url mode functions
			client := tcp.NewClient()
			resp, err := client.Get(url)
			if err != nil {
				ui.Print("request failed: %v\n", err)
				return
			}

			ui.Print("response:\n%s\n", string(resp))

			return
		}
	}

	if flags.Search != "" {
		// call search mode functions
	}
}
