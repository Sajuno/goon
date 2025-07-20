CREATE SCHEMA IF NOT EXISTS rag;
SET SEARCH_PATH = 'rag', 'public';

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS code_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol_name TEXT NOT NULL,  -- the name as which the symbol was declared
    symbol_type TEXT NOT NULL,
    file_path TEXT NOT NULL,
    package TEXT NOT NULL,
    start_line INT NOT NULL,
    end_line INT NOT NULL,
    content TEXT NOT NULL,      -- Code itself as raw text
    doc TEXT,                   -- Optional comments
    receiver_name TEXT,         -- Optional receiver for methods

    -- 1536 represents the vector dimensions of GPT's text-embedding-3-small
    embedding vector(1536),
    token_count INT NOT NULL,
    sha256 TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
