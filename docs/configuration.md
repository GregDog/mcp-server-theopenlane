# Configuration

Configuration is environment-based.

| Variable | Required | Default | Purpose |
| --- | --- | --- | --- |
| `OPENLANE_API_TOKEN` | yes | | Openlane API token (`tola_`) or PAT (`tolp_`) |
| `OPENLANE_BASE_URL` | no | `https://api.theopenlane.io` | API base URL for Cloud or self-hosted Openlane |
| `OPENLANE_ORGANIZATION_ID` | no | | Organization ID for multi-org PATs |
| `OPENLANE_MCP_LOG_LEVEL` | no | `info` | `debug`, `info`, `warn`, or `error` |
| `OPENLANE_ALLOW_WRITE` | no | `false` | Enable write MCP tools when `true`, `1`, `yes`, or `on` |
| `OPENLANE_ALLOW_DELETE` | no | `false` | Enable delete MCP tools when `true`, `1`, `yes`, or `on` |
| `OPENLANE_MCP_TRANSPORT` | no | `stdio` | `stdio` or `http` |
| `OPENLANE_MCP_HTTP_ADDR` | no | `127.0.0.1:8090` | Loopback listen address when `OPENLANE_MCP_TRANSPORT=http` |
| `OPENLANE_MCP_HTTP_JSON` | no | `false` | Use `application/json` responses for HTTP transport |
| `OPENLANE_MCP_HTTP_MAX_BODY_BYTES` | no | `33554432` | Max HTTP request body size (32 MiB) |
| `OPENLANE_MCP_MAX_UPLOAD_BYTES` | no | `10485760` | Max decoded evidence upload per file (10 MiB) |
| `OPENLANE_MCP_UPLOAD_TIMEOUT` | no | `2m` | Openlane API timeout for evidence uploads |

You can also pass `--allow-write`, `--allow-delete`, `--transport`, `--http-addr`, or `--http-json` to `openlane-mcp serve`. Write and delete are independent; enabling writes does not enable deletes.

Create tokens in the Openlane console under developer settings. The token value is shown once.

GraphQL is served at `{OPENLANE_BASE_URL}/query`. The Go client appends that path; set `OPENLANE_BASE_URL` to the API origin only, for example `https://api.theopenlane.io` or `https://openlane.example.internal`.

When `OPENLANE_ORGANIZATION_ID` is set, the server passes it to `openlane.WithOrganizationHeader`. That is required for personal access tokens that can access more than one organization.

The API token is never written to logs.
