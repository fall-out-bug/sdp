# Beads Migration Runbook

**Purpose:** Migrate SDP workstreams (markdown) to Beads tasks for 100% integration.

## Prerequisites

- Beads CLI installed: `bd --version`
- Beads initialized: `bd init` (creates `.beads/`)
- Project in SDP format with workstreams in `docs/workstreams/`

## Migration Command

```bash
# Migrate backlog (work to be done)
poetry run sdp beads migrate docs/workstreams/backlog/ --real

# Migrate completed (historical)
poetry run sdp beads migrate docs/workstreams/completed/ --real
```

**Note:** Use `--real` for real Beads CLI. Without it, uses mock (no actual tasks created).

## Workstream Requirements

For migration to succeed, workstream must have YAML frontmatter with:

- `ws_id` — PP-FFF-SS format (e.g., 00-032-01)
- `feature` — Feature ID (e.g., F032)
- `status` — backlog | active | completed | blocked
- `size` — SMALL | MEDIUM | LARGE

**Excluded from migration:**
- `00-032-00-*` — Feature overviews (no ws_id)
- `BEADS-001-*` — Epics (different format)
- Files without frontmatter (P0/P1 runbooks)
- Legacy WS with `id:` instead of `ws_id:`

## Output

- **Mapping file:** `.beads-sdp-mapping.jsonl` (ws_id ↔ beads_id, deduplicated on each run)
- **Beads tasks:** Created in `.beads/issues.jsonl`

**Note:** Migration calls `persist_mapping()` at the end to deduplicate the mapping file (fixes legacy append-duplicates).

## Verification

```bash
# Count migrated workstreams
wc -l .beads-sdp-mapping.jsonl

# List Beads tasks
bd list

# Check ready tasks
bd ready

# SDP status
poetry run sdp beads status
```

## Troubleshooting

**"title must be 500 characters or less"** — Fixed in sync: title truncated to 500 chars.

**"Missing required fields"** — Add ws_id, feature, status, size to frontmatter.

**"No frontmatter found"** — File uses blockquote format; add YAML frontmatter or exclude.

**Dependencies not linked** — Migrate in dependency order, or run migration twice (second run updates deps).
