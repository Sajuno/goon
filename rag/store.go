package rag

import (
	"context"
	"github.com/sajuno/goon/golang"
)

type Store interface {
	SaveChunks(ctx context.Context, chunks []EmbeddedChunk) error
}

type EmbeddedChunk struct {
	golang.Chunk

	Vector []float32

	// https://platform.openai.com/tokenizer
	Tokens int
}
