# SDP — Codex setup

This project uses [Spec-Driven Protocol (SDP)](https://github.com/fall-out-bug/sdp). Codex reads this file for setup.

## Project-level skills

Project skills source of truth lives in `prompts/skills/` (this repo). Tool folders (`.codex`, `.claude`, `.cursor`, `.opencode`) use symlinks to this source.

## Quick start

1. Install SDP into your project with Codex integration:
   ```bash
   SDP_IDE=codex curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
   ```
2. Run project init:
   ```bash
   sdp init --auto
   ```
3. Use `@build 00-XXX-YY` or `sdp plan`, `sdp apply`, `sdp log trace` per [CLAUDE.md](../CLAUDE.md).

If you want the CLI only, use:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
```

## Directory layout

```
.codex/
├── INSTALL.md   # This file (read by Codex)
├── agents/      # Project-level agent symlink
└── skills/
    ├── README.md
    └── sdp/     # Project-level skills sourced from prompts/skills

~/.codex/
└── skills/      # User-level skills (persistent)
```

## Beads (optional)

If Beads is installed (`bd --version`), use `bd ready`, `bd update`, `bd close` for task tracking. See [AGENTS.md](../AGENTS.md).

## Updates

Rerun the same installer command to refresh the vendored `sdp/` checkout and managed Codex links after upstream changes.
