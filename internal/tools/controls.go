package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type controlItem struct {
	ID                     string                         `json:"id"`
	DisplayID              string                         `json:"display_id,omitempty"`
	RefCode                string                         `json:"ref_code,omitempty"`
	Title                  string                         `json:"title,omitempty"`
	Status                 string                         `json:"status,omitempty"`
	ReferenceFramework     string                         `json:"reference_framework,omitempty"`
	Category               string                         `json:"category,omitempty"`
	Description            string                         `json:"description,omitempty"`
	Subcategory            string                         `json:"subcategory,omitempty"`
	StandardID             string                         `json:"standard_id,omitempty"`
	OwnerID                string                         `json:"owner_id,omitempty"`
	ControlOwnerID         string                         `json:"control_owner_id,omitempty"`
	ControlKindName        string                         `json:"control_kind_name,omitempty"`
	ImplementationGuidance any                            `json:"implementation_guidance,omitempty"`
	AssessmentMethods      any                            `json:"assessment_methods,omitempty"`
	AssessmentObjectives   any                            `json:"assessment_objectives,omitempty"`
	Tags                   []string                       `json:"tags,omitempty"`
	Programs               *relSummary[idNameRef]         `json:"programs,omitempty"`
	Evidence               *relSummary[idNameRef]         `json:"evidence,omitempty"`
	Findings               *relSummary[findingRef]        `json:"findings,omitempty"`
	Risks                  *relSummary[idNameRef]         `json:"risks,omitempty"`
	Implementations        *relSummary[implementationRef] `json:"implementations,omitempty"`
}

func registerControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_controls_list",
		Title:       "List Openlane controls",
		Description: "List controls in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listControls)

	addTool(server, &mcp.Tool{
		Name:        "openlane_controls_search",
		Title:       "Search Openlane controls",
		Description: "Search controls by ref code, title, or description using Openlane where-filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.searchControls)

	addTool(server, &mcp.Tool{
		Name:        "openlane_control_get",
		Title:       "Get an Openlane control",
		Description: "Get a single control by ID with assessment metadata and bounded summaries of linked programs, evidence, findings, risks, and implementations.",
		Annotations: readOnly(),
	}, h.getControl)
}

func (h *handlers) listControls(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControls(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[controlItem]{}, openlane.APIError(err)
	}
	return nil, mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount), nil
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
	return nil, mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount), nil
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
	item := controlItem{
		ID:                     c.ID,
		DisplayID:              c.DisplayID,
		RefCode:                c.RefCode,
		Title:                  openlane.Deref(c.Title),
		Status:                 openlane.Format(c.Status),
		ReferenceFramework:     openlane.Deref(c.ReferenceFramework),
		Category:               openlane.Deref(c.Category),
		Description:            openlane.Deref(c.Description),
		Subcategory:            openlane.Deref(c.Subcategory),
		StandardID:             openlane.Deref(c.StandardID),
		OwnerID:                openlane.Deref(c.OwnerID),
		ControlOwnerID:         openlane.Deref(c.ControlOwnerID),
		ControlKindName:        openlane.Deref(c.ControlKindName),
		ImplementationGuidance: c.ImplementationGuidance,
		AssessmentMethods:      c.AssessmentMethods,
		AssessmentObjectives:   c.AssessmentObjectives,
		Tags:                   c.Tags,
	}
	controlWhere := []*graphclient.ControlWhereInput{{ID: &in.ID}}
	runRelJobs(
		func() {
			item.Programs = h.fetchPrograms(ctx, &graphclient.ProgramWhereInput{HasControlsWith: controlWhere})
		},
		func() {
			item.Evidence = h.fetchEvidences(ctx, &graphclient.EvidenceWhereInput{HasControlsWith: controlWhere})
		},
		func() {
			item.Findings = h.fetchFindings(ctx, &graphclient.FindingWhereInput{HasControlsWith: controlWhere})
		},
		func() { item.Risks = h.fetchRisks(ctx, &graphclient.RiskWhereInput{HasControlsWith: controlWhere}) },
		func() {
			item.Implementations = h.fetchImplementations(ctx, &graphclient.ControlImplementationWhereInput{HasControlsWith: controlWhere})
		},
	)
	return nil, item, nil
}

func mapControlPage(edges []*graphclient.GetControls_Controls_Edges, hasMore bool, endCursor *string, total int64) openlane.Page[controlItem] {
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
