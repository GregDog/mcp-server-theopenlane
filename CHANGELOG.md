# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]

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
