# Security

Read tools are always available. Write tools are opt-in via `OPENLANE_ALLOW_WRITE` or `--allow-write`. Delete tools are opt-in via `OPENLANE_ALLOW_DELETE` or `--allow-delete`. Enabling writes does not enable deletes.

A call succeeds only if:

```text
the MCP tool is enabled by server configuration (for writes/deletes)
        AND
the Openlane token is valid and permitted for that object
        AND
Openlane authorization allows the operation
```

Organization API tokens are scoped as `object:action`, for example `control:read`. PATs inherit the creating user's permissions. Grant the minimum scopes needed.

## Token handling

- Tokens are read from `OPENLANE_API_TOKEN`.
- Values matching `tola_` or `tolp_` prefixes are redacted from error strings.
- Tokens are not used as slog attributes.
- Do not put real tokens in git, issue reports, or example files.

## HTTP transport

The Streamable HTTP transport (`OPENLANE_MCP_TRANSPORT=http`) is intended for **local development and trusted networks only**.

- **Default bind:** `127.0.0.1:8090` (loopback). The server does not default to `0.0.0.0` or bare `:port` addresses, which would listen on all interfaces.
- **No built-in authentication:** HTTP mode does not validate MCP client identity. Anyone who can reach the endpoint can invoke tools using the server's configured Openlane token.
- **Do not expose directly to the public internet.** A publicly reachable, unauthenticated endpoint would grant Openlane API access to anonymous clients.
- **Remote or network deployment** requires a trusted external layer (reverse proxy, VPN, or private network) that performs authentication and access control before traffic reaches `openlane-mcp`.
- **Non-loopback binds** (for example `0.0.0.0:8090`) log a startup warning. Built-in HTTP authentication is not enabled.
- Request body size is capped (`OPENLANE_MCP_HTTP_MAX_BODY_BYTES`, default 32 MiB). Server read/write/idle timeouts apply.
- No debug or pprof endpoints are registered; only the MCP Streamable HTTP handler is served.
- Logs redact Openlane tokens, `Authorization` values, and `content_base64` upload payloads.

## Request bounds

- HTTP timeout: 30 seconds (2 minutes for evidence uploads)
- Default page size: 20
- Maximum page size: 50
- Default evidence upload size: 10 MiB decoded per file
- No automatic retries

## Reporting vulnerabilities

See [SECURITY.md](../SECURITY.md). Enable GitHub private vulnerability reporting and secret scanning on the repository.
