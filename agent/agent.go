package agent

import (
	"github.com/sajuno/goon/language/lsp"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
)

type AssistantConfig struct {
	ID string
}
type Agent struct {
	cfg AssistantConfig

	openai   *openai.Client
	ragStore rag.Store
	lsp      *lsp.Client
}

func New(openai *openai.Client, ragStore rag.Store, cfg AssistantConfig, lsp *lsp.Client) *Agent {
	return &Agent{cfg: cfg, openai: openai, ragStore: ragStore, lsp: lsp}
}
