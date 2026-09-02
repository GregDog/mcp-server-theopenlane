package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type workflowAssignmentActionInput struct {
	ID      string `json:"id" jsonschema:"Workflow assignment ID. Do not pass a policy ID."`
	Reason  string `json:"reason,omitempty" jsonschema:"Required for reject; optional for request_changes."`
	Confirm bool   `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type workflowAssignmentReassignInput struct {
	ID           string `json:"id" jsonschema:"Existing workflow assignment ID to extend."`
	TargetUserID string `json:"target_user_id,omitempty" jsonschema:"User ID to assign an additional approval to."`
	TargetUser   string `json:"target_user,omitempty" jsonschema:"User email or display name to resolve (errors on 0 or many matches)."`
	Confirm      bool   `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type workflowAssignmentActionResult struct {
	Confirmed            bool   `json:"confirmed"`
	Error                string `json:"error,omitempty"`
	ID                   string `json:"id"`
	DisplayID            string `json:"display_id,omitempty"`
	Status               string `json:"status,omitempty"`
	WorkflowInstanceID   string `json:"workflow_instance_id,omitempty"`
	WorkflowDefinitionID string `json:"workflow_definition_id,omitempty"`
	ObjectType           string `json:"object_type,omitempty"`
	ObjectID             string `json:"object_id,omitempty"`
	ResultStatus         string `json:"result_status,omitempty"`
	Summary              string `json:"summary,omitempty"`
}

func registerWriteWorkflowAssignments(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignment_approve",
		Title:       "Approve an Openlane workflow assignment",
		Description: "Approve a WorkflowAssignment by ID using approveWorkflowAssignment. Use only when openlane_policies_awaiting_approval returns path workflow or a real WorkflowInstance/Assignment exists. For native InternalPolicy NEEDS_APPROVAL with no assignment, use openlane_policy_approve instead. Refuses bulk approve. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.approveWorkflowAssignment)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignment_reject",
		Title:       "Reject an Openlane workflow assignment",
		Description: "Reject a WorkflowAssignment by ID. Requires reason. Use only when a real WorkflowInstance/Assignment exists. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.rejectWorkflowAssignment)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignment_request_changes",
		Title:       "Request changes on an Openlane workflow assignment",
		Description: "Request changes on a WorkflowAssignment (custom GraphQL; not in go-client v0.14.0). Use only when a real assignment exists. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.requestChangesWorkflowAssignment)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignment_reassign",
		Title:       "Reassign an Openlane workflow assignment",
		Description: "Add another user to a WorkflowAssignment via reassignWorkflowAssignment (custom GraphQL; not in go-client v0.14.0). Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.reassignWorkflowAssignment)
}

func (h *handlers) approveWorkflowAssignment(ctx context.Context, _ *mcp.CallToolRequest, in workflowAssignmentActionInput) (*mcp.CallToolResult, workflowAssignmentActionResult, error) {
	preview, err := h.previewAssignmentAction(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, err
	}
	if preview.Status != string(enums.WorkflowAssignmentStatusPending) {
		return nil, workflowAssignmentActionResult{}, fmt.Errorf("assignment status is %s, not PENDING", preview.Status)
	}
	if !in.Confirm {
		preview.Error = errConfirmationRequired
		preview.Summary = "Would approve " + preview.Summary
		return nil, preview, nil
	}
	if _, err := h.api.ApproveWorkflowAssignment(ctx, in.ID); err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	return h.verifyAssignmentAction(ctx, in.ID, preview, string(enums.WorkflowAssignmentStatusApproved), "Approved")
}

func (h *handlers) rejectWorkflowAssignment(ctx context.Context, _ *mcp.CallToolRequest, in workflowAssignmentActionInput) (*mcp.CallToolResult, workflowAssignmentActionResult, error) {
	if strings.TrimSpace(in.Reason) == "" {
		return nil, workflowAssignmentActionResult{}, fmt.Errorf("reason is required")
	}
	preview, err := h.previewAssignmentAction(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, err
	}
	if preview.Status != string(enums.WorkflowAssignmentStatusPending) {
		return nil, workflowAssignmentActionResult{}, fmt.Errorf("assignment status is %s, not PENDING", preview.Status)
	}
	if !in.Confirm {
		preview.Error = errConfirmationRequired
		preview.Summary = "Would reject " + preview.Summary
		return nil, preview, nil
	}
	reason := strings.TrimSpace(in.Reason)
	if _, err := h.api.RejectWorkflowAssignment(ctx, in.ID, &reason); err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	return h.verifyAssignmentAction(ctx, in.ID, preview, string(enums.WorkflowAssignmentStatusRejected), "Rejected")
}

func (h *handlers) requestChangesWorkflowAssignment(ctx context.Context, _ *mcp.CallToolRequest, in workflowAssignmentActionInput) (*mcp.CallToolResult, workflowAssignmentActionResult, error) {
	preview, err := h.previewAssignmentAction(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, err
	}
	if preview.Status != string(enums.WorkflowAssignmentStatusPending) {
		return nil, workflowAssignmentActionResult{}, fmt.Errorf("assignment status is %s, not PENDING", preview.Status)
	}
	if !in.Confirm {
		preview.Error = errConfirmationRequired
		preview.Summary = "Would request changes on " + preview.Summary
		return nil, preview, nil
	}
	var reason *string
	if s := strings.TrimSpace(in.Reason); s != "" {
		reason = &s
	}
	if err := h.api.RequestChangesWorkflowAssignment(ctx, in.ID, reason, nil); err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	return h.verifyAssignmentAction(ctx, in.ID, preview, string(enums.WorkflowAssignmentStatusChangesRequested), "Requested changes on")
}

func (h *handlers) reassignWorkflowAssignment(ctx context.Context, _ *mcp.CallToolRequest, in workflowAssignmentReassignInput) (*mcp.CallToolResult, workflowAssignmentActionResult, error) {
	preview, err := h.previewAssignmentAction(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, err
	}
	target := strings.TrimSpace(in.TargetUserID)
	if target == "" {
		target, err = h.resolveUserID(ctx, in.TargetUser)
		if err != nil {
			return nil, workflowAssignmentActionResult{}, err
		}
	}
	if !in.Confirm {
		preview.Error = errConfirmationRequired
		preview.Summary = "Would reassign " + preview.Summary + " to user " + target
		return nil, preview, nil
	}
	if _, err := h.api.ReassignWorkflowAssignment(ctx, in.ID, target); err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	got, err := h.api.GetWorkflowAssignmentDetail(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	preview.Confirmed = true
	preview.ResultStatus = got.Status
	preview.Summary = "Reassigned " + preview.Summary
	return nil, preview, nil
}

func (h *handlers) previewAssignmentAction(ctx context.Context, id string) (workflowAssignmentActionResult, error) {
	if id == "" {
		return workflowAssignmentActionResult{}, errIDRequired
	}
	d, err := h.api.GetWorkflowAssignmentDetail(ctx, id)
	if err != nil {
		return workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	item := mapWorkflowAssignmentDetail(d)
	return workflowAssignmentActionResult{
		ID:                   item.ID,
		DisplayID:            item.DisplayID,
		Status:               item.Status,
		WorkflowInstanceID:   item.WorkflowInstanceID,
		WorkflowDefinitionID: item.WorkflowDefinitionID,
		ObjectType:           item.ObjectType,
		ObjectID:             item.ObjectID,
		Summary:              item.Summary,
	}, nil
}

func (h *handlers) verifyAssignmentAction(ctx context.Context, id string, preview workflowAssignmentActionResult, wantStatus, verb string) (*mcp.CallToolResult, workflowAssignmentActionResult, error) {
	got, err := h.api.GetWorkflowAssignmentDetail(ctx, id)
	if err != nil {
		return nil, workflowAssignmentActionResult{}, openlane.APIError(err)
	}
	if got.Status != wantStatus {
		return nil, workflowAssignmentActionResult{}, fmt.Errorf("assignment persisted but status is %s, expected %s", got.Status, wantStatus)
	}
	preview.Confirmed = true
	preview.ResultStatus = got.Status
	preview.Status = got.Status
	preview.Summary = verb + " " + preview.Summary
	return nil, preview, nil
}
