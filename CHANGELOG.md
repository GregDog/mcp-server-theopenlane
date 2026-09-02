# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

## [0.4.0] - 2026-09-02

### Added

- Six new read tools: findings, assessments, and control implementations (`openlane_findings_*`, `openlane_assessments_*`, `openlane_control_implementations_*`)
- Enriched get responses for programs, controls, entities, risks, and evidence with compact relationship summaries
- List filters on entities, risks, findings, evidence, programs, assessments, and control implementations
- Optional `.env` auto-load when `OPENLANE_API_TOKEN` is unset (searches `OPENLANE_ENV_FILE`, cwd, and binary-relative paths)
- `scripts/mcp-verify.sh` to confirm MCP tool registration after `make build`

### Fixed

- Omit per-tool `outputSchema` in `tools/list` so Cursor does not silently drop oversized tool lists (connected with 0 tools)
- MCP Registry publish step in release workflow (`mcp-publisher` v1.8.1 download URL)

## [0.3.0] - 2026-09-02

### Added

- Read tools for tasks, entities (vendors), assets, and contacts (`openlane_tasks_*`, `openlane_entities_*`, `openlane_assets_*`, `openlane_contacts_*`)
- Evidence file uploads on `openlane_evidence_create` and `openlane_evidence_update` via base64 `files[]` payloads
- Streamable HTTP transport (`OPENLANE_MCP_TRANSPORT=http`, `--transport http`, `--http-addr`)
- MCP Registry manifest (`server.json`) and automated publishing from release workflow
- `scripts/mcp-http.sh` and `examples/cursor-http.mcp.json` for HTTP transport testing
- `docs/publishing.md` for maintainer registry steps

### Changed

- HTTP transport defaults to loopback-only bind (`127.0.0.1:8090`); server timeouts, request size limits, log redaction for auth headers and upload payloads; startup warning when binding to a non-loopback address

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
