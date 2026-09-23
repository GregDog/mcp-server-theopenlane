package tools

import (
	"context"
	"fmt"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func controlWhereLinkableOnly(linkableOnly bool) *graphclient.ControlWhereInput {
	if !linkableOnly {
		return nil
	}
	ownerSet := true
	return &graphclient.ControlWhereInput{
		OwnerIDNotNil: &ownerSet,
	}
}

func mergeControlWhere(base, extra *graphclient.ControlWhereInput) *graphclient.ControlWhereInput {
	if base == nil {
		return extra
	}
	if extra == nil {
		return base
	}
	merged := *base
	if extra.OwnerIDNotNil != nil {
		merged.OwnerIDNotNil = extra.OwnerIDNotNil
	}
	if len(extra.SourceNotIn) > 0 {
		merged.SourceNotIn = extra.SourceNotIn
	}
	if len(extra.SourceIn) > 0 {
		merged.SourceIn = extra.SourceIn
	}
	if extra.SourceNeq != nil {
		merged.SourceNeq = extra.SourceNeq
	}
	if extra.Source != nil {
		merged.Source = extra.Source
	}
	return &merged
}

func mapControlListNode(n *graphclient.GetControls_Controls_Edges_Node) controlItem {
	return controlItem{
		ID:                 n.ID,
		DisplayID:          n.DisplayID,
		RefCode:            n.RefCode,
		Title:              openlane.Deref(n.Title),
		Status:             openlane.Format(n.Status),
		Source:             openlane.Format(n.Source),
		OwnerID:            openlane.Deref(n.OwnerID),
		ControlOwnerID:     openlane.Deref(n.ControlOwnerID),
		DelegateID:         openlane.Deref(n.DelegateID),
		LinkableToEvidence: openlane.ControlLinkableToEvidence(n.OwnerID),
		ReferenceFramework: openlane.Deref(n.ReferenceFramework),
		Category:           openlane.Deref(n.Category),
	}
}

func mapControlDetail(c graphclient.GetControlByID_Control) controlItem {
	return controlItem{
		ID:                     c.ID,
		DisplayID:              c.DisplayID,
		RefCode:                c.RefCode,
		Title:                  openlane.Deref(c.Title),
		Status:                 openlane.Format(c.Status),
		Source:                 openlane.Format(c.Source),
		OwnerID:                openlane.Deref(c.OwnerID),
		LinkableToEvidence:     openlane.ControlLinkableToEvidence(c.OwnerID),
		ReferenceFramework:     openlane.Deref(c.ReferenceFramework),
		Category:               openlane.Deref(c.Category),
		Description:            openlane.Deref(c.Description),
		Subcategory:            openlane.Deref(c.Subcategory),
		StandardID:             openlane.Deref(c.StandardID),
		ControlOwnerID:         openlane.Deref(c.ControlOwnerID),
		DelegateID:             openlane.Deref(c.DelegateID),
		ControlKindName:        openlane.Deref(c.ControlKindName),
		ImplementationGuidance: c.ImplementationGuidance,
		AssessmentMethods:      c.AssessmentMethods,
		AssessmentObjectives:   c.AssessmentObjectives,
		Tags:                   c.Tags,
	}
}

func (h *handlers) validateEvidenceControlIDs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		resp, err := h.api.GetControlByID(ctx, id)
		if err != nil {
			return openlane.APIError(err)
		}
		c := resp.Control
		if !openlane.ControlLinkableToEvidence(c.OwnerID) {
			title := openlane.Deref(c.Title)
			return fmt.Errorf(
				"control %s (%s %q) is a catalog/system control and cannot be linked to evidence; use the org-owned copy from your program (owner_id set, display_id like CTL-...). Search with linkable_only: true",
				c.ID,
				c.RefCode,
				title,
			)
		}
	}
	return nil
}
