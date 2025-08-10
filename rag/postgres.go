package rag

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"github.com/sajuno/goon/rag/sqlc/pg"
)

type PGStore struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (s *PGStore) SaveChunks(ctx context.Context, chunks []Chunk) error {
	var params []pg.CreateChunksParams
	for _, chunk := range chunks {
		params = append(params, pg.CreateChunksParams{
			SymbolName: chunk.Name,
			SymbolType: chunk.Kind.String(),
			Package:    chunk.Package,
			FilePath:   chunk.FilePath,
			EndLine:    int32(chunk.EndLine),
			Content:    chunk.Content,
			Doc:        pgtype.Text{String: chunk.Doc},
			Embedding:  pgvector.NewVector(chunk.Vector),
			TokenCount: int32(chunk.Tokens),
			Sha256:     chunk.Sha256(),
		})
	}

	_, err := s.queries.CreateChunks(ctx, params)
	if err != nil {
		return err
	}

	q := `
DROP INDEX IF EXISTS code_chunks_embedding_idx;
CREATE INDEX code_chunks_embedding_idx 
ON code_chunks 
USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);`

	_, err = s.pool.Exec(ctx, q)
	if err != nil {
		return fmt.Errorf("failed to recreate embedding index: %w", err)
	}

	return nil
}

func (s *PGStore) FindSimilarChunks(ctx context.Context, vector []float32) ([]SimilarChunk, error) {
	res, err := s.queries.FindSimilarChunks(ctx, pg.FindSimilarChunksParams{
		Embedding: pgvector.NewVector(vector),
		Limit:     50,
	})
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return unmarshalSimilarChunks(res), nil
}
