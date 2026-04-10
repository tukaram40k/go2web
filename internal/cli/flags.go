package cli

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	URL    string
	Search string
	Help   bool
}

// parse flags and return struct
func GetFlags() *Config {
	return GetFlagsFromSet(flag.CommandLine, os.Args[1:])
}

// helper
func GetFlagsFromSet(fs *flag.FlagSet, args []string) *Config {
	c := &Config{}

	fs.StringVar(&c.URL, "url", "http://example.com", "Website URL")
	fs.StringVar(&c.URL, "u", "http://example.com", "")

	fs.StringVar(&c.Search, "search", "search term", "Search term")
	fs.StringVar(&c.Search, "s", "search term", "")

	fs.BoolVar(&c.Help, "help", false, "Display help information")
	fs.BoolVar(&c.Help, "h", false, "")

	fs.Usage = func() {
		PrintHelp()
	}

	fs.Parse(args)
	return c
}

func PrintHelp() {
	fmt.Fprintf(os.Stderr, "Options:\n")
	fmt.Fprintf(os.Stderr, "  -u, -url <URL>               # make an HTTP request to the specified URL and print the response\n")
	fmt.Fprintf(os.Stderr, "  -s, -search <search-term>    # make an HTTP request to search the term using your favorite search engine and print top 10 results\n")
	fmt.Fprintf(os.Stderr, "  -h, -help                    # show this help\n")
}

func ValidateURL(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("empty url given")
	}

	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	if u.Host == "" {
		return "", fmt.Errorf("invalid url: missing host")
	}

	// Validate port if present
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return "", fmt.Errorf("invalid port number")
		}
	}

	hostname := u.Hostname()

	// Check for IPv4 literals and validate octets
	if strings.Contains(hostname, ".") && strings.Count(hostname, ".") == 3 {
		parts := strings.Split(hostname, ".")
		allNumeric := true
		for _, part := range parts {
			if _, err := strconv.Atoi(part); err != nil {
				allNumeric = false
				break
			}
		}
		if allNumeric {
			// It's an IPv4 address — validate octets
			for _, part := range parts {
				n, _ := strconv.Atoi(part)
				if n < 0 || n > 255 {
					return "", fmt.Errorf("invalid IP address")
				}
			}
		} else {
			// It's a domain name — validate as hostname
			if len(parts) < 2 {
				return "", fmt.Errorf("invalid domain name")
			}
			for _, part := range parts {
				if part == "" || strings.HasPrefix(part, "-") || part == "http" || part == "https" {
					return "", fmt.Errorf("invalid domain name")
				}
			}
		}
	} else {
		// Not an IPv4 address, validate as hostname
		parts := strings.Split(hostname, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid domain name")
		}
		for _, part := range parts {
			if part == "" || strings.HasPrefix(part, "-") || part == "http" || part == "https" {
				return "", fmt.Errorf("invalid domain name")
			}
		}
	}

	return s, nil
}
