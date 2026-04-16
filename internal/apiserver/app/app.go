package app

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"

	"mube/internal/apiserver/config"
	apiserverhttp "mube/internal/apiserver/http"
	"mube/internal/apiserver/http/handlers"
	"mube/internal/apiserver/store"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load api server config: %w", err)
	}

	registered := make([]store.RegisteredNode, 0, len(cfg.Nodes))
	for _, n := range cfg.Nodes {
		registered = append(registered, store.RegisteredNode{
			Name: n.Name,
			IP:   n.IP,
			Port: n.Port,
		})
	}

	nodeStore := store.NewMemoryNodeStore(registered, cfg.NodeNotReadyTimeout, time.Now())
	nodeHandler := handlers.NewNodeHandler(nodeStore)
	router := apiserverhttp.NewRouter(nodeHandler)

	srv := &nethttp.Server{
		Addr:              cfg.ListenAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			errCh <- fmt.Errorf("listen and serve api server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown api server: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}
