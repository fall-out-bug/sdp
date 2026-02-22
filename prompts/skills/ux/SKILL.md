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

## See Also

- `@feature` - Orchestrator that auto-triggers @ux
- `@design` - Reads docs/ux/ when present, adds UX acceptance criteria to workstreams
- `@idea` - Produces input that @ux analyzes for UX risks
