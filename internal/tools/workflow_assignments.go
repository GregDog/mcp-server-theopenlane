package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type workflowAssignmentItem struct {
	ID                   string                              `json:"id"`
	DisplayID            string                              `json:"display_id,omitempty"`
	Status               string                              `json:"status,omitempty"`
	Role                 string                              `json:"role,omitempty"`
	Label                string                              `json:"label,omitempty"`
	Required             bool                                `json:"required"`
	DueAt                string                              `json:"due_at,omitempty"`
	Notes                string                              `json:"notes,omitempty"`
	WorkflowInstanceID   string                              `json:"workflow_instance_id,omitempty"`
	WorkflowDefinitionID string                              `json:"workflow_definition_id,omitempty"`
	ActorUserID          string                              `json:"actor_user_id,omitempty"`
	ActorGroupID         string                              `json:"actor_group_id,omitempty"`
	DecidedAt            string                              `json:"decided_at,omitempty"`
	ObjectType           string                              `json:"object_type,omitempty"`
	ObjectID             string                              `json:"object_id,omitempty"`
	InternalPolicyID     string                              `json:"internal_policy_id,omitempty"`
	ControlID            string                              `json:"control_id,omitempty"`
	EvidenceID           string                              `json:"evidence_id,omitempty"`
	InstanceState        string                              `json:"instance_state,omitempty"`
	Targets              []openlane.WorkflowAssignmentTarget `json:"targets,omitempty"`
	Summary              string                              `json:"summary,omitempty"`
}

func registerWorkflowAssignments(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignments_list",
		Title:       "List my Openlane workflow assignments",
		Description: "List workflow approval assignments for the current user. Use for WorkflowDefinition-driven approvals on policies, controls, and other objects. For native InternalPolicy status NEEDS_APPROVAL (no WorkflowAssignment), use openlane_policies_awaiting_approval and openlane_policy_approve instead. Defaults to outstanding approvals; filter by status or instance ID. Results are paginated.",
		Annotations: readOnly(),
	}, h.listWorkflowAssignments)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_assignment_get",
		Title:       "Get an Openlane workflow assignment",
		Description: "Get a workflow assignment by ID including approval targets, due date, and linked workflow instance/object context.",
		Annotations: readOnly(),
	}, h.getWorkflowAssignment)
}

func (h *handlers) listWorkflowAssignments(ctx context.Context, _ *mcp.CallToolRequest, in workflowAssignmentListInput) (*mcp.CallToolResult, openlane.Page[workflowAssignmentItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	where := buildWorkflowAssignmentWhere(in)
	resp, err := h.api.GetMyWorkflowAssignments(ctx, &first, after, where)
	if err != nil {
		return nil, openlane.Page[workflowAssignmentItem]{}, openlane.APIError(err)
	}
	items := make([]workflowAssignmentItem, 0, len(resp.MyWorkflowAssignments.Edges))
	for _, e := range resp.MyWorkflowAssignments.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapMyWorkflowAssignmentNode(*e.Node))
	}
	return nil, openlane.Page[workflowAssignmentItem]{
		Items:      items,
		NextCursor: resp.MyWorkflowAssignments.PageInfo.EndCursor,
		HasMore:    resp.MyWorkflowAssignments.PageInfo.HasNextPage,
		TotalCount: resp.MyWorkflowAssignments.TotalCount,
	}, nil
}

func (h *handlers) getWorkflowAssignment(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, workflowAssignmentItem, error) {
	if in.ID == "" {
		return nil, workflowAssignmentItem{}, errIDRequired
	}
	detail, err := h.api.GetWorkflowAssignmentDetail(ctx, in.ID)
	if err != nil {
		return nil, workflowAssignmentItem{}, openlane.APIError(err)
	}
	item := mapWorkflowAssignmentDetail(detail)
	return nil, item, nil
}

func mapMyWorkflowAssignmentNode(n graphclient.GetMyWorkflowAssignments_MyWorkflowAssignments_Edges_Node) workflowAssignmentItem {
	item := workflowAssignmentItem{
		ID:                 n.ID,
		DisplayID:          n.DisplayID,
		Status:             openlane.Format(n.Status),
		Role:               n.Role,
		Label:              openlane.Deref(n.Label),
		Required:           n.Required,
		Notes:              openlane.Deref(n.Notes),
		WorkflowInstanceID: n.WorkflowInstanceID,
		ActorUserID:        openlane.Deref(n.ActorUserID),
		ActorGroupID:       openlane.Deref(n.ActorGroupID),
		DecidedAt:          openlane.Format(n.DecidedAt),
	}
	item.Summary = summarizeWorkflowAssignment(item)
	return item
}

func mapWorkflowAssignmentDetail(d *openlane.WorkflowAssignmentDetail) workflowAssignmentItem {
	item := workflowAssignmentItem{
		ID:                   d.ID,
		DisplayID:            d.DisplayID,
		Status:               d.Status,
		Role:                 d.Role,
		Label:                d.Label,
		Required:             d.Required,
		DueAt:                d.DueAt,
		Notes:                d.Notes,
		WorkflowInstanceID:   d.WorkflowInstanceID,
		WorkflowDefinitionID: d.InstanceWorkflowDefID,
		ActorUserID:          d.ActorUserID,
		ActorGroupID:         d.ActorGroupID,
		DecidedAt:            d.DecidedAt,
		InstanceState:        d.InstanceState,
		InternalPolicyID:     d.InstanceInternalPolicyID,
		ControlID:            d.InstanceControlID,
		EvidenceID:           d.InstanceEvidenceID,
		Targets:              d.Targets,
	}
	item.ObjectType, item.ObjectID = assignmentObject(d)
	item.Summary = summarizeWorkflowAssignment(item)
	return item
}

func assignmentObject(d *openlane.WorkflowAssignmentDetail) (typ, id string) {
	if d.InstanceInternalPolicyID != "" {
		return "InternalPolicy", d.InstanceInternalPolicyID
	}
	if d.InstanceControlID != "" {
		return "Control", d.InstanceControlID
	}
	if d.InstanceEvidenceID != "" {
		return "Evidence", d.InstanceEvidenceID
	}
	return "", ""
}

func summarizeWorkflowAssignment(item workflowAssignmentItem) string {
	s := "Assignment " + item.DisplayID
	if item.Status != "" {
		s += " (" + item.Status + ")"
	}
	if item.ObjectType != "" {
		s += " for " + item.ObjectType
	}
	if item.DueAt != "" {
		s += "; due " + item.DueAt
	}
	if len(item.Targets) > 0 {
		s += "; " + summarizeAssignmentTargets(item.Targets)
	}
	return s
}

func summarizeAssignmentTargets(targets []openlane.WorkflowAssignmentTarget) string {
	var parts []string
	for _, t := range targets {
		switch t.TargetType {
		case "RESOLVER":
			parts = append(parts, "RESOLVER:"+t.ResolverKey)
		case "GROUP":
			parts = append(parts, "GROUP")
		case "USER":
			parts = append(parts, "USER")
		default:
			parts = append(parts, t.TargetType)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += ", " + p
	}
	return out
}
