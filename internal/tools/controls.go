package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type controlItem struct {
	ID                 string   `json:"id"`
	DisplayID          string   `json:"display_id,omitempty"`
	RefCode            string   `json:"ref_code,omitempty"`
	Title              string   `json:"title,omitempty"`
	Status             string   `json:"status,omitempty"`
	ReferenceFramework string   `json:"reference_framework,omitempty"`
	Category           string   `json:"category,omitempty"`
	Description        string   `json:"description,omitempty"`
	Subcategory        string   `json:"subcategory,omitempty"`
	StandardID         string   `json:"standard_id,omitempty"`
	OwnerID            string   `json:"owner_id,omitempty"`
	Tags               []string `json:"tags,omitempty"`
}

func registerControls(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_controls_list",
		Title:       "List Openlane controls",
		Description: "List controls in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listControls)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_controls_search",
		Title:       "Search Openlane controls",
		Description: "Search controls by ref code, title, or description using Openlane where-filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.searchControls)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_control_get",
		Title:       "Get an Openlane control",
		Description: "Get a single control by ID.",
		Annotations: readOnly(),
	}, h.getControl)
}

func (h *handlers) listControls(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControls(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[controlItem]{}, openlane.APIError(err)
	}
	return nil, mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount, false), nil
}

func (h *handlers) searchControls(ctx context.Context, _ *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	if in.Query == "" {
		return nil, openlane.Page[controlItem]{}, errQueryRequired
	}
	q := in.Query
	where := &graphclient.ControlWhereInput{
		Or: []*graphclient.ControlWhereInput{
			{RefCodeContainsFold: &q},
			{TitleContainsFold: &q},
			{DescriptionContainsFold: &q},
		},
	}
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControls(ctx, &first, after, where)
	if err != nil {
		return nil, openlane.Page[controlItem]{}, openlane.APIError(err)
	}
	return nil, mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount, false), nil
}

func (h *handlers) getControl(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, controlItem, error) {
	if in.ID == "" {
		return nil, controlItem{}, errIDRequired
	}
	resp, err := h.api.GetControlByID(ctx, in.ID)
	if err != nil {
		return nil, controlItem{}, openlane.APIError(err)
	}
	c := resp.Control
	return nil, controlItem{
		ID:                 c.ID,
		DisplayID:          c.DisplayID,
		RefCode:            c.RefCode,
		Title:              openlane.Deref(c.Title),
		Status:             openlane.Format(c.Status),
		ReferenceFramework: openlane.Deref(c.ReferenceFramework),
		Category:           openlane.Deref(c.Category),
		Description:        openlane.Deref(c.Description),
		Subcategory:        openlane.Deref(c.Subcategory),
		StandardID:         openlane.Deref(c.StandardID),
		OwnerID:            openlane.Deref(c.OwnerID),
		Tags:               c.Tags,
	}, nil
}

func mapControlPage(edges []*graphclient.GetControls_Controls_Edges, hasMore bool, endCursor *string, total int64, _ bool) openlane.Page[controlItem] {
	items := make([]controlItem, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, controlItem{
			ID:                 n.ID,
			DisplayID:          n.DisplayID,
			RefCode:            n.RefCode,
			Title:              openlane.Deref(n.Title),
			Status:             openlane.Format(n.Status),
			ReferenceFramework: openlane.Deref(n.ReferenceFramework),
			Category:           openlane.Deref(n.Category),
		})
	}
	return openlane.Page[controlItem]{
		Items:      items,
		NextCursor: endCursor,
		HasMore:    hasMore,
		TotalCount: total,
	}
}
