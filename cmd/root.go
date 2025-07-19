package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"github.com/sajuno/goon/agent"
	"github.com/sajuno/goon/rag"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
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

			dsn := "postgresql://postgres:postgres@localhost:5432/goon?search_path=public,rag"
			pgCfg, err := pgxpool.ParseConfig(dsn)
			if err != nil {
				return err
			}
			pgCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
				return pgxvec.RegisterTypes(ctx, conn)
			}
			pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
			if err != nil {
				return err
			}

			if err = pool.Ping(ctx); err != nil {
				return fmt.Errorf("postgres not ready: %w", err)
			}

			ag = agent.New(openai.NewClient(cfg.APIKey), rag.NewPGStore(pool), agent.AssistantConfig{ID: cfg.AssistantID})
			return nil
		},
	}

	cmd.AddCommand(goonTest(ctx))
	cmd.AddCommand(goonIndex(ctx))

	return cmd
}
