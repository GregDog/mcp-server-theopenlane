package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createEvidenceInput struct {
	Name                string   `json:"name" jsonschema:"Evidence name."`
	Description         string   `json:"description,omitempty" jsonschema:"Evidence description."`
	Source              string   `json:"source,omitempty" jsonschema:"Evidence source system."`
	CollectionProcedure string   `json:"collection_procedure,omitempty" jsonschema:"How the evidence was collected."`
	URL                 string   `json:"url,omitempty" jsonschema:"External evidence URL when not uploaded as a file."`
	Tags                []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
}

type updateEvidenceInput struct {
	ID                  string   `json:"id" jsonschema:"Evidence ID to update."`
	Name                string   `json:"name,omitempty" jsonschema:"Updated name."`
	Description         string   `json:"description,omitempty" jsonschema:"Updated description."`
	Source              string   `json:"source,omitempty" jsonschema:"Updated source system."`
	CollectionProcedure string   `json:"collection_procedure,omitempty" jsonschema:"Updated collection procedure."`
	URL                 string   `json:"url,omitempty" jsonschema:"Updated external evidence URL."`
	Tags                []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
}

func registerWriteEvidence(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_evidence_create",
		Title:       "Create Openlane evidence metadata",
		Description: "Create evidence metadata in Openlane. File uploads are not supported. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createEvidence)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_evidence_update",
		Title:       "Update Openlane evidence metadata",
		Description: "Update evidence metadata by ID. File uploads are not supported. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateEvidence)
}

func (h *handlers) createEvidence(ctx context.Context, _ *mcp.CallToolRequest, in createEvidenceInput) (*mcp.CallToolResult, evidenceItem, error) {
	if in.Name == "" {
		return nil, evidenceItem{}, errNameRequired
	}
	input := graphclient.CreateEvidenceInput{
		Name: in.Name,
		Tags: in.Tags,
	}
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.Source != "" {
		input.Source = &in.Source
	}
	if in.CollectionProcedure != "" {
		input.CollectionProcedure = &in.CollectionProcedure
	}
	if in.URL != "" {
		input.URL = &in.URL
	}

	resp, err := h.api.CreateEvidence(ctx, input)
	if err != nil {
		return nil, evidenceItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedEvidence(resp.CreateEvidence.Evidence), nil
}

func (h *handlers) updateEvidence(ctx context.Context, _ *mcp.CallToolRequest, in updateEvidenceInput) (*mcp.CallToolResult, evidenceItem, error) {
	if in.ID == "" {
		return nil, evidenceItem{}, errIDRequired
	}
	input := graphclient.UpdateEvidenceInput{}
	if in.Name != "" {
		input.Name = &in.Name
	}
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.Source != "" {
		input.Source = &in.Source
	}
	if in.CollectionProcedure != "" {
		input.CollectionProcedure = &in.CollectionProcedure
	}
	if in.URL != "" {
		input.URL = &in.URL
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if isEmptyUpdateEvidence(input) {
		return nil, evidenceItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateEvidence(ctx, in.ID, input)
	if err != nil {
		return nil, evidenceItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedEvidence(resp.UpdateEvidence.Evidence), nil
}

func mapCreatedEvidence(e graphclient.CreateEvidence_CreateEvidence_Evidence) evidenceItem {
	return evidenceItem{
		ID:           e.ID,
		DisplayID:    e.DisplayID,
		Name:         e.Name,
		Status:       openlane.Format(e.Status),
		Source:       openlane.Deref(e.Source),
		CreationDate: openlane.Format(e.CreationDate),
		Description:  openlane.Deref(e.Description),
		RenewalDate:  openlane.Format(e.RenewalDate),
		IsAutomated:  e.IsAutomated,
	}
}

func mapUpdatedEvidence(e graphclient.UpdateEvidence_UpdateEvidence_Evidence) evidenceItem {
	return evidenceItem{
		ID:           e.ID,
		DisplayID:    e.DisplayID,
		Name:         e.Name,
		Status:       openlane.Format(e.Status),
		Source:       openlane.Deref(e.Source),
		CreationDate: openlane.Format(e.CreationDate),
		Description:  openlane.Deref(e.Description),
		RenewalDate:  openlane.Format(e.RenewalDate),
		IsAutomated:  e.IsAutomated,
	}
}

func isEmptyUpdateEvidence(in graphclient.UpdateEvidenceInput) bool {
	return in.Name == nil &&
		in.Description == nil &&
		in.Source == nil &&
		in.CollectionProcedure == nil &&
		in.URL == nil &&
		len(in.Tags) == 0
}
