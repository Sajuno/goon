package agent

import (
	"fmt"
	"github.com/sajuno/goon/rag"
	"strings"
)

func buildPromptContext(chunks []rag.Chunk, maxTokens int) string {
	var (
		sb         strings.Builder
		tokensUsed int
	)

	sb.WriteString("# Code Context\n\n")

	for _, chunk := range chunks {
		// generate code context up to max token limit
		if tokensUsed+chunk.Tokens >= maxTokens {
			break
		}

		sb.WriteString(fmt.Sprintf("## %s", chunk.Name))
		if chunk.FilePath != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", chunk.FilePath))
		}
		sb.WriteString("\n\n```go\n")
		sb.WriteString(chunk.Content)
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("\n```\n\n")

	return sb.String()
}
