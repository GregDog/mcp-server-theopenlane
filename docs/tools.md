# Tools

Read tools are always registered. Write and delete tools are registered only when their mode is enabled.

With all modes enabled there are **66 tools** (40 read, 20 write, 6 delete).

Arguments are validated by the MCP Go SDK from Go structs.

## Pagination

List and search tools accept:

| Field | Default | Notes |
| --- | --- | --- |
| `limit` | 20 | Clamped to 1–50 |
| `cursor` | | Opaque `next_cursor` from a previous page |

List responses:

```json
{
  "items": [],
  "next_cursor": null,
  "has_more": false,
  "total_count": 0
}
```

`total_count` is the value Openlane already returns on the connection. The server does not issue a second query to compute it.

Get tools require `id`.

## Relationship summaries

Enriched get tools return compact relationship blocks:

```json
{
  "findings": {
    "count": 3,
    "items": [
      {"id": "...", "name": "...", "severity": "high", "open": true}
    ]
  }
}
```

Each get tool fetches at most eight related items per relationship (plus `total_count` from Openlane when available). Omitted summaries mean the relation could not be loaded or is empty.

## Controls

| Tool | Openlane client |
| --- | --- |
| `openlane_controls_list` | `GetControls` |
| `openlane_controls_search` | `GetControls` with `ControlWhereInput` |
| `openlane_control_get` | `GetControlByID` + bounded relation lists |

`openlane_control_get` includes assessment methods/objectives and summaries of programs, evidence, findings, risks, and implementations.

## Programs

| Tool | Openlane client |
| --- | --- |
| `openlane_programs_list` | `GetPrograms` |
| `openlane_program_get` | `GetProgramByID` + bounded relation lists |

`openlane_programs_list` filters:

| Field | Maps to |
| --- | --- |
| `name` | `NameContainsFold` |

`openlane_program_get` includes auditor metadata and summaries of evidence, findings, risks, tasks, and remediations.

## Evidence

| Tool | Openlane client |
| --- | --- |
| `openlane_evidence_list` | `GetEvidences` |
| `openlane_evidence_get` | `GetEvidenceByID` + program summary |

`openlane_evidence_list` filters:

| Field | Maps to |
| --- | --- |
| `program_id` | `HasProgramsWith` |
| `control_id` | `HasControlsWith` |

List responses omit file contents. Get responses include `file_ids` and file metadata (`id` only) but not presigned URLs.

## Policies

| Tool | Openlane client |
| --- | --- |
| `openlane_policies_list` | `GetInternalPolicies` |
| `openlane_policies_awaiting_approval` | `GetInternalPolicies` + `GetMyWorkflowAssignments` (+ detail lookups) |
| `openlane_policy_get` | `GetInternalPolicyByID` |

MCP names use “policy”; the Openlane type is `InternalPolicy`.

### Native vs workflow policy approval

Openlane supports two approval paths for policies:

| Path | How to detect | Approve with |
| --- | --- | --- |
| **Native** | `InternalPolicy.status` is `NEEDS_APPROVAL` (no `WorkflowAssignment`) | `openlane_policy_approve` |
| **Workflow** | Pending `WorkflowAssignment` on a `WorkflowInstance` linked to the policy | `openlane_workflow_assignment_approve` |

Use `openlane_policies_awaiting_approval` to list both paths in one call. `openlane_policy_submit_for_approval` starts the native path only; it does not create workflow instances.

`openlane_policies_list` filters:

| Field | Maps to |
| --- | --- |
| `status` | `Status` (`DRAFT`, `NEEDS_APPROVAL`, `APPROVED`, `PUBLISHED`, `ARCHIVED`) |

## Risks

| Tool | Openlane client |
| --- | --- |
| `openlane_risks_list` | `GetRisks` |
| `openlane_risk_get` | `GetRiskByID` + bounded relation lists |

`openlane_risks_list` filters:

| Field | Maps to |
| --- | --- |
| `program_id` | `HasProgramsWith` |
| `entity_id` | `HasEntitiesWith` |
| `control_id` | `HasControlsWith` |
| `status` | `RiskStatus` |

`openlane_risk_get` includes owner/stakeholder IDs and summaries of controls, programs, findings, and remediations.

## Findings

| Tool | Openlane client |
| --- | --- |
| `openlane_findings_list` | `GetFindings` |
| `openlane_finding_get` | `GetFindingByID` + bounded control summary |

`openlane_findings_list` filters:

| Field | Maps to |
| --- | --- |
| `program_id` | `HasProgramsWith` |
| `assessment_id` | `AssessmentID` |
| `open` | `Open` |
| `status` | `FindingStatusName` |
| `severity` | `Severity` |

`openlane_finding_get` includes description, impact, recommendation, remediations, vulnerabilities, and linked controls. Raw scanner payloads (`rawPayload`, `metadata`) are omitted from list responses.

## Assessments

| Tool | Openlane client |
| --- | --- |
| `openlane_assessments_list` | `GetAssessments` |
| `openlane_assessment_get` | `GetAssessmentByID` + bounded responses/findings |

`openlane_assessments_list` filters:

| Field | Maps to |
| --- | --- |
| `assessment_type` | `AssessmentType` |

## Control implementations

| Tool | Openlane client |
| --- | --- |
| `openlane_control_implementations_list` | `GetControlImplementations` |
| `openlane_control_implementation_get` | `GetControlImplementationByID` |

`openlane_control_implementations_list` filters:

| Field | Maps to |
| --- | --- |
| `control_id` | `HasControlsWith` |
| `status` | `Status` |
| `verified` | `Verified` |

Implementation get responses include control/subcontrol `ref_code` values (not control IDs in the generated client selection set).

## Standards

| Tool | Openlane client |
| --- | --- |
| `openlane_standards_list` | `GetStandards` |
| `openlane_standard_get` | `GetStandardByID` |

## Tasks

| Tool | Openlane client |
| --- | --- |
| `openlane_tasks_list` | `GetTasks` |
| `openlane_task_get` | `GetTaskByID` |

## Entities

| Tool | Openlane client |
| --- | --- |
| `openlane_entities_list` | `GetEntities` |
| `openlane_entity_get` | `GetEntityByID` + bounded assets/risks/findings |

Entities represent vendors and other third parties in Openlane. There is no separate vendor API in the current Go client.

`openlane_entities_list` filters:

| Field | Maps to |
| --- | --- |
| `risk_rating` | `RiskRating` |
| `tier` | `Tier` |
| `approved_for_use` | `ApprovedForUse` |
| `questionnaire_status` | `EntitySecurityQuestionnaireStatusName` |
| `next_review_before` | `NextReviewAtLTE` |
| `has_soc2` | `HasSoc2` |
| `sso_enforced` | `SsoEnforced` |
| `mfa_enforced` | `MfaEnforced` |

`openlane_entity_get` exposes vendor/security/commercial fields (SOC 2, SSO/MFA, contract, spend, reviews, owner) and bounded relationship summaries.

## Assets

| Tool | Openlane client |
| --- | --- |
| `openlane_assets_list` | `GetAssets` |
| `openlane_asset_get` | `GetAssetByID` |

## Contacts

| Tool | Openlane client |
| --- | --- |
| `openlane_contacts_list` | `GetContacts` |
| `openlane_contact_get` | `GetContactByID` |

## Groups and users

Lookup tools for name→ID resolution when authoring workflows or reassigning assignments. Always-on (read).

| Tool | Openlane client |
| --- | --- |
| `openlane_groups_list` | `GetGroups` |
| `openlane_group_get` | `GetGroupByID` |
| `openlane_users_list` | `GetUsers` |
| `openlane_user_get` | `GetUserByID` |

`openlane_groups_list` filters:

| Field | Maps to |
| --- | --- |
| `name` | `NameContainsFold` **or** `DisplayNameContainsFold` |

`openlane_users_list` filters:

| Field | Maps to |
| --- | --- |
| `name` | `DisplayNameContainsFold` **or** `FirstNameContainsFold` **or** `LastNameContainsFold` |
| `email` | `EmailContainsFold` |

Workflow write handlers resolve group/user names the same way and **error** on 0 or many matches (no silent pick).

## Workflows

Workflow tools expose Openlane's approval and automation engine. A **WorkflowDefinition** describes triggers, conditions, and actions. A **WorkflowInstance** is a running execution against a target object. A **WorkflowAssignment** is an approval task routed to a user or group.

| Tool | Openlane client |
| --- | --- |
| `openlane_workflows_list` | `GetWorkflowDefinitions` |
| `openlane_workflows_search` | `GetWorkflowDefinitions` with `WorkflowDefinitionWhereInput` |
| `openlane_workflow_get` | `GetWorkflowDefinitionByID` |
| `openlane_workflow_instances_list` | `GetWorkflowInstances` |
| `openlane_workflow_instance_get` | `GetWorkflowInstanceByID` + custom GraphQL for `currentActionIndex` / `proposalPreview` |
| `openlane_workflow_assignments_list` | `GetMyWorkflowAssignments` |
| `openlane_workflow_assignment_get` | Custom GraphQL (go-client Get omits `dueAt` and targets) |
| `openlane_workflow_metadata_get` | Custom GraphQL `workflowMetadata` query |

`openlane_workflows_list` filters:

| Field | Maps to |
| --- | --- |
| `name` | `NameContainsFold` |
| `schema_type` | `SchemaTypeContainsFold` |
| `workflow_kind` | `WorkflowKind` |
| `active` | `Active` |
| `draft` | `Draft` |
| `tracked_field` | `TrackedFieldsHas` |

`openlane_workflow_instances_list` filters:

| Field | Maps to |
| --- | --- |
| `definition_id` | `WorkflowDefinitionID` |
| `state` | `State` (`RUNNING`, `PAUSED`, `COMPLETED`, `FAILED`) |
| `object_type` | Object-type predicate on `WorkflowInstanceWhereInput` (e.g. `internalPolicyIDNotNil` for `InternalPolicy` or `Policy`) |
| `object_id` | Specific object ID field (requires `object_type`) |

`openlane_workflow_assignments_list` filters:

| Field | Maps to |
| --- | --- |
| `status` | `Status` (`PENDING`, `APPROVED`, `REJECTED`, `CHANGES_REQUESTED`) |
| `instance_id` | `WorkflowInstanceID` |

Get responses include a `summary` field with a concise plain-English interpretation. Definition get responses include the full `definition_json` document.

### Example MCP calls

List active policy workflows:

```json
{"name": "openlane_workflows_list", "arguments": {"schema_type": "InternalPolicy", "active": true}}
```

List policies awaiting approval (native + workflow):

```json
{"name": "openlane_policies_awaiting_approval", "arguments": {"limit": 50}}
```

List policies in NEEDS_APPROVAL:

```json
{"name": "openlane_policies_list", "arguments": {"status": "NEEDS_APPROVAL"}}
```

List policy workflow instances (type only):

```json
{"name": "openlane_workflow_instances_list", "arguments": {"object_type": "InternalPolicy", "state": "PAUSED"}}
```

Show my pending approvals:

```json
{"name": "openlane_workflow_assignments_list", "arguments": {"status": "PENDING"}}
```

Check which fields can be approval-gated on policies:

```json
{"name": "openlane_workflow_metadata_get", "arguments": {"schema_type": "InternalPolicy"}}
```

Explain a workflow definition:

```json
{"name": "openlane_workflow_get", "arguments": {"id": "WFD..."}}
```

Approve a native InternalPolicy (no workflow assignment):

```json
{"name": "openlane_policy_approve", "arguments": {"id": "01…", "confirm": true}}
```

Approve a workflow assignment:

```json
{"name": "openlane_workflow_assignment_approve", "arguments": {"id": "01…", "confirm": true}}
```

### PRE_COMMIT vs POST_COMMIT

Approval workflows may use `approvalTiming` in `definition_json`:

- `PRE_COMMIT` — eligible field changes are staged in a proposal before applying
- `POST_COMMIT` — changes commit first; approval runs afterward

### Built-in resolver targeting

Resolver keys are returned per object type from `openlane_workflow_metadata_get` (for example `POLICY_APPROVER`, `CONTROL_OWNER`, `OBJECT_CREATOR`). Do not hard-code eligibility; always query metadata before authoring workflows.

### Native InternalPolicy approval vs WorkflowAssignment approval

These are **two different paths**. Do not mix mutations.

**Native InternalPolicy lifecycle** (no WorkflowDefinition required):

- `openlane_policy_submit_for_approval` — `updateInternalPolicy` status `NEEDS_APPROVAL` (notifies the Approver group)
- `openlane_policy_approve` — `updateInternalPolicy` status `APPROVED` (core `HookStatusApproval` requires Approver or Delegate group membership)
- `openlane_policy_publish` — `updateInternalPolicy` status `PUBLISHED`
- `openlane_policy_return_to_draft` — `updateInternalPolicy` status `DRAFT` (there is no `rejectInternalPolicy` mutation)

There is no native approve/reject mutation or task. Natural language “Approve this policy” should use `openlane_policy_approve` when the InternalPolicy is `NEEDS_APPROVAL` and there is **no** WorkflowAssignment.

**WorkflowAssignment** (only when a real WorkflowInstance/Assignment exists):

- `openlane_workflow_assignment_approve` — `approveWorkflowAssignment`
- `openlane_workflow_assignment_reject` — `rejectWorkflowAssignment` (reason required)
- `openlane_workflow_assignment_request_changes` — custom GraphQL (`requestChangesWorkflowAssignment`; not in go-client v0.14.0)
- `openlane_workflow_assignment_reassign` — custom GraphQL (`reassignWorkflowAssignment`)

Natural language “Approve this workflow assignment” uses `openlane_workflow_assignment_approve`. Never call assignment mutations against a policy that has no instance/assignment. Do not create a WorkflowDefinition to emulate native document approval.

## Write tools (opt-in)

Enabled with `OPENLANE_ALLOW_WRITE=true` or `openlane-mcp serve --allow-write`.

| Tool | Openlane client |
| --- | --- |
| `openlane_control_create` | `CreateControl` |
| `openlane_control_update` | `UpdateControl` |
| `openlane_evidence_create` | `CreateEvidence` (optional `files[]` base64 uploads) |
| `openlane_evidence_update` | `UpdateEvidence` (optional `files[]` base64 uploads) |
| `openlane_policy_create` | `CreateInternalPolicy` |
| `openlane_policy_update` | `UpdateInternalPolicy` |
| `openlane_policy_submit_for_approval` | `UpdateInternalPolicy` status `NEEDS_APPROVAL` (`confirm` required) |
| `openlane_policy_approve` | `UpdateInternalPolicy` status `APPROVED` (`confirm` required) |
| `openlane_policy_publish` | `UpdateInternalPolicy` status `PUBLISHED` (`confirm` required) |
| `openlane_policy_return_to_draft` | `UpdateInternalPolicy` status `DRAFT` (`confirm` required) |
| `openlane_risk_create` | `CreateRisk` |
| `openlane_risk_update` | `UpdateRisk` |
| `openlane_task_create` | `CreateTask` |
| `openlane_task_update` | `UpdateTask` |
| `openlane_workflow_create` | `CreateWorkflowDefinition` (`confirm` required; validated against `workflowMetadata` + copied JSON schema) |
| `openlane_workflow_update` | `UpdateWorkflowDefinition` Get-then-patch (`confirm` required) |
| `openlane_workflow_assignment_approve` | `ApproveWorkflowAssignment` (`confirm` required) |
| `openlane_workflow_assignment_reject` | `RejectWorkflowAssignment` (`confirm` required) |
| `openlane_workflow_assignment_request_changes` | Custom GraphQL `requestChangesWorkflowAssignment` (`confirm` required) |
| `openlane_workflow_assignment_reassign` | Custom GraphQL `reassignWorkflowAssignment` (`confirm` required) |

Create tools require the Openlane mandatory fields (`ref_code` for controls, `name` for evidence/policies/risks, `title` for tasks). Update tools require `id` and at least one field to change.

Workflow and native policy lifecycle writes additionally require `confirm: true`. If `confirm` is false or omitted, the tool returns a before/after (or current/requested status) preview with `error: confirmation required` and does not mutate.

Example: create a PRE_COMMIT policy publication approval workflow without hand-authoring JSON:

```json
{
  "name": "openlane_workflow_create",
  "arguments": {
    "name": "Policy publication approval",
    "schema_type": "InternalPolicy",
    "workflow_kind": "APPROVAL",
    "approval_timing": "PRE_COMMIT",
    "triggers": [{"operation": "UPDATE", "fields": ["status"]}],
    "condition_cel": "object.status == \"PUBLISHED\"",
    "actions": [{"type": "REQUEST_APPROVAL", "targets": [{"type": "RESOLVER", "resolver_key": "POLICY_APPROVER"}]}],
    "confirm": true
  }
}
```

Native policy approval (no workflow):

```json
{"name": "openlane_policy_approve", "arguments": {"id": "01…", "confirm": true}}
```

Approve a workflow assignment:

```json
{"name": "openlane_workflow_assignment_approve", "arguments": {"id": "01…", "confirm": true}}
```

Evidence file uploads use `files[]` objects with `filename`, optional `content_type`, and `content_base64`. Default max decoded size is 10 MiB per file (`OPENLANE_MCP_MAX_UPLOAD_BYTES`). Presigned download URLs are still omitted from read responses.

`cancelWorkflowInstance` / `forceCompleteWorkflowInstance` exist in the GraphQL schema but have no go-client methods and are **not** exposed as MCP tools.

## Delete tools (opt-in)

Enabled with `OPENLANE_ALLOW_DELETE=true` or `openlane-mcp serve --allow-delete`. Independent of write mode.

| Tool | Openlane client |
| --- | --- |
| `openlane_control_delete` | `DeleteControl` |
| `openlane_evidence_delete` | `DeleteEvidence` |
| `openlane_policy_delete` | `DeleteInternalPolicy` |
| `openlane_risk_delete` | `DeleteRisk` |
| `openlane_task_delete` | `DeleteTask` |
| `openlane_workflow_delete` | `DeleteWorkflowDefinition` (`confirm` required; re-reads to verify deletion) |

Delete tools require `id` and return `deleted_id`. `openlane_workflow_delete` also requires `confirm: true`.

## Generated Get-query gaps (go-client v0.14.0)

The following fields exist on Openlane core models or WhereInputs but are **not** returned by the stock `Get*ByID` queries in `github.com/theopenlane/go-client` v0.14.0:

- Program observation period and fieldwork dates
- Entity domains, provided services, risk score coverage, review frequency, entity type name
- Risk decision, residual score, review fields; linked entities/assets from `risk_get` (no reverse WhereInput)
- Control implementation status on control get; nested owner objects
- Evidence review frequency
- Assessment questionnaire config; responses not nested in assessment Get
- Finding `displayID` and review fields
- Control implementation related control IDs (ref codes only)
- Workflow instance `currentActionIndex` and `proposalPreview` on stock `GetWorkflowInstanceByID` (MCP uses custom GraphQL)
- Workflow assignment `dueAt` and `workflowAssignmentTargets` on stock `GetWorkflowAssignmentByID` (MCP uses custom GraphQL)
- `workflowMetadata` query (no go-client method in v0.14.0; MCP uses custom GraphQL)
- `requestChangesWorkflowAssignment` / `reassignWorkflowAssignment` (schema yes; no go-client methods in v0.14.0; MCP uses custom GraphQL)

Do not expect MCP tools to surface these until the official client Get selection sets include them.

Verified against `github.com/theopenlane/go-client` v0.14.0.
