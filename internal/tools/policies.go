package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type policyItem struct {
	ID        string `json:"id"`
	DisplayID string `json:"display_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Revision  string `json:"revision,omitempty"`
	Kind      string `json:"kind,omitempty"`
	ReviewDue string `json:"review_due,omitempty"`
}

func registerPolicies(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_policies_list",
		Title:       "List Openlane policies",
		Description: "List internal policies in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listPolicies)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_get",
		Title:       "Get an Openlane policy",
		Description: "Get a single internal policy by ID.",
		Annotations: readOnly(),
	}, h.getPolicy)
}

func (h *handlers) listPolicies(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[policyItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetInternalPolicies(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[policyItem]{}, openlane.APIError(err)
	}
	items := make([]policyItem, 0, len(resp.InternalPolicies.Edges))
	for _, e := range resp.InternalPolicies.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, policyItem{
			ID:        n.ID,
			DisplayID: n.DisplayID,
			Name:      n.Name,
			Status:    openlane.Format(n.Status),
			Summary:   openlane.Deref(n.Summary),
		})
	}
	return nil, openlane.Page[policyItem]{
		Items:      items,
		NextCursor: resp.InternalPolicies.PageInfo.EndCursor,
		HasMore:    resp.InternalPolicies.PageInfo.HasNextPage,
		TotalCount: resp.InternalPolicies.TotalCount,
	}, nil
}

func (h *handlers) getPolicy(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, policyItem, error) {
	if in.ID == "" {
		return nil, policyItem{}, errIDRequired
	}
	resp, err := h.api.GetInternalPolicyByID(ctx, in.ID)
	if err != nil {
		return nil, policyItem{}, openlane.APIError(err)
	}
	p := resp.InternalPolicy
	return nil, policyItem{
		ID:        p.ID,
		DisplayID: p.DisplayID,
		Name:      p.Name,
		Status:    openlane.Format(p.Status),
		Summary:   openlane.Deref(p.Summary),
		Revision:  openlane.Deref(p.Revision),
		Kind:      openlane.Deref(p.InternalPolicyKindName),
		ReviewDue: openlane.Format(p.ReviewDue),
	}, nil
}
