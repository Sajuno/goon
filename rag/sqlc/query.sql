-- name: CreateChunk :one
INSERT INTO code_chunks (symbol_name, symbol_type, file_path, start_line, end_line, content, doc, embedding, token_count, sha256)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
    sha256
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
    @sha256
);
