package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"mube/internal/mubelet/agent"
	"mube/internal/mubelet/config"
)

func main() {
	cfg := config.Load()
	a := agent.New(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		log.Fatalf("mubelet failed: %v", err)
	}
}
