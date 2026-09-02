package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createPolicyInput struct {
	Name    string   `json:"name" jsonschema:"Policy name."`
	Status  string   `json:"status,omitempty" jsonschema:"Document status enum value."`
	Details string   `json:"details,omitempty" jsonschema:"Policy details."`
	Tags    []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
}

type updatePolicyInput struct {
	ID      string   `json:"id" jsonschema:"Policy ID to update."`
	Name    string   `json:"name,omitempty" jsonschema:"Updated name."`
	Status  string   `json:"status,omitempty" jsonschema:"Updated document status enum value."`
	Details string   `json:"details,omitempty" jsonschema:"Updated details."`
	Tags    []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
}

func registerWritePolicies(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_create",
		Title:       "Create an Openlane policy",
		Description: "Create an internal policy in Openlane. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createPolicy)

	addTool(server, &mcp.Tool{
		Name:        "openlane_policy_update",
		Title:       "Update an Openlane policy",
		Description: "Update an internal policy by ID. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updatePolicy)
}

func (h *handlers) createPolicy(ctx context.Context, _ *mcp.CallToolRequest, in createPolicyInput) (*mcp.CallToolResult, policyItem, error) {
	if in.Name == "" {
		return nil, policyItem{}, errNameRequired
	}
	input := graphclient.CreateInternalPolicyInput{
		Name: in.Name,
		Tags: in.Tags,
	}
	if in.Status != "" {
		input.Status = documentStatus(in.Status)
	}
	if in.Details != "" {
		input.Details = &in.Details
	}

	resp, err := h.api.CreateInternalPolicy(ctx, input)
	if err != nil {
		return nil, policyItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedPolicy(resp.CreateInternalPolicy.InternalPolicy), nil
}

func (h *handlers) updatePolicy(ctx context.Context, _ *mcp.CallToolRequest, in updatePolicyInput) (*mcp.CallToolResult, policyItem, error) {
	if in.ID == "" {
		return nil, policyItem{}, errIDRequired
	}
	input := graphclient.UpdateInternalPolicyInput{}
	if in.Name != "" {
		input.Name = &in.Name
	}
	if in.Status != "" {
		input.Status = documentStatus(in.Status)
	}
	if in.Details != "" {
		input.Details = &in.Details
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if isEmptyUpdatePolicy(input) {
		return nil, policyItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateInternalPolicy(ctx, in.ID, input)
	if err != nil {
		return nil, policyItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedPolicy(resp.UpdateInternalPolicy.InternalPolicy), nil
}

func mapCreatedPolicy(p graphclient.CreateInternalPolicy_CreateInternalPolicy_InternalPolicy) policyItem {
	return policyItem{
		ID:        p.ID,
		DisplayID: p.DisplayID,
		Name:      p.Name,
		Status:    openlane.Format(p.Status),
		Summary:   openlane.Deref(p.Summary),
		Revision:  openlane.Deref(p.Revision),
		Kind:      openlane.Deref(p.InternalPolicyKindName),
		ReviewDue: openlane.Format(p.ReviewDue),
	}
}

func mapUpdatedPolicy(p graphclient.UpdateInternalPolicy_UpdateInternalPolicy_InternalPolicy) policyItem {
	return policyItem{
		ID:        p.ID,
		DisplayID: p.DisplayID,
		Name:      p.Name,
		Status:    openlane.Format(p.Status),
		Summary:   openlane.Deref(p.Summary),
		Revision:  openlane.Deref(p.Revision),
		ReviewDue: openlane.Format(p.ReviewDue),
	}
}

func isEmptyUpdatePolicy(in graphclient.UpdateInternalPolicyInput) bool {
	return in.Name == nil &&
		in.Status == nil &&
		in.Details == nil &&
		len(in.Tags) == 0
}
