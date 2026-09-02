package tools

import (
	"fmt"
	"strings"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"
)

type workflowListInput struct {
	listInput
	Name         string `json:"name,omitempty" jsonschema:"Filter definitions whose name contains this text (case-insensitive)."`
	SchemaType   string `json:"schema_type,omitempty" jsonschema:"Filter by workflow schema type (e.g. InternalPolicy, Control)."`
	WorkflowKind string `json:"workflow_kind,omitempty" jsonschema:"Filter by workflow kind: APPROVAL, LIFECYCLE, or NOTIFICATION."`
	Active       *bool  `json:"active,omitempty" jsonschema:"Filter by active state."`
	Draft        *bool  `json:"draft,omitempty" jsonschema:"Filter by draft state."`
	TrackedField string `json:"tracked_field,omitempty" jsonschema:"Filter definitions that track this field name."`
}

type workflowInstanceListInput struct {
	listInput
	DefinitionID string `json:"definition_id,omitempty" jsonschema:"Filter instances for this workflow definition ID."`
	State        string `json:"state,omitempty" jsonschema:"Filter by instance state: RUNNING, PAUSED, COMPLETED, or FAILED."`
	ObjectType   string `json:"object_type,omitempty" jsonschema:"Filter by target object type (e.g. InternalPolicy, Control, Evidence)."`
	ObjectID     string `json:"object_id,omitempty" jsonschema:"Filter by target object ID. When set, object_type is required. When omitted, object_type alone matches any instance linked to that object type."`
}

type workflowAssignmentListInput struct {
	listInput
	Status     string `json:"status,omitempty" jsonschema:"Filter by assignment status: PENDING, APPROVED, REJECTED, or CHANGES_REQUESTED."`
	InstanceID string `json:"instance_id,omitempty" jsonschema:"Filter assignments for this workflow instance ID."`
}

type workflowMetadataInput struct {
	SchemaType string `json:"schema_type,omitempty" jsonschema:"Return metadata for this object type only (e.g. InternalPolicy)."`
}

func buildWorkflowDefinitionWhere(in workflowListInput) *graphclient.WorkflowDefinitionWhereInput {
	var w graphclient.WorkflowDefinitionWhereInput
	has := false
	if s := strings.TrimSpace(in.Name); s != "" {
		w.NameContainsFold = &s
		has = true
	}
	if s := strings.TrimSpace(in.SchemaType); s != "" {
		w.SchemaTypeContainsFold = &s
		has = true
	}
	if s := strings.TrimSpace(in.WorkflowKind); s != "" {
		k := enums.WorkflowKind(s)
		w.WorkflowKind = &k
		has = true
	}
	if in.Active != nil {
		w.Active = in.Active
		has = true
	}
	if in.Draft != nil {
		w.Draft = in.Draft
		has = true
	}
	if s := strings.TrimSpace(in.TrackedField); s != "" {
		w.TrackedFieldsHas = &s
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func buildWorkflowDefinitionSearchWhere(query string) *graphclient.WorkflowDefinitionWhereInput {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil
	}
	return &graphclient.WorkflowDefinitionWhereInput{
		Or: []*graphclient.WorkflowDefinitionWhereInput{
			{NameContainsFold: &q},
			{DescriptionContainsFold: &q},
		},
	}
}

func buildWorkflowInstanceWhere(in workflowInstanceListInput) (*graphclient.WorkflowInstanceWhereInput, error) {
	var w graphclient.WorkflowInstanceWhereInput
	has := false
	if s := strings.TrimSpace(in.DefinitionID); s != "" {
		w.WorkflowDefinitionID = &s
		has = true
	}
	if s := strings.TrimSpace(in.State); s != "" {
		st := enums.WorkflowInstanceState(s)
		w.State = &st
		has = true
	}
	objType := strings.TrimSpace(in.ObjectType)
	objID := strings.TrimSpace(in.ObjectID)
	if objID != "" && objType == "" {
		return nil, fmt.Errorf("object_type is required when object_id is set")
	}
	if objType != "" {
		if objID == "" {
			if err := applyWorkflowInstanceObjectTypeFilter(&w, objType); err != nil {
				return nil, err
			}
		} else if err := applyWorkflowInstanceObjectFilter(&w, objType, objID); err != nil {
			return nil, err
		}
		has = true
	}
	if !has {
		return nil, nil
	}
	return &w, nil
}

func applyWorkflowInstanceObjectTypeFilter(w *graphclient.WorkflowInstanceWhereInput, objectType string) error {
	notNil := true
	switch normalizeWorkflowObjectType(objectType) {
	case "InternalPolicy", "Policy":
		w.InternalPolicyIDNotNil = &notNil
	case "Control":
		w.ControlIDNotNil = &notNil
	case "Evidence":
		w.EvidenceIDNotNil = &notNil
	case "Risk":
		w.RiskIDNotNil = &notNil
	case "Finding":
		w.FindingIDNotNil = &notNil
	case "Task":
		w.TaskIDNotNil = &notNil
	case "Assessment":
		w.AssessmentIDNotNil = &notNil
	case "Subcontrol":
		w.SubcontrolIDNotNil = &notNil
	case "ActionPlan":
		w.ActionPlanIDNotNil = &notNil
	case "Procedure":
		w.ProcedureIDNotNil = &notNil
	case "Campaign":
		w.CampaignIDNotNil = &notNil
	case "CampaignTarget":
		w.CampaignTargetIDNotNil = &notNil
	case "Platform":
		w.PlatformIDNotNil = &notNil
	case "Remediation":
		w.RemediationIDNotNil = &notNil
	case "Vulnerability":
		w.VulnerabilityIDNotNil = &notNil
	case "IdentityHolder":
		w.IdentityHolderIDNotNil = &notNil
	case "AssessmentResponse":
		w.AssessmentResponseIDNotNil = &notNil
	default:
		return fmt.Errorf("unsupported object_type %q", objectType)
	}
	return nil
}

func applyWorkflowInstanceObjectFilter(w *graphclient.WorkflowInstanceWhereInput, objectType, objectID string) error {
	switch normalizeWorkflowObjectType(objectType) {
	case "InternalPolicy", "Policy":
		w.InternalPolicyID = &objectID
	case "Control":
		w.ControlID = &objectID
	case "Evidence":
		w.EvidenceID = &objectID
	case "Risk":
		w.RiskID = &objectID
	case "Finding":
		w.FindingID = &objectID
	case "Task":
		w.TaskID = &objectID
	case "Assessment":
		w.AssessmentID = &objectID
	case "Subcontrol":
		w.SubcontrolID = &objectID
	case "ActionPlan":
		w.ActionPlanID = &objectID
	case "Procedure":
		w.ProcedureID = &objectID
	case "Campaign":
		w.CampaignID = &objectID
	case "CampaignTarget":
		w.CampaignTargetID = &objectID
	case "Platform":
		w.PlatformID = &objectID
	case "Remediation":
		w.RemediationID = &objectID
	case "Vulnerability":
		w.VulnerabilityID = &objectID
	case "IdentityHolder":
		w.IdentityHolderID = &objectID
	case "AssessmentResponse":
		w.AssessmentResponseID = &objectID
	default:
		return fmt.Errorf("unsupported object_type %q", objectType)
	}
	return nil
}

func normalizeWorkflowObjectType(t string) string {
	t = strings.TrimSpace(t)
	if strings.EqualFold(t, "policy") {
		return "InternalPolicy"
	}
	return t
}

func buildWorkflowAssignmentWhere(in workflowAssignmentListInput) *graphclient.WorkflowAssignmentWhereInput {
	var w graphclient.WorkflowAssignmentWhereInput
	has := false
	if s := strings.TrimSpace(in.Status); s != "" {
		st := enums.WorkflowAssignmentStatus(s)
		w.Status = &st
		has = true
	}
	if s := strings.TrimSpace(in.InstanceID); s != "" {
		w.WorkflowInstanceID = &s
		has = true
	}
	if !has {
		return nil
	}
	return &w
}
