# 00-052-03: Update CLAUDE.md with @vision

> **Beads ID:** sdp-luo8
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1A - Strategic Skills (@vision)
> **Size:** SMALL
> **Duration:** 1-2 days
> **Dependencies:**
> - 00-052-01 (@vision Skill Structure)

## Goal

Add @vision to CLAUDE.md decision tree and examples.

## Acceptance Criteria

- **AC1:** CLAUDE.md has "Four-Level Planning Model" section (@vision/@reality/@feature/@oneshot)
- **AC2:** @vision added to decision trees
- **AC3:** Example workflow showing @vision → @reality → @feature → @oneshot
- **AC4:** When to use each level documented
- **AC5:** Git commit with conventional commit message

## Files

**Modify:**
- `CLAUDE.md` - Add @vision documentation

## Steps

### Step 1: Read Current CLAUDE.md

Review existing decision tree section to understand current structure.

### Step 2: Add Four-Level Planning Model Section

Add new section after "Quick Start" in CLAUDE.md:

```markdown
## Decision Tree: @vision → @reality → @feature → @oneshot

### Four-Level Planning Model

SDP has four orchestrators for different planning levels:

| Level | Purpose | Input | Output | Duration |
|-------|---------|-------|--------|----------|
| **@vision** | Strategic product planning | Product idea | PRODUCT_VISION.md, PRD.md, ROADMAP.md | Quarterly/annual |
| **@reality** | Codebase analysis | Project directory | Reality report (health, gaps, debt) | Per project or quarterly |
| **@feature** | Feature planning | Feature description | Workstreams (00-FFF-SS.md) | Per feature |
| **@oneshot** | Feature execution | Workstreams | Implemented code + deployed feature | Per feature |

### When to Use @vision

Use @vision when:
- ✅ Starting a new project or product
- ✅ Quarterly strategic review
- ✅ Major pivot or direction change
- ✅ Need comprehensive product analysis
- ✅ Want expert analysis across 7 dimensions (product, market, technical, UX, business, growth, risk)

**Skip @vision if:**
- Product vision already exists (PRODUCT_VISION.md present)
- Working on existing product (not new project)
- Incremental feature (not major pivot)

### When to Use @reality

Use @reality when:
- ✅ New to project (what's actually here?)
- ✅ Before @feature (what can we build on?)
- ✅ After @vision (how do docs match code?)
- ✅ Quarterly review (track tech debt and quality trends)
- ✅ Debugging mysteries (why doesn't this work?)
- ✅ Want 8-expert codebase analysis (architecture, quality, testing, security, performance, docs, debt, standards)

### Typical Full Flow

```bash
# Step 1: Strategic planning (quarterly or new project)
@vision "AI-powered task manager for remote teams"
# → 7 expert agents analyze product, market, technical, UX, business, growth, risk
# → Generates: PRODUCT_VISION.md, PRD.md, ROADMAP.md

# Step 2: Feature planning (per feature)
@feature "Add user authentication"
# → @idea gathers requirements
# → @design creates workstreams

# Step 3: Execution (autonomous)
@oneshot F01
# → @build executes workstreams
# → @review checks quality
# → @deploy merges to main
```
```

### Step 3: Update Quick Start Example

Update quick start to include @vision:

```markdown
### Quick Start (4 Commands)

```bash
@vision "AI-powered task manager"     # Strategic planning
@reality --quick                      # What's actually there?
@feature "Add user authentication"    # Plan feature
@oneshot F01                          # Execute workstreams
```
```

### Step 4: Verify Changes

```bash
git diff CLAUDE.md
```

Expected:
- New section: "Decision Tree: @vision → @reality → @feature → @oneshot"
- Table comparing 4 levels
- When to use @vision documented
- Skip @vision criteria documented

### Step 5: Commit Changes

```bash
git add CLAUDE.md
git commit -m "docs: add @vision to CLAUDE.md decision tree

- Add four-level planning model section
- Document when to use @vision vs @feature
- Add full workflow example
- Update quick start to include @vision"
```

## Quality Gates

- CLAUDE.md valid markdown
- No broken links
- Examples are accurate
- Consistent with existing documentation style

## Success Metrics

- User can understand when to use @vision
- Decision tree is clear
- Examples work as documented
