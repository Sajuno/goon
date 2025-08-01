// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: copyfrom.go

package pg

import (
	"context"
)

// iteratorForCreateChunks implements pgx.CopyFromSource.
type iteratorForCreateChunks struct {
	rows                 []CreateChunksParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateChunks) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateChunks) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].SymbolName,
		r.rows[0].SymbolType,
		r.rows[0].StartLine,
		r.rows[0].EndLine,
		r.rows[0].Content,
		r.rows[0].Doc,
		r.rows[0].ReceiverName,
		r.rows[0].Embedding,
		r.rows[0].TokenCount,
		r.rows[0].Sha256,
		r.rows[0].Package,
	}, nil
}

func (r iteratorForCreateChunks) Err() error {
	return nil
}

func (q *Queries) CreateChunks(ctx context.Context, arg []CreateChunksParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"code_chunks"}, []string{"symbol_name", "symbol_type", "start_line", "end_line", "content", "doc", "receiver_name", "embedding", "token_count", "sha256", "package"}, &iteratorForCreateChunks{rows: arg})
}
