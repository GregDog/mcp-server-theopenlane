package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type listInput struct {
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of items to return. Default 20, maximum 50."`
	Cursor string `json:"cursor,omitempty" jsonschema:"Opaque cursor from a previous list response."`
}

type getInput struct {
	ID string `json:"id" jsonschema:"Openlane object ID."`
}

type searchInput struct {
	Query  string `json:"query" jsonschema:"Text to match against control ref code, title, or description."`
	Limit  int    `json:"limit,omitempty" jsonschema:"Number of items to return. Default 20, maximum 50."`
	Cursor string `json:"cursor,omitempty" jsonschema:"Opaque cursor from a previous search response."`
}

type programListInput struct {
	listInput
	Name string `json:"name,omitempty" jsonschema:"Filter programs whose name contains this text (case-insensitive)."`
}

type entityListInput struct {
	listInput
	RiskRating          string `json:"risk_rating,omitempty" jsonschema:"Filter by vendor risk rating."`
	Tier                string `json:"tier,omitempty" jsonschema:"Filter by vendor tier."`
	ApprovedForUse      *bool  `json:"approved_for_use,omitempty" jsonschema:"Filter by approved-for-use flag."`
	QuestionnaireStatus string `json:"questionnaire_status,omitempty" jsonschema:"Filter by security questionnaire status name."`
	NextReviewBefore    string `json:"next_review_before,omitempty" jsonschema:"Filter entities with next review on or before this date (RFC3339 or YYYY-MM-DD)."`
	HasSoc2             *bool  `json:"has_soc2,omitempty" jsonschema:"Filter by SOC 2 assurance flag."`
	SsoEnforced         *bool  `json:"sso_enforced,omitempty" jsonschema:"Filter by SSO enforced flag."`
	MfaEnforced         *bool  `json:"mfa_enforced,omitempty" jsonschema:"Filter by MFA enforced flag."`
}

type riskListInput struct {
	listInput
	ProgramID string `json:"program_id,omitempty" jsonschema:"Filter risks linked to this program ID."`
	EntityID  string `json:"entity_id,omitempty" jsonschema:"Filter risks linked to this entity (vendor) ID."`
	ControlID string `json:"control_id,omitempty" jsonschema:"Filter risks linked to this control ID."`
	Status    string `json:"status,omitempty" jsonschema:"Filter by risk status."`
}

type findingListInput struct {
	listInput
	ProgramID    string `json:"program_id,omitempty" jsonschema:"Filter findings linked to this program ID."`
	AssessmentID string `json:"assessment_id,omitempty" jsonschema:"Filter findings for this assessment ID."`
	Open         *bool  `json:"open,omitempty" jsonschema:"Filter by open state."`
	Status       string `json:"status,omitempty" jsonschema:"Filter by finding status name."`
	Severity     string `json:"severity,omitempty" jsonschema:"Filter by severity."`
}

type implementationListInput struct {
	listInput
	ControlID string `json:"control_id,omitempty" jsonschema:"Filter implementations linked to this control ID."`
	Status    string `json:"status,omitempty" jsonschema:"Filter by implementation status."`
	Verified  *bool  `json:"verified,omitempty" jsonschema:"Filter by verified flag."`
}

type assessmentListInput struct {
	listInput
	AssessmentType string `json:"assessment_type,omitempty" jsonschema:"Filter by assessment type."`
}

type evidenceListInput struct {
	listInput
	ProgramID string `json:"program_id,omitempty" jsonschema:"Filter evidence linked to this program ID."`
	ControlID string `json:"control_id,omitempty" jsonschema:"Filter evidence linked to this control ID."`
}

func readOnly() *mcp.ToolAnnotations {
	destructive := false
	openWorld := true
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    true,
		DestructiveHint: &destructive,
		OpenWorldHint:   &openWorld,
		IdempotentHint:  true,
	}
}

func pageArgs(limit int, cursor string) (first int64, after *string) {
	return openlane.ClampLimit(limit), openlane.CursorPtr(cursor)
}
