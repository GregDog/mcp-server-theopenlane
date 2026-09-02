package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
	"github.com/GregDog/mcp-server-theopenlane/internal/tools"
)

func newMCPServer(cfg config.Config, client openlane.GraphAPI) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "openlane-mcp",
		Title:   "Openlane MCP Server",
		Version: version,
	}, nil)
	tools.Register(server, client, tools.Options{
		AllowWrite:     cfg.AllowWrite,
		AllowDelete:    cfg.AllowDelete,
		MaxUploadBytes: cfg.MaxUploadBytes,
	})
	return server
}
