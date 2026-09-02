# Development

Requires Go 1.27.

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
| `make vuln` | `govulncheck ./...` |

Do not commit `.env` files or real tokens.

## Local testing

1. Copy `.env.example` to `.env` and set `OPENLANE_API_TOKEN` (and `OPENLANE_ORGANIZATION_ID` for multi-org PATs).
2. Build the binary: `make build`
3. Verify API access: `make test-access`
4. Open this project in Cursor — `.cursor/mcp.json` runs `scripts/mcp-serve.sh`, which sources `.env` and starts `bin/openlane-mcp serve`.

In Cursor: **Settings → MCP** (or the MCP panel) and confirm `openlane` is enabled. Logs appear under **Output → MCP Logs**.

Example prompts to try in chat:

- "List the first 5 Openlane controls"
- "Get control OL-12.06"
- "List Openlane programs"

MCP logs from Cursor appear in the Output panel under MCP Logs. The server writes structured logs to stderr.

## Maintainer commits

GitHub shows commits as Verified when they are SSH-signed (or GPG/S/MIME). Contributors are not required to sign commits.
