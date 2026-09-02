# Openlane MCP Server

[![CI](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/ci.yml/badge.svg)](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![CodeQL](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/codeql.yml/badge.svg)](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/codeql.yml)

A secure, open-source [Model Context Protocol](https://modelcontextprotocol.io/) server for the [Openlane](https://www.theopenlane.io/) GRC platform.

This project is **not** an official Openlane product and is not endorsed by theopenlane, Inc.

## Overview

`openlane-mcp` lets MCP clients such as Cursor and Claude Desktop query Openlane over **stdio** (default) or **Streamable HTTP**. It talks to Openlane Cloud or a self-hosted instance through the [official Openlane Go client](https://github.com/theopenlane/go-client).

```text
MCP Client
    → Openlane MCP Server (stdio or HTTP)
        → Official Openlane Go Client
            → Openlane API
```

The server is read-only by default. Write and delete tools are opt-in and independent. Openlane authorization still applies to every request.

## Features

- Openlane MCP access for programs, controls, evidence, policies, risks, standards, tasks, entities (vendors), assets, contacts, findings, assessments, control implementations, groups, users, and workflows
- Enriched get tools with vendor/security fields and compact relationship summaries
- List filters on entities, risks, findings, evidence, programs, assessments, implementations, and workflows
- Opt-in create/update tools for controls, evidence, policies, risks, tasks, workflow definitions, workflow assignments, and native policy lifecycle
- Opt-in delete tools for the same domains plus workflow definitions (except programs and standards)
- Openlane Cloud and self-hosted Openlane (configurable base URL)
- stdio transport (default) and opt-in Streamable HTTP transport
- Native Go binary
- Docker image (published with GitHub Releases)
- MCP Registry listing on tagged releases (`io.github.GregDog/mcp-server-theopenlane`)

## Quick Start

Create an Openlane API token or PAT in console developer settings. Organization tokens start with `tola_`. Personal access tokens start with `tolp_`.

```bash
export OPENLANE_API_TOKEN="tola_..."
# Optional for multi-org PATs:
export OPENLANE_ORGANIZATION_ID="..."

openlane-mcp serve
```

Then connect an MCP client. See [Client configuration](#client-configuration).

## Installation

### From source

```bash
go install github.com/GregDog/mcp-server-theopenlane/cmd/openlane-mcp@latest
```

Requires Go 1.27 or later.

### GitHub Releases

Binary archives will be published on tagged GitHub Releases (`linux`/`darwin` amd64+arm64, `windows` amd64) with SHA256 checksums.

### Docker

```bash
docker run --rm -i \
  -e OPENLANE_API_TOKEN \
  -e OPENLANE_ORGANIZATION_ID \
  ghcr.io/gregdog/mcp-server-theopenlane serve
```

Images are published with GitHub Releases to `ghcr.io/gregdog/mcp-server-theopenlane`.

## Client configuration

### Cursor (stdio, recommended)

This repository's `.cursor/mcp.json` uses `scripts/mcp-serve.sh`, which loads `.env` and runs the local binary:

```json
{
  "mcpServers": {
    "openlane": {
      "command": "bash",
      "args": ["${workspaceFolder}/scripts/mcp-serve.sh"]
    }
  }
}
```

Or install globally and pass env vars directly:

```json
{
  "mcpServers": {
    "openlane": {
      "command": "openlane-mcp",
      "args": ["serve"],
      "env": {
        "OPENLANE_API_TOKEN": "${env:OPENLANE_API_TOKEN}",
        "OPENLANE_ORGANIZATION_ID": "${env:OPENLANE_ORGANIZATION_ID}"
      }
    }
  }
}
```

### Cursor (HTTP)

Start the server with `bash scripts/mcp-http.sh` (see [HTTP transport](#http-transport)), then use `examples/cursor-http.mcp.json` as a template.

### Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "openlane": {
      "command": "openlane-mcp",
      "args": ["serve"],
      "env": {
        "OPENLANE_API_TOKEN": "tola_your_token_here"
      }
    }
  }
}
```

Prefer environment substitution or a local secrets store over committing tokens. Examples in this repository use fictional values only.

## Security

Read tools are always available. Write tools require `OPENLANE_ALLOW_WRITE=true` or `--allow-write`. Delete tools require `OPENLANE_ALLOW_DELETE=true` or `--allow-delete`.

A successful read still requires:

1. An Openlane token with the relevant `object:read` scope (or equivalent PAT permissions)
2. Openlane authorization for that object in the selected organization

Writes and deletes additionally require server opt-in (`OPENLANE_ALLOW_WRITE` / `OPENLANE_ALLOW_DELETE`) and matching Openlane token permissions. Workflow definition writes, workflow assignment actions, native policy lifecycle actions, and workflow deletes also require `confirm: true` on the tool call.

Tokens are never logged. See [docs/security.md](docs/security.md).

## Tool coverage

| Tool | Description |
| --- | --- |
| `openlane_controls_list` | List controls |
| `openlane_controls_search` | Search controls by ref code, title, or description |
| `openlane_control_get` | Get a control by ID (with relationship summaries) |
| `openlane_programs_list` | List programs (optional name filter) |
| `openlane_program_get` | Get a program by ID (with relationship summaries) |
| `openlane_evidence_list` | List evidence metadata (optional program/control filters) |
| `openlane_evidence_get` | Get evidence metadata by ID |
| `openlane_policies_list` | List internal policies (optional status filter) |
| `openlane_policies_awaiting_approval` | List policies awaiting approval (native NEEDS_APPROVAL + your pending workflow assignments) |
| `openlane_policy_get` | Get a policy by ID |
| `openlane_risks_list` | List risks (optional program/entity/control/status filters) |
| `openlane_risk_get` | Get a risk by ID (with relationship summaries) |
| `openlane_findings_list` | List findings (optional program/assessment/open/status/severity filters) |
| `openlane_finding_get` | Get a finding by ID |
| `openlane_assessments_list` | List assessments |
| `openlane_assessment_get` | Get an assessment by ID |
| `openlane_control_implementations_list` | List control implementations |
| `openlane_control_implementation_get` | Get a control implementation by ID |
| `openlane_standards_list` | List standards / frameworks |
| `openlane_standard_get` | Get a standard by ID |
| `openlane_tasks_list` | List tasks |
| `openlane_task_get` | Get a task by ID |
| `openlane_entities_list` | List entities (vendors; optional risk/tier/review/security filters) |
| `openlane_entity_get` | Get an entity by ID (vendor/security/commercial fields) |
| `openlane_assets_list` | List assets |
| `openlane_asset_get` | Get an asset by ID |
| `openlane_contacts_list` | List contacts |
| `openlane_contact_get` | Get a contact by ID |
| `openlane_groups_list` | List groups (optional name filter) |
| `openlane_group_get` | Get a group by ID |
| `openlane_users_list` | List users (optional name/email filters) |
| `openlane_user_get` | Get a user by ID |
| `openlane_workflows_list` | List workflow definitions (optional schema/kind/active filters) |
| `openlane_workflows_search` | Search workflow definitions by name or description |
| `openlane_workflow_get` | Get a workflow definition by ID (with plain-English summary) |
| `openlane_workflow_instances_list` | List workflow instances (optional definition/state/object filters) |
| `openlane_workflow_instance_get` | Get a workflow instance by ID (assignments, events, proposal preview) |
| `openlane_workflow_assignments_list` | List my workflow approval assignments |
| `openlane_workflow_assignment_get` | Get a workflow assignment by ID (targets, due date, object context) |
| `openlane_workflow_metadata_get` | Get workflow-eligible fields, edges, and resolver keys per object type |

Write tools (require `OPENLANE_ALLOW_WRITE=true` or `--allow-write`):

| Tool | Description |
| --- | --- |
| `openlane_control_create` / `openlane_control_update` | Create or update a control |
| `openlane_evidence_create` / `openlane_evidence_update` | Create or update evidence; optional base64 file uploads |
| `openlane_policy_create` / `openlane_policy_update` | Create or update an internal policy |
| `openlane_policy_submit_for_approval` / `openlane_policy_approve` / `openlane_policy_publish` / `openlane_policy_return_to_draft` | Native InternalPolicy status transitions (`confirm` required) |
| `openlane_risk_create` / `openlane_risk_update` | Create or update a risk |
| `openlane_task_create` / `openlane_task_update` | Create or update a task |
| `openlane_workflow_create` / `openlane_workflow_update` | Create or update a WorkflowDefinition (`confirm` required) |
| `openlane_workflow_assignment_approve` / `openlane_workflow_assignment_reject` | Approve or reject a WorkflowAssignment (`confirm` required) |
| `openlane_workflow_assignment_request_changes` / `openlane_workflow_assignment_reassign` | Request changes or reassign an assignment (`confirm` required) |

Delete tools (require `OPENLANE_ALLOW_DELETE=true` or `--allow-delete`):

| Tool | Description |
| --- | --- |
| `openlane_control_delete` | Delete a control by ID |
| `openlane_evidence_delete` | Delete evidence by ID |
| `openlane_policy_delete` | Delete a policy by ID |
| `openlane_risk_delete` | Delete a risk by ID |
| `openlane_task_delete` | Delete a task by ID |
| `openlane_workflow_delete` | Delete a workflow definition by ID (`confirm` required) |

See [docs/tools.md](docs/tools.md) for full details. With all modes enabled there are **66 tools** (40 read, 20 write, 6 delete).

Enriched get tools return bounded relationship summaries (`count` + `items`) so agents can answer program, vendor, control, and finding questions without chaining dozens of shallow calls.

List responses are paginated (`items`, `next_cursor`, `has_more`, `total_count`). Default page size is 20; maximum is 50.

There is no dedicated control search GraphQL operation in the current Openlane Go client. `openlane_controls_search` uses official `ControlWhereInput` contains-filters.

File contents and presigned download URLs are not returned from read tools. Write tools accept optional base64-encoded `files[]` on evidence create/update (default max 10 MiB decoded per file).

## HTTP transport

By default the server uses stdio. For local HTTP testing:

```bash
export OPENLANE_MCP_TRANSPORT=http
export OPENLANE_MCP_HTTP_ADDR=127.0.0.1:8090
openlane-mcp serve
```

The default bind is **loopback only** (`127.0.0.1:8090`). Do not use `0.0.0.0` or bare `:port` addresses unless you understand the exposure.

**HTTP mode has no built-in authentication.** Do not expose it directly to the public internet. Anyone who can reach the endpoint can use the server's Openlane token. For remote or shared-network deployment, place a trusted authentication reverse proxy (or equivalent private network controls) in front of the server. See [docs/security.md](docs/security.md).

Increase `OPENLANE_MCP_HTTP_MAX_BODY_BYTES` when uploading large evidence files over HTTP.

## MCP Registry

Published to the [MCP Registry](https://registry.modelcontextprotocol.io/) as `io.github.GregDog/mcp-server-theopenlane` on each tagged release. The OCI package is `ghcr.io/gregdog/mcp-server-theopenlane`.

## Development

```bash
make test
make build
./bin/openlane-mcp version
```

See [docs/development.md](docs/development.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Licence

Apache License 2.0. See [LICENSE](LICENSE).

Openlane names are used only to describe compatibility. Do not copy Openlane logos or imply official endorsement.
