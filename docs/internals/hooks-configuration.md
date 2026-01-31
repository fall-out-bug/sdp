# Hook Configuration

Hooks auto-detect project configuration for project-agnostic operation.

## Project Root Detection

Hooks auto-detect project root using these markers (in order of precedence):

1. **`.sdp-root`** — Explicit marker file at project root
2. **`docs/workstreams/`** — Standard workstream directory
3. **`.git` + `pyproject.toml`** — Git root with `[tool.sdp]` section

## Workstream Directory

Search order for workstream directory:

1. **`quality-gate.toml`** — `[workstreams] dir = "path"`
2. **`SDP_WORKSTREAM_DIR`** — Environment variable (absolute or relative)
3. **`docs/workstreams/`** — Default
4. **`workstreams/`** — Legacy fallback

## Configuration Example

```toml
# quality-gate.toml
[workstreams]
dir = "docs/workstreams"  # Relative to project root

[quality]
coverage_min = 80
file_size_max = 200
```

## Environment Variables

| Variable | Purpose |
|---------|---------|
| `SDP_WORKSTREAM_DIR` | Override workstream directory path |
| `SDP_HARD_PUSH` | Set to `1` to block push on test failures |

## Backward Compatibility

- **hw_checker layout**: `tools/hw_checker/docs/workstreams` still supported as fallback
- **SDP layout**: `docs/workstreams` is the default
- No breaking changes for existing projects
