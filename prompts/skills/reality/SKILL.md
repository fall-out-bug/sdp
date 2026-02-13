---
name: reality
description: Codebase analysis and architecture validation - what's actually there vs documented
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @reality - Codebase Analysis & Architecture Validation

**Analyze what's actually in your codebase (vs. what's documented).**

---

## EXECUTE THIS NOW

When user invokes `@reality`:

### Step 0: Auto-Detect Project Type

```bash
# Detect language/framework
if [ -f "go.mod" ]; then PROJECT_TYPE="go"
elif [ -f "pyproject.toml" ] || [ -f "requirements.txt" ]; then PROJECT_TYPE="python"
elif [ -f "pom.xml" ] || [ -f "build.gradle" ]; then PROJECT_TYPE="java"
elif [ -f "package.json" ]; then PROJECT_TYPE="nodejs"
else PROJECT_TYPE="unknown"
fi
```

### Step 1: Quick Scan (--quick mode)

**Analysis:**
1. Project size (lines of code, file count)
2. Architecture (layer violations, circular dependencies)
3. Test coverage (if tests exist, estimate %)
4. Documentation (doc coverage, drift detection)
5. Quick smell check (TODO/FIXME/HACK comments, long files)

**Output:** Health Score X/100 + Top 5 Issues

### Step 2: Deep Analysis (--deep mode)

Run 8 parallel expert analyses:

1. **Architecture Expert** - Layer mapping, dependencies, violations
2. **Code Quality Expert** - File size, complexity, duplication
3. **Testing Expert** - Coverage, test quality, frameworks
4. **Security Expert** - Secrets, OWASP, dependencies
5. **Performance Expert** - Bottlenecks, caching, scalability
6. **Documentation Expert** - Coverage, drift, quality
7. **Technical Debt Expert** - TODO/FIXME, code smells
8. **Standards Expert** - Conventions, error handling, types

### Step 3: Synthesize Report

Create comprehensive report with:
- Executive Summary with Health Score
- Critical Issues (Fix Now)
- Quick Wins (Fix Today)
- Detailed Analysis from each expert
- Action Items (This Week / This Month / This Quarter)

---

## When to Use

- **New to project** - "What's actually here?"
- **Before @feature** - "What can we build on?"
- **After @vision** - "How do docs match code?"
- **Quarterly review** - Track tech debt and quality trends
- **Debugging mysteries** - "Why doesn't this work?"

---

## Modes

| Mode | Duration | Purpose |
|------|----------|---------|
| `--quick` | 5-10 min | Health check + top issues |
| `--deep` | 30-60 min | Comprehensive with 8 experts |
| `--focus=topic` | Varies | Single expert deep dive |

**Focus topics:** security, architecture, testing, performance

---

## Output

```
## Reality Check: {project_name}

### Quick Stats
- Language: {detected}
- Size: {LOC} lines, {N} files
- Architecture: {layers detected}
- Tests: {coverage if available}

### Top 5 Issues
1. {issue} - {severity}
   - Location: {file:line}
   - Impact: {why it matters}
   - Fix: {recommendation}

### Health Score: {X}/100
```

---

## Vision Integration

If PRODUCT_VISION.md exists, compare reality to vision:

| Feature | PRD Status | Reality Status | Gap |
|---------|------------|----------------|-----|
| Feature 1 | P0 | Implemented | None |
| Feature 2 | P1 | Partial | Missing X |
| Feature 3 | P0 | Not found | Not started |

---

## Examples

```bash
@reality --quick              # Quick health check
@reality --deep               # Deep analysis
@reality --focus=security     # Security only
@reality --deep --output=docs/reality/check.md  # Save report
```

---

## See Also

- `@vision` - Strategic planning
- `@feature` - Feature planning
- `@idea` - Requirements gathering
