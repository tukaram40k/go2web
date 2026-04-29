package app

import (
	"strings"
	"time"

	"go2web/internal/cache"
	"go2web/internal/cli"
	"go2web/internal/parser"
	"go2web/internal/search"
	"go2web/internal/tcp"
	"go2web/internal/ui"
)

func Run() {
	flags := cli.GetFlags()

	if flags.Help || (flags.URL == "" && flags.Search == "") || (flags.URL != "" && flags.Search != "") {
		cli.PrintHelp()
		return
	}

	if flags.URL != "" {
		url, err := cli.ValidateURL(flags.URL)
		if err != nil {
			ui.Print("error: %v\n", err)
			return
		} else {
			resp, redirectCount, cached, cacheAgeMinutes, cacheErr, err := fetchWithCache(url)
			if err != nil {
				ui.Print("request failed: %v\n", err)
				return
			}

			if cacheErr != nil {
				ui.Print("failed to save cache: %v\n", cacheErr)
			}

			parsedResp, err := parser.ParseWithRedirectInfo(resp, redirectCount)
			if err != nil {
				ui.Print("failed to parse response: %v\n", err)
				return
			}
			parsedResp.Cached = cached
			parsedResp.CacheAgeMinutes = cacheAgeMinutes

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

		resp, redirectCount, cached, cacheAgeMinutes, cacheErr, err := fetchWithCache(searchURL)
		if err != nil {
			ui.Print("search request failed: %v\n", err)
			return
		}

		if cacheErr != nil {
			ui.Print("failed to save cache: %v\n", cacheErr)
		}

		parsedResp, err := parser.ParseWithRedirectInfo(resp, redirectCount)
		if err != nil {
			ui.Print("failed to parse search response: %v\n", err)
			return
		}
		parsedResp.Cached = cached
		parsedResp.CacheAgeMinutes = cacheAgeMinutes

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

		logPath, err := ui.LogSearchResults(searchTerm, parsedResp, results)
		if err != nil {
			ui.Print("failed to save search log: %v\n", err)
			return
		}

		ui.Print("saved search log to: %s\n", logPath)
		return
	}
}

func fetchWithCache(rawURL string) ([]byte, int, bool, int, error, error) {
	if cachedBody, entry, found, err := cache.Load(rawURL); err != nil {
		return nil, 0, false, 0, nil, err
	} else if found {
		ageMinutes := 0
		if !entry.CachedAt.IsZero() {
			age := time.Since(entry.CachedAt)
			if age > 0 {
				ageMinutes = int(age.Minutes())
			}
		}
		return cachedBody, entry.RedirectCount, true, ageMinutes, nil, nil
	}

	client := tcp.NewClient()
	body, redirectCount, err := client.GetWithMeta(rawURL)
	if err != nil {
		return nil, 0, false, 0, nil, err
	}

	cacheErr := cache.Store(rawURL, body, redirectCount)
	return body, redirectCount, false, 0, cacheErr, nil
}
