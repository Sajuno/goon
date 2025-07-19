package rag

import (
	"context"
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

func (s *PGStore) SaveChunks(ctx context.Context, chunks []EmbeddedChunk) error {
	var params []pg.CreateChunksParams
	for _, chunk := range chunks {
		params = append(params, pg.CreateChunksParams{
			SymbolName:   chunk.Name,
			SymbolType:   chunk.Kind.String(),
			FilePath:     chunk.FilePath,
			EndLine:      int32(chunk.EndLine),
			Content:      chunk.Content,
			Doc:          pgtype.Text{String: chunk.Doc},
			ReceiverName: pgtype.Text{String: chunk.ReceiverName},
			Embedding:    pgvector.NewVector(chunk.Vector),
			TokenCount:   int32(chunk.Tokens),
			Sha256:       chunk.Sha256(),
		})
	}

	_, err := s.queries.CreateChunks(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
