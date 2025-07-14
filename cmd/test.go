package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func goonTest(ctx context.Context) *cobra.Command {
	var pkgName string

	cmd := &cobra.Command{
		Use:   "test <go-function-name>",
		Short: "Generate a test for a Go function",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			funcName := args[0]

			err := ag.WriteTest(ctx, funcName, pkgName)
			if err != nil {
				return fmt.Errorf("failed to generate goonTest: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&pkgName, "pkg", "", "Optional Go package name to narrow search scope")

	return cmd
}
