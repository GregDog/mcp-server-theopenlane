# Architecture

```text
MCP Client (stdio or HTTP)
    → Openlane MCP Server
        → Official Openlane Go Client
            → Openlane API (GraphQL)
```

`openlane-mcp serve` speaks MCP over **stdio by default** using the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk). Set `OPENLANE_MCP_TRANSPORT=http` (or `--transport http`) for **Streamable HTTP** instead. Logs go to stderr so they do not corrupt the stdio protocol on stdout.

The server constructs an [official Openlane Go client](https://github.com/theopenlane/go-client) at startup. It does not call Openlane until a tool runs. Authentication, connectivity, and organisation-context errors are returned from that tool call.

```text
cmd/openlane-mcp        CLI: serve, version; stdio or HTTP transport
internal/config         Environment configuration
internal/openlane       Client wrapper, uploads, pagination, redaction
internal/tools          MCP tool handlers (read / write / delete)
```

Evidence file uploads decode base64 MCP payloads and pass `graphql.Upload` values to `CreateEvidence` / `UpdateEvidence`. There is no separate upload API.

There is no built-in MCP authentication in HTTP mode. The Openlane API token is process-wide configuration. Remote deployment requires a trusted external authentication reverse-proxy layer.

## Security model

```text
read tools are always available
        AND
write/delete tools only when explicitly enabled
        AND
Openlane token permissions
        AND
Openlane authorization
```

This server does not bypass Openlane privacy rules or FGA checks. A token that cannot read a control will not see it here either. Write and delete operations require both server opt-in and matching Openlane permissions.

## Distribution

- **GitHub Releases** — binaries, checksums, SBOMs
- **GHCR** — `ghcr.io/gregdog/mcp-server-theopenlane` Docker images
- **MCP Registry** — `io.github.GregDog/mcp-server-theopenlane` (see [publishing.md](publishing.md))
