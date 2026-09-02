package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

const errConfirmationRequired = "confirmation required"

type workflowTriggerSpec struct {
	Operation  string   `json:"operation,omitempty" jsonschema:"Trigger operation: CREATE, UPDATE, or DELETE."`
	Fields     []string `json:"fields,omitempty" jsonschema:"Field names that trigger evaluation (must be eligible in workflowMetadata)."`
	Edges      []string `json:"edges,omitempty" jsonschema:"Edge names that trigger evaluation."`
	Expression string   `json:"expression,omitempty" jsonschema:"Optional CEL expression gate for this trigger."`
}

type workflowTargetSpec struct {
	Type        string `json:"type" jsonschema:"Target type: USER, GROUP, ROLE, RESOLVER, or CHANNEL."`
	ID          string `json:"id,omitempty" jsonschema:"User or group ID. Prefer this over a name."`
	GroupName   string `json:"group_name,omitempty" jsonschema:"Group name to resolve to an ID (errors on 0 or many matches)."`
	UserEmail   string `json:"user_email,omitempty" jsonschema:"User email to resolve to an ID (errors on 0 or many matches)."`
	ResolverKey string `json:"resolver_key,omitempty" jsonschema:"Resolver key for RESOLVER targets (must be in workflowMetadata)."`
}

type workflowActionSpec struct {
	Key           string               `json:"key,omitempty" jsonschema:"Unique action key within the workflow."`
	Type          string               `json:"type" jsonschema:"Action type, e.g. REQUEST_APPROVAL, NOTIFY."`
	Description   string               `json:"description,omitempty" jsonschema:"Human-readable action description."`
	When          string               `json:"when,omitempty" jsonschema:"Optional CEL expression that gates this action."`
	Targets       []workflowTargetSpec `json:"targets,omitempty" jsonschema:"Approval or notify targets."`
	RequiredCount *int                 `json:"required_count,omitempty" jsonschema:"Quorum count for REQUEST_APPROVAL (0 means all targets)."`
	Required      *bool                `json:"required,omitempty" jsonschema:"Whether this approval is required. Defaults to true for REQUEST_APPROVAL."`
	Label         string               `json:"label,omitempty" jsonschema:"Optional display label for the approval."`
}

type createWorkflowInput struct {
	Name            string                `json:"name" jsonschema:"Workflow definition name."`
	Description     string                `json:"description,omitempty" jsonschema:"Optional description."`
	SchemaType      string                `json:"schema_type" jsonschema:"Target schema type from workflowMetadata (e.g. InternalPolicy)."`
	WorkflowKind    string                `json:"workflow_kind,omitempty" jsonschema:"APPROVAL, LIFECYCLE, or NOTIFICATION. Defaults to APPROVAL."`
	ApprovalTiming  string                `json:"approval_timing,omitempty" jsonschema:"PRE_COMMIT or POST_COMMIT. Defaults to PRE_COMMIT for gating."`
	CooldownSeconds *int64                `json:"cooldown_seconds,omitempty" jsonschema:"Suppress duplicate triggers within this window."`
	Draft           *bool                 `json:"draft,omitempty" jsonschema:"Create as a draft."`
	Active          *bool                 `json:"active,omitempty" jsonschema:"Whether the definition is active."`
	Triggers        []workflowTriggerSpec `json:"triggers,omitempty" jsonschema:"Triggers. Ignored when definition_json is set."`
	ConditionCEL    string                `json:"condition_cel,omitempty" jsonschema:"Optional CEL condition. Ignored when definition_json is set."`
	Actions         []workflowActionSpec  `json:"actions,omitempty" jsonschema:"Actions. Ignored when definition_json is set."`
	DefinitionJSON  json.RawMessage       `json:"definition_json,omitempty" jsonschema:"Raw WorkflowDefinitionDocument JSON for power users. Still validated."`
	Confirm         bool                  `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type updateWorkflowInput struct {
	ID               string              `json:"id" jsonschema:"Workflow definition ID to update."`
	Name             string              `json:"name,omitempty" jsonschema:"Updated name."`
	Description      string              `json:"description,omitempty" jsonschema:"Updated description."`
	Active           *bool               `json:"active,omitempty" jsonschema:"Set active (enable) or inactive (disable)."`
	Draft            *bool               `json:"draft,omitempty" jsonschema:"Updated draft flag."`
	CooldownSeconds  *int64              `json:"cooldown_seconds,omitempty" jsonschema:"Updated cooldown in seconds."`
	AddTarget        *workflowTargetSpec `json:"add_target,omitempty" jsonschema:"Append a target to the first REQUEST_APPROVAL action."`
	SetRequiredCount *int                `json:"set_required_count,omitempty" jsonschema:"Set required_count on the first REQUEST_APPROVAL action."`
	AddNotifyAction  *workflowActionSpec `json:"add_notify_action,omitempty" jsonschema:"Append a NOTIFY action."`
	ReplaceAction    *workflowActionSpec `json:"replace_action,omitempty" jsonschema:"Replace the action whose key matches."`
	DefinitionJSON   json.RawMessage     `json:"definition_json,omitempty" jsonschema:"Replace the entire definition document (still validated)."`
	Confirm          bool                `json:"confirm" jsonschema:"Must be true to persist. If false, returns a preview only."`
}

type workflowWriteResult struct {
	Confirmed bool                    `json:"confirmed"`
	Error     string                  `json:"error,omitempty"`
	Before    *workflowDefinitionItem `json:"before,omitempty"`
	After     *workflowDefinitionItem `json:"after,omitempty"`
	Summary   string                  `json:"summary,omitempty"`
}

func registerWriteWorkflows(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_create",
		Title:       "Create an Openlane workflow definition",
		Description: "Create a generic WorkflowDefinition (triggers, conditions, actions). This is not how native InternalPolicy approval works — use openlane_policy_submit_for_approval for that. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.createWorkflow)

	addTool(server, &mcp.Tool{
		Name:        "openlane_workflow_update",
		Title:       "Update an Openlane workflow definition",
		Description: "Update a workflow definition by ID. Supports scalar patches, activate/disable via active, and structured edits (add_target, set_required_count, add_notify_action, replace_action). Always Get-then-patch. Requires write mode and confirm: true.",
		Annotations: writeAnnotations(),
	}, h.updateWorkflow)
}

func (h *handlers) createWorkflow(ctx context.Context, _ *mcp.CallToolRequest, in createWorkflowInput) (*mcp.CallToolResult, workflowWriteResult, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, workflowWriteResult{}, errNameRequired
	}
	meta, err := h.api.GetWorkflowMetadata(ctx)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	doc, err := h.buildCreateDocument(ctx, in)
	if err != nil {
		return nil, workflowWriteResult{}, err
	}
	if err := openlane.ValidateWorkflowDefinition(doc, meta); err != nil {
		return nil, workflowWriteResult{}, err
	}
	obj, err := resolveSchemaType(meta, in.SchemaType, doc.SchemaType)
	if err != nil {
		return nil, workflowWriteResult{}, err
	}
	doc.SchemaType = obj.Type
	preview := workflowDefinitionItem{
		Name:            in.Name,
		Description:     firstNonEmpty(in.Description, doc.Description),
		SchemaType:      obj.Type,
		WorkflowKind:    firstNonEmpty(in.WorkflowKind, string(doc.WorkflowKind)),
		Active:          openlane.Deref(in.Active),
		Draft:           openlane.Deref(in.Draft),
		CooldownSeconds: openlane.Deref(in.CooldownSeconds),
		TrackedFields:   trackedFieldsFromDoc(doc),
		DefinitionJSON:  doc,
		Summary:         summarizeWorkflowDefinition(in.Name, obj.Type, firstNonEmpty(in.WorkflowKind, string(doc.WorkflowKind)), openlane.Deref(in.Active), openlane.Deref(in.Draft), 0, openlane.Deref(in.CooldownSeconds), trackedFieldsFromDoc(doc), doc),
	}
	out := workflowWriteResult{After: &preview, Summary: "Would create: " + preview.Summary}
	if !in.Confirm {
		out.Error = errConfirmationRequired
		return nil, out, nil
	}

	kind := enums.WorkflowKind(firstNonEmpty(in.WorkflowKind, string(doc.WorkflowKind), string(enums.WorkflowKindApproval)))
	input := graphclient.CreateWorkflowDefinitionInput{
		Name:            in.Name,
		SchemaType:      obj.Type,
		WorkflowKind:    kind,
		DefinitionJSON:  doc,
		TrackedFields:   trackedFieldsFromDoc(doc),
		Active:          in.Active,
		Draft:           in.Draft,
		CooldownSeconds: in.CooldownSeconds,
	}
	if d := firstNonEmpty(in.Description, doc.Description); d != "" {
		input.Description = &d
	}
	created, err := h.api.CreateWorkflowDefinition(ctx, input)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	got, err := h.api.GetWorkflowDefinitionByID(ctx, created.CreateWorkflowDefinition.WorkflowDefinition.ID)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	item := mapWorkflowDefinition(got.WorkflowDefinition)
	if item.Name != in.Name || item.SchemaType != obj.Type {
		return nil, workflowWriteResult{}, fmt.Errorf("workflow persisted but name or schema_type did not match")
	}
	out.Confirmed = true
	out.After = &item
	out.Summary = "Created: " + item.Summary
	return nil, out, nil
}

func (h *handlers) updateWorkflow(ctx context.Context, _ *mcp.CallToolRequest, in updateWorkflowInput) (*mcp.CallToolResult, workflowWriteResult, error) {
	if in.ID == "" {
		return nil, workflowWriteResult{}, errIDRequired
	}
	existing, err := h.api.GetWorkflowDefinitionByID(ctx, in.ID)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	before := mapWorkflowDefinition(existing.WorkflowDefinition)
	meta, err := h.api.GetWorkflowMetadata(ctx)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	doc, err := cloneWorkflowDoc(existing.WorkflowDefinition.DefinitionJSON)
	if err != nil {
		return nil, workflowWriteResult{}, err
	}
	if len(in.DefinitionJSON) > 0 {
		var replacement models.WorkflowDefinitionDocument
		if err := json.Unmarshal(in.DefinitionJSON, &replacement); err != nil {
			return nil, workflowWriteResult{}, fmt.Errorf("definition_json: %w", err)
		}
		doc = &replacement
	}
	if err := h.applyWorkflowPatches(ctx, doc, in); err != nil {
		return nil, workflowWriteResult{}, err
	}
	name := firstNonEmpty(in.Name, before.Name)
	doc.Name = name
	if err := openlane.ValidateWorkflowDefinition(doc, meta); err != nil {
		return nil, workflowWriteResult{}, err
	}

	patch := graphclient.UpdateWorkflowDefinitionInput{DefinitionJSON: doc}
	if in.Name != "" {
		patch.Name = &in.Name
	}
	if in.Description != "" {
		patch.Description = &in.Description
	}
	if in.Active != nil {
		patch.Active = in.Active
	}
	if in.Draft != nil {
		patch.Draft = in.Draft
	}
	if in.CooldownSeconds != nil {
		patch.CooldownSeconds = in.CooldownSeconds
	}
	if fields := trackedFieldsFromDoc(doc); len(fields) > 0 {
		patch.TrackedFields = fields
	}

	afterPreview := before
	afterPreview.Name = name
	if in.Description != "" {
		afterPreview.Description = in.Description
	}
	if in.Active != nil {
		afterPreview.Active = *in.Active
	}
	if in.Draft != nil {
		afterPreview.Draft = *in.Draft
	}
	if in.CooldownSeconds != nil {
		afterPreview.CooldownSeconds = *in.CooldownSeconds
	}
	afterPreview.DefinitionJSON = doc
	afterPreview.TrackedFields = trackedFieldsFromDoc(doc)
	afterPreview.Summary = summarizeWorkflowDefinition(afterPreview.Name, afterPreview.SchemaType, afterPreview.WorkflowKind, afterPreview.Active, afterPreview.Draft, afterPreview.Revision, afterPreview.CooldownSeconds, afterPreview.TrackedFields, doc)

	out := workflowWriteResult{Before: &before, After: &afterPreview, Summary: "Would update: " + afterPreview.Summary}
	if !in.Confirm {
		out.Error = errConfirmationRequired
		return nil, out, nil
	}

	if _, err := h.api.UpdateWorkflowDefinition(ctx, in.ID, patch); err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	got, err := h.api.GetWorkflowDefinitionByID(ctx, in.ID)
	if err != nil {
		return nil, workflowWriteResult{}, openlane.APIError(err)
	}
	item := mapWorkflowDefinition(got.WorkflowDefinition)
	if in.Active != nil && item.Active != *in.Active {
		return nil, workflowWriteResult{}, fmt.Errorf("workflow persisted but active=%v, expected %v", item.Active, *in.Active)
	}
	if in.Name != "" && item.Name != in.Name {
		return nil, workflowWriteResult{}, fmt.Errorf("workflow persisted but name did not match")
	}
	out.Confirmed = true
	out.After = &item
	out.Summary = "Updated: " + item.Summary
	return nil, out, nil
}

func (h *handlers) buildCreateDocument(ctx context.Context, in createWorkflowInput) (*models.WorkflowDefinitionDocument, error) {
	if len(in.DefinitionJSON) > 0 {
		var doc models.WorkflowDefinitionDocument
		if err := json.Unmarshal(in.DefinitionJSON, &doc); err != nil {
			return nil, fmt.Errorf("definition_json: %w", err)
		}
		if doc.Name == "" {
			doc.Name = in.Name
		}
		if doc.SchemaType == "" {
			doc.SchemaType = in.SchemaType
		}
		return &doc, nil
	}
	if strings.TrimSpace(in.SchemaType) == "" {
		return nil, fmt.Errorf("schema_type is required")
	}
	if len(in.Triggers) == 0 {
		return nil, fmt.Errorf("triggers are required unless definition_json is set")
	}
	if len(in.Actions) == 0 {
		return nil, fmt.Errorf("actions are required unless definition_json is set")
	}
	kind := firstNonEmpty(in.WorkflowKind, string(enums.WorkflowKindApproval))
	timing := firstNonEmpty(in.ApprovalTiming, string(enums.WorkflowApprovalTimingPreCommit))
	doc := &models.WorkflowDefinitionDocument{
		Name:           in.Name,
		Description:    in.Description,
		SchemaType:     in.SchemaType,
		WorkflowKind:   enums.WorkflowKind(kind),
		ApprovalTiming: enums.WorkflowApprovalTiming(timing),
	}
	for _, t := range in.Triggers {
		doc.Triggers = append(doc.Triggers, models.WorkflowTrigger{
			Operation:  t.Operation,
			Fields:     t.Fields,
			Edges:      t.Edges,
			Expression: t.Expression,
		})
	}
	if strings.TrimSpace(in.ConditionCEL) != "" {
		doc.Conditions = []models.WorkflowCondition{{Expression: in.ConditionCEL}}
	}
	for i, a := range in.Actions {
		act, err := h.buildAction(ctx, a, i)
		if err != nil {
			return nil, err
		}
		doc.Actions = append(doc.Actions, act)
	}
	return doc, nil
}

func (h *handlers) applyWorkflowPatches(ctx context.Context, doc *models.WorkflowDefinitionDocument, in updateWorkflowInput) error {
	if in.AddTarget != nil {
		idx := firstApprovalAction(doc)
		if idx < 0 {
			return fmt.Errorf("add_target requires a REQUEST_APPROVAL action")
		}
		params, err := decodeApprovalParams(doc.Actions[idx].Params)
		if err != nil {
			return err
		}
		tgt, err := h.resolveTarget(ctx, *in.AddTarget)
		if err != nil {
			return err
		}
		params.Targets = append(params.Targets, tgt)
		raw, err := json.Marshal(params)
		if err != nil {
			return err
		}
		doc.Actions[idx].Params = raw
	}
	if in.SetRequiredCount != nil {
		idx := firstApprovalAction(doc)
		if idx < 0 {
			return fmt.Errorf("set_required_count requires a REQUEST_APPROVAL action")
		}
		params, err := decodeApprovalParams(doc.Actions[idx].Params)
		if err != nil {
			return err
		}
		params.RequiredCount = in.SetRequiredCount
		raw, err := json.Marshal(params)
		if err != nil {
			return err
		}
		doc.Actions[idx].Params = raw
	}
	if in.AddNotifyAction != nil {
		spec := *in.AddNotifyAction
		if spec.Type == "" {
			spec.Type = string(enums.WorkflowActionTypeNotification)
		}
		act, err := h.buildAction(ctx, spec, len(doc.Actions))
		if err != nil {
			return err
		}
		doc.Actions = append(doc.Actions, act)
	}
	if in.ReplaceAction != nil {
		if strings.TrimSpace(in.ReplaceAction.Key) == "" {
			return fmt.Errorf("replace_action.key is required")
		}
		act, err := h.buildAction(ctx, *in.ReplaceAction, 0)
		if err != nil {
			return err
		}
		found := false
		for i, existing := range doc.Actions {
			if existing.Key == act.Key {
				doc.Actions[i] = act
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("no action with key %q", in.ReplaceAction.Key)
		}
	}
	return nil
}

type approvalParams struct {
	Targets       []openlane.WorkflowActionTarget `json:"targets,omitempty"`
	Required      *bool                           `json:"required,omitempty"`
	Label         string                          `json:"label,omitempty"`
	RequiredCount *int                            `json:"required_count,omitempty"`
}

func decodeApprovalParams(raw json.RawMessage) (approvalParams, error) {
	var p approvalParams
	if len(raw) == 0 {
		return p, nil
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return p, fmt.Errorf("action params: %w", err)
	}
	return p, nil
}

func firstApprovalAction(doc *models.WorkflowDefinitionDocument) int {
	for i, a := range doc.Actions {
		if a.Type == string(enums.WorkflowActionTypeApproval) {
			return i
		}
	}
	return -1
}

func (h *handlers) buildAction(ctx context.Context, spec workflowActionSpec, index int) (models.WorkflowAction, error) {
	key := strings.TrimSpace(spec.Key)
	if key == "" {
		key = fmt.Sprintf("action_%d", index+1)
	}
	act := models.WorkflowAction{
		Key:         key,
		Type:        spec.Type,
		When:        spec.When,
		Description: spec.Description,
	}
	if len(spec.Targets) == 0 && spec.RequiredCount == nil && spec.Label == "" && spec.Required == nil {
		return act, nil
	}
	params := approvalParams{
		Required:      spec.Required,
		Label:         spec.Label,
		RequiredCount: spec.RequiredCount,
	}
	for _, t := range spec.Targets {
		tgt, err := h.resolveTarget(ctx, t)
		if err != nil {
			return models.WorkflowAction{}, err
		}
		params.Targets = append(params.Targets, tgt)
	}
	raw, err := json.Marshal(params)
	if err != nil {
		return models.WorkflowAction{}, err
	}
	act.Params = raw
	return act, nil
}

func (h *handlers) resolveTarget(ctx context.Context, spec workflowTargetSpec) (openlane.WorkflowActionTarget, error) {
	t := strings.ToUpper(strings.TrimSpace(spec.Type))
	tgt := openlane.WorkflowActionTarget{Type: t, ResolverKey: spec.ResolverKey}
	switch t {
	case string(enums.WorkflowTargetTypeGroup):
		id := spec.ID
		if id == "" && spec.GroupName != "" {
			resolved, err := h.resolveGroupID(ctx, spec.GroupName)
			if err != nil {
				return tgt, err
			}
			id = resolved
		}
		if id == "" {
			return tgt, fmt.Errorf("GROUP target requires id or group_name")
		}
		tgt.ID = id
	case string(enums.WorkflowTargetTypeUser):
		id := spec.ID
		if id == "" && spec.UserEmail != "" {
			resolved, err := h.resolveUserID(ctx, spec.UserEmail)
			if err != nil {
				return tgt, err
			}
			id = resolved
		}
		if id == "" {
			return tgt, fmt.Errorf("USER target requires id or user_email")
		}
		tgt.ID = id
	case string(enums.WorkflowTargetTypeResolver):
		if tgt.ResolverKey == "" {
			return tgt, fmt.Errorf("RESOLVER target requires resolver_key")
		}
	case string(enums.WorkflowTargetTypeRole), string(enums.WorkflowTargetTypeChannel):
		tgt.ID = spec.ID
	default:
		if t == "" {
			return tgt, fmt.Errorf("target type is required")
		}
	}
	return tgt, nil
}

func cloneWorkflowDoc(doc *models.WorkflowDefinitionDocument) (*models.WorkflowDefinitionDocument, error) {
	if doc == nil {
		return &models.WorkflowDefinitionDocument{}, nil
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var out models.WorkflowDefinitionDocument
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func trackedFieldsFromDoc(doc *models.WorkflowDefinitionDocument) []string {
	if doc == nil {
		return nil
	}
	seen := map[string]struct{}{}
	var out []string
	for _, t := range doc.Triggers {
		for _, f := range t.Fields {
			if _, ok := seen[f]; ok {
				continue
			}
			seen[f] = struct{}{}
			out = append(out, f)
		}
	}
	return out
}

func resolveSchemaType(meta *openlane.WorkflowMetadata, requested, fromDoc string) (openlane.WorkflowObjectTypeMetadata, error) {
	want := firstNonEmpty(requested, fromDoc)
	if meta == nil {
		return openlane.WorkflowObjectTypeMetadata{}, fmt.Errorf("workflowMetadata returned no object types")
	}
	for _, t := range meta.ObjectTypes {
		if strings.EqualFold(t.Type, want) {
			return t, nil
		}
		if strings.EqualFold(want, "Policy") && strings.EqualFold(t.Type, "InternalPolicy") {
			return t, nil
		}
	}
	return openlane.WorkflowObjectTypeMetadata{}, fmt.Errorf("unknown schema_type %q (not in workflowMetadata)", want)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
