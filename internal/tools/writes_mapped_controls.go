package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type createMappedControlInput struct {
	MappingType           string   `json:"mapping_type" jsonschema:"Mapping type: EQUAL, SUPERSET, SUBSET, INTERSECT, or PARTIAL."`
	Relation              string   `json:"relation,omitempty" jsonschema:"Description of how the controls are related."`
	Confidence            *int64   `json:"confidence,omitempty" jsonschema:"Confidence score from 0 to 100."`
	Source                string   `json:"source,omitempty" jsonschema:"Mapping source: MANUAL, SUGGESTED, or IMPORTED."`
	Tags                  []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
	FromControlIDs        []string `json:"from_control_ids,omitempty" jsonschema:"Source-side control IDs."`
	ToControlIDs          []string `json:"to_control_ids,omitempty" jsonschema:"Target-side control IDs."`
	FromSubcontrolIDs     []string `json:"from_subcontrol_ids,omitempty" jsonschema:"Source-side subcontrol IDs."`
	ToSubcontrolIDs       []string `json:"to_subcontrol_ids,omitempty" jsonschema:"Target-side subcontrol IDs."`
	FromControlRefCodes   []string `json:"from_control_ref_codes,omitempty" jsonschema:"Source-side control ref codes prefixed with standard, e.g. PCI DSS::12.6.2."`
	ToControlRefCodes     []string `json:"to_control_ref_codes,omitempty" jsonschema:"Target-side control ref codes prefixed with standard."`
	FromSubcontrolRefCodes []string `json:"from_subcontrol_ref_codes,omitempty" jsonschema:"Source-side subcontrol ref codes prefixed with standard."`
	ToSubcontrolRefCodes  []string `json:"to_subcontrol_ref_codes,omitempty" jsonschema:"Target-side subcontrol ref codes prefixed with standard."`
}

type updateMappedControlInput struct {
	ID                    string   `json:"id" jsonschema:"MappedControl ID to update."`
	MappingType           string   `json:"mapping_type,omitempty" jsonschema:"Updated mapping type."`
	Relation              string   `json:"relation,omitempty" jsonschema:"Updated relation description."`
	Confidence            *int64   `json:"confidence,omitempty" jsonschema:"Updated confidence score from 0 to 100."`
	Source                string   `json:"source,omitempty" jsonschema:"Updated mapping source."`
	Tags                  []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	AddFromControlIDs     []string `json:"add_from_control_ids,omitempty" jsonschema:"Control IDs to add on the source side."`
	RemoveFromControlIDs  []string `json:"remove_from_control_ids,omitempty" jsonschema:"Control IDs to remove from the source side."`
	AddToControlIDs       []string `json:"add_to_control_ids,omitempty" jsonschema:"Control IDs to add on the target side."`
	RemoveToControlIDs    []string `json:"remove_to_control_ids,omitempty" jsonschema:"Control IDs to remove from the target side."`
	AddFromSubcontrolIDs  []string `json:"add_from_subcontrol_ids,omitempty" jsonschema:"Subcontrol IDs to add on the source side."`
	RemoveFromSubcontrolIDs []string `json:"remove_from_subcontrol_ids,omitempty" jsonschema:"Subcontrol IDs to remove from the source side."`
	AddToSubcontrolIDs    []string `json:"add_to_subcontrol_ids,omitempty" jsonschema:"Subcontrol IDs to add on the target side."`
	RemoveToSubcontrolIDs []string `json:"remove_to_subcontrol_ids,omitempty" jsonschema:"Subcontrol IDs to remove from the target side."`
}

func registerWriteMappedControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_mapped_control_create",
		Title:       "Create an Openlane control mapping",
		Description: "Create a MappedControl that links controls or subcontrols across frameworks or within the org. Provide control IDs or prefixed ref codes on both from and to sides. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createMappedControl)

	addTool(server, &mcp.Tool{
		Name:        "openlane_mapped_control_update",
		Title:       "Update an Openlane control mapping",
		Description: "Update a MappedControl by ID. Use add/remove from/to control or subcontrol ID fields to change linked endpoints. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateMappedControl)
}

func (h *handlers) createMappedControl(ctx context.Context, _ *mcp.CallToolRequest, in createMappedControlInput) (*mcp.CallToolResult, mappedControlItem, error) {
	if strings.TrimSpace(in.MappingType) == "" {
		return nil, mappedControlItem{}, fmt.Errorf("mapping_type is required")
	}
	if !hasMappedControlEndpoints(in.FromControlIDs, in.FromSubcontrolIDs, in.FromControlRefCodes, in.FromSubcontrolRefCodes) {
		return nil, mappedControlItem{}, fmt.Errorf("at least one from_control_ids, from_subcontrol_ids, from_control_ref_codes, or from_subcontrol_ref_codes value is required")
	}
	if !hasMappedControlEndpoints(in.ToControlIDs, in.ToSubcontrolIDs, in.ToControlRefCodes, in.ToSubcontrolRefCodes) {
		return nil, mappedControlItem{}, fmt.Errorf("at least one to_control_ids, to_subcontrol_ids, to_control_ref_codes, or to_subcontrol_ref_codes value is required")
	}
	input := graphclient.CreateMappedControlInput{
		Tags:                  in.Tags,
		MappingType:           mappingType(in.MappingType),
		FromControlIDs:        in.FromControlIDs,
		ToControlIDs:          in.ToControlIDs,
		FromSubcontrolIDs:     in.FromSubcontrolIDs,
		ToSubcontrolIDs:       in.ToSubcontrolIDs,
		FromControlRefCodes:   in.FromControlRefCodes,
		ToControlRefCodes:     in.ToControlRefCodes,
		FromSubcontrolRefCodes: in.FromSubcontrolRefCodes,
		ToSubcontrolRefCodes:  in.ToSubcontrolRefCodes,
	}
	if in.Relation != "" {
		input.Relation = &in.Relation
	}
	if in.Confidence != nil {
		input.Confidence = in.Confidence
	}
	if in.Source != "" {
		input.Source = mappingSource(in.Source)
	}

	resp, err := h.api.CreateMappedControl(ctx, input)
	if err != nil {
		return nil, mappedControlItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedMappedControl(resp.CreateMappedControl.MappedControl), nil
}

func (h *handlers) updateMappedControl(ctx context.Context, _ *mcp.CallToolRequest, in updateMappedControlInput) (*mcp.CallToolResult, mappedControlItem, error) {
	if in.ID == "" {
		return nil, mappedControlItem{}, errIDRequired
	}
	input := graphclient.UpdateMappedControlInput{}
	if in.MappingType != "" {
		input.MappingType = mappingType(in.MappingType)
	}
	if in.Relation != "" {
		input.Relation = &in.Relation
	}
	if in.Confidence != nil {
		input.Confidence = in.Confidence
	}
	if in.Source != "" {
		input.Source = mappingSource(in.Source)
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if len(in.AddFromControlIDs) > 0 {
		input.AddFromControlIDs = in.AddFromControlIDs
	}
	if len(in.RemoveFromControlIDs) > 0 {
		input.RemoveFromControlIDs = in.RemoveFromControlIDs
	}
	if len(in.AddToControlIDs) > 0 {
		input.AddToControlIDs = in.AddToControlIDs
	}
	if len(in.RemoveToControlIDs) > 0 {
		input.RemoveToControlIDs = in.RemoveToControlIDs
	}
	if len(in.AddFromSubcontrolIDs) > 0 {
		input.AddFromSubcontrolIDs = in.AddFromSubcontrolIDs
	}
	if len(in.RemoveFromSubcontrolIDs) > 0 {
		input.RemoveFromSubcontrolIDs = in.RemoveFromSubcontrolIDs
	}
	if len(in.AddToSubcontrolIDs) > 0 {
		input.AddToSubcontrolIDs = in.AddToSubcontrolIDs
	}
	if len(in.RemoveToSubcontrolIDs) > 0 {
		input.RemoveToSubcontrolIDs = in.RemoveToSubcontrolIDs
	}
	if isEmptyUpdateMappedControl(input) {
		return nil, mappedControlItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateMappedControl(ctx, in.ID, input)
	if err != nil {
		return nil, mappedControlItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedMappedControl(resp.UpdateMappedControl.MappedControl), nil
}

func hasMappedControlEndpoints(controlIDs, subcontrolIDs, controlRefCodes, subcontrolRefCodes []string) bool {
	return len(controlIDs) > 0 || len(subcontrolIDs) > 0 || len(controlRefCodes) > 0 || len(subcontrolRefCodes) > 0
}

func mappingType(s string) *enums.MappingType {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return nil
	}
	mt := enums.MappingType(s)
	return &mt
}

func mappingSource(s string) *enums.MappingSource {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return nil
	}
	src := enums.MappingSource(s)
	return &src
}

func isEmptyUpdateMappedControl(in graphclient.UpdateMappedControlInput) bool {
	return in.MappingType == nil &&
		in.Relation == nil &&
		in.Confidence == nil &&
		in.Source == nil &&
		len(in.Tags) == 0 &&
		len(in.AddFromControlIDs) == 0 &&
		len(in.RemoveFromControlIDs) == 0 &&
		len(in.AddToControlIDs) == 0 &&
		len(in.RemoveToControlIDs) == 0 &&
		len(in.AddFromSubcontrolIDs) == 0 &&
		len(in.RemoveFromSubcontrolIDs) == 0 &&
		len(in.AddToSubcontrolIDs) == 0 &&
		len(in.RemoveToSubcontrolIDs) == 0
}
