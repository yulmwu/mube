package agent

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"
	"sync"
	"time"

	"mube/internal/mubelet/client"
	"mube/internal/mubelet/config"
	mubelethttp "mube/internal/mubelet/http"
	v1 "mube/pkg/api/v1"
)

type Agent struct {
	cfg    config.Config
	client *client.APIServerClient
}

func New(cfg config.Config) *Agent {
	return &Agent{
		cfg:    cfg,
		client: client.NewAPIServerClient(cfg.APIServerEndpoint),
	}
}

func (a *Agent) Run(ctx context.Context) error {
	router := mubelethttp.NewRouter(a.cfg)
	srv := &nethttp.Server{
		Addr:              a.cfg.HealthListen,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			errCh <- fmt.Errorf("run mubelet http server: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(a.cfg.HeartbeatInterval)
		defer ticker.Stop()

		if err := a.sendHeartbeat(ctx); err != nil {
			log.Printf("initial heartbeat failed: %v", err)
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := a.sendHeartbeat(ctx); err != nil {
					log.Printf("heartbeat failed: %v", err)
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown mubelet http server: %w", err)
		}
		wg.Wait()
		return nil
	case err := <-errCh:
		return err
	}
}

func (a *Agent) sendHeartbeat(ctx context.Context) error {
	req := v1.NodeHeartbeatRequest{
		Name:     a.cfg.NodeName,
		Runtime:  a.cfg.Runtime,
		Version:  a.cfg.Version,
		Capacity: a.cfg.Capacity,
	}

	return a.client.SendHeartbeat(ctx, req)
}
