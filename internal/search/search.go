package search

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

const (
	duckDuckGoHTMLBaseURL = "https://duckduckgo.com/html/?q="
	defaultResultsLimit   = 10
)

var (
	resultAnchorRE = regexp.MustCompile(`(?is)<a[^>]*class=["'][^"']*result__a[^"']*["'][^>]*href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	tagRE          = regexp.MustCompile(`(?is)<[^>]+>`)
)

type Result struct {
	Title string
	URL   string
}

func BuildURL(term string) (string, error) {
	q := strings.TrimSpace(term)
	if q == "" {
		return "", fmt.Errorf("empty search term")
	}

	return duckDuckGoHTMLBaseURL + url.QueryEscape(q), nil
}

func sanitizeTitle(raw string) string {
	text := tagRE.ReplaceAllString(raw, "")
	text = html.UnescapeString(text)
	text = strings.TrimSpace(text)

	return strings.Join(strings.Fields(text), " ")
}

func normalizeResultURL(raw string) string {
	href := strings.TrimSpace(html.UnescapeString(raw))
	if href == "" {
		return ""
	}

	if strings.HasPrefix(href, "//") {
		href = "https:" + href
	}

	parsed, err := url.Parse(href)
	if err != nil {
		return ""
	}

	if parsed.Scheme == "http" || parsed.Scheme == "https" {
		if strings.Contains(parsed.Hostname(), "duckduckgo.com") && strings.HasPrefix(parsed.Path, "/l/") {
			if uddg := parsed.Query().Get("uddg"); uddg != "" {
				decoded, err := url.QueryUnescape(uddg)
				if err == nil {
					if target, targetErr := url.Parse(decoded); targetErr == nil && (target.Scheme == "http" || target.Scheme == "https") {
						return target.String()
					}
				}
			}
		}

		return parsed.String()
	}

	if strings.HasPrefix(parsed.Path, "/l/") {
		if uddg := parsed.Query().Get("uddg"); uddg != "" {
			decoded, err := url.QueryUnescape(uddg)
			if err == nil {
				if target, targetErr := url.Parse(decoded); targetErr == nil && (target.Scheme == "http" || target.Scheme == "https") {
					return target.String()
				}
			}
		}
	}

	return ""
}

func ExtractResults(body []byte, limit int) ([]Result, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("empty search response body")
	}

	if limit <= 0 {
		limit = defaultResultsLimit
	}

	matches := resultAnchorRE.FindAllSubmatch(body, -1)
	results := make([]Result, 0, limit)
	seenURLs := make(map[string]struct{})

	for _, m := range matches {
		if len(m) < 3 {
			continue
		}

		normalizedURL := normalizeResultURL(string(m[1]))
		title := sanitizeTitle(string(m[2]))
		if normalizedURL == "" || title == "" {
			continue
		}

		if _, seen := seenURLs[normalizedURL]; seen {
			continue
		}

		seenURLs[normalizedURL] = struct{}{}
		results = append(results, Result{Title: title, URL: normalizedURL})
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}
