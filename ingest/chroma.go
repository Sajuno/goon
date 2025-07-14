package ingest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Chroma struct {
	baseURL    string
	collection string

	http *http.Client
}

func NewChroma(baseURL string, collection string) *Chroma {
	return &Chroma{
		baseURL:    baseURL,
		collection: collection,
		http: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   10 * time.Second,
		},
	}
}

func (c *Chroma) EnsureCollection() error {
	reqBody := map[string]string{"name": c.collection}
	data, _ := json.Marshal(reqBody)

	resp, err := c.http.Post(c.baseURL+"/api/v1/collections", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create collection: %s", resp.Status)
	}
	return nil
}

func (c *Chroma) upsertChunks(chunks []EmbeddedChunk) error {
	var ids []string
	var docs []string
	var vectors [][]float64
	var metas []map[string]string

	for _, ch := range chunks {
		ids = append(ids, ch.ID)
		docs = append(docs, ch.Content)
		vectors = append(vectors, ch.Vector)
		metas = append(metas, map[string]string{
			"name":     ch.Name,
			"kind":     string(ch.Kind),
			"package":  ch.Package,
			"filePath": ch.FilePath,
		})
	}

	req := map[string]interface{}{
		"ids":        ids,
		"documents":  docs,
		"embeddings": vectors,
		"metadatas":  metas,
	}

	body, _ := json.Marshal(req)
	url := fmt.Sprintf("%s/api/v1/collections/%s/add", c.baseURL, c.collection)
	resp, err := c.http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upsert failed: %s", resp.Status)
	}
	return nil
}
