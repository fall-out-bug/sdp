---
name: idea
description: Interactive requirements gathering through deep interviewing using AskUserQuestion. Creates feature draft with goals, scope, and open questions.
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
---

# /idea - Requirements Gathering

Deep, interactive interviewing to capture comprehensive feature requirements using AskUserQuestion tool.

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive requirements document
- Want to explore tradeoffs and edge cases

## Invocation

```bash
/idea "feature description"
# or with existing spec
/idea "feature description" --spec path/to/SPEC.md
```

## Workflow

**IMPORTANT:** Use AskUserQuestion for deep, continuous interviewing until requirements are complete.

### Step 1: Read Context (If Exists)

If user provides existing spec or similar features exist:

```bash
# Read existing spec
@path/to/SPEC.md

# Search for similar features
Glob("docs/specs/**/*")
Grep("similar feature keywords")
```

### Step 2: Initial Interview

**Ask foundational questions using AskUserQuestion:**

```markdown
AskUserQuestion({
  "questions": [{
    "question": "What is the primary problem this feature solves?",
    "header": "Problem",
    "options": [
      {"label": "User pain point", "description": "Addresses frustration or inefficiency"},
      {"label": "Business requirement", "description": "Enables new revenue or reduces cost"},
      {"label": "Technical debt", "description": "Improves maintainability or performance"},
      {"label": "Competitive parity", "description": "Matches competitor capabilities"}
    ],
    "multiSelect": false
  }, {
    "question": "Who are the primary users of this feature?",
    "header": "Users",
    "options": [
      {"label": "End users", "description": "Direct product users"},
      {"label": "Administrators", "description": "System managers and ops teams"},
      {"label": "Developers", "description": "Engineering team integration"},
      {"label": "API consumers", "description": "External integrations"}
    ],
    "multiSelect": true
  }]
})
```

### Step 3: Deep Dive Interview

**Continue with progressively detailed questions. Be VERY in-depth:**

**Technical Implementation:**
```markdown
AskUserQuestion({
  "questions": [{
    "question": "How should data be persisted for this feature?",
    "header": "Data Storage",
    "options": [
      {"label": "Relational DB (Recommended)", "description": "ACID guarantees, complex queries. Adds DB dependency."},
      {"label": "NoSQL", "description": "Flexible schema, horizontal scaling. Eventual consistency."},
      {"label": "In-memory cache", "description": "Fast, ephemeral. Data loss on restart."},
      {"label": "File system", "description": "Simple, no DB. Limited query capabilities."}
    ],
    "multiSelect": false
  }, {
    "question": "What happens if this feature's service is unavailable?",
    "header": "Failure Mode",
    "options": [
      {"label": "Graceful degradation", "description": "Fallback to limited functionality"},
      {"label": "Fail closed", "description": "Block operation entirely"},
      {"label": "Queue and retry", "description": "Async processing with eventual completion"},
      {"label": "Best-effort", "description": "Continue with warnings"}
    ],
    "multiSelect": false
  }]
})
```

**UI & UX:**
```markdown
AskUserQuestion({
  "questions": [{
    "question": "Where should this feature be accessible in the UI?",
    "header": "UI Location",
    "options": [
      {"label": "Main navigation", "description": "High visibility, always accessible"},
      {"label": "Settings page", "description": "One-time configuration"},
      {"label": "Contextual menu", "description": "Appears when relevant"},
      {"label": "Dashboard widget", "description": "Summary view with drill-down"}
    ],
    "multiSelect": true
  }, {
    "question": "How should users discover this feature?",
    "header": "Discoverability",
    "options": [
      {"label": "Onboarding tutorial", "description": "Guided intro for new users"},
      {"label": "In-app tooltip", "description": "Passive education on first encounter"},
      {"label": "Documentation", "description": "Self-serve learning"},
      {"label": "Assume familiarity", "description": "Power user feature, no hand-holding"}
    ],
    "multiSelect": false
  }]
})
```

**Concerns & Tradeoffs:**
```markdown
AskUserQuestion({
  "questions": [{
    "question": "What is your primary concern about this feature?",
    "header": "Main Concern",
    "options": [
      {"label": "Performance impact", "description": "May slow down existing operations"},
      {"label": "Complexity", "description": "Increases codebase maintenance burden"},
      {"label": "Security risk", "description": "New attack surface or data exposure"},
      {"label": "User confusion", "description": "May not be intuitive, requires support"}
    ],
    "multiSelect": true
  }, {
    "question": "If scope must be reduced, what is negotiable?",
    "header": "Scope Priority",
    "options": [
      {"label": "Advanced features (Recommended)", "description": "Start with MVP, iterate based on usage"},
      {"label": "UI polish", "description": "Functional but basic interface initially"},
      {"label": "Edge case handling", "description": "Handle 80% of cases, document limitations"},
      {"label": "Nothing", "description": "All requirements are critical"}
    ],
    "multiSelect": false
  }]
})
```

**Continue interviewing** until you have comprehensive answers about:
- Technical implementation details
- UI/UX specifics
- Error handling approach
- Performance considerations
- Security implications
- Testing strategy
- Deployment approach
- Monitoring and observability

**DON'T ask obvious questions.** Instead, ask about:
- Ambiguities in requirements
- Hidden assumptions
- Edge cases
- Failure modes
- Integration points
- Tradeoffs between approaches

### Step 4: Load Master Prompt

After interviewing is complete:

```bash
cat prompts/commands/idea.md
```

### Step 5: Create Draft

Write comprehensive spec to `docs/drafts/idea-{slug}.md` including:
- All answers from AskUserQuestion
- Technical decisions made
- Tradeoffs discussed
- Open questions remaining

## Output

Draft: `docs/drafts/idea-{slug}.md`

**Comprehensive sections:**
- **Context & Problem** — from interview about problem type
- **Goals & Non-Goals** — from scope priority questions
- **User Stories** — from primary users and use cases
- **Technical Approach** — from implementation questions (storage, failure modes, etc.)
- **UI/UX Design** — from UI location, discoverability questions
- **Concerns & Risks** — from concerns interview
- **Tradeoffs** — from tradeoff questions (performance vs simplicity, etc.)
- **Open Questions** — remaining ambiguities

## Key Principles

**Interviewing Strategy:**
1. **Start broad, go deep** — foundational questions first, then drill into details
2. **No obvious questions** — don't ask "should we test?" Ask "integration tests or just unit tests?"
3. **Expose tradeoffs** — every option should show pros/cons in description
4. **Continue until complete** — keep asking until no ambiguities remain
5. **Capture decisions** — record why certain approaches were chosen/rejected

**Question Design:**
- **Header**: 8-12 chars, clear category
- **Question**: Specific, unambiguous
- **Options**: 2-4 distinct choices
- **Descriptions**: Show tradeoffs, not just features
- **MultiSelect**: Use when choices aren't mutually exclusive

**Example of GOOD question:**
```markdown
{
  "question": "How should rate limiting be enforced?",
  "header": "Rate Limit",
  "options": [
    {"label": "Token bucket", "description": "Smooth traffic, allows bursts. Requires distributed state."},
    {"label": "Fixed window", "description": "Simple, stateless. Allows burst at window boundary."},
    {"label": "Sliding window", "description": "Precise, fair. Higher memory/compute cost."}
  ],
  "multiSelect": false
}
```

**Example of BAD question:**
```markdown
{
  "question": "Should we add logging?",  // ← TOO OBVIOUS
  "header": "Logging",
  "options": [
    {"label": "Yes"},  // ← No tradeoffs shown
    {"label": "No"}
  ]
}
```

## Quick Reference

**Input:** Feature idea (+ optional existing spec)
**Output:** `docs/drafts/idea-{slug}.md` (comprehensive)
**Next:** `/design idea-{slug}`

**Typical flow:**
1. 2-3 rounds of AskUserQuestion (6-12 total questions)
2. ~15-20 minutes of user time
3. Results in 2-3 page spec with clear decisions
