# 00-041-01: Plugin Package Structure

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** MEDIUM
> **Created:** 2026-02-02

## Goal

Create Claude Plugin package format with .claude/ prompts for zero-Python distribution.

## Acceptance Criteria

- AC1: Plugin manifest (`plugin.json`) created with metadata and permissions
- AC2: All 18 skills copied to `sdp-plugin/prompts/skills/`
- AC3: All 11 agents copied to `sdp-plugin/prompts/agents/`
- AC4: README.md with zero-Python installation instructions
- AC5: Plugin package loads in Claude Code successfully

## Scope

### Input Files
- `.claude/skills/*/SKILL.md` (18 skills)
- `.claude/agents/*.md` (11 agents)
- Existing SDP documentation for reference

### Output Files
- `sdp-plugin/plugin.json` (NEW)
- `sdp-plugin/prompts/skills/*.md` (COPIED)
- `sdp-plugin/prompts/agents/*.md` (COPIED)
- `sdp-plugin/README.md` (NEW)
- `sdp-plugin/.claudeignore` (NEW)

### Out of Scope
- Modifying existing .claude/ files (Python SDP remains untouched)
- Creating validators (WS-00-041-04)
- Go binary development (WS-00-041-05)

## Implementation Steps

1. **Research Claude Plugin Schema**
   - Check Claude Plugin Marketplace documentation
   - Understand plugin.json format
   - Identify required permissions

2. **Create Plugin Manifest**
   - `sdp-plugin/plugin.json`
   - Include: name, version, description, permissions
   - Define prompts location

3. **Copy Prompts**
   - Copy all 18 skills from `.claude/skills/*/SKILL.md` → `sdp-plugin/prompts/skills/{name}.md`
   - Copy all 11 agents from `.claude/agents/*.md` → `sdp-plugin/prompts/agents/{name}.md`
   - Preserve directory structure

4. **Create README**
   - Zero-Python quick start
   - Plugin installation instructions
   - Skills overview
   - Language support note (future work)

5. **Create .claudeignore**
   - Exclude unnecessary files
   - Documentation files
   - Test files

## Verification

```bash
# Test plugin structure
ls -la sdp-plugin/prompts/skills/  # Should show 18 .md files
ls -la sdp-plugin/prompts/agents/  # Should show 11 .md files

# Validate plugin.json format
cat sdp-plugin/plugin.json | python -m json.tool  # Should parse without errors

# Manual test in Claude Code (when available)
# claude plugin load sdp-plugin/plugin.json
# Expected: All 18 skills, 11 agents available
```

## Quality Gates

- plugin.json valid JSON schema
- All 18 skills present (grep count)
- All 11 agents present (grep count)
- README has installation section
- No Python dependencies mentioned in README

## Dependencies

None

## Blocks

- 00-041-03 (Remove Python Dependencies from Skills)
- 00-041-04 (AI-Based Validation Prompts)
