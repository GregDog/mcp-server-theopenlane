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

// EntityDetail is an entity (vendor) with logo fields omitted from go-client v0.14.0 Get/mutation selection sets.
type EntityDetail struct {
	ID                                    string
	Name                                  *string
	DisplayName                           *string
	Description                           *string
	EntityTypeID                          *string
	EntitySourceTypeName                  *string
	EntityRelationshipStateName           *string
	EntitySecurityQuestionnaireStatusName *string
	EnvironmentName                       *string
	Tier                                  *enums.VendorTier
	RiskRating                            *string
	RiskScore                             *int64
	ApprovedForUse                        *bool
	HasSoc2                               *bool
	Soc2PeriodEnd                         *models.DateTime
	SsoEnforced                           *bool
	MfaSupported                          *bool
	MfaEnforced                           *bool
	LastReviewedAt                        *models.DateTime
	NextReviewAt                          *models.DateTime
	ContractStartDate                     *models.DateTime
	ContractEndDate                       *models.DateTime
	ContractRenewalAt                     *models.DateTime
	AutoRenews                            *bool
	TerminationNoticeDays                 *int64
	AnnualSpend                           *float64
	SpendCurrency                         *string
	BillingModel                          *string
	RenewalRisk                           *string
	InternalOwner                         *string
	InternalOwnerUserID                   *string
	InternalOwnerGroupID                  *string
	Tags                                  []string
	LogoRemoteURL                         *string
	LogoFileID                            *string
	CreatedAt                             *time.Time
	UpdatedAt                             *time.Time
}

const entityDetailFields = `
      id
      name
      displayName
      description
      entityTypeID
      entitySourceTypeName
      entityRelationshipStateName
      entitySecurityQuestionnaireStatusName
      environmentName
      tier
      riskRating
      riskScore
      approvedForUse
      hasSoc2
      soc2PeriodEnd
      ssoEnforced
      mfaSupported
      mfaEnforced
      lastReviewedAt
      nextReviewAt
      contractStartDate
      contractEndDate
      contractRenewalAt
      autoRenews
      terminationNoticeDays
      annualSpend
      spendCurrency
      billingModel
      renewalRisk
      internalOwner
      internalOwnerUserID
      internalOwnerGroupID
      tags
      logoRemoteURL
      logoFileID
      createdAt
      updatedAt
`

const createEntityWithLogoMutation = `mutation CreateEntityWithLogo($input: CreateEntityInput!, $entityTypeName: String, $logoFile: Upload) {
  createEntity(input: $input, entityTypeName: $entityTypeName, logoFile: $logoFile) {
    entity {
` + entityDetailFields + `
    }
  }
}`

const updateEntityWithLogoMutation = `mutation UpdateEntityWithLogo($updateEntityId: ID!, $input: UpdateEntityInput!, $logoFile: Upload) {
  updateEntity(id: $updateEntityId, input: $input, logoFile: $logoFile) {
    entity {
` + entityDetailFields + `
    }
  }
}`

const entityDetailQuery = `query EntityDetail($entityId: ID!) {
  entity(id: $entityId) {
` + entityDetailFields + `
  }
}`

type entityDetailNode struct {
	ID                                    string            `json:"id"`
	Name                                  *string           `json:"name,omitempty"`
	DisplayName                           *string           `json:"displayName,omitempty"`
	Description                           *string           `json:"description,omitempty"`
	EntityTypeID                          *string           `json:"entityTypeID,omitempty"`
	EntitySourceTypeName                  *string           `json:"entitySourceTypeName,omitempty"`
	EntityRelationshipStateName           *string           `json:"entityRelationshipStateName,omitempty"`
	EntitySecurityQuestionnaireStatusName *string           `json:"entitySecurityQuestionnaireStatusName,omitempty"`
	EnvironmentName                       *string           `json:"environmentName,omitempty"`
	Tier                                  *enums.VendorTier `json:"tier,omitempty"`
	RiskRating                            *string           `json:"riskRating,omitempty"`
	RiskScore                             *int64            `json:"riskScore,omitempty"`
	ApprovedForUse                        *bool             `json:"approvedForUse,omitempty"`
	HasSoc2                               *bool             `json:"hasSoc2,omitempty"`
	Soc2PeriodEnd                         *models.DateTime  `json:"soc2PeriodEnd,omitempty"`
	SsoEnforced                           *bool             `json:"ssoEnforced,omitempty"`
	MfaSupported                          *bool             `json:"mfaSupported,omitempty"`
	MfaEnforced                           *bool             `json:"mfaEnforced,omitempty"`
	LastReviewedAt                        *models.DateTime  `json:"lastReviewedAt,omitempty"`
	NextReviewAt                          *models.DateTime  `json:"nextReviewAt,omitempty"`
	ContractStartDate                     *models.DateTime  `json:"contractStartDate,omitempty"`
	ContractEndDate                       *models.DateTime  `json:"contractEndDate,omitempty"`
	ContractRenewalAt                     *models.DateTime  `json:"contractRenewalAt,omitempty"`
	AutoRenews                            *bool             `json:"autoRenews,omitempty"`
	TerminationNoticeDays                 *int64            `json:"terminationNoticeDays,omitempty"`
	AnnualSpend                           *float64          `json:"annualSpend,omitempty"`
	SpendCurrency                         *string           `json:"spendCurrency,omitempty"`
	BillingModel                          *string           `json:"billingModel,omitempty"`
	RenewalRisk                           *string           `json:"renewalRisk,omitempty"`
	InternalOwner                         *string           `json:"internalOwner,omitempty"`
	InternalOwnerUserID                   *string           `json:"internalOwnerUserID,omitempty"`
	InternalOwnerGroupID                  *string           `json:"internalOwnerGroupID,omitempty"`
	Tags                                  []string          `json:"tags,omitempty"`
	LogoRemoteURL                         *string           `json:"logoRemoteURL,omitempty"`
	LogoFileID                            *string           `json:"logoFileID,omitempty"`
	CreatedAt                             *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt                             *time.Time        `json:"updatedAt,omitempty"`
}

func mapEntityDetailNode(n entityDetailNode) *EntityDetail {
	return &EntityDetail{
		ID:                                    n.ID,
		Name:                                  n.Name,
		DisplayName:                           n.DisplayName,
		Description:                           n.Description,
		EntityTypeID:                          n.EntityTypeID,
		EntitySourceTypeName:                  n.EntitySourceTypeName,
		EntityRelationshipStateName:           n.EntityRelationshipStateName,
		EntitySecurityQuestionnaireStatusName: n.EntitySecurityQuestionnaireStatusName,
		EnvironmentName:                       n.EnvironmentName,
		Tier:                                  n.Tier,
		RiskRating:                            n.RiskRating,
		RiskScore:                             n.RiskScore,
		ApprovedForUse:                        n.ApprovedForUse,
		HasSoc2:                               n.HasSoc2,
		Soc2PeriodEnd:                         n.Soc2PeriodEnd,
		SsoEnforced:                           n.SsoEnforced,
		MfaSupported:                          n.MfaSupported,
		MfaEnforced:                           n.MfaEnforced,
		LastReviewedAt:                        n.LastReviewedAt,
		NextReviewAt:                          n.NextReviewAt,
		ContractStartDate:                     n.ContractStartDate,
		ContractEndDate:                       n.ContractEndDate,
		ContractRenewalAt:                     n.ContractRenewalAt,
		AutoRenews:                            n.AutoRenews,
		TerminationNoticeDays:                 n.TerminationNoticeDays,
		AnnualSpend:                           n.AnnualSpend,
		SpendCurrency:                         n.SpendCurrency,
		BillingModel:                          n.BillingModel,
		RenewalRisk:                           n.RenewalRisk,
		InternalOwner:                         n.InternalOwner,
		InternalOwnerUserID:                   n.InternalOwnerUserID,
		InternalOwnerGroupID:                  n.InternalOwnerGroupID,
		Tags:                                  n.Tags,
		LogoRemoteURL:                         n.LogoRemoteURL,
		LogoFileID:                            n.LogoFileID,
		CreatedAt:                             n.CreatedAt,
		UpdatedAt:                             n.UpdatedAt,
	}
}

type createEntityWithLogoResponse struct {
	CreateEntity struct {
		Entity entityDetailNode `json:"entity"`
	} `json:"createEntity"`
}

type updateEntityWithLogoResponse struct {
	UpdateEntity struct {
		Entity entityDetailNode `json:"entity"`
	} `json:"updateEntity"`
}

type entityDetailResponse struct {
	Entity entityDetailNode `json:"entity"`
}

func (a *api) CreateEntity(ctx context.Context, input graphclient.CreateEntityInput, entityTypeName *string, logoFile *graphql.Upload) (*EntityDetail, error) {
	return a.withUploadTimeout(ctx, logoFile != nil, func(ctx context.Context) (*EntityDetail, error) {
		gc, err := a.graphClient()
		if err != nil {
			return nil, err
		}
		vars := map[string]any{
			"input":          input,
			"entityTypeName": entityTypeName,
			"logoFile":       logoFile,
		}
		var res createEntityWithLogoResponse
		if err := gc.Client.Post(ctx, "CreateEntityWithLogo", createEntityWithLogoMutation, &res, vars); err != nil {
			return nil, RedactError(err)
		}
		if res.CreateEntity.Entity.ID == "" {
			return nil, fmt.Errorf("create entity: empty response")
		}
		return mapEntityDetailNode(res.CreateEntity.Entity), nil
	})
}

func (a *api) UpdateEntity(ctx context.Context, id string, input graphclient.UpdateEntityInput, logoFile *graphql.Upload) (*EntityDetail, error) {
	return a.withUploadTimeout(ctx, logoFile != nil, func(ctx context.Context) (*EntityDetail, error) {
		gc, err := a.graphClient()
		if err != nil {
			return nil, err
		}
		vars := map[string]any{
			"updateEntityId": id,
			"input":          input,
			"logoFile":       logoFile,
		}
		var res updateEntityWithLogoResponse
		if err := gc.Client.Post(ctx, "UpdateEntityWithLogo", updateEntityWithLogoMutation, &res, vars); err != nil {
			return nil, RedactError(err)
		}
		if res.UpdateEntity.Entity.ID == "" {
			return nil, fmt.Errorf("update entity: empty response")
		}
		return mapEntityDetailNode(res.UpdateEntity.Entity), nil
	})
}

func (a *api) GetEntityDetail(ctx context.Context, id string) (*EntityDetail, error) {
	gc, err := a.graphClient()
	if err != nil {
		return nil, err
	}
	vars := map[string]any{"entityId": id}
	var res entityDetailResponse
	if err := gc.Client.Post(ctx, "EntityDetail", entityDetailQuery, &res, vars); err != nil {
		return nil, RedactError(err)
	}
	if res.Entity.ID == "" {
		return nil, fmt.Errorf("entity not found")
	}
	return mapEntityDetailNode(res.Entity), nil
}
