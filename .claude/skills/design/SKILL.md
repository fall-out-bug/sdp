---
name: design
description: System design with progressive disclosure
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 4.0.0
---

# @design - System Design with Progressive Disclosure

Multi-agent system design (Arch + Security + SRE) with progressive discovery blocks.

## When to Use

- After @idea requirements gathering
- Need architecture decisions
- Creating workstream breakdown

## Invocation

```bash
@design <task_id>
@design <task_id> --quiet    # Minimal design blocks
@design "feature description"  # Skip @idea, design directly
```

## Progressive Discovery Workflow

### Overview

**Discovery Blocks:** 3-5 focused blocks (not one big questionnaire)

**Block Structure:**
- Each block: 3 questions
- After each block: trigger point (Continue / Skip block / Done)
- User can skip blocks not relevant to feature

### Discovery Blocks

**Block 1: Data & Storage (3 questions)**
- Data models?
- Storage requirements?
- Persistence strategy?

**Block 2: API & Integration (3 questions)**
- API endpoints?
- External integrations?
- Authentication/authorization?

**Block 3: Architecture (3 questions)**
- Component structure?
- Layer boundaries?
- Error handling strategy?

**Block 4: Security (3 questions)**
- Input validation?
- Sensitive data handling?
- Rate limiting?

**Block 5: Operations (3 questions)**
- Monitoring?
- Deployment?
- Rollback strategy?

### Phase 5: Contract Synthesis (CRITICAL - Before Implementation)

**Purpose:** Multi-agent contract agreement BEFORE parallel implementation begins.

**Why Now:** Prevents component integration failures (404 Not Found) when agents work in parallel.

**Workflow:**

1. **Architect Proposes Initial Contract**
   ```bash
   sdp contract synthesize \
     --feature=<feature-name> \
     --requirements=<idea-doc> \
     --output=.contracts/<feature-name>.yaml
   ```
   - Analyzes requirements from @idea
   - Proposes OpenAPI 3.0 contract
   - Defines endpoints, methods, request/response schemas

2. **Multi-Agent Review (Parallel)**
   - Frontend Agent: "Need /batch endpoint"
   - Backend Agent: "Works for us"
   - SDK Agent: "Matches our method naming"
   - Security Agent: "Add authentication headers"

3. **Synthesizer Resolves Conflicts**
   - Unanimous agreement → Contract locked
   - Domain expertise veto → Escalate to human
   - Merge suggestions → Update contract
   - Escalation → Human decides

4. **Lock Contract**
   ```bash
   sdp contract lock \
     --contract=.contracts/<feature-name>.yaml \
     --reason="Multi-agent agreement complete"
   ```
   - Creates .lock file
   - Stores SHA256 checksum
   - Prevents modifications during implementation

5. **Create Validation Workstreams**
   - Add WS: Contract validation (post-implementation)
   - Add WS: Integration testing

**Exit Criteria:**
- [ ] Contract file exists: `.contracts/<feature-name>.yaml`
- [ ] Contract locked: `.contracts/<feature-name>.yaml.lock`
- [ ] No ERROR-level conflicts
- [ ] Validation workstream created

**Skip Only If:**
- Feature has NO API contracts (pure computation)
- Feature has single component (no integration risk)

### After Each Block: Trigger Point

```markdown
AskUserQuestion({
  "questions": [
    {"question": "Block complete. Continue to next block?",
     "header": "Discovery",
     "options": [
       {"label": "Continue (Recommended)", "description": "Next discovery block"},
       {"label": "Skip block", "description": "Skip remaining blocks"},
       {"label": "Done", "description": "Generate workstreams with current info"}
     ]}
  ]
})
```

## Integration with @idea

```python
# Uses requirements from @idea
idea_result = load_idea_result(task_id)

# Skip already covered topics
skip_topics = idea_result.covered_topics

# Focus on design-specific questions
design_blocks = filter_blocks(skip_topics)
```

## --quiet Mode

Minimal blocks (2 blocks, 6 questions):
1. Data & Storage
2. Core Architecture

## Output

**Primary:** Workstream files in `docs/workstreams/backlog/`

**Secondary:**
- `docs/drafts/<task_id>-design.md` - Design document

## Next Steps

```bash
@oneshot <feature>  # Execute all workstreams (includes contract validation)
@build <ws_id>      # Execute workstream
```

**Note:** Contract validation workstream (created in Phase 5) runs AFTER implementation to detect drift.

---

**Version:** 5.0.0 (Contract Synthesis Phase)
**See Also:** `@idea`, `@build`, `@oneshot`
