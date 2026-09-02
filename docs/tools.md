# Tools

Read tools are always registered. Write and delete tools are registered only when their mode is enabled.

Arguments are validated by the MCP Go SDK from Go structs.

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

## Controls

| Tool | Openlane client |
| --- | --- |
| `openlane_controls_list` | `GetControls` |
| `openlane_controls_search` | `GetControls` with `ControlWhereInput` |
| `openlane_control_get` | `GetControlByID` |

Search matches ref code, title, or description (`ContainsFold`). The current client has no `ControlSearch` operation.

## Programs

| Tool | Openlane client |
| --- | --- |
| `openlane_programs_list` | `GetPrograms` |
| `openlane_program_get` | `GetProgramByID` |

## Evidence

| Tool | Openlane client |
| --- | --- |
| `openlane_evidence_list` | `GetEvidences` |
| `openlane_evidence_get` | `GetEvidenceByID` |

File contents and presigned URLs are omitted.

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
| `openlane_risk_get` | `GetRiskByID` |

## Standards

| Tool | Openlane client |
| --- | --- |
| `openlane_standards_list` | `GetStandards` |
| `openlane_standard_get` | `GetStandardByID` |

## Write tools (opt-in)

Enabled with `OPENLANE_ALLOW_WRITE=true` or `openlane-mcp serve --allow-write`. Evidence writes are metadata only; file uploads are not supported.

| Tool | Openlane client |
| --- | --- |
| `openlane_control_create` | `CreateControl` |
| `openlane_control_update` | `UpdateControl` |
| `openlane_evidence_create` | `CreateEvidence` (no files) |
| `openlane_evidence_update` | `UpdateEvidence` (no files) |
| `openlane_policy_create` | `CreateInternalPolicy` |
| `openlane_policy_update` | `UpdateInternalPolicy` |
| `openlane_risk_create` | `CreateRisk` |
| `openlane_risk_update` | `UpdateRisk` |
| `openlane_task_create` | `CreateTask` |
| `openlane_task_update` | `UpdateTask` |

Create tools require the Openlane mandatory fields (`ref_code` for controls, `name` for evidence/policies/risks, `title` for tasks). Update tools require `id` and at least one field to change.

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

Verified against `github.com/theopenlane/go-client` v0.14.0, the latest stable release at implementation time.
