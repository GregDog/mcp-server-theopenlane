package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type controlImplementationItem struct {
	ID                 string   `json:"id"`
	Status             string   `json:"status,omitempty"`
	Verified           *bool    `json:"verified,omitempty"`
	Details            string   `json:"details,omitempty"`
	ImplementationDate string   `json:"implementation_date,omitempty"`
	VerificationDate   string   `json:"verification_date,omitempty"`
	OwnerID            string   `json:"owner_id,omitempty"`
	ControlRefCodes    []string `json:"control_ref_codes,omitempty"`
	SubcontrolRefCodes []string `json:"subcontrol_ref_codes,omitempty"`
	Tags               []string `json:"tags,omitempty"`
}

func registerImplementations(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_control_implementations_list",
		Title:       "List Openlane control implementations",
		Description: "List control implementations in the configured Openlane organization. Supports control, status, and verified filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.listControlImplementations)

	addTool(server, &mcp.Tool{
		Name:        "openlane_control_implementation_get",
		Title:       "Get an Openlane control implementation",
		Description: "Get a single control implementation by ID, including linked control and subcontrol ref codes.",
		Annotations: readOnly(),
	}, h.getControlImplementation)
}

func (h *handlers) listControlImplementations(ctx context.Context, _ *mcp.CallToolRequest, in implementationListInput) (*mcp.CallToolResult, openlane.Page[controlImplementationItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControlImplementations(ctx, &first, after, buildImplementationWhere(in))
	if err != nil {
		return nil, openlane.Page[controlImplementationItem]{}, openlane.APIError(err)
	}
	items := make([]controlImplementationItem, 0, len(resp.ControlImplementations.Edges))
	for _, e := range resp.ControlImplementations.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListImplementation(e.Node))
	}
	return nil, openlane.Page[controlImplementationItem]{
		Items:      items,
		NextCursor: resp.ControlImplementations.PageInfo.EndCursor,
		HasMore:    resp.ControlImplementations.PageInfo.HasNextPage,
		TotalCount: resp.ControlImplementations.TotalCount,
	}, nil
}

func (h *handlers) getControlImplementation(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, controlImplementationItem, error) {
	if in.ID == "" {
		return nil, controlImplementationItem{}, errIDRequired
	}
	resp, err := h.api.GetControlImplementationByID(ctx, in.ID)
	if err != nil {
		return nil, controlImplementationItem{}, openlane.APIError(err)
	}
	return nil, mapGetImplementation(resp.ControlImplementation), nil
}

func mapListImplementation(n *graphclient.GetControlImplementations_ControlImplementations_Edges_Node) controlImplementationItem {
	return controlImplementationItem{
		ID:              n.ID,
		Status:          openlane.Format(n.Status),
		Verified:        n.Verified,
		ControlRefCodes: controlRefCodesFromEdges(n.Controls.Edges),
	}
}

func mapGetImplementation(n graphclient.GetControlImplementationByID_ControlImplementation) controlImplementationItem {
	return controlImplementationItem{
		ID:                 n.ID,
		Status:             openlane.Format(n.Status),
		Verified:           n.Verified,
		Details:            openlane.Deref(n.Details),
		ImplementationDate: openlane.Format(n.ImplementationDate),
		VerificationDate:   openlane.Format(n.VerificationDate),
		OwnerID:            openlane.Deref(n.OwnerID),
		ControlRefCodes:    controlRefCodesFromGetEdges(n.Controls.Edges),
		SubcontrolRefCodes: subcontrolRefCodesFromGetEdges(n.Subcontrols.Edges),
		Tags:               n.Tags,
	}
}

func controlRefCodesFromEdges(edges []*graphclient.GetControlImplementations_ControlImplementations_Edges_Node_Controls_Edges) []string {
	out := make([]string, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil || e.Node.RefCode == "" {
			continue
		}
		out = append(out, e.Node.RefCode)
	}
	return out
}

func controlRefCodesFromGetEdges(edges []*graphclient.GetControlImplementationByID_ControlImplementation_Controls_Edges) []string {
	out := make([]string, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil || e.Node.RefCode == "" {
			continue
		}
		out = append(out, e.Node.RefCode)
	}
	return out
}

func subcontrolRefCodesFromGetEdges(edges []*graphclient.GetControlImplementationByID_ControlImplementation_Subcontrols_Edges) []string {
	out := make([]string, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil || e.Node.RefCode == "" {
			continue
		}
		out = append(out, e.Node.RefCode)
	}
	return out
}
