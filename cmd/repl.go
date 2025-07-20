package cmd

import (
	"context"
	"github.com/sajuno/goon/repl"
	"github.com/spf13/cobra"
)

func goonRepl(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repl",
		Short: "Starts Goon's repl for interactive assistance",
		RunE: func(cmd *cobra.Command, args []string) error {
			repl.Start(ctx, ag)
			return nil
		},
	}

	return cmd
}
