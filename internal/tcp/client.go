package tcp

import (
	"bufio"
	"fmt"
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
	path := u.Path
	if path == "" {
		path = "/"
	}

	addr := host + ":80"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
		path,
		host,
	)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadBytes(0) // temporary; we'll fix later
	if err != nil {
		// EOF is expected when connection closes
	}

	return resp, nil
}