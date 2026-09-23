package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

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
