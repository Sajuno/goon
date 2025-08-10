package gopls

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
)

type Server struct {
	stdin  io.WriteCloser
	stdout *bufio.Reader

	stderr *bufio.Scanner

	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func Start(ctx context.Context) (*Server, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
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

	pipe := &Server{
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		stderr: bufio.NewScanner(stderrPipe),
		cmd:    cmd,
		cancel: cancel,
	}

	// small side car process for error logging
	go pipe.logStderr()

	return pipe, nil
}

func (c *Server) logStderr() {
	// TODO: pipe this to some kind of error channel instead
	for c.stderr.Scan() {
		log.Printf(c.stderr.Text())
	}
}

func (c *Server) Close() error {
	c.cancel()
	return c.cmd.Wait()
}

func (c *Server) Stdin() io.WriteCloser {
	return c.stdin
}

func (c *Server) Stdout() *bufio.Reader {
	return c.stdout
}
