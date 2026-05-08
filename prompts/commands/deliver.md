# /deliver — Autonomous Feature Delivery

Invoke `@delivery-loop` with no arguments.

The skill handles end-to-end:

1. Feature selection (`bd ready -n 50`, pick highest-priority epic/feature).
2. Workstream identification (cross-reference `docs/workstreams/backlog/` with beads children).
3. Claim + worktree + checkpoint bootstrap.
4. Build → review → fix loop (bounded).
5. PR creation (after local quality gates pass).
6. Codex review loop (bounded, stable-N exit).
7. Closeout (bead close, worktree teardown, beads transport push).

## Recovery

- **Resume after compaction:** `@delivery-loop --resume`
- **Abort mid-loop:** `@delivery-loop --abort`
  (cleans claim, worktree, checkpoint, and lock; stashes uncommitted work)

## Escalation policy

Do **not** stop for routine fix/rebuild decisions. **Do** stop to escalate:
- Tests fail unrelated to feature code
- Merge conflicts
- Ambiguous findings with no clear fix strategy
- Phase-1 cap hit at cycle 5 (operator must paste deferred-P3 list into spin-out bead)

See `.agents/skills/delivery-loop.md` for the full state machine and `docs/plans/2026-04-22-deliver-skill-review-design.md` for the design rationale.
