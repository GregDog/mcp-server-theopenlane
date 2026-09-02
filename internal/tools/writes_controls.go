package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createControlInput struct {
	RefCode            string   `json:"ref_code" jsonschema:"Control reference code."`
	Title              string   `json:"title,omitempty" jsonschema:"Control title."`
	Description        string   `json:"description,omitempty" jsonschema:"Control description."`
	Status             string   `json:"status,omitempty" jsonschema:"Openlane control status enum value."`
	ReferenceFramework string   `json:"reference_framework,omitempty" jsonschema:"Reference framework name."`
	Category           string   `json:"category,omitempty" jsonschema:"Control category."`
	Subcategory        string   `json:"subcategory,omitempty" jsonschema:"Control subcategory."`
	StandardID         string   `json:"standard_id,omitempty" jsonschema:"Associated standard ID."`
	OwnerID            string   `json:"owner_id,omitempty" jsonschema:"Control owner user ID."`
	Tags               []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
}

type updateControlInput struct {
	ID          string   `json:"id" jsonschema:"Control ID to update."`
	Title       string   `json:"title,omitempty" jsonschema:"Updated title."`
	Description string   `json:"description,omitempty" jsonschema:"Updated description."`
	Status      string   `json:"status,omitempty" jsonschema:"Updated control status enum value."`
	Category    string   `json:"category,omitempty" jsonschema:"Updated category."`
	Subcategory string   `json:"subcategory,omitempty" jsonschema:"Updated subcategory."`
	StandardID  string   `json:"standard_id,omitempty" jsonschema:"Updated standard ID."`
	OwnerID     string   `json:"owner_id,omitempty" jsonschema:"Updated control owner user ID."`
	Tags        []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
}

func registerWriteControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_control_create",
		Title:       "Create an Openlane control",
		Description: "Create a control in the configured Openlane organization. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createControl)

	addTool(server, &mcp.Tool{
		Name:        "openlane_control_update",
		Title:       "Update an Openlane control",
		Description: "Update a control by ID. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateControl)
}

func (h *handlers) createControl(ctx context.Context, _ *mcp.CallToolRequest, in createControlInput) (*mcp.CallToolResult, controlItem, error) {
	if in.RefCode == "" {
		return nil, controlItem{}, errRefCodeRequired
	}
	input := graphclient.CreateControlInput{
		RefCode: in.RefCode,
		Tags:    in.Tags,
	}
	if in.Title != "" {
		input.Title = &in.Title
	}
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.Status != "" {
		input.Status = controlStatus(in.Status)
	}
	if in.ReferenceFramework != "" {
		input.ReferenceFramework = &in.ReferenceFramework
	}
	if in.Category != "" {
		input.Category = &in.Category
	}
	if in.Subcategory != "" {
		input.Subcategory = &in.Subcategory
	}
	if in.StandardID != "" {
		input.StandardID = &in.StandardID
	}
	if in.OwnerID != "" {
		input.ControlOwnerID = &in.OwnerID
	}

	resp, err := h.api.CreateControl(ctx, input)
	if err != nil {
		return nil, controlItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedControl(resp.CreateControl.Control), nil
}

func (h *handlers) updateControl(ctx context.Context, _ *mcp.CallToolRequest, in updateControlInput) (*mcp.CallToolResult, controlItem, error) {
	if in.ID == "" {
		return nil, controlItem{}, errIDRequired
	}
	input := graphclient.UpdateControlInput{}
	if in.Title != "" {
		input.Title = &in.Title
	}
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.Status != "" {
		input.Status = controlStatus(in.Status)
	}
	if in.Category != "" {
		input.Category = &in.Category
	}
	if in.Subcategory != "" {
		input.Subcategory = &in.Subcategory
	}
	if in.StandardID != "" {
		input.StandardID = &in.StandardID
	}
	if in.OwnerID != "" {
		input.ControlOwnerID = &in.OwnerID
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if isEmptyUpdateControl(input) {
		return nil, controlItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateControl(ctx, in.ID, input)
	if err != nil {
		return nil, controlItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedControl(resp.UpdateControl.Control), nil
}

func mapCreatedControl(c graphclient.CreateControl_CreateControl_Control) controlItem {
	return controlItem{
		ID:                 c.ID,
		DisplayID:          c.DisplayID,
		RefCode:            c.RefCode,
		Title:              openlane.Deref(c.Title),
		Status:             openlane.Format(c.Status),
		ReferenceFramework: openlane.Deref(c.ReferenceFramework),
		Category:           openlane.Deref(c.Category),
		Description:        openlane.Deref(c.Description),
		Subcategory:        openlane.Deref(c.Subcategory),
		StandardID:         openlane.Deref(c.StandardID),
		OwnerID:            openlane.Deref(c.OwnerID),
		Tags:               c.Tags,
	}
}

func mapUpdatedControl(c graphclient.UpdateControl_UpdateControl_Control) controlItem {
	return controlItem{
		ID:                 c.ID,
		DisplayID:          c.DisplayID,
		RefCode:            c.RefCode,
		Title:              openlane.Deref(c.Title),
		Status:             openlane.Format(c.Status),
		ReferenceFramework: openlane.Deref(c.ReferenceFramework),
		Category:           openlane.Deref(c.Category),
		Description:        openlane.Deref(c.Description),
		Subcategory:        openlane.Deref(c.Subcategory),
		StandardID:         openlane.Deref(c.StandardID),
		OwnerID:            openlane.Deref(c.OwnerID),
		Tags:               c.Tags,
	}
}

func isEmptyUpdateControl(in graphclient.UpdateControlInput) bool {
	return in.Title == nil &&
		in.Description == nil &&
		in.Status == nil &&
		in.Category == nil &&
		in.Subcategory == nil &&
		in.StandardID == nil &&
		in.ControlOwnerID == nil &&
		len(in.Tags) == 0
}
