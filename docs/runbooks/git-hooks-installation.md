# Git Hooks Installation Runbook

**Purpose:** Install cross-platform Git hooks for SDP quality checks.

**Prerequisites:**
- Git repository initialized
- Write permissions to `.git/hooks/`
- SDP CLI installed (recommended) or Bash access

## Why Git Hooks?

SDP uses cross-platform Git hooks that work across all IDEs:

✅ **Universal** - Works in Claude Code, Cursor, OpenCode
✅ **Pre-commit** - Quality checks before commit
✅ **Pre-push** - Regression tests before pushing
✅ **Post-merge** - Actions after merging
✅ **Post-checkout** - Actions after switching branches
✅ **Consistent** - Same behavior across all IDEs

## Available Hooks

| Hook | When Runs | What It Checks |
|------|-----------|----------------|
| **pre-commit** | Before commit | Quality checks: time estimates, code quality, Python code, Clean Architecture, WS format, breaking changes, test quality |
| **commit-msg** | During commit message validation | Conventional format + provenance trailers (`SDP-Agent`, `SDP-Model`, `SDP-Task`) |
| **post-commit** | After commit | Emits commit provenance event to evidence log; optionally comments on linked GitHub issues |
| **pre-push** | Before push | Regression tests and coverage checks |
| **post-merge** | After git merge | Post-merge actions |
| **post-checkout** | After git checkout | Post-checkout actions |

## Installation

### Recommended: SDP CLI (Go Implementation)

The SDP CLI provides canonical hook installation with marker-based ownership tracking:

```bash
# From repository root
sdp hooks install

# Include commit provenance hooks
sdp hooks install --with-provenance
```

**What it does:**
1. Installs 4 canonical SDP hooks by default: pre-commit, pre-push, post-merge, post-checkout
2. Optionally installs provenance hooks with `--with-provenance`: commit-msg, post-commit
3. Marks each hook with `# SDP-MANAGED-HOOK` for ownership tracking
4. Preserves existing non-SDP hooks
5. Makes hooks executable
6. Hooks check for `sdp` binary and warn if missing

**Expected output:**
```
✓ Installed pre-commit
✓ Installed pre-push
✓ Installed post-merge
✓ Installed post-checkout

Git hooks installed!
Customize hooks in .git/hooks/ if needed
```

For provenance validation after install:

```bash
sdp doctor hooks-provenance
```

**Re-running install (idempotent):**
- Existing SDP hooks are updated (if marker present)
- Non-SDP hooks are skipped
- No duplication occurs

### Alternative: Bash Script (Deprecated)

The bash script is kept for backward compatibility but will redirect to the Go implementation if available:

```bash
# From repository root
bash hooks/install-hooks.sh
```

**Note:** If `sdp` CLI is available, this script will delegate to `sdp hooks install`.

### Manual Installation

If automatic installation fails, you can manually create hooks with the SDP marker:

```bash
# 1. Navigate to repository root
cd /path/to/repository

# 2. Create hooks directory if needed
mkdir -p .git/hooks

# 3. Create hooks with SDP marker
for hook in pre-commit pre-push post-merge post-checkout; do
    cat > .git/hooks/$hook << 'EOF'
#!/bin/sh
# SDP-MANAGED-HOOK
# SDP Git Hook
# This hook is managed by SDP. Do not edit manually.

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks: https://github.com/fall-out-bug/sdp"
    exit 0
fi

# Add your custom checks here
EOF
    chmod +x .git/hooks/$hook
done
```

### Installation in Worktrees

If using git worktrees:

```bash
# Worktrees don't have their own .git/hooks/
# Hooks are shared with main repository

# Install from main repo:
cd /path/to/repository
sdp hooks install

# Hooks will work in all worktrees automatically
```

## Hook Ownership and Marker System

SDP hooks are marked with `# SDP-MANAGED-HOOK` to track ownership:

**Benefits:**
- **Safe uninstall** - Only removes SDP-managed hooks
- **Preserves custom hooks** - Non-SDP hooks are never touched
- **Idempotent install** - Re-running safely updates SDP hooks

**Marker location:**
- Second line of hook script (after shebang)
- Used by install/uninstall to identify SDP hooks

**Example:**
```sh
#!/bin/sh
# SDP-MANAGED-HOOK
# SDP Git Hook
# This hook is managed by SDP. Do not edit manually.
```

## Verification

### Check Hooks Installed

```bash
# List installed hooks
ls -la .git/hooks/{pre-commit,pre-push,post-merge,post-checkout}

# Check for SDP marker
grep "# SDP-MANAGED-HOOK" .git/hooks/pre-commit
```

### Check Hooks Executable

```bash
# Test each hook
test -x .git/hooks/pre-commit && echo "pre-commit: executable ✅" || echo "pre-commit: NOT executable ❌"
test -x .git/hooks/pre-push && echo "pre-push: executable ✅" || echo "pre-push: NOT executable ❌"
test -x .git/hooks/post-merge && echo "post-merge: executable ✅" || echo "post-merge: NOT executable ❌"
test -x .git/hooks/post-checkout && echo "post-checkout: executable ✅" || echo "post-checkout: NOT executable ❌"
```

### Test Hooks

```bash
# Create a test commit
touch test_file.txt
git add test_file.txt
git commit -m "test: verify hooks work"

# Expected: pre-commit hook runs and checks pass
```

## GitHub Integration (Optional)

The post-commit hook can comment on GitHub issues if environment variables are set:

### Setup

```bash
# Export in shell
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_REPO="fall-out-bug/msu_ai_masters"

# Or add to .env file
echo "GITHUB_TOKEN=ghp_your_token_here" >> .env
echo "GITHUB_REPO=fall-out-bug/msu_ai_masters" >> .env
```

### Token Scopes

GitHub token must have:
- `repo` - Full control of private repositories
- `project` - Access to GitHub projects (if using project boards)

### Create GitHub Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes: `repo`, `project`
4. Copy token and set as `GITHUB_TOKEN`

### How It Works

When you commit with a WS ID in the message (e.g., `feat: WS-060-01 - Create submission`):

1. post-commit hook extracts WS ID from commit message
2. Finds corresponding WS file in `workstreams/`
3. Extracts `github_issue` from WS frontmatter
4. Posts commit comment on GitHub issue

## Hook Behavior

### pre-commit

**What it checks:**

| Check | Description | Fails If |
|-------|-------------|-----------|
| Branch check | Warns if committing to main/master | N/A (warning only) |
| Time estimates | No days/hours/weeks in WS files | Found time-based estimates |
| Code quality markers | No TODO/FIXME/HACK in code | Found code quality markers |
| Python code quality | No bare except, specific exceptions | Bare except found |
| Clean Architecture | No infra imports in domain layer | Layer violation found |
| WS format | Goal + Acceptance Criteria exist | Missing required sections |
| Breaking changes | Auto-detects and generates docs | Breaking changes detected |
| Test quality | Test files meet quality standards | Test violations found |

**What to do on failure:**

1. Read error message
2. Fix the issue
3. Stage changes: `git add .`
4. Try again: `git commit`

### post-commit

**What it does:**

1. Gets latest commit hash and message
2. Reads provenance trailers (`SDP-Agent`, `SDP-Model`, `SDP-Task`)
3. Emits evidence event (`skill=commit`, `type=generation`) into `.sdp/log/events.jsonl`
4. Extracts WS ID (format: `WS-XXX-YY`)
5. Finds WS file in `workstreams/`
6. Extracts `github_issue` from WS frontmatter
7. Posts comment to GitHub issue

**Silent failures:**
- If `GITHUB_TOKEN` not set → silently skip
- If `GITHUB_REPO` not set → silently skip
- If no WS ID in commit message → silently skip
- If WS file not found → log warning, continue

### commit-msg

**What it does:**

1. Validates conventional commit first line
2. Auto-adds provenance trailers if missing:
   - `SDP-Agent` (`opencode` when `OPENCODE=1`, otherwise `human`)
   - `SDP-Model` (from `SDP_MODEL_ID`, `OPENCODE_MODEL`, `ANTHROPIC_MODEL`, `OPENAI_MODEL`, or `MODEL`)
   - `SDP-Task` (from `SDP_TASK_ID` or commit text tokens like `WS-123-45`, `00-123-45`, `F123`)

This metadata is used by post-commit evidence logging for commit provenance.

### pre-push

**What it checks:**

| Check | Description | Fails If |
|-------|-------------|-----------|
| Regression tests | Run `pytest -m fast` | Tests failed |
| Coverage | Check >= 80% | Coverage below 80% |

**Behavior Modes:**

**Default (WARNING mode):**
- Warns about failures but doesn't block push
- Allows flexibility during development
- Shows clear remediation steps

**Hard blocking mode (SDP_HARD_PUSH=1):**
- Blocks push on regression test failures
- Blocks push on coverage < 80%
- Enforces quality standards
- Same remediation steps

**Enable hard blocking:**

```bash
# Enable for current session
export SDP_HARD_PUSH=1

# Enable permanently (add to ~/.bashrc or ~/.zshrc)
echo 'export SDP_HARD_PUSH=1' >> ~/.bashrc
source ~/.bashrc
```

**What to do on failure:**

1. Read error/warning message
2. Run suggested remediation command
3. Fix failing tests or add coverage
4. Commit the fixes
5. Push again

**Bypass hooks (emergency only):**

```bash
# Not recommended - use only for critical hotfixes
git push --no-verify
```

## Uninstallation

### Recommended: SDP CLI (Safe Uninstall)

The SDP CLI safely removes only SDP-managed hooks:

```bash
# Remove SDP hooks (preserves non-SDP hooks)
sdp hooks uninstall
```

**What it does:**
- Only removes hooks with `# SDP-MANAGED-HOOK` marker
- Preserves any custom hooks you've added
- Safe to re-run (idempotent)

### Manual Uninstall

To manually remove all hooks:

```bash
# Remove all hooks (SDP and custom)
rm .git/hooks/{pre-commit,pre-push,post-merge,post-checkout}

# Verify removed
ls -la .git/hooks/{pre-commit,pre-push,post-merge,post-checkout}
# Expected: No such file or directory
```

## Troubleshooting

### Hooks not executing

**Problem:** Commits succeed without running hooks

**Solution:**

```bash
# 1. Check hooks installed
ls -la .git/hooks/{pre-commit,post-commit,pre-push}

# 2. Check hooks executable
test -x .git/hooks/pre-commit
test -x .git/hooks/post-commit
test -x .git/hooks/pre-push

# 3. Check hooks are symlinks
file .git/hooks/pre-commit
# Expected: symbolic link to sdp/hooks/pre-commit.sh

# 4. Re-install
bash sdp/hooks/install-hooks.sh
```

### Hooks causing build to fail

**Problem:** Pre-commit hook errors on build artifacts

**Solution:**

Add build artifacts to `.gitignore`:

```bash
echo "*.pyc" >> .gitignore
echo "__pycache__/" >> .gitignore
echo ".coverage" >> .gitignore
echo "dist/" >> .gitignore
```

### GitHub comments not posting

**Problem:** post-commit hook not commenting on issues

**Solution:**

```bash
# 1. Check environment variables
echo $GITHUB_TOKEN
echo $GITHUB_REPO

# 2. Check commit message has WS ID
git log -1 --pretty=%B | grep WS-

# 3. Check WS file has github_issue
grep "github_issue:" tools/hw_checker/docs/workstreams/backlog/WS-XXX-XX.md

# 4. Test curl manually
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/repos/$GITHUB_REPO/issues/1
```

### Pre-push blocking critical fixes

**Problem:** Need to push hotfix but pre-push checks fail

**Solution:**

```bash
# Option 1: Bypass hooks (emergency only, not recommended)
git push --no-verify

# Option 2: Disable hard blocking temporarily (recommended)
export SDP_HARD_PUSH=0
git push
export SDP_HARD_PUSH=1  # Re-enable after push

# Option 3: Fix issues first (best practice)
# Fix tests/coverage
git add .
git commit
git push

# Option 4: Disable hook temporarily (not recommended)
mv .git/hooks/pre-push .git/hooks/pre-push.disabled
git push
mv .git/hooks/pre-push.disabled .git/hooks/pre-push
```

**Note:** If hard blocking is enabled (SDP_HARD_PUSH=1), the hook will:
1. Block push on test failures
2. Block push on coverage < 80%
3. Show clear remediation steps
4. Require fix or explicit bypass (--no-verify)

## IDE-Specific Notes

### Claude Code

- PreToolUse/PostToolUse/Stop hooks disabled in `.claude/settings.json`
- Git hooks provide same quality checks
- Both don't run (no duplication)

### Cursor

- No native hooks API
- Uses Git hooks only
- Works seamlessly

### OpenCode

- No native hooks API
- Uses Git hooks only
- Works seamlessly

## Best Practices

1. **Install once** - Hooks work across all IDEs
2. **Read error messages** - Hooks explain what's wrong
3. **Fix before commit** - Don't try to bypass hooks
4. **Test hooks** - Create test commit after installation
5. **Keep hooks updated** - Run install after pulling changes

## References

- Pre-commit hook: `sdp/hooks/pre-commit.sh`
- Post-commit hook: `sdp/hooks/post-commit.sh`
- Pre-push hook: `sdp/hooks/pre-push.sh`
- Install script: `sdp/hooks/install-hooks.sh`
- PROJECT_MAP: `tools/hw_checker/docs/PROJECT_MAP.md` → Git Hooks section
