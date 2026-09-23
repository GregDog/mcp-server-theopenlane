package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createRiskInput struct {
	Name          string   `json:"name" jsonschema:"Risk name."`
	Status        string   `json:"status,omitempty" jsonschema:"Risk status enum value."`
	Impact        string   `json:"impact,omitempty" jsonschema:"Risk impact enum value."`
	Likelihood    string   `json:"likelihood,omitempty" jsonschema:"Risk likelihood enum value."`
	Details       string   `json:"details,omitempty" jsonschema:"Risk details."`
	Mitigation    string   `json:"mitigation,omitempty" jsonschema:"Risk mitigation."`
	Tags          []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
	EntityIDs     []string `json:"entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to associate with this risk on create."`
	StakeholderID string   `json:"stakeholder_id,omitempty" jsonschema:"Risk stakeholder: user ID, email, name, or group ID. User identifiers resolve to the member's managed personal group (Openlane stakeholderID)."`
	DelegateID    string   `json:"delegate_id,omitempty" jsonschema:"Risk delegate: user ID, email, name, or group ID. User identifiers resolve to the member's managed personal group (Openlane delegateID)."`
}

type updateRiskInput struct {
	ID              string   `json:"id" jsonschema:"Risk ID to update."`
	Name            string   `json:"name,omitempty" jsonschema:"Updated name."`
	Status          string   `json:"status,omitempty" jsonschema:"Updated risk status enum value."`
	Impact          string   `json:"impact,omitempty" jsonschema:"Updated impact enum value."`
	Likelihood      string   `json:"likelihood,omitempty" jsonschema:"Updated likelihood enum value."`
	Details         string   `json:"details,omitempty" jsonschema:"Updated details."`
	Mitigation      string   `json:"mitigation,omitempty" jsonschema:"Updated mitigation."`
	Tags            []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	AddEntityIDs    []string `json:"add_entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to associate with this risk."`
	RemoveEntityIDs []string `json:"remove_entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to unlink from this risk without deleting the risk record."`
	StakeholderID   string   `json:"stakeholder_id,omitempty" jsonschema:"Updated risk stakeholder: user ID, email, name, or group ID. User identifiers resolve to the member's managed personal group (Openlane stakeholderID)."`
	DelegateID      string   `json:"delegate_id,omitempty" jsonschema:"Updated risk delegate: user ID, email, name, or group ID. User identifiers resolve to the member's managed personal group (Openlane delegateID)."`
}

func registerWriteRisks(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_risk_create",
		Title:       "Create an Openlane risk",
		Description: "Create a risk in Openlane. Optional entity_ids links the risk to vendors. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createRisk)

	addTool(server, &mcp.Tool{
		Name:        "openlane_risk_update",
		Title:       "Update an Openlane risk",
		Description: "Update a risk by ID. Use add_entity_ids to link vendors or remove_entity_ids to unlink without deleting the risk. Requires write mode.",
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
	if len(in.EntityIDs) > 0 {
		input.EntityIDs = in.EntityIDs
	}
	if in.StakeholderID != "" {
		stakeholderGroupID, err := h.resolveGroupAssigneeGroupID(ctx, in.StakeholderID)
		if err != nil {
			return nil, riskItem{}, err
		}
		input.StakeholderID = &stakeholderGroupID
	}
	if in.DelegateID != "" {
		delegateGroupID, err := h.resolveGroupAssigneeGroupID(ctx, in.DelegateID)
		if err != nil {
			return nil, riskItem{}, err
		}
		input.DelegateID = &delegateGroupID
	}

	resp, err := h.api.CreateRisk(ctx, input)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	item := mapCreatedRisk(resp.CreateRisk.Risk)
	if err := h.enrichRiskAssignees(ctx, &item, nil); err != nil {
		return nil, riskItem{}, err
	}
	return nil, item, nil
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
	if len(in.AddEntityIDs) > 0 {
		input.AddEntityIDs = in.AddEntityIDs
	}
	if len(in.RemoveEntityIDs) > 0 {
		input.RemoveEntityIDs = in.RemoveEntityIDs
	}
	if in.StakeholderID != "" {
		stakeholderGroupID, err := h.resolveGroupAssigneeGroupID(ctx, in.StakeholderID)
		if err != nil {
			return nil, riskItem{}, err
		}
		input.StakeholderID = &stakeholderGroupID
	}
	if in.DelegateID != "" {
		delegateGroupID, err := h.resolveGroupAssigneeGroupID(ctx, in.DelegateID)
		if err != nil {
			return nil, riskItem{}, err
		}
		input.DelegateID = &delegateGroupID
	}
	if isEmptyUpdateRisk(input) {
		return nil, riskItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateRisk(ctx, in.ID, input)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	item := mapUpdatedRisk(resp.UpdateRisk.Risk)
	if err := h.enrichRiskAssignees(ctx, &item, nil); err != nil {
		return nil, riskItem{}, err
	}
	return nil, item, nil
}

func mapCreatedRisk(r graphclient.CreateRisk_CreateRisk_Risk) riskItem {
	return riskItem{
		ID:            r.ID,
		DisplayID:     r.DisplayID,
		Name:          r.Name,
		Status:        openlane.Format(r.Status),
		Impact:        openlane.Format(r.Impact),
		Likelihood:    openlane.Format(r.Likelihood),
		Score:         r.Score,
		Details:       openlane.Deref(r.Details),
		Mitigation:    openlane.Deref(r.Mitigation),
		StakeholderID: openlane.Deref(r.StakeholderID),
		DelegateID:    openlane.Deref(r.DelegateID),
	}
}

func mapUpdatedRisk(r graphclient.UpdateRisk_UpdateRisk_Risk) riskItem {
	return riskItem{
		ID:            r.ID,
		DisplayID:     r.DisplayID,
		Name:          r.Name,
		Status:        openlane.Format(r.Status),
		Impact:        openlane.Format(r.Impact),
		Likelihood:    openlane.Format(r.Likelihood),
		Score:         r.Score,
		Details:       openlane.Deref(r.Details),
		Mitigation:    openlane.Deref(r.Mitigation),
		StakeholderID: openlane.Deref(r.StakeholderID),
		DelegateID:    openlane.Deref(r.DelegateID),
	}
}

func isEmptyUpdateRisk(in graphclient.UpdateRiskInput) bool {
	return in.Name == nil &&
		in.Status == nil &&
		in.Impact == nil &&
		in.Likelihood == nil &&
		in.Details == nil &&
		in.Mitigation == nil &&
		len(in.Tags) == 0 &&
		len(in.AddEntityIDs) == 0 &&
		len(in.RemoveEntityIDs) == 0 &&
		in.StakeholderID == nil &&
		in.DelegateID == nil
}
