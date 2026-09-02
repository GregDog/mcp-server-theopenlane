package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type riskItem struct {
	ID            string                      `json:"id"`
	DisplayID     string                      `json:"display_id,omitempty"`
	Name          string                      `json:"name,omitempty"`
	Status        string                      `json:"status,omitempty"`
	Impact        string                      `json:"impact,omitempty"`
	Likelihood    string                      `json:"likelihood,omitempty"`
	Score         *int64                      `json:"score,omitempty"`
	Details       string                      `json:"details,omitempty"`
	Mitigation    string                      `json:"mitigation,omitempty"`
	BusinessCosts string                      `json:"business_costs,omitempty"`
	RiskKindName  string                      `json:"risk_kind_name,omitempty"`
	OwnerID       string                      `json:"owner_id,omitempty"`
	StakeholderID string                      `json:"stakeholder_id,omitempty"`
	DelegateID    string                      `json:"delegate_id,omitempty"`
	Tags          []string                    `json:"tags,omitempty"`
	Controls      *relSummary[controlRef]     `json:"controls,omitempty"`
	Programs      *relSummary[idNameRef]      `json:"programs,omitempty"`
	Findings      *relSummary[findingRef]     `json:"findings,omitempty"`
	Remediations  *relSummary[remediationRef] `json:"remediations,omitempty"`
}

func registerRisks(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_risks_list",
		Title:       "List Openlane risks",
		Description: "List risks in the configured Openlane organization. Supports program, entity, control, and status filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.listRisks)

	addTool(server, &mcp.Tool{
		Name:        "openlane_risk_get",
		Title:       "Get an Openlane risk",
		Description: "Get a single risk by ID with scoring/treatment fields and bounded summaries of linked controls, programs, entities, assets, findings, and remediations.",
		Annotations: readOnly(),
	}, h.getRisk)
}

func (h *handlers) listRisks(ctx context.Context, _ *mcp.CallToolRequest, in riskListInput) (*mcp.CallToolResult, openlane.Page[riskItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetRisks(ctx, &first, after, buildRiskWhere(in))
	if err != nil {
		return nil, openlane.Page[riskItem]{}, openlane.APIError(err)
	}
	items := make([]riskItem, 0, len(resp.Risks.Edges))
	for _, e := range resp.Risks.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListRisk(e.Node))
	}
	return nil, openlane.Page[riskItem]{
		Items:      items,
		NextCursor: resp.Risks.PageInfo.EndCursor,
		HasMore:    resp.Risks.PageInfo.HasNextPage,
		TotalCount: resp.Risks.TotalCount,
	}, nil
}

func (h *handlers) getRisk(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, riskItem, error) {
	if in.ID == "" {
		return nil, riskItem{}, errIDRequired
	}
	resp, err := h.api.GetRiskByID(ctx, in.ID)
	if err != nil {
		return nil, riskItem{}, openlane.APIError(err)
	}
	item := mapGetRisk(resp.Risk)
	riskWhere := []*graphclient.RiskWhereInput{{ID: &in.ID}}
	runRelJobs(
		func() { item.Controls = h.fetchControls(ctx, &graphclient.ControlWhereInput{HasRisksWith: riskWhere}) },
		func() { item.Programs = h.fetchPrograms(ctx, &graphclient.ProgramWhereInput{HasRisksWith: riskWhere}) },
		func() { item.Findings = h.fetchFindings(ctx, &graphclient.FindingWhereInput{HasRisksWith: riskWhere}) },
		func() {
			item.Remediations = h.fetchRemediations(ctx, &graphclient.RemediationWhereInput{HasRisksWith: riskWhere})
		},
	)
	return nil, item, nil
}

func mapListRisk(n *graphclient.GetRisks_Risks_Edges_Node) riskItem {
	return riskItem{
		ID:         n.ID,
		DisplayID:  n.DisplayID,
		Name:       n.Name,
		Status:     openlane.Format(n.Status),
		Impact:     openlane.Format(n.Impact),
		Likelihood: openlane.Format(n.Likelihood),
		Score:      n.Score,
	}
}

func mapGetRisk(r graphclient.GetRiskByID_Risk) riskItem {
	return riskItem{
		ID:            r.ID,
		DisplayID:     r.DisplayID,
		Name:          r.Name,
		Status:        openlane.Format(r.Status),
		Impact:        openlane.Format(r.Impact),
		Likelihood:    openlane.Format(r.Likelihood),
		Score:         r.Score,
		Details:       openlane.Deref(r.Details),
		Mitigation:    openlane.Deref(r.Mitigation),
		BusinessCosts: openlane.Deref(r.BusinessCosts),
		RiskKindName:  openlane.Deref(r.RiskKindName),
		OwnerID:       openlane.Deref(r.OwnerID),
		StakeholderID: openlane.Deref(r.StakeholderID),
		DelegateID:    openlane.Deref(r.DelegateID),
		Tags:          r.Tags,
	}
}
