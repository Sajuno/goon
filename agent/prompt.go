package agent

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
	"strings"
	"time"
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

		tokensUsed += chunk.Tokens
	}
	sb.WriteString("\n```\n\n")

	return sb.String()
}

func (a *Agent) promptAI(ctx context.Context, prompt string) (string, error) {
	thread, err := a.openai.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create new thread: %w", err)
	}

	_, err = a.openai.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new message: %w", err)
	}

	run, err := a.openai.CreateRun(ctx, thread.ID, openai.RunRequest{
		AssistantID: a.cfg.ID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new run: %w", err)
	}

	// TODO: add timeout
	for {
		time.Sleep(time.Second)
		runStatus, _ := a.openai.RetrieveRun(ctx, thread.ID, run.ID)

		switch runStatus.Status {
		case openai.RunStatusCompleted:
			break
		case openai.RunStatusFailed:
			return "", fmt.Errorf("run failed, last error: %s", runStatus.LastError.Message)
		default:
		}
		if runStatus.Status == openai.RunStatusCompleted {
			break
		}
	}

	res, err := a.openai.ListMessage(ctx, thread.ID, nil, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list message: %w", err)
	}

	return res.Messages[0].Content[0].Text.Value, nil
}
