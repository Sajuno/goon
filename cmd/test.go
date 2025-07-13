package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func goonTest(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "test <go-function-name>",
		Short: "Generate a test for a Go function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			funcName := args[0]
			err := ag.WriteTest(ctx, funcName)
			if err != nil {
				return fmt.Errorf("failed to generate test: %w", err)
			}

			return nil
		},
	}
}
