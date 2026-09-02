package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// addTool registers a tool handler without an output schema in tools/list.
// Cursor silently drops tool lists when the JSON payload is too large; omitting
// per-tool output schemas keeps the list small while preserving typed handlers.
func addTool[In, Out any](server *mcp.Server, t *mcp.Tool, h func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)) {
	mcp.AddTool(server, t, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, any, error) {
		return h(ctx, req, in)
	})
}
