package cli

import "flag"

type Config struct {
	URL string
	Search int
	Help bool
}

func GetFlags() *Config {
	c := &Config{}

	flag.StringVar(&c.URL, "url", "http://example.com", "Website URL")
	flag.IntVar(&c.Search, "search", 10, "Number of search results")
	flag.BoolVar(&c.Help, "help", false, "Display help information")

	flag.Parse()
	return c
}
