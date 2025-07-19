package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

func goonIndex(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index <(root)path>",
		Short: "indexes all go code from path recursively to make your goon context aware",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			err := ag.IndexRepository(ctx, path)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
