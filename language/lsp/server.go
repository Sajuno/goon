package lsp

import (
	"bufio"
	"io"
)

type Server interface {
	Close() error
	Stdin() io.WriteCloser
	Stdout() *bufio.Reader
}
