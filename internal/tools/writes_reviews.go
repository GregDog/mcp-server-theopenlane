package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createVendorRiskReviewInput struct {
	Title     string   `json:"title" jsonschema:"Risk Review title shown on the vendor record."`
	Body      string   `json:"body,omitempty" jsonschema:"Risk Review body text (stored as review details)."`
	Summary   string   `json:"summary,omitempty" jsonschema:"Optional short summary for the review."`
	EntityIDs []string `json:"entity_ids" jsonschema:"Entity (vendor) IDs to attach this Risk Review to. At least one required."`
	Category  string   `json:"category,omitempty" jsonschema:"Optional review category."`
	Status    string   `json:"status,omitempty" jsonschema:"Optional review status enum value."`
	State     string   `json:"state,omitempty" jsonschema:"Optional review state label."`
	Tags      []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
	FileIDs   []string `json:"file_ids,omitempty" jsonschema:"Optional file IDs to attach to the review (for example uploaded evidence artifacts)."`
}

type updateVendorRiskReviewInput struct {
	ID              string   `json:"id" jsonschema:"Review ID to update."`
	Title           string   `json:"title,omitempty" jsonschema:"Updated Risk Review title."`
	Body            string   `json:"body,omitempty" jsonschema:"Updated Risk Review body text (stored as review details)."`
	Summary         string   `json:"summary,omitempty" jsonschema:"Updated summary."`
	Category        string   `json:"category,omitempty" jsonschema:"Updated category."`
	Status          string   `json:"status,omitempty" jsonschema:"Updated review status enum value."`
	State           string   `json:"state,omitempty" jsonschema:"Updated review state label."`
	Tags            []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	AddEntityIDs    []string `json:"add_entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to link to this review."`
	RemoveEntityIDs []string `json:"remove_entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to unlink from this review."`
	AddFileIDs      []string `json:"add_file_ids,omitempty" jsonschema:"File IDs to attach to the review."`
	RemoveFileIDs   []string `json:"remove_file_ids,omitempty" jsonschema:"File IDs to detach from the review."`
}

func registerWriteReviews(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_vendor_risk_review_create",
		Title:       "Create a vendor Risk Review",
		Description: "Create a vendor Risk Review record (Openlane Review object) with title and body, linked to one or more entities (vendors). Use this instead of standalone risks for vendor onboarding assessments. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createVendorRiskReview)

	addTool(server, &mcp.Tool{
		Name:        "openlane_vendor_risk_review_update",
		Title:       "Update a vendor Risk Review",
		Description: "Update a vendor Risk Review record by ID. Supports title, body (details), summary, and entity/file links. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateVendorRiskReview)
}

func (h *handlers) createVendorRiskReview(ctx context.Context, _ *mcp.CallToolRequest, in createVendorRiskReviewInput) (*mcp.CallToolResult, reviewItem, error) {
	if strings.TrimSpace(in.Title) == "" {
		return nil, reviewItem{}, errTitleRequired
	}
	if len(in.EntityIDs) == 0 {
		return nil, reviewItem{}, errEntityIDRequired
	}

	input := graphclient.CreateReviewInput{
		Title:     strings.TrimSpace(in.Title),
		Tags:      in.Tags,
		EntityIDs: in.EntityIDs,
		FileIDs:   in.FileIDs,
	}
	if s := strings.TrimSpace(in.Body); s != "" {
		input.Details = &s
	}
	if s := strings.TrimSpace(in.Summary); s != "" {
		input.Summary = &s
	}
	if s := strings.TrimSpace(in.Category); s != "" {
		input.Category = &s
	}
	if in.Status != "" {
		input.Status = reviewStatus(in.Status)
	}
	if s := strings.TrimSpace(in.State); s != "" {
		input.State = &s
	}

	resp, err := h.api.CreateReview(ctx, input)
	if err != nil {
		return nil, reviewItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedReview(resp.CreateReview.Review), nil
}

func (h *handlers) updateVendorRiskReview(ctx context.Context, _ *mcp.CallToolRequest, in updateVendorRiskReviewInput) (*mcp.CallToolResult, reviewItem, error) {
	if in.ID == "" {
		return nil, reviewItem{}, errIDRequired
	}

	input := graphclient.UpdateReviewInput{}
	if s := strings.TrimSpace(in.Title); s != "" {
		input.Title = &s
	}
	if s := strings.TrimSpace(in.Body); s != "" {
		input.Details = &s
	}
	if s := strings.TrimSpace(in.Summary); s != "" {
		input.Summary = &s
	}
	if s := strings.TrimSpace(in.Category); s != "" {
		input.Category = &s
	}
	if in.Status != "" {
		input.Status = reviewStatus(in.Status)
	}
	if s := strings.TrimSpace(in.State); s != "" {
		input.State = &s
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if len(in.AddEntityIDs) > 0 {
		input.AddEntityIDs = in.AddEntityIDs
	}
	if len(in.RemoveEntityIDs) > 0 {
		input.RemoveEntityIDs = in.RemoveEntityIDs
	}
	if len(in.AddFileIDs) > 0 {
		input.AddFileIDs = in.AddFileIDs
	}
	if len(in.RemoveFileIDs) > 0 {
		input.RemoveFileIDs = in.RemoveFileIDs
	}
	if isEmptyUpdateReview(input) {
		return nil, reviewItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateReview(ctx, in.ID, input)
	if err != nil {
		return nil, reviewItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedReview(resp.UpdateReview.Review), nil
}

func isEmptyUpdateReview(in graphclient.UpdateReviewInput) bool {
	return in.Title == nil &&
		in.Details == nil &&
		in.Summary == nil &&
		in.Category == nil &&
		in.Status == nil &&
		in.State == nil &&
		len(in.Tags) == 0 &&
		len(in.AddEntityIDs) == 0 &&
		len(in.RemoveEntityIDs) == 0 &&
		len(in.AddFileIDs) == 0 &&
		len(in.RemoveFileIDs) == 0
}
