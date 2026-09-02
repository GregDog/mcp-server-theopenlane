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
| `make test-access` | live API smoke test (requires `.env`) |
| `make vuln` | `govulncheck ./...` |

Do not commit `.env` files, `.local/` test artifacts, or real tokens.

## Local testing (stdio)

1. Copy `.env.example` to `.env` and set `OPENLANE_API_TOKEN` (and `OPENLANE_ORGANIZATION_ID` for multi-org PATs).
2. Build the binary: `make build`
3. Verify API access: `make test-access`
4. Open this project in Cursor — `.cursor/mcp.json` runs `scripts/mcp-serve.sh`, which sources `.env` and starts `bin/openlane-mcp serve` over stdio.

To test write or delete tools, also set in `.env`:

```bash
OPENLANE_ALLOW_WRITE=true   # create/update tools
OPENLANE_ALLOW_DELETE=true  # delete tools (independent of write)
```

Then rebuild and reload Cursor (**Developer: Reload Window**). Check **Output → MCP Logs** for `allow_write=true` and `allow_delete=true` on startup.

In Cursor: **Settings → MCP** and confirm `openlane` is enabled.

Example prompts:

- "List the first 5 Openlane controls"
- "List entities (vendors) with tier HIGH"
- "Get control OL-12.06"
- "Create evidence named MCP test with a text file" (requires write mode; see evidence uploads below)
- "Delete evidence `<id>`" (requires delete mode)

## HTTP transport testing

For Streamable HTTP instead of stdio:

1. Set in `.env`:

```bash
OPENLANE_MCP_TRANSPORT=http
OPENLANE_MCP_HTTP_ADDR=127.0.0.1:8090   # loopback only; do not use 0.0.0.0 without a trusted auth proxy
```

2. Start the server in a separate terminal:

```bash
bash scripts/mcp-http.sh
```

3. Point Cursor at the URL — see `examples/cursor-http.mcp.json`:

```json
{
  "mcpServers": {
    "openlane": {
      "url": "http://127.0.0.1:8090"
    }
  }
}
```

HTTP mode has no built-in MCP authentication. Do not expose it directly to the public internet. Use a trusted reverse proxy for remote deployment. See [security.md](security.md).

## Evidence upload testing

Evidence file uploads are tested via MCP write tools (`openlane_evidence_create` / `openlane_evidence_update`) with a `files[]` array:

```json
{
  "filename": "sample.txt",
  "content_type": "text/plain",
  "content_base64": "..."
}
```

Local test artifacts can live under `.local/` (gitignored). Do not commit test files or API tokens.

MCP read tools return `file_ids` but not presigned download URLs. Verify downloads in the Openlane UI or with a direct API client.

## Tool inventory

| Mode | Count |
| --- | --- |
| Read (always on) | 21 |
| Write (opt-in) | 10 |
| Delete (opt-in) | 5 |
| **Total** | **37** |

See [tools.md](tools.md) for the full mapping to Openlane client methods.

## Maintainer commits

GitHub shows commits as Verified when they are SSH-signed (or GPG/S/MIME). Contributors are not required to sign commits.

Registry publishing is documented in [publishing.md](publishing.md).
