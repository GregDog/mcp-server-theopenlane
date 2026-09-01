package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type riskItem struct {
	ID         string `json:"id"`
	DisplayID  string `json:"display_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Status     string `json:"status,omitempty"`
	Impact     string `json:"impact,omitempty"`
	Likelihood string `json:"likelihood,omitempty"`
	Score      *int64 `json:"score,omitempty"`
	Details    string `json:"details,omitempty"`
	Mitigation string `json:"mitigation,omitempty"`
}

func registerRisks(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_risks_list",
		Title:       "List Openlane risks",
		Description: "List risks in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listRisks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_risk_get",
		Title:       "Get an Openlane risk",
		Description: "Get a single risk by ID.",
		Annotations: readOnly(),
	}, h.getRisk)
}

func (h *handlers) listRisks(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[riskItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetRisks(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[riskItem]{}, openlane.APIError(err)
	}
	items := make([]riskItem, 0, len(resp.Risks.Edges))
	for _, e := range resp.Risks.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, riskItem{
			ID:         n.ID,
			DisplayID:  n.DisplayID,
			Name:       n.Name,
			Status:     openlane.Format(n.Status),
			Impact:     openlane.Format(n.Impact),
			Likelihood: openlane.Format(n.Likelihood),
			Score:      n.Score,
		})
	}
	return nil, openlane.Page[riskItem]{
		Items:      items,
		NextCursor: resp.Risks.PageInfo.EndCursor,
		HasMore:    resp.Risks.PageInfo.HasNextPage,
		TotalCount: resp.Risks.TotalCount,
	}, nil
}

func (h *handlers) getRisk(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, riskItem, error) {
	if in.ID == "" {
		return nil, riskItem{}, errIDRequired
	}
	resp, err := h.api.GetRiskByID(ctx, in.ID)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	r := resp.Risk
	return nil, riskItem{
		ID:         r.ID,
		DisplayID:  r.DisplayID,
		Name:       r.Name,
		Status:     openlane.Format(r.Status),
		Impact:     openlane.Format(r.Impact),
		Likelihood: openlane.Format(r.Likelihood),
		Score:      r.Score,
		Details:    openlane.Deref(r.Details),
		Mitigation: openlane.Deref(r.Mitigation),
	}, nil
}
