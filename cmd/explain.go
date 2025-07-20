package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func goonExplain(ctx context.Context) *cobra.Command {
	var pkgName string

	cmd := &cobra.Command{
		Use:   "explain <query>",
		Short: "explain will attempt to answer a question about the code base",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := strings.Join(args, " ")
			response, err := ag.Explain(ctx, prompt)
			if err != nil {
				return fmt.Errorf("failed to generate goonExplain: %w", err)
			}

			fmt.Println(response)

			return nil
		},
	}

	cmd.Flags().StringVar(&pkgName, "pkg", "", "Optional Go package name to narrow search scope")

	return cmd
}
