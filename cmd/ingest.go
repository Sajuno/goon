package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func goonIngest(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest <(root)path>",
		Short: "ingest all go code recursively to make your goon context aware",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			err := ragStore.Everything(ctx, path)
			if err != nil {
				return fmt.Errorf("failed to generate goonTest: %w", err)
			}

			return nil
		},
	}

	return cmd
}
