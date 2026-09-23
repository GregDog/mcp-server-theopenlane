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
	Source                 string                         `json:"source,omitempty"`
	LinkableToEvidence     bool                           `json:"linkable_to_evidence"`
	ReferenceFramework     string                         `json:"reference_framework,omitempty"`
	Category               string                         `json:"category,omitempty"`
	Description            string                         `json:"description,omitempty"`
	Subcategory            string                         `json:"subcategory,omitempty"`
	StandardID             string                         `json:"standard_id,omitempty"`
	OwnerID                string                         `json:"owner_id,omitempty"`
	ControlOwnerID         string                         `json:"control_owner_id,omitempty"`
	ControlOwner           *groupAssigneeSummary          `json:"control_owner,omitempty"`
	DelegateID             string                         `json:"delegate_id,omitempty"`
	Delegate               *groupAssigneeSummary          `json:"delegate,omitempty"`
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
	RelatedControls        *relSummary[relatedControlRef] `json:"related_controls,omitempty"`
}

func registerControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_controls_list",
		Title:       "List Openlane controls",
		Description: "List controls in the configured Openlane organization. Use linkable_only to return org-owned controls suitable for evidence linking (excludes system catalog copies). Results are paginated.",
		Annotations: readOnly(),
	}, h.listControls)

	addTool(server, &mcp.Tool{
		Name:        "openlane_controls_search",
		Title:       "Search Openlane controls",
		Description: "Search controls by ref code, title, or description. Use linkable_only when resolving controls for evidence linking (org-owned program controls only). Results are paginated.",
		Annotations: readOnly(),
	}, h.searchControls)

	addTool(server, &mcp.Tool{
		Name:        "openlane_control_get",
		Title:       "Get an Openlane control",
		Description: "Get a single control by ID with assessment metadata and bounded summaries of linked programs, evidence, findings, risks, implementations, and cross-framework control mappings.",
		Annotations: readOnly(),
	}, h.getControl)
}

func (h *handlers) listControls(ctx context.Context, _ *mcp.CallToolRequest, in controlListInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	where := controlWhereLinkableOnly(in.LinkableOnly)
	resp, err := h.api.GetControls(ctx, &first, after, where)
	if err != nil {
		return nil, openlane.Page[controlItem]{}, openlane.APIError(err)
	}
	page := mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount)
	if err := enrichControlPageAssignees(ctx, h, page.Items); err != nil {
		return nil, openlane.Page[controlItem]{}, err
	}
	return nil, page, nil
}

func (h *handlers) searchControls(ctx context.Context, _ *mcp.CallToolRequest, in controlSearchInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	if in.Query == "" {
		return nil, openlane.Page[controlItem]{}, errQueryRequired
	}
	q := in.Query
	where := mergeControlWhere(&graphclient.ControlWhereInput{
		Or: []*graphclient.ControlWhereInput{
			{RefCodeContainsFold: &q},
			{TitleContainsFold: &q},
			{DescriptionContainsFold: &q},
		},
	}, controlWhereLinkableOnly(in.LinkableOnly))
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControls(ctx, &first, after, where)
	if err != nil {
		return nil, openlane.Page[controlItem]{}, openlane.APIError(err)
	}
	page := mapControlPage(resp.Controls.Edges, resp.Controls.PageInfo.HasNextPage, resp.Controls.PageInfo.EndCursor, resp.Controls.TotalCount)
	if err := enrichControlPageAssignees(ctx, h, page.Items); err != nil {
		return nil, openlane.Page[controlItem]{}, err
	}
	return nil, page, nil
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
	item := mapControlDetail(c)
	if err := h.enrichControlAssignees(ctx, &item, nil); err != nil {
		return nil, controlItem{}, err
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
		func() { item.RelatedControls = h.fetchRelatedControls(ctx, in.ID) },
	)
	return nil, item, nil
}

func mapControlPage(edges []*graphclient.GetControls_Controls_Edges, hasMore bool, endCursor *string, total int64) openlane.Page[controlItem] {
	items := make([]controlItem, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapControlListNode(e.Node))
	}
	return openlane.Page[controlItem]{
		Items:      items,
		NextCursor: endCursor,
		HasMore:    hasMore,
		TotalCount: total,
	}
}
