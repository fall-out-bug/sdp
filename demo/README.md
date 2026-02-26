# SDP Demo

Terminal recording of the SDP workflow with OpenCode CLI.

## Clean environment (recommended)

```bash
cd sdp/demo
./run-demo.sh
```

**What it does:**
1. Creates a temp dir (e.g. `/var/folders/.../sdp-demo-XXXXXX`)
2. Clones sdp there
3. Runs `vhs demo/demo.tape` from the clone
4. Copies `demo.gif` to `demo/demo.gif`

**Env vars:** `SDP_REPO`, `SDP_REF` (default: main)

## Quick run (from existing clone)

```bash
cd /path/to/sdp    # sdp repo root (e.g. sdp_dev/sdp)
vhs demo/demo.tape
```

Creates `demo-beads-viz/` in the current dir (gitignored). Output: `demo/demo.gif`

## Requirements

- `vhs` — `brew install vhs`
- `opencode` — OpenCode CLI
- `go`, `curl`
- `bd` (beads) — for `@feature` (optional: `bd init` runs if available)

**Note:** If running from Cursor's integrated terminal, the sandbox may block `opencode` access to the temp dir. Run `./run-demo.sh` from a regular terminal, or approve directory access when prompted.

## What the demo shows

1. Install SDP (per README)
2. `sdp init --auto`
3. `opencode run '@feature beads viz in opencode'`
4. `opencode run '@oneshot F001' --continue`
5. `opencode run '@review F001' --continue`
6. `sdp verify 00-001-01`
7. Result (bd list or workstreams)
