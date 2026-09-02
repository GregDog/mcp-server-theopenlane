# Publishing

Tagged releases publish:

1. Go binaries and checksums (GoReleaser)
2. Docker images to `ghcr.io/gregdog/mcp-server-theopenlane`
3. MCP Registry entry via `mcp-publisher` in `.github/workflows/release.yml`

## MCP Registry

Manifest: [`server.json`](../server.json)

- **Name:** `io.github.GregDog/mcp-server-theopenlane`
- **Package:** OCI image `ghcr.io/gregdog/mcp-server-theopenlane:v<version>` (GoReleaser uses the git tag, including the `v` prefix)
- **Ownership:** Docker label `io.modelcontextprotocol.server.name` in `Dockerfile.release`

The release workflow syncs `server.json` version and image tag from the git tag, then runs:

```bash
mcp-publisher login github-oidc
mcp-publisher publish
```

Manual publish (maintainers):

```bash
# After updating server.json version and ensuring the OCI image exists
mcp-publisher login github
mcp-publisher publish
```

See the [MCP Registry publishing guide](https://github.com/modelcontextprotocol/registry/blob/main/docs/modelcontextprotocol-io/quickstart.mdx).
