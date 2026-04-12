package tcp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Client struct{}

const defaultMaxRedirects = 10

func NewClient() *Client {
	return &Client{}
}

func (c *Client) dial(u *url.URL) (net.Conn, error) {
	hostname := u.Hostname()
	if hostname == "" {
		return nil, fmt.Errorf("invalid url: missing host")
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme == "" {
		scheme = "http"
	}

	port := u.Port()

	switch scheme {
	case "http":
		if port == "" {
			port = "80"
		}

		addr := net.JoinHostPort(hostname, port)
		return net.Dial("tcp", addr)
	case "https":
		if port == "" {
			port = "443"
		}

		addr := net.JoinHostPort(hostname, port)
		dialer := &net.Dialer{}
		return tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: hostname})
	default:
		return nil, fmt.Errorf("unsupported url scheme: %s", scheme)
	}
}

func (c *Client) getOnce(u *url.URL) ([]byte, error) {
	host := u.Host
	requestTarget := u.EscapedPath()
	if requestTarget == "" {
		requestTarget = "/"
	}

	if u.RawQuery != "" {
		requestTarget += "?" + u.RawQuery
	}

	conn, err := c.dial(u)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
		requestTarget,
		host,
	)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return nil, err
	}

	resp, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func parseStatusAndLocation(resp []byte) (int, string, error) {
	headerEnd := bytes.Index(resp, []byte("\r\n\r\n"))
	lineSeparator := "\r\n"
	if headerEnd == -1 {
		headerEnd = bytes.Index(resp, []byte("\n\n"))
		lineSeparator = "\n"
	}

	if headerEnd == -1 {
		return 0, "", fmt.Errorf("invalid http response: missing header separator")
	}

	lines := strings.Split(string(resp[:headerEnd]), lineSeparator)
	if len(lines) == 0 {
		return 0, "", fmt.Errorf("invalid http response: missing status line")
	}

	statusParts := strings.SplitN(strings.TrimSpace(lines[0]), " ", 3)
	if len(statusParts) < 2 {
		return 0, "", fmt.Errorf("invalid http status line: %q", lines[0])
	}

	statusCode, err := strconv.Atoi(statusParts[1])
	if err != nil {
		return 0, "", fmt.Errorf("invalid status code in response: %w", err)
	}

	location := ""
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.EqualFold(strings.TrimSpace(parts[0]), "Location") {
			location = strings.TrimSpace(parts[1])
			break
		}
	}

	return statusCode, location, nil
}

func isRedirectStatus(statusCode int) bool {
	switch statusCode {
	case 301, 302, 303, 307, 308:
		return true
	default:
		return false
	}
}

func resolveRedirectURL(current *url.URL, location string) (*url.URL, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return nil, fmt.Errorf("empty redirect location")
	}

	nextURL, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("invalid redirect location %q: %w", location, err)
	}

	return current.ResolveReference(nextURL), nil
}

func (c *Client) GetWithRedirects(rawURL string, maxRedirects int) ([]byte, error) {
	if maxRedirects < 0 {
		return nil, fmt.Errorf("max redirects cannot be negative")
	}

	currentURL := rawURL
	visited := make(map[string]struct{})

	for redirectsFollowed := 0; ; redirectsFollowed++ {
		u, err := url.Parse(currentURL)
		if err != nil {
			return nil, err
		}

		normalizedURL := u.String()
		if _, seen := visited[normalizedURL]; seen {
			return nil, fmt.Errorf("redirect loop detected at %s", normalizedURL)
		}
		visited[normalizedURL] = struct{}{}

		resp, err := c.getOnce(u)
		if err != nil {
			return nil, err
		}

		statusCode, location, err := parseStatusAndLocation(resp)
		if err != nil {
			return resp, nil
		}

		if !isRedirectStatus(statusCode) {
			return resp, nil
		}

		if location == "" {
			return nil, fmt.Errorf("redirect response missing Location header")
		}

		if redirectsFollowed >= maxRedirects {
			return nil, fmt.Errorf("maximum redirects exceeded (%d)", maxRedirects)
		}

		nextURL, err := resolveRedirectURL(u, location)
		if err != nil {
			return nil, err
		}

		currentURL = nextURL.String()
	}
}

func (c *Client) Get(rawURL string) ([]byte, error) {
	return c.GetWithRedirects(rawURL, defaultMaxRedirects)
}
