package gopls

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sajuno/goon/language/lsp"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Client struct {
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr *bufio.Scanner

	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func NewClient(ctx context.Context) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(ctx, "gopls", "-mode=stdio")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	client := &Client{
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		stderr: bufio.NewScanner(stderrPipe),
		cmd:    cmd,
		cancel: cancel,
	}

	// small side car process for error logging
	go client.logStderr()

	if err := client.initialize(); err != nil {
		cancel()
		return nil, err
	}

	return client, nil
}

func (c *Client) logStderr() {
	for c.stderr.Scan() {
		log.Printf("[gopls] %s", c.stderr.Text())
	}
}

func (c *Client) Close() error {
	c.cancel()
	return c.cmd.Wait()
}

func (c *Client) send(msg *lsp.Message) error {
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

func (c *Client) read() (*lsp.Message, error) {
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

	var msg lsp.Message
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
		"capabilities": map[string]interface{}{},
	}
	paramBytes, _ := json.Marshal(params)
	msg := &lsp.Message{
		ID:     uuid.NewString(),
		Method: "initialize",
		Params: paramBytes,
	}
	if err := c.send(msg); err != nil {
		return err
	}

	_, err := c.read()
	if err != nil {
		return err
	}

	return c.send(&lsp.Message{Method: "initialized"})
}
