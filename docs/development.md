# Development

Requires **Go 1.27** or later ([install](https://go.dev/dl/)). See [CONTRIBUTING.md](../CONTRIBUTING.md) for the fork → branch → pull request workflow.

Unit tests (`make test`) do not require an Openlane API token. Live smoke tests are optional — see `make test-access` below.

```bash
git clone https://github.com/GregDog/mcp-server-theopenlane.git
cd mcp-server-theopenlane
make test
make build
```

The binary is written to `bin/openlane-mcp`.

```bash
export OPENLANE_API_TOKEN="tola_..."
./bin/openlane-mcp serve
```

Useful targets:

| Target | Action |
| --- | --- |
| `make test` | `go test ./...` |
| `make vet` | `go vet ./...` |
| `make fmt` | fail if `gofmt` would change files |
| `make build` | build `openlane-mcp` |
| `make test-access` | live API smoke test (requires `.env`) |
| `make vuln` | `govulncheck ./...` |

Do not commit `.env` files, `.local/` test artifacts, or real tokens.

## Local testing (stdio)

1. Copy `.env.example` to `.env` and set `OPENLANE_API_TOKEN` (and `OPENLANE_ORGANIZATION_ID` for multi-org PATs).
2. Build the binary: `make build`
3. Verify API access: `make test-access`
4. Open this project in Cursor — `.cursor/mcp.json` runs `bin/openlane-mcp serve` over stdio (the binary auto-loads `.env` from the repo root). Run `bash scripts/mcp-verify.sh` to confirm tool registration.

To test write or delete tools, also set in `.env`:

```bash
OPENLANE_ALLOW_WRITE=true   # create/update tools
OPENLANE_ALLOW_DELETE=true  # delete tools (independent of write)
```

Then rebuild and reload Cursor (**Developer: Reload Window**). Check **Output → MCP Logs** for `allow_write=true` and `allow_delete=true` on startup.

In Cursor: **Settings → MCP** and confirm `openlane` is enabled.

Example prompts:

- "List the first 5 Openlane controls"
- "List entities with high risk rating that do not enforce MFA"
- "List policies awaiting approval" (`openlane_policies_awaiting_approval`)
- "List policies in NEEDS_APPROVAL" (`openlane_policies_list` with `status`)
- "Approve this policy" (native InternalPolicy in NEEDS_APPROVAL; requires write mode and confirm)
- "Approve this workflow assignment" (requires a real WorkflowAssignment; write mode and confirm)
- "Submit this policy for approval"
- "List my pending workflow assignments"
- "What fields can be approval-gated on policies?" (`openlane_workflow_metadata_get`)
- "Get program `<id>` — what evidence, findings, and tasks are linked?"
- "List open high-severity findings for program `<id>`"
- "Get entity `<id>` — SOC 2, SSO/MFA, contract dates, and owner"
- "Get control `<id>` — what evidence and implementations support it?"
- "Create evidence named MCP test with a text file" (requires write mode; see evidence uploads below)
- "Create vendor Acme Corp with tier HIGH" (requires write mode; `openlane_entity_create`)
- "Upload a logo for entity `<id>`" (requires write mode; `openlane_entity_update` with `logo` base64)
- "Delete evidence `<id>`" (requires delete mode)

## Compliance context manual testing

After `make build` and configuring `.env`, exercise the enriched read tools:

1. **Program context** — `openlane_programs_list` with `{"name":"PCI"}` then `openlane_program_get` with a program ID. Confirm `evidence`, `findings`, `risks`, `tasks`, and `remediations` summaries appear with `count` and up to eight `items`.
2. **Vendor review** — `openlane_entities_list` with `{"risk_rating":"high"}` or `{"mfa_enforced":false}`. Then `openlane_entity_get` and confirm SOC 2, SSO/MFA, contract, spend, and review fields.
3. **Findings** — `openlane_findings_list` with `{"program_id":"<id>","open":true,"severity":"high"}`. Then `openlane_finding_get` for remediations and vulnerabilities.
4. **Control implementation** — `openlane_control_implementations_list` with `{"control_id":"<id>"}` then `openlane_control_implementation_get`.
5. **Assessment** — `openlane_assessments_list` then `openlane_assessment_get` for campaigns, responses, and findings.

## Workflow and policy approval manual testing

After `make build` and configuring `.env` (enable write mode for lifecycle/assignment actions):

1. **Policies awaiting approval** — `openlane_policies_awaiting_approval` with `{"limit":50}`. Confirm `path` is `native` or `workflow` and the summary counts match.
2. **Native policy filter** — `openlane_policies_list` with `{"status":"NEEDS_APPROVAL"}`.
3. **Workflow instances by type** — `openlane_workflow_instances_list` with `{"object_type":"InternalPolicy","state":"PAUSED"}` (no `object_id` required).
4. **My assignments** — `openlane_workflow_assignments_list` with `{"status":"PENDING"}`.
5. **Metadata** — `openlane_workflow_metadata_get` with `{"schema_type":"InternalPolicy"}` before authoring workflows.
6. **Native lifecycle** (write mode, `confirm: true`) — `openlane_policy_submit_for_approval`, `openlane_policy_approve`, `openlane_policy_return_to_draft`, `openlane_policy_publish` on a test policy.
7. **Workflow assignment** (write mode, `confirm: true`) — only when a real assignment exists: `openlane_workflow_assignment_approve` or `openlane_workflow_assignment_reject` with `reason`.
8. **Group/user lookup** — `openlane_groups_list` / `openlane_users_list` with name filters before `openlane_workflow_assignment_reassign` or `openlane_control_update` `owner_id` / `delegate_id` (user ids resolve to managed personal groups).

In Cursor, reload MCP after rebuilding (`make build`, then **Developer: Reload Window** or restart the `openlane` server under **Settings → MCP**). Use **Output → MCP Logs** if a tool fails. Compare results with the Openlane UI for the same object IDs.

## HTTP transport testing

For Streamable HTTP instead of stdio:

1. Set in `.env`:

```bash
OPENLANE_MCP_TRANSPORT=http
OPENLANE_MCP_HTTP_ADDR=127.0.0.1:8090   # loopback only; do not use 0.0.0.0 without a trusted auth proxy
```

2. Start the server in a separate terminal:

```bash
bash scripts/mcp-http.sh
```

3. Point Cursor at the URL — see `examples/cursor-http.mcp.json`:

```json
{
  "mcpServers": {
    "openlane": {
      "url": "http://127.0.0.1:8090"
    }
  }
}
```

HTTP mode has no built-in MCP authentication. Do not expose it directly to the public internet. Use a trusted reverse proxy for remote deployment. See [security.md](security.md).

## Evidence and entity logo upload testing

Evidence file uploads are tested via MCP write tools (`openlane_evidence_create` / `openlane_evidence_update`) with a `files[]` array. Entity logos use `openlane_entity_create` / `openlane_entity_update` with a single `logo` object or `logo_remote_url`:

```json
{
  "id": "01…",
  "logo": {
    "filename": "acme.png",
    "content_type": "image/png",
    "content_base64": "..."
  }
}
```

Evidence uploads (`files[]` array):

```json
{
  "filename": "sample.txt",
  "content_type": "text/plain",
  "content_base64": "..."
}
```

Local test artifacts can live under `.local/` (gitignored). Do not commit test files or API tokens.

MCP read tools return `file_ids` but not presigned download URLs. Verify downloads in the Openlane UI or with a direct API client.

### Linking evidence to controls

`openlane_evidence_create` accepts optional `control_ids` (org-owned controls only). `openlane_controls_search` with `linkable_only: true` filters to controls with `owner_id` set. System catalog copies are rejected ([core#1647](https://github.com/theopenlane/core/pull/1647)). Probes in this repo: `scripts/test-org-control-link`, `scripts/test-program-controls`, `scripts/test-linkable-controls`.

When testing CSV or other large file uploads through an agent, use a direct MCP client or script for the create call — relaying `content_base64` through multi-step agent tool chains often corrupts the payload.

## Tool inventory

| Mode | Count |
| --- | --- |
| Read (always on) | 40 |
| Write (opt-in) | 22 |
| Delete (opt-in) | 6 |
| **Total** | **68** |

See [tools.md](tools.md) for the full mapping to Openlane client methods.

## Maintainer commits

GitHub shows commits as Verified when they are SSH-signed (or GPG/S/MIME). Contributors are not required to sign commits.

Registry publishing is documented in [publishing.md](publishing.md).
