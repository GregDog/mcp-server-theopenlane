package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type workflowDefinitionItem struct {
	ID              string                             `json:"id"`
	DisplayID       string                             `json:"display_id,omitempty"`
	Name            string                             `json:"name,omitempty"`
	Description     string                             `json:"description,omitempty"`
	SchemaType      string                             `json:"schema_type,omitempty"`
	WorkflowKind    string                             `json:"workflow_kind,omitempty"`
	Revision        int64                              `json:"revision,omitempty"`
	Active          bool                               `json:"active"`
	Draft           bool                               `json:"draft"`
	CooldownSeconds int64                              `json:"cooldown_seconds,omitempty"`
	TrackedFields   []string                           `json:"tracked_fields,omitempty"`
	DefinitionJSON  *models.WorkflowDefinitionDocument `json:"definition_json,omitempty"`
	Summary         string                             `json:"summary,omitempty"`
}

func registerWorkflows(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflows_list",
		Title:       "List Openlane workflow definitions",
		Description: "List workflow definitions in the configured organization. Filter by schema type, kind, active/draft state, or tracked field. Results are paginated.",
		Annotations: readOnly(),
	}, h.listWorkflows)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflows_search",
		Title:       "Search Openlane workflow definitions",
		Description: "Search workflow definitions by name or description. Results are paginated.",
		Annotations: readOnly(),
	}, h.searchWorkflows)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_get",
		Title:       "Get an Openlane workflow definition",
		Description: "Get a workflow definition by ID including triggers, conditions, actions, and a plain-English summary.",
		Annotations: readOnly(),
	}, h.getWorkflow)
}

func (h *handlers) listWorkflows(ctx context.Context, _ *mcp.CallToolRequest, in workflowListInput) (*mcp.CallToolResult, openlane.Page[workflowDefinitionItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetWorkflowDefinitions(ctx, &first, after, buildWorkflowDefinitionWhere(in))
	if err != nil {
		return nil, openlane.Page[workflowDefinitionItem]{}, openlane.APIError(err)
	}
	return nil, mapWorkflowDefinitionPage(resp), nil
}

func (h *handlers) searchWorkflows(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, openlane.Page[workflowDefinitionItem], error) {
	if in.Query == "" {
		return nil, openlane.Page[workflowDefinitionItem]{}, errQueryRequired
	}
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetWorkflowDefinitions(ctx, &first, after, buildWorkflowDefinitionSearchWhere(in.Query))
	if err != nil {
		return nil, openlane.Page[workflowDefinitionItem]{}, openlane.APIError(err)
	}
	return nil, mapWorkflowDefinitionPage(resp), nil
}

func (h *handlers) getWorkflow(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, workflowDefinitionItem, error) {
	if in.ID == "" {
		return nil, workflowDefinitionItem{}, errIDRequired
	}
	resp, err := h.api.GetWorkflowDefinitionByID(ctx, in.ID)
	if err != nil {
		return nil, workflowDefinitionItem{}, openlane.APIError(err)
	}
	return nil, mapWorkflowDefinition(resp.WorkflowDefinition), nil
}

func mapWorkflowDefinitionPage(resp *graphclient.GetWorkflowDefinitions) openlane.Page[workflowDefinitionItem] {
	items := make([]workflowDefinitionItem, 0, len(resp.WorkflowDefinitions.Edges))
	for _, e := range resp.WorkflowDefinitions.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapWorkflowDefinitionNode(*e.Node))
	}
	return openlane.Page[workflowDefinitionItem]{
		Items:      items,
		NextCursor: resp.WorkflowDefinitions.PageInfo.EndCursor,
		HasMore:    resp.WorkflowDefinitions.PageInfo.HasNextPage,
		TotalCount: resp.WorkflowDefinitions.TotalCount,
	}
}

func mapWorkflowDefinition(w graphclient.GetWorkflowDefinitionByID_WorkflowDefinition) workflowDefinitionItem {
	return mapWorkflowDefinitionFields(
		w.ID, w.DisplayID, w.Name, openlane.Deref(w.Description), w.SchemaType,
		openlane.Format(w.WorkflowKind), w.Revision, w.Active, w.Draft,
		w.CooldownSeconds, w.TrackedFields, w.DefinitionJSON,
	)
}

func mapWorkflowDefinitionNode(w graphclient.GetWorkflowDefinitions_WorkflowDefinitions_Edges_Node) workflowDefinitionItem {
	return mapWorkflowDefinitionFields(
		w.ID, w.DisplayID, w.Name, openlane.Deref(w.Description), w.SchemaType,
		openlane.Format(w.WorkflowKind), w.Revision, w.Active, w.Draft,
		w.CooldownSeconds, w.TrackedFields, w.DefinitionJSON,
	)
}

func mapWorkflowDefinitionFields(
	id, displayID, name, description, schemaType, kind string,
	revision int64, active, draft bool, cooldown int64,
	tracked []string, doc *models.WorkflowDefinitionDocument,
) workflowDefinitionItem {
	return workflowDefinitionItem{
		ID:              id,
		DisplayID:       displayID,
		Name:            name,
		Description:     description,
		SchemaType:      schemaType,
		WorkflowKind:    kind,
		Revision:        revision,
		Active:          active,
		Draft:           draft,
		CooldownSeconds: cooldown,
		TrackedFields:   tracked,
		DefinitionJSON:  doc,
		Summary:         summarizeWorkflowDefinition(name, schemaType, kind, active, draft, revision, cooldown, tracked, doc),
	}
}
