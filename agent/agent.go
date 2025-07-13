package agent

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/ingest"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

type AssistantConfig struct {
	ID string
}
type Agent struct {
	cfg    AssistantConfig
	client *openai.Client
}

func New(client *openai.Client, cfg AssistantConfig) *Agent {
	return &Agent{cfg: cfg, client: client}
}

func (a *Agent) WriteTest(ctx context.Context, funcName, pkgName string) error {
	fn, err := ingest.FindFunction(".", ingest.FindFunctionQuery{
		Name:    funcName,
		Package: pkgName,
	})
	if err != nil {
		return err
	}

	thread, err := a.client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return fmt.Errorf("failed to create new thread: %w", err)
	}

	_, err = a.client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role: openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(`
Generate a test for the following function. It lives inside of the '%s' package.

**Return only the Go test code, and nothing else. No explanations. No comments. No Markdown. No formatting.**

%s
`, fn.Package, fn.Source),
	})
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	run, err := a.client.CreateRun(ctx, thread.ID, openai.RunRequest{
		AssistantID: a.cfg.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create new run: %w", err)
	}

	// TODO: add timeout
	for {
		time.Sleep(time.Second)
		runStatus, _ := a.client.RetrieveRun(ctx, thread.ID, run.ID)

		switch runStatus.Status {
		case openai.RunStatusCompleted:
			break
		case openai.RunStatusFailed:
			return fmt.Errorf("run failed: %s", runStatus.Status)
		default:
			log.Println("Run still in progress, waiting...")
		}
		if runStatus.Status == openai.RunStatusCompleted {
			break
		}
	}

	res, err := a.client.ListMessage(ctx, thread.ID, nil, nil, nil, nil, nil)
	if err != nil {
		log.Fatalf("Error listing message: %v", err)
	}

	testCode := res.Messages[0].Content[0].Text.Value
	log.Println("GENERATED:")
	log.Println(testCode)

	return nil
}
