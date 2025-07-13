package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sajuno/goon/cmd"
	"log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := cmd.NewRootCmd(ctx).Execute(); err != nil {
		log.Printf("Command error: %v", err)
		os.Exit(1)
	}
}
