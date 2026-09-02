package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type deleteResult struct {
	DeletedID string `json:"deleted_id"`
}

func registerDeletes(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_control_delete",
		Title:       "Delete an Openlane control",
		Description: "Permanently delete a control by ID. Requires delete mode.",
		Annotations: deleteAnnotations(),
	}, h.deleteControl)

	addTool(server, &mcp.Tool{
		Name:        "openlane_evidence_delete",
		Title:       "Delete Openlane evidence",
		Description: "Permanently delete an evidence record by ID. Requires delete mode.",
		Annotations: deleteAnnotations(),
	}, h.deleteEvidence)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_delete",
		Title:       "Delete an Openlane policy",
		Description: "Permanently delete an internal policy by ID. Requires delete mode.",
		Annotations: deleteAnnotations(),
	}, h.deletePolicy)

	addTool(server, &mcp.Tool{
		Name:        "openlane_risk_delete",
		Title:       "Delete an Openlane risk",
		Description: "Permanently delete a risk by ID. Requires delete mode.",
		Annotations: deleteAnnotations(),
	}, h.deleteRisk)

	addTool(server, &mcp.Tool{
		Name:        "openlane_task_delete",
		Title:       "Delete an Openlane task",
		Description: "Permanently delete a task by ID. Requires delete mode.",
		Annotations: deleteAnnotations(),
	}, h.deleteTask)
}

func (h *handlers) deleteControl(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, deleteResult, error) {
	return h.deleteByID(ctx, in.ID, h.api.DeleteControl)
}

func (h *handlers) deleteEvidence(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, deleteResult, error) {
	return h.deleteByID(ctx, in.ID, h.api.DeleteEvidence)
}

func (h *handlers) deletePolicy(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, deleteResult, error) {
	return h.deleteByID(ctx, in.ID, h.api.DeleteInternalPolicy)
}

func (h *handlers) deleteRisk(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, deleteResult, error) {
	return h.deleteByID(ctx, in.ID, h.api.DeleteRisk)
}

func (h *handlers) deleteTask(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, deleteResult, error) {
	return h.deleteByID(ctx, in.ID, h.api.DeleteTask)
}

type deleteFn func(context.Context, string) (string, error)

func (h *handlers) deleteByID(ctx context.Context, id string, fn deleteFn) (*mcp.CallToolResult, deleteResult, error) {
	if id == "" {
		return nil, deleteResult{}, errIDRequired
	}
	deletedID, err := fn(ctx, id)
	if err != nil {
		return nil, deleteResult{}, openlane.APIError(err)
	}
	return nil, deleteResult{DeletedID: deletedID}, nil
}
