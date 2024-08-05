package client

import (
	"net"
	"strings"
)

type Client struct {
	socket net.Conn
}

func NewClient(path string) (*Client, error) {
	s, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	return &Client{socket: s}, nil
}

func (c *Client) Request(host string) (keypath string, err error) {
	buf := make([]byte, 4096)

	_, err = c.socket.Write([]byte(host))
	if err != nil {
		return "", err
	}

	n, err := c.socket.Read(buf)
	if err != nil {
		return "", err
	}

	return strings.Replace(string(buf[:n]), "\n", "", -1), nil
}

func (c *Client) Acknowledge() error {
	_, err := c.socket.Write([]byte("OK"))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Failed() error {
	_, err := c.socket.Write([]byte("FAILED"))
	if err != nil {
		return err
	}

	return nil
}
