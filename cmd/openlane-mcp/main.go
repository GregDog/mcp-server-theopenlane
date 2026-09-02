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
		flags, err := parseServeArgs(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			usage()
			os.Exit(2)
		}
		if err := runServe(flags); err != nil {
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
  openlane-mcp serve [--allow-write] [--allow-delete]
  openlane-mcp version

Environment:
  OPENLANE_API_TOKEN          Required Openlane API token or PAT
  OPENLANE_BASE_URL           Openlane API base URL (default https://api.theopenlane.io)
  OPENLANE_ORGANIZATION_ID    Organization ID when using a multi-org PAT
  OPENLANE_MCP_LOG_LEVEL      debug, info, warn, or error (default info)
  OPENLANE_ALLOW_WRITE        Enable write MCP tools when true (default false)
  OPENLANE_ALLOW_DELETE       Enable delete MCP tools when true (default false)
`)
}

type serveFlags struct {
	allowWrite  bool
	allowDelete bool
}

func parseServeArgs(args []string) (serveFlags, error) {
	flags := serveFlags{}
	for _, arg := range args {
		switch arg {
		case "--allow-write":
			flags.allowWrite = true
		case "--allow-delete":
			flags.allowDelete = true
		default:
			return serveFlags{}, fmt.Errorf("unknown serve flag %q", arg)
		}
	}
	return flags, nil
}

func runServe(flags serveFlags) error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}
	if flags.allowWrite {
		cfg.AllowWrite = true
	}
	if flags.allowDelete {
		cfg.AllowDelete = true
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
	tools.Register(server, client, tools.Options{
		AllowWrite:  cfg.AllowWrite,
		AllowDelete: cfg.AllowDelete,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("starting openlane mcp server",
		"transport", "stdio",
		"version", version,
		"allow_write", cfg.AllowWrite,
		"allow_delete", cfg.AllowDelete,
	)
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
