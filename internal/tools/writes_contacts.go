package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createContactInput struct {
	FullName    string   `json:"full_name,omitempty" jsonschema:"Contact full name."`
	Email       string   `json:"email,omitempty" jsonschema:"Contact email address."`
	Title       string   `json:"title,omitempty" jsonschema:"Contact job title."`
	Company     string   `json:"company,omitempty" jsonschema:"Contact company name."`
	PhoneNumber string   `json:"phone_number,omitempty" jsonschema:"Contact phone number."`
	Address     string   `json:"address,omitempty" jsonschema:"Contact postal address."`
	Status      string   `json:"status,omitempty" jsonschema:"Contact status enum value."`
	Tags        []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
	EntityIDs   []string `json:"entity_ids,omitempty" jsonschema:"Entity (vendor) IDs to associate with this contact on create."`
}

func registerWriteContacts(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_contact_create",
		Title:       "Create an Openlane contact",
		Description: "Create a contact in Openlane. Optionally associate with entities (vendors) via entity_ids. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createContact)
}

func (h *handlers) createContact(ctx context.Context, _ *mcp.CallToolRequest, in createContactInput) (*mcp.CallToolResult, contactItem, error) {
	if strings.TrimSpace(in.FullName) == "" && strings.TrimSpace(in.Email) == "" {
		return nil, contactItem{}, errContactIdentifierRequired
	}

	input := graphclient.CreateContactInput{
		Tags:      in.Tags,
		EntityIDs: in.EntityIDs,
	}
	if s := strings.TrimSpace(in.FullName); s != "" {
		input.FullName = &s
	}
	if s := strings.TrimSpace(in.Email); s != "" {
		input.Email = &s
	}
	if s := strings.TrimSpace(in.Title); s != "" {
		input.Title = &s
	}
	if s := strings.TrimSpace(in.Company); s != "" {
		input.Company = &s
	}
	if s := strings.TrimSpace(in.PhoneNumber); s != "" {
		input.PhoneNumber = &s
	}
	if s := strings.TrimSpace(in.Address); s != "" {
		input.Address = &s
	}
	if in.Status != "" {
		input.Status = userStatus(in.Status)
	}

	resp, err := h.api.CreateContact(ctx, input)
	if err != nil {
		return nil, contactItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedContact(resp.CreateContact.Contact), nil
}

func mapCreatedContact(c graphclient.CreateContact_CreateContact_Contact) contactItem {
	return contactItem{
		ID:          c.ID,
		FullName:    openlane.Deref(c.FullName),
		Email:       openlane.Deref(c.Email),
		Company:     openlane.Deref(c.Company),
		Title:       openlane.Deref(c.Title),
		Status:      openlane.Format(c.Status),
		PhoneNumber: openlane.Deref(c.PhoneNumber),
		Address:     openlane.Deref(c.Address),
		Tags:        c.Tags,
	}
}
