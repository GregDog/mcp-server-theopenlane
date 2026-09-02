package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type entityItem struct {
	ID                                    string                  `json:"id"`
	Name                                  string                  `json:"name,omitempty"`
	DisplayName                           string                  `json:"display_name,omitempty"`
	Description                           string                  `json:"description,omitempty"`
	EntityTypeID                          string                  `json:"entity_type_id,omitempty"`
	EntitySourceTypeName                  string                  `json:"entity_source_type_name,omitempty"`
	EntityRelationshipStateName           string                  `json:"entity_relationship_state_name,omitempty"`
	EntitySecurityQuestionnaireStatusName string                  `json:"entity_security_questionnaire_status_name,omitempty"`
	EnvironmentName                       string                  `json:"environment_name,omitempty"`
	Tier                                  string                  `json:"tier,omitempty"`
	RiskRating                            string                  `json:"risk_rating,omitempty"`
	RiskScore                             *int64                  `json:"risk_score,omitempty"`
	ApprovedForUse                        *bool                   `json:"approved_for_use,omitempty"`
	HasSoc2                               *bool                   `json:"has_soc2,omitempty"`
	Soc2PeriodEnd                         string                  `json:"soc2_period_end,omitempty"`
	SsoEnforced                           *bool                   `json:"sso_enforced,omitempty"`
	MfaSupported                          *bool                   `json:"mfa_supported,omitempty"`
	MfaEnforced                           *bool                   `json:"mfa_enforced,omitempty"`
	LastReviewedAt                        string                  `json:"last_reviewed_at,omitempty"`
	NextReviewAt                          string                  `json:"next_review_at,omitempty"`
	ContractStartDate                     string                  `json:"contract_start_date,omitempty"`
	ContractEndDate                       string                  `json:"contract_end_date,omitempty"`
	ContractRenewalAt                     string                  `json:"contract_renewal_at,omitempty"`
	AutoRenews                            *bool                   `json:"auto_renews,omitempty"`
	TerminationNoticeDays                 *int64                  `json:"termination_notice_days,omitempty"`
	AnnualSpend                           *float64                `json:"annual_spend,omitempty"`
	SpendCurrency                         string                  `json:"spend_currency,omitempty"`
	BillingModel                          string                  `json:"billing_model,omitempty"`
	RenewalRisk                           string                  `json:"renewal_risk,omitempty"`
	InternalOwner                         string                  `json:"internal_owner,omitempty"`
	InternalOwnerUserID                   string                  `json:"internal_owner_user_id,omitempty"`
	InternalOwnerGroupID                  string                  `json:"internal_owner_group_id,omitempty"`
	Tags                                  []string                `json:"tags,omitempty"`
	Assets                                *relSummary[idNameRef]  `json:"assets,omitempty"`
	Risks                                 *relSummary[idNameRef]  `json:"risks,omitempty"`
	Findings                              *relSummary[findingRef] `json:"findings,omitempty"`
}

func registerEntities(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_entities_list",
		Title:       "List Openlane entities",
		Description: "List entities in the configured Openlane organization. Entities represent vendors and other third parties. Supports risk, tier, questionnaire, review, and security filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.listEntities)

	addTool(server, &mcp.Tool{
		Name:        "openlane_entity_get",
		Title:       "Get an Openlane entity",
		Description: "Get a single entity by ID with vendor/security/commercial fields and bounded summaries of linked assets, risks, and findings.",
		Annotations: readOnly(),
	}, h.getEntity)
}

func (h *handlers) listEntities(ctx context.Context, _ *mcp.CallToolRequest, in entityListInput) (*mcp.CallToolResult, openlane.Page[entityItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetEntities(ctx, &first, after, buildEntityWhere(in))
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
	item := mapGetEntity(resp.Entity)
	entityWhere := []*graphclient.EntityWhereInput{{ID: &in.ID}}
	runRelJobs(
		func() { item.Assets = h.fetchAssets(ctx, &graphclient.AssetWhereInput{HasEntitiesWith: entityWhere}) },
		func() { item.Risks = h.fetchRisks(ctx, &graphclient.RiskWhereInput{HasEntitiesWith: entityWhere}) },
		func() {
			item.Findings = h.fetchFindings(ctx, &graphclient.FindingWhereInput{HasEntitiesWith: entityWhere})
		},
	)
	return nil, item, nil
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
		ApprovedForUse:                        n.ApprovedForUse,
		HasSoc2:                               n.HasSoc2,
		SsoEnforced:                           n.SsoEnforced,
		MfaEnforced:                           n.MfaEnforced,
		NextReviewAt:                          openlane.Format(n.NextReviewAt),
	}
}

func mapGetEntity(e graphclient.GetEntityByID_Entity) entityItem {
	return entityItem{
		ID:                                    e.ID,
		Name:                                  openlane.Deref(e.Name),
		DisplayName:                           openlane.Deref(e.DisplayName),
		Description:                           openlane.Deref(e.Description),
		EntityTypeID:                          openlane.Deref(e.EntityTypeID),
		EntitySourceTypeName:                  openlane.Deref(e.EntitySourceTypeName),
		EntityRelationshipStateName:           openlane.Deref(e.EntityRelationshipStateName),
		EntitySecurityQuestionnaireStatusName: openlane.Deref(e.EntitySecurityQuestionnaireStatusName),
		EnvironmentName:                       openlane.Deref(e.EnvironmentName),
		Tier:                                  openlane.Format(e.Tier),
		RiskRating:                            openlane.Deref(e.RiskRating),
		RiskScore:                             e.RiskScore,
		ApprovedForUse:                        e.ApprovedForUse,
		HasSoc2:                               e.HasSoc2,
		Soc2PeriodEnd:                         openlane.Format(e.Soc2PeriodEnd),
		SsoEnforced:                           e.SsoEnforced,
		MfaSupported:                          e.MfaSupported,
		MfaEnforced:                           e.MfaEnforced,
		LastReviewedAt:                        openlane.Format(e.LastReviewedAt),
		NextReviewAt:                          openlane.Format(e.NextReviewAt),
		ContractStartDate:                     openlane.Format(e.ContractStartDate),
		ContractEndDate:                       openlane.Format(e.ContractEndDate),
		ContractRenewalAt:                     openlane.Format(e.ContractRenewalAt),
		AutoRenews:                            e.AutoRenews,
		TerminationNoticeDays:                 e.TerminationNoticeDays,
		AnnualSpend:                           e.AnnualSpend,
		SpendCurrency:                         openlane.Deref(e.SpendCurrency),
		BillingModel:                          openlane.Deref(e.BillingModel),
		RenewalRisk:                           openlane.Deref(e.RenewalRisk),
		InternalOwner:                         openlane.Deref(e.InternalOwner),
		InternalOwnerUserID:                   openlane.Deref(e.InternalOwnerUserID),
		InternalOwnerGroupID:                  openlane.Deref(e.InternalOwnerGroupID),
		Tags:                                  e.Tags,
	}
}
