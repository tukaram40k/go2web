package tcp

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
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

	addr := host + ":80"

	conn, err := net.Dial("tcp", addr)
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
