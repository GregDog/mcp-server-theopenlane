package openlane

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
)

// workflow.definition.json is copied from github.com/theopenlane/core/v2@v2.3.0
// jsonschema/workflow.definition.json. Validation below mirrors that schema's
// enums and required shape, then checks live workflowMetadata eligibility.
//
//go:embed workflow.definition.json
var WorkflowDefinitionSchemaJSON []byte

var (
	workflowKinds = map[string]struct{}{
		string(enums.WorkflowKindApproval):     {},
		string(enums.WorkflowKindLifecycle):    {},
		string(enums.WorkflowKindNotification): {},
	}
	workflowTimings = map[string]struct{}{
		string(enums.WorkflowApprovalTimingPreCommit):  {},
		string(enums.WorkflowApprovalTimingPostCommit): {},
	}
	workflowActionTypes = map[string]struct{}{
		string(enums.WorkflowActionTypeApproval):         {},
		string(enums.WorkflowActionTypeNotification):     {},
		string(enums.WorkflowActionTypeWebhook):          {},
		string(enums.WorkflowActionTypeFieldUpdate):      {},
		string(enums.WorkflowActionTypeIntegration):      {},
		string(enums.WorkflowActionTypeReassignApproval): {},
		string(enums.WorkflowActionTypeSendEmail):        {},
		string(enums.WorkflowActionTypeCreateObject):     {},
		string(enums.WorkflowActionTypeReview):           {},
	}
	workflowTargetTypes = map[string]struct{}{
		string(enums.WorkflowTargetTypeUser):     {},
		string(enums.WorkflowTargetTypeGroup):    {},
		string(enums.WorkflowTargetTypeRole):     {},
		string(enums.WorkflowTargetTypeResolver): {},
		string(enums.WorkflowTargetTypeChannel):  {},
	}
)

// ValidateWorkflowDefinition checks a document against the copied JSON schema
// enums and against live workflowMetadata. Unknown types/fields/edges/resolvers
// are rejected; nothing is substituted.
func ValidateWorkflowDefinition(doc *models.WorkflowDefinitionDocument, meta *WorkflowMetadata) error {
	if doc == nil {
		return fmt.Errorf("definition_json is required")
	}
	if strings.TrimSpace(doc.Name) == "" {
		return fmt.Errorf("workflow name is required")
	}
	if strings.TrimSpace(doc.SchemaType) == "" {
		return fmt.Errorf("schema_type is required")
	}
	kind := string(doc.WorkflowKind)
	if kind == "" {
		return fmt.Errorf("workflow_kind is required")
	}
	if _, ok := workflowKinds[kind]; !ok {
		return fmt.Errorf("unsupported workflow_kind %q", kind)
	}
	if t := string(doc.ApprovalTiming); t != "" {
		if _, ok := workflowTimings[t]; !ok {
			return fmt.Errorf("unsupported approval_timing %q", t)
		}
	}
	if len(doc.Triggers) == 0 {
		return fmt.Errorf("at least one trigger is required")
	}
	if len(doc.Actions) == 0 {
		return fmt.Errorf("at least one action is required")
	}

	obj, err := lookupWorkflowObjectType(meta, doc.SchemaType)
	if err != nil {
		return err
	}

	for i, t := range doc.Triggers {
		if strings.TrimSpace(t.Operation) == "" && strings.TrimSpace(t.Interval) == "" {
			return fmt.Errorf("triggers[%d]: operation or interval is required", i)
		}
		for _, f := range t.Fields {
			if !eligibleField(obj, f) {
				return fmt.Errorf("unknown eligible field %q for %s", f, obj.Type)
			}
		}
		for _, e := range t.Edges {
			if !eligibleEdge(obj, e) {
				return fmt.Errorf("unknown eligible edge %q for %s", e, obj.Type)
			}
		}
	}

	for i, a := range doc.Actions {
		if strings.TrimSpace(a.Key) == "" {
			return fmt.Errorf("actions[%d]: key is required", i)
		}
		if _, ok := workflowActionTypes[a.Type]; !ok {
			return fmt.Errorf("actions[%d]: unsupported action type %q", i, a.Type)
		}
		targets, err := actionTargets(a.Params)
		if err != nil {
			return fmt.Errorf("actions[%d]: %w", i, err)
		}
		for _, tgt := range targets {
			if _, ok := workflowTargetTypes[tgt.Type]; !ok {
				return fmt.Errorf("actions[%d]: unsupported target type %q", i, tgt.Type)
			}
			if tgt.Type == string(enums.WorkflowTargetTypeResolver) {
				if tgt.ResolverKey == "" {
					return fmt.Errorf("actions[%d]: resolver_key is required for RESOLVER targets", i)
				}
				if !eligibleResolver(obj, tgt.ResolverKey) {
					return fmt.Errorf("unknown resolver_key %q for %s", tgt.ResolverKey, obj.Type)
				}
			}
		}
	}
	return nil
}

func lookupWorkflowObjectType(meta *WorkflowMetadata, schemaType string) (WorkflowObjectTypeMetadata, error) {
	if meta == nil || len(meta.ObjectTypes) == 0 {
		return WorkflowObjectTypeMetadata{}, fmt.Errorf("workflowMetadata returned no object types")
	}
	want := strings.TrimSpace(schemaType)
	for _, t := range meta.ObjectTypes {
		if strings.EqualFold(t.Type, want) {
			return t, nil
		}
		if strings.EqualFold(want, "Policy") && strings.EqualFold(t.Type, "InternalPolicy") {
			return t, nil
		}
	}
	return WorkflowObjectTypeMetadata{}, fmt.Errorf("unknown schema_type %q (not in workflowMetadata)", schemaType)
}

func eligibleField(obj WorkflowObjectTypeMetadata, name string) bool {
	for _, f := range obj.EligibleFields {
		if strings.EqualFold(f.Name, name) {
			return true
		}
	}
	return false
}

func eligibleEdge(obj WorkflowObjectTypeMetadata, name string) bool {
	for _, e := range obj.EligibleEdges {
		if strings.EqualFold(e, name) {
			return true
		}
	}
	return false
}

func eligibleResolver(obj WorkflowObjectTypeMetadata, key string) bool {
	for _, k := range obj.ResolverKeys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

// WorkflowActionTarget is a USER/GROUP/RESOLVER target inside action params.
type WorkflowActionTarget struct {
	Type        string `json:"type"`
	ID          string `json:"id,omitempty"`
	ResolverKey string `json:"resolver_key,omitempty"`
	Destination string `json:"destination,omitempty"`
}

func actionTargets(params json.RawMessage) ([]WorkflowActionTarget, error) {
	if len(params) == 0 {
		return nil, nil
	}
	var p struct {
		Targets []WorkflowActionTarget `json:"targets"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("params: %w", err)
	}
	return p.Targets, nil
}
