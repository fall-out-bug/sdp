---
name: design
description: System design with progressive disclosure
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 5.1.0
changes:
  - Added Phase 5: Cross-feature boundary detection
  - Added sdp collision detect --deep integration
  - Added sdp contract generate for shared boundaries
  - AC4 implementation for F060
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

1. **Check for Cross-Feature Boundaries** (NEW - AC4 for F060)
   ```bash
   sdp collision detect
   ```
   - Analyzes scope files for shared types/interfaces across parallel features
   - Reports: shared types, fields needed by each feature, merge recommendations
   - If boundaries found → suggest shared contracts
   - JSON output: `sdp collision detect --output-json`

2. **Generate Shared Contracts** (if boundaries detected)
   ```bash
   sdp contract generate --features=F054,F055
   ```
   - Creates `.contracts/<type>.yaml` files for shared boundaries
   - Contract includes: typeName, fields, requiredBy features, status
   - Example: `.contracts/User.yaml`

3. **Lock Shared Contracts**
   ```bash
   sdp contract lock .contracts/User.yaml
   ```
   - Creates .lock file with SHA256 checksum
   - Prevents modifications during implementation

4. **API Contract Synthesis** (if applicable)
   ```bash
   sdp contract synthesize \
     --feature=<feature-name> \
     --requirements=<idea-doc>
   ```
   - OpenAPI 3.0 contract for API endpoints
   - Endpoints, methods, request/response schemas

5. **Multi-Agent Review** (if contracts exist)
   - Frontend Agent: "Need /batch endpoint"
   - Backend Agent: "Works for us"
   - Security Agent: "Add authentication headers"
   - Synthesizer resolves conflicts

**Exit Criteria:**
- [ ] Cross-feature boundaries checked (if parallel features)
- [ ] Shared contracts generated (if boundaries found)
- [ ] Shared contracts locked (if generated)
- [ ] API contracts generated (if applicable)
- [ ] No ERROR-level conflicts

**Skip Only If:**
- Feature has NO shared boundaries (single feature)
- Feature has NO API contracts (pure computation)

### Phase 6: Workstream Generation

**Generate workstreams based on:**
- Shared contracts (from Phase 5)
- API contracts (from Phase 5)
- Architecture decisions (from discovery blocks)

**Output:** Workstream files in `docs/workstreams/backlog/`

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

**Version:** 5.1.0 (Cross-Feature Contract Detection - F060 AC4)
**See Also:** `@idea`, `@build`, `@oneshot`
**Related Features:** F054 (scope collision), F060 (cross-feature boundaries)
