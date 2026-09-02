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

List responses omit file contents. Get responses include `file_ids` but not presigned URLs.

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

## Tasks

| Tool | Openlane client |
| --- | --- |
| `openlane_tasks_list` | `GetTasks` |
| `openlane_task_get` | `GetTaskByID` |

## Entities

| Tool | Openlane client |
| --- | --- |
| `openlane_entities_list` | `GetEntities` |
| `openlane_entity_get` | `GetEntityByID` |

Entities represent vendors and other third parties in Openlane. There is no separate vendor API in the current Go client.

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

Verified against `github.com/theopenlane/go-client` v0.14.0, the latest stable release at implementation time.
