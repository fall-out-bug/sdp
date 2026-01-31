# SDP Overview for Team Leads

**High-level introduction for technical leaders evaluating SDP.**

---

## What is SDP?

SDP (Spec-Driven Protocol) is a **workstream-driven framework** that transforms AI coding tools (Claude Code, Cursor, OpenCode) into a **structured software development process**.

Think of it as: **TDD + CI/CD + AI Agents** combined into one coherent system.

---

## Key Benefits

### 1. Structure for AI-Generated Code

**Problem:** AI agents generate code quickly, but without process it becomes unmanageable.

**SDP Solution:**
- Break features into atomic **workstreams** (500-1500 LOC each)
- Execute one workstream per AI session
- Track progress with task systems (Beads CLI)
- Quality gates enforced on every workstream

**Result:** Predictable, manageable AI-powered development.

---

### 2. Multi-Agent Coordination

**Problem:** AI agents need specialized roles (planner, builder, reviewer).

**SDP Solution:**
```python
# Spawn specialized agents
spawner.spawn_agent(AgentConfig(
    name="planner",
    prompt="Break features into workstreams...",
))

# Route messages between agents
router.send_message(Message(
    sender="orchestrator",
    content="Plan feature F24",
    recipient=planner_id,
))
```

**Result:** Orchestrate multiple AI agents automatically.

---

### 3. Progress Tracking

**Problem:** Lost track of what's done, what's blocked.

**SDP Solution:**
- **Beads CLI** integration for task tracking
- Hash-based task IDs (bd-0001, bd-0001.1)
- Dependency DAG (WS-02 blocked by WS-01)
- Ready detection (`bd ready` shows what to work on next)

**Result:** Clear visibility into project state.

---

### 4. Quality Gates Enforced

**Problem:** AI code can be buggy, untested, or hard to maintain.

**SDP Solution:**
- **TDD required** - Tests first, code second
- **Coverage ≥80%** - Enforced on all files
- **Type hints** - Full mypy --strict compliance
- **Linting** - ruff for code quality
- **File size <200 LOC** - Keep code focused

**Result:** Consistent code quality across all AI-generated code.

---

## When to Use SDP

### ✅ Good Fit

- **Solo developers** using Claude Code/Cursor/OpenCode
- **Small teams** (1-5 developers) with AI assistants
- **Projects** with 5-500 workstreams
- **Need** for structured process and quality gates

### ❌ Not Ideal

- **Large teams** (>10 developers) - might need heavier process
- **Simple scripts** - overhead not worth it
- **Legacy codebases** - requires refactoring to fit SDP

---

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         SDP Orchestrator               │
│  ┌──────────┐  ┌──────────┐  ┌──────┐  │
│  │ Planner  │  │ Builder  │  │Reviewer│  │
│  └──────────┘  └──────────┘  └──────┘  │
│         │              │         │        │
│         ▼              ▼         ▼        │
│  ┌─────────────────────────────────┐ │
│  │      Beads Task Tracker          │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

### Workflow

```
1. @feature "Add user authentication"
   ↓ (interviews you)
2. @design beads-auth
   ↓ (plans workstreams)
3. @build WS-AUTH-01
   ↓ (TDD: Red → Green → Refactor)
4. @review beads-auth
   ↓ (quality check)
5. @deploy beads-auth
   ↓ (production)
```

---

## Metrics & ROI

### Time to First Value
- **Installation:** 5 minutes (pip install)
- **First feature:** 30 minutes (@feature → @oneshot)
- **ROI:** Immediate if you already use AI coding tools

### Quality Impact
- **Coverage:** 91% (enforced)
- **Type safety:** 100% (mypy --strict)
- **Code quality:** ruff enforced
- **File size:** <200 LOC (maintainable)

### Productivity
- **Workstreams:** Executed in one AI session
- **Checkpoints:** Resume after interruption
- **Parallel execution:** Independent WS can run in parallel

---

## Comparison

### Without SDP
```
❌ AI generates code quickly → becomes unmanageable
❌ No clear process → chaotic PRs
❌ Quality varies → depends on prompt
❌ Lost progress → what's done, what's next?
```

### With SDP
```
✅ Structured workstreams → manageable chunks
✅ Quality gates enforced → consistent quality
✅ Progress tracking → clear visibility
✅ Multi-agent coordination → specialized expertise
✅ Checkpoint system → fault tolerance
```

---

## Getting Started

### For Your Team

1. **Try it locally:**
   ```bash
   pipx install sdp-cli
   @feature "Add user comments"
   ```

2. **Evaluate for 1 week:**
   - Use SDP for 1-2 features
   - Compare quality/speed vs current process
   - Get team feedback

3. **Adopt if:**
   - Quality improves
   - Team likes structure
   - AI tools become more productive

---

## Common Questions

**Q: Does this replace CI/CD?**
A: No, SDP complements CI/CD. SDP is for AI-assisted development, CI/CD is for automation.

**Q: Can we use our own agents?**
A: Yes, SDP is agent-agnostic. Create custom roles in `.claude/agents/`.

**Q: What about existing code?**
A: SDP works best for new features. Existing code needs refactoring first.

**Q: Is this only for solo developers?**
A: Designed for 1-5 person teams. Larger teams might need heavier process.

---

## Next Steps

1. **Read the tutorial:** [docs/TUTORIAL.md](docs/TUTORIAL.md) (15 min)
2. **Try the example:** `@feature "test feature"`
3. **Decide:** Is SDP right for your team?

**For technical details:** [PROTOCOL.md](PROTOCOL.md)

---

*SDP v0.5.0 - Unified Workflow*
*Updated: 2026-01-29*
