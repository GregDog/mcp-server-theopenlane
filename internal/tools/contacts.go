package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type contactItem struct {
	ID          string   `json:"id"`
	FullName    string   `json:"full_name,omitempty"`
	Email       string   `json:"email,omitempty"`
	Company     string   `json:"company,omitempty"`
	Title       string   `json:"title,omitempty"`
	Status      string   `json:"status,omitempty"`
	PhoneNumber string   `json:"phone_number,omitempty"`
	Address     string   `json:"address,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func registerContacts(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_contacts_list",
		Title:       "List Openlane contacts",
		Description: "List contacts in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listContacts)

	addTool(server, &mcp.Tool{
		Name:        "openlane_contact_get",
		Title:       "Get an Openlane contact",
		Description: "Get a single contact by ID.",
		Annotations: readOnly(),
	}, h.getContact)
}

func (h *handlers) listContacts(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[contactItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetContacts(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[contactItem]{}, openlane.APIError(err)
	}
	items := make([]contactItem, 0, len(resp.Contacts.Edges))
	for _, e := range resp.Contacts.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListContact(e.Node))
	}
	return nil, openlane.Page[contactItem]{
		Items:      items,
		NextCursor: resp.Contacts.PageInfo.EndCursor,
		HasMore:    resp.Contacts.PageInfo.HasNextPage,
		TotalCount: resp.Contacts.TotalCount,
	}, nil
}

func (h *handlers) getContact(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, contactItem, error) {
	if in.ID == "" {
		return nil, contactItem{}, errIDRequired
	}
	resp, err := h.api.GetContactByID(ctx, in.ID)
	if err != nil {
		return nil, contactItem{}, openlane.APIError(err)
	}
	return nil, mapGetContact(resp.Contact), nil
}

func mapListContact(n *graphclient.GetContacts_Contacts_Edges_Node) contactItem {
	return contactItem{
		ID:       n.ID,
		FullName: openlane.Deref(n.FullName),
		Email:    openlane.Deref(n.Email),
		Company:  openlane.Deref(n.Company),
		Title:    openlane.Deref(n.Title),
		Status:   openlane.Format(n.Status),
	}
}

func mapGetContact(c graphclient.GetContactByID_Contact) contactItem {
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
