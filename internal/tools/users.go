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
	OrgRole     string `json:"org_role,omitempty"`
}

func registerUsers(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_users_list",
		Title:       "List Openlane users",
		Description: "List users who are members of the configured organization. Filter by name or email. Use to resolve a person to a user ID before creating workflows, assigning vendor owners, or reassigning assignments.",
		Annotations: readOnly(),
	}, h.listUsers)

	addTool(server, &mcp.Tool{
		Name:        "openlane_user_get",
		Title:       "Get an Openlane user",
		Description: "Get an organization member by user ID.",
		Annotations: readOnly(),
	}, h.getUser)
}

func (h *handlers) listUsers(ctx context.Context, _ *mcp.CallToolRequest, in userListInput) (*mcp.CallToolResult, openlane.Page[userItem], error) {
	where, err := buildOrgMemberWhere(h.organizationID, in)
	if err != nil {
		return nil, openlane.Page[userItem]{}, err
	}
	resp, err := h.api.GetOrgMembers(ctx, where)
	if err != nil {
		return nil, openlane.Page[userItem]{}, openlane.APIError(err)
	}

	all := make([]userItem, 0, len(resp.OrgMemberships.Edges))
	for _, e := range resp.OrgMemberships.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		all = append(all, mapOrgMemberUser(*e.Node))
	}

	limit := openlane.ClampLimit(in.Limit)
	start := 0
	if c := strings.TrimSpace(in.Cursor); c != "" {
		if off, err := parseOffsetCursor(c); err == nil && off >= 0 {
			start = off
		}
	}
	end := start + int(limit)
	if end > len(all) {
		end = len(all)
	}
	items := all[start:end]
	hasMore := end < len(all)
	var nextCursor *string
	if hasMore {
		nextCursor = openlane.CursorPtr(formatOffsetCursor(end))
	}

	return nil, openlane.Page[userItem]{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		TotalCount: int64(len(all)),
	}, nil
}

func (h *handlers) getUser(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, userItem, error) {
	if in.ID == "" {
		return nil, userItem{}, errIDRequired
	}
	if strings.TrimSpace(h.organizationID) == "" {
		return nil, userItem{}, errOrganizationRequired
	}
	orgID := h.organizationID
	resp, err := h.api.GetOrgMembers(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
		UserID:         &in.ID,
	})
	if err != nil {
		return nil, userItem{}, openlane.APIError(err)
	}
	for _, e := range resp.OrgMemberships.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		return nil, mapOrgMemberUser(*e.Node), nil
	}
	return nil, userItem{}, fmt.Errorf("no organization user matched id %q", in.ID)
}

func mapOrgMemberUser(n graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node) userItem {
	u := n.User
	return userItem{
		ID:          u.ID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		FirstName:   openlane.Deref(u.FirstName),
		LastName:    openlane.Deref(u.LastName),
		OrgRole:     openlane.Format(n.Role),
	}
}

func buildOrgMemberWhere(orgID string, in userListInput) (*graphclient.OrgMembershipWhereInput, error) {
	orgID = strings.TrimSpace(orgID)
	if orgID == "" {
		return nil, errOrganizationRequired
	}
	where := &graphclient.OrgMembershipWhereInput{OrganizationID: &orgID}

	email := strings.TrimSpace(in.Email)
	name := strings.TrimSpace(in.Name)
	if email == "" && name == "" {
		return where, nil
	}

	var userOr []*graphclient.UserWhereInput
	if email != "" {
		userOr = append(userOr, &graphclient.UserWhereInput{EmailContainsFold: &email})
	}
	if name != "" {
		userOr = append(userOr,
			&graphclient.UserWhereInput{DisplayNameContainsFold: &name},
			&graphclient.UserWhereInput{FirstNameContainsFold: &name},
			&graphclient.UserWhereInput{LastNameContainsFold: &name},
		)
	}
	userWhere := userOr[0]
	if len(userOr) > 1 {
		userWhere = &graphclient.UserWhereInput{Or: userOr}
	}
	where.HasUserWith = []*graphclient.UserWhereInput{userWhere}
	return where, nil
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
	where, err := buildOrgMemberWhere(h.organizationID, in)
	if err != nil {
		return "", err
	}
	resp, err := h.api.GetOrgMembers(ctx, where)
	if err != nil {
		return "", openlane.APIError(err)
	}
	var matches []string
	for _, e := range resp.OrgMemberships.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		u := e.Node.User
		if strings.EqualFold(u.Email, s) || strings.EqualFold(u.DisplayName, s) {
			matches = append(matches, u.ID)
		}
	}
	if len(matches) == 0 {
		for _, e := range resp.OrgMemberships.Edges {
			if e != nil && e.Node != nil && e.Node.User.ID != "" {
				matches = append(matches, e.Node.User.ID)
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

func parseOffsetCursor(cursor string) (int, error) {
	const prefix = "offset:"
	if !strings.HasPrefix(cursor, prefix) {
		return 0, fmt.Errorf("invalid cursor")
	}
	var off int
	_, err := fmt.Sscanf(cursor, prefix+"%d", &off)
	return off, err
}

func formatOffsetCursor(offset int) string {
	return fmt.Sprintf("offset:%d", offset)
}
