package agent

import (
	"context"
	"github.com/sajuno/goon/ingest"
	"log"
)

type AssistantConfig struct {
	ID string
}
type Agent struct {
	cfg AssistantConfig
}

func New(cfg AssistantConfig) *Agent {
	return &Agent{cfg: cfg}
}

func (a *Agent) WriteTest(ctx context.Context, funcName string) error {
	fn, err := ingest.FindFunction(".", funcName)
	if err != nil {
		return err
	}
	log.Println(fn.Source)
	return nil
}
