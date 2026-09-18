package tools

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type entityLogoInput struct {
	Filename      string `json:"filename" jsonschema:"Original filename including extension."`
	ContentType   string `json:"content_type,omitempty" jsonschema:"MIME type, for example image/png."`
	ContentBase64 string `json:"content_base64" jsonschema:"Base64-encoded logo file contents."`
}

type createEntityInput struct {
	Name                  string           `json:"name" jsonschema:"Entity (vendor) name."`
	DisplayName           string           `json:"display_name,omitempty" jsonschema:"Friendly display name."`
	Description           string           `json:"description,omitempty" jsonschema:"Entity description."`
	EntityTypeID          string           `json:"entity_type_id,omitempty" jsonschema:"Entity type ID."`
	EntityTypeName        string           `json:"entity_type_name,omitempty" jsonschema:"Entity type name (alternative to entity_type_id)."`
	Tags                  []string         `json:"tags,omitempty" jsonschema:"Tags to apply."`
	Domains               []string         `json:"domains,omitempty" jsonschema:"Domains associated with the entity."`
	Tier                  string           `json:"tier,omitempty" jsonschema:"Vendor tier enum value."`
	RiskRating            string           `json:"risk_rating,omitempty" jsonschema:"Vendor risk rating label."`
	RiskScore             *int64           `json:"risk_score,omitempty" jsonschema:"Vendor risk score."`
	ApprovedForUse        *bool            `json:"approved_for_use,omitempty" jsonschema:"Whether the entity is approved for use."`
	NextReviewAt          string           `json:"next_review_at,omitempty" jsonschema:"Next review date (RFC3339 or YYYY-MM-DD)."`
	HasSoc2               *bool            `json:"has_soc2,omitempty" jsonschema:"Whether the entity has an active SOC 2 report."`
	Soc2PeriodEnd         string           `json:"soc2_period_end,omitempty" jsonschema:"SOC 2 period end date (RFC3339 or YYYY-MM-DD)."`
	SsoEnforced           *bool            `json:"sso_enforced,omitempty" jsonschema:"Whether SSO is enforced."`
	MfaSupported          *bool            `json:"mfa_supported,omitempty" jsonschema:"Whether MFA is supported."`
	MfaEnforced           *bool            `json:"mfa_enforced,omitempty" jsonschema:"Whether MFA is enforced."`
	ContractStartDate     string           `json:"contract_start_date,omitempty" jsonschema:"Contract start date (RFC3339 or YYYY-MM-DD)."`
	ContractEndDate       string           `json:"contract_end_date,omitempty" jsonschema:"Contract end date (RFC3339 or YYYY-MM-DD)."`
	ContractRenewalAt     string           `json:"contract_renewal_at,omitempty" jsonschema:"Contract renewal date (RFC3339 or YYYY-MM-DD)."`
	AutoRenews            *bool            `json:"auto_renews,omitempty" jsonschema:"Whether the contract auto-renews."`
	TerminationNoticeDays *int64           `json:"termination_notice_days,omitempty" jsonschema:"Termination notice period in days."`
	AnnualSpend           *float64         `json:"annual_spend,omitempty" jsonschema:"Annual spend amount."`
	SpendCurrency         string           `json:"spend_currency,omitempty" jsonschema:"Spend currency code."`
	BillingModel          string           `json:"billing_model,omitempty" jsonschema:"Billing model label."`
	RenewalRisk           string           `json:"renewal_risk,omitempty" jsonschema:"Renewal risk rating."`
	InternalOwner         string           `json:"internal_owner,omitempty" jsonschema:"Internal owner name when no user or group is linked."`
	InternalOwnerUserID   string           `json:"internal_owner_user_id,omitempty" jsonschema:"Internal owner user ID."`
	InternalOwnerGroupID  string           `json:"internal_owner_group_id,omitempty" jsonschema:"Internal owner group ID."`
	LogoRemoteURL         string           `json:"logo_remote_url,omitempty" jsonschema:"Remote URL for the entity logo."`
	Logo                  *entityLogoInput `json:"logo,omitempty" jsonschema:"Optional logo file uploaded as base64."`
}

type updateEntityInput struct {
	ID                    string           `json:"id" jsonschema:"Entity ID to update."`
	Name                  string           `json:"name,omitempty" jsonschema:"Updated name."`
	DisplayName           string           `json:"display_name,omitempty" jsonschema:"Updated display name."`
	Description           string           `json:"description,omitempty" jsonschema:"Updated description."`
	EntityTypeID          string           `json:"entity_type_id,omitempty" jsonschema:"Updated entity type ID."`
	Tags                  []string         `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	Domains               []string         `json:"domains,omitempty" jsonschema:"Replace domains with this list."`
	Tier                  string           `json:"tier,omitempty" jsonschema:"Updated vendor tier enum value."`
	RiskRating            string           `json:"risk_rating,omitempty" jsonschema:"Updated risk rating label."`
	RiskScore             *int64           `json:"risk_score,omitempty" jsonschema:"Updated risk score."`
	ApprovedForUse        *bool            `json:"approved_for_use,omitempty" jsonschema:"Updated approved-for-use flag."`
	NextReviewAt          string           `json:"next_review_at,omitempty" jsonschema:"Updated next review date (RFC3339 or YYYY-MM-DD)."`
	HasSoc2               *bool            `json:"has_soc2,omitempty" jsonschema:"Updated SOC 2 flag."`
	Soc2PeriodEnd         string           `json:"soc2_period_end,omitempty" jsonschema:"Updated SOC 2 period end date."`
	SsoEnforced           *bool            `json:"sso_enforced,omitempty" jsonschema:"Updated SSO enforced flag."`
	MfaSupported          *bool            `json:"mfa_supported,omitempty" jsonschema:"Updated MFA supported flag."`
	MfaEnforced           *bool            `json:"mfa_enforced,omitempty" jsonschema:"Updated MFA enforced flag."`
	ContractStartDate     string           `json:"contract_start_date,omitempty" jsonschema:"Updated contract start date."`
	ContractEndDate       string           `json:"contract_end_date,omitempty" jsonschema:"Updated contract end date."`
	ContractRenewalAt     string           `json:"contract_renewal_at,omitempty" jsonschema:"Updated contract renewal date."`
	AutoRenews            *bool            `json:"auto_renews,omitempty" jsonschema:"Updated auto-renew flag."`
	TerminationNoticeDays *int64           `json:"termination_notice_days,omitempty" jsonschema:"Updated termination notice days."`
	AnnualSpend           *float64         `json:"annual_spend,omitempty" jsonschema:"Updated annual spend."`
	SpendCurrency         string           `json:"spend_currency,omitempty" jsonschema:"Updated spend currency."`
	BillingModel          string           `json:"billing_model,omitempty" jsonschema:"Updated billing model."`
	RenewalRisk           string           `json:"renewal_risk,omitempty" jsonschema:"Updated renewal risk."`
	InternalOwner         string           `json:"internal_owner,omitempty" jsonschema:"Updated internal owner name."`
	InternalOwnerUserID   string           `json:"internal_owner_user_id,omitempty" jsonschema:"Updated internal owner user ID."`
	InternalOwnerGroupID  string           `json:"internal_owner_group_id,omitempty" jsonschema:"Updated internal owner group ID."`
	AddContactIDs         []string         `json:"add_contact_ids,omitempty" jsonschema:"Contact IDs to associate with this entity (vendor)."`
	AddReviewIDs          []string         `json:"add_review_ids,omitempty" jsonschema:"Review IDs to associate with this entity (vendor), including Risk Reviews."`
	RemoveReviewIDs       []string         `json:"remove_review_ids,omitempty" jsonschema:"Review IDs to unlink from this entity (vendor)."`
	LogoRemoteURL         string           `json:"logo_remote_url,omitempty" jsonschema:"Updated remote logo URL."`
	Logo                  *entityLogoInput `json:"logo,omitempty" jsonschema:"Optional logo file uploaded as base64."`
}

func registerWriteEntities(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_entity_create",
		Title:       "Create an Openlane entity",
		Description: "Create an entity (vendor) in Openlane. Optional logo upload via base64 or logo_remote_url. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createEntity)

	addTool(server, &mcp.Tool{
		Name:        "openlane_entity_update",
		Title:       "Update an Openlane entity",
		Description: "Update an entity (vendor) by ID. Optional logo upload via base64 or logo_remote_url. Use add_contact_ids to associate contacts. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateEntity)
}

func (h *handlers) createEntity(ctx context.Context, _ *mcp.CallToolRequest, in createEntityInput) (*mcp.CallToolResult, entityItem, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, entityItem{}, errNameRequired
	}
	input := graphclient.CreateEntityInput{
		Name: &in.Name,
		Tags: in.Tags,
	}
	if err := applyCreateEntityFields(&input, in); err != nil {
		return nil, entityItem{}, err
	}

	logoFile, err := h.decodeEntityLogo(in.Logo)
	if err != nil {
		return nil, entityItem{}, err
	}

	var entityTypeName *string
	if s := strings.TrimSpace(in.EntityTypeName); s != "" {
		entityTypeName = &s
	}

	resp, err := h.api.CreateEntity(ctx, input, entityTypeName, logoFile)
	if err != nil {
		return nil, entityItem{}, openlane.APIError(err)
	}
	return nil, mapEntityDetail(*resp), nil
}

func (h *handlers) updateEntity(ctx context.Context, _ *mcp.CallToolRequest, in updateEntityInput) (*mcp.CallToolResult, entityItem, error) {
	if in.ID == "" {
		return nil, entityItem{}, errIDRequired
	}
	input := graphclient.UpdateEntityInput{}
	if err := applyUpdateEntityFields(&input, in); err != nil {
		return nil, entityItem{}, err
	}

	logoFile, err := h.decodeEntityLogo(in.Logo)
	if err != nil {
		return nil, entityItem{}, err
	}
	if isEmptyUpdateEntity(input) && logoFile == nil {
		return nil, entityItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateEntity(ctx, in.ID, input, logoFile)
	if err != nil {
		return nil, entityItem{}, openlane.APIError(err)
	}
	return nil, mapEntityDetail(*resp), nil
}

func applyCreateEntityFields(input *graphclient.CreateEntityInput, in createEntityInput) error {
	if s := strings.TrimSpace(in.DisplayName); s != "" {
		input.DisplayName = &s
	}
	if s := strings.TrimSpace(in.Description); s != "" {
		input.Description = &s
	}
	if s := strings.TrimSpace(in.EntityTypeID); s != "" {
		input.EntityTypeID = &s
	}
	if len(in.Domains) > 0 {
		input.Domains = in.Domains
	}
	if s := strings.TrimSpace(in.Tier); s != "" {
		input.Tier = vendorTier(s)
	}
	if s := strings.TrimSpace(in.RiskRating); s != "" {
		input.RiskRating = &s
	}
	if in.RiskScore != nil {
		input.RiskScore = in.RiskScore
	}
	if in.ApprovedForUse != nil {
		input.ApprovedForUse = in.ApprovedForUse
	}
	if err := setOptionalDate(&input.NextReviewAt, in.NextReviewAt); err != nil {
		return err
	}
	if in.HasSoc2 != nil {
		input.HasSoc2 = in.HasSoc2
	}
	if err := setOptionalDate(&input.Soc2PeriodEnd, in.Soc2PeriodEnd); err != nil {
		return err
	}
	if in.SsoEnforced != nil {
		input.SsoEnforced = in.SsoEnforced
	}
	if in.MfaSupported != nil {
		input.MfaSupported = in.MfaSupported
	}
	if in.MfaEnforced != nil {
		input.MfaEnforced = in.MfaEnforced
	}
	if err := setOptionalDate(&input.ContractStartDate, in.ContractStartDate); err != nil {
		return err
	}
	if err := setOptionalDate(&input.ContractEndDate, in.ContractEndDate); err != nil {
		return err
	}
	if err := setOptionalDate(&input.ContractRenewalAt, in.ContractRenewalAt); err != nil {
		return err
	}
	if in.AutoRenews != nil {
		input.AutoRenews = in.AutoRenews
	}
	if in.TerminationNoticeDays != nil {
		input.TerminationNoticeDays = in.TerminationNoticeDays
	}
	if in.AnnualSpend != nil {
		input.AnnualSpend = in.AnnualSpend
	}
	if s := strings.TrimSpace(in.SpendCurrency); s != "" {
		input.SpendCurrency = &s
	}
	if s := strings.TrimSpace(in.BillingModel); s != "" {
		input.BillingModel = &s
	}
	if s := strings.TrimSpace(in.RenewalRisk); s != "" {
		input.RenewalRisk = &s
	}
	if s := strings.TrimSpace(in.InternalOwner); s != "" {
		input.InternalOwner = &s
	}
	if s := strings.TrimSpace(in.InternalOwnerUserID); s != "" {
		input.InternalOwnerUserID = &s
	}
	if s := strings.TrimSpace(in.InternalOwnerGroupID); s != "" {
		input.InternalOwnerGroupID = &s
	}
	if s := strings.TrimSpace(in.LogoRemoteURL); s != "" {
		input.LogoRemoteURL = &s
	}
	return nil
}

func applyUpdateEntityFields(input *graphclient.UpdateEntityInput, in updateEntityInput) error {
	if s := strings.TrimSpace(in.Name); s != "" {
		input.Name = &s
	}
	if s := strings.TrimSpace(in.DisplayName); s != "" {
		input.DisplayName = &s
	}
	if s := strings.TrimSpace(in.Description); s != "" {
		input.Description = &s
	}
	if s := strings.TrimSpace(in.EntityTypeID); s != "" {
		input.EntityTypeID = &s
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if len(in.Domains) > 0 {
		input.Domains = in.Domains
	}
	if s := strings.TrimSpace(in.Tier); s != "" {
		input.Tier = vendorTier(s)
	}
	if s := strings.TrimSpace(in.RiskRating); s != "" {
		input.RiskRating = &s
	}
	if in.RiskScore != nil {
		input.RiskScore = in.RiskScore
	}
	if in.ApprovedForUse != nil {
		input.ApprovedForUse = in.ApprovedForUse
	}
	if err := setOptionalDate(&input.NextReviewAt, in.NextReviewAt); err != nil {
		return err
	}
	if in.HasSoc2 != nil {
		input.HasSoc2 = in.HasSoc2
	}
	if err := setOptionalDate(&input.Soc2PeriodEnd, in.Soc2PeriodEnd); err != nil {
		return err
	}
	if in.SsoEnforced != nil {
		input.SsoEnforced = in.SsoEnforced
	}
	if in.MfaSupported != nil {
		input.MfaSupported = in.MfaSupported
	}
	if in.MfaEnforced != nil {
		input.MfaEnforced = in.MfaEnforced
	}
	if err := setOptionalDate(&input.ContractStartDate, in.ContractStartDate); err != nil {
		return err
	}
	if err := setOptionalDate(&input.ContractEndDate, in.ContractEndDate); err != nil {
		return err
	}
	if err := setOptionalDate(&input.ContractRenewalAt, in.ContractRenewalAt); err != nil {
		return err
	}
	if in.AutoRenews != nil {
		input.AutoRenews = in.AutoRenews
	}
	if in.TerminationNoticeDays != nil {
		input.TerminationNoticeDays = in.TerminationNoticeDays
	}
	if in.AnnualSpend != nil {
		input.AnnualSpend = in.AnnualSpend
	}
	if s := strings.TrimSpace(in.SpendCurrency); s != "" {
		input.SpendCurrency = &s
	}
	if s := strings.TrimSpace(in.BillingModel); s != "" {
		input.BillingModel = &s
	}
	if s := strings.TrimSpace(in.RenewalRisk); s != "" {
		input.RenewalRisk = &s
	}
	if s := strings.TrimSpace(in.InternalOwner); s != "" {
		input.InternalOwner = &s
	}
	if s := strings.TrimSpace(in.InternalOwnerUserID); s != "" {
		input.InternalOwnerUserID = &s
	}
	if s := strings.TrimSpace(in.InternalOwnerGroupID); s != "" {
		input.InternalOwnerGroupID = &s
	}
	if s := strings.TrimSpace(in.LogoRemoteURL); s != "" {
		input.LogoRemoteURL = &s
	}
	if len(in.AddContactIDs) > 0 {
		input.AddContactIDs = in.AddContactIDs
	}
	if len(in.AddReviewIDs) > 0 {
		input.AddReviewIDs = in.AddReviewIDs
	}
	if len(in.RemoveReviewIDs) > 0 {
		input.RemoveReviewIDs = in.RemoveReviewIDs
	}
	return nil
}

func (h *handlers) decodeEntityLogo(logo *entityLogoInput) (*graphql.Upload, error) {
	if logo == nil {
		return nil, nil
	}
	maxBytes := h.maxUploadBytes
	if maxBytes <= 0 {
		maxBytes = openlane.DefaultMaxUploadBytes
	}
	uploads, err := openlane.DecodeEvidenceUploads([]openlane.EvidenceFile{{
		Filename:      logo.Filename,
		ContentType:   logo.ContentType,
		ContentBase64: logo.ContentBase64,
	}}, maxBytes)
	if err != nil {
		return nil, err
	}
	if len(uploads) == 0 {
		return nil, nil
	}
	return uploads[0], nil
}

func vendorTier(value string) *enums.VendorTier {
	return enumPtr[enums.VendorTier](value)
}

func setOptionalDate(dst **models.DateTime, value string) error {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	dt, err := parseDateTime(value)
	if err != nil {
		return err
	}
	*dst = dt
	return nil
}

func isEmptyUpdateEntity(in graphclient.UpdateEntityInput) bool {
	return in.Name == nil &&
		in.DisplayName == nil &&
		in.Description == nil &&
		in.EntityTypeID == nil &&
		len(in.Tags) == 0 &&
		len(in.Domains) == 0 &&
		in.Tier == nil &&
		in.RiskRating == nil &&
		in.RiskScore == nil &&
		in.ApprovedForUse == nil &&
		in.NextReviewAt == nil &&
		in.HasSoc2 == nil &&
		in.Soc2PeriodEnd == nil &&
		in.SsoEnforced == nil &&
		in.MfaSupported == nil &&
		in.MfaEnforced == nil &&
		in.ContractStartDate == nil &&
		in.ContractEndDate == nil &&
		in.ContractRenewalAt == nil &&
		in.AutoRenews == nil &&
		in.TerminationNoticeDays == nil &&
		in.AnnualSpend == nil &&
		in.SpendCurrency == nil &&
		in.BillingModel == nil &&
		in.RenewalRisk == nil &&
		in.InternalOwner == nil &&
		in.InternalOwnerUserID == nil &&
		in.InternalOwnerGroupID == nil &&
		in.LogoRemoteURL == nil &&
		len(in.AddContactIDs) == 0 &&
		len(in.AddReviewIDs) == 0 &&
		len(in.RemoveReviewIDs) == 0
}
