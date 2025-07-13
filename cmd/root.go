package cmd

import (
	"context"
	"fmt"
	"github.com/sajuno/goon/agent"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// global agent instance
var ag *agent.Agent

func NewRootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goon",
		Short: "goon's root cmd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}
			ag = agent.New(agent.AssistantConfig{ID: cfg.AssistantID})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx)
		},
	}

	cmd.AddCommand(goonTest(ctx))

	return cmd
}

func run(ctx context.Context) error {
	done := make(chan error, 1)
	go func() {
		log.Println("WriteTest work..")
		defer log.Println("Job's done")

		select {
		case <-time.After(1 * time.Second):
			done <- nil
		case <-ctx.Done():
			done <- ctx.Err()
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Received shutdown signal")
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return fmt.Errorf("job failed: %v", err)
		}
		log.Println("Job completed successfully")
		return nil
	}
}
