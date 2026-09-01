package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type listInput struct {
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of items to return. Default 20, maximum 50."`
	Cursor string `json:"cursor,omitempty" jsonschema:"Opaque cursor from a previous list response."`
}

type getInput struct {
	ID string `json:"id" jsonschema:"Openlane object ID."`
}

type searchInput struct {
	Query  string `json:"query" jsonschema:"Text to match against control ref code, title, or description."`
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of items to return. Default 20, maximum 50."`
	Cursor string `json:"cursor,omitempty" jsonschema:"Opaque cursor from a previous search response."`
}

func readOnly() *mcp.ToolAnnotations {
	destructive := false
	openWorld := true
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
		IdempotentHint:  true,
	}
}

func pageArgs(limit int, cursor string) (first int64, after *string) {
	return openlane.ClampLimit(limit), openlane.CursorPtr(cursor)
}
