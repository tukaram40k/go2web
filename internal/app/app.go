package app

import (
	"strings"

	"go2web/internal/cli"
	"go2web/internal/tcp"
	"go2web/internal/parser"
	"go2web/internal/search"
	"go2web/internal/ui"
)

func Run() {
	flags := cli.GetFlags()

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
		searchTerm := strings.TrimSpace(flags.Search)
		searchURL, err := search.BuildURL(searchTerm)
		if err != nil {
			ui.Print("invalid search term: %v\n", err)
			return
		}

		client := tcp.NewClient()
		resp, redirectCount, err := client.GetWithMeta(searchURL)
		if err != nil {
			ui.Print("search request failed: %v\n", err)
			return
		}

		parsedResp, err := parser.ParseWithRedirectInfo(resp, redirectCount)
		if err != nil {
			ui.Print("failed to parse search response: %v\n", err)
			return
		}

		if !parsedResp.ResponseIsOK {
			ui.Print("search request was not successful\n")
			ui.PrintParsedResponse(parsedResp)
			return
		}

		if !strings.Contains(strings.ToLower(parsedResp.ContentType), "text/html") {
			ui.Print("unexpected search response content type: %s\n", parsedResp.ContentType)
			return
		}

		results, err := search.ExtractResults(parsedResp.Body, 10)
		if err != nil {
			ui.Print("failed to extract search results: %v\n", err)
			return
		}

		ui.PrintSearchResults(searchTerm, parsedResp, results)
		return
	}
}
