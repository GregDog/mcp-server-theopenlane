package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type userListInput struct {
	listInput
	Name  string `json:"name,omitempty" jsonschema:"Filter users whose display name, first name, or last name contains this text (case-insensitive)."`
	Email string `json:"email,omitempty" jsonschema:"Filter users whose email contains this text (case-insensitive)."`
}

type userItem struct {
	ID          string `json:"id"`
	DisplayID   string `json:"display_id,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
}

func registerUsers(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_users_list",
		Title:       "List Openlane users",
		Description: "List users in the configured organization. Filter by name or email. Use to resolve a person to a user ID before creating workflows or reassigning assignments.",
		Annotations: readOnly(),
	}, h.listUsers)

	addTool(server, &mcp.Tool{
		Name:        "openlane_user_get",
		Title:       "Get an Openlane user",
		Description: "Get a user by ID.",
		Annotations: readOnly(),
	}, h.getUser)
}

func (h *handlers) listUsers(ctx context.Context, _ *mcp.CallToolRequest, in userListInput) (*mcp.CallToolResult, openlane.Page[userItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetUsers(ctx, &first, after, buildUserWhere(in))
	if err != nil {
		return nil, openlane.Page[userItem]{}, openlane.APIError(err)
	}
	items := make([]userItem, 0, len(resp.Users.Edges))
	for _, e := range resp.Users.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapUserNode(*e.Node))
	}
	return nil, openlane.Page[userItem]{
		Items:      items,
		NextCursor: resp.Users.PageInfo.EndCursor,
		HasMore:    resp.Users.PageInfo.HasNextPage,
		TotalCount: resp.Users.TotalCount,
	}, nil
}

func (h *handlers) getUser(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, userItem, error) {
	if in.ID == "" {
		return nil, userItem{}, errIDRequired
	}
	resp, err := h.api.GetUserByID(ctx, in.ID)
	if err != nil {
		return nil, userItem{}, openlane.APIError(err)
	}
	u := resp.User
	return nil, userItem{
		ID:          u.ID,
		DisplayID:   u.DisplayID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		FirstName:   openlane.Deref(u.FirstName),
		LastName:    openlane.Deref(u.LastName),
	}, nil
}

func mapUserNode(n graphclient.GetUsers_Users_Edges_Node) userItem {
	return userItem{
		ID:          n.ID,
		DisplayID:   n.DisplayID,
		DisplayName: n.DisplayName,
		Email:       n.Email,
		FirstName:   openlane.Deref(n.FirstName),
		LastName:    openlane.Deref(n.LastName),
	}
}

func buildUserWhere(in userListInput) *graphclient.UserWhereInput {
	email := strings.TrimSpace(in.Email)
	name := strings.TrimSpace(in.Name)
	if email == "" && name == "" {
		return nil
	}
	var or []*graphclient.UserWhereInput
	if email != "" {
		or = append(or, &graphclient.UserWhereInput{EmailContainsFold: &email})
	}
	if name != "" {
		or = append(or,
			&graphclient.UserWhereInput{DisplayNameContainsFold: &name},
			&graphclient.UserWhereInput{FirstNameContainsFold: &name},
			&graphclient.UserWhereInput{LastNameContainsFold: &name},
		)
	}
	if len(or) == 1 {
		return or[0]
	}
	return &graphclient.UserWhereInput{Or: or}
}

func (h *handlers) resolveUserID(ctx context.Context, nameEmailOrID string) (string, error) {
	s := strings.TrimSpace(nameEmailOrID)
	if s == "" {
		return "", fmt.Errorf("user name, email, or id is required")
	}
	if looksLikeOpenlaneID(s) {
		return s, nil
	}
	in := userListInput{}
	if strings.Contains(s, "@") {
		in.Email = s
	} else {
		in.Name = s
	}
	first := int64(10)
	resp, err := h.api.GetUsers(ctx, &first, nil, buildUserWhere(in))
	if err != nil {
		return "", openlane.APIError(err)
	}
	var matches []string
	for _, e := range resp.Users.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		if strings.EqualFold(n.Email, s) || strings.EqualFold(n.DisplayName, s) {
			matches = append(matches, n.ID)
		}
	}
	if len(matches) == 0 {
		for _, e := range resp.Users.Edges {
			if e != nil && e.Node != nil {
				matches = append(matches, e.Node.ID)
			}
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("no user matched %q", s)
	default:
		return "", fmt.Errorf("multiple users matched %q; use an id", s)
	}
}
