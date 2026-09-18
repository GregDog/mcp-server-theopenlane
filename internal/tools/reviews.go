package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type reviewItem struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Details        string   `json:"details,omitempty"`
	Status         string   `json:"status,omitempty"`
	State          string   `json:"state,omitempty"`
	Category       string   `json:"category,omitempty"`
	Classification string   `json:"classification,omitempty"`
	Reporter       string   `json:"reporter,omitempty"`
	ReviewerID     string   `json:"reviewer_id,omitempty"`
	Approved       *bool    `json:"approved,omitempty"`
	ReviewedAt     string   `json:"reviewed_at,omitempty"`
	ReportedAt     string   `json:"reported_at,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

type reviewListInput struct {
	listInput
	EntityID string `json:"entity_id,omitempty" jsonschema:"Filter reviews linked to this entity (vendor) ID."`
	Title    string `json:"title,omitempty" jsonschema:"Filter reviews whose title contains this text (case-insensitive)."`
}

func registerReviews(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_reviews_list",
		Title:       "List Openlane reviews",
		Description: "List review records in the configured organization, including vendor Risk Review entries. Filter by entity (vendor) ID or title. Results are paginated.",
		Annotations: readOnly(),
	}, h.listReviews)

	addTool(server, &mcp.Tool{
		Name:        "openlane_review_get",
		Title:       "Get an Openlane review",
		Description: "Get a single review record by ID, including vendor Risk Review title and body (details).",
		Annotations: readOnly(),
	}, h.getReview)
}

func (h *handlers) listReviews(ctx context.Context, _ *mcp.CallToolRequest, in reviewListInput) (*mcp.CallToolResult, openlane.Page[reviewItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetReviews(ctx, &first, after, buildReviewWhere(in))
	if err != nil {
		return nil, openlane.Page[reviewItem]{}, openlane.APIError(err)
	}
	items := make([]reviewItem, 0, len(resp.Reviews.Edges))
	for _, e := range resp.Reviews.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListReview(e.Node))
	}
	return nil, openlane.Page[reviewItem]{
		Items:      items,
		NextCursor: resp.Reviews.PageInfo.EndCursor,
		HasMore:    resp.Reviews.PageInfo.HasNextPage,
		TotalCount: resp.Reviews.TotalCount,
	}, nil
}

func (h *handlers) getReview(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, reviewItem, error) {
	if in.ID == "" {
		return nil, reviewItem{}, errIDRequired
	}
	resp, err := h.api.GetReviewByID(ctx, in.ID)
	if err != nil {
		return nil, reviewItem{}, openlane.APIError(err)
	}
	return nil, mapGetReview(resp.Review), nil
}

func mapListReview(n *graphclient.GetReviews_Reviews_Edges_Node) reviewItem {
	return reviewItem{
		ID:       n.ID,
		Title:    n.Title,
		Summary:  openlane.Deref(n.Summary),
		Details:  openlane.Deref(n.Details),
		State:    openlane.Deref(n.State),
		Category: openlane.Deref(n.Category),
		Tags:     n.Tags,
	}
}

func mapGetReview(r graphclient.GetReviewByID_Review) reviewItem {
	return reviewItem{
		ID:             r.ID,
		Title:          r.Title,
		Summary:        openlane.Deref(r.Summary),
		Details:        openlane.Deref(r.Details),
		State:          openlane.Deref(r.State),
		Category:       openlane.Deref(r.Category),
		Classification: openlane.Deref(r.Classification),
		Reporter:       openlane.Deref(r.Reporter),
		ReviewerID:     openlane.Deref(r.ReviewerID),
		Approved:       r.Approved,
		ReviewedAt:     openlane.Format(r.ReviewedAt),
		ReportedAt:     openlane.Format(r.ReportedAt),
		Tags:           r.Tags,
	}
}

func mapCreatedReview(r graphclient.CreateReview_CreateReview_Review) reviewItem {
	return reviewItem{
		ID:             r.ID,
		Title:          r.Title,
		Summary:        openlane.Deref(r.Summary),
		Details:        openlane.Deref(r.Details),
		State:          openlane.Deref(r.State),
		Category:       openlane.Deref(r.Category),
		Classification: openlane.Deref(r.Classification),
		Reporter:       openlane.Deref(r.Reporter),
		ReviewerID:     openlane.Deref(r.ReviewerID),
		Approved:       r.Approved,
		ReviewedAt:     openlane.Format(r.ReviewedAt),
		ReportedAt:     openlane.Format(r.ReportedAt),
		Tags:           r.Tags,
	}
}

func mapUpdatedReview(r graphclient.UpdateReview_UpdateReview_Review) reviewItem {
	return reviewItem{
		ID:             r.ID,
		Title:          r.Title,
		Summary:        openlane.Deref(r.Summary),
		Details:        openlane.Deref(r.Details),
		State:          openlane.Deref(r.State),
		Category:       openlane.Deref(r.Category),
		Classification: openlane.Deref(r.Classification),
		Reporter:       openlane.Deref(r.Reporter),
		ReviewerID:     openlane.Deref(r.ReviewerID),
		Approved:       r.Approved,
		ReviewedAt:     openlane.Format(r.ReviewedAt),
		ReportedAt:     openlane.Format(r.ReportedAt),
		Tags:           r.Tags,
	}
}
