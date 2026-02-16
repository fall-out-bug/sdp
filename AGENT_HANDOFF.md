# Agent Handoff — F059 Observability Bridge Design

**For agent:** Open this repository in Cursor. Workspace path: `/Users/fall_out_bug/projects/vibe_coding/sdp-F059`

---

## Context

- **Worktree:** `sdp-F059` (branch `feature/F059`)
- **Feature:** F059 — Observability Bridge Design (design document, not implementation)
- **Beads ID:** sdp-pom6

## Workstreams (in execution order)

| WS ID | Title | Depends |
|-------|-------|---------|
| 00-059-01 | Observability Bridge Design Document | 00-054-03 |
| 00-059-02 | OTel Semantic Convention Draft | 00-059-01 |

## How to Work

1. **Open worktree in Cursor:**
   ```
   File → Open Folder → /Users/fall_out_bug/projects/vibe_coding/sdp-F059
   ```

2. **Before starting:**
   ```bash
   export BEADS_NO_DAEMON=1
   sdp guard activate 00-059-01
   bd update sdp-pom6 --status in_progress
   ```

3. **Execution:**
   ```bash
   @build 00-059-01
   # after completion:
   sdp guard activate 00-059-02
   @build 00-059-02
   ```

4. **After completion:**
   ```bash
   bd sync
   bd close sdp-pom6
   git push -u origin feature/F059
   ```

## Important

- **Make all changes in this worktree** — not in the main repo `/Users/fall_out_bug/projects/vibe_coding/sdp`
- Workstream specs: `docs/workstreams/backlog/00-059-01.md`, `00-059-02.md`
- Summary: `docs/workstreams/backlog/F059-WORKSTREAMS-SUMMARY.md`
