package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type mappedControlEndpointRef struct {
	ID      string `json:"id"`
	RefCode string `json:"ref_code,omitempty"`
}

type mappedControlItem struct {
	ID              string                     `json:"id"`
	MappingType     string                     `json:"mapping_type"`
	Source          string                     `json:"source,omitempty"`
	Relation        string                     `json:"relation,omitempty"`
	Confidence      *int64                     `json:"confidence,omitempty"`
	Tags            []string                   `json:"tags,omitempty"`
	FromControls    []mappedControlEndpointRef `json:"from_controls,omitempty"`
	ToControls      []mappedControlEndpointRef `json:"to_controls,omitempty"`
	FromSubcontrols []mappedControlEndpointRef `json:"from_subcontrols,omitempty"`
	ToSubcontrols   []mappedControlEndpointRef `json:"to_subcontrols,omitempty"`
}

type relatedControlRef struct {
	MappingID   string `json:"mapping_id"`
	MappingType string `json:"mapping_type,omitempty"`
	Source      string `json:"source,omitempty"`
	Relation    string `json:"relation,omitempty"`
	Direction   string `json:"direction"`
	ID          string `json:"id"`
	RefCode     string `json:"ref_code,omitempty"`
}

func registerMappedControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_mapped_controls_list",
		Title:       "List Openlane control mappings",
		Description: "List MappedControl records that link controls and subcontrols across frameworks or within the org. Filter by control_id to see mappings where the control appears on either side. Results are paginated.",
		Annotations: readOnly(),
	}, h.listMappedControls)

	addTool(server, &mcp.Tool{
		Name:        "openlane_mapped_control_get",
		Title:       "Get an Openlane control mapping",
		Description: "Get a single MappedControl by ID, including from/to control and subcontrol ref codes.",
		Annotations: readOnly(),
	}, h.getMappedControl)
}

func (h *handlers) listMappedControls(ctx context.Context, _ *mcp.CallToolRequest, in mappedControlListInput) (*mcp.CallToolResult, openlane.Page[mappedControlItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetMappedControls(ctx, &first, after, buildMappedControlWhere(in))
	if err != nil {
		return nil, openlane.Page[mappedControlItem]{}, openlane.APIError(err)
	}
	items := make([]mappedControlItem, 0, len(resp.MappedControls.Edges))
	for _, e := range resp.MappedControls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListMappedControl(e.Node))
	}
	return nil, openlane.Page[mappedControlItem]{
		Items:      items,
		NextCursor: resp.MappedControls.PageInfo.EndCursor,
		HasMore:    resp.MappedControls.PageInfo.HasNextPage,
		TotalCount: resp.MappedControls.TotalCount,
	}, nil
}

func (h *handlers) getMappedControl(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, mappedControlItem, error) {
	if in.ID == "" {
		return nil, mappedControlItem{}, errIDRequired
	}
	resp, err := h.api.GetMappedControlByID(ctx, in.ID)
	if err != nil {
		return nil, mappedControlItem{}, openlane.APIError(err)
	}
	return nil, mapGetMappedControl(resp.MappedControl), nil
}

func (h *handlers) fetchRelatedControls(ctx context.Context, controlID string) *relSummary[relatedControlRef] {
	controlID = strings.TrimSpace(controlID)
	if controlID == "" {
		return nil
	}
	cw := []*graphclient.ControlWhereInput{{ID: &controlID}}
	first := relFirst()
	resp, err := h.api.GetMappedControls(ctx, &first, nil, &graphclient.MappedControlWhereInput{
		Or: []*graphclient.MappedControlWhereInput{
			{HasFromControlsWith: cw},
			{HasToControlsWith: cw},
		},
	})
	if err != nil {
		return nil
	}
	items := make([]relatedControlRef, 0, len(resp.MappedControls.Edges))
	for _, e := range resp.MappedControls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, relatedControlsForMapping(controlID, e.Node)...)
	}
	return &relSummary[relatedControlRef]{Count: resp.MappedControls.TotalCount, Items: items}
}

func relatedControlsForMapping(controlID string, m *graphclient.GetMappedControls_MappedControls_Edges_Node) []relatedControlRef {
	if m == nil {
		return nil
	}
	for _, ref := range m.FromControls.Edges {
		if ref != nil && ref.Node != nil && ref.Node.ID == controlID {
			out := make([]relatedControlRef, 0)
			for _, e := range m.ToControls.Edges {
				if e == nil || e.Node == nil {
					continue
				}
				out = append(out, newRelatedControlRef(m, "outgoing", e.Node.ID, e.Node.RefCode))
			}
			for _, e := range m.ToSubcontrols.Edges {
				if e == nil || e.Node == nil {
					continue
				}
				out = append(out, newRelatedControlRef(m, "outgoing", e.Node.ID, e.Node.RefCode))
			}
			return out
		}
	}
	for _, ref := range m.ToControls.Edges {
		if ref != nil && ref.Node != nil && ref.Node.ID == controlID {
			out := make([]relatedControlRef, 0)
			for _, e := range m.FromControls.Edges {
				if e == nil || e.Node == nil {
					continue
				}
				out = append(out, newRelatedControlRef(m, "incoming", e.Node.ID, e.Node.RefCode))
			}
			for _, e := range m.FromSubcontrols.Edges {
				if e == nil || e.Node == nil {
					continue
				}
				out = append(out, newRelatedControlRef(m, "incoming", e.Node.ID, e.Node.RefCode))
			}
			return out
		}
	}
	return nil
}

func newRelatedControlRef(m *graphclient.GetMappedControls_MappedControls_Edges_Node, direction, id, refCode string) relatedControlRef {
	return relatedControlRef{
		MappingID:   m.ID,
		MappingType: openlane.Format(m.MappingType),
		Source:      openlane.Format(m.Source),
		Relation:    openlane.Deref(m.Relation),
		Direction:   direction,
		ID:          id,
		RefCode:     refCode,
	}
}

func mapListMappedControl(n *graphclient.GetMappedControls_MappedControls_Edges_Node) mappedControlItem {
	return mappedControlItem{
		ID:              n.ID,
		MappingType:     openlane.Format(n.MappingType),
		Source:          openlane.Format(n.Source),
		Relation:        openlane.Deref(n.Relation),
		Confidence:      n.Confidence,
		Tags:            n.Tags,
		FromControls:    endpointRefsFromListFromControls(n.FromControls.Edges),
		ToControls:      endpointRefsFromListToControls(n.ToControls.Edges),
		FromSubcontrols: endpointRefsFromListFromSubcontrols(n.FromSubcontrols.Edges),
		ToSubcontrols:   endpointRefsFromListToSubcontrols(n.ToSubcontrols.Edges),
	}
}

func mapGetMappedControl(n graphclient.GetMappedControlByID_MappedControl) mappedControlItem {
	return mappedControlItem{
		ID:              n.ID,
		MappingType:     openlane.Format(n.MappingType),
		Source:          openlane.Format(n.Source),
		Relation:        openlane.Deref(n.Relation),
		Confidence:      n.Confidence,
		Tags:            n.Tags,
		FromControls:    endpointRefsFromGetFromControls(n.FromControls.Edges),
		ToControls:      endpointRefsFromGetToControls(n.ToControls.Edges),
		FromSubcontrols: endpointRefsFromGetFromSubcontrols(n.FromSubcontrols.Edges),
		ToSubcontrols:   endpointRefsFromGetToSubcontrols(n.ToSubcontrols.Edges),
	}
}

func mapCreatedMappedControl(n graphclient.CreateMappedControl_CreateMappedControl_MappedControl) mappedControlItem {
	return mappedControlItem{
		ID:              n.ID,
		MappingType:     openlane.Format(n.MappingType),
		Source:          openlane.Format(n.Source),
		Relation:        openlane.Deref(n.Relation),
		Confidence:      n.Confidence,
		Tags:            n.Tags,
		FromControls:    endpointRefsFromCreateFromControls(n.FromControls.Edges),
		ToControls:      endpointRefsFromCreateToControls(n.ToControls.Edges),
		FromSubcontrols: endpointRefsFromCreateFromSubcontrols(n.FromSubcontrols.Edges),
		ToSubcontrols:   endpointRefsFromCreateToSubcontrols(n.ToSubcontrols.Edges),
	}
}

func mapUpdatedMappedControl(n graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl) mappedControlItem {
	return mappedControlItem{
		ID:              n.ID,
		MappingType:     openlane.Format(n.MappingType),
		Source:          openlane.Format(n.Source),
		Relation:        openlane.Deref(n.Relation),
		Confidence:      n.Confidence,
		Tags:            n.Tags,
		FromControls:    endpointRefsFromUpdateFromControls(n.FromControls.Edges),
		ToControls:      endpointRefsFromUpdateToControls(n.ToControls.Edges),
		FromSubcontrols: endpointRefsFromUpdateFromSubcontrols(n.FromSubcontrols.Edges),
		ToSubcontrols:   endpointRefsFromUpdateToSubcontrols(n.ToSubcontrols.Edges),
	}
}

func endpointRefsFromListFromControls(edges []*graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromListToControls(edges []*graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromListFromSubcontrols(edges []*graphclient.GetMappedControls_MappedControls_Edges_Node_FromSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromListToSubcontrols(edges []*graphclient.GetMappedControls_MappedControls_Edges_Node_ToSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromGetFromControls(edges []*graphclient.GetMappedControlByID_MappedControl_FromControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromGetToControls(edges []*graphclient.GetMappedControlByID_MappedControl_ToControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromGetFromSubcontrols(edges []*graphclient.GetMappedControlByID_MappedControl_FromSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromGetToSubcontrols(edges []*graphclient.GetMappedControlByID_MappedControl_ToSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromCreateFromControls(edges []*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromCreateToControls(edges []*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromCreateFromSubcontrols(edges []*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromCreateToSubcontrols(edges []*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromUpdateFromControls(edges []*graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl_FromControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromUpdateToControls(edges []*graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl_ToControls_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromUpdateFromSubcontrols(edges []*graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl_FromSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}

func endpointRefsFromUpdateToSubcontrols(edges []*graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl_ToSubcontrols_Edges) []mappedControlEndpointRef {
	out := make([]mappedControlEndpointRef, 0, len(edges))
	for _, e := range edges {
		if e == nil || e.Node == nil {
			continue
		}
		out = append(out, mappedControlEndpointRef{ID: e.Node.ID, RefCode: e.Node.RefCode})
	}
	return out
}
