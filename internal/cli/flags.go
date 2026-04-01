package cli

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	URL    string
	Search string
	Help   bool
}

func GetFlags() *Config {
	c := &Config{}

	flag.StringVar(&c.URL, "url", "http://example.com", "Website URL")
	flag.StringVar(&c.URL, "u", "http://example.com", "")

	flag.StringVar(&c.Search, "search", "search term", "Search term")
	flag.StringVar(&c.Search, "s", "search term", "")

	flag.BoolVar(&c.Help, "help", false, "Display help information")
	flag.BoolVar(&c.Help, "h", false, "")

	// custom usage function
	flag.Usage = func() {
		// fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -u, -url <URL>               # make an HTTP request to the specified URL and print the response\n")
		fmt.Fprintf(os.Stderr, "  -s, -search <search-term>    # make an HTTP request to search the term using your favorite search engine and print top 10 results\n")
		fmt.Fprintf(os.Stderr, "  -h, -help                    # show this help\n")
	}

	flag.Parse()
	return c
}