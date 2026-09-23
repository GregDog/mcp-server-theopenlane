package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

// groupAssigneeSummary resolves Openlane controlOwnerID/delegateID group ids to names and users.
type groupAssigneeSummary struct {
	GroupID          string `json:"group_id,omitempty"`
	GroupName        string `json:"group_name,omitempty"`
	GroupDisplayName string `json:"group_display_name,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	UserDisplayName  string `json:"user_display_name,omitempty"`
	UserEmail        string `json:"user_email,omitempty"`
}

// resolveControlAssigneeGroupID maps a user id/email/name or group id/name to the
// Openlane group id required by controlOwnerID and delegateID on UpdateControl.
func (h *handlers) resolveControlAssigneeGroupID(ctx context.Context, assignee string) (string, error) {
	s := strings.TrimSpace(assignee)
	if s == "" {
		return "", fmt.Errorf("assignee is required")
	}
	if looksLikeOpenlaneID(s) {
		if groupID, err := h.lookupGroupByID(ctx, s); err == nil {
			return groupID, nil
		}
		return h.resolvePersonalManagedGroupForUser(ctx, s)
	}
	if strings.Contains(s, "@") {
		userID, err := h.resolveUserID(ctx, s)
		if err != nil {
			return "", err
		}
		return h.resolvePersonalManagedGroupForUser(ctx, userID)
	}
	userID, userErr := h.resolveUserID(ctx, s)
	if userErr == nil {
		return h.resolvePersonalManagedGroupForUser(ctx, userID)
	}
	return h.resolveGroupID(ctx, s)
}

func (h *handlers) lookupGroupByID(ctx context.Context, id string) (string, error) {
	resp, err := h.api.GetGroupByID(ctx, id)
	if err != nil {
		return "", openlane.APIError(err)
	}
	if resp.Group.ID == "" {
		return "", fmt.Errorf("group not found")
	}
	return resp.Group.ID, nil
}

func (h *handlers) resolvePersonalManagedGroupForUser(ctx context.Context, userID string) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", fmt.Errorf("user id is required")
	}
	managed := true
	first := int64(20)
	resp, err := h.api.GetGroups(ctx, &first, nil, &graphclient.GroupWhereInput{
		IsManaged: &managed,
		HasMembersWith: []*graphclient.GroupMembershipWhereInput{
			{UserID: &userID},
		},
	})
	if err != nil {
		return "", openlane.APIError(err)
	}
	nodes := make([]graphclient.GetGroups_Groups_Edges_Node, 0, len(resp.Groups.Edges))
	for _, e := range resp.Groups.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		nodes = append(nodes, *e.Node)
	}
	return pickPersonalManagedGroup(userID, nodes)
}

func (h *handlers) resolveGroupAssigneeSummary(ctx context.Context, groupID string, cache map[string]groupAssigneeSummary) (groupAssigneeSummary, error) {
	groupID = strings.TrimSpace(groupID)
	if groupID == "" {
		return groupAssigneeSummary{}, nil
	}
	if cache != nil {
		if cached, ok := cache[groupID]; ok {
			return cached, nil
		}
	}

	resp, err := h.api.GetGroupByID(ctx, groupID)
	if err != nil {
		return groupAssigneeSummary{}, openlane.APIError(err)
	}
	g := resp.Group
	out := groupAssigneeSummary{
		GroupID:          g.ID,
		GroupName:        g.Name,
		GroupDisplayName: g.DisplayName,
	}
	if openlane.Deref(g.IsManaged) {
		out.UserID = userIDFromManagedGroupName(g.Name)
	}
	if out.UserID != "" {
		displayName, email := h.lookupOrgMemberContact(ctx, out.UserID)
		out.UserDisplayName = displayName
		out.UserEmail = email
	}
	if cache != nil {
		cache[groupID] = out
	}
	return out, nil
}

func (h *handlers) lookupOrgMemberContact(ctx context.Context, userID string) (displayName, email string) {
	userID = strings.TrimSpace(userID)
	if userID == "" || strings.TrimSpace(h.organizationID) == "" {
		return "", ""
	}
	orgID := h.organizationID
	resp, err := h.api.GetOrgMembers(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
		UserID:         &userID,
	})
	if err != nil {
		return "", ""
	}
	for _, e := range resp.OrgMemberships.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		u := e.Node.User
		return u.DisplayName, u.Email
	}
	return "", ""
}

func userIDFromManagedGroupName(name string) string {
	const sep = " - "
	i := strings.LastIndex(name, sep)
	if i < 0 {
		return ""
	}
	id := strings.TrimSpace(name[i+len(sep):])
	if looksLikeOpenlaneID(id) {
		return id
	}
	return ""
}

func (h *handlers) enrichControlAssignees(ctx context.Context, item *controlItem, cache map[string]groupAssigneeSummary) error {
	if item == nil {
		return nil
	}
	if item.ControlOwnerID != "" {
		summary, err := h.resolveGroupAssigneeSummary(ctx, item.ControlOwnerID, cache)
		if err != nil {
			return err
		}
		if summary.GroupID != "" {
			item.ControlOwner = &summary
		}
	}
	if item.DelegateID != "" {
		summary, err := h.resolveGroupAssigneeSummary(ctx, item.DelegateID, cache)
		if err != nil {
			return err
		}
		if summary.GroupID != "" {
			item.Delegate = &summary
		}
	}
	return nil
}

func enrichControlPageAssignees(ctx context.Context, h *handlers, items []controlItem) error {
	cache := make(map[string]groupAssigneeSummary)
	for i := range items {
		if err := h.enrichControlAssignees(ctx, &items[i], cache); err != nil {
			return err
		}
	}
	return nil
}

func pickPersonalManagedGroup(userID string, groups []graphclient.GetGroups_Groups_Edges_Node) (string, error) {
	if len(groups) == 0 {
		return "", fmt.Errorf("no managed group found for user %s", userID)
	}
	suffix := " - " + userID
	var suffixMatches []string
	for _, g := range groups {
		if strings.HasSuffix(g.Name, suffix) {
			suffixMatches = append(suffixMatches, g.ID)
		}
	}
	switch len(suffixMatches) {
	case 1:
		return suffixMatches[0], nil
	case 0:
		if len(groups) == 1 {
			return groups[0].ID, nil
		}
		return "", fmt.Errorf("multiple managed groups for user %s; use a group id", userID)
	default:
		return "", fmt.Errorf("multiple managed groups matched user %s; use a group id", userID)
	}
}
