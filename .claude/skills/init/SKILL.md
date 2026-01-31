---
name: init
description: Initialize SDP in current project (interactive wizard)
tools: Read, Write, Bash, AskUserQuestion
---

# /init - SDP Project Setup Wizard

Interactive setup wizard for SDP projects.

## When to Use

- Setting up SDP in a new project
- Reconfiguring SDP in existing project
- Verifying SDP installation

## Workflow

### Step 1: Collect Project Metadata

Prompt for:
- **Project name**: Default from directory name
- **Description**: Brief project description
- **Author**: Project author/maintainer

### Step 2: Detect Optional Dependencies

Auto-detect:
- Beads CLI (task tracking)
- GitHub CLI (gh)
- Telegram (notifications)

### Step 3: Create Directory Structure

Create standard directories:
```
docs/
├── workstreams/
│   ├── INDEX.md
│   ├── TEMPLATE.md
│   └── backlog/
├── PROJECT_MAP.md
└── drafts/
sdp.local/
```

### Step 4: Generate Quality Gate Config

Create `quality-gate.toml` with:
- Coverage settings (80% minimum)
- Complexity limits (CC < 10)
- File size limits (200 LOC)
- Type hint requirements
- Error handling rules
- Architecture constraints

### Step 5: Create .env Template

Generate `.env.template` with placeholders for detected dependencies:
- Telegram bot token/chat ID
- GitHub token/repo
- Beads API URL

### Step 6: Install Git Hooks

Install pre-commit hook for SDP validation.

### Step 7: Run Doctor

Execute `sdp doctor` to validate setup.

## Usage

```bash
# Interactive setup (prompts for values)
sdp init

# Use defaults
sdp init --non-interactive

# Target specific directory
sdp init --path /path/to/project

# Overwrite existing files
sdp init --force
```

## Output

- `docs/PROJECT_MAP.md` — Project decision log
- `docs/workstreams/INDEX.md` — Workstream index
- `docs/workstreams/TEMPLATE.md` — Workstream template
- `quality-gate.toml` — Quality gate configuration
- `.env.template` — Environment variable template
- `.git/hooks/pre-commit` — Git hook (if git repo)

## Next Steps

After setup:
1. Edit `docs/PROJECT_MAP.md` with project info
2. Run `sdp extension list` to see available extensions
3. Start with `/idea "your first feature"`
