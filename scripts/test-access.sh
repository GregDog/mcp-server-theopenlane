#!/usr/bin/env bash
# Quick live API check using .env credentials. Does not print tokens.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ ! -f .env ]]; then
	echo "error: .env not found (copy .env.example and add your token)" >&2
	exit 1
fi

if [[ ! -x bin/openlane-mcp ]]; then
	echo "building bin/openlane-mcp..."
	make build
fi

set -a
# shellcheck disable=SC1091
source .env
set +a

if [[ -z "${OPENLANE_API_TOKEN:-}" ]]; then
	echo "error: OPENLANE_API_TOKEN is empty in .env" >&2
	exit 1
fi

echo "openlane-mcp: $(./bin/openlane-mcp version)"
echo "base_url: ${OPENLANE_BASE_URL:-https://api.theopenlane.io}"
echo "org_id: ${OPENLANE_ORGANIZATION_ID:-<not set>}"
echo "token: set (prefix $(printf '%.4s' "$OPENLANE_API_TOKEN")..., length ${#OPENLANE_API_TOKEN})"
echo
echo "calling GetControls (limit 1)..."

go run ./scripts/test-access/main.go
