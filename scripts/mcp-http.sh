#!/usr/bin/env bash
# Start the MCP server in Streamable HTTP mode (sources .env; requires OPENLANE_MCP_TRANSPORT=http).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="$ROOT/.env"
BIN="$ROOT/bin/openlane-mcp"

if [[ ! -f "$ENV_FILE" ]]; then
	echo "error: $ENV_FILE not found (copy .env.example and add your token)" >&2
	exit 1
fi

if [[ ! -x "$BIN" ]]; then
	echo "error: $BIN not found; run: make build" >&2
	exit 1
fi

set -a
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a

export OPENLANE_MCP_TRANSPORT=http
export OPENLANE_MCP_HTTP_ADDR="${OPENLANE_MCP_HTTP_ADDR:-127.0.0.1:8090}"

if [[ -z "${OPENLANE_API_TOKEN:-}" ]]; then
	echo "error: OPENLANE_API_TOKEN is empty in $ENV_FILE" >&2
	exit 1
fi

echo "starting openlane-mcp HTTP transport on ${OPENLANE_MCP_HTTP_ADDR}" >&2
exec "$BIN" serve
