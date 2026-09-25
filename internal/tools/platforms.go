package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type platformDiagramRef struct {
	ID       string `json:"id"`
	Filename string `json:"filename,omitempty"`
}

type platformRelItem struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	RefCode   string `json:"ref_code,omitempty"`
	Title     string `json:"title,omitempty"`
	Status    string `json:"status,omitempty"`
	DisplayID string `json:"display_id,omitempty"`
}

type platformItem struct {
	ID                             string                       `json:"id"`
	DisplayID                      string                       `json:"display_id,omitempty"`
	Name                           string                       `json:"name,omitempty"`
	Description                    string                       `json:"description,omitempty"`
	BusinessPurpose                string                       `json:"business_purpose,omitempty"`
	ScopeStatement                 string                       `json:"scope_statement,omitempty"`
	TrustBoundaryDescription       string                       `json:"trust_boundary_description,omitempty"`
	DataFlowSummary                string                       `json:"data_flow_summary,omitempty"`
	Status                         string                       `json:"status,omitempty"`
	PhysicalLocation               string                       `json:"physical_location,omitempty"`
	Region                         string                       `json:"region,omitempty"`
	SourceType                     string                       `json:"source_type,omitempty"`
	SourceIdentifier               string                       `json:"source_identifier,omitempty"`
	CostCenter                     string                       `json:"cost_center,omitempty"`
	EstimatedMonthlyCost           *float64                     `json:"estimated_monthly_cost,omitempty"`
	PurchaseDate                   string                       `json:"purchase_date,omitempty"`
	ExternalReferenceID            string                       `json:"external_reference_id,omitempty"`
	Tags                           []string                     `json:"tags,omitempty"`
	PlatformKindName               string                       `json:"platform_kind_name,omitempty"`
	PlatformDataClassificationName string                       `json:"platform_data_classification_name,omitempty"`
	EnvironmentName                string                       `json:"environment_name,omitempty"`
	ScopeName                      string                       `json:"scope_name,omitempty"`
	AccessModelName                string                       `json:"access_model_name,omitempty"`
	EncryptionStatusName           string                       `json:"encryption_status_name,omitempty"`
	SecurityTierName               string                       `json:"security_tier_name,omitempty"`
	CriticalityName                string                       `json:"criticality_name,omitempty"`
	BusinessOwner                  *platformOwnerSummary        `json:"business_owner,omitempty"`
	TechnicalOwner                 *platformOwnerSummary        `json:"technical_owner,omitempty"`
	SecurityOwner                  *platformOwnerSummary        `json:"security_owner,omitempty"`
	InternalOwner                  *platformOwnerSummary        `json:"internal_owner,omitempty"`
	Assets                         *relSummary[platformRelItem] `json:"assets,omitempty"`
	Entities                       *relSummary[platformRelItem] `json:"entities,omitempty"`
	SourceAssets                   *relSummary[platformRelItem] `json:"source_assets,omitempty"`
	SourceEntities                 *relSummary[platformRelItem] `json:"source_entities,omitempty"`
	OutOfScopeAssets               *relSummary[platformRelItem] `json:"out_of_scope_assets,omitempty"`
	OutOfScopeVendors              *relSummary[platformRelItem] `json:"out_of_scope_vendors,omitempty"`
	Controls                       *relSummary[platformRelItem] `json:"controls,omitempty"`
	Evidence                       *relSummary[platformRelItem] `json:"evidence,omitempty"`
	ApplicableFrameworks           *relSummary[platformRelItem] `json:"applicable_frameworks,omitempty"`
	Risks                          *relSummary[platformRelItem] `json:"risks,omitempty"`
	Tasks                          *relSummary[platformRelItem] `json:"tasks,omitempty"`
	ArchitectureDiagrams           []platformDiagramRef         `json:"architecture_diagrams,omitempty"`
	DataFlowDiagrams               []platformDiagramRef         `json:"data_flow_diagrams,omitempty"`
	TrustBoundaryDiagrams          []platformDiagramRef         `json:"trust_boundary_diagrams,omitempty"`
	CreatedAt                      string                       `json:"created_at,omitempty"`
	UpdatedAt                      string                       `json:"updated_at,omitempty"`
}

func registerPlatforms(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_platforms_list",
		Title:       "List Openlane platforms",
		Description: "List platforms (system boundaries) in the configured Openlane organization. Supports name, display id, status, environment, region, and criticality filters. Results are paginated.",
		Annotations: readOnly(),
	}, h.listPlatforms)

	addTool(server, &mcp.Tool{
		Name:        "openlane_platform_get",
		Title:       "Get an Openlane platform",
		Description: "Get a single platform by ID with narrative fields, owner roles, bounded scope links, and diagram metadata.",
		Annotations: readOnly(),
	}, h.getPlatform)
}

func (h *handlers) listPlatforms(ctx context.Context, _ *mcp.CallToolRequest, in platformListInput) (*mcp.CallToolResult, openlane.Page[platformItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetPlatforms(ctx, &first, after, buildPlatformWhere(in))
	if err != nil {
		return nil, openlane.Page[platformItem]{}, openlane.APIError(err)
	}
	items := make([]platformItem, 0, len(resp.Platforms.Edges))
	cache := make(map[string]groupAssigneeSummary)
	for _, e := range resp.Platforms.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		item := mapListPlatform(e.Node)
		if err := h.enrichListPlatformOwners(ctx, &item, e.Node, cache); err != nil {
			return nil, openlane.Page[platformItem]{}, err
		}
		items = append(items, item)
	}
	return nil, openlane.Page[platformItem]{
		Items:      items,
		NextCursor: resp.Platforms.PageInfo.EndCursor,
		HasMore:    resp.Platforms.PageInfo.HasNextPage,
		TotalCount: resp.Platforms.TotalCount,
	}, nil
}

func (h *handlers) getPlatform(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, platformItem, error) {
	if in.ID == "" {
		return nil, platformItem{}, errIDRequired
	}
	detail, err := h.api.GetPlatformDetail(ctx, in.ID)
	if err != nil {
		return nil, platformItem{}, openlane.APIError(err)
	}
	return nil, mapPlatformDetail(*detail), nil
}

func (h *handlers) enrichListPlatformOwners(ctx context.Context, item *platformItem, n *graphclient.GetPlatforms_Platforms_Edges_Node, cache map[string]groupAssigneeSummary) error {
	business, err := h.enrichPlatformOwnerFromScalars(ctx, openlane.Deref(n.BusinessOwnerUserID), openlane.Deref(n.BusinessOwnerGroupID), openlane.Deref(n.BusinessOwner), cache)
	if err != nil {
		return err
	}
	technical, err := h.enrichPlatformOwnerFromScalars(ctx, openlane.Deref(n.TechnicalOwnerUserID), openlane.Deref(n.TechnicalOwnerGroupID), openlane.Deref(n.TechnicalOwner), cache)
	if err != nil {
		return err
	}
	security, err := h.enrichPlatformOwnerFromScalars(ctx, openlane.Deref(n.SecurityOwnerUserID), openlane.Deref(n.SecurityOwnerGroupID), openlane.Deref(n.SecurityOwner), cache)
	if err != nil {
		return err
	}
	internal, err := h.enrichPlatformOwnerFromScalars(ctx, openlane.Deref(n.InternalOwnerUserID), openlane.Deref(n.InternalOwnerGroupID), openlane.Deref(n.InternalOwner), cache)
	if err != nil {
		return err
	}
	item.BusinessOwner = business
	item.TechnicalOwner = technical
	item.SecurityOwner = security
	item.InternalOwner = internal
	return nil
}

func mapListPlatform(n *graphclient.GetPlatforms_Platforms_Edges_Node) platformItem {
	return platformItem{
		ID:              n.ID,
		DisplayID:       n.DisplayID,
		Name:            n.Name,
		Status:          openlane.Format(n.Status),
		EnvironmentName: openlane.Deref(n.EnvironmentName),
		CriticalityName: openlane.Deref(n.CriticalityName),
		Region:          openlane.Deref(n.Region),
	}
}

func mapPlatformDetail(d openlane.PlatformDetail) platformItem {
	return platformItem{
		ID:                             d.ID,
		DisplayID:                      d.DisplayID,
		Name:                           d.Name,
		Description:                    d.Description,
		BusinessPurpose:                d.BusinessPurpose,
		ScopeStatement:                 d.ScopeStatement,
		TrustBoundaryDescription:       d.TrustBoundaryDescription,
		DataFlowSummary:                d.DataFlowSummary,
		Status:                         openlane.Format(d.Status),
		PhysicalLocation:               d.PhysicalLocation,
		Region:                         d.Region,
		SourceType:                     openlane.Format(d.SourceType),
		SourceIdentifier:               d.SourceIdentifier,
		CostCenter:                     d.CostCenter,
		EstimatedMonthlyCost:           d.EstimatedMonthlyCost,
		PurchaseDate:                   openlane.Format(d.PurchaseDate),
		ExternalReferenceID:            d.ExternalReferenceID,
		Tags:                           d.Tags,
		PlatformKindName:               d.PlatformKindName,
		PlatformDataClassificationName: d.PlatformDataClassificationName,
		EnvironmentName:                d.EnvironmentName,
		ScopeName:                      d.ScopeName,
		AccessModelName:                d.AccessModelName,
		EncryptionStatusName:           d.EncryptionStatusName,
		SecurityTierName:               d.SecurityTierName,
		CriticalityName:                d.CriticalityName,
		BusinessOwner:                  mapPlatformOwnerRole(d.BusinessOwner),
		TechnicalOwner:                 mapPlatformOwnerRole(d.TechnicalOwner),
		SecurityOwner:                  mapPlatformOwnerRole(d.SecurityOwner),
		InternalOwner:                  mapPlatformOwnerRole(d.InternalOwner),
		Assets:                         mapPlatformRelSummary(d.Assets),
		Entities:                       mapPlatformRelSummary(d.Entities),
		SourceAssets:                   mapPlatformRelSummary(d.SourceAssets),
		SourceEntities:                 mapPlatformRelSummary(d.SourceEntities),
		OutOfScopeAssets:               mapPlatformRelSummary(d.OutOfScopeAssets),
		OutOfScopeVendors:              mapPlatformRelSummary(d.OutOfScopeVendors),
		Controls:                       mapPlatformRelSummary(d.Controls),
		Evidence:                       mapPlatformRelSummary(d.Evidence),
		ApplicableFrameworks:           mapPlatformRelSummary(d.ApplicableFrameworks),
		Risks:                          mapPlatformRelSummary(d.Risks),
		Tasks:                          mapPlatformRelSummary(d.Tasks),
		ArchitectureDiagrams:           mapPlatformDiagrams(d.ArchitectureDiagrams),
		DataFlowDiagrams:               mapPlatformDiagrams(d.DataFlowDiagrams),
		TrustBoundaryDiagrams:          mapPlatformDiagrams(d.TrustBoundaryDiagrams),
		CreatedAt:                      openlane.Format(d.CreatedAt),
		UpdatedAt:                      openlane.Format(d.UpdatedAt),
	}
}

func mapPlatformRelSummary(s openlane.PlatformRelSummary) *relSummary[platformRelItem] {
	if s.Count == 0 && len(s.Items) == 0 {
		return nil
	}
	items := make([]platformRelItem, 0, len(s.Items))
	for _, n := range s.Items {
		items = append(items, platformRelItem{
			ID:        n.ID,
			Name:      n.Name,
			RefCode:   n.RefCode,
			Title:     n.Title,
			Status:    n.Status,
			DisplayID: n.DisplayID,
		})
	}
	return &relSummary[platformRelItem]{Count: s.Count, Items: items}
}

func mapPlatformDiagrams(diagrams []openlane.PlatformDiagramRef) []platformDiagramRef {
	if len(diagrams) == 0 {
		return nil
	}
	out := make([]platformDiagramRef, 0, len(diagrams))
	for _, d := range diagrams {
		out = append(out, platformDiagramRef{ID: d.ID, Filename: d.Filename})
	}
	return out
}
