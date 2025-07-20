-- name: CreateChunk :one
INSERT INTO code_chunks (symbol_name, symbol_type, file_path, start_line, end_line, content, doc, embedding, token_count, sha256, package)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: CreateChunks :copyfrom
INSERT INTO code_chunks (
    symbol_name,
    symbol_type,
    file_path,
    start_line,
    end_line,
    content,
    doc,
    receiver_name,
    embedding,
    token_count,
    sha256,
    package
) VALUES (
    @symbol_name,
    @symbol_type,
    @file_path,
    @start_line,
    @end_line,
    @content,
    @doc,
    @receiver_name,
    @embedding,
    @token_count,
    @sha256,
    @package
);

-- name: FindSimilarChunks :many
SELECT *,
       embedding <-> $1 AS distance
FROM code_chunks
ORDER BY embedding <-> $1
LIMIT $2;
