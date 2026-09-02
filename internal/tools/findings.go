package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type findingItem struct {
	ID              string                        `json:"id"`
	Name            string                        `json:"name,omitempty"`
	Description     string                        `json:"description,omitempty"`
	Open            *bool                         `json:"open,omitempty"`
	Status          string                        `json:"status,omitempty"`
	Severity        string                        `json:"severity,omitempty"`
	Impact          *float64                      `json:"impact,omitempty"`
	AssessmentID    string                        `json:"assessment_id,omitempty"`
	Recommendation  string                        `json:"recommendation,omitempty"`
	Remediations    *relSummary[remediationRef]   `json:"remediations,omitempty"`
	Vulnerabilities *relSummary[vulnerabilityRef] `json:"vulnerabilities,omitempty"`
	Controls        *relSummary[controlRef]       `json:"controls,omitempty"`
}

func registerFindings(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_findings_list",
		Title:       "List Openlane findings",
		Description: "List findings in the configured Openlane organization. Supports program, assessment, open, status, and severity filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.listFindings)

	addTool(server, &mcp.Tool{
		Name:        "openlane_finding_get",
		Title:       "Get an Openlane finding",
		Description: "Get a single finding by ID with description, impact, remediations, vulnerabilities, and a bounded control summary.",
		Annotations: readOnly(),
	}, h.getFinding)
}

func (h *handlers) listFindings(ctx context.Context, _ *mcp.CallToolRequest, in findingListInput) (*mcp.CallToolResult, openlane.Page[findingItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetFindings(ctx, &first, after, buildFindingWhere(in))
	if err != nil {
		return nil, openlane.Page[findingItem]{}, openlane.APIError(err)
	}
	items := make([]findingItem, 0, len(resp.Findings.Edges))
	for _, e := range resp.Findings.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListFinding(e.Node))
	}
	return nil, openlane.Page[findingItem]{
		Items:      items,
		NextCursor: resp.Findings.PageInfo.EndCursor,
		HasMore:    resp.Findings.PageInfo.HasNextPage,
		TotalCount: resp.Findings.TotalCount,
	}, nil
}

func (h *handlers) getFinding(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, findingItem, error) {
	if in.ID == "" {
		return nil, findingItem{}, errIDRequired
	}
	resp, err := h.api.GetFindingByID(ctx, in.ID)
	if err != nil {
		return nil, findingItem{}, openlane.APIError(err)
	}
	item := mapGetFinding(resp.Finding)
	findingWhere := []*graphclient.FindingWhereInput{{ID: &in.ID}}
	item.Controls = h.fetchControls(ctx, &graphclient.ControlWhereInput{HasFindingsWith: findingWhere})
	return nil, item, nil
}

func mapListFinding(n *graphclient.GetFindings_Findings_Edges_Node) findingItem {
	return findingItem{
		ID:           n.ID,
		Name:         openlane.Deref(n.DisplayName),
		Open:         n.Open,
		Status:       openlane.Deref(n.FindingStatusName),
		Severity:     openlane.Deref(n.Severity),
		AssessmentID: openlane.Deref(n.AssessmentID),
	}
}

func mapGetFinding(f graphclient.GetFindingByID_Finding) findingItem {
	return findingItem{
		ID:              f.ID,
		Name:            openlane.Deref(f.DisplayName),
		Description:     openlane.Deref(f.Description),
		Open:            f.Open,
		Status:          openlane.Deref(f.FindingStatusName),
		Severity:        openlane.Deref(f.Severity),
		Impact:          f.Impact,
		AssessmentID:    openlane.Deref(f.AssessmentID),
		Recommendation:  openlane.Deref(f.Recommendation),
		Remediations:    remediationsFromGetFinding(f.Remediations),
		Vulnerabilities: vulnerabilitiesFromGetFinding(f.Vulnerabilities),
	}
}

func remediationsFromGetFinding(r graphclient.GetFindingByID_Finding_Remediations) *relSummary[remediationRef] {
	if len(r.Edges) == 0 {
		return nil
	}
	items := make([]remediationRef, 0, len(r.Edges))
	for _, e := range r.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, remediationRef{
			ID:     n.ID,
			Intent: openlane.Deref(n.Intent),
			DueAt:  openlane.Format(n.DueAt),
		})
	}
	return &relSummary[remediationRef]{Count: int64(len(items)), Items: items}
}

func vulnerabilitiesFromGetFinding(v graphclient.GetFindingByID_Finding_Vulnerabilities) *relSummary[vulnerabilityRef] {
	if len(v.Edges) == 0 {
		return nil
	}
	items := make([]vulnerabilityRef, 0, len(v.Edges))
	for _, e := range v.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, vulnerabilityRef{
			ID:          n.ID,
			DisplayName: openlane.Deref(n.DisplayName),
			CveID:       openlane.Deref(n.CveID),
			Severity:    openlane.Deref(n.Severity),
		})
	}
	return &relSummary[vulnerabilityRef]{Count: int64(len(items)), Items: items}
}
