package agent

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
)

func (a *Agent) Explain(ctx context.Context, query string) (string, error) {
	resp, err := a.openai.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: query,
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create embeddings for user query: %w", err)
	}
	vec := resp.Data[0].Embedding

	simChunks, err := a.ragStore.FindSimilarChunks(ctx, vec)
	if err != nil {
		return "", fmt.Errorf("failed to find similar chunks: %w", err)
	}

	chunks := make([]rag.Chunk, 0, len(simChunks))
	for _, chunk := range simChunks {
		chunks = append(chunks, chunk.Chunk)
	}

	promptContext := buildPromptContext(chunks, 25000)
	prompt := fmt.Sprintf(`
%s

The user asked for an explanation about a codebase using the following query: "%s"

Attempt to generate an answer to that question as well as you can given the code context provided above.
This context represents a number of code chunks that were deemed relevant to the user's prompt.
Summarize the relevant behavior and responsibilities of the shown functions, types, and interactions. 
Focus on functionality, structure, and intent, not low-level implementation details.
Be concise, accurate, and if you can an asshole about it, please do so.
`, promptContext, query)

	response, err := a.promptAI(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM prompt failed: %w", err)
	}

	return response, nil
}
