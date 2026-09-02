# Openlane MCP Server

[![CI](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/ci.yml/badge.svg)](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![CodeQL](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/codeql.yml/badge.svg)](https://github.com/GregDog/mcp-server-theopenlane/actions/workflows/codeql.yml)

A secure, open-source [Model Context Protocol](https://modelcontextprotocol.io/) server for the [Openlane](https://www.theopenlane.io/) GRC platform.

This project is **not** an official Openlane product and is not endorsed by theopenlane, Inc.

## Overview

`openlane-mcp` lets MCP clients such as Cursor and Claude Desktop query Openlane over stdio. It talks to Openlane Cloud or a self-hosted instance through the [official Openlane Go client](https://github.com/theopenlane/go-client).

```text
MCP Client
    → Openlane MCP Server (stdio)
        → Official Openlane Go Client
            → Openlane API
```

The server is read-only by default. Write tools are opt-in. Openlane authorization still applies to every request.

## Features

- Openlane MCP access for programs, controls, evidence, policies, risks, and standards
- Openlane Cloud and self-hosted Openlane (configurable base URL)
- Read-only tools
- stdio transport
- Native Go binary
- Docker image (published with GitHub Releases)

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

### Cursor

Add to `~/.cursor/mcp.json` or `.cursor/mcp.json`:

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

List responses are paginated (`items`, `next_cursor`, `has_more`, `total_count`). Default page size is 20; maximum is 50.

There is no dedicated control search GraphQL operation in the current Openlane Go client. `openlane_controls_search` uses official `ControlWhereInput` contains-filters.

File contents and presigned download URLs are not returned.

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
