package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createRiskInput struct {
	Name       string   `json:"name" jsonschema:"Risk name."`
	Status     string   `json:"status,omitempty" jsonschema:"Risk status enum value."`
	Impact     string   `json:"impact,omitempty" jsonschema:"Risk impact enum value."`
	Likelihood string   `json:"likelihood,omitempty" jsonschema:"Risk likelihood enum value."`
	Details    string   `json:"details,omitempty" jsonschema:"Risk details."`
	Mitigation string   `json:"mitigation,omitempty" jsonschema:"Risk mitigation."`
	Tags       []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
}

type updateRiskInput struct {
	ID         string   `json:"id" jsonschema:"Risk ID to update."`
	Name       string   `json:"name,omitempty" jsonschema:"Updated name."`
	Status     string   `json:"status,omitempty" jsonschema:"Updated risk status enum value."`
	Impact     string   `json:"impact,omitempty" jsonschema:"Updated impact enum value."`
	Likelihood string   `json:"likelihood,omitempty" jsonschema:"Updated likelihood enum value."`
	Details    string   `json:"details,omitempty" jsonschema:"Updated details."`
	Mitigation string   `json:"mitigation,omitempty" jsonschema:"Updated mitigation."`
	Tags       []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
}

func registerWriteRisks(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_risk_create",
		Title:       "Create an Openlane risk",
		Description: "Create a risk in Openlane. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createRisk)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_risk_update",
		Title:       "Update an Openlane risk",
		Description: "Update a risk by ID. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateRisk)
}

func (h *handlers) createRisk(ctx context.Context, _ *mcp.CallToolRequest, in createRiskInput) (*mcp.CallToolResult, riskItem, error) {
	if in.Name == "" {
		return nil, riskItem{}, errNameRequired
	}
	input := graphclient.CreateRiskInput{
		Name: in.Name,
		Tags: in.Tags,
	}
	if in.Status != "" {
		input.Status = riskStatus(in.Status)
	}
	if in.Impact != "" {
		input.Impact = riskImpact(in.Impact)
	}
	if in.Likelihood != "" {
		input.Likelihood = riskLikelihood(in.Likelihood)
	}
	if in.Details != "" {
		input.Details = &in.Details
	}
	if in.Mitigation != "" {
		input.Mitigation = &in.Mitigation
	}

	resp, err := h.api.CreateRisk(ctx, input)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedRisk(resp.CreateRisk.Risk), nil
}

func (h *handlers) updateRisk(ctx context.Context, _ *mcp.CallToolRequest, in updateRiskInput) (*mcp.CallToolResult, riskItem, error) {
	if in.ID == "" {
		return nil, riskItem{}, errIDRequired
	}
	input := graphclient.UpdateRiskInput{}
	if in.Name != "" {
		input.Name = &in.Name
	}
	if in.Status != "" {
		input.Status = riskStatus(in.Status)
	}
	if in.Impact != "" {
		input.Impact = riskImpact(in.Impact)
	}
	if in.Likelihood != "" {
		input.Likelihood = riskLikelihood(in.Likelihood)
	}
	if in.Details != "" {
		input.Details = &in.Details
	}
	if in.Mitigation != "" {
		input.Mitigation = &in.Mitigation
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if isEmptyUpdateRisk(input) {
		return nil, riskItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateRisk(ctx, in.ID, input)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedRisk(resp.UpdateRisk.Risk), nil
}

func mapCreatedRisk(r graphclient.CreateRisk_CreateRisk_Risk) riskItem {
	return riskItem{
		ID:         r.ID,
		DisplayID:  r.DisplayID,
		Name:       r.Name,
		Status:     openlane.Format(r.Status),
		Impact:     openlane.Format(r.Impact),
		Likelihood: openlane.Format(r.Likelihood),
		Score:      r.Score,
		Details:    openlane.Deref(r.Details),
		Mitigation: openlane.Deref(r.Mitigation),
	}
}

func mapUpdatedRisk(r graphclient.UpdateRisk_UpdateRisk_Risk) riskItem {
	return riskItem{
		ID:         r.ID,
		DisplayID:  r.DisplayID,
		Name:       r.Name,
		Status:     openlane.Format(r.Status),
		Impact:     openlane.Format(r.Impact),
		Likelihood: openlane.Format(r.Likelihood),
		Score:      r.Score,
		Details:    openlane.Deref(r.Details),
		Mitigation: openlane.Deref(r.Mitigation),
	}
}

func isEmptyUpdateRisk(in graphclient.UpdateRiskInput) bool {
	return in.Name == nil &&
		in.Status == nil &&
		in.Impact == nil &&
		in.Likelihood == nil &&
		in.Details == nil &&
		in.Mitigation == nil &&
		len(in.Tags) == 0
}
