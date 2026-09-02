package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type assessmentItem struct {
	ID             string                             `json:"id"`
	Name           string                             `json:"name,omitempty"`
	AssessmentType string                             `json:"assessment_type,omitempty"`
	TemplateID     string                             `json:"template_id,omitempty"`
	OwnerID        string                             `json:"owner_id,omitempty"`
	Tags           []string                           `json:"tags,omitempty"`
	Campaigns      *relSummary[campaignRef]           `json:"campaigns,omitempty"`
	Responses      *relSummary[assessmentResponseRef] `json:"responses,omitempty"`
	Findings       *relSummary[findingRef]            `json:"findings,omitempty"`
}

func registerAssessments(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_assessments_list",
		Title:       "List Openlane assessments",
		Description: "List assessments in the configured Openlane organization. Optional assessment_type filter. Results are paginated.",
		Annotations: readOnly(),
	}, h.listAssessments)

	addTool(server, &mcp.Tool{
		Name:        "openlane_assessment_get",
		Title:       "Get an Openlane assessment",
		Description: "Get a single assessment by ID with campaign summaries and bounded response and finding summaries.",
		Annotations: readOnly(),
	}, h.getAssessment)
}

func (h *handlers) listAssessments(ctx context.Context, _ *mcp.CallToolRequest, in assessmentListInput) (*mcp.CallToolResult, openlane.Page[assessmentItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetAssessments(ctx, &first, after, buildAssessmentWhere(in))
	if err != nil {
		return nil, openlane.Page[assessmentItem]{}, openlane.APIError(err)
	}
	items := make([]assessmentItem, 0, len(resp.Assessments.Edges))
	for _, e := range resp.Assessments.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListAssessment(e.Node))
	}
	return nil, openlane.Page[assessmentItem]{
		Items:      items,
		NextCursor: resp.Assessments.PageInfo.EndCursor,
		HasMore:    resp.Assessments.PageInfo.HasNextPage,
		TotalCount: resp.Assessments.TotalCount,
	}, nil
}

func (h *handlers) getAssessment(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, assessmentItem, error) {
	if in.ID == "" {
		return nil, assessmentItem{}, errIDRequired
	}
	resp, err := h.api.GetAssessmentByID(ctx, in.ID)
	if err != nil {
		return nil, assessmentItem{}, openlane.APIError(err)
	}
	item := mapGetAssessment(resp.Assessment)
	runRelJobs(
		func() { item.Responses = h.fetchAssessmentResponses(ctx, in.ID) },
		func() { item.Findings = h.fetchFindings(ctx, &graphclient.FindingWhereInput{AssessmentID: &in.ID}) },
	)
	return nil, item, nil
}

func mapListAssessment(n *graphclient.GetAssessments_Assessments_Edges_Node) assessmentItem {
	return assessmentItem{
		ID:             n.ID,
		Name:           n.Name,
		AssessmentType: openlane.Format(n.AssessmentType),
		TemplateID:     openlane.Deref(n.TemplateID),
		OwnerID:        openlane.Deref(n.OwnerID),
		Tags:           n.Tags,
		Campaigns:      campaignsFromListNode(n.Campaigns),
	}
}

func mapGetAssessment(a graphclient.GetAssessmentByID_Assessment) assessmentItem {
	return assessmentItem{
		ID:             a.ID,
		Name:           a.Name,
		AssessmentType: openlane.Format(a.AssessmentType),
		TemplateID:     openlane.Deref(a.TemplateID),
		OwnerID:        openlane.Deref(a.OwnerID),
		Tags:           a.Tags,
		Campaigns:      campaignsFromGetNode(a.Campaigns),
	}
}

func campaignsFromListNode(c graphclient.GetAssessments_Assessments_Edges_Node_Campaigns) *relSummary[campaignRef] {
	if len(c.Edges) == 0 {
		return nil
	}
	items := make([]campaignRef, 0, len(c.Edges))
	for _, e := range c.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, campaignRef{ID: n.ID, DisplayID: n.DisplayID, Name: n.Name, Status: openlane.Format(n.Status)})
	}
	return &relSummary[campaignRef]{Count: int64(len(items)), Items: items}
}

func campaignsFromGetNode(c graphclient.GetAssessmentByID_Assessment_Campaigns) *relSummary[campaignRef] {
	if len(c.Edges) == 0 {
		return nil
	}
	items := make([]campaignRef, 0, len(c.Edges))
	for _, e := range c.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, campaignRef{ID: n.ID, DisplayID: n.DisplayID, Name: n.Name, Status: openlane.Format(n.Status)})
	}
	return &relSummary[campaignRef]{Count: int64(len(items)), Items: items}
}
