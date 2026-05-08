---
name: init
description: Initialize SDP in a new or existing project
version: 1.0.0
changes:
  - Initial release: project detection, config scaffolding, harness setup
---

# @init

Initialize SDP (Spec-Driven Protocol) in the current project directory.

## Workflow

When user invokes `@init` or `sdp init`:

### Step 1: Project Detection

Auto-detect project characteristics:

```bash
# Detect language
if [ -f "go.mod" ]; then LANG="go"
elif [ -f "pyproject.toml" ] || [ -f "requirements.txt" ]; then LANG="python"
elif [ -f "pom.xml" ] || [ -f "build.gradle" ]; then LANG="java"
elif [ -f "package.json" ]; then LANG="nodejs"
elif [ -f "Cargo.toml" ]; then LANG="rust"
else LANG="unknown"
fi

# Detect framework (language-specific)
# Detect existing SDP artifacts
```

### Step 2: Ask Configuration Questions

1. **Project name** — default: directory name
2. **Primary language** — default: detected
3. **AI harness(es)** — which AI tools the team uses (claude, codex, cursor, opencode, zed, warp, other). Default: detect from existing config files.
4. **Issue tracker** — beads (default), github-issues, linear, none

### Step 3: Scaffold SDP Structure

Create the following (skip existing):

```
docs/
  roadmap/ROADMAP.md     # Feature roadmap
  workstreams/            # Workstream tracking
  drafts/                 # Discovery drafts
.sdp/                     # SDP internal state
  checkpoints/
AGENTS.md                 # Agent instructions (harness-neutral)
```

### Step 4: Configure Harnesses

For each selected harness, create appropriate config:

- **Claude Code** — `.claude/settings.json` (hooks), `.claude/commands/` (slash commands)
- **Codex** — `.codex/` with symlinks to `prompts/skills/`
- **Cursor** — `.cursor/` with `.cursorrules` and skill symlinks
- **OpenCode** — `.opencode/opencode.json` with agent cards

All harness configs point to the same canonical source: `prompts/skills/` and `prompts/agents/`.

### Step 5: Configure Quality Gates

Based on detected language, set up appropriate quality gate commands:

| Language | Build | Test | Lint |
|----------|-------|------|------|
| Go | `go build ./...` | `go test ./...` | `go vet ./...` |
| Python | `pip install .` | `pytest` | `ruff check .` |
| Node.js | `npm run build` | `npm test` | `npm run lint` |
| Rust | `cargo build` | `cargo test` | `cargo clippy` |
| Java | `mvn compile` | `mvn test` | `mvn checkstyle:check` |

Write the detected gates into `AGENTS.md`.

### Step 6: Initialize Issue Tracker

If beads selected:
```bash
bd init
```

### Step 7: Verify

- All configured harness symlinks resolve
- AGENTS.md exists with quality gates
- Issue tracker is functional (if selected)
- `sdp health` passes

## Flags

| Flag | Description |
|------|-------------|
| `--auto` | Skip all questions, accept defaults from detection |
| `--lang <lang>` | Force specific language |
| `--harness <list>` | Comma-separated list of harnesses to configure |

## When to Use

- New project starting with SDP
- Existing project adopting SDP
- Adding a new AI harness to existing SDP project
- After cloning an SDP project (verify/setup)

## Output

```
SDP initialized: {project_name}
Language: {detected}
Harnesses: {configured list}
Quality gates: {build}, {test}, {lint}
Issue tracker: {selected}

Next steps:
  @vision "your product idea"   # Start from scratch
  @reality --quick              # Analyze existing codebase
  @feature "add X"              # Plan a feature
```

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Re-run `@init` with `--auto` flag to skip prompts; verify you have write permissions in current directory |
| Harness config not created | Check harness CLI is installed; re-run `@init --harness <name>` for missing harness |
| "bd init fails" | Verify beads CLI installed (`bd --version`); run `bd init` manually |
| Quality gates not detected | Manually specify: `@init --lang <lang>` |

## See Also

- @vision -- Strategic planning
- @reality -- Codebase analysis
- @feature -- Feature planning
