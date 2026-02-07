# 00-052-06: Update CLAUDE.md with @reality

> **Beads ID:** sdp-kxif
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1B - Analysis Skills (@reality)
> **Size:** SMALL
> **Duration:** 1-2 days
> **Dependencies:**
> - 00-052-04 (@reality Skill Structure)

## Goal

Add @reality examples and integration to CLAUDE.md.

## Acceptance Criteria

- **AC1:** CLAUDE.md has @reality examples
- **AC2:** @reality added to decision tree examples
- **AC3:** Comparison table: @vision vs @reality vs @feature vs @oneshot
- **AC4:** Typical workflow with @reality documented
- **AC5:** Git commit with conventional commit message

## Files

**Modify:**
- `CLAUDE.md` - Add @reality documentation

## Steps

### Step 1: Read Current CLAUDE.md

Review to understand existing structure (already updated with @vision in 00-052-03).

### Step 2: Add @reality Examples

Add section after @vision documentation:

```markdown
### When to Use @reality

Use @reality when:
- ✅ New to project (what's actually here?)
- ✅ Before @feature (what can we build on?)
- ✅ After @vision (how do docs match code?)
- ✅ Quarterly review (track tech debt and quality trends)
- ✅ Debugging mysteries (why doesn't this work?)
- ✅ Want 8-expert codebase analysis (architecture, quality, testing, security, performance, docs, debt, standards)

**Skip @reality if:**
- Just ran @reality recently (no significant changes)
- Only making trivial modifications
- Working on greenfield project (no code to analyze)

### @reality Modes

**Quick mode (5-10 min):**
```bash
@reality --quick
```
Fast health check: language, frameworks, structure, obvious issues.

**Deep mode (30-60 min):**
```bash
@reality --deep
```
Comprehensive analysis with 8 expert agents.

**Focus mode (10-15 min):**
```bash
@reality --focus=security
@reality --focus=architecture
@reality --focus=testing
```
Deep dive on specific aspect.
```

### Step 3: Update Comparison Table

Update four-level comparison table to include @reality details:

```markdown
| Aspect | @vision | @reality | @feature | @oneshot |
|--------|---------|----------|----------|----------|
| **Purpose** | Strategic product planning | Codebase analysis | Feature planning (requirements + workstreams) | Execute workstreams |
| **Input** | Product idea | Project directory | Feature description | Feature ID (F01) or WS list |
| **Output** | PRODUCT_VISION.md, PRD.md, ROADMAP.md | Reality report (health, gaps, debt) | Workstreams (00-FFF-SS.md) | Implemented code + deployed feature |
| **Modes** | Standard | --quick, --deep, --focus | Interactive | Autonomous |
| **Experts** | 7 experts (product, market, technical, UX, business, growth, risk) | 8 experts (architecture, quality, testing, security, performance, docs, debt, standards) | @idea + @design | @build + @review + @deploy |
| **Duration** | Quarterly/annual | Per project or quarterly | Per feature | Per feature |
```

### Step 4: Add @vision → @reality Integration Example

```markdown
### Vision vs Reality Gap Analysis

After running both @vision and @reality, compare:

```bash
# Step 1: Create vision
@vision "AI-powered task manager"
# → Generates PRODUCT_VISION.md, PRD.md

# Step 2: Check reality (for existing project)
@reality --deep
# → Generates REALITY_REPORT.md

# Step 3: Compare gap
# Vision says: "Clean architecture"
# Reality shows: "Layer violations in 3 files"
# → Action: Refactor before adding features
```
```

### Step 5: Commit Changes

```bash
git add CLAUDE.md
git commit -m "docs: add @reality examples to CLAUDE.md

- Document @reality modes (--quick, --deep, --focus)
- Update comparison table
- Add vision vs reality gap analysis example"
```

## Quality Gates

- CLAUDE.md valid markdown
- Examples are accurate and tested
- Consistent with existing documentation style
- No broken links

## Success Metrics

- User understands when to use @reality
- Examples work as documented
- Clear distinction from @vision
