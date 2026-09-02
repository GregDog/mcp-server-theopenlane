package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type entityItem struct {
	ID                                    string   `json:"id"`
	Name                                  string   `json:"name,omitempty"`
	DisplayName                           string   `json:"display_name,omitempty"`
	Description                           string   `json:"description,omitempty"`
	EntitySourceTypeName                  string   `json:"entity_source_type_name,omitempty"`
	EntityRelationshipStateName           string   `json:"entity_relationship_state_name,omitempty"`
	EntitySecurityQuestionnaireStatusName string   `json:"entity_security_questionnaire_status_name,omitempty"`
	EnvironmentName                       string   `json:"environment_name,omitempty"`
	Tier                                  string   `json:"tier,omitempty"`
	RiskRating                            string   `json:"risk_rating,omitempty"`
	RiskScore                             *int64   `json:"risk_score,omitempty"`
	ApprovedForUse                        *bool    `json:"approved_for_use,omitempty"`
	ContractStartDate                     string   `json:"contract_start_date,omitempty"`
	ContractEndDate                       string   `json:"contract_end_date,omitempty"`
	Tags                                  []string `json:"tags,omitempty"`
}

func registerEntities(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_entities_list",
		Title:       "List Openlane entities",
		Description: "List entities in the configured Openlane organization. Entities represent vendors and other third parties. Results are paginated.",
		Annotations: readOnly(),
	}, h.listEntities)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_entity_get",
		Title:       "Get an Openlane entity",
		Description: "Get a single entity by ID. Entities represent vendors and other third parties.",
		Annotations: readOnly(),
	}, h.getEntity)
}

func (h *handlers) listEntities(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[entityItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetEntities(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[entityItem]{}, openlane.APIError(err)
	}
	items := make([]entityItem, 0, len(resp.Entities.Edges))
	for _, e := range resp.Entities.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListEntity(e.Node))
	}
	return nil, openlane.Page[entityItem]{
		Items:      items,
		NextCursor: resp.Entities.PageInfo.EndCursor,
		HasMore:    resp.Entities.PageInfo.HasNextPage,
		TotalCount: resp.Entities.TotalCount,
	}, nil
}

func (h *handlers) getEntity(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, entityItem, error) {
	if in.ID == "" {
		return nil, entityItem{}, errIDRequired
	}
	resp, err := h.api.GetEntityByID(ctx, in.ID)
	if err != nil {
		return nil, entityItem{}, openlane.APIError(err)
	}
	return nil, mapGetEntity(resp.Entity), nil
}

func mapListEntity(n *graphclient.GetEntities_Entities_Edges_Node) entityItem {
	return entityItem{
		ID:                                    n.ID,
		Name:                                  openlane.Deref(n.Name),
		DisplayName:                           openlane.Deref(n.DisplayName),
		Description:                           openlane.Deref(n.Description),
		EntitySourceTypeName:                  openlane.Deref(n.EntitySourceTypeName),
		EntityRelationshipStateName:           openlane.Deref(n.EntityRelationshipStateName),
		EntitySecurityQuestionnaireStatusName: openlane.Deref(n.EntitySecurityQuestionnaireStatusName),
		EnvironmentName:                       openlane.Deref(n.EnvironmentName),
		Tier:                                  openlane.Format(n.Tier),
		RiskRating:                            openlane.Deref(n.RiskRating),
		RiskScore:                             n.RiskScore,
	}
}

func mapGetEntity(e graphclient.GetEntityByID_Entity) entityItem {
	approved := e.ApprovedForUse
	return entityItem{
		ID:                                    e.ID,
		Name:                                  openlane.Deref(e.Name),
		DisplayName:                           openlane.Deref(e.DisplayName),
		Description:                           openlane.Deref(e.Description),
		EntitySourceTypeName:                  openlane.Deref(e.EntitySourceTypeName),
		EntityRelationshipStateName:           openlane.Deref(e.EntityRelationshipStateName),
		EntitySecurityQuestionnaireStatusName: openlane.Deref(e.EntitySecurityQuestionnaireStatusName),
		EnvironmentName:                       openlane.Deref(e.EnvironmentName),
		Tier:                                  openlane.Format(e.Tier),
		RiskRating:                            openlane.Deref(e.RiskRating),
		RiskScore:                             e.RiskScore,
		ApprovedForUse:                        approved,
		ContractStartDate:                     openlane.Format(e.ContractStartDate),
		ContractEndDate:                       openlane.Format(e.ContractEndDate),
		Tags:                                  e.Tags,
	}
}
