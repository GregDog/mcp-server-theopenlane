package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type evidenceControlRef struct {
	ID      string `json:"id"`
	RefCode string `json:"ref_code,omitempty"`
}

type evidenceItem struct {
	ID           string               `json:"id"`
	DisplayID    string               `json:"display_id,omitempty"`
	Name         string               `json:"name,omitempty"`
	Status       string               `json:"status,omitempty"`
	Source       string               `json:"source,omitempty"`
	CreationDate string               `json:"creation_date,omitempty"`
	Description  string               `json:"description,omitempty"`
	RenewalDate  string               `json:"renewal_date,omitempty"`
	IsAutomated  *bool                `json:"is_automated,omitempty"`
	Controls     []evidenceControlRef `json:"controls,omitempty"`
}

func registerEvidence(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_evidence_list",
		Title:       "List Openlane evidence",
		Description: "List evidence records in the configured Openlane organization. Results are paginated. File contents and download URLs are not included.",
		Annotations: readOnly(),
	}, h.listEvidence)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_evidence_get",
		Title:       "Get Openlane evidence",
		Description: "Get a single evidence record by ID. File contents and presigned URLs are not returned.",
		Annotations: readOnly(),
	}, h.getEvidence)
}

func (h *handlers) listEvidence(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[evidenceItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetEvidences(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[evidenceItem]{}, openlane.APIError(err)
	}
	items := make([]evidenceItem, 0, len(resp.Evidences.Edges))
	for _, e := range resp.Evidences.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, evidenceItem{
			ID:           n.ID,
			DisplayID:    n.DisplayID,
			Name:         n.Name,
			Status:       openlane.Format(n.Status),
			Source:       openlane.Deref(n.Source),
			CreationDate: openlane.Format(n.CreationDate),
		})
	}
	return nil, openlane.Page[evidenceItem]{
		Items:      items,
		NextCursor: resp.Evidences.PageInfo.EndCursor,
		HasMore:    resp.Evidences.PageInfo.HasNextPage,
		TotalCount: resp.Evidences.TotalCount,
	}, nil
}

func (h *handlers) getEvidence(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, evidenceItem, error) {
	if in.ID == "" {
		return nil, evidenceItem{}, errIDRequired
	}
	resp, err := h.api.GetEvidenceByID(ctx, in.ID)
	if err != nil {
		return nil, evidenceItem{}, openlane.APIError(err)
	}
	e := resp.Evidence
	var controls []evidenceControlRef
	for _, edge := range e.Controls.Edges {
		if edge == nil || edge.Node == nil {
			continue
		}
		controls = append(controls, evidenceControlRef{ID: edge.Node.ID, RefCode: edge.Node.RefCode})
	}
	return nil, evidenceItem{
		ID:           e.ID,
		DisplayID:    e.DisplayID,
		Name:         e.Name,
		Status:       openlane.Format(e.Status),
		Source:       openlane.Deref(e.Source),
		CreationDate: openlane.Format(e.CreationDate),
		Description:  openlane.Deref(e.Description),
		RenewalDate:  openlane.Format(e.RenewalDate),
		IsAutomated:  e.IsAutomated,
		Controls:     controls,
	}, nil
}
