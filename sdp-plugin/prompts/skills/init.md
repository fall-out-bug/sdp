---
name: init
description: Initialize SDP in current project (interactive wizard)
tools: Read, Write, Bash, AskUserQuestion
---

# /init - SDP Project Setup Wizard

Interactive setup wizard for SDP projects.

## When to Use

- Setting up SDP in a new project
- Reconfiguring existing project
- Verifying SDP installation

## Workflow

### Step 1: Collect Project Metadata

Prompt for:
- **Project name** (default: directory name)
- **Description**
- **Author**

### Step 2: Detect Dependencies

Auto-detect:
- Beads CLI (task tracking)
- GitHub CLI (gh)
- Telegram (notifications)

### Step 3: Create Directory Structure

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

Create `quality-gate.toml`:
- Coverage: 80% minimum
- Complexity: CC < 10
- File size: 200 LOC max
- Type hints required

### Step 5: Create .env Template

Generate `.env.template` with placeholders for detected dependencies.

### Step 6: Install Git Hooks

Install pre-commit hook for SDP validation.

### Step 7: Run Doctor

Execute `sdp doctor` to validate setup.

## Usage

```bash
sdp init                    # Interactive
sdp init --non-interactive  # Use defaults
sdp init --path /project    # Target directory
sdp init --force            # Overwrite existing
```

## Output

- `docs/PROJECT_MAP.md`
- `docs/workstreams/INDEX.md`
- `docs/workstreams/TEMPLATE.md`
- `quality-gate.toml`
- `.env.template`
- `.git/hooks/pre-commit`

## Next Steps

After setup:
1. Edit `docs/PROJECT_MAP.md`
2. Copy `.env.template` to `.env`
3. Run `@idea "your first feature"`
