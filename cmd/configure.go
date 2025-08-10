package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

func configure(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Sets up Open AI for usage with Goon",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Configure(ctx)
		},
	}

	return cmd
}
