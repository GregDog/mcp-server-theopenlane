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

MCP logs from Cursor appear in the Output panel under MCP Logs. The server writes structured logs to stderr.

## Maintainer commits

GitHub shows commits as Verified when they are SSH-signed (or GPG/S/MIME). Contributors are not required to sign commits.
