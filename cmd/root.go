package cmd

import (
	"context"
	"github.com/sajuno/goon/agent"
	"github.com/sajuno/goon/ingest"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// global agent instance
var ag *agent.Agent
var ragStore *ingest.RAGStore

func NewRootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goon",
		Short: "goon's root cmd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			ag = agent.New(openai.NewClient(cfg.APIKey), agent.AssistantConfig{ID: cfg.AssistantID})
			chroma := ingest.NewChroma("http://localhost:8000", "chunks")
			// TODO: pass context
			err := chroma.EnsureCollection()
			if err != nil {
				return err
			}
			ragStore = ingest.NewRAGStore(chroma)
			return nil
		},
	}

	cmd.AddCommand(goonTest(ctx))
	cmd.AddCommand(goonIngest(ctx))

	return cmd
}
