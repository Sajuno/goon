package agent

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

type AssistantConfig struct {
	ID string
}
type Agent struct {
	cfg AssistantConfig

	client   *openai.Client
	ragStore rag.Store
}

func New(client *openai.Client, ragStore rag.Store, cfg AssistantConfig) *Agent {
	return &Agent{cfg: cfg, client: client, ragStore: ragStore}
}

func (a *Agent) promptAI(ctx context.Context, prompt string) (string, error) {
	thread, err := a.client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create new thread: %w", err)
	}

	_, err = a.client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new message: %w", err)
	}

	run, err := a.client.CreateRun(ctx, thread.ID, openai.RunRequest{
		AssistantID: a.cfg.ID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create new run: %w", err)
	}

	// TODO: add timeout
	for {
		time.Sleep(time.Second)
		runStatus, _ := a.client.RetrieveRun(ctx, thread.ID, run.ID)

		switch runStatus.Status {
		case openai.RunStatusCompleted:
			break
		case openai.RunStatusFailed:
			return "", fmt.Errorf("run failed: %s", runStatus.Status)
		default:
			log.Println("Run still in progress, waiting...")
		}
		if runStatus.Status == openai.RunStatusCompleted {
			break
		}
	}

	res, err := a.client.ListMessage(ctx, thread.ID, nil, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list message: %w", err)
	}

	return res.Messages[0].Content[0].Text.Value, nil
}
