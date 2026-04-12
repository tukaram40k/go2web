package tcp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) dial(u *url.URL) (net.Conn, error) {
	hostname := u.Hostname()
	if hostname == "" {
		return nil, fmt.Errorf("invalid url: missing host")
	}

	scheme := u.Scheme
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

func (c *Client) Get(rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

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
