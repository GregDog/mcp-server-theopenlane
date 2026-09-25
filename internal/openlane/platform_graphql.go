package openlane

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"
)

// PlatformDiagramUploads holds optional diagram file uploads for create/update.
type PlatformDiagramUploads struct {
	Architecture  []*graphql.Upload
	DataFlow      []*graphql.Upload
	TrustBoundary []*graphql.Upload
}

func (d PlatformDiagramUploads) hasUploads() bool {
	return len(d.Architecture) > 0 || len(d.DataFlow) > 0 || len(d.TrustBoundary) > 0
}

// PlatformOwnerRole is a resolved business, technical, security, or internal owner.
type PlatformOwnerRole struct {
	Kind             string
	UserID           string
	UserEmail        string
	UserDisplayName  string
	GroupID          string
	GroupName        string
	GroupDisplayName string
	Name             string
}

// PlatformDiagramRef is diagram metadata without file bytes.
type PlatformDiagramRef struct {
	ID       string
	Filename string
}

// PlatformRelNode is a compact linked object on a platform.
type PlatformRelNode struct {
	ID        string
	Name      string
	RefCode   string
	Title     string
	Status    string
	DisplayID string
}

// PlatformRelSummary is a bounded relationship collection.
type PlatformRelSummary struct {
	Count int64
	Items []PlatformRelNode
}

// PlatformDetail is a platform with owner edges and bounded relationships.
type PlatformDetail struct {
	ID                             string
	DisplayID                      string
	Name                           string
	Description                    string
	BusinessPurpose                string
	ScopeStatement                 string
	TrustBoundaryDescription       string
	DataFlowSummary                string
	Status                         enums.PlatformStatus
	PhysicalLocation               string
	Region                         string
	SourceType                     enums.SourceType
	SourceIdentifier               string
	CostCenter                     string
	EstimatedMonthlyCost           *float64
	PurchaseDate                   *models.DateTime
	ExternalReferenceID            string
	Tags                           []string
	PlatformKindName               string
	PlatformDataClassificationName string
	EnvironmentName                string
	ScopeName                      string
	AccessModelName                string
	EncryptionStatusName           string
	SecurityTierName               string
	CriticalityName                string
	BusinessOwner                  PlatformOwnerRole
	TechnicalOwner                 PlatformOwnerRole
	SecurityOwner                  PlatformOwnerRole
	InternalOwner                  PlatformOwnerRole
	Assets                         PlatformRelSummary
	Entities                       PlatformRelSummary
	SourceAssets                   PlatformRelSummary
	SourceEntities                 PlatformRelSummary
	OutOfScopeAssets               PlatformRelSummary
	OutOfScopeVendors              PlatformRelSummary
	Controls                       PlatformRelSummary
	Evidence                       PlatformRelSummary
	ApplicableFrameworks           PlatformRelSummary
	Risks                          PlatformRelSummary
	Tasks                          PlatformRelSummary
	ArchitectureDiagrams           []PlatformDiagramRef
	DataFlowDiagrams               []PlatformDiagramRef
	TrustBoundaryDiagrams          []PlatformDiagramRef
	CreatedAt                      *time.Time
	UpdatedAt                      *time.Time
}

const platformDetailFields = `
      id
      displayID
      name
      description
      businessPurpose
      scopeStatement
      trustBoundaryDescription
      dataFlowSummary
      status
      physicalLocation
      region
      sourceType
      sourceIdentifier
      costCenter
      estimatedMonthlyCost
      purchaseDate
      externalReferenceID
      tags
      platformKindName
      platformDataClassificationName
      environmentName
      scopeName
      accessModelName
      encryptionStatusName
      securityTierName
      criticalityName
      businessOwner
      technicalOwner
      securityOwner
      internalOwner
      businessOwnerUser { id displayName email }
      businessOwnerGroup { id name displayName isManaged }
      technicalOwnerUser { id displayName email }
      technicalOwnerGroup { id name displayName isManaged }
      securityOwnerUser { id displayName email }
      securityOwnerGroup { id name displayName isManaged }
      internalOwnerUser { id displayName email }
      internalOwnerGroup { id name displayName isManaged }
      assets(first: 8) { totalCount edges { node { id name } } }
      entities(first: 8) { totalCount edges { node { id name displayName } } }
      sourceAssets(first: 8) { totalCount edges { node { id name } } }
      sourceEntities(first: 8) { totalCount edges { node { id name displayName } } }
      outOfScopeAssets(first: 8) { totalCount edges { node { id name } } }
      outOfScopeVendors(first: 8) { totalCount edges { node { id name displayName } } }
      controls(first: 8) { totalCount edges { node { id refCode title status } } }
      evidence(first: 8) { totalCount edges { node { id name status } } }
      applicableFrameworks(first: 8) { totalCount edges { node { id name framework } } }
      risks(first: 8) { totalCount edges { node { id name status } } }
      tasks(first: 8) { totalCount edges { node { id displayID title status } } }
      architectureDiagrams(first: 8) { totalCount edges { node { id providedFileName } } }
      dataFlowDiagrams(first: 8) { totalCount edges { node { id providedFileName } } }
      trustBoundaryDiagrams(first: 8) { totalCount edges { node { id providedFileName } } }
      createdAt
      updatedAt
`

const platformDetailQuery = `query PlatformDetail($platformId: ID!) {
  platform(id: $platformId) {
` + platformDetailFields + `
  }
}`

const createPlatformWithDiagramsMutation = `mutation CreatePlatformWithDiagrams($input: CreatePlatformInput!, $architectureDiagrams: [Upload!], $dataFlowDiagrams: [Upload!], $trustBoundaryDiagrams: [Upload!]) {
  createPlatform(input: $input, architectureDiagrams: $architectureDiagrams, dataFlowDiagrams: $dataFlowDiagrams, trustBoundaryDiagrams: $trustBoundaryDiagrams) {
    platform {
` + platformDetailFields + `
    }
  }
}`

const updatePlatformWithDiagramsMutation = `mutation UpdatePlatformWithDiagrams($updatePlatformId: ID!, $input: UpdatePlatformInput!, $architectureDiagrams: [Upload!], $dataFlowDiagrams: [Upload!], $trustBoundaryDiagrams: [Upload!]) {
  updatePlatform(id: $updatePlatformId, input: $input, architectureDiagrams: $architectureDiagrams, dataFlowDiagrams: $dataFlowDiagrams, trustBoundaryDiagrams: $trustBoundaryDiagrams) {
    platform {
` + platformDetailFields + `
    }
  }
}`

type platformUserNode struct {
	ID          string  `json:"id"`
	DisplayName *string `json:"displayName,omitempty"`
	Email       *string `json:"email,omitempty"`
}

type platformGroupNode struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	IsManaged   *bool  `json:"isManaged,omitempty"`
}

type platformRelConn struct {
	TotalCount int64 `json:"totalCount"`
	Edges      []struct {
		Node *platformRelNode `json:"node"`
	} `json:"edges"`
}

type platformRelNode struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	RefCode     string  `json:"refCode,omitempty"`
	Title       *string `json:"title,omitempty"`
	Status      any     `json:"status,omitempty"`
	Framework   *string `json:"framework,omitempty"`
	DisplayID   string  `json:"displayID,omitempty"`
}

type platformDiagramConn struct {
	TotalCount int64 `json:"totalCount"`
	Edges      []struct {
		Node *struct {
			ID               string `json:"id"`
			ProvidedFileName string `json:"providedFileName"`
		} `json:"node"`
	} `json:"edges"`
}

type platformDetailNode struct {
	ID                             string               `json:"id"`
	DisplayID                      string               `json:"displayID"`
	Name                           string               `json:"name"`
	Description                    *string              `json:"description,omitempty"`
	BusinessPurpose                *string              `json:"businessPurpose,omitempty"`
	ScopeStatement                 *string              `json:"scopeStatement,omitempty"`
	TrustBoundaryDescription       *string              `json:"trustBoundaryDescription,omitempty"`
	DataFlowSummary                *string              `json:"dataFlowSummary,omitempty"`
	Status                         enums.PlatformStatus `json:"status"`
	PhysicalLocation               *string              `json:"physicalLocation,omitempty"`
	Region                         *string              `json:"region,omitempty"`
	SourceType                     enums.SourceType     `json:"sourceType"`
	SourceIdentifier               *string              `json:"sourceIdentifier,omitempty"`
	CostCenter                     *string              `json:"costCenter,omitempty"`
	EstimatedMonthlyCost           *float64             `json:"estimatedMonthlyCost,omitempty"`
	PurchaseDate                   *models.DateTime     `json:"purchaseDate,omitempty"`
	ExternalReferenceID            *string              `json:"externalReferenceID,omitempty"`
	Tags                           []string             `json:"tags,omitempty"`
	PlatformKindName               *string              `json:"platformKindName,omitempty"`
	PlatformDataClassificationName *string              `json:"platformDataClassificationName,omitempty"`
	EnvironmentName                *string              `json:"environmentName,omitempty"`
	ScopeName                      *string              `json:"scopeName,omitempty"`
	AccessModelName                *string              `json:"accessModelName,omitempty"`
	EncryptionStatusName           *string              `json:"encryptionStatusName,omitempty"`
	SecurityTierName               *string              `json:"securityTierName,omitempty"`
	CriticalityName                *string              `json:"criticalityName,omitempty"`
	BusinessOwner                  *string              `json:"businessOwner,omitempty"`
	TechnicalOwner                 *string              `json:"technicalOwner,omitempty"`
	SecurityOwner                  *string              `json:"securityOwner,omitempty"`
	InternalOwner                  *string              `json:"internalOwner,omitempty"`
	BusinessOwnerUser              *platformUserNode    `json:"businessOwnerUser,omitempty"`
	BusinessOwnerGroup             *platformGroupNode   `json:"businessOwnerGroup,omitempty"`
	TechnicalOwnerUser             *platformUserNode    `json:"technicalOwnerUser,omitempty"`
	TechnicalOwnerGroup            *platformGroupNode   `json:"technicalOwnerGroup,omitempty"`
	SecurityOwnerUser              *platformUserNode    `json:"securityOwnerUser,omitempty"`
	SecurityOwnerGroup             *platformGroupNode   `json:"securityOwnerGroup,omitempty"`
	InternalOwnerUser              *platformUserNode    `json:"internalOwnerUser,omitempty"`
	InternalOwnerGroup             *platformGroupNode   `json:"internalOwnerGroup,omitempty"`
	Assets                         platformRelConn      `json:"assets"`
	Entities                       platformRelConn      `json:"entities"`
	SourceAssets                   platformRelConn      `json:"sourceAssets"`
	SourceEntities                 platformRelConn      `json:"sourceEntities"`
	OutOfScopeAssets               platformRelConn      `json:"outOfScopeAssets"`
	OutOfScopeVendors              platformRelConn      `json:"outOfScopeVendors"`
	Controls                       platformRelConn      `json:"controls"`
	Evidence                       platformRelConn      `json:"evidence"`
	ApplicableFrameworks           platformRelConn      `json:"applicableFrameworks"`
	Risks                          platformRelConn      `json:"risks"`
	Tasks                          platformRelConn      `json:"tasks"`
	ArchitectureDiagrams           platformDiagramConn  `json:"architectureDiagrams"`
	DataFlowDiagrams               platformDiagramConn  `json:"dataFlowDiagrams"`
	TrustBoundaryDiagrams          platformDiagramConn  `json:"trustBoundaryDiagrams"`
	CreatedAt                      *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt                      *time.Time           `json:"updatedAt,omitempty"`
}

type platformDetailResponse struct {
	Platform platformDetailNode `json:"platform"`
}

type createPlatformWithDiagramsResponse struct {
	CreatePlatform struct {
		Platform platformDetailNode `json:"platform"`
	} `json:"createPlatform"`
}

type updatePlatformWithDiagramsResponse struct {
	UpdatePlatform struct {
		Platform platformDetailNode `json:"platform"`
	} `json:"updatePlatform"`
}

func (a *api) GetPlatforms(ctx context.Context, first *int64, after *string, where *graphclient.PlatformWhereInput) (*graphclient.GetPlatforms, error) {
	// GetPlatforms puts orderBy before where (unlike most list queries). Do not pass a
	// trailing nil — it becomes a nil variadic interceptor and panics the gqlgenc chain.
	return a.c.GetPlatforms(ctx, first, nil, after, nil, nil, where)
}

func (a *api) GetPlatformByID(ctx context.Context, id string) (*graphclient.GetPlatformByID, error) {
	return a.c.GetPlatformByID(ctx, id)
}

func (a *api) GetPlatformDetail(ctx context.Context, id string) (*PlatformDetail, error) {
	gc, err := a.graphClient()
	if err != nil {
		return nil, err
	}
	vars := map[string]any{"platformId": id}
	var res platformDetailResponse
	if err := gc.Client.Post(ctx, "PlatformDetail", platformDetailQuery, &res, vars); err != nil {
		return nil, RedactError(err)
	}
	if res.Platform.ID == "" {
		return nil, fmt.Errorf("platform not found")
	}
	return mapPlatformDetailNode(res.Platform), nil
}

func (a *api) CreatePlatform(ctx context.Context, input graphclient.CreatePlatformInput, diagrams PlatformDiagramUploads) (*PlatformDetail, error) {
	return a.platformMutation(ctx, diagrams.hasUploads(), func(ctx context.Context) (*PlatformDetail, error) {
		gc, err := a.graphClient()
		if err != nil {
			return nil, err
		}
		vars := map[string]any{
			"input":                 input,
			"architectureDiagrams":  diagrams.Architecture,
			"dataFlowDiagrams":      diagrams.DataFlow,
			"trustBoundaryDiagrams": diagrams.TrustBoundary,
		}
		var res createPlatformWithDiagramsResponse
		if err := gc.Client.Post(ctx, "CreatePlatformWithDiagrams", createPlatformWithDiagramsMutation, &res, vars); err != nil {
			return nil, RedactError(err)
		}
		if res.CreatePlatform.Platform.ID == "" {
			return nil, fmt.Errorf("create platform: empty response")
		}
		return mapPlatformDetailNode(res.CreatePlatform.Platform), nil
	})
}

func (a *api) UpdatePlatform(ctx context.Context, id string, input graphclient.UpdatePlatformInput, diagrams PlatformDiagramUploads) (*PlatformDetail, error) {
	return a.platformMutation(ctx, diagrams.hasUploads(), func(ctx context.Context) (*PlatformDetail, error) {
		gc, err := a.graphClient()
		if err != nil {
			return nil, err
		}
		vars := map[string]any{
			"updatePlatformId":      id,
			"input":                 input,
			"architectureDiagrams":  diagrams.Architecture,
			"dataFlowDiagrams":      diagrams.DataFlow,
			"trustBoundaryDiagrams": diagrams.TrustBoundary,
		}
		var res updatePlatformWithDiagramsResponse
		if err := gc.Client.Post(ctx, "UpdatePlatformWithDiagrams", updatePlatformWithDiagramsMutation, &res, vars); err != nil {
			return nil, RedactError(err)
		}
		if res.UpdatePlatform.Platform.ID == "" {
			return nil, fmt.Errorf("update platform: empty response")
		}
		return mapPlatformDetailNode(res.UpdatePlatform.Platform), nil
	})
}

func (a *api) platformMutation(ctx context.Context, uploading bool, fn func(context.Context) (*PlatformDetail, error)) (*PlatformDetail, error) {
	return a.withUploadTimeout(ctx, uploading, fn)
}

func mapPlatformDetailNode(n platformDetailNode) *PlatformDetail {
	return &PlatformDetail{
		ID:                             n.ID,
		DisplayID:                      n.DisplayID,
		Name:                           n.Name,
		Description:                    Deref(n.Description),
		BusinessPurpose:                Deref(n.BusinessPurpose),
		ScopeStatement:                 Deref(n.ScopeStatement),
		TrustBoundaryDescription:       Deref(n.TrustBoundaryDescription),
		DataFlowSummary:                Deref(n.DataFlowSummary),
		Status:                         n.Status,
		PhysicalLocation:               Deref(n.PhysicalLocation),
		Region:                         Deref(n.Region),
		SourceType:                     n.SourceType,
		SourceIdentifier:               Deref(n.SourceIdentifier),
		CostCenter:                     Deref(n.CostCenter),
		EstimatedMonthlyCost:           n.EstimatedMonthlyCost,
		PurchaseDate:                   n.PurchaseDate,
		ExternalReferenceID:            Deref(n.ExternalReferenceID),
		Tags:                           n.Tags,
		PlatformKindName:               Deref(n.PlatformKindName),
		PlatformDataClassificationName: Deref(n.PlatformDataClassificationName),
		EnvironmentName:                Deref(n.EnvironmentName),
		ScopeName:                      Deref(n.ScopeName),
		AccessModelName:                Deref(n.AccessModelName),
		EncryptionStatusName:           Deref(n.EncryptionStatusName),
		SecurityTierName:               Deref(n.SecurityTierName),
		CriticalityName:                Deref(n.CriticalityName),
		BusinessOwner:                  mapPlatformOwnerRole(n.BusinessOwner, n.BusinessOwnerUser, n.BusinessOwnerGroup),
		TechnicalOwner:                 mapPlatformOwnerRole(n.TechnicalOwner, n.TechnicalOwnerUser, n.TechnicalOwnerGroup),
		SecurityOwner:                  mapPlatformOwnerRole(n.SecurityOwner, n.SecurityOwnerUser, n.SecurityOwnerGroup),
		InternalOwner:                  mapPlatformOwnerRole(n.InternalOwner, n.InternalOwnerUser, n.InternalOwnerGroup),
		Assets:                         mapPlatformRelSummary(n.Assets, relNameOnly),
		Entities:                       mapPlatformRelSummary(n.Entities, relEntityName),
		SourceAssets:                   mapPlatformRelSummary(n.SourceAssets, relNameOnly),
		SourceEntities:                 mapPlatformRelSummary(n.SourceEntities, relEntityName),
		OutOfScopeAssets:               mapPlatformRelSummary(n.OutOfScopeAssets, relNameOnly),
		OutOfScopeVendors:              mapPlatformRelSummary(n.OutOfScopeVendors, relEntityName),
		Controls:                       mapPlatformRelSummary(n.Controls, relControl),
		Evidence:                       mapPlatformRelSummary(n.Evidence, relEvidence),
		ApplicableFrameworks:           mapPlatformRelSummary(n.ApplicableFrameworks, relFramework),
		Risks:                          mapPlatformRelSummary(n.Risks, relRisk),
		Tasks:                          mapPlatformRelSummary(n.Tasks, relTask),
		ArchitectureDiagrams:           mapPlatformDiagrams(n.ArchitectureDiagrams),
		DataFlowDiagrams:               mapPlatformDiagrams(n.DataFlowDiagrams),
		TrustBoundaryDiagrams:          mapPlatformDiagrams(n.TrustBoundaryDiagrams),
		CreatedAt:                      n.CreatedAt,
		UpdatedAt:                      n.UpdatedAt,
	}
}

func mapPlatformOwnerRole(freeText *string, user *platformUserNode, group *platformGroupNode) PlatformOwnerRole {
	if user != nil && user.ID != "" {
		return PlatformOwnerRole{
			Kind:            "user",
			UserID:          user.ID,
			UserEmail:       Deref(user.Email),
			UserDisplayName: Deref(user.DisplayName),
		}
	}
	if group != nil && group.ID != "" {
		return PlatformOwnerRole{
			Kind:             "group",
			GroupID:          group.ID,
			GroupName:        group.Name,
			GroupDisplayName: group.DisplayName,
		}
	}
	if freeText != nil && *freeText != "" {
		return PlatformOwnerRole{Kind: "name", Name: *freeText}
	}
	return PlatformOwnerRole{}
}

type relMapper func(*platformRelNode) PlatformRelNode

func relNameOnly(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	return PlatformRelNode{ID: n.ID, Name: n.Name}
}

func relEntityName(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	name := n.Name
	if name == "" {
		name = Deref(n.DisplayName)
	}
	return PlatformRelNode{ID: n.ID, Name: name}
}

func relControl(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	return PlatformRelNode{
		ID:      n.ID,
		RefCode: n.RefCode,
		Title:   Deref(n.Title),
		Status:  formatAnyStatus(n.Status),
	}
}

func relEvidence(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	return PlatformRelNode{ID: n.ID, Name: n.Name, Status: formatAnyStatus(n.Status)}
}

func relFramework(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	name := n.Name
	if fw := Deref(n.Framework); fw != "" {
		name = fw
	}
	return PlatformRelNode{ID: n.ID, Name: name}
}

func relRisk(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	return PlatformRelNode{ID: n.ID, Name: n.Name, Status: formatAnyStatus(n.Status)}
}

func formatAnyStatus(v any) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprintf("%v", v)
	}
}

func relTask(n *platformRelNode) PlatformRelNode {
	if n == nil {
		return PlatformRelNode{}
	}
	return PlatformRelNode{
		ID:        n.ID,
		DisplayID: n.DisplayID,
		Title:     Deref(n.Title),
		Status:    formatAnyStatus(n.Status),
	}
}

func mapPlatformRelSummary(conn platformRelConn, mapNode relMapper) PlatformRelSummary {
	items := make([]PlatformRelNode, 0, len(conn.Edges))
	for _, e := range conn.Edges {
		if e.Node == nil {
			continue
		}
		items = append(items, mapNode(e.Node))
	}
	return PlatformRelSummary{Count: conn.TotalCount, Items: items}
}

func mapPlatformDiagrams(conn platformDiagramConn) []PlatformDiagramRef {
	out := make([]PlatformDiagramRef, 0, len(conn.Edges))
	for _, e := range conn.Edges {
		if e.Node == nil {
			continue
		}
		out = append(out, PlatformDiagramRef{ID: e.Node.ID, Filename: e.Node.ProvidedFileName})
	}
	return out
}
