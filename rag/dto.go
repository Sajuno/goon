package rag

import (
	"github.com/sajuno/goon/golang"
	"github.com/sajuno/goon/rag/sqlc/pg"
)

func unmarshalChunk(chunk pg.CodeChunk) Chunk {
	return Chunk{
		Chunk: golang.Chunk{
			ID:           chunk.ID.String(),
			Content:      chunk.Content,
			FilePath:     chunk.FilePath,
			Package:      chunk.Package,
			Kind:         golang.ChunkKind(chunk.SymbolType),
			Name:         chunk.SymbolName,
			StartLine:    int(chunk.StartLine),
			EndLine:      int(chunk.EndLine),
			Doc:          chunk.Doc.String,
			ReceiverName: chunk.ReceiverName.String,
		},
		Vector: chunk.Embedding.Slice(),
		Tokens: int(chunk.TokenCount),
	}
}

func unmarshalSimilarChunks(chunks []pg.FindSimilarChunksRow) []SimilarChunk {
	out := make([]SimilarChunk, 0, len(chunks))
	for _, chunk := range chunks {
		out = append(out, SimilarChunk{
			Chunk: Chunk{
				Chunk: golang.Chunk{
					ID:           chunk.ID.String(),
					Content:      chunk.Content,
					FilePath:     chunk.FilePath,
					Package:      chunk.Package,
					Kind:         golang.ChunkKind(chunk.SymbolType),
					Name:         chunk.SymbolName,
					StartLine:    int(chunk.StartLine),
					EndLine:      int(chunk.EndLine),
					Doc:          chunk.Doc.String,
					ReceiverName: chunk.ReceiverName.String,
				},
				Vector: chunk.Embedding.Slice(),
				Tokens: int(chunk.TokenCount),
			},
			Distance: chunk.Distance.(float64),
		})
	}
	return out
}
