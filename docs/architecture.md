# Architecture

```text
MCP Client
    → Openlane MCP Server
        → Official Openlane Go Client
            → Openlane API
```

`openlane-mcp serve` speaks MCP over stdio using the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk). Logs go to stderr so they do not corrupt the protocol on stdout.

The server constructs an [official Openlane Go client](https://github.com/theopenlane/go-client) at startup. It does not call Openlane until a tool runs. Authentication, connectivity, and organisation-context errors are returned from that tool call.

```text
cmd/openlane-mcp        CLI: serve, version
internal/config         Environment configuration
internal/openlane       Client construction, pagination, redaction
internal/tools          MCP tool handlers
```

There is no Streamable HTTP transport in this release.

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
