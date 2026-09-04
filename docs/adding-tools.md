# Adding a tool

This guide covers the usual path for adding or extending an Openlane MCP tool. Read [architecture.md](architecture.md) first.

## Rules

- Verify behaviour against the [official Openlane Go client](https://github.com/theopenlane/go-client). Do not invent GraphQL fields or mutations.
- Read tools are always registered. Write and delete tools must stay behind `OPENLANE_ALLOW_WRITE` / `OPENLANE_ALLOW_DELETE` (see `internal/tools/register.go`).
- New write or delete tools require updates to [tools.md](tools.md), [security.md](security.md), and the tool tables in [README.md](../README.md).
- Destructive or workflow actions that change production data should require `confirm: true` on the tool input, matching existing write tools.

## Layout

| Path | Role |
| --- | --- |
| `internal/tools/<domain>.go` | Read tools (list, get, search) |
| `internal/tools/writes_<domain>.go` | Create/update tools |
| `internal/tools/deletes.go` | Delete tools |
| `internal/tools/register.go` | Registration and write/delete gating |
| `internal/tools/addtool.go` | Typed `addTool` helper |
| `internal/openlane/client.go` | `GraphAPI` interface — add methods here when the go-client exposes new calls |
| `internal/tools/*_test.go` | Handler tests with `fakeAPI` stubs |

## Read tool checklist

1. **Extend `GraphAPI`** in `internal/openlane/client.go` if the go-client method is not already wrapped.
2. **Add handler(s)** in `internal/tools/<domain>.go`:
   - Define a JSON output struct (e.g. `controlItem`).
   - Use `addTool` with `readOnly()` annotations for list/get/search tools.
   - Reuse `listInput` / pagination via `openlane.ClampLimit` and `openlane.CursorPtr` from `internal/tools/common.go`.
3. **Register** in `register<Domain>(server, h)` and call it from `Register` in `register.go`.
4. **Test** with a `fakeAPI` stub in `internal/tools/<domain>_test.go` or an existing `fake_api_*_test.go` file. Tests must pass without a live Openlane token.
5. **Document** the tool name, parameters, and behaviour in [tools.md](tools.md) and the README tool table.

## Write / delete tool checklist

1. Follow the read-tool steps, but register inside the `if opts.AllowWrite` or `if opts.AllowDelete` block in `register.go`.
2. Add input validation and map go-client responses to stable JSON output structs.
3. For file uploads, follow `internal/tools/writes_evidence.go` (base64 decode, size limits from `handlers.maxUploadBytes`).
4. Update [security.md](security.md) with any new risk or confirmation requirement.
5. Add tests in `fake_api_writes_test.go` or `fake_api_writes_test.go` / `deletes_test.go` as appropriate.

## Example pattern

Registration (from `controls.go`):

```go
func registerControls(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_controls_list",
		Title:       "List Openlane controls",
		Description: "List controls in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listControls)
}
```

Handler:

```go
func (h *handlers) listControls(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[controlItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetControls(ctx, first, after, where)
	// map resp → openlane.Page[controlItem]
}
```

## Testing without Openlane credentials

```bash
make test    # unit tests only — no API token required
make check   # fmt + vet + test
```

Optional live smoke test (maintainers or contributors with a sandbox org):

```bash
cp .env.example .env   # set OPENLANE_API_TOKEN
make test-access
```

## Pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for the full workflow. Tool PRs should include:

- [ ] `make check` passes
- [ ] Tests for new handler logic
- [ ] [tools.md](tools.md) updated
- [ ] README tool table updated (if adding a public tool)
- [ ] [security.md](security.md) updated (if adding write/delete)
