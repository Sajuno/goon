package repl

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/agent"
	"strings"
)

type commandHandler struct {
	agent *agent.Agent
}

func newCommandHandler(agent *agent.Agent) *commandHandler {
	return &commandHandler{agent: agent}
}

func (h *commandHandler) handleCommand(ctx context.Context, line string) error {
	switch {
	case strings.HasPrefix(line, "explain "):
		prompt := strings.TrimPrefix(line, "explain ")
		fmt.Print("Thinking...")
		return h.explain(ctx, prompt)
	default:
		fmt.Println("Unknown command. Try :help")
	}
	return nil
}

func (h *commandHandler) explain(ctx context.Context, prompt string) error {
	response, err := h.agent.Explain(ctx, prompt)
	fmt.Print("\r\033[2K") // clear 'spinner'
	if err != nil {
		return fmt.Errorf(`failed to explain "%s": %w`, prompt, err)
	}

	fmt.Println(response)
	return nil
}
