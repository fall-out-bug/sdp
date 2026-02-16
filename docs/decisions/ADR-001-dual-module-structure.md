# ADR-001: Dual Go Module Structure

**Status:** Accepted
**Date:** 2026-02-13
**Workstream:** WS-067-10

## Context

SDP has two Go module trees:

| Location | Module Path | Purpose |
|----------|-------------|---------|
| `go.mod` (root) | `github.com/fall-out-bug/sdp` | Core modules (src/sdp/) |
| `sdp-plugin/go.mod` | `github.com/fall-out-bug/sdp` | CLI implementation |

Both modules use the **same module path**, which prevents using Go workspaces (`go.work`).

## Decision

**Keep dual modules without workspace** for now. The consolidation will be deferred to a future task.

### Rationale

1. **Name conflict**: Both modules use `github.com/fall-out-bug/sdp` as their module path. Go workspaces require unique module paths.

2. **Import impact**: Renaming `sdp-plugin` to `github.com/fall-out-bug/sdp/sdp-plugin` would require updating 152 import statements across 50+ files.

3. **Risk**: Large-scale import changes risk breaking builds and require extensive testing.

4. **Current state works**: Both modules build and test independently.

### Future Path

When consolidation becomes necessary:

1. **Option A (Recommended)**: Rename sdp-plugin module path
   ```bash
   # Update sdp-plugin/go.mod
   module github.com/fall-out-bug/sdp/sdp-plugin
   
   # Bulk update imports
   find sdp-plugin -name "*.go" -exec sed -i '' 's|github.com/fall-out-bug/sdp/internal|github.com/fall-out-bug/sdp/sdp-plugin/internal|g' {} \;
   ```

2. **Option B**: Merge src/sdp/ into sdp-plugin/internal/
   - Larger change, but creates single module
   - All code in one place

3. **Option C**: Deprecate root module
   - If src/sdp/ is not actively used, move to sdp-plugin or remove

## Consequences

- Developers must build/test each module separately
- No unified `go test ./...` from root
- IDEs may need separate configuration per module
- Dependabot configured for both modules (AC6-AC7 of WS-067-12)

## References

- [CONTRIBUTING.md](../../CONTRIBUTING.md) - Module structure documentation
- [DEVELOPMENT.md](../../DEVELOPMENT.md) - Build commands per module
