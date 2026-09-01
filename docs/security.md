# Security

v0.1 is read-only. There are no write or delete MCP tools, and no `--allow-write` / `--allow-delete` flags.

A call succeeds only if:

```text
the MCP tool exists and is read-only
        AND
the Openlane token is valid and permitted for that object
        AND
Openlane authorization allows the read
```

Organization API tokens are scoped as `object:action`, for example `control:read`. PATs inherit the creating user's permissions. Grant the minimum scopes needed.

## Token handling

- Tokens are read from `OPENLANE_API_TOKEN`.
- Values matching `tola_` or `tolp_` prefixes are redacted from error strings.
- Tokens are not used as slog attributes.
- Do not put real tokens in git, issue reports, or example files.

## Request bounds

- HTTP timeout: 30 seconds
- Default page size: 20
- Maximum page size: 50
- No automatic retries

## Reporting vulnerabilities

See [SECURITY.md](../SECURITY.md). Enable GitHub private vulnerability reporting and secret scanning on the repository.
