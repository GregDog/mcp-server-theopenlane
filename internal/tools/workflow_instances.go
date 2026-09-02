package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type workflowAssignmentRef struct {
	ID        string `json:"id"`
	DisplayID string `json:"display_id,omitempty"`
	Status    string `json:"status,omitempty"`
	Role      string `json:"role,omitempty"`
	Label     string `json:"label,omitempty"`
	Required  bool   `json:"required"`
}

type workflowInstanceItem struct {
	ID                   string                            `json:"id"`
	DisplayID            string                            `json:"display_id,omitempty"`
	State                string                            `json:"state,omitempty"`
	WorkflowDefinitionID string                            `json:"workflow_definition_id,omitempty"`
	WorkflowProposalID   string                            `json:"workflow_proposal_id,omitempty"`
	CreatedBy            string                            `json:"created_by,omitempty"`
	CurrentActionIndex   *int64                            `json:"current_action_index,omitempty"`
	ObjectType           string                            `json:"object_type,omitempty"`
	ObjectID             string                            `json:"object_id,omitempty"`
	InternalPolicyID     string                            `json:"internal_policy_id,omitempty"`
	ControlID            string                            `json:"control_id,omitempty"`
	EvidenceID           string                            `json:"evidence_id,omitempty"`
	Assignments          []workflowAssignmentRef           `json:"assignments,omitempty"`
	ProposalPreview      *openlane.WorkflowProposalPreview `json:"proposal_preview,omitempty"`
	Events               *relSummary[workflowEventRef]     `json:"events,omitempty"`
	Summary              string                            `json:"summary,omitempty"`
}

type workflowEventRef struct {
	ID        string `json:"id"`
	EventType string `json:"event_type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func registerWorkflowInstances(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_instances_list",
		Title:       "List Openlane workflow instances",
		Description: "List workflow instances. Filter by definition, state, or target object type (and optional object ID). object_type alone matches any instance linked to that type. Results are paginated.",
		Annotations: readOnly(),
	}, h.listWorkflowInstances)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_instance_get",
		Title:       "Get an Openlane workflow instance",
		Description: "Get a workflow instance by ID including assignments, proposal preview, and a bounded event timeline.",
		Annotations: readOnly(),
	}, h.getWorkflowInstance)
}

func (h *handlers) listWorkflowInstances(ctx context.Context, _ *mcp.CallToolRequest, in workflowInstanceListInput) (*mcp.CallToolResult, openlane.Page[workflowInstanceItem], error) {
	where, err := buildWorkflowInstanceWhere(in)
	if err != nil {
		return nil, openlane.Page[workflowInstanceItem]{}, err
	}
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetWorkflowInstances(ctx, &first, after, where)
	if err != nil {
		return nil, openlane.Page[workflowInstanceItem]{}, openlane.APIError(err)
	}
	items := make([]workflowInstanceItem, 0, len(resp.WorkflowInstances.Edges))
	for _, e := range resp.WorkflowInstances.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapWorkflowInstanceNode(*e.Node))
	}
	return nil, openlane.Page[workflowInstanceItem]{
		Items:      items,
		NextCursor: resp.WorkflowInstances.PageInfo.EndCursor,
		HasMore:    resp.WorkflowInstances.PageInfo.HasNextPage,
		TotalCount: resp.WorkflowInstances.TotalCount,
	}, nil
}

func (h *handlers) getWorkflowInstance(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, workflowInstanceItem, error) {
	if in.ID == "" {
		return nil, workflowInstanceItem{}, errIDRequired
	}
	resp, err := h.api.GetWorkflowInstanceByID(ctx, in.ID)
	if err != nil {
		return nil, workflowInstanceItem{}, openlane.APIError(err)
	}
	item := mapWorkflowInstance(resp.WorkflowInstance)

	detail, err := h.api.GetWorkflowInstanceDetail(ctx, in.ID)
	if err == nil && detail != nil {
		item.CurrentActionIndex = detail.CurrentActionIndex
		item.ProposalPreview = detail.ProposalPreview
	}

	events, err := h.fetchWorkflowEvents(ctx, in.ID)
	if err == nil {
		item.Events = events
	}

	item.Summary = summarizeWorkflowInstance(item)
	return nil, item, nil
}

func (h *handlers) fetchWorkflowEvents(ctx context.Context, instanceID string) (*relSummary[workflowEventRef], error) {
	first := relFirst()
	where := &graphclient.WorkflowEventWhereInput{
		WorkflowInstanceID: &instanceID,
	}
	resp, err := h.api.GetWorkflowEvents(ctx, &first, nil, where)
	if err != nil {
		return nil, err
	}
	items := make([]workflowEventRef, 0, len(resp.WorkflowEvents.Edges))
	for _, e := range resp.WorkflowEvents.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, workflowEventRef{
			ID:        n.ID,
			EventType: openlane.Format(n.EventType),
			CreatedAt: openlane.Format(n.CreatedAt),
		})
	}
	return &relSummary[workflowEventRef]{
		Count: resp.WorkflowEvents.TotalCount,
		Items: items,
	}, nil
}

func mapWorkflowInstance(w graphclient.GetWorkflowInstanceByID_WorkflowInstance) workflowInstanceItem {
	objType, objID := workflowInstanceObjectIDs(
		w.InternalPolicyID, w.ControlID, w.EvidenceID,
		w.SubcontrolID, w.ActionPlanID, w.ProcedureID,
		w.CampaignID, w.CampaignTargetID, w.PlatformID, w.IdentityHolderID,
	)
	item := workflowInstanceItem{
		ID:                   w.ID,
		DisplayID:            w.DisplayID,
		State:                openlane.Format(w.State),
		WorkflowDefinitionID: w.WorkflowDefinitionID,
		WorkflowProposalID:   openlane.Deref(w.WorkflowProposalID),
		CreatedBy:            openlane.Deref(w.CreatedBy),
		ObjectType:           objType,
		ObjectID:             objID,
		InternalPolicyID:     openlane.Deref(w.InternalPolicyID),
		ControlID:            openlane.Deref(w.ControlID),
		EvidenceID:           openlane.Deref(w.EvidenceID),
	}
	for _, e := range w.WorkflowAssignments.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		item.Assignments = append(item.Assignments, workflowAssignmentRef{
			ID:        n.ID,
			DisplayID: n.DisplayID,
			Status:    openlane.Format(n.Status),
			Role:      n.Role,
			Label:     openlane.Deref(n.Label),
			Required:  n.Required,
		})
	}
	return item
}

func mapWorkflowInstanceNode(w graphclient.GetWorkflowInstances_WorkflowInstances_Edges_Node) workflowInstanceItem {
	objType, objID := workflowInstanceObjectIDs(
		w.InternalPolicyID, w.ControlID, w.EvidenceID,
		w.SubcontrolID, w.ActionPlanID, w.ProcedureID,
		w.CampaignID, w.CampaignTargetID, w.PlatformID, w.IdentityHolderID,
	)
	item := workflowInstanceItem{
		ID:                   w.ID,
		DisplayID:            w.DisplayID,
		State:                openlane.Format(w.State),
		WorkflowDefinitionID: w.WorkflowDefinitionID,
		WorkflowProposalID:   openlane.Deref(w.WorkflowProposalID),
		CreatedBy:            openlane.Deref(w.CreatedBy),
		ObjectType:           objType,
		ObjectID:             objID,
		InternalPolicyID:     openlane.Deref(w.InternalPolicyID),
		ControlID:            openlane.Deref(w.ControlID),
		EvidenceID:           openlane.Deref(w.EvidenceID),
	}
	for _, e := range w.WorkflowAssignments.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		item.Assignments = append(item.Assignments, workflowAssignmentRef{
			ID:        n.ID,
			DisplayID: n.DisplayID,
			Status:    openlane.Format(n.Status),
			Role:      n.Role,
			Label:     openlane.Deref(n.Label),
			Required:  n.Required,
		})
	}
	return item
}

func workflowInstanceObjectIDs(
	internalPolicyID, controlID, evidenceID,
	subcontrolID, actionPlanID, procedureID,
	campaignID, campaignTargetID, platformID, identityHolderID *string,
) (typ, id string) {
	type idPair struct {
		typ string
		id  *string
	}
	candidates := []idPair{
		{"InternalPolicy", internalPolicyID},
		{"Control", controlID},
		{"Evidence", evidenceID},
		{"Subcontrol", subcontrolID},
		{"ActionPlan", actionPlanID},
		{"Procedure", procedureID},
		{"Campaign", campaignID},
		{"CampaignTarget", campaignTargetID},
		{"Platform", platformID},
		{"IdentityHolder", identityHolderID},
	}
	for _, c := range candidates {
		if c.id != nil && *c.id != "" {
			return c.typ, *c.id
		}
	}
	return "", ""
}

func summarizeWorkflowInstance(item workflowInstanceItem) string {
	s := "Workflow instance " + item.DisplayID
	if item.State != "" {
		s += " is " + item.State
	}
	if item.ObjectType != "" {
		s += " on " + item.ObjectType + " " + item.ObjectID
	}
	if item.WorkflowDefinitionID != "" {
		s += " (definition " + item.WorkflowDefinitionID + ")"
	}
	if len(item.Assignments) > 0 {
		pending := 0
		for _, a := range item.Assignments {
			if a.Status == "PENDING" {
				pending++
			}
		}
		s += fmt.Sprintf("; %d assignment(s), %d pending", len(item.Assignments), pending)
	}
	return s
}
