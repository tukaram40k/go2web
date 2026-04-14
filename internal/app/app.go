package app

import (
	"go2web/internal/cli"
	"go2web/internal/parser"
	"go2web/internal/tcp"
	"go2web/internal/ui"
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
			resp, redirectCount, err := client.GetWithMeta(url)
			if err != nil {
				ui.Print("request failed: %v\n", err)
				return
			}

			parsedResp, err := parser.ParseWithRedirectInfo(resp, redirectCount)
			if err != nil {
				ui.Print("failed to parse response: %v\n", err)
				return
			}

			ui.PrintParsedResponse(parsedResp)

			logPath, err := ui.Log(parsedResp)
			if err != nil {
				ui.Print("failed to save response log: %v\n", err)
				return
			}

			ui.Print("saved response to: %s\n", logPath)

			return
		}
	}

	if flags.Search != "" {
		// call search mode functions
	}
}
