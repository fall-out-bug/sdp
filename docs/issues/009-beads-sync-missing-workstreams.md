# 009: Workstreams not synced to Beads

**Source:** /issue (2026-01-31)  
**Status:** Fixed  
**Priority:** P3 (LOW)  
**Route:** /bugfix (applied)

## Problem

Workstream files exist in `docs/workstreams/` but are not tracked in Beads.

**Example:**
```bash
# File exists
ls docs/workstreams/completed/00-020-03-ws-complete-tests.md
# ✓ exists

# But Beads doesn't know about it
bd list | grep "020-03"
# (empty)

sdp guard activate 00-020-03
# ❌ WS not found: 00-020-03
```

## Root Cause

Workstreams created manually (markdown files) without running:
- `bd create` — to create Beads task first
- `bd sync` — to sync markdown → Beads after creation

The Beads-first workflow wasn't followed:

```
❌ Wrong:  Create .md file → (forget bd sync)
✅ Right:  bd create → work → bd sync → commit
```

## Impact

- `sdp guard activate <ws_id>` fails for unsynced WS
- Beads dashboard incomplete
- Traceability broken (can't track WS status via `bd`)

## Affected Workstreams

Need to audit which WS files exist but aren't in Beads:

```bash
# List all WS files
find docs/workstreams -name "*.md" -type f | grep -E "[0-9]{2}-[0-9]{3}-[0-9]{2}"

# Compare with bd list
bd list | grep -oE "[0-9]{2}-[0-9]{3}-[0-9]{2}"
```

## Action

**Immediate:** Run migration for existing WS:
```bash
sdp beads migrate docs/workstreams/backlog/ --real
sdp beads migrate docs/workstreams/completed/ --real
bd sync
```

**Process:** Update workflow documentation to enforce Beads-first:
1. `@design` should use `bd create` not just write .md
2. Pre-commit hook to warn if WS file added without Beads task
3. Add to CLAUDE.md / PROTOCOL.md

## Severity

- **P3 LOW** — workaround exists (manual `bd sync`)
- Doesn't block development
- Process/tooling gap, not a bug

## Resolution (2026-01-31)

- **Migration:** Ran `sdp beads migrate docs/workstreams/backlog/ --real` and `docs/workstreams/completed/ --real` (71+5 success)
- **Guard fix:** Added `resolve_ws_id_to_beads_id()` — guard activate now accepts ws_id (00-020-03) and resolves via `.beads-sdp-mapping.jsonl`
- **Docs:** Updated `docs/runbooks/beads-migration.md` with Beads-first workflow; added note to CLAUDE.md
- Branch: `bugfix/009-beads-sync-missing-workstreams`

## Related

- Issue 007: guard activate NameError (separate bug, fixed)
- CLAUDE.md: Beads workflow documentation
