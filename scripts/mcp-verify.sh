#!/usr/bin/env bash
# Quick sanity check: MCP stdio server lists tools (expects make build + .env).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ ! -x "$ROOT/bin/openlane-mcp" ]]; then
	echo "error: run make build first" >&2
	exit 1
fi
if [[ ! -f "$ROOT/.env" ]]; then
	echo "error: .env not found (copy .env.example)" >&2
	exit 1
fi

go run ./scripts/mcp-verify/main.go
