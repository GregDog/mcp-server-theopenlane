# Tools

Read tools are always registered. Write and delete tools are registered only when their mode is enabled.

With all modes enabled there are **42 tools** (27 read, 10 write, 5 delete).

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
| `openlane_policy_get` | `GetInternalPolicyByID` |

MCP names use “policy”; the Openlane type is `InternalPolicy`.

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
| `openlane_risk_create` | `CreateRisk` |
| `openlane_risk_update` | `UpdateRisk` |
| `openlane_task_create` | `CreateTask` |
| `openlane_task_update` | `UpdateTask` |

Create tools require the Openlane mandatory fields (`ref_code` for controls, `name` for evidence/policies/risks, `title` for tasks). Update tools require `id` and at least one field to change.

Evidence file uploads use `files[]` objects with `filename`, optional `content_type`, and `content_base64`. Default max decoded size is 10 MiB per file (`OPENLANE_MCP_MAX_UPLOAD_BYTES`). Presigned download URLs are still omitted from read responses.

## Delete tools (opt-in)

Enabled with `OPENLANE_ALLOW_DELETE=true` or `openlane-mcp serve --allow-delete`. Independent of write mode.

| Tool | Openlane client |
| --- | --- |
| `openlane_control_delete` | `DeleteControl` |
| `openlane_evidence_delete` | `DeleteEvidence` |
| `openlane_policy_delete` | `DeleteInternalPolicy` |
| `openlane_risk_delete` | `DeleteRisk` |
| `openlane_task_delete` | `DeleteTask` |

Delete tools require `id` and return `deleted_id`.

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

Do not expect MCP tools to surface these until the official client Get selection sets include them.

Verified against `github.com/theopenlane/go-client` v0.14.0.
