package tools

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type evidenceFileInput struct {
	Filename      string `json:"filename" jsonschema:"Original filename including extension."`
	ContentType   string `json:"content_type,omitempty" jsonschema:"MIME type, for example application/pdf."`
	ContentBase64 string `json:"content_base64" jsonschema:"Base64-encoded file contents."`
}

type createEvidenceInput struct {
	Name                string              `json:"name" jsonschema:"Evidence name."`
	Description         string              `json:"description,omitempty" jsonschema:"Evidence description."`
	Source              string              `json:"source,omitempty" jsonschema:"Evidence source system."`
	CollectionProcedure string              `json:"collection_procedure,omitempty" jsonschema:"How the evidence was collected."`
	URL                 string              `json:"url,omitempty" jsonschema:"External evidence URL when not uploaded as a file."`
	Tags                []string            `json:"tags,omitempty" jsonschema:"Tags to apply."`
	Files               []evidenceFileInput `json:"files,omitempty" jsonschema:"Optional files to upload with the evidence record."`
}

type updateEvidenceInput struct {
	ID                  string              `json:"id" jsonschema:"Evidence ID to update."`
	Name                string              `json:"name,omitempty" jsonschema:"Updated name."`
	Description         string              `json:"description,omitempty" jsonschema:"Updated description."`
	Source              string              `json:"source,omitempty" jsonschema:"Updated source system."`
	CollectionProcedure string              `json:"collection_procedure,omitempty" jsonschema:"Updated collection procedure."`
	URL                 string              `json:"url,omitempty" jsonschema:"Updated external evidence URL."`
	Tags                []string            `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	Files               []evidenceFileInput `json:"files,omitempty" jsonschema:"Optional files to attach to the evidence record."`
}

func registerWriteEvidence(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_evidence_create",
		Title:       "Create Openlane evidence",
		Description: "Create evidence in Openlane. Optional files are uploaded as base64 payloads. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createEvidence)

	addTool(server, &mcp.Tool{
		Name:        "openlane_evidence_update",
		Title:       "Update Openlane evidence",
		Description: "Update evidence by ID. Optional files are uploaded as base64 payloads. Requires write mode.",
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

	uploads, err := h.decodeEvidenceFiles(in.Files)
	if err != nil {
		return nil, evidenceItem{}, err
	}

	resp, err := h.api.CreateEvidence(ctx, input, uploads)
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

	uploads, err := h.decodeEvidenceFiles(in.Files)
	if err != nil {
		return nil, evidenceItem{}, err
	}
	if isEmptyUpdateEvidence(input) && len(uploads) == 0 {
		return nil, evidenceItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateEvidence(ctx, in.ID, input, uploads)
	if err != nil {
		return nil, evidenceItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedEvidence(resp.UpdateEvidence.Evidence), nil
}

func (h *handlers) decodeEvidenceFiles(files []evidenceFileInput) ([]*graphql.Upload, error) {
	if len(files) == 0 {
		return nil, nil
	}
	maxBytes := h.maxUploadBytes
	if maxBytes <= 0 {
		maxBytes = openlane.DefaultMaxUploadBytes
	}
	payload := make([]openlane.EvidenceFile, 0, len(files))
	for _, f := range files {
		payload = append(payload, openlane.EvidenceFile{
			Filename:      f.Filename,
			ContentType:   f.ContentType,
			ContentBase64: f.ContentBase64,
		})
	}
	return openlane.DecodeEvidenceUploads(payload, maxBytes)
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
		FileIDs:      evidenceFileIDsFromCreate(e.Files),
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
		FileIDs:      evidenceFileIDsFromUpdate(e.Files),
	}
}

func evidenceFileIDsFromCreate(files graphclient.CreateEvidence_CreateEvidence_Evidence_Files) []string {
	return fileIDsFromEdges(len(files.Edges), func(i int) string {
		if files.Edges[i] == nil || files.Edges[i].Node == nil {
			return ""
		}
		return files.Edges[i].Node.ID
	})
}

func evidenceFileIDsFromUpdate(files graphclient.UpdateEvidence_UpdateEvidence_Evidence_Files) []string {
	return fileIDsFromEdges(len(files.Edges), func(i int) string {
		if files.Edges[i] == nil || files.Edges[i].Node == nil {
			return ""
		}
		return files.Edges[i].Node.ID
	})
}

func fileIDsFromEdges(n int, idAt func(int) string) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if id := idAt(i); id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func isEmptyUpdateEvidence(in graphclient.UpdateEvidenceInput) bool {
	return in.Name == nil &&
		in.Description == nil &&
		in.Source == nil &&
		in.CollectionProcedure == nil &&
		in.URL == nil &&
		len(in.Tags) == 0
}
