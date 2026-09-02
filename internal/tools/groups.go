package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type groupListInput struct {
	listInput
	Name string `json:"name,omitempty" jsonschema:"Filter groups whose name or display name contains this text (case-insensitive)."`
}

type groupItem struct {
	ID          string `json:"id"`
	DisplayID   string `json:"display_id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	IsManaged   *bool  `json:"is_managed,omitempty"`
}

func registerGroups(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_groups_list",
		Title:       "List Openlane groups",
		Description: "List groups in the configured organization. Filter by name or display name. Use to resolve a group name to an ID before creating workflows.",
		Annotations: readOnly(),
	}, h.listGroups)

	addTool(server, &mcp.Tool{
		Name:        "openlane_group_get",
		Title:       "Get an Openlane group",
		Description: "Get a group by ID.",
		Annotations: readOnly(),
	}, h.getGroup)
}

func (h *handlers) listGroups(ctx context.Context, _ *mcp.CallToolRequest, in groupListInput) (*mcp.CallToolResult, openlane.Page[groupItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetGroups(ctx, &first, after, buildGroupWhere(in.Name))
	if err != nil {
		return nil, openlane.Page[groupItem]{}, openlane.APIError(err)
	}
	items := make([]groupItem, 0, len(resp.Groups.Edges))
	for _, e := range resp.Groups.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapGroupNode(*e.Node))
	}
	return nil, openlane.Page[groupItem]{
		Items:      items,
		NextCursor: resp.Groups.PageInfo.EndCursor,
		HasMore:    resp.Groups.PageInfo.HasNextPage,
		TotalCount: resp.Groups.TotalCount,
	}, nil
}

func (h *handlers) getGroup(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, groupItem, error) {
	if in.ID == "" {
		return nil, groupItem{}, errIDRequired
	}
	resp, err := h.api.GetGroupByID(ctx, in.ID)
	if err != nil {
		return nil, groupItem{}, openlane.APIError(err)
	}
	g := resp.Group
	return nil, groupItem{
		ID:          g.ID,
		DisplayID:   g.DisplayID,
		Name:        g.Name,
		DisplayName: g.DisplayName,
		Description: openlane.Deref(g.Description),
		IsManaged:   g.IsManaged,
	}, nil
}

func mapGroupNode(n graphclient.GetGroups_Groups_Edges_Node) groupItem {
	return groupItem{
		ID:          n.ID,
		DisplayID:   n.DisplayID,
		Name:        n.Name,
		DisplayName: n.DisplayName,
		Description: openlane.Deref(n.Description),
		IsManaged:   n.IsManaged,
	}
}

func buildGroupWhere(name string) *graphclient.GroupWhereInput {
	s := strings.TrimSpace(name)
	if s == "" {
		return nil
	}
	return &graphclient.GroupWhereInput{
		Or: []*graphclient.GroupWhereInput{
			{NameContainsFold: &s},
			{DisplayNameContainsFold: &s},
		},
	}
}

func (h *handlers) resolveGroupID(ctx context.Context, nameOrID string) (string, error) {
	s := strings.TrimSpace(nameOrID)
	if s == "" {
		return "", fmt.Errorf("group name or id is required")
	}
	if looksLikeOpenlaneID(s) {
		return s, nil
	}
	first := int64(10)
	resp, err := h.api.GetGroups(ctx, &first, nil, buildGroupWhere(s))
	if err != nil {
		return "", openlane.APIError(err)
	}
	var matches []string
	for _, e := range resp.Groups.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		if strings.EqualFold(n.Name, s) || strings.EqualFold(n.DisplayName, s) {
			matches = append(matches, n.ID)
		}
	}
	if len(matches) == 0 {
		// fallback: any contains-fold hit
		for _, e := range resp.Groups.Edges {
			if e != nil && e.Node != nil {
				matches = append(matches, e.Node.ID)
			}
		}
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return "", fmt.Errorf("no group matched %q", s)
	default:
		return "", fmt.Errorf("multiple groups matched %q; use an id", s)
	}
}

func looksLikeOpenlaneID(s string) bool {
	return strings.HasPrefix(s, "01") && len(s) >= 20
}
