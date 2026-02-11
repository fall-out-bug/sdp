# SDP — Codex setup

This project uses [Spec-Driven Protocol (SDP)](https://github.com/fall-out-bug/sdp). Codex reads this file for setup.

## Project-level skills

Project skills source of truth lives in `prompts/skills/` (this repo). Tool folders (`.codex`, `.claude`, `.cursor`, `.opencode`) use symlinks to this source.

## Quick start

1. Install SDP CLI (from repo root):
   ```bash
   cd sdp-plugin && go build -o sdp ./cmd/sdp && mv sdp ../
   ```
2. Ensure `sdp` is on PATH.
3. Use `@build 00-XXX-YY` or `sdp plan`, `sdp apply`, `sdp log trace` per [CLAUDE.md](../CLAUDE.md).

## Directory layout

```
.codex/
├── INSTALL.md   # This file (read by Codex)
└── skills/      # Project-level symlinks to prompts/skills

~/.codex/
└── skills/      # User-level skills (persistent)
```

## Beads (optional)

If Beads is installed (`bd --version`), use `bd ready`, `bd update`, `bd close` for task tracking. See [AGENTS.md](../AGENTS.md).
