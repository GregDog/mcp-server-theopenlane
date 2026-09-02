package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type policyAwaitingApprovalItem struct {
	Path                string `json:"path"`
	PolicyID            string `json:"policy_id"`
	DisplayID           string `json:"display_id,omitempty"`
	Name                string `json:"name,omitempty"`
	Status              string `json:"status,omitempty"`
	AssignmentID        string `json:"assignment_id,omitempty"`
	AssignmentDisplayID string `json:"assignment_display_id,omitempty"`
	WorkflowInstanceID  string `json:"workflow_instance_id,omitempty"`
	DueAt               string `json:"due_at,omitempty"`
	Summary             string `json:"summary,omitempty"`
}

type policiesAwaitingApprovalResult struct {
	Items         []policyAwaitingApprovalItem `json:"items"`
	NativeCount   int                          `json:"native_count"`
	WorkflowCount int                          `json:"workflow_count"`
	TotalCount    int                          `json:"total_count"`
	Summary       string                       `json:"summary"`
}

type policyRef struct {
	name      string
	displayID string
}

func registerPoliciesAwaitingApproval(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:  "openlane_policies_awaiting_approval",
		Title: "List Openlane policies awaiting approval",
		Description: "List policies that need approval via either path: native InternalPolicy status NEEDS_APPROVAL (org-wide), or your pending WorkflowAssignments on InternalPolicy objects. Use openlane_policy_approve for native items; use openlane_workflow_assignment_approve for workflow items. Results are bounded by limit (default 20, max 50 per path).",
		Annotations: readOnly(),
	}, h.listPoliciesAwaitingApproval)
}

func (h *handlers) listPoliciesAwaitingApproval(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, policiesAwaitingApprovalResult, error) {
	first := openlane.ClampLimit(in.Limit)
	out := policiesAwaitingApprovalResult{Items: []policyAwaitingApprovalItem{}}

	native, err := h.listNativePoliciesAwaitingApproval(ctx, first)
	if err != nil {
		return nil, policiesAwaitingApprovalResult{}, err
	}
	out.Items = append(out.Items, native...)
	out.NativeCount = len(native)

	workflow, err := h.listWorkflowPoliciesAwaitingApproval(ctx, first)
	if err != nil {
		return nil, policiesAwaitingApprovalResult{}, err
	}
	out.Items = append(out.Items, workflow...)
	out.WorkflowCount = len(workflow)
	out.TotalCount = len(out.Items)
	out.Summary = fmt.Sprintf("%d policy approval item(s): %d native (NEEDS_APPROVAL), %d workflow assignment(s) pending for you",
		out.TotalCount, out.NativeCount, out.WorkflowCount)
	return nil, out, nil
}

func (h *handlers) listNativePoliciesAwaitingApproval(ctx context.Context, first int64) ([]policyAwaitingApprovalItem, error) {
	status := enums.DocumentNeedsApproval
	resp, err := h.api.GetInternalPolicies(ctx, &first, nil, &graphclient.InternalPolicyWhereInput{Status: &status})
	if err != nil {
		return nil, openlane.APIError(err)
	}
	items := make([]policyAwaitingApprovalItem, 0, len(resp.InternalPolicies.Edges))
	for _, e := range resp.InternalPolicies.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, policyAwaitingApprovalItem{
			Path:      "native",
			PolicyID:  n.ID,
			DisplayID: n.DisplayID,
			Name:      n.Name,
			Status:    openlane.Format(n.Status),
			Summary:   fmt.Sprintf("Native approval: %s (%s) is NEEDS_APPROVAL — use openlane_policy_approve", n.Name, n.DisplayID),
		})
	}
	return items, nil
}

func (h *handlers) listWorkflowPoliciesAwaitingApproval(ctx context.Context, first int64) ([]policyAwaitingApprovalItem, error) {
	pending := enums.WorkflowAssignmentStatusPending
	hasPolicy := true
	where := &graphclient.WorkflowAssignmentWhereInput{
		Status: &pending,
		HasWorkflowInstanceWith: []*graphclient.WorkflowInstanceWhereInput{{
			InternalPolicyIDNotNil: &hasPolicy,
		}},
	}
	resp, err := h.api.GetMyWorkflowAssignments(ctx, &first, nil, where)
	if err != nil {
		return nil, openlane.APIError(err)
	}
	items := make([]policyAwaitingApprovalItem, 0, len(resp.MyWorkflowAssignments.Edges))
	policies := map[string]policyRef{}
	for _, e := range resp.MyWorkflowAssignments.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		detail, err := h.api.GetWorkflowAssignmentDetail(ctx, n.ID)
		if err != nil {
			return nil, openlane.APIError(err)
		}
		policyID := detail.InstanceInternalPolicyID
		ref := policies[policyID]
		if policyID != "" && ref.name == "" {
			if got, err := h.api.GetInternalPolicyByID(ctx, policyID); err == nil {
				ref = policyRef{name: got.InternalPolicy.Name, displayID: got.InternalPolicy.DisplayID}
				policies[policyID] = ref
			}
		}
		items = append(items, policyAwaitingApprovalItem{
			Path:                "workflow",
			PolicyID:            policyID,
			DisplayID:           ref.displayID,
			Name:                ref.name,
			Status:              detail.Status,
			AssignmentID:        detail.ID,
			AssignmentDisplayID: detail.DisplayID,
			WorkflowInstanceID:  detail.WorkflowInstanceID,
			DueAt:               detail.DueAt,
			Summary:             fmt.Sprintf("Workflow approval: assignment %s for policy %s — use openlane_workflow_assignment_approve", detail.DisplayID, nameOrID(ref.name, policyID)),
		})
	}
	return items, nil
}

func nameOrID(name, id string) string {
	if name != "" {
		return name
	}
	return id
}
