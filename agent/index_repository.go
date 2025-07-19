package agent

import (
	"context"
	"github.com/sajuno/goon/golang"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
)

func (a *Agent) IndexRepository(ctx context.Context, path string) error {
	chunks, err := golang.ChunkRepository(path)
	if err != nil {
		return err
	}

	var embeddedChunks []rag.EmbeddedChunk
	for _, chunk := range chunks {
		resp, err := a.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Input: chunk.Content,
			Model: openai.SmallEmbedding3, // This is `text-embedding-3-small`
		})
		if err != nil {
			return err
		}
		embedding := resp.Data[0].Embedding
		embeddedChunks = append(embeddedChunks, rag.EmbeddedChunk{
			Chunk:  chunk,
			Vector: embedding,
		})
	}

	err = a.ragStore.SaveChunks(ctx, embeddedChunks)
	if err != nil {
		return err
	}

	return nil
}
