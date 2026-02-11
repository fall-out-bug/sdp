# Agent Handoff — F059 Observability Bridge Design

**Для агента:** Открой этот репозиторий в Cursor. Workspace path: `/Users/fall_out_bug/projects/vibe_coding/sdp-F059`

---

## Контекст

- **Worktree:** `sdp-F059` (ветка `feature/F059`)
- **Feature:** F059 — Observability Bridge Design (дизайн-док, не имплементация)
- **Beads ID:** sdp-pom6

## Workstreams (в порядке выполнения)

| WS ID | Title | Depends |
|-------|-------|---------|
| 00-059-01 | Observability Bridge Design Document | 00-054-03 |
| 00-059-02 | OTel Semantic Convention Draft | 00-059-01 |

## Как работать

1. **Открой worktree в Cursor:**
   ```
   File → Open Folder → /Users/fall_out_bug/projects/vibe_coding/sdp-F059
   ```

2. **Перед работой:**
   ```bash
   export BEADS_NO_DAEMON=1
   sdp guard activate 00-059-01
   bd update sdp-pom6 --status in_progress
   ```

3. **Выполнение:**
   ```bash
   @build 00-059-01
   # после завершения:
   sdp guard activate 00-059-02
   @build 00-059-02
   ```

4. **После завершения:**
   ```bash
   bd sync
   bd close sdp-pom6
   git push -u origin feature/F059
   ```

## Важно

- **Все изменения делай в этом worktree** — не в основном репо `/Users/fall_out_bug/projects/vibe_coding/sdp`
- Workstream specs: `docs/workstreams/backlog/00-059-01.md`, `00-059-02.md`
- Summary: `docs/workstreams/backlog/F059-WORKSTREAMS-SUMMARY.md`
