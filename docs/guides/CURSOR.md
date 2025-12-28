# Cursor IDE 2.0 Integration Guide

**Updated for Cursor 2.0** (Released October 29, 2025)

This guide explains how to use the Consensus Workflow with [Cursor 2.0](https://cursor.com), leveraging its revolutionary multi-agent architecture, Composer model, and parallel execution capabilities.

## 🆕 What's New in Cursor 2.0

Cursor 2.0 represents a fundamental shift from traditional AI coding assistants to a **true multi-agent system** - perfect for Consensus Workflow!

### Key Features

**Composer** - Cursor's proprietary coding model
- 4x faster than similarly intelligent models
- Built for agentic coding with RL techniques
- Completes most turns in <30 seconds
- Optimized for low-latency iterations

**Multi-Agent Mode** (Up to 8 Agents)
- Run agents in parallel on single prompt
- Git worktrees for conflict prevention
- Automatic best-solution selection
- Agent recommendation with explanations

**Agent-Centric Interface**
- Agents, plans, and runs as first-class objects
- Dedicated sidebar for agent coordination
- Four layouts: agent, editor, zen, browser
- Keyboard shortcuts for layout switching

**Auto Context Gathering**
- Agents self-gather context without manual @-mentions
- Semantic codebase search built-in
- Structured feedback loops

**Sandboxed Terminals** (GA on macOS)
- Secure command execution
- Read/write access to workspace only
- No internet access by default

**Embedded Browser**
- In-editor web browsing
- DOM element selection
- Forward information to agents

## 🚀 Quick Start

### Prerequisites

1. **Cursor IDE 2.0+** installed
2. **API keys configured** for your chosen providers
3. **Repository opened** in Cursor

### Basic Setup

1. Open Settings (Cmd/Ctrl + ,)
2. Go to "Models" tab
3. Configure your preferred providers:

```json
{
  "models.default": "gemini-3.0-flash",
  "models.providers": {
    "anthropic": {
      "apiKey": "sk-ant-...",
      "enabled": true
    },
    "google": {
      "apiKey": "...",
      "enabled": true
    },
    "openai": {
      "apiKey": "sk-...",
      "enabled": true
    }
  }
}
```

## 🎯 Recommended Model Configuration

Based on [MODELS.md](../../MODELS.md) SWE-bench Verified data (December 2025):

### Workspace Settings (.cursor/settings.json)

```json
{
  "chat.defaultModel": "gemini-3.0-flash",

  "chat.modelsByRole": {
    "analyst": "claude-opus-4-5-20251101",
    "architect": "claude-opus-4-5-20251101",
    "tech_lead": "gemini-3.0-flash",
    "developer": "gemini-3.0-flash",
    "qa": "gemini-3.0-flash",
    "devops": "qwen3-coder-480b",
    "sre": "gemini-3.0-flash",
    "security": "claude-opus-4-5-20251101"
  },

  "composer.enabled": true,
  "composer.multiAgent": true,
  "composer.maxAgents": 6,

  "terminal.sandboxed": true,

  "chat.contextFiles": [
    "CLAUDE.md",
    "RULES_COMMON.md"
  ],

  "indexing.include": [
    "docs/specs/**",
    "prompts/**",
    "src/**"
  ]
}
```

### Model Selection Matrix

| Role | Primary (Best) | Alternative (Budget) | Open Source (Free) |
|------|----------------|---------------------|-------------------|
| **Analyst** | Opus 4.5 (80.9%) | Gemini 3 Pro (74.2%) | Kimi K2 Thinking (71.3%) |
| **Architect** | Opus 4.5 (80.9%) | Gemini 3 Pro (74.2%) | Kimi K2 Thinking (71.3%) |
| **Tech Lead** | **Gemini 3 Flash** (76-78%) ⭐ | Sonnet 4.5 (77.2%) | Qwen3-Coder (69.6%) |
| **Developer** | **Gemini 3 Flash** (76-78%) ⭐ | Composer / Haiku 4.5 | Kimi K2 Thinking (71.3%) |
| **QA** | **Gemini 3 Flash** (76-78%) ⭐ | Composer / Haiku 4.5 | Qwen3-Coder (69.6%) |
| **DevOps** | **Gemini 3 Flash** (76-78%) ⭐ | Composer | Qwen3-Coder (69.6%) |
| **Security** | Opus 4.5 (80.9%) | GPT-5.2 (71.8%) | Kimi K2 Thinking (71.3%) |

**💡 Why Gemini 3 Flash as default:**
- 76-78% on SWE-bench (beats Haiku 4.5's 73.3%)
- **13x cheaper** than Haiku ($0.075/$0.30 vs $1/$5)
- 4-5x faster than Sonnet 4.5
- Perfect for rapid iteration

**🤖 Cursor's Composer:**
- Use for Developer/QA when you want Cursor-native experience
- 4x faster than similar models
- Built for low-latency coding
- No external API costs

## 🎨 Multi-Agent Workflow (Cursor 2.0)

### Method 1: Parallel Agent Execution (RECOMMENDED)

Cursor 2.0's killer feature - run multiple agents simultaneously!

#### Setup: Epic with Multiple Workstreams

```
Epic: User Authentication API

Workstreams:
1. Login endpoint + JWT
2. Registration endpoint + validation
3. Password reset flow
4. Session management
```

#### Running Multi-Agent Mode

1. **Open Composer** (Cmd/Ctrl + I)

2. **Enable Multi-Agent Mode** (click "Multi-Agent" toggle)

3. **Configure Agents:**
   ```
   Agent 1 (Developer): Implement login endpoint
   Agent 2 (Developer): Implement registration endpoint
   Agent 3 (Developer): Implement password reset
   Agent 4 (QA): Write integration tests for all endpoints
   Agent 5 (DevOps): Create deployment config
   Agent 6 (Security): Review auth implementation
   ```

4. **Set Prompt:**
   ```
   @prompts/developer_prompt.md
   @docs/specs/epic_XX/implementation.md

   You are working on the User Authentication API epic.

   Agent 1-3: Each implement ONE endpoint with TDD
   Agent 4: Write integration tests covering all endpoints
   Agent 5: Create Docker + CI/CD config
   Agent 6: Security review of all auth code

   Use Gemini 3 Flash for speed.
   Follow Clean Architecture boundaries.
   No duplications - search codebase first!
   ```

5. **Run** - Cursor will:
   - Spawn 6 agents in parallel
   - Use git worktrees to prevent conflicts
   - Execute agents simultaneously
   - Evaluate all solutions
   - Recommend best approach per agent
   - Show diff comparison

6. **Review Results:**
   - Check Agent sidebar for completion status
   - Review recommended solutions (marked with ✓)
   - Read explanations for why each was chosen
   - Merge selected solutions

**⚡ Time Savings:**
- Sequential: Agents run one at a time with context switching overhead
- Parallel (Cursor 2.0): Multiple agents work simultaneously on independent workstreams

### Method 2: Sequential Agent Windows

Traditional approach for complex epics requiring human checkpoints:

#### Window Layout

```
┌──────────────────┬──────────────────┐
│  Analyst         │  Architect       │
│  (Opus 4.5)      │  (Opus 4.5)      │
├──────────────────┼──────────────────┤
│  Tech Lead       │  Developer       │
│  (Flash)         │  (Flash/Composer)│
├──────────────────┼──────────────────┤
│  QA              │  DevOps          │
│  (Flash)         │  (Flash/Qwen3)   │
└──────────────────┴──────────────────┘
```

#### Opening Multiple Windows

```bash
# From terminal
cursor docs/specs/epic_XX/ --new-window  # Analyst
cursor docs/specs/epic_XX/ --new-window  # Architect
cursor docs/specs/epic_XX/ --new-window  # Tech Lead
cursor docs/specs/epic_XX/ --new-window  # Developer
cursor docs/specs/epic_XX/ --new-window  # QA
cursor docs/specs/epic_XX/ --new-window  # DevOps
```

Or use Cursor's built-in window management (View → New Window).

## 🤖 Running Each Agent

### Analyst (Strategic - Opus 4.5)

**In Analyst window chat:**
```
Model: claude-opus-4-5-20251101

@prompts/analyst_prompt.md
@docs/specs/epic_XX/epic.md

You are the Analyst agent. Analyze this epic and create requirements.

Output:
1. docs/specs/epic_XX/consensus/artifacts/requirements.json
2. Message to architect's inbox
3. Decision log entry

Critical: ALL text must be in English.
```

**What Opus 4.5 does well:**
- Identifying edge cases and ambiguities
- Formulating testable acceptance criteria
- Detecting scope creep
- Understanding stakeholder intent

### Architect (Strategic - Opus 4.5)

**In Architect window chat:**
```
Model: claude-opus-4-5-20251101

@prompts/architect_prompt.md
@docs/specs/epic_XX/consensus/artifacts/requirements.json
@docs/specs/epic_XX/consensus/messages/inbox/architect/

You are the Architect agent. Review requirements and create architecture.

VETO if Clean Architecture violations detected:
- Domain depends on Infrastructure ❌
- Presentation accesses Domain directly ❌
- Missing Port/Adapter abstractions ❌

Output architecture.json with layer boundaries clearly defined.
```

**Veto Example:**
If Analyst specified "API calls database directly" → VETO with explanation.

### Tech Lead (Planning - Gemini 3 Flash)

**In Tech Lead window chat:**
```
Model: gemini-3.0-flash

@prompts/tech_lead_prompt.md
@docs/specs/epic_XX/consensus/artifacts/architecture.json

You are the Tech Lead. Create implementation plan.

Before planning:
1. Search codebase for similar implementations (DRY)
2. Check for existing abstractions to reuse
3. Identify cross-epic dependencies

Break into workstreams (max 150 LOC each).
Define testing strategy per workstream.
```

**Why Gemini 3 Flash here:**
- 76-78% on SWE-bench (sufficient for planning)
- 1-2s latency (instant feedback)
- 13x cheaper than alternatives
- Can iterate quickly on plan

### Developer (Implementation - Gemini 3 Flash or Composer) ⭐

**In Developer window chat:**

#### Option A: Gemini 3 Flash (Recommended)
```
Model: gemini-3.0-flash

@prompts/quick/developer_quick.md
@docs/specs/epic_XX/implementation.md
@docs/specs/epic_XX/consensus/messages/inbox/developer/

You are the Developer. Implement Workstream 1 with TDD.

Steps:
1. Search codebase for existing implementations
2. Write failing tests first (RED)
3. Implement minimal code to pass (GREEN)
4. Refactor if needed (REFACTOR)
5. Verify ≥80% coverage
6. Send message to QA inbox

Use Gemini 3 Flash for speed (1-2s response time).
```

#### Option B: Cursor Composer (Alternative)
```
Model: composer

@prompts/quick/developer_quick.md
@docs/specs/epic_XX/implementation.md

Enable Composer's agentic mode for multi-step implementation:
1. Read implementation.md
2. Search for duplications automatically
3. Write tests
4. Implement code
5. Run tests in sandboxed terminal
6. Iterate until passing

Composer will handle the full TDD cycle autonomously.
```

**Multi-Agent Developer Mode:**
For complex features, use Cursor 2.0 multi-agent:
```
Launch 3 Developer agents in parallel:
- Agent 1: Implement domain logic
- Agent 2: Implement infrastructure adapters
- Agent 3: Implement presentation layer

Cursor will coordinate via git worktrees.
```

### QA (Verification - Gemini 3 Flash)

**In QA window chat:**
```
Model: gemini-3.0-flash

@prompts/quick/qa_quick.md
@docs/specs/epic_XX/testing.md
@docs/specs/epic_XX/consensus/messages/inbox/qa/

You are the QA agent. Verify implementation quality.

Checklist:
1. All acceptance criteria met ✓
2. Test coverage ≥80% ✓
3. No regressions ✓
4. Performance acceptable ✓
5. Error handling complete ✓

Run tests in sandboxed terminal:
```
/workspace pytest tests/ --cov=src --cov-report=term
```

Output: test_results.md + message to DevOps
```

**Sandboxed Terminal:**
Cursor 2.0's sandboxed terminal ensures QA can run tests safely without affecting system.

### DevOps (Deployment - Gemini 3 Flash or Qwen3-Coder)

**In DevOps window chat:**
```
Model: gemini-3.0-flash  # or qwen3-coder-480b for free

@prompts/quick/devops_quick.md
@docs/specs/epic_XX/deployment.md

You are the DevOps agent. Create deployment plan.

Deliverables:
1. Dockerfile (if needed)
2. CI/CD pipeline config (.github/workflows/)
3. Environment variables documentation
4. Rollback procedure

Test pipeline in sandboxed terminal:
```
/workspace docker build -t app:test .
/workspace docker run --rm app:test pytest
```
```

**Why Qwen3-Coder option:**
- FREE (self-hosted via Ollama)
- 69.6% on SWE-bench (sufficient for scripts)
- Good at infrastructure code
- Budget-friendly for non-critical path

## 🎛️ Cursor 2.0 Interface Features

### Agent Layout (New in 2.0)

Switch to Agent layout: `Cmd/Ctrl + Shift + A`

**Sidebar shows:**
- Active agents (running/completed)
- Agent plans (step-by-step breakdown)
- Agent runs (execution history)
- Recommendations (best solutions)

**Main panel:**
- Agent conversation
- Code diffs
- File changes

### Composer Features

**Plan vs Build Separation:**
```
1. Create plan with strategic model:
   Model: claude-opus-4-5-20251101
   Prompt: "Design authentication system architecture"
   → Generates detailed plan

2. Build with implementation model:
   Model: gemini-3.0-flash
   Prompt: "Execute plan in foreground"
   → Implements code rapidly

Result: Best of both worlds (strategic + fast)
```

**Background vs Foreground Building:**
- Foreground: See progress live, can interrupt
- Background: Continue working, check result later
- For Consensus: Use foreground for critical path

### Auto Context Gathering

**Old way (manual):**
```
@file1.ts @file2.ts @docs/spec.md
Implement feature X
```

**New way (automatic):**
```
Implement feature X for the authentication epic

→ Cursor automatically finds:
  - docs/specs/epic_XX/epic.md
  - src/auth/*.ts
  - tests/auth/*.test.ts
  - Related dependencies
```

**For Consensus:**
Still explicitly reference prompts and protocol:
```
@prompts/developer_prompt.md
@RULES_COMMON.md

Implement Workstream 2

→ Agent knows to follow Consensus protocol
```

### Sandboxed Terminal

**Enable in settings:**
```json
{
  "terminal.sandboxed": true,
  "terminal.allowedCommands": [
    "pytest",
    "npm",
    "docker",
    "git"
  ]
}
```

**Security benefits for Consensus:**
- QA runs tests safely
- DevOps tests pipelines without side effects
- No accidental production deployments
- Logs all command execution

### Embedded Browser

**Use cases:**
1. **Documentation lookup** (in-editor)
2. **API testing** (see responses live)
3. **Visual verification** (UI changes)

**For Consensus:**
```
Developer agent:
1. Implements API endpoint
2. Uses embedded browser to test:
   - POST /api/login
   - See response in browser panel
   - Forward to chat: "Fix: status should be 201, not 200"
3. Agent corrects code
4. Retest immediately
```

## 🚄 Performance Optimization

### Speed Comparison (per agent task)

| Model | Setup | Execute | Total | Cost |
|-------|-------|---------|-------|------|
| **Gemini 3 Flash** | 1s | 2-3s | **3-4s** ⭐ | **$0.02** |
| Composer | 0.5s | 3-5s | 3.5-5.5s | $0 (included) |
| Haiku 4.5 | 1s | 3-4s | 4-5s | $0.25 |
| Sonnet 4.5 | 2s | 8-12s | 10-14s | $0.75 |
| Opus 4.5 | 3s | 25-35s | 28-38s | $3.50 |

**Epic Workflow (6 agents):**
```
Strategy 1: All Gemini 3 Flash ⭐ RECOMMENDED
- Speed: Very fast (Flash latency: 1-3s per agent)
- Cost: Very low ($0.075/$0.30 per 1M tokens)
- Quality: 76-78% (SWE-bench verified)

Strategy 2: Mixed (Opus strategic, Flash implementation)
- Speed: Medium (Opus: 25-35s, Flash: 1-3s)
- Cost: Medium (strategic premium model usage)
- Quality: 79% balanced

Strategy 3: All Opus 4.5
- Speed: Slow (Opus latency: 25-35s per agent)
- Cost: High ($15/$75 per 1M tokens)
- Quality: 80.9% (diminishing returns)
```

### Parallel Execution Speedup

**Sequential (traditional):**
```
Analyst    → Architect → Tech Lead → Developer (×3) → QA → DevOps
Each agent waits for previous to complete
High context switching overhead
```

**Parallel (Cursor 2.0 multi-agent):**
```
Phase 1: Analyst + Architect (parallel, independent)
Phase 2: Tech Lead (depends on Phase 1)
Phase 3: Developer×3 + QA + DevOps (parallel workstreams)

Significantly faster through parallel execution
```

## 💰 Cost Optimization Strategies

### 1. Default to Gemini 3 Flash

```json
{
  "chat.defaultModel": "gemini-3.0-flash",
  "chat.escalationRules": {
    "iterations": {
      "threshold": 2,
      "escalateTo": "claude-sonnet-4-5-20250929"
    },
    "complexity": {
      "multiFile": 5,
      "escalateTo": "claude-sonnet-4-5-20250929"
    },
    "strategic": {
      "roles": ["analyst", "architect", "security"],
      "useModel": "claude-opus-4-5-20251101"
    }
  }
}
```

**Result:**
- 80% of tasks use Flash ($0.075/$0.30)
- Auto-escalate when needed
- Strategic roles always get Opus

### 2. Use Composer for Routine Tasks

Cursor's Composer is **free** (included in subscription):
- Developer: Routine implementations
- QA: Test generation
- DevOps: Simple scripts

**When to escalate from Composer:**
- Complex business logic (→ Gemini 3 Flash)
- Architectural decisions (→ Opus 4.5)
- Security-critical (→ Opus 4.5)

### 3. Open Source for Non-Critical

```json
{
  "chat.modelsByRole": {
    "devops": "qwen3-coder-480b",
    "documentation": "qwen3-coder-480b"
  }
}
```

**Setup Qwen3-Coder locally:**
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull Qwen3-Coder
ollama pull qwen2.5-coder:32b

# Configure Cursor to use local model
# Settings → Models → Add Local Model
# URL: http://localhost:11434
# Model: qwen2.5-coder:32b
```

**Savings:**
- DevOps: FREE vs $0.02/task
- Documentation: FREE vs $0.01/task
- 50+ tasks/month → Save $50-100

### 4. Batch Operations in Multi-Agent

Instead of running 6 separate chat sessions:
```
Old way: 6 sessions × $0.02 = $0.12
New way: 1 multi-agent run = $0.02 (shared context)

10 epics/month: Save $1.00 (83% reduction)
```

## 🎯 Best Practices

### 1. Start with Gemini 3 Flash, Escalate if Needed

**Workflow:**
```
Task: Implement user authentication

Try 1: Gemini 3 Flash (3s, $0.02)
  → Implements basic flow
  → Missing edge case: concurrent logins

Try 2: Gemini 3 Flash with clarification (3s, $0.02)
  → Adds mutex for concurrent safety
  → ✓ Passes all tests

Total: 6s, $0.04
Alternative (Opus from start): 30s, $3.50
```

Only escalate if:
- 2+ iterations without progress
- Complex trade-offs emerge (use Sonnet 4.5)
- Architectural ambiguity (use Opus 4.5)

### 2. Use Multi-Agent for Independent Workstreams

**Good for multi-agent:**
```
Epic: REST API with 4 endpoints

Agent 1: GET /users
Agent 2: POST /users
Agent 3: PUT /users/:id
Agent 4: DELETE /users/:id

→ Perfect! All independent, can run parallel
```

**Bad for multi-agent:**
```
Epic: Implement MVC pattern

Agent 1: Model layer
Agent 2: View layer (depends on Model)
Agent 3: Controller (depends on both)

→ Dependencies! Run sequentially instead
```

### 3. Leverage Sandboxed Terminal for Testing

```
Developer agent:
1. Writes code
2. Writes tests
3. Runs in sandbox:
   ```
   /workspace pytest tests/test_auth.py -v
   ```
4. If failures → fix → rerun
5. If pass → send to QA

QA agent:
1. Runs full suite in sandbox
2. Checks coverage:
   ```
   /workspace pytest --cov=src --cov-report=html
   ```
3. Reviews htmlcov/index.html in embedded browser
4. Approves or vetoes
```

**Safety:** Sandboxed terminal prevents:
- Accidental deployments
- File system corruption
- Network requests (unless explicitly allowed)

### 4. Use Composer for Plan, Flash for Build

```
Step 1: Planning (Opus 4.5)
  "Design authentication system with OAuth, JWT, refresh tokens"
  → Detailed architecture plan generated

Step 2: Build (Gemini 3 Flash in background)
  "Execute plan with TDD, run tests automatically"
  → Rapid implementation while you work on other tasks
  → Get notification when complete

Step 3: Review
  → Check generated code
  → Run in sandboxed terminal
  → Approve or iterate
```

### 5. Keep Context Focused Per Agent

**Good:**
```
Developer window:
  @prompts/developer_prompt.md
  @docs/specs/epic_XX/implementation.md
  @src/auth/  (only relevant code)

QA window:
  @prompts/qa_prompt.md
  @tests/
  @docs/specs/epic_XX/testing.md

→ Each agent sees only what it needs
```

**Bad:**
```
All agents:
  @docs/specs/epic_XX/**/*
  @src/**/*
  @tests/**/*

→ Too much context, slower, expensive
```

### 6. Verify Consensus Messages

Before sending messages between agents:
```
Developer → QA handoff:

In Developer chat:
"Create message for QA's inbox:
- Workstream 1 complete
- Tests passing
- Coverage verified
- File: src/auth/login.ts

Verify message format:
- JSON with compact keys (d, st, r, epic, sm, nx)
- All text in English
- No emojis (unless user requested)
"

Save to: docs/specs/epic_XX/consensus/messages/inbox/qa/2025-12-29-workstream1-ready.json
```

## 🔧 Advanced Configuration

### Custom Rules per Role

`.cursor/rules/developer.mdc`:
```markdown
---
description: Developer agent rules for Consensus workflow
globs: ["src/**", "tests/**"]
---

# Developer Agent Rules

## TDD Cycle
1. Write failing test (RED)
2. Implement minimal code (GREEN)
3. Refactor (REFACTOR)

## Code Quality
- Functions ≤15 LOC when practical
- No silent fallbacks (`except: pass` forbidden)
- DRY: Search before implementing
- SOLID principles

## Before Implementation
Run: grep -r "similar_function_name" src/
Check: Does this already exist?

## After Implementation
Run in sandboxed terminal:
```bash
pytest tests/ -v --cov=src
```

Coverage must be ≥80% in touched areas.

## Handoff to QA
Create message: docs/specs/{epic}/consensus/messages/inbox/qa/{date}-{subject}.json
Format: JSON with compact keys (d, st, r, epic, sm, nx, artifacts)
```

### Model Switching Automation

`.cursor/settings.json`:
```json
{
  "chat.autoSwitch": true,
  "chat.switchRules": [
    {
      "pattern": "architect|security|veto",
      "model": "claude-opus-4-5-20251101",
      "reason": "Strategic decision required"
    },
    {
      "pattern": "refactor.*\\d{5,}",
      "model": "claude-sonnet-4-5-20250929",
      "reason": "Large refactoring (5+ files)"
    },
    {
      "pattern": "implement|fix|test",
      "model": "gemini-3.0-flash",
      "reason": "Standard implementation task"
    }
  ]
}
```

### Integration with External Tools

**GitHub Copilot + Cursor:**
```json
{
  "github.copilot.enable": true,
  "github.copilot.useCursorComposer": true
}
```

**Use Copilot for:**
- Inline completions (fast, free with subscription)
- Boilerplate generation

**Use Cursor Composer for:**
- Multi-step refactoring
- Agent-based implementation
- Complex logic

**Use Gemini 3 Flash for:**
- TDD cycles
- Test generation
- Bug fixing

## 🐛 Troubleshooting

### Agent Not Following Protocol

**Symptom:** Agent creates messages in wrong format

**Solution:**
```
Add to chat:

CRITICAL RULES:
1. ALL messages must be in English
2. JSON format with compact keys: d, st, r, epic, sm, nx, artifacts
3. Save to: docs/specs/{epic}/consensus/messages/inbox/{target_agent}/

@RULES_COMMON.md

Verify message before saving. Show me the JSON first.
```

### Multi-Agent Conflicts

**Symptom:** Agents modify same files, git conflicts

**Solution 1:** Use Cursor's automatic git worktrees
```json
{
  "composer.multiAgent": true,
  "composer.useWorktrees": true
}
```

**Solution 2:** Partition by file boundaries
```
Agent 1: src/auth/login.ts
Agent 2: src/auth/register.ts
Agent 3: src/auth/reset.ts

→ No conflicts possible
```

### Model Not Available

**Symptom:** "Model gemini-3.0-flash not found"

**Solution:**
1. Check API key: Settings → Models → Google
2. Verify model name (not `gemini-3-flash`, use `gemini-3.0-flash`)
3. Check regional availability
4. Fallback to Haiku 4.5:
   ```json
   {
     "chat.fallbackModel": "claude-haiku-4-5-20241022"
   }
   ```

### Slow Performance

**Symptom:** Agents taking >30s per task

**Check:**
1. Using right model? (Should be Flash for most tasks)
2. Context too large? (Remove unnecessary @-mentions)
3. Multi-agent disabled? (Enable for parallel speed)

**Fix:**
```json
{
  "chat.defaultModel": "gemini-3.0-flash",  // Fast
  "chat.maxContextFiles": 10,                // Limit context
  "composer.multiAgent": true                 // Enable parallel
}
```

### High Costs

**Symptom:** >$100/month on API costs

**Audit:**
```bash
# Check settings
cat .cursor/settings.json | jq '.chat.defaultModel'

# Should show: "gemini-3.0-flash" (not "claude-opus-4-5-20251101")
```

**Optimize:**
```json
{
  "chat.defaultModel": "gemini-3.0-flash",
  "chat.modelsByRole": {
    "analyst": "claude-opus-4-5-20251101",    // Only strategic
    "architect": "claude-opus-4-5-20251101",  // Only strategic
    "developer": "gemini-3.0-flash",          // Everything else
    "qa": "gemini-3.0-flash",
    "devops": "qwen3-coder-480b"              // Free!
  }
}
```

**Result:** Significant cost reduction through strategic model selection

## 📊 Example Workflow

### Complete Epic: User Authentication API

This example shows the recommended workflow structure. Actual timing and costs vary by epic complexity.

**Epic Setup:**
```
docs/specs/epic_12_auth_api/epic.md
Goal: Implement OAuth2 + JWT authentication
Workstreams: 4 (login, register, refresh, logout)
```

**Recommended Agent Sequence:**

1. **Analyst (Opus 4.5):**
   ```
   Model: claude-opus-4-5-20251101
   @prompts/analyst_prompt.md
   @docs/specs/epic_12_auth_api/epic.md

   Analyze and create requirements.json
   ```

2. **Architect (Opus 4.5):**
   ```
   Model: claude-opus-4-5-20251101
   @prompts/architect_prompt.md
   @requirements.json

   Design Clean Architecture with Port/Adapter pattern
   ```

3. **Tech Lead (Gemini 3 Flash):**
   ```
   Model: gemini-3.0-flash
   @prompts/tech_lead_prompt.md
   @architecture.json

   Break into workstreams, define test strategy
   ```

4. **Developer (Gemini 3 Flash - Multi-Agent Mode):**
   ```
   Cursor Multi-Agent Mode (parallel agents)
   Model: gemini-3.0-flash

   Agent 1: Implement login endpoint with JWT
   Agent 2: Implement register endpoint with validation
   Agent 3: Implement refresh token rotation
   Agent 4: Implement logout with token invalidation
   ```

5. **QA (Gemini 3 Flash):**
   ```
   Model: gemini-3.0-flash
   @prompts/qa_prompt.md

   Run integration tests, verify coverage
   ```

6. **DevOps (Qwen3-Coder via Ollama - FREE):**
   ```
   Model: qwen3-coder-480b
   @prompts/devops_prompt.md

   Create Dockerfile + GitHub Actions CI/CD
   ```

**Benefits of Multi-Agent Parallel Execution:**
- Faster iteration (parallel vs sequential)
- Lower cost (Gemini 3 Flash for implementation vs premium models)
- Consistent quality (structured prompts ensure quality gates)

## 🎓 Summary Recommendations

### For Most Teams (Optimal) ⭐

```json
{
  "chat.defaultModel": "gemini-3.0-flash",
  "chat.modelsByRole": {
    "analyst": "claude-opus-4-5-20251101",
    "architect": "claude-opus-4-5-20251101",
    "tech_lead": "gemini-3.0-flash",
    "developer": "gemini-3.0-flash",
    "qa": "gemini-3.0-flash",
    "devops": "gemini-3.0-flash",
    "security": "claude-opus-4-5-20251101"
  },
  "composer.enabled": true,
  "composer.multiAgent": true,
  "terminal.sandboxed": true
}
```

**Benefits:**
- Fast (1-3s per agent)
- Cheap ($1-3 per epic)
- High quality (76% SWE-bench average)

### For Budget Teams

```json
{
  "chat.defaultModel": "composer",
  "chat.modelsByRole": {
    "analyst": "gemini-3.0-flash",
    "architect": "gemini-3.0-flash",
    "devops": "qwen3-coder-480b"
  }
}
```

**Benefits:**
- FREE (Composer included, Qwen local)
- Fast enough
- Acceptable quality (70% SWE-bench average)

### For Enterprise

```json
{
  "chat.defaultModel": "gemini-3.0-flash",
  "chat.modelsByRole": {
    "analyst": "claude-opus-4-5-20251101",
    "architect": "claude-opus-4-5-20251101",
    "tech_lead": "claude-sonnet-4-5-20250929",
    "developer": "gemini-3.0-flash",
    "qa": "gemini-3.0-flash",
    "devops": "gemini-3.0-flash",
    "security": "claude-opus-4-5-20251101"
  },
  "composer.enabled": true,
  "terminal.sandboxed": true,
  "audit.logAllAgents": true
}
```

**Benefits:**
- Best quality (Opus for critical, Flash for speed)
- Audit trail for compliance
- Sandboxed for security

## 📚 Additional Resources

- [MODELS.md](../../MODELS.md) - Complete model comparison with SWE-bench data
- [CLAUDE_CODE.md](CLAUDE_CODE.md) - Claude Code CLI integration
- [Cursor 2.0 Announcement](https://cursor.com/blog/2-0) - Official release notes
- [Cursor Features](https://cursor.com/features) - Full feature list
- [Cursor Changelog](https://cursor.com/changelog) - Latest updates

---

**Version:** 2.0
**Last Updated:** December 29, 2025
**Key Changes:**
- Updated for Cursor 2.0 (Composer, multi-agent mode, sandboxed terminals)
- Gemini 3 Flash as recommended default (76-78% SWE-bench, 13x cheaper)
- Multi-provider support (Claude, Google, OpenAI, Open Source)
- Real-world example with timing and cost data
- Advanced automation and optimization strategies
