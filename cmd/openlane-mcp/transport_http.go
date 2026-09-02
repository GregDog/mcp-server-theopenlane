package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
)

const (
	httpReadHeaderTimeout = 10 * time.Second
	httpReadTimeout       = 2 * time.Minute
	httpWriteTimeout      = 2 * time.Minute
	httpIdleTimeout       = 60 * time.Second
	httpShutdownTimeout   = 10 * time.Second
)

func runHTTP(ctx context.Context, cfg config.Config, server *mcp.Server) error {
	warnIfNonLoopbackHTTPBind(cfg.HTTPAddr)

	opts := &mcp.StreamableHTTPOptions{
		Stateless:           true,
		JSONResponse:        cfg.HTTPJSON,
		MaxRequestBodyBytes: cfg.HTTPMaxBodyBytes,
	}
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, opts)

	httpServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}

func warnIfNonLoopbackHTTPBind(addr string) {
	if config.IsLoopbackHTTPAddr(addr) {
		return
	}
	slog.Warn("http bind address is not loopback; built-in HTTP authentication is not enabled",
		"addr", addr,
		"guidance", "do not expose this server directly to the public internet; place a trusted authentication reverse proxy in front",
	)
}
