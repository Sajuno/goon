package agent

import (
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
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
