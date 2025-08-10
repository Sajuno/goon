package lsp

import (
	"context"
	"github.com/sajuno/goon/language/lsp/gopls"
)

func NewGoplsClient(ctx context.Context) (*Client, error) {
	s, err := gopls.Start(ctx)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(s)
	if err != nil {
		return nil, err
	}

	return client, nil
}
