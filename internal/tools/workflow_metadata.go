package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type workflowMetadataResult struct {
	ObjectTypes []openlane.WorkflowObjectTypeMetadata `json:"object_types"`
	Extensions  map[string]any                        `json:"extensions,omitempty"`
	Summary     string                                `json:"summary,omitempty"`
}

func registerWorkflowMetadata(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_metadata_get",
		Title:       "Get Openlane workflow metadata",
		Description: "Return workflow-eligible fields, edges, and resolver keys per object type from the live workflowMetadata query. Use before authoring or explaining approval-gated fields.",
		Annotations: readOnly(),
	}, h.getWorkflowMetadata)
}

func (h *handlers) getWorkflowMetadata(ctx context.Context, _ *mcp.CallToolRequest, in workflowMetadataInput) (*mcp.CallToolResult, workflowMetadataResult, error) {
	meta, err := h.api.GetWorkflowMetadata(ctx)
	if err != nil {
		return nil, workflowMetadataResult{}, openlane.APIError(err)
	}
	types := meta.ObjectTypes
	if s := strings.TrimSpace(in.SchemaType); s != "" {
		types = filterWorkflowObjectTypes(meta.ObjectTypes, s)
	}
	result := workflowMetadataResult{
		ObjectTypes: types,
		Extensions:  meta.Extensions,
		Summary:     summarizeWorkflowMetadata(types),
	}
	return nil, result, nil
}

func filterWorkflowObjectTypes(types []openlane.WorkflowObjectTypeMetadata, schemaType string) []openlane.WorkflowObjectTypeMetadata {
	norm := normalizeWorkflowObjectType(schemaType)
	var out []openlane.WorkflowObjectTypeMetadata
	for _, t := range types {
		if strings.EqualFold(t.Type, norm) || strings.EqualFold(t.Type, schemaType) {
			out = append(out, t)
		}
	}
	return out
}

func summarizeWorkflowMetadata(types []openlane.WorkflowObjectTypeMetadata) string {
	if len(types) == 0 {
		return "No workflow object types matched."
	}
	var parts []string
	for _, t := range types {
		parts = append(parts, fmt.Sprintf("%s (%d eligible fields, %d resolvers)", t.Type, len(t.EligibleFields), len(t.ResolverKeys)))
	}
	return "Workflow metadata: " + strings.Join(parts, "; ") + "."
}
