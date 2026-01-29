# Multi-IDE SDP Parity

**Version:** 1.0.0
**Updated:** 2026-01-23
**Status:** ✅ Complete

---

## Overview

SDP (Spec-Driven Protocol) supports three AI coding tools: Claude Code, Cursor, OpenCode.

All IDEs use the same master prompts from `sdp/prompts/commands/`, ensuring consistent behavior across platforms.

## Parity Matrix

| Feature | Claude Code | Cursor | OpenCode |
|---------|-------------|--------|----------|
| **Slash Commands** | | | |
| /idea | ✅ `.claude/skills/idea/SKILL.md` | ✅ `.cursor/commands/idea.md` | ✅ `.opencode/commands/idea.md` |
| /design | ✅ `.claude/skills/design/SKILL.md` | ✅ `.cursor/commands/design.md` | ✅ `.opencode/commands/design.md` |
| /build | ✅ `.claude/skills/build/SKILL.md` | ✅ Uses `.claude/agents/` | ✅ Uses `.claude/agents/` |
| /oneshot | ✅ `.claude/skills/oneshot/SKILL.md` | ✅ `.cursor/commands/oneshot.md` | ✅ `.opencode/commands/oneshot-simple.md` |
| /test | ✅ `.claude/skills/test/SKILL.md` | ✅ `.cursor/commands/test.md` | ✅ `.opencode/commands/test.md` |
| /debug | ✅ `.claude/skills/debug/SKILL.md` | ✅ `.cursor/commands/debug.md` | ✅ `.opencode/commands/debug.md` |
| /issue | ✅ `.claude/skills/issue/SKILL.md` | ✅ Uses master prompts | ✅ Uses master prompts |
| /hotfix | ✅ `.claude/skills/hotfix/SKILL.md` | ✅ Uses master prompts | ✅ Uses master prompts |
| /bugfix | ✅ `.claude/skills/bugfix/SKILL.md` | ✅ Uses master prompts | ✅ Uses master prompts |
| /codereview | ✅ `.claude/skills/codereview/SKILL.md` | ✅ Uses master prompts | ✅ Uses master prompts |
| /deploy | ✅ `.claude/skills/deploy/SKILL.md` | ✅ Uses master prompts | ✅ Uses master prompts |
| **Agents** | | | |
| builder | ✅ `.claude/agents/builder.md` | ✅ `.cursor/agents/builder.md` | ✅ `opencode.json` → agent: builder |
| reviewer | ✅ `.claude/agents/reviewer.md` | ✅ `.cursor/agents/reviewer.md` | ✅ `opencode.json` → agent: reviewer |
| planner | ✅ `.claude/agents/planner.md` | ✅ `.cursor/agents/planner.md` | ✅ `opencode.json` → agent: planner |
| deployer | ✅ `.claude/agents/deployer.md` | ✅ `.cursor/agents/deployer.md` | ✅ `opencode.json` → agent: deployer |
| orchestrator | ✅ `.claude/agents/orchestrator.md` | ✅ `.cursor/agents/orchestrator.md` | ✅ `opencode.json` → agent: orchestrator |
| **Hooks** | | | |
| PreToolUse | ✅ `.claude/settings.json` | ❌ No hooks API | ❌ No hooks API |
| PostToolUse | ✅ `.claude/settings.json` | ❌ No hooks API | ❌ No hooks API |
| Stop | ✅ `.claude/settings.json` | ❌ No hooks API | ❌ No hooks API |
| Git hooks (pre-commit) | ✅ | ✅ | ✅ |
| Git hooks (post-commit) | ✅ | ✅ | ✅ |
| Git hooks (pre-push) | ✅ | ✅ | ✅ |
| **Configuration** | | | |
| settings.json | ✅ `.claude/settings.json` | ❌ | ❌ |
| cursorrules | ❌ | ✅ `.cursorrules` | ❌ |
| opencode commands/ | ❌ | ❌ | ✅ `.opencode/commands/` |
| opencode.json | ❌ | ❌ | ✅ `.opencode/opencode.json` |

## Architecture

### Master Prompts

All IDEs use the same master prompts from `sdp/prompts/commands/`:

✅ **Single source of truth**
✅ **Universal workflows**
✅ **Consistent behavior**

### Hooks

**Cross-platform Git hooks** work in all IDEs:

- **Claude Code:** Git hooks (PreToolUse/PostToolUse/Stop disabled in settings.json to avoid duplication)
- **Cursor:** Git hooks only (no native hooks API)
- **OpenCode:** Git hooks only (no native hooks API)

See: [Git Hooks Installation Runbook](runbooks/git-hooks-installation.md)

### Agents

**Claude Code:** `.claude/agents/` with SKILL.md files

**Cursor:** `.cursor/agents/` with delegation to master prompts

**OpenCode:** `opencode.json` with agent configuration:
```json
{
  "agent": {
    "builder": {
      "mode": "primary",
      "description": "...",
      "prompt": "See @sdp/prompts/commands/build.md",
      "tools": { ... }
    }
  }
}
```

### Commands Format

**Claude Code/Cursor:**
```markdown
---
name: command_name
description: Description
tools: [...]
model: inherit
---
[Prompt content]
```

**OpenCode:**
```markdown
---
description: Description
agent: agent_name (optional)
---
[Prompt template]
```

**⚠️ Critical:** OpenCode commands MUST NOT include `model: inherit` field (causes "Model not found: inherit" error)

## Installation

### Claude Code

Already configured in `.claude/settings.json`:

```bash
# Git hooks
bash sdp/hooks/install-hooks.sh

# Nothing additional needed
```

### Cursor

```bash
# Install Git hooks
bash sdp/hooks/install-hooks.sh

# Enable cursorrules (if not enabled)
cp .cursor/rules/cursorrules-unified.md .cursorrules
```

### OpenCode

```bash
# Install Git hooks
bash sdp/hooks/install-hooks.sh

# Commands already created in .opencode/commands/
# See OpenCode documentation for more details
```

## Usage

### Runbooks

See runbooks for detailed workflows:

- [Oneshot Runbook](runbooks/oneshot-runbook.md) - Autonomous feature execution
- [Debug Runbook](runbooks/debug-runbook.md) - Systematic debugging workflow
- [Test Runbook](runbooks/test-runbook.md) - Contract-driven test generation
- [Git Hooks Installation](runbooks/git-hooks-installation.md) - Cross-platform hooks

### Quick Start

#### For New Users

1. **Choose your IDE:**
   - **Claude Code:** Full support (PreToolUse/PostToolUse/Stop hooks + agents + skills)
   - **Cursor:** Good support (git hooks + agents + commands)
   - **OpenCode:** Basic support (git hooks + commands + opencode.json agents)

2. **Install dependencies:**
   ```bash
   cd sdp
   poetry install
   ```

3. **Install hooks:**
   ```bash
   bash sdp/hooks/install-hooks.sh
   ```

4. **Start with /idea:**
   ```bash
   /idea "my feature description"
   ```

#### Example Workflow

```bash
# 1. Gather requirements
/idea "add LMS integration"

# 2. Design workstreams
/design idea-lms-integration

# 3. Generate test contract
/test WS-060-01

# 4. Implement first WS
/build WS-060-01

# 5. Implement next WS
/test WS-060-02
/build WS-060-02

# 6. Review all WS
/codereview F60

# 7. Human UAT
# (manual testing)

# 8. Deploy
/deploy F60
```

## OpenCode Specifics

### Command Format

OpenCode commands use a **different frontmatter format:**

```markdown
---
description: Command description
agent: agent_name
---
[Prompt content]
```

**Critical Differences:**
- ❌ NO `name:` field
- ❌ NO `model:` field (causes "Model not found: inherit" error)
- ✅ Only `description:` and `agent:` fields

### Agent Configuration

OpenCode uses **JSON config** for all agents (`.opencode/opencode.json`):

```json
{
  "$schema": "https://opencode.ai/config.json",
  "agent": {
    "builder": {
      "mode": "primary",
      "description": "TDD execution agent",
      "prompt": "See @sdp/prompts/commands/build.md",
      "tools": {
        "write": true,
        "edit": true,
        "bash": true,
        "read": true,
        "glob": true,
        "grep": true,
        "webfetch": true
      }
    }
  }
}
```

**Note:** In JSON config, `model: inherit` IS allowed (different from command files)

### /oneshot Command

OpenCode uses `/oneshot-simple` instead of `/oneshot` because:
- `/oneshot` with `model: inherit` causes error in command files
- Both commands delegate to same master prompt (`sdp/prompts/commands/oneshot.md`)

## Troubleshooting

### /oneshot doesn't work in Cursor

**Problem:** `/oneshot` command not found

**Solution:**
1. Check `.cursor/commands/oneshot.md` exists
2. Restart Cursor
3. Check Cursor logs for errors

### /debug doesn't work in OpenCode

**Problem:** `/debug` command not found

**Solution:**
1. Check `.opencode/commands/debug.md` exists
2. Verify NO `model:` field in frontmatter
3. Restart OpenCode

### Git hooks not running

**Problem:** Hooks don't execute on commit

**Solution:**
1. Run `bash sdp/hooks/install-hooks.sh`
2. Verify hooks installed: `ls -la .git/hooks/`
3. Check permissions: `test -x .git/hooks/pre-commit`

### Model not found: inherit (OpenCode)

**Problem:** Error when using OpenCode commands

**Solution:**
1. Check command file for `model: inherit` field
2. Remove `model:` field from frontmatter
3. Save and restart OpenCode

## Contributing

When adding new features:

1. ✅ Update parity matrix (above)
2. ✅ Create runbook if needed
3. ✅ Test in all 3 IDEs
4. ✅ Update this document
5. ✅ Respect OpenCode format constraints

## Known Limitations

### OpenCode

- `/oneshot` must use `/oneshot-simple` (avoid `model:` field)
- No PreToolUse/PostToolUse/Stop hooks (use Git hooks)
- Agent configuration uses JSON instead of markdown files
- `/idea`, `/design` commands not yet implemented

### Cursor

- `/idea` command not yet implemented
- No native hooks API (use Git hooks)

## Future Work

- [ ] Implement `/idea` command in Cursor and OpenCode
- [ ] Implement `/design` command in OpenCode
- [ ] Add OpenCode-specific documentation
- [ ] Improve OpenCode agent discoverability

---

**See also:**
- [SDP Protocol](../PROTOCOL.md)
- [HW Checker Patterns](../HW_CHECKER_PATTERNS.md)
- [PROJECT_MAP](../tools/hw_checker/docs/PROJECT_MAP.md)
- [F201 Workstreams](../tools/hw_checker/docs/workstreams/backlog/)
