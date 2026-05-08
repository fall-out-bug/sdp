---
name: ux
description: UX research with mental model elicitation and codebase pattern scan
version: 1.0.0
changes:
  - Initial release: 6-question listening session + autonomous codebase research
---

# @ux - UX Research

**Structured UX findings that @design consumes as acceptance criteria.** Not "what UI do you want?" — mental model elicitation.

---

## EXECUTE THIS NOW

When user invokes `@ux {feature-id}` or when `@feature` auto-triggers it (user-facing features only):

### Phase 1: Listening Session (6 Questions)

These are mental model elicitation questions, not UI specification:

1. **Context of reach:** "What is the user doing in the 10 minutes before they encounter this feature? What problem are they mid-solving?"

2. **Mental model gap (Don Norman: Gulf of Execution):** "What will the user *think* happens when they perform the primary action? Where does that model likely diverge from what the system actually does?"

3. **Workaround reality:** "What do users do today without this feature? The workaround reveals the existing mental model."

4. **Friction prediction:** "At which step will most users pause, hesitate, or abandon? What makes that moment hard?"

5. **Thinking style spectrum (Indi Young):** "Who is the cautious user who double-checks everything vs. the fast mover who skips instructions? Does the design need to serve both?"

6. **Accessibility context:** "Who might be excluded by the obvious implementation? (screen reader, keyboard-only, low bandwidth, cognitive load under stress)"

---

### Phase 2: Autonomous Codebase Research

Use codebase access (no human UX researcher has this):

- Scan for existing features with similar user-visible surfaces → find established patterns to follow
- Check for existing accessibility patterns in the codebase
- Cross-reference stated pain points against current error handling → flag "user sees generic error when X happens"
- Flag technical decisions in @idea's output that will create Gulf of Execution/Evaluation problems
- Generate **UX Risk Register**: a ranked list of user-visible failure modes

---

### Output: `docs/ux/{feature}.md`

Create file with YAML frontmatter and prose sections. @design reads this when present and converts `friction_points` and `ux_risks` into acceptance criteria.

```yaml
---
user_context: "[description of the moment the user reaches for this feature]"
mental_model_gap: "[where user belief ≠ system reality]"
friction_points:
  - step: "[step name]"
    risk: high|medium|low
    description: "[what makes this moment hard]"
    recommendation: "[design mitigation]"
accessibility_notes:
  - "[specific exclusion risk and mitigation]"
thinking_styles:
  cautious_user: "[how design must accommodate them]"
  fast_user: "[how design must accommodate them]"
ux_risks:
  - "[ranked list of user-visible failure modes]"
validated_workaround: "[what users do today]"
---

## Summary

[Brief prose summary for human readers]
```

---

## Auto-Trigger Heuristic (when invoked by @feature)

**Run @ux when:**
- @idea output contains user-facing keywords: `ui`, `user`, `interface`, `dashboard`, `form`, `flow`, `UX`, `screen`, `page`, `button`
- AND absent: `K8s`, `CRD`, `reconciler`, `stream`, `JetStream`, `CLI-only` (explicit infra signals)

**Skip @ux when:**
- `@feature "..." --infra` flag is set
- Feature is clearly infrastructure-only (no user-visible surface)

---

## When to Use

- **Standalone:** `@ux user-authentication` — UX research for any existing feature or idea
- **Via @feature:** Auto-triggered between @idea and @design for user-facing features

---

## Output

**Primary:** `docs/ux/{feature}.md` with typed YAML schema

**Consumed by:** @design — converts `friction_points` and `ux_risks` into workstream acceptance criteria

---

## Write Plan (F101)

Before creating the UX research file, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (UX findings doc, event log).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"ux"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @ux <feature>:
  CREATE: docs/ux/{feature}.md — UX research findings with YAML frontmatter
  MODIFY: .sdp/log/events.jsonl — Write plan event log

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

---

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- `@feature` - Orchestrator that auto-triggers @ux
- `@design` - Reads docs/ux/ when present, adds UX acceptance criteria to workstreams
- `@idea` - Produces input that @ux analyzes for UX risks
