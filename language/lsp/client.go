package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Client struct {
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

func NewClient(s Server) (*Client, error) {
	client := &Client{
		stdin:  s.Stdin(),
		stdout: s.Stdout(),
	}

	if err := client.initialize(); err != nil {
		s.Close()
		return nil, err
	}

	return client, nil
}

func (c *Client) send(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	head := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := c.stdin.Write([]byte(head)); err != nil {
		return err
	}

	_, err = c.stdin.Write(data)
	return err
}

func (c *Client) read() (*Message, error) {
	header := ""
	for {
		line, err := c.stdout.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			break
		}
		header += line
	}

	headers := parseHeaders(header)
	length, _ := strconv.Atoi(headers["Content-Length"])
	body := make([]byte, length)
	_, err := io.ReadFull(c.stdout, body)
	if err != nil {
		return nil, err
	}

	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func parseHeaders(h string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(h, "\r\n")

	for _, l := range lines {
		parts := strings.SplitN(l, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	return headers
}

func (c *Client) initialize() error {
	params := map[string]interface{}{
		"processId":    nil,
		"rootUri":      nil,
		"capabilities": map[string]any{},
	}
	paramBytes, _ := json.Marshal(params)
	if err := c.send(newMessage("initialize", paramBytes)); err != nil {
		return err
	}

	_, err := c.read()
	if err != nil {
		return err
	}

	return c.send(newMessage("initialize", nil))
}
