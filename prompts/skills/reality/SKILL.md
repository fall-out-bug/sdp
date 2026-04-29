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

## Workflow

When user invokes `@reality`:

1. Auto-detect project type
2. Run scan based on mode (--quick, --deep, --focus)
3. Spawn expert agents in parallel using Task tool
4. Synthesize report with health score

---

## Step 0: Auto-Detect Project Type

```bash
# Detect language/framework
if [ -f "go.mod" ]; then PROJECT_TYPE="go"
elif [ -f "pyproject.toml" ] || [ -f "requirements.txt" ]; then PROJECT_TYPE="python"
elif [ -f "pom.xml" ] || [ -f "build.gradle" ]; then PROJECT_TYPE="java"
elif [ -f "package.json" ]; then PROJECT_TYPE="nodejs"
else PROJECT_TYPE="unknown"
fi
```

## Step 1: Quick Scan (--quick mode)

**Analysis:**
1. Project size (lines of code, file count)
2. Architecture (layer violations, circular dependencies)
3. Test coverage (if tests exist, estimate %)
4. Documentation (doc coverage, drift detection)
5. Quick smell check (TODO/FIXME/HACK comments, long files)

**Output:** Health Score X/100 + Top 5 Issues

## Step 2: Deep Analysis (--deep mode)

Spawn 8 parallel expert analyses using Task tool with subagent_type:

```
Task(subagent_type="general-purpose", prompt="Analyze ARCHITECTURE...")
Task(subagent_type="general-purpose", prompt="Analyze CODE QUALITY...")
Task(subagent_type="general-purpose", prompt="Analyze TESTING...")
Task(subagent_type="general-purpose", prompt="Analyze SECURITY...")
Task(subagent_type="general-purpose", prompt="Analyze PERFORMANCE...")
Task(subagent_type="general-purpose", prompt="Analyze DOCUMENTATION...")
Task(subagent_type="general-purpose", prompt="Analyze TECHNICAL DEBT...")
Task(subagent_type="general-purpose", prompt="Analyze STANDARDS...")
```

Expert agents:
1. ARCHITECTURE expert - Layer mapping, dependencies, violations
2. CODE QUALITY expert - File size, complexity, duplication
3. TESTING expert - Coverage, test quality, frameworks
4. SECURITY expert - Secrets, OWASP, dependencies
5. PERFORMANCE expert - Bottlenecks, caching, scalability
6. DOCUMENTATION expert - Coverage, drift, quality
7. TECHNICAL DEBT expert - TODO/FIXME, code smells
8. STANDARDS expert - Conventions, error handling, types

## Step 3: Synthesize Report

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
| `@reality --quick` | 5-10 min | Health check + top issues |
| `@reality --deep` | 30-60 min | Comprehensive with 8 experts |
| `@reality --focus=security` | Varies | Security expert deep dive |
| `@reality --focus=architecture` | Varies | Architecture expert deep dive |
| `@reality --focus=testing` | Varies | Testing expert deep dive |
| `@reality --focus=performance` | Varies | Performance expert deep dive |

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

If PRODUCT_VISION.md exists, compare reality to vision with Vision vs Reality Gap analysis:

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

## Write Plan (F101)

Before writing the reality-check report, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (report file, event log).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"reality"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @reality <project>:
  CREATE: docs/reality/check.md — Reality-check report with health score
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

- `@vision` - Strategic planning
- `@feature` - Feature planning
- `@idea` - Requirements gathering
