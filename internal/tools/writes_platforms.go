package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type platformDiagramInput struct {
	Filename      string `json:"filename" jsonschema:"Original filename including extension."`
	ContentType   string `json:"content_type,omitempty" jsonschema:"MIME type, for example image/png."`
	ContentBase64 string `json:"content_base64" jsonschema:"Base64-encoded file contents."`
}

type createPlatformInput struct {
	Name                           string                 `json:"name" jsonschema:"Platform name."`
	Description                    string                 `json:"description,omitempty" jsonschema:"Platform boundary description."`
	BusinessPurpose                string                 `json:"business_purpose,omitempty" jsonschema:"Business purpose of the platform."`
	ScopeStatement                 string                 `json:"scope_statement,omitempty" jsonschema:"Scope statement for audit narrative."`
	TrustBoundaryDescription       string                 `json:"trust_boundary_description,omitempty" jsonschema:"Trust boundary description."`
	DataFlowSummary                string                 `json:"data_flow_summary,omitempty" jsonschema:"Data flow summary."`
	Status                         string                 `json:"status,omitempty" jsonschema:"Lifecycle status: ACTIVE, INACTIVE, or RETIRED."`
	PhysicalLocation               string                 `json:"physical_location,omitempty" jsonschema:"Physical location when applicable."`
	Region                         string                 `json:"region,omitempty" jsonschema:"Hosting or operating region."`
	SourceType                     string                 `json:"source_type,omitempty" jsonschema:"Source type, for example MANUAL."`
	SourceIdentifier               string                 `json:"source_identifier,omitempty" jsonschema:"Identifier from the upstream source system."`
	CostCenter                     string                 `json:"cost_center,omitempty" jsonschema:"Cost center code."`
	EstimatedMonthlyCost           *float64               `json:"estimated_monthly_cost,omitempty" jsonschema:"Estimated monthly cost."`
	PurchaseDate                   string                 `json:"purchase_date,omitempty" jsonschema:"Purchase date (RFC3339 or YYYY-MM-DD)."`
	ExternalReferenceID            string                 `json:"external_reference_id,omitempty" jsonschema:"External inventory reference id."`
	Tags                           []string               `json:"tags,omitempty" jsonschema:"Tags to apply."`
	PlatformKindName               string                 `json:"platform_kind_name,omitempty" jsonschema:"Platform kind custom enum name."`
	PlatformDataClassificationName string                 `json:"platform_data_classification_name,omitempty" jsonschema:"Data classification custom enum name."`
	EnvironmentName                string                 `json:"environment_name,omitempty" jsonschema:"Environment custom enum name."`
	ScopeName                      string                 `json:"scope_name,omitempty" jsonschema:"Scope custom enum name."`
	AccessModelName                string                 `json:"access_model_name,omitempty" jsonschema:"Access model custom enum name."`
	EncryptionStatusName           string                 `json:"encryption_status_name,omitempty" jsonschema:"Encryption status custom enum name."`
	SecurityTierName               string                 `json:"security_tier_name,omitempty" jsonschema:"Security tier custom enum name."`
	CriticalityName                string                 `json:"criticality_name,omitempty" jsonschema:"Criticality custom enum name."`
	BusinessOwner                  string                 `json:"business_owner,omitempty" jsonschema:"Business owner: user id, email, name, or group id/name."`
	BusinessOwnerName              string                 `json:"business_owner_name,omitempty" jsonschema:"Business owner free-text name when no user or group is linked."`
	TechnicalOwner                 string                 `json:"technical_owner,omitempty" jsonschema:"Technical owner: user id, email, name, or group id/name."`
	TechnicalOwnerName             string                 `json:"technical_owner_name,omitempty" jsonschema:"Technical owner free-text name when no user or group is linked."`
	SecurityOwner                  string                 `json:"security_owner,omitempty" jsonschema:"Security owner: user id, email, name, or group id/name."`
	SecurityOwnerName              string                 `json:"security_owner_name,omitempty" jsonschema:"Security owner free-text name when no user or group is linked."`
	InternalOwner                  string                 `json:"internal_owner,omitempty" jsonschema:"Internal owner: user id, email, name, or group id/name."`
	InternalOwnerName              string                 `json:"internal_owner_name,omitempty" jsonschema:"Internal owner free-text name when no user or group is linked."`
	AssetIDs                       []string               `json:"asset_ids,omitempty" jsonschema:"Asset IDs in scope for this platform."`
	EntityIDs                      []string               `json:"entity_ids,omitempty" jsonschema:"Vendor (entity) IDs in scope for this platform."`
	SourceAssetIDs                 []string               `json:"source_asset_ids,omitempty" jsonschema:"Source asset IDs linked to this platform."`
	SourceEntityIDs                []string               `json:"source_entity_ids,omitempty" jsonschema:"Source vendor IDs linked to this platform."`
	OutOfScopeAssetIDs             []string               `json:"out_of_scope_asset_ids,omitempty" jsonschema:"Asset IDs explicitly out of scope."`
	OutOfScopeVendorIDs            []string               `json:"out_of_scope_vendor_ids,omitempty" jsonschema:"Vendor IDs explicitly out of scope."`
	ControlIDs                     []string               `json:"control_ids,omitempty" jsonschema:"Control IDs linked to this platform."`
	EvidenceIDs                    []string               `json:"evidence_ids,omitempty" jsonschema:"Evidence IDs linked to this platform."`
	ApplicableFrameworkIDs         []string               `json:"applicable_framework_ids,omitempty" jsonschema:"Framework (standard) IDs applicable to this platform."`
	RiskIDs                        []string               `json:"risk_ids,omitempty" jsonschema:"Risk IDs linked to this platform."`
	TaskIDs                        []string               `json:"task_ids,omitempty" jsonschema:"Task IDs linked to this platform."`
	ArchitectureDiagrams           []platformDiagramInput `json:"architecture_diagrams,omitempty" jsonschema:"Architecture diagram files uploaded as base64."`
	DataFlowDiagrams               []platformDiagramInput `json:"data_flow_diagrams,omitempty" jsonschema:"Data flow diagram files uploaded as base64."`
	TrustBoundaryDiagrams          []platformDiagramInput `json:"trust_boundary_diagrams,omitempty" jsonschema:"Trust boundary diagram files uploaded as base64."`
}

type updatePlatformInput struct {
	ID                             string                 `json:"id" jsonschema:"Platform ID to update."`
	Name                           string                 `json:"name,omitempty" jsonschema:"Updated platform name."`
	Description                    string                 `json:"description,omitempty" jsonschema:"Updated description."`
	BusinessPurpose                string                 `json:"business_purpose,omitempty" jsonschema:"Updated business purpose."`
	ScopeStatement                 string                 `json:"scope_statement,omitempty" jsonschema:"Updated scope statement."`
	TrustBoundaryDescription       string                 `json:"trust_boundary_description,omitempty" jsonschema:"Updated trust boundary description."`
	DataFlowSummary                string                 `json:"data_flow_summary,omitempty" jsonschema:"Updated data flow summary."`
	Status                         string                 `json:"status,omitempty" jsonschema:"Updated status: ACTIVE, INACTIVE, or RETIRED."`
	PhysicalLocation               string                 `json:"physical_location,omitempty" jsonschema:"Updated physical location."`
	Region                         string                 `json:"region,omitempty" jsonschema:"Updated region."`
	SourceType                     string                 `json:"source_type,omitempty" jsonschema:"Updated source type."`
	SourceIdentifier               string                 `json:"source_identifier,omitempty" jsonschema:"Updated source identifier."`
	CostCenter                     string                 `json:"cost_center,omitempty" jsonschema:"Updated cost center."`
	EstimatedMonthlyCost           *float64               `json:"estimated_monthly_cost,omitempty" jsonschema:"Updated estimated monthly cost."`
	PurchaseDate                   string                 `json:"purchase_date,omitempty" jsonschema:"Updated purchase date (RFC3339 or YYYY-MM-DD)."`
	ExternalReferenceID            string                 `json:"external_reference_id,omitempty" jsonschema:"Updated external reference id."`
	Tags                           []string               `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
	PlatformKindName               string                 `json:"platform_kind_name,omitempty" jsonschema:"Updated platform kind name."`
	PlatformDataClassificationName string                 `json:"platform_data_classification_name,omitempty" jsonschema:"Updated data classification name."`
	EnvironmentName                string                 `json:"environment_name,omitempty" jsonschema:"Updated environment name."`
	ScopeName                      string                 `json:"scope_name,omitempty" jsonschema:"Updated scope name."`
	AccessModelName                string                 `json:"access_model_name,omitempty" jsonschema:"Updated access model name."`
	EncryptionStatusName           string                 `json:"encryption_status_name,omitempty" jsonschema:"Updated encryption status name."`
	SecurityTierName               string                 `json:"security_tier_name,omitempty" jsonschema:"Updated security tier name."`
	CriticalityName                string                 `json:"criticality_name,omitempty" jsonschema:"Updated criticality name."`
	BusinessOwner                  string                 `json:"business_owner,omitempty" jsonschema:"Updated business owner: user id, email, name, or group id/name."`
	BusinessOwnerName              string                 `json:"business_owner_name,omitempty" jsonschema:"Updated business owner free-text name."`
	TechnicalOwner                 string                 `json:"technical_owner,omitempty" jsonschema:"Updated technical owner: user id, email, name, or group id/name."`
	TechnicalOwnerName             string                 `json:"technical_owner_name,omitempty" jsonschema:"Updated technical owner free-text name."`
	SecurityOwner                  string                 `json:"security_owner,omitempty" jsonschema:"Updated security owner: user id, email, name, or group id/name."`
	SecurityOwnerName              string                 `json:"security_owner_name,omitempty" jsonschema:"Updated security owner free-text name."`
	InternalOwner                  string                 `json:"internal_owner,omitempty" jsonschema:"Updated internal owner: user id, email, name, or group id/name."`
	InternalOwnerName              string                 `json:"internal_owner_name,omitempty" jsonschema:"Updated internal owner free-text name."`
	AddAssetIDs                    []string               `json:"add_asset_ids,omitempty" jsonschema:"Asset IDs to link."`
	RemoveAssetIDs                 []string               `json:"remove_asset_ids,omitempty" jsonschema:"Asset IDs to unlink."`
	AddEntityIDs                   []string               `json:"add_entity_ids,omitempty" jsonschema:"Vendor IDs to link."`
	RemoveEntityIDs                []string               `json:"remove_entity_ids,omitempty" jsonschema:"Vendor IDs to unlink."`
	AddSourceAssetIDs              []string               `json:"add_source_asset_ids,omitempty" jsonschema:"Source asset IDs to link."`
	RemoveSourceAssetIDs           []string               `json:"remove_source_asset_ids,omitempty" jsonschema:"Source asset IDs to unlink."`
	AddSourceEntityIDs             []string               `json:"add_source_entity_ids,omitempty" jsonschema:"Source vendor IDs to link."`
	RemoveSourceEntityIDs          []string               `json:"remove_source_entity_ids,omitempty" jsonschema:"Source vendor IDs to unlink."`
	AddOutOfScopeAssetIDs          []string               `json:"add_out_of_scope_asset_ids,omitempty" jsonschema:"Out-of-scope asset IDs to link."`
	RemoveOutOfScopeAssetIDs       []string               `json:"remove_out_of_scope_asset_ids,omitempty" jsonschema:"Out-of-scope asset IDs to unlink."`
	AddOutOfScopeVendorIDs         []string               `json:"add_out_of_scope_vendor_ids,omitempty" jsonschema:"Out-of-scope vendor IDs to link."`
	RemoveOutOfScopeVendorIDs      []string               `json:"remove_out_of_scope_vendor_ids,omitempty" jsonschema:"Out-of-scope vendor IDs to unlink."`
	AddControlIDs                  []string               `json:"add_control_ids,omitempty" jsonschema:"Control IDs to link."`
	RemoveControlIDs               []string               `json:"remove_control_ids,omitempty" jsonschema:"Control IDs to unlink."`
	AddEvidenceIDs                 []string               `json:"add_evidence_ids,omitempty" jsonschema:"Evidence IDs to link."`
	RemoveEvidenceIDs              []string               `json:"remove_evidence_ids,omitempty" jsonschema:"Evidence IDs to unlink."`
	AddApplicableFrameworkIDs      []string               `json:"add_applicable_framework_ids,omitempty" jsonschema:"Framework IDs to link."`
	RemoveApplicableFrameworkIDs   []string               `json:"remove_applicable_framework_ids,omitempty" jsonschema:"Framework IDs to unlink."`
	AddRiskIDs                     []string               `json:"add_risk_ids,omitempty" jsonschema:"Risk IDs to link."`
	RemoveRiskIDs                  []string               `json:"remove_risk_ids,omitempty" jsonschema:"Risk IDs to unlink."`
	AddTaskIDs                     []string               `json:"add_task_ids,omitempty" jsonschema:"Task IDs to link."`
	RemoveTaskIDs                  []string               `json:"remove_task_ids,omitempty" jsonschema:"Task IDs to unlink."`
	ArchitectureDiagrams           []platformDiagramInput `json:"architecture_diagrams,omitempty" jsonschema:"Architecture diagram files to upload as base64."`
	DataFlowDiagrams               []platformDiagramInput `json:"data_flow_diagrams,omitempty" jsonschema:"Data flow diagram files to upload as base64."`
	TrustBoundaryDiagrams          []platformDiagramInput `json:"trust_boundary_diagrams,omitempty" jsonschema:"Trust boundary diagram files to upload as base64."`
}

func registerWritePlatforms(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_platform_create",
		Title:       "Create an Openlane platform",
		Description: "Create a platform (system boundary) in Openlane. Owner fields accept user id, email, name, or group id/name. Optional diagram files upload as base64. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createPlatform)

	addTool(server, &mcp.Tool{
		Name:        "openlane_platform_update",
		Title:       "Update an Openlane platform",
		Description: "Update a platform by ID. Use add/remove id lists to manage scope links. Owner fields accept user id, email, name, or group id/name. Optional diagram files upload as base64. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updatePlatform)
}

func (h *handlers) createPlatform(ctx context.Context, _ *mcp.CallToolRequest, in createPlatformInput) (*mcp.CallToolResult, platformItem, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, platformItem{}, errNameRequired
	}
	input := graphclient.CreatePlatformInput{Name: strings.TrimSpace(in.Name)}
	if err := applyCreatePlatformFields(ctx, h, &input, in); err != nil {
		return nil, platformItem{}, err
	}
	diagrams, err := decodePlatformDiagrams(in.ArchitectureDiagrams, in.DataFlowDiagrams, in.TrustBoundaryDiagrams, h.maxUploadBytes)
	if err != nil {
		return nil, platformItem{}, err
	}
	detail, err := h.api.CreatePlatform(ctx, input, diagrams)
	if err != nil {
		return nil, platformItem{}, openlane.APIError(err)
	}
	return nil, mapPlatformDetail(*detail), nil
}

func (h *handlers) updatePlatform(ctx context.Context, _ *mcp.CallToolRequest, in updatePlatformInput) (*mcp.CallToolResult, platformItem, error) {
	if strings.TrimSpace(in.ID) == "" {
		return nil, platformItem{}, errIDRequired
	}
	input := graphclient.UpdatePlatformInput{}
	if err := applyUpdatePlatformFields(ctx, h, &input, in); err != nil {
		return nil, platformItem{}, err
	}
	diagrams, err := decodePlatformDiagrams(in.ArchitectureDiagrams, in.DataFlowDiagrams, in.TrustBoundaryDiagrams, h.maxUploadBytes)
	if err != nil {
		return nil, platformItem{}, err
	}
	detail, err := h.api.UpdatePlatform(ctx, in.ID, input, diagrams)
	if err != nil {
		return nil, platformItem{}, openlane.APIError(err)
	}
	return nil, mapPlatformDetail(*detail), nil
}

func applyCreatePlatformFields(ctx context.Context, h *handlers, input *graphclient.CreatePlatformInput, in createPlatformInput) error {
	if err := applyPlatformScalarFieldsCreate(input, in); err != nil {
		return err
	}
	for _, role := range []struct {
		key  string
		val  string
		name string
	}{
		{"business", in.BusinessOwner, in.BusinessOwnerName},
		{"technical", in.TechnicalOwner, in.TechnicalOwnerName},
		{"security", in.SecurityOwner, in.SecurityOwnerName},
		{"internal", in.InternalOwner, in.InternalOwnerName},
	} {
		write, err := h.buildPlatformOwnerWrite(ctx, role.val, role.name)
		if err != nil {
			return err
		}
		if err := setPlatformOwnerCreate(input, role.key, write); err != nil {
			return err
		}
	}
	input.AssetIDs = in.AssetIDs
	input.EntityIDs = in.EntityIDs
	input.SourceAssetIDs = in.SourceAssetIDs
	input.SourceEntityIDs = in.SourceEntityIDs
	input.OutOfScopeAssetIDs = in.OutOfScopeAssetIDs
	input.OutOfScopeVendorIDs = in.OutOfScopeVendorIDs
	input.ControlIDs = in.ControlIDs
	input.EvidenceIDs = in.EvidenceIDs
	input.ApplicableFrameworkIDs = in.ApplicableFrameworkIDs
	input.RiskIDs = in.RiskIDs
	input.TaskIDs = in.TaskIDs
	return nil
}

func applyUpdatePlatformFields(ctx context.Context, h *handlers, input *graphclient.UpdatePlatformInput, in updatePlatformInput) error {
	if err := applyPlatformScalarFieldsUpdate(input, in); err != nil {
		return err
	}
	for _, role := range []struct {
		key  string
		val  string
		name string
	}{
		{"business", in.BusinessOwner, in.BusinessOwnerName},
		{"technical", in.TechnicalOwner, in.TechnicalOwnerName},
		{"security", in.SecurityOwner, in.SecurityOwnerName},
		{"internal", in.InternalOwner, in.InternalOwnerName},
	} {
		write, err := h.buildPlatformOwnerWrite(ctx, role.val, role.name)
		if err != nil {
			return err
		}
		if write.userID != nil || write.groupID != nil || write.name != nil {
			if err := setPlatformOwnerUpdate(input, role.key, write); err != nil {
				return err
			}
		}
	}
	input.AddAssetIDs = in.AddAssetIDs
	input.RemoveAssetIDs = in.RemoveAssetIDs
	input.AddEntityIDs = in.AddEntityIDs
	input.RemoveEntityIDs = in.RemoveEntityIDs
	input.AddSourceAssetIDs = in.AddSourceAssetIDs
	input.RemoveSourceAssetIDs = in.RemoveSourceAssetIDs
	input.AddSourceEntityIDs = in.AddSourceEntityIDs
	input.RemoveSourceEntityIDs = in.RemoveSourceEntityIDs
	input.AddOutOfScopeAssetIDs = in.AddOutOfScopeAssetIDs
	input.RemoveOutOfScopeAssetIDs = in.RemoveOutOfScopeAssetIDs
	input.AddOutOfScopeVendorIDs = in.AddOutOfScopeVendorIDs
	input.RemoveOutOfScopeVendorIDs = in.RemoveOutOfScopeVendorIDs
	input.AddControlIDs = in.AddControlIDs
	input.RemoveControlIDs = in.RemoveControlIDs
	input.AddEvidenceIDs = in.AddEvidenceIDs
	input.RemoveEvidenceIDs = in.RemoveEvidenceIDs
	input.AddApplicableFrameworkIDs = in.AddApplicableFrameworkIDs
	input.RemoveApplicableFrameworkIDs = in.RemoveApplicableFrameworkIDs
	input.AddRiskIDs = in.AddRiskIDs
	input.RemoveRiskIDs = in.RemoveRiskIDs
	input.AddTaskIDs = in.AddTaskIDs
	input.RemoveTaskIDs = in.RemoveTaskIDs
	return nil
}

func applyPlatformScalarFieldsCreate(input *graphclient.CreatePlatformInput, in createPlatformInput) error {
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.BusinessPurpose != "" {
		input.BusinessPurpose = &in.BusinessPurpose
	}
	if in.ScopeStatement != "" {
		input.ScopeStatement = &in.ScopeStatement
	}
	if in.TrustBoundaryDescription != "" {
		input.TrustBoundaryDescription = &in.TrustBoundaryDescription
	}
	if in.DataFlowSummary != "" {
		input.DataFlowSummary = &in.DataFlowSummary
	}
	if in.Status != "" {
		status, err := parsePlatformStatus(in.Status)
		if err != nil {
			return err
		}
		input.Status = status
	}
	if in.PhysicalLocation != "" {
		input.PhysicalLocation = &in.PhysicalLocation
	}
	if in.Region != "" {
		input.Region = &in.Region
	}
	if in.SourceType != "" {
		st := enums.SourceType(in.SourceType)
		input.SourceType = &st
	}
	if in.SourceIdentifier != "" {
		input.SourceIdentifier = &in.SourceIdentifier
	}
	if in.CostCenter != "" {
		input.CostCenter = &in.CostCenter
	}
	if in.EstimatedMonthlyCost != nil {
		input.EstimatedMonthlyCost = in.EstimatedMonthlyCost
	}
	if in.PurchaseDate != "" {
		dt, err := parseDateTime(in.PurchaseDate)
		if err != nil {
			return err
		}
		input.PurchaseDate = dt
	}
	if in.ExternalReferenceID != "" {
		input.ExternalReferenceID = &in.ExternalReferenceID
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	setOptionalString(&input.PlatformKindName, in.PlatformKindName)
	setOptionalString(&input.PlatformDataClassificationName, in.PlatformDataClassificationName)
	setOptionalString(&input.EnvironmentName, in.EnvironmentName)
	setOptionalString(&input.ScopeName, in.ScopeName)
	setOptionalString(&input.AccessModelName, in.AccessModelName)
	setOptionalString(&input.EncryptionStatusName, in.EncryptionStatusName)
	setOptionalString(&input.SecurityTierName, in.SecurityTierName)
	setOptionalString(&input.CriticalityName, in.CriticalityName)
	return nil
}

func applyPlatformScalarFieldsUpdate(input *graphclient.UpdatePlatformInput, in updatePlatformInput) error {
	if in.Name != "" {
		input.Name = &in.Name
	}
	if in.Description != "" {
		input.Description = &in.Description
	}
	if in.BusinessPurpose != "" {
		input.BusinessPurpose = &in.BusinessPurpose
	}
	if in.ScopeStatement != "" {
		input.ScopeStatement = &in.ScopeStatement
	}
	if in.TrustBoundaryDescription != "" {
		input.TrustBoundaryDescription = &in.TrustBoundaryDescription
	}
	if in.DataFlowSummary != "" {
		input.DataFlowSummary = &in.DataFlowSummary
	}
	if in.Status != "" {
		status, err := parsePlatformStatus(in.Status)
		if err != nil {
			return err
		}
		input.Status = status
	}
	if in.PhysicalLocation != "" {
		input.PhysicalLocation = &in.PhysicalLocation
	}
	if in.Region != "" {
		input.Region = &in.Region
	}
	if in.SourceType != "" {
		st := enums.SourceType(in.SourceType)
		input.SourceType = &st
	}
	if in.SourceIdentifier != "" {
		input.SourceIdentifier = &in.SourceIdentifier
	}
	if in.CostCenter != "" {
		input.CostCenter = &in.CostCenter
	}
	if in.EstimatedMonthlyCost != nil {
		input.EstimatedMonthlyCost = in.EstimatedMonthlyCost
	}
	if in.PurchaseDate != "" {
		dt, err := parseDateTime(in.PurchaseDate)
		if err != nil {
			return err
		}
		input.PurchaseDate = dt
	}
	if in.ExternalReferenceID != "" {
		input.ExternalReferenceID = &in.ExternalReferenceID
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	setOptionalString(&input.PlatformKindName, in.PlatformKindName)
	setOptionalString(&input.PlatformDataClassificationName, in.PlatformDataClassificationName)
	setOptionalString(&input.EnvironmentName, in.EnvironmentName)
	setOptionalString(&input.ScopeName, in.ScopeName)
	setOptionalString(&input.AccessModelName, in.AccessModelName)
	setOptionalString(&input.EncryptionStatusName, in.EncryptionStatusName)
	setOptionalString(&input.SecurityTierName, in.SecurityTierName)
	setOptionalString(&input.CriticalityName, in.CriticalityName)
	return nil
}

func parsePlatformStatus(s string) (*enums.PlatformStatus, error) {
	status := enums.ToPlatformStatus(strings.TrimSpace(s))
	if status == nil || *status == enums.PlatformStatusInvalid {
		return nil, fmt.Errorf("invalid platform status %q (expected ACTIVE, INACTIVE, or RETIRED)", s)
	}
	return status, nil
}

func setOptionalString(target **string, value string) {
	if strings.TrimSpace(value) != "" {
		v := strings.TrimSpace(value)
		*target = &v
	}
}

func decodePlatformDiagrams(architecture, dataFlow, trustBoundary []platformDiagramInput, maxBytes int64) (openlane.PlatformDiagramUploads, error) {
	arch, err := decodeDiagramUploads(architecture, "architecture_diagrams", maxBytes)
	if err != nil {
		return openlane.PlatformDiagramUploads{}, err
	}
	flow, err := decodeDiagramUploads(dataFlow, "data_flow_diagrams", maxBytes)
	if err != nil {
		return openlane.PlatformDiagramUploads{}, err
	}
	trust, err := decodeDiagramUploads(trustBoundary, "trust_boundary_diagrams", maxBytes)
	if err != nil {
		return openlane.PlatformDiagramUploads{}, err
	}
	return openlane.PlatformDiagramUploads{
		Architecture:  arch,
		DataFlow:      flow,
		TrustBoundary: trust,
	}, nil
}

func decodeDiagramUploads(files []platformDiagramInput, field string, maxBytes int64) ([]*graphql.Upload, error) {
	if len(files) == 0 {
		return nil, nil
	}
	evidenceFiles := make([]openlane.EvidenceFile, 0, len(files))
	for _, f := range files {
		evidenceFiles = append(evidenceFiles, openlane.EvidenceFile{
			Filename:      f.Filename,
			ContentType:   f.ContentType,
			ContentBase64: f.ContentBase64,
		})
	}
	uploads, err := openlane.DecodeEvidenceUploads(evidenceFiles, maxBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", field, err)
	}
	return uploads, nil
}
