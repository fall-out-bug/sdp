---
name: issue
description: Analyze bugs, classify severity (P0-P3), route to appropriate fix command (/hotfix, /bugfix, or backlog).
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /issue - Analyze & Route Issues

Systematic bug analysis with severity classification and routing.

## Invocation

```bash
/issue "description"
/issue "description" --logs=error.log
```

## Master Prompt

📄 **.claude/skills/issue/SKILL.md** (640+ lines)

**Contains:**
- 5-phase systematic debugging workflow
- Hypothesis formation and testing
- Root cause isolation
- Severity classification (P0/P1/P2/P3)
- Routing rules (hotfix/bugfix/backlog)
- Issue file format
- GitHub integration

## Workflow

1. **Systematic Debugging:**
   - Phase 1: Symptom Documentation
   - Phase 2: Hypothesis Formation (ranked)
   - Phase 3: Systematic Elimination
   - Phase 4: Root Cause Isolation
   - Phase 5: Impact Chain Analysis

2. **Severity Classification:**
   - P0 (CRITICAL): Production down → `/hotfix`
   - P1 (HIGH): Feature broken → `/bugfix`
   - P2 (MEDIUM): Edge case → New WS
   - P3 (LOW): Cosmetic → Defer

3. **Route to Fix**

## Key Output

Issue file: `docs/issues/{ID}-{slug}.md`  
GitHub issue (if gh available)  
Routing recommendation

## Quick Reference

**Input:** Bug description  
**Output:** Issue file + Routing  
**Next:** `/hotfix` or `/bugfix` or schedule WS
