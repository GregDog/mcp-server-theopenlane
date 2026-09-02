package tools

import (
	"context"
	"fmt"

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

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_delete",
		Title:       "Delete an Openlane workflow definition",
		Description: "Permanently delete a workflow definition by ID. Requires delete mode and confirm: true.",
		Annotations: deleteAnnotations(),
	}, h.deleteWorkflow)
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

type deleteWorkflowInput struct {
	ID      string `json:"id" jsonschema:"Workflow definition ID to delete."`
	Confirm bool   `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type deleteWorkflowResult struct {
	Confirmed bool   `json:"confirmed"`
	Error     string `json:"error,omitempty"`
	DeletedID string `json:"deleted_id,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Summary   string `json:"summary,omitempty"`
}

func (h *handlers) deleteWorkflow(ctx context.Context, _ *mcp.CallToolRequest, in deleteWorkflowInput) (*mcp.CallToolResult, deleteWorkflowResult, error) {
	if in.ID == "" {
		return nil, deleteWorkflowResult{}, errIDRequired
	}
	existing, err := h.api.GetWorkflowDefinitionByID(ctx, in.ID)
	if err != nil {
		return nil, deleteWorkflowResult{}, openlane.APIError(err)
	}
	item := mapWorkflowDefinition(existing.WorkflowDefinition)
	out := deleteWorkflowResult{ID: item.ID, Name: item.Name, Summary: "Would delete workflow " + item.Name}
	if !in.Confirm {
		out.Error = errConfirmationRequired
		return nil, out, nil
	}
	deletedID, err := h.api.DeleteWorkflowDefinition(ctx, in.ID)
	if err != nil {
		return nil, deleteWorkflowResult{}, openlane.APIError(err)
	}
	if _, err := h.api.GetWorkflowDefinitionByID(ctx, in.ID); err == nil {
		return nil, deleteWorkflowResult{}, fmt.Errorf("workflow delete reported success but definition still exists")
	}
	out.Confirmed = true
	out.DeletedID = deletedID
	out.Summary = "Deleted workflow " + item.Name
	return nil, out, nil
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
