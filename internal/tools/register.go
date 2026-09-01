package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

// Register adds the verified read-only Openlane tools to the MCP server.
func Register(server *mcp.Server, api openlane.GraphAPI) {
	h := &handlers{api: api}
	registerControls(server, h)
	registerPrograms(server, h)
	registerEvidence(server, h)
	registerPolicies(server, h)
	registerRisks(server, h)
	registerStandards(server, h)
}

type handlers struct {
	api openlane.GraphAPI
}
