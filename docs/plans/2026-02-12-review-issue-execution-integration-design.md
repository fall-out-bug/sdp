# Review → Issue → Execution Integration Design

> **Status:** Research complete
> **Date:** 2026-02-12
> **Goal:** Design unified system for registering and executing tasks with optional Beads support

---

## Table of Contents

1. [Overview](#overview)
2. [Source of Truth](#1-source-of-truth)
3. [Unified Task Format](#2-unified-task-format)
4. [Execution Skill Interface](#3-execution-skill-interface)
5. [Beads ↔ Workstream Linking](#4-beads--workstream-linking)
6. [Fallback without Beads](#5-fallback-without-beads)
7. [/issue Skill Design](#6-issue-skill-design)
8. [Implementation Plan](#implementation-plan)

---

## Overview

### Goals

1. **Unified task registration** — @review and /issue create tasks consistently
2. **Beads-optional** — Full functionality with or without Beads installed
3. **Executable artifacts** — All tasks can be executed via /build or /bugfix
4. **No data loss** — Transitions between Beads/no-Beads preserve all data

### Key Decisions

| Aspect | Decision |
|--------|----------|
| Source of Truth | Workstream MD files (git-tracked) |
| Beads role | Optional tracking layer, not identity source |
| Task types | Unified format with `type: bug\|task` |
| Execution interface | Auto-detect ID type via unified resolver |
| Fallback | `docs/issues/` directory + index file |

---

## 1. Source of Truth

> **Experts:** Sam Newman (Architecture)

### Decision: Workstream MD Files as Primary

```
docs/workstreams/backlog/*.md  <---- SOURCE OF TRUTH (execution)
         |
         +---> /build (direct read)
         +---> /bugfix (direct read)
         |
         +---> Beads (optional sync for tracking)
         |
         +---> .sdp/log/events.jsonl (evidence)
```

### Rationale

| Concern | Why Workstream MD Wins |
|---------|------------------------|
| **Zero dependencies** | Works without Beads, without network |
| **Human-readable** | MD files are reviewable, mergeable |
| **Git-native** | Versioned alongside code |
| **Self-contained** | File includes all execution context |
| **Graceful degradation** | System works even if external tools fail |

### Beads Role

Beads is a **projection/cache** for team dashboards and metrics, NOT the source of truth.

```bash
# Beads sync is one-way: MD -> Beads
bd sync  # Updates beads from workstream state
```

---

## 2. Unified Task Format

> **Experts:** Theo Browne (API Design)

### Decision: Extend Workstream Format with `type` Field

```yaml
---
ws_id: 99-F063-0001           # Workstream ID (99 prefix = fix task)
feature_id: F063              # Parent feature
type: bug                     # NEW: bug | task | hotfix
priority: 1                   # 0=P0, 1=P1, 2=P2, 3=P3
title: "FIX: CI Go version"
status: backlog
depends_on: []
blocks: []
project_id: sdp
beads_id: sdp-uhsp            # Link to beads (optional)
branch_base: dev              # dev (default) | main (hotfix)
---

## Goal
{What needs to be done}

## Context
{Source: review | issue}

## Acceptance Criteria
- [ ] AC1: ...
- [ ] AC2: ...

## Scope Files
- path/to/file.go

## Notes
- Beads: sdp-uhsp
```

### Type Routing

| Type | Priority | Execution | Branch From |
|------|----------|-----------|-------------|
| `hotfix` | P0 | /hotfix | main |
| `bug` | P1-P2 | /bugfix | dev or feature |
| `task` | P1-P3 | /build | dev or feature |

### Validation Rules

```yaml
type: bug:
  required: [severity, reproduction_steps]

type: task:
  required: [feature_id, acceptance_criteria]

type: hotfix:
  required: [severity]
  priority: 0  # Must be P0
```

---

## 3. Execution Skill Interface

> **Experts:** Martin Fowler (Refactoring)

### Decision: Auto-Detect ID Type with Unified Resolver

```bash
# All execution skills use same interface
/bugfix <id>    # Works with beads ID OR workstream ID
/build <id>     # Works with beads ID OR workstream ID
/hotfix <id>    # Works with beads ID OR workstream ID
```

### ID Detection Logic

```python
def resolve_identifier(id: str) -> dict:
    """
    Pattern Match:
      ^\d{2}-\d{3}-\d{2}$      -> workstream (00-063-01)
      ^\d{2}-[A-Z]\d{3}-\d{4}$ -> workstream fix (99-F063-0001)
      ^[a-z]{3}-[a-z0-9]+$     -> beads ID (sdp-uhsp)
    """
    if re.match(r'^\d{2}-\d{3}-\d{2}$', id):
        return {"type": "workstream", "ws_id": id}
    elif re.match(r'^\d{2}-[A-Z]\d{3}-\d{4}$', id):
        return {"type": "workstream", "ws_id": id}
    elif re.match(r'^[a-z]{3}-[a-z0-9]+$', id):
        # Look up linked workstream
        ws_file = find_workstream_by_beads_id(id)
        return {"type": "beads", "beads_id": id, "ws_file": ws_file}
    else:
        raise ValueError(f"Unknown ID format: {id}")
```

### Resolver Location

```
.claude/
  lib/
    resolver.sh       # ID resolution functions
    resolver.py       # Python implementation (optional)
```

### Example Flow

```
User: /bugfix sdp-uhsp
  |
  v
Resolver: beads ID detected
  |--> Find workstream file with beads_id: sdp-uhsp
  |--> Return: docs/workstreams/backlog/99-F063-0002.md
  v
/bugfix reads workstream, executes fix
```

---

## 4. Beads ↔ Workstream Linking

> **Experts:** Sam Newman (Architecture)

### Decision: Workstream as Master, Bidirectional Reference

```yaml
# In workstream MD file (source of truth)
beads_id: sdp-uhsp            # Reference to beads

# In beads issue (projection)
bd update sdp-uhsp --notes="Workstream: docs/workstreams/backlog/99-F063-0002.md"
```

### Linking Protocol

```
1. Workstream created → If Beads enabled, create beads issue
2. Beads ID written to workstream frontmatter
3. Beads notes include "Workstream: path/to/file.md"
4. On workstream ID change → Update beads notes
5. On beads delete → Workstream keeps beads_id (marked orphaned)
```

### Conflict Resolution

**Workstream always wins.** Beads is a projection, not the source.

```python
def resolve_conflict(workstream, beads_issue):
    if workstream.beads_id != beads_issue.key:
        # Beads might have been recreated
        new_beads = beads.find_by_ws_id(workstream.id)
        if new_beads:
            workstream.beads_id = new_beads.key
        else:
            workstream.beads_id = None  # Orphan
    return workstream
```

---

## 5. Fallback without Beads

> **Experts:** Kelsey Hightower (DevOps)

### Complete Fallback Architecture

```
docs/
├── issues/                        # Beads-free issue storage
│   ├── ISSUE-{SEQ}.md            # Sequential issue files
│   └── ISSUE-0001-auth-error.md
│
├── workstreams/
│   ├── backlog/                  # All tasks (bugs + features)
│   └── completed/

.sdp/
├── log/events.jsonl              # Evidence log (universal)
├── config.yml                    # Project config
└── issues-index.jsonl            # Issue index for fast lookup
```

### Issue File Format (Beads-Free)

```yaml
---
issue_id: ISSUE-0001
title: "Authentication fails"
status: open
priority: 1
type: bug
created_at: 2026-02-12T10:00:00Z
beads_id: null                   # null if no Beads
feature_id: null
---

## Symptom
{User description}

## Classification
- Severity: P1
- Route: /bugfix

## Scope Files
- path/to/file.go
```

### Index File (`.sdp/issues-index.jsonl`)

```json
{"issue_id": "ISSUE-0001", "title": "Auth error", "status": "open", "priority": 1, "file": "docs/issues/ISSUE-0001.md"}
```

### Transition: Beads → No Beads

```bash
# Migrate beads issues to docs/issues/
for issue in $(bd list --format=json); do
  create_issue_file "$issue"
done

# Archive mapping
mv .beads-sdp-mapping.jsonl .sdp/beads-archive.jsonl
```

### Transition: No Beads → Beads

```bash
# Sync issues up to beads
for file in docs/issues/ISSUE-*.md; do
  beads_id=$(bd create --title="..." --type=bug)
  update_frontmatter "$file" "beads_id: $beads_id"
done
```

---

## 6. /issue Skill Design

> **Experts:** Nir Eyal (UX/Product)

### Input Format - Human-First

```bash
# Simple
/issue "Login button doesn't work on mobile"

# With context
/issue "Login broken" --context="Safari iOS 17"

# With evidence
/issue "Crash on upload" --logs=error.log
```

### Auto-Classification Logic

| Severity | Keyword Signals | Route |
|----------|----------------|-------|
| **P0** | "production down", "crash", "blocked", "security" | /hotfix |
| **P1** | "doesn't work", "failing", "error", "broken" | /bugfix |
| **P2** | "edge case", "sometimes", "inconsistently" | backlog |
| **P3** | "cosmetic", "typo", "minor" | defer |

### Workflow

```
1. PARSE INPUT
   └─ Extract keywords, affected components

2. AUTO-CLASSIFY
   └─ Severity rubric + confidence scoring

3. CONFIRM WITH USER (if confidence < 70%)
   └─ "I classified this as P1. Correct?"

4. CREATE ARTIFACT
   ├─ Beads enabled: bd create + workstream MD
   └─ Beads disabled: docs/issues/ISSUE-{SEQ}.md

5. ROUTE
   ├─ P0 → /hotfix (immediate)
   ├─ P1-P2 → /bugfix
   └─ P3 → Defer/Backlog
```

---

## Implementation Plan

### Phase 1: Core Infrastructure

- [ ] Create `docs/issues/` directory structure
- [ ] Create `.sdp/issues-index.jsonl` index file
- [ ] Create `.claude/lib/resolver.sh` ID resolver

### Phase 2: Update @review Skill

- [ ] Always create workstream MD for ALL findings (bugs + tasks)
- [ ] Add `type: bug|task` to workstream frontmatter
- [ ] Add `beads_id` field with bidirectional link
- [ ] Create beads issues if beads enabled

### Phase 3: Update /bugfix Skill

- [ ] Accept workstream ID (99-XXX-XXXX)
- [ ] Accept beads ID (sdp-xxx) → resolve to workstream
- [ ] Read task definition from workstream MD
- [ ] Update both workstream status AND beads status

### Phase 4: Update /issue Skill

- [ ] Human-friendly input parsing
- [ ] Auto-classification logic
- [ ] Dual-track artifact creation
- [ ] User confirmation flow

### Phase 5: Fallback Tools

- [ ] `sdp migrate beads-to-issues` command
- [ ] `sdp migrate issues-to-beads` command
- [ ] Index rebuild command

---

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Tasks work without beads | Partial | 100% |
| ID confusion | High | Zero (unified resolver) |
| Data loss on transition | Possible | Impossible |
| Issue registration time | Manual | Auto-classified |

---

## Summary

The architecture follows these principles:

1. **Workstream MD = Source of Truth** for execution
2. **Beads = Optional tracking layer** for team dashboards
3. **Unified resolver** handles all ID types transparently
4. **Full fallback** via `docs/issues/` when beads unavailable
5. **Bidirectional linking** with workstream as master

This ensures the system works identically with or without Beads, while providing enhanced team features when Beads is available.
