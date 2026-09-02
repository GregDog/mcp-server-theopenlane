package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

// Register adds Openlane MCP tools to the server.
func Register(server *mcp.Server, api openlane.GraphAPI, opts Options) {
	h := &handlers{api: api, allowWrite: opts.AllowWrite, maxUploadBytes: opts.MaxUploadBytes}
	registerControls(server, h)
	registerPrograms(server, h)
	registerEvidence(server, h)
	registerPolicies(server, h)
	registerRisks(server, h)
	registerStandards(server, h)
	registerTasks(server, h)
	registerEntities(server, h)
	registerAssets(server, h)
	registerContacts(server, h)
	if opts.AllowWrite {
		registerWriteControls(server, h)
		registerWriteEvidence(server, h)
		registerWritePolicies(server, h)
		registerWriteRisks(server, h)
		registerWriteTasks(server, h)
	}
	if opts.AllowDelete {
		registerDeletes(server, h)
	}
}

type handlers struct {
	api            openlane.GraphAPI
	allowWrite     bool
	maxUploadBytes int64
}
