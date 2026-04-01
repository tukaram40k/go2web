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
	return GetFlagsFromSet(flag.CommandLine, os.Args[1:])
}

func GetFlagsFromSet(fs *flag.FlagSet, args []string) *Config {
	c := &Config{}

	fs.StringVar(&c.URL, "url", "http://example.com", "Website URL")
	fs.StringVar(&c.URL, "u", "http://example.com", "")

	fs.StringVar(&c.Search, "search", "search term", "Search term")
	fs.StringVar(&c.Search, "s", "search term", "")

	fs.BoolVar(&c.Help, "help", false, "Display help information")
	fs.BoolVar(&c.Help, "h", false, "")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -u, -url <URL>               # make an HTTP request to the specified URL and print the response\n")
		fmt.Fprintf(os.Stderr, "  -s, -search <search-term>    # make an HTTP request to search the term using your favorite search engine and print top 10 results\n")
		fmt.Fprintf(os.Stderr, "  -h, -help                    # show this help\n")
	}

	fs.Parse(args)
	return c
}
