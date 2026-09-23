# Openlane assignee and owner IDs

Openlane uses **three different ID meanings** that MCP tools and agents often conflate. Mixing them produces `UNAUTHORIZED` on writes or failed `openlane_user_get` on reads.

## The three meanings

| MCP / response field | Openlane field | ID type | Meaning |
| --- | --- | --- | --- |
| `owner_id` on controls, evidence, programs, risks (read) | `ownerID` | **Organization ID** | Which org owns the record. On controls: set on program-imported copies; required for evidence linking (`linkable_to_evidence: true`). **Not a person.** |
| `control_owner_id`, `delegate_id` on controls | `controlOwnerID`, `delegateID` | **Group ID** | Responsibility assignment. Usually a **managed personal group** (`{displayName} - {userID}`). **Not a user ID.** |
| `owner_id`, `delegate_id`, `stakeholder_id` on risks (read) | `stakeholderID`, `delegateID` | **Group ID** | Risk oversight / delegate groups. MCP does not expose write fields for these yet. |
| `approver_id`, `delegate_id` on policies (lifecycle preview) | `approverID`, `delegateID` | **Group ID** | Native policy approval. MCP lifecycle tools do not set these. |
| `openlane_users_list` / `openlane_user_get` | User | **User ID** | Org members only. **Cannot look up group IDs.** |
| Workflow `USER` target | user id | **User ID** | Correct as-is. |
| Workflow `GROUP` target | group id | **Group ID** | Use `openlane_groups_list` or `group_name`. |

Passing a **user ID** to Openlane `controlOwnerID` / `delegateID` returns `UNAUTHORIZED`. Passing a **group ID** to `openlane_user_get` returns no match.

## What this MCP implements today

| Surface | Write: user → group resolution | Read: group → user enrichment |
| --- | --- | --- |
| **Controls** (`openlane_control_create` / `openlane_control_update`) | Yes — `owner_id` / `delegate_id` accept user id, email, name, or group id | Yes — `control_owner` and `delegate` objects on list/search/get/update |
| **Risks** | No assignee write fields | No — raw `delegate_id` / `stakeholder_id` only |
| **Policies** | No approver/delegate write fields | Lifecycle preview exposes raw group ids only |
| **Workflows** | USER vs GROUP targets resolved separately | N/A |

Helpers live in `internal/tools/control_assignees.go`. When adding risk/policy assignee writes or read enrichment, **reuse** `resolveControlAssigneeGroupID` and `resolveGroupAssigneeSummary` — do not duplicate the pattern.

## Control writes (agents)

Prefer user email or name; MCP resolves to the managed personal group:

```json
{
  "id": "01KTY1WZBZZ17NK5C7CJXGJW97",
  "owner_id": "abhishek.gupta@nomupay.com",
  "delegate_id": "greg.knell@nomupay.com"
}
```

## Control reads (agents)

Use enriched objects — **do not** call `openlane_user_get` on `control_owner_id` or `delegate_id`:

```json
{
  "control_owner_id": "01KTY1J6WEC5EXBVSJJWW6QP68",
  "control_owner": {
    "group_display_name": "abhishek.gupta",
    "user_id": "01KRGYEMTKH5XJABCQZ5BJMXGC",
    "user_email": "abhishek.gupta@nomupay.com"
  },
  "delegate_id": "01M1EN3JHTHT3X5P2FJRTGBH00",
  "delegate": {
    "group_display_name": "Greg",
    "user_email": "greg.knell@nomupay.com"
  },
  "owner_id": "01KTY1J6M1WTQAPYC2T2MKREKD"
}
```

`owner_id` here is the **org** — ignore for “who owns this control” questions; use `control_owner`.

## Managed personal groups

Openlane creates one managed group per user, named like `abhishek.gupta - 01KRGYEMTKH5XJABCQZ5BJMXGC`. Resolution:

1. **Write path:** `openlane_users_list` → `GetGroups` with `isManaged` + `hasMembersWith.userID`.
2. **Read path:** parse `userID` from group name suffix → `GetOrgMembers` for email/display name.

## Probes

| Script | Purpose |
| --- | --- |
| `scripts/test-control-owner-update` | User id vs group id write behavior |
| `scripts/apply-control-owner-update` | Batch assign owner/delegate by group id |
| `scripts/test-org-control-link` | Evidence linking requires org-owned control |

## Related

- [tools.md](tools.md) — control list/get fields
- [development.md](development.md) — manual MCP verification
- Sibling `security-ai-platform` runbook: [openlane-assignee-ids.md](https://github.com/GregDog/security-ai-platform/blob/main/docs/runbooks/openlane-assignee-ids.md) (GRaCe operators)
