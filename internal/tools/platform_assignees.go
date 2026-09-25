package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

// platformOwnerSummary is a resolved platform owner role for MCP responses.
// See docs/openlane-assignee-ids.md.
type platformOwnerSummary struct {
	Kind             string `json:"kind,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	UserEmail        string `json:"user_email,omitempty"`
	UserDisplayName  string `json:"user_display_name,omitempty"`
	GroupID          string `json:"group_id,omitempty"`
	GroupName        string `json:"group_name,omitempty"`
	GroupDisplayName string `json:"group_display_name,omitempty"`
	Name             string `json:"name,omitempty"`
}

type platformOwnerWrite struct {
	userID  *string
	groupID *string
	name    *string
}

func (h *handlers) resolvePlatformOwner(ctx context.Context, assignee string) (platformOwnerWrite, error) {
	s := strings.TrimSpace(assignee)
	if s == "" {
		return platformOwnerWrite{}, fmt.Errorf("owner is required")
	}
	if looksLikeOpenlaneID(s) {
		if groupID, err := h.lookupGroupByID(ctx, s); err == nil {
			return platformOwnerWrite{groupID: &groupID}, nil
		}
		userID := s
		return platformOwnerWrite{userID: &userID}, nil
	}
	if strings.Contains(s, "@") {
		userID, err := h.resolveUserID(ctx, s)
		if err != nil {
			return platformOwnerWrite{}, err
		}
		return platformOwnerWrite{userID: &userID}, nil
	}
	userID, userErr := h.resolveUserID(ctx, s)
	if userErr == nil {
		return platformOwnerWrite{userID: &userID}, nil
	}
	groupID, err := h.resolveGroupID(ctx, s)
	if err != nil {
		return platformOwnerWrite{}, err
	}
	return platformOwnerWrite{groupID: &groupID}, nil
}

func (h *handlers) buildPlatformOwnerWrite(ctx context.Context, assignee, ownerName string) (platformOwnerWrite, error) {
	assignee = strings.TrimSpace(assignee)
	ownerName = strings.TrimSpace(ownerName)
	if assignee != "" && ownerName != "" {
		return platformOwnerWrite{}, fmt.Errorf("set either owner or owner_name, not both")
	}
	if ownerName != "" {
		return platformOwnerWrite{name: &ownerName}, nil
	}
	if assignee == "" {
		return platformOwnerWrite{}, nil
	}
	return h.resolvePlatformOwner(ctx, assignee)
}

func setPlatformOwnerCreate(input *graphclient.CreatePlatformInput, role string, write platformOwnerWrite) error {
	switch role {
	case "business":
		input.BusinessOwnerUserID = write.userID
		input.BusinessOwnerGroupID = write.groupID
		input.BusinessOwner = write.name
	case "technical":
		input.TechnicalOwnerUserID = write.userID
		input.TechnicalOwnerGroupID = write.groupID
		input.TechnicalOwner = write.name
	case "security":
		input.SecurityOwnerUserID = write.userID
		input.SecurityOwnerGroupID = write.groupID
		input.SecurityOwner = write.name
	case "internal":
		input.InternalOwnerUserID = write.userID
		input.InternalOwnerGroupID = write.groupID
		input.InternalOwner = write.name
	default:
		return fmt.Errorf("unknown owner role %q", role)
	}
	return nil
}

func setPlatformOwnerUpdate(input *graphclient.UpdatePlatformInput, role string, write platformOwnerWrite) error {
	switch role {
	case "business":
		input.BusinessOwnerUserID = write.userID
		input.BusinessOwnerGroupID = write.groupID
		input.BusinessOwner = write.name
	case "technical":
		input.TechnicalOwnerUserID = write.userID
		input.TechnicalOwnerGroupID = write.groupID
		input.TechnicalOwner = write.name
	case "security":
		input.SecurityOwnerUserID = write.userID
		input.SecurityOwnerGroupID = write.groupID
		input.SecurityOwner = write.name
	case "internal":
		input.InternalOwnerUserID = write.userID
		input.InternalOwnerGroupID = write.groupID
		input.InternalOwner = write.name
	default:
		return fmt.Errorf("unknown owner role %q", role)
	}
	return nil
}

func (h *handlers) enrichPlatformOwnerFromScalars(ctx context.Context, userID, groupID, freeText string, cache map[string]groupAssigneeSummary) (*platformOwnerSummary, error) {
	userID = strings.TrimSpace(userID)
	groupID = strings.TrimSpace(groupID)
	freeText = strings.TrimSpace(freeText)
	if userID != "" {
		displayName, email := h.lookupOrgMemberContact(ctx, userID)
		return &platformOwnerSummary{
			Kind:            "user",
			UserID:          userID,
			UserEmail:       email,
			UserDisplayName: displayName,
		}, nil
	}
	if groupID != "" {
		summary, err := h.resolveGroupAssigneeSummary(ctx, groupID, cache)
		if err != nil {
			return nil, err
		}
		if summary.GroupID == "" {
			return nil, nil
		}
		out := &platformOwnerSummary{
			Kind:             "group",
			GroupID:          summary.GroupID,
			GroupName:        summary.GroupName,
			GroupDisplayName: summary.GroupDisplayName,
		}
		if summary.UserID != "" {
			out.UserID = summary.UserID
			out.UserEmail = summary.UserEmail
			out.UserDisplayName = summary.UserDisplayName
		}
		return out, nil
	}
	if freeText != "" {
		return &platformOwnerSummary{Kind: "name", Name: freeText}, nil
	}
	return nil, nil
}

func mapPlatformOwnerRole(role openlane.PlatformOwnerRole) *platformOwnerSummary {
	if role.Kind == "" {
		return nil
	}
	return &platformOwnerSummary{
		Kind:             role.Kind,
		UserID:           role.UserID,
		UserEmail:        role.UserEmail,
		UserDisplayName:  role.UserDisplayName,
		GroupID:          role.GroupID,
		GroupName:        role.GroupName,
		GroupDisplayName: role.GroupDisplayName,
		Name:             role.Name,
	}
}
