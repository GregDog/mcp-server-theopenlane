# Contributing

Thanks for considering a contribution.

This is a small MCP server. Prefer a narrow change that is easy to review.

## Governance

The repository is maintained by [@GregDog](https://github.com/GregDog). Pull requests are reviewed on a best-effort basis — typically within a few business days. Larger changes should start with an issue so we can agree on approach before you invest time.

## How to contribute

1. **Check existing work** — search [open issues](https://github.com/GregDog/mcp-server-theopenlane/issues), [Discussions](https://github.com/GregDog/mcp-server-theopenlane/discussions), and open PRs. Issues labeled `good first issue` or `help wanted` are good entry points.
2. **Fork** the repository on GitHub.
3. **Clone** your fork and create a branch from `main`:
   ```bash
   git clone https://github.com/<you>/mcp-server-theopenlane.git
   cd mcp-server-theopenlane
   git checkout -b my-feature
   ```
4. **Install Go 1.27** — required by `go.mod`. See [go.dev/dl](https://go.dev/dl/).
5. **Develop and test** locally (see [Development](#development) below).
6. **Push** to your fork and **open a pull request** against `main` on `GregDog/mcp-server-theopenlane`.
7. Ensure **CI passes** on your PR (`gofmt`, `go vet`, `go test`, `go build`, `govulncheck`).

For new MCP tools, follow [docs/adding-tools.md](docs/adding-tools.md).

## Before you start

- Read [docs/architecture.md](docs/architecture.md) and [docs/development.md](docs/development.md).
- Do not invent Openlane API behaviour. Verify against the current [official Go client](https://github.com/theopenlane/go-client).
- Write and delete tools must stay opt-in. Do not enable them by default or add new mutating tools without updating `docs/tools.md` and `docs/security.md`.
- Do not commit secrets, `.env`, or `.local/` test artifacts.

## Development

Requires **Go 1.27** or later.

```bash
make test   # unit tests — no Openlane token required
make fmt
make vet
make check  # fmt + vet + test
make build  # writes bin/openlane-mcp
```

### Testing without an Openlane account

You do **not** need an Openlane API token to contribute. All unit tests use fakes and stubs under `internal/tools/*_test.go` and run via `make test`.

Optional live smoke test (only if you have a sandbox org and token):

```bash
cp .env.example .env
# set OPENLANE_API_TOKEN and optionally OPENLANE_ORGANIZATION_ID
make test-access
```

See [docs/development.md](docs/development.md) for Cursor/MCP local testing.

## Pull requests

- Keep the diff small.
- Include tests for config, redaction, pagination, or handler changes.
- Update `docs/tools.md` if you add or change a tool.
- Update [docs/publishing.md](docs/publishing.md) if you change `server.json` or the release workflow.
- Fill in the PR template checklist.
- Use a clear commit message.

Commit signing is optional for contributors. Maintainers should use GitHub-verified commits (SSH signing is fine).

## Code of conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
