-- name: CreateChunk :one
INSERT INTO code_chunks (symbol_name, symbol_type, start_line, end_line, content, doc, embedding, token_count, sha256, package, file_path)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: CreateChunks :copyfrom
INSERT INTO code_chunks (
    symbol_name,
    symbol_type,
    start_line,
    end_line,
    content,
    doc,
    embedding,
    token_count,
    sha256,
    package,
    file_path
) VALUES (
    @symbol_name,
    @symbol_type,
    @start_line,
    @end_line,
    @content,
    @doc,
    @embedding,
    @token_count,
    @sha256,
    @package,
    @file_path
);

-- name: FindSimilarChunks :many
SELECT *,
       embedding <-> $1 AS distance
FROM code_chunks
ORDER BY embedding <-> $1
LIMIT $2;
