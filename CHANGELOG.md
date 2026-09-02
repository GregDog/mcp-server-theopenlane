# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

## [0.3.0] - 2026-09-02

### Added

- Read tools for tasks, entities (vendors), assets, and contacts (`openlane_tasks_*`, `openlane_entities_*`, `openlane_assets_*`, `openlane_contacts_*`)
- Evidence file uploads on `openlane_evidence_create` and `openlane_evidence_update` via base64 `files[]` payloads
- Streamable HTTP transport (`OPENLANE_MCP_TRANSPORT=http`, `--transport http`, `--http-addr`)
- MCP Registry manifest (`server.json`) and automated publishing from release workflow
- `scripts/mcp-http.sh` and `examples/cursor-http.mcp.json` for HTTP transport testing
- `docs/publishing.md` for maintainer registry steps

### Fixed

- `openlane-mcp serve` no longer overrides `OPENLANE_MCP_TRANSPORT` and `OPENLANE_MCP_HTTP_ADDR` from the environment when CLI flags are omitted

## [0.2.0] - 2026-09-02

### Added

- Opt-in write tools for controls, evidence metadata, policies, risks, and tasks
- Opt-in delete tools for controls, evidence, policies, risks, and tasks
- `OPENLANE_ALLOW_WRITE` / `--allow-write` and `OPENLANE_ALLOW_DELETE` / `--allow-delete` (both disabled by default)

## [0.1.0] - 2026-09-02

### Added

- Initial stdio MCP server (`openlane-mcp serve` / `openlane-mcp version`)
- Read-only tools for controls, programs, evidence, policies, risks, and standards
- Environment configuration, token redaction, and bounded pagination
- CI, CodeQL, Dependabot, GoReleaser, and Docker packaging
- Local testing helpers: `make test-access`, `scripts/mcp-serve.sh`, and project `.cursor/mcp.json`
