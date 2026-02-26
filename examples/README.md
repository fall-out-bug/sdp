# SDP Examples

Minimal projects for testing SDP commands.

## beads-viz-demo

Minimal Go project with one workstream. Use for testing `sdp verify`, `sdp status`, `sdp build`.

```bash
cd examples/beads-viz-demo
# Ensure sdp is in PATH (from sdp/sdp-plugin or install)
sdp verify 00-001-01
```

**Structure:**
- `main.go`, `main_test.go` — minimal Go code
- `docs/workstreams/backlog/00-001-01.md` — workstream with scope_files, verification_commands
