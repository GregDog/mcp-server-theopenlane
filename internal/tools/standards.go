package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type standardItem struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	ShortName     string `json:"short_name,omitempty"`
	Framework     string `json:"framework,omitempty"`
	Version       string `json:"version,omitempty"`
	Status        string `json:"status,omitempty"`
	Description   string `json:"description,omitempty"`
	StandardType  string `json:"standard_type,omitempty"`
	GoverningBody string `json:"governing_body,omitempty"`
	ControlCount  *int64 `json:"control_count,omitempty"`
}

func registerStandards(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_standards_list",
		Title:       "List Openlane standards",
		Description: "List standards and frameworks available in Openlane. Results are paginated.",
		Annotations: readOnly(),
	}, h.listStandards)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_standard_get",
		Title:       "Get an Openlane standard",
		Description: "Get a single standard by ID.",
		Annotations: readOnly(),
	}, h.getStandard)
}

func (h *handlers) listStandards(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[standardItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetStandards(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[standardItem]{}, openlane.APIError(err)
	}
	items := make([]standardItem, 0, len(resp.Standards.Edges))
	for _, e := range resp.Standards.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, standardItem{
			ID:        n.ID,
			Name:      n.Name,
			ShortName: openlane.Deref(n.ShortName),
			Framework: openlane.Deref(n.Framework),
			Version:   openlane.Deref(n.Version),
			Status:    openlane.Format(n.Status),
		})
	}
	return nil, openlane.Page[standardItem]{
		Items:      items,
		NextCursor: resp.Standards.PageInfo.EndCursor,
		HasMore:    resp.Standards.PageInfo.HasNextPage,
		TotalCount: resp.Standards.TotalCount,
	}, nil
}

func (h *handlers) getStandard(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, standardItem, error) {
	if in.ID == "" {
		return nil, standardItem{}, errIDRequired
	}
	resp, err := h.api.GetStandardByID(ctx, in.ID)
	if err != nil {
		return nil, standardItem{}, openlane.APIError(err)
	}
	s := resp.Standard
	count := s.Controls.TotalCount
	return nil, standardItem{
		ID:            s.ID,
		Name:          s.Name,
		ShortName:     openlane.Deref(s.ShortName),
		Framework:     openlane.Deref(s.Framework),
		Version:       openlane.Deref(s.Version),
		Status:        openlane.Format(s.Status),
		Description:   openlane.Deref(s.Description),
		StandardType:  openlane.Deref(s.StandardType),
		GoverningBody: openlane.Deref(s.GoverningBody),
		ControlCount:  &count,
	}, nil
}
