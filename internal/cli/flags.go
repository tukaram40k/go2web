package cli

import (
	"flag"
	"net/url"
	"os"
	"strconv"
	"strings"

	"go2web/internal/ui"
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

	fs.StringVar(&c.URL, "url", "", "Website URL")
	fs.StringVar(&c.URL, "u", "", "")

	fs.StringVar(&c.Search, "search", "", "Search term")
	fs.StringVar(&c.Search, "s", "", "")

	fs.BoolVar(&c.Help, "help", false, "Display help information")
	fs.BoolVar(&c.Help, "h", false, "")

	fs.Usage = func() {
		PrintHelp()
	}

	fs.Parse(args)
	return c
}

func PrintHelp() {
	ui.Print("Options:\n")
	ui.Print("  -u, -url <URL>               # make an HTTP request to the specified URL and print the response\n")
	ui.Print("  -s, -search <search-term>    # make an HTTP request to search the term using your favorite search engine and print top 10 results\n")
	ui.Print("  -h, -help                    # show this help\n")
}

func ValidateURL(s string) (string, error) {
	if s == "" {
		return "", ui.Error("empty url given")
	}

	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", ui.Error("invalid url: %w", err)
	}

	if u.Host == "" {
		return "", ui.Error("invalid url: missing host")
	}

	// Validate port if present
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return "", ui.Error("invalid port number")
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
					return "", ui.Error("invalid IP address")
				}
			}
		} else {
			// It's a domain name — validate as hostname
			if len(parts) < 2 {
				return "", ui.Error("invalid domain name")
			}
			for _, part := range parts {
				if part == "" || strings.HasPrefix(part, "-") || part == "http" || part == "https" {
					return "", ui.Error("invalid domain name")
				}
			}
		}
	} else {
		// Not an IPv4 address, validate as hostname
		parts := strings.Split(hostname, ".")
		if len(parts) < 2 {
			return "", ui.Error("invalid domain name")
		}
		for _, part := range parts {
			if part == "" || strings.HasPrefix(part, "-") || part == "http" || part == "https" {
				return "", ui.Error("invalid domain name")
			}
		}
	}

	return s, nil
}
