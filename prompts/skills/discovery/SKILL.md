---
name: discovery
description: Pre-requirements product discovery gate (roadmap check, research loop, feature brief)
version: 1.0.0
depends_on: "@feature v8"
changes:
  - Initial release: 4 phases, 3 routing tracks (Obvious / Competitive / Novel)
---

# @discovery - Product Discovery Gate

**Validate before specifying.** Answer "should we build this?" before "how should we build this?"

---

## EXECUTE THIS NOW

When user invokes `@discovery "feature description"` or when `@feature` invokes it (unless `--quick`):

### Phase 1: Roadmap Pre-Check

1. Extract 3-5 high-signal keywords from the feature description (nouns + domain verbs — NOT generic terms like "add", "update", "implement").
2. Search docs for overlap:

```bash
rg "<keyword1>|<keyword2>|<keyword3>" docs/ --type-add 'md:*.md' -t md -l
```

3. Analyze results for:
   - Features in ROADMAP.md covering same domain terms
   - Workstream files with matching Scope Files or goals
   - Existing docs/drafts/idea-*.md that cover similar territory

4. Present **Overlap Report** (HIGH and MEDIUM confidence only; log LOW to file):
   ```
   Found N potentially related items:
   [HIGH] F005 Rework Loop — covers [summary]. Similarity reason: [1 sentence]
   [MEDIUM] 00-008-02 — touches [same module]. Overlap type: [data model / API / user flow]
   ```

5. User resolution (single question):
   - A) These are different — proceed to Phase 2
   - B) This extends F005 — incorporate and modify existing workstream
   - C) This supersedes F005 — flag for later review (propose: set F005 status to 'deferred')
   - D) Show me more detail before deciding

**Gate:** Proceed only after user resolves.

**Mode `--quiet`:** Phase 1 only, then stop. Output: overlap report only.

---

### Phase 2: Signal Check (~30 seconds)

1. Ask 2 questions:
   - "What user problem does this solve and for whom?"
   - "Do you know of existing solutions (libraries, tools, competitors)?"

2. Run web search: `"{feature_name} existing solutions 2026"`

3. Route to track:

| Condition | Track |
|-----------|-------|
| User answers both confidently AND search finds ≥1 clear prior art | **OBVIOUS** |
| User answers but search shows competitive landscape | **COMPETITIVE** |
| User uncertain on Q1 OR search shows no clear prior art | **NOVEL** |

4. Soft override for OBVIOUS: "You're on the Obvious track. Type 'research' to switch to Competitive."

**Mode `--skip-research`:** Phase 1+2 only, then stop. Output: overlap report + route decision.

---

### Phase 3: Product Research (track-dependent)

#### OBVIOUS Track

- Skip to `@idea --quiet` (invoked by @feature).
- No discovery brief generated.

#### COMPETITIVE Track (single research pass)

1. Web searches:

   ```
   "{feature_name} best practices {year}"
   "{feature_category} open source alternatives"
   "how does [top competitor] implement {feature_name}"
   ```

2. Synthesize:
   - Alternatives comparison table (≥3 alternatives)
   - Build-vs-adopt recommendation with rationale
   - Primary differentiator in one sentence

3. Ask 3 targeted questions:
   - Differentiation (what makes yours different?)
   - Constraints (what rules out the adopt option?)
   - Must-haves vs nice-to-haves

4. Convergence criteria (all 3 must be met):
   - ✓ ≥3 alternatives identified
   - ✓ Build-vs-adopt decision stated
   - ✓ Primary differentiator articulated in one sentence

#### NOVEL Track (iterative loop, max 3 iterations)

Each iteration targets one of Cagan's four risks:

| Iteration | Expert Role | Risk | Web Search Focus |
|-----------|-------------|------|------------------|
| 1 | Product PM | Value risk — is this a real problem worth solving? | user pain points, demand signals |
| 2 | Tech Lead | Feasibility risk — can we build this well? | technical patterns, implementation complexity |
| 3 | DevRel/Strategist | Strategic fit risk — does this belong in the product? | roadmap alignment, user segment fit |

Per iteration: form hypothesis → 1-2 web searches → simulate expert (use @think internally) → ask user ONE clarifying question → update risk score.

**JTBD Convergence:** Loop stops when the user can articulate the feature in Jobs-to-be-Done format: `"When [situation], I want to [motivation], so I can [outcome]"` AND all 3 risk scores ≥ 3/5 (total ≥ 9/15).

---

### Phase 4: Feature Brief (COMPETITIVE and NOVEL tracks only)

Generate `docs/drafts/discovery-{slug}.md`:

```markdown
## Feature Brief: {Name}

### Opportunity Statement
When [situation], [user segment] want to [motivation], so they can [outcome].

### Market Context
- Existing alternatives: [table]
- Build rationale: [why build vs adopt]
- Differentiation: [one sentence]

### Validated Assumptions
- Value risk: [score]/5 — [evidence]
- Feasibility risk: [score]/5 — [evidence]
- Strategic fit: [score]/5 — [evidence]

### Open Questions
- [Q1]: [answer or "unresolved"]

### Research Context (for @idea)
- Alternatives: [list]
- Key constraints: [list]
- Pre-answered cycles: [Vision ✓, Problem ✓, ...]
```

---

## Modes

| Mode | Phases | Output |
|------|--------|--------|
| Default | 1–4 | Full discovery brief |
| `--quiet` | 1 only | Overlap report |
| `--skip-research` | 1+2 | Overlap report + route decision |

---

## When to Use

- **Standalone:** `@discovery "auth"` — pre-check only, produces discovery brief, stops
- **Via @feature:** `@feature "auth"` invokes @discovery before @idea (unless `--quick`)

---

## Output

**Primary:** `docs/drafts/discovery-{slug}.md` (COMPETITIVE / NOVEL tracks)

**Secondary:** Overlap report (presented to user); route decision (OBVIOUS / COMPETITIVE / NOVEL)

---

## Next Steps

- **OBVIOUS:** @feature continues to @idea --quiet
- **COMPETITIVE / NOVEL:** @feature passes discovery brief to @idea with `--spec docs/drafts/discovery-{slug}.md`

---

## See Also

- `@feature` - Orchestrator that invokes @discovery
- `@idea` - Requirements gathering (receives discovery output as --spec)
- `@think` - Internal expert simulation for NOVEL track
