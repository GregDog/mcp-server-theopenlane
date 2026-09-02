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

- Openlane MCP access for programs, controls, evidence, policies, risks, standards, tasks, entities (vendors), assets, and contacts
- Opt-in create/update tools for controls, evidence, policies, risks, and tasks
- Opt-in delete tools for the same domains (except programs and standards)
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

Writes and deletes additionally require server opt-in (`OPENLANE_ALLOW_WRITE` / `OPENLANE_ALLOW_DELETE`) and matching Openlane token permissions.

Tokens are never logged. See [docs/security.md](docs/security.md).

## Tool coverage

| Tool | Description |
| --- | --- |
| `openlane_controls_list` | List controls |
| `openlane_controls_search` | Search controls by ref code, title, or description |
| `openlane_control_get` | Get a control by ID |
| `openlane_programs_list` | List programs |
| `openlane_program_get` | Get a program by ID |
| `openlane_evidence_list` | List evidence metadata |
| `openlane_evidence_get` | Get evidence metadata by ID |
| `openlane_policies_list` | List internal policies |
| `openlane_policy_get` | Get a policy by ID |
| `openlane_risks_list` | List risks |
| `openlane_risk_get` | Get a risk by ID |
| `openlane_standards_list` | List standards / frameworks |
| `openlane_standard_get` | Get a standard by ID |
| `openlane_tasks_list` | List tasks |
| `openlane_task_get` | Get a task by ID |
| `openlane_entities_list` | List entities (vendors and third parties) |
| `openlane_entity_get` | Get an entity by ID |
| `openlane_assets_list` | List assets |
| `openlane_asset_get` | Get an asset by ID |
| `openlane_contacts_list` | List contacts |
| `openlane_contact_get` | Get a contact by ID |

Write tools (require `OPENLANE_ALLOW_WRITE=true` or `--allow-write`):

| Tool | Description |
| --- | --- |
| `openlane_control_create` / `openlane_control_update` | Create or update a control |
| `openlane_evidence_create` / `openlane_evidence_update` | Create or update evidence; optional base64 file uploads |
| `openlane_policy_create` / `openlane_policy_update` | Create or update an internal policy |
| `openlane_risk_create` / `openlane_risk_update` | Create or update a risk |
| `openlane_task_create` / `openlane_task_update` | Create or update a task |

Delete tools (require `OPENLANE_ALLOW_DELETE=true` or `--allow-delete`):

| Tool | Description |
| --- | --- |
| `openlane_control_delete` | Delete a control by ID |
| `openlane_evidence_delete` | Delete evidence by ID |
| `openlane_policy_delete` | Delete a policy by ID |
| `openlane_risk_delete` | Delete a risk by ID |
| `openlane_task_delete` | Delete a task by ID |

See [docs/tools.md](docs/tools.md) for full details. With all modes enabled there are **37 tools** (21 read, 10 write, 5 delete).

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
