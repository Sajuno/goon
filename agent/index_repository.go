package agent

import (
	"context"
	"fmt"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sajuno/goon/golang"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
	"log"
)

func (a *Agent) IndexRepository(ctx context.Context, path string) error {
	chunks, err := golang.ChunkRepository(path)
	if err != nil {
		return err
	}

	embeddedChunks, err := a.batchEmbedChunks(ctx, chunks)
	if err != nil {
		return fmt.Errorf("failed to get embeddings for chunks: %w", err)
	}

	return a.ragStore.SaveChunks(ctx, embeddedChunks)
}

func (a *Agent) batchEmbedChunks(ctx context.Context, chunks []golang.Chunk) ([]rag.Chunk, error) {
	enc, err := tiktoken.EncodingForModel("text-embedding-3-small")
	if err != nil {
		return nil, fmt.Errorf("invalid tiktoken encoding: %w", err)
	}

	// OpenAI's limit is 300k, but we're leaving some room for token inflation etc.
	maxTokens := 200_000
	var (
		batchContents  []string
		batchChunks    []golang.Chunk
		batchTokens    int
		embeddedChunks []rag.Chunk
	)

	flushBatch := func() error {
		if len(batchContents) == 0 {
			return nil
		}

		resp, err := a.client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
			Input: batchContents,
			Model: openai.SmallEmbedding3,
		})
		if err != nil {
			return fmt.Errorf("embedding request failed: %w", err)
		}

		for _, e := range resp.Data {
			sourceChunk := batchChunks[e.Index]
			tokenCount := len(enc.Encode(sourceChunk.Content, nil, nil))

			embeddedChunks = append(embeddedChunks, rag.Chunk{
				Chunk:  sourceChunk,
				Vector: e.Embedding,
				Tokens: tokenCount,
			})
		}

		// reset
		batchContents = batchContents[:0]
		batchChunks = batchChunks[:0]
		batchTokens = 0

		return nil
	}

	maxContentTokens := 8192
	for _, chunk := range chunks {
		tokens := len(enc.Encode(chunk.Content, nil, nil))

		// we can't create embeddings for chunk contents larger than 8192 tokens
		// this represents like small book worth of characters though so it's probably generated code or otherwise not very relevant
		if tokens > maxContentTokens {
			log.Printf("skipping chunk %s: %d tokens exceeds input token limit\n", chunk.Name, tokens)
			continue
		}

		batchContents = append(batchContents, chunk.Content)
		batchChunks = append(batchChunks, chunk)
		batchTokens += tokens

		// flush if tokens for batch is exceeded
		if batchTokens > maxTokens {
			if err := flushBatch(); err != nil {
				return nil, err
			}
		}
	}

	// Final flush
	if err := flushBatch(); err != nil {
		return nil, err
	}

	return embeddedChunks, nil
}
