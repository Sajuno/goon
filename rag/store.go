package rag

import (
	"context"
	"github.com/sajuno/goon/language/golang"
)

type Store interface {
	SaveChunks(ctx context.Context, chunks []Chunk) error
	FindSimilarChunks(ctx context.Context, vector []float32) ([]SimilarChunk, error)
}

type Chunk struct {
	golang.Chunk

	Vector []float32

	// https://platform.openai.com/tokenizer
	Tokens int
}

// SimilarChunk is returned from FindSimilarChunks
// In addition to the Chunk itself, it contains the distance to the prompt
type SimilarChunk struct {
	Chunk

	Distance float64
}
