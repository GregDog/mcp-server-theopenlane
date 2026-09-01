# Tools

All tools are read-only. Arguments are validated by the MCP Go SDK from Go structs.

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

Verified against `github.com/theopenlane/go-client` v0.14.0, the latest stable release at implementation time.
