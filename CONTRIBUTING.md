# Contributing

Thanks for considering a contribution.

This is a small MCP server. Prefer a narrow change that is easy to review.

## Before you start

- Read [docs/architecture.md](docs/architecture.md) and [docs/development.md](docs/development.md).
- Do not invent Openlane API behaviour. Verify against the current official Go client.
- Do not add write or delete tools unless that work is explicitly in scope.
- Do not commit secrets.

## Development

```bash
make test
make fmt
make vet
```

## Pull requests

- Keep the diff small.
- Include tests for config, redaction, pagination, or handler changes.
- Update `docs/tools.md` if you add or skip a tool.
- Use a clear commit message.

Commit signing is optional for contributors. Maintainers should use GitHub-verified commits (SSH signing is fine).

## Code of conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
