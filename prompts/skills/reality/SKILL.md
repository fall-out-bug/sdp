---
name: reality
description: Codebase analysis and architecture validation - what's actually there vs documented
tools: Read, Glob, Grep, Bash, Task, Write, Edit
version: 1.0.0
---

# @reality - Codebase Analysis & Architecture Validation

**Analyze what's actually in your codebase (vs. what's documented).**

## When to Use
- **New to project** - "What's actually here?"
- **Before @feature** - "What can we build on?"
- **After @vision** - "How do docs match code?"
- **Quarterly review** - `@reality --review` - track tech debt and quality trends
- **Debugging mysteries** - "Why doesn't this work?"

## Verbosity Tiers

```bash
@reality --quick --quiet     # Exit status only: ✅
@reality --quick             # Summary: ✅ Reality Check: 15K LOC, Health 72/100, 3 critical issues
@reality --quick --verbose   # Step-by-step progress
@reality --quick --debug     # Internal state + API calls
```

**Examples:**

```bash
# Quiet mode
@reality --quick --quiet
# Output: ✅

# Default mode
@reality --quick
# Output: ✅ Reality Check: 15K LOC, Health 72/100, 3 critical issues

# Verbose mode
@reality --quick --verbose
# Output:
# → Project scan: 15,234 LOC, 127 files (Go)
# → Architecture: Layer violations detected (2)
# → Testing: Coverage ~65% (below 80% threshold)
# → Documentation: 4 drift issues found
# → Quick smell check: 15 TODO/FIXME comments
# ⚠️ 3 critical issues found
# ✅ COMPLETE

# Debug mode
@reality --quick --debug
# Output:
# [DEBUG] Mode: quick
# [DEBUG] Detecting project type...
# [DEBUG] Checking for go.mod...
# [DEBUG] Found go.mod → Go project
# → Project scan: 15,234 LOC, 127 files (Go)
# [DEBUG] Scanning src/ directory...
# [DEBUG] Found 87 .go files
# [DEBUG] Analyzing imports...
# [DEBUG] Checking for circular dependencies...
# → Architecture: Layer violations detected (2)
# [DEBUG] Scanning tests/ directory...
# [DEBUG] Found 45 test files
# [DEBUG] Calculating coverage...
# → Testing: Coverage ~65% (below 80% threshold)
# [DEBUG] Comparing docs to code...
# [DEBUG] Checking 23 documentation files
# → Documentation: 4 drift issues found
# [DEBUG] Searching for TODO/FIXME/HACK...
# → Quick smell check: 15 TODO/FIXME comments
# ⚠️ 3 critical issues found
# ✅ COMPLETE
```

## Modes

```bash
@reality                    # Auto-detect mode (based on project size)
@reality --quick            # Quick scan (5-10 min): health check + top issues
@reality --deep             # Deep analysis (30-60 min): comprehensive with all experts
@reality --focus=security   # Focused analysis: single expert deep dive
@reality --focus=architecture
@reality --focus=testing
@reality --focus=performance
```

## Workflow

### Step 0: Auto-Detect Project Type

```bash
# Detect language/framework
if [ -f "go.mod" ]; then
    PROJECT_TYPE="go"
elif [ -f "pyproject.toml" ] || [ -f "requirements.txt" ]; then
    PROJECT_TYPE="python"
elif [ -f "pom.xml" ] || [ -f "build.gradle" ]; then
    PROJECT_TYPE="java"
elif [ -f "package.json" ]; then
    PROJECT_TYPE="nodejs"
else
    PROJECT_TYPE="unknown"
fi
```

### Step 1: Quick Scan (--quick mode, 5-10 min)

**Goal:** Health check + top 5 issues

**Analysis:**
1. **Project size** (lines of code, file count)
2. **Architecture** (layer violations, circular dependencies)
3. **Test coverage** (if tests exist, estimate %)
4. **Documentation** (doc coverage, drift detection)
5. **Quick smell check** (TODO/FIXME/HACK comments, long files)

**Output:**
```markdown
## Reality Check: {project_name}

### Quick Stats
- **Language:** {detected}
- **Size:** {LOC} lines, {N} files
- **Architecture:** {layers detected}
- **Tests:** {coverage if available}

### Top 5 Issues
1. ⚠️ **{issue}** - {severity}
   - Location: {file:line}
   - Impact: {why it matters}
   - Fix: {recommendation}

2. ⚠️ **{issue}** - {severity}
   ...

### Health Score: {X}/100
```

### Step 2: Deep Analysis (--deep mode, 30-60 min)

**Goal:** Comprehensive analysis with 8 expert agents

Spawn parallel expert agents via Task tool:

```python
# Agent 1: Architecture
Task(
    subagent_type="general-purpose",
    prompt="""You are the ARCHITECTURE expert.

PROJECT: {project_name}
TYPE: {detected language}

Your task:
1. Map layers (domain/, application/, infrastructure/ or equivalent)
2. Analyze dependencies:
   - Import graph analysis
   - Circular dependencies
   - Layer violations (domain → infrastructure is BAD)
3. Identify architectural patterns:
   - Clean architecture? DDD? MVC?
   - Service boundaries?
   - Module organization?

Output:
## Architecture Report
**Pattern:** {detected pattern}
**Layers:** {count}
**Violations:** {count}

### Issues
1. **{layer} imports {layer}** (violation)
   - File: {path}
   - Fix: {recommendation}

### Strengths
- {what's working well}

### Verdict:** ✅ HEALTHY / ⚠️ ISSUES / ❌ CRITICAL
""",
    description="Architecture analysis"
)

# Agent 2: Code Quality
Task(
    subagent_type="general-purpose",
    prompt="""You are the CODE QUALITY expert.

PROJECT: {project_name}

Your task:
1. File size analysis:
   - Count lines per file
   - List files >200 LOC (threshold)
   - Identify monsters >500 LOC

2. Complexity analysis:
   - Count functions per file
   - Identify functions >50 LOC
   - Find deeply nested code (>4 levels)

3. Code duplication:
   - Search for similar patterns (copy-paste)
   - Identify repeated logic
   - Find DRY violations

4. Anti-patterns:
   - God classes/functions
   - Magic numbers/strings
   - Bare exceptions
   - Hard-coded paths

Output:
## Code Quality Report
**Avg LOC/file:** {number}
**Files >200 LOC:** {count}
**Complexity:** {LOW/MEDIUM/HIGH}

### Issues
1. **{file}** ({N} LOC) - Too long
   - Recommendation: Split into {N} functions

2. **{file}** - Duplicated code
   - Duplicates: {list}
   - Recommendation: Extract to shared function

### Verdict:** ✅ CLEAN / ⚠️ REFACTOR NEEDED / ❌ CRITICAL
""",
    description="Code quality analysis"
)

# Agent 3: Testing
Task(
    subagent_type="general-purpose",
    prompt="""You are the TESTING expert.

PROJECT: {project_name}
LANGUAGE: {detected}

Your task:
1. Test discovery:
   - Find all test files (tests/, test/, *_test.go, *.test.ts)
   - Count test functions/cases
   - Identify test types (unit, integration, e2e)

2. Coverage estimation:
   - Map source files → test files
   - Identify untested modules
   - Estimate coverage % (tested functions / total functions)

3. Test quality:
   - Check for test anti-patterns (assert True, no assertions)
   - Identify brittle tests (timeouts, hard-coded paths)
   - Find missing edge cases

4. Test infrastructure:
   - Framework detection (pytest, jest, go test)
   - CI/CD integration (GitHub Actions, GitLab CI)
   - Test data management (fixtures, factories)

Output:
## Testing Report
**Tests:** {count} tests across {N} files
**Coverage:** ~{X}% ({tested}/{total} functions)
**Framework:** {detected}

### Untested
- {module/file} - {functions} untested

### Test Quality Issues
1. **{test_file}** - Brittle (hard-coded paths)
2. **{test_file}** - No assertions (smoke test only)

### Verdict:** ✅ GOOD / ⚠️ COVERAGE LOW / ❌ NO TESTS
""",
    description="Testing analysis"
)

# Agent 4: Security
Task(
    subagent_type="general-purpose",
    prompt="""You are the SECURITY expert.

PROJECT: {project_name}

Your task:
1. Secrets detection:
   - Search for API keys, passwords, tokens
   - Find hard-coded credentials
   - Check .env files, config files

2. OWASP Top 10:
   - SQL injection (string concatenation in queries)
   - XSS (user input echoed without escaping)
   - Auth bypass (missing checks)
   - CSRF (no tokens in forms)

3. Dependencies:
   - Check for known vulnerabilities (if possible)
   - Identify outdated dependencies
   - Find unmaintained packages

4. Security controls:
   - Input validation
   - Output encoding
   - Authentication/authorization
   - HTTPS enforcement

Output:
## Security Report
**Secrets found:** {count}
**OWASP issues:** {count}
**Dependency vulns:** {count (if detectable)}

### Critical Issues
1. **{file}:{line}** - Hard-coded API key
   - Action: Rotate credential, use env variable

2. **{file}** - SQL injection risk
   - Issue: Query built with string concatenation
   - Fix: Use parameterized queries

### Verdict:** ✅ SECURE / ⚠️ ISSUES / ❌ VULNERABLE
""",
    description="Security analysis"
)

# Agent 5: Performance
Task(
    subagent_type="general-purpose",
    prompt="""You are the PERFORMANCE expert.

PROJECT: {project_name}

Your task:
1. Bottleneck detection:
   - N+1 queries (database calls in loops)
   - Missing indexes (WHERE without index)
   - Large file reads (reading entire files)
   - Inefficient algorithms (O(n²) where O(n) possible)

2. Resource usage:
   - Memory leaks (unclosed connections, growing caches)
   - Connection pool exhaustion
   - Disk I/O (excessive file operations)

3. Caching:
   - Missing cache opportunities
   - Cache invalidation issues
   - Stale data

4. Scalability:
   - Synchronous operations (could be async)
   - Single-threaded bottlenecks
   - Lack of pagination

Output:
## Performance Report
**Bottlenecks:** {count}
**Memory issues:** {count}
**Scalability:** {GOOD/CONCERNS}

### Issues
1. **{file}:{line}** - N+1 query
   - Impact: O(N) database calls
   - Fix: Batch query or use JOIN

2. **{file}** - No pagination
   - Impact: Loads all records (could be 100K+)
   - Fix: Add LIMIT/OFFSET or cursor

### Opportunities
- {function} - Could cache result ({X} ms → {Y} ms)

### Verdict:** ✅ OPTIMIZED / ⚠️ IMPROVEMENTS NEEDED / ❌ CRITICAL
""",
    description="Performance analysis"
)

# Agent 6: Documentation
Task(
    subagent_type="general-purpose",
    prompt="""You are the DOCUMENTATION expert.

PROJECT: {project_name}

Your task:
1. Doc coverage:
   - README.md exists?
   - API documentation?
   - Code comments (% of functions with docstrings)

2. Documentation accuracy:
   - Compare docs to code (drift detection)
   - Find outdated examples
   - Identify missing docs for public APIs

3. Doc quality:
   - Clear installation instructions?
   - Usage examples?
   - Architecture documentation?
   - Contribution guide?

4. Generated docs:
   - Swagger/OpenAPI?
   - Godoc/PyDoc/Javadoc?
   - Auto-generated from comments?

Output:
## Documentation Report
**Coverage:** {X}% (README + API + code comments)
**Drift:** {count} inaccuracies found
**Quality:** {GOOD/FAIR/POOR}

### Missing Docs
- {module/function} - Public API undocumented
- README - No installation instructions

### Drift Detected
1. **{doc}** says {X}, but code does {Y}
   - Fix: Update doc

### Verdict:** ✅ WELL DOCUMENTED / ⚠️ GAPS / ❌ NO DOCS
""",
    description="Documentation analysis"
)

# Agent 7: Technical Debt
Task(
    subagent_type="general-purpose",
    prompt="""You are the TECHNICAL DEBT expert.

PROJECT: {project_name}

Your task:
1. Debt markers:
   - Search for TODO, FIXME, HACK, XXX comments
   - Count occurrences
   - Categorize by severity

2. Code smells:
   - Dead code (unused functions, commented-out code)
   - Flag/flag arguments (boolean parameters)
   - Feature envy (methods using other objects excessively)
   - Long parameter lists (>5 params)

3. Design debt:
   - God classes (classes doing too much)
   - Feature creep (classes with too many responsibilities)
   - Wrong abstractions (leaky implementations)

4. Infrastructure debt:
   - Deprecated dependencies
   - Unmaintained libraries
   - Outdated tooling

Output:
## Technical Debt Report
**TODO/FIXME:** {count}
**Code smells:** {count}
**Design debt:** {count}

### High Priority Debt
1. **{file}:{line}** - FIXME: {description}
   - Severity: CRITICAL
   - Recommendation: {action}

2. **{file}** - Dead code: {unused function}
   - Action: Delete or document why kept

### Debt Trend
- New debt this month: {count}
- Resolved this month: {count}

### Verdict:** ✅ MANAGEABLE / ⚠️ ACCUMULATING / ❌ CRISIS
""",
    description="Technical debt analysis"
)

# Agent 8: Standards Compliance
Task(
    subagent_type="general-purpose",
    prompt="""You are the STANDARDS expert.

PROJECT: {project_name}
LANGUAGE: {detected}

Your task:
1. Language/framework conventions:
   - Python: PEP 8 compliance
   - Go: Effective Go guidelines
   - Java: Java Code Conventions
   - JavaScript: Airbnb/Standard style

2. Error handling:
   - Proper exception handling (not bare except)
   - Error propagation patterns
   - Error logging

3. Type safety:
   - Type hints (Python)
   - Interface definitions (Go, Java)
   - Type assertions (TypeScript)

4. Best practices:
   - Naming conventions
   - File organization
   - Import management (no circular imports)
   - Resource cleanup (defer, finally, using)

Output:
## Standards Report
**Convention:** {detected style}
**Compliance:** {X}%
**Issues:** {count}

### Violations
1. **{file}:{line}** - Naming: {function should be snake_case}
2. **{file}** - Imports: {circular dependency}

### Strengths
- {what follows conventions well}

### Verdict:** ✅ COMPLIANT / ⚠️ VIOLATIONS / ❌ NONCOMPLIANT
""",
    description="Standards compliance analysis"
)
```

### Step 3: Synthesize Report

After all experts complete, create comprehensive report:

```markdown
# Reality Check: {project_name}

**Date:** {timestamp}
**Mode:** Deep Analysis
**Language:** {detected}
**Size:** {LOC} lines, {N} files

## Executive Summary

### Health Score: {X}/100

- Architecture: {score}/100
- Code Quality: {score}/100
- Testing: {score}/100
- Security: {score}/100
- Performance: {score}/100
- Documentation: {score}/100
- Technical Debt: {score}/100
- Standards: {score}/100

### Critical Issues (Fix Now)
1. **{issue}** - {severity}
   - Expert: {which agent found it}
   - Impact: {why it matters}
   - Fix: {recommendation}

### Quick Wins (Fix Today)
1. **{issue}** - Easy fix, high impact
   - Estimated effort: {time}
   - Fix: {one-liner or simple refactor}

## Detailed Analysis

{Insert all expert reports here}

## Trends (if --review mode)

Compare to previous reality check:
- Health score: {before} → {after} ({change})
- New debt: {count}
- Resolved debt: {count}

## Action Items

### This Week
- [ ] Fix critical issues
- [ ] Address quick wins

### This Month
- [ ] Refactor large files
- [ ] Improve test coverage
- [ ] Update documentation

### This Quarter
- [ ] Address technical debt
- [ ] Performance optimization
- [ ] Security hardening
```

### Step 4: Integration with @vision

**Optional:** If PRODUCT_VISION.md or PRD.md exist, compare reality to vision:

```python
# Read vision artifacts
if [ -f "PRODUCT_VISION.md" ]; then
    VISION=$(cat PRODUCT_VISION.md)
fi

if [ -f "docs/prd/PRD.md" ]; then
    PRD=$(cat docs/prd/PRD.md)
fi

# Compare
GAP_ANALYSIS=$(compare_vision_to_reality "$VISION" "$PRD" "$REALITY_REPORT")
```

**Output:**
```markdown
## Vision vs Reality Gap

### Planned vs Implemented
| Feature | PRD Status | Reality Status | Gap |
|---------|------------|----------------|-----|
| Feature 1 | P0 | ✅ Implemented | None |
| Feature 2 | P1 | ⚠️ Partial | Missing X |
| Feature 3 | P0 | ❌ Not found | Not started |

### Architecture Drift
**Planned:** Clean architecture, domain-driven design
**Reality:** Mixed, some layer violations
**Gap:** {description}

### Recommendations
1. Align implementation with PRD priorities
2. Address missing P0 features first
3. Refactor layer violations
```

## Focused Analysis (--focus=topic)

```bash
@reality --focus=security     # Only security expert
@reality --focus=architecture # Only architecture expert
@reality --focus=testing      # Only testing expert
@reality --focus=performance  # Only performance expert
```

**Output:** Single expert report (same format as deep analysis, but only one section)

## Universal Mode

**Works without SDP:** Just run `@reality` in any project directory.

**Auto-detects:**
- Language (Go, Python, Java, Node.js, Ruby)
- Framework (Spring, Express, Django, Gin, etc.)
- Test framework (pytest, jest, go test, JUnit)
- Architecture (layered, DDD, MVC, microservices)

## Examples

```bash
# Quick health check
@reality --quick

# Deep analysis before quarterly planning
@reality --deep

# Focus on security (after vulnerability report)
@reality --focus=security

# Compare to vision (after @vision)
@reality --deep --compare-to-vision
```

## Output Files

**Optional:** Save reports for trend analysis

```bash
@reality --deep --output=docs/reality/2026-02-07-reality-check.md
```

Creates timestamped report for historical tracking.

## Integration with @vision

**Recommended workflow:**

```bash
# Step 1: Strategic planning
@vision "AI-powered task manager"
# → Creates PRODUCT_VISION.md, PRD.md, ROADMAP.md

# Step 2: Reality check (where are we now?)
@reality --quick
# → Analyzes current codebase
# → Identifies gaps to vision

# Step 3: Plan features to bridge gap
@feature "Implement P0 features from PRD"
# → Creates workstreams
```

## Version

**1.0.0** - Initial release with 8 expert agents
