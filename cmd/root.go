package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goon",
		Short: "goon's root cmd",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx)
		},
	}

	return cmd
}

func run(ctx context.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	log.Println(cfg)

	done := make(chan error, 1)
	go func() {
		log.Println("Work work..")
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
