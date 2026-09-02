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

type serveFlags struct {
	allowWrite  bool
	allowDelete bool
	transport   string
	httpAddr    string
	httpJSON    bool
}

func parseServeArgs(args []string) (serveFlags, error) {
	flags := serveFlags{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--allow-write":
			flags.allowWrite = true
		case "--allow-delete":
			flags.allowDelete = true
		case "--http-json":
			flags.httpJSON = true
		case "--transport", "--http-addr":
			if i+1 >= len(args) {
				return serveFlags{}, fmt.Errorf("flag %s requires a value", arg)
			}
			i++
			if arg == "--transport" {
				flags.transport = args[i]
			} else {
				flags.httpAddr = args[i]
			}
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
	if flags.transport != "" {
		cfg.Transport = flags.transport
	}
	if flags.httpAddr != "" {
		cfg.HTTPAddr = flags.httpAddr
	}
	if flags.httpJSON {
		cfg.HTTPJSON = true
	}

	setupLogger(cfg.LogLevel)

	client, err := openlane.New(cfg)
	if err != nil {
		return err
	}

	server := newMCPServer(cfg, client)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	transport := strings.ToLower(cfg.Transport)
	slog.Info("starting openlane mcp server",
		"transport", transport,
		"version", version,
		"allow_write", cfg.AllowWrite,
		"allow_delete", cfg.AllowDelete,
	)
	switch transport {
	case "stdio":
		return server.Run(ctx, &mcp.StdioTransport{})
	case "http":
		slog.Info("http transport configured",
			"addr", cfg.HTTPAddr,
			"json_response", cfg.HTTPJSON,
			"max_body_bytes", cfg.HTTPMaxBodyBytes,
		)
		return runHTTP(ctx, cfg, server)
	default:
		return fmt.Errorf("unsupported transport %q", cfg.Transport)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `openlane-mcp is a Model Context Protocol server for the Openlane GRC platform.

Usage:
  openlane-mcp serve [--allow-write] [--allow-delete] [--transport stdio|http] [--http-addr ADDR] [--http-json]
  openlane-mcp version

Environment:
  OPENLANE_API_TOKEN              Required Openlane API token or PAT
  OPENLANE_BASE_URL               Openlane API base URL (default https://api.theopenlane.io)
  OPENLANE_ORGANIZATION_ID        Organization ID when using a multi-org PAT
  OPENLANE_MCP_LOG_LEVEL          debug, info, warn, or error (default info)
  OPENLANE_ALLOW_WRITE            Enable write MCP tools when true (default false)
  OPENLANE_ALLOW_DELETE           Enable delete MCP tools when true (default false)
  OPENLANE_MCP_TRANSPORT          stdio or http (default stdio)
  OPENLANE_MCP_HTTP_ADDR          HTTP listen address when transport=http (default :8080)
  OPENLANE_MCP_HTTP_JSON          Use application/json responses for HTTP transport
  OPENLANE_MCP_HTTP_MAX_BODY_BYTES  Max HTTP request body size (default 33554432)
  OPENLANE_MCP_MAX_UPLOAD_BYTES   Max decoded evidence upload size (default 10485760)
  OPENLANE_MCP_UPLOAD_TIMEOUT     Openlane API timeout for uploads (default 2m)
`)
}
