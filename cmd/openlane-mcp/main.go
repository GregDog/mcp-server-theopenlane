package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
	"github.com/GregDog/mcp-server-theopenlane/internal/tools"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "serve":
		if err := runServe(); err != nil {
			fmt.Fprintln(os.Stderr, openlane.Redact(err.Error()))
			os.Exit(1)
		}
	case "version", "--version", "-v":
		fmt.Printf("openlane-mcp %s (commit %s, built %s)\n", version, commit, date)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `openlane-mcp is a Model Context Protocol server for the Openlane GRC platform.

Usage:
  openlane-mcp serve
  openlane-mcp version

Environment:
  OPENLANE_API_TOKEN          Required Openlane API token or PAT
  OPENLANE_BASE_URL           Openlane API base URL (default https://api.theopenlane.io)
  OPENLANE_ORGANIZATION_ID    Organization ID when using a multi-org PAT
  OPENLANE_MCP_LOG_LEVEL      debug, info, warn, or error (default info)
`)
}

func runServe() error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	setupLogger(cfg.LogLevel)

	client, err := openlane.New(cfg)
	if err != nil {
		return err
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "openlane-mcp",
		Title:   "Openlane MCP Server",
		Version: version,
	}, nil)
	tools.Register(server, client)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting openlane mcp server", "transport", "stdio", "version", version)
	return server.Run(ctx, &mcp.StdioTransport{})
}

func setupLogger(level string) {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})))
}
