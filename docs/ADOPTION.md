# SDP Adoption Guide — Brownfield Projects

> Use this guide when adding SDP to an existing project. If you are starting a new project, start with [QUICKSTART.md](QUICKSTART.md) instead.

> **Note:** This guide describes the planned SDP CLI experience. Commands and flags are
> specification-level and will be verified against implementation as features land.

SDP is designed to layer onto existing work without disrupting it. This guide walks through the safe adoption path: assess first, try on a branch, then adopt when ready.

## Install Method Decision Matrix

Choose your install method based on your situation:

| Method | When to use | What it does | Risk level |
|--------|-------------|--------------|------------|
| `sdp assess` | You want to see what SDP would change before touching anything | Read-only scan of your repo. Produces a report. No files modified. | None |
| `sdp try "task"` | You want to test SDP on real work in isolation | Creates a temporary branch, runs the task, shows results. Branch is disposable. | Minimal |
| `install.sh` (project mode) | You are ready to install SDP prompts and hooks into your repo | Installs prompts, hooks, config into your project. Backs up existing files. | Low (backup + preview) |
| `install.sh --binary-only` | You want the CLI without touching project files | Installs `sdp` binary to `~/.local/bin/`. No project changes. | None |
| `git submodule add` | You want SDP vendored and version-locked inside your repo | Adds `sdp/` as a git submodule. Full control over updates. | Low |

> All methods install the `sdp` command-line binary. Package names differ by registry.

### Environment Variables

The installer respects these environment variables:

| Variable | Default | Purpose |
|----------|---------|---------|
| `SDP_IDE` | `auto` | Force a specific IDE: `claude`, `cursor`, `opencode`, `codex` |
| `SDP_DIR` | `sdp` | Change the SDP checkout directory name |
| `SDP_REF` | `main` | Install from a specific branch or tag |
| `SDP_REPO` | `fall-out-bug/sdp` | Install from a fork |

## Step 1: Assess Your Project

Before installing anything, run a read-only assessment:

```bash
sdp assess
```

Or point it at a specific path:

```bash
sdp assess /path/to/project
```

**What it does:**

- Scans your repo structure, languages, and tooling
- Reports which SDP surfaces would activate (skills, agents, hooks)
- Identifies potential conflicts with existing IDE config
- Suggests the best install method for your setup
- Writes nothing to disk

**Typical output includes:**

```
Project: my-app (Go + TypeScript)
Detected IDEs: Claude Code, Cursor
Recommended install: project mode (install.sh)

Conflicts:
  .claude/settings.json — will be merged, not overwritten

Suggested next step:
  sdp try "Add a health endpoint"
```

If `sdp` is not installed yet, use the binary-only install first:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
```

## Step 2: Try SDP on a Disposable Task

When the assessment looks good, test SDP on real work without committing to it:

```bash
sdp try "Add a health check endpoint"
```

**What it does:**

1. Creates a temporary branch (`sdp-try-<timestamp>`)
2. Initializes SDP config on that branch
3. Runs the task using SDP skills and agents
4. Shows the result: what changed, evidence produced, quality gates passed
5. Leaves you on the branch so you can inspect the result

**You control what happens next:**

```bash
# Inspect the changes
git diff main

# Keep the work if you like it
git checkout main
git merge sdp-try-<timestamp>

# Or discard it entirely
git checkout main
git branch -D sdp-try-<timestamp>
```

**Key flags:**

| Flag | Behavior |
|------|----------|
| `sdp try "task"` | Interactive: shows progress, asks before finalizing |
| `sdp try "task" --auto` | Fully automated: runs end to end without prompts |

## Step 3: Install SDP Into Your Project

Once you have tried SDP and want to keep it:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
```

### Preview Before Applying

Always preview first:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --preview
```

Preview mode shows every file that would be created, modified, or merged, without writing anything.

### What Gets Installed

The installer adds these to your project:

1. **`sdp/`** — SDP checkout (can be a git submodule or shallow clone)
2. **IDE adapters** — symlinks in `.claude/`, `.cursor/`, `.opencode/`, or `.codex/` pointing to prompts and agents
3. **`settings.json` merge** — SDP hooks are merged into your existing IDE settings. Your existing keys are preserved.
4. **Git hooks** — `pre-commit` and `pre-push` hooks (if not already present)
5. **`.gitignore` entries** — SDP-managed paths are added

### How Existing Config Is Protected

The installer never overwrites your existing config without consent:

- **Backup:** Before any modification, the original file is copied to `.sdp/backup/` with a timestamp.
- **Merge, not overwrite:** If `.claude/settings.json` already exists, SDP hooks are appended to arrays using `jq` deep merge. Your existing keys, hooks, and settings are untouched.
- **No `jq` installed?** The installer preserves your existing settings entirely and prints instructions for manual merge.
- **`--no-overwrite-config` flag:** Skip config file changes altogether. Prompts and agents still install.

### IDE Detection

The installer auto-detects which IDE you use by checking for `.claude/`, `.cursor/`, `.opencode/`, and `.codex/` directories, plus PATH entries for `claude`, `cursor`, `opencode`, `windsurf`, and `codex`.

Override detection with:

```bash
SDP_IDE=claude curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
```

## Step 4: Transition to Full SDP

After installing, SDP runs in **adoption mode** by default. Quality gates (file size limits, coverage thresholds, TDD enforcement) are disabled. Evidence logging stays active but is non-blocking.

### Adoption Mode

```bash
sdp init --auto
```

This creates `.sdp/config.yml` with `adoption_mode: true`. In this mode:

- `@build` and `@oneshot` execute without quality gate failures on legacy code
- Evidence events are still logged to `.sdp/log/events.jsonl`
- `sdp doctor` shows the current adoption mode

> **Tip:** For monorepo setups, see [FAQ: Does SDP work with monorepos?](#6-faq)

### Enable Full Gates

When your project is ready for full SDP enforcement:

```bash
sdp adopt --full
```

This sets `adoption_mode: false` in `.sdp/config.yml` and enables all quality gates:

- File size limits
- Test coverage thresholds
- TDD enforcement (test-first)
- Full evidence validation

**When to enable full gates:** after your project has passing tests, the codebase structure is stable, and you have at least one successful SDP feature delivery in adoption mode.

### Verify Your Setup

```bash
sdp doctor
```

`doctor` checks your environment and reports:

- Git status
- SDP CLI version
- IDE integration status
- Config file health
- Adoption mode (on or off)

## Step 5: Uninstall SDP

SDP can be removed cleanly. User data is preserved by default.

### Standard Uninstall

Removes SDP artifacts (symlinks, hooks, gitignore entries) while preserving your config and the SDP checkout:

```bash
sh sdp/scripts/uninstall.sh
```

**What is removed:**

- Symlinks in `.claude/skills`, `.claude/agents`, `.cursor/skills`, `.cursor/agents`, `.opencode/skills`, `.opencode/agents`, `.codex/skills/sdp`, `.codex/agents`
- SDP-installed files: `.claude/commands.json`, `.codex/INSTALL.md`, `.codex/skills/README.md`
- Git hooks pointing to SDP (`pre-commit`, `pre-push`)
- SDP entries from `.gitignore`

**What is preserved:**

- `.claude/settings.json` (your IDE settings, including non-SDP keys)
- `.sdp/` directory (project config and backups)
- `sdp/` checkout directory

### Dry Run

See what would be removed without removing anything:

```bash
sh sdp/scripts/uninstall.sh --dry-run
```

### Purge Everything

Remove all SDP-related files, including the checkout, backups, and config:

```bash
sh sdp/scripts/uninstall.sh --purge
```

Even in purge mode, `.claude/settings.json` is preserved because it likely contains your own customizations beyond SDP.

Skip the confirmation prompt:

```bash
sh sdp/scripts/uninstall.sh --purge -y
```

### Reinstall After Uninstall

```bash
sh sdp/scripts/install-project.sh
```

Or from remote:

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
```

## Quick Reference

```bash
# Assess before installing
sdp assess

# Try a task in isolation
sdp try "Add a health endpoint"

# Install (preview first)
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --preview
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh

# Verify setup
sdp doctor

# Enable full quality gates
sdp adopt --full

# Uninstall (dry run first)
sh sdp/scripts/uninstall.sh --dry-run
sh sdp/scripts/uninstall.sh

# Full purge
sh sdp/scripts/uninstall.sh --purge
```

## Related Documentation

| Document | Content |
|----------|---------|
| [QUICKSTART.md](QUICKSTART.md) | 5-minute guide for new projects |
| [CLI_REFERENCE.md](CLI_REFERENCE.md) | Full command reference |
| [PROTOCOL.md](PROTOCOL.md) | Protocol overview |
| [PRODUCT_CONTRACT.md](PRODUCT_CONTRACT.md) | Product definition and mode policy |
