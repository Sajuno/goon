package ingest

import (
	"context"
	"log"
)

type RAGStore struct {
	chroma *Chroma
}

func NewRAGStore(chroma *Chroma) *RAGStore {
	return &RAGStore{chroma: chroma}
}

// Everything ingests every go file and node type, and converts it to a Chunk
func (s *RAGStore) Everything(ctx context.Context, path string) error {
	files, err := findGoFiles(path)
	if err != nil {
		return err
	}

	var all []Chunk
	for _, f := range files {
		chunks, err := chunkFile(f)
		if err != nil {
			log.Printf("❌ Error parsing %s: %v", f, err)
			continue
		}
		all = append(all, chunks...)
	}

	log.Printf("✅ Ingested %d chunks from %d files", len(all), len(files))
	for _, chunk := range all {
		log.Printf("%+v\n", chunk)
	}

	return nil
}
