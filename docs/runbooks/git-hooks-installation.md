# Git Hooks Installation Runbook

**Purpose:** Install cross-platform Git hooks for SDP quality checks.

**Prerequisites:**
- Git repository initialized
- Bash access
- Write permissions to `.git/hooks/`

## Why Git Hooks?

SDP uses cross-platform Git hooks that work across all IDEs:

✅ **Universal** - Works in Claude Code, Cursor, OpenCode
✅ **Pre-commit** - Quality checks before commit
✅ **Post-commit** - GitHub issue sync (if GITHUB_TOKEN set)
✅ **Pre-push** - Regression tests before pushing
✅ **Consistent** - Same behavior across all IDEs

## Available Hooks

| Hook | When Runs | What It Checks |
|------|-----------|----------------|
| **pre-commit** | Before commit | Quality checks: time estimates, code quality, Python code, Clean Architecture, WS format, breaking changes, test quality |
| **post-commit** | After commit | GitHub issue sync (comments on linked issues if GITHUB_TOKEN set) |
| **pre-push** | Before push | Regression tests and coverage checks |

## Installation

### Automatic Installation

```bash
# From repository root
bash sdp/hooks/install-hooks.sh
```

**What it does:**
1. Creates symlinks from `.git/hooks/` to `sdp/hooks/*.sh`
2. Makes hooks executable (`chmod +x`)
3. Verifies installation

**Expected output:**
```
Installing post-commit hook...
✓ Installed: post-commit

Installing pre-commit hook...
✓ Installed: pre-commit

Installing pre-push hook...
✓ Installed: pre-push

✅ Git hooks installed successfully

Hooks:
  - pre-commit:  quality checks (time, code quality, Python, Clean Arch, WS format)
  - post-commit: GitHub issue sync (if GITHUB_TOKEN set)
  - pre-push:    regression tests

Required environment variables for GitHub integration:
  GITHUB_TOKEN  - GitHub personal access token
  GITHUB_REPO   - Repository in format 'owner/repo'

Add to .env file or export in shell
```

### Manual Installation

If automatic installation fails:

```bash
# 1. Navigate to repository root
cd /path/to/msu_ai_masters

# 2. Create symlinks manually
ln -sf sdp/hooks/pre-commit.sh .git/hooks/pre-commit
ln -sf sdp/hooks/post-commit.sh .git/hooks/post-commit
ln -sf sdp/hooks/pre-push.sh .git/hooks/pre-push

# 3. Make executable
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/post-commit
chmod +x .git/hooks/pre-push

# 4. Verify installation
ls -la .git/hooks/{pre-commit,post-commit,pre-push}
```

### Installation in Worktrees

If using git worktrees:

```bash
# Worktrees don't have their own .git/hooks/
# Hooks are shared with main repository

# Install from main repo:
cd /path/to/msu_ai_masters
bash sdp/hooks/install-hooks.sh

# Hooks will work in all worktrees automatically
```

## Verification

### Check Hooks Installed

```bash
# List installed hooks
ls -la .git/hooks/{pre-commit,post-commit,pre-push}

# Expected output:
# lrwxrwxrwx 1 user user 30 Jan 23 10:00 pre-commit -> sdp/hooks/pre-commit.sh
# lrwxrwxrwx 1 user user 32 Jan 23 10:00 post-commit -> sdp/hooks/post-commit.sh
# lrwxrwxrwx 1 user user 28 Jan 23 10:00 pre-push -> sdp/hooks/pre-push.sh
```

### Check Hooks Executable

```bash
# Test each hook
test -x .git/hooks/pre-commit && echo "pre-commit: executable ✅" || echo "pre-commit: NOT executable ❌"
test -x .git/hooks/post-commit && echo "post-commit: executable ✅" || echo "post-commit: NOT executable ❌"
test -x .git/hooks/pre-push && echo "pre-push: executable ✅" || echo "pre-push: NOT executable ❌"
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
2. Extracts WS ID (format: `WS-XXX-YY`)
3. Finds WS file in `workstreams/`
4. Extracts `github_issue` from WS frontmatter
5. Posts comment to GitHub issue

**Silent failures:**
- If `GITHUB_TOKEN` not set → silently skip
- If `GITHUB_REPO` not set → silently skip
- If no WS ID in commit message → silently skip
- If WS file not found → log warning, continue

### pre-push

**What it checks:**

| Check | Description | Fails If |
|-------|-------------|-----------|
| Regression tests | Run `pytest -m fast` | Tests failed |
| Coverage | Check >= 80% | Coverage below 80% |

**Note:** pre-push **warns** but doesn't block pushing

**What to do on warning:**

1. Read warning message
2. Fix issues if critical
3. Push: `git push`
4. Or fix and try again

## Uninstallation

```bash
# Remove all hooks
rm .git/hooks/pre-commit
rm .git/hooks/post-commit
rm .git/hooks/pre-push

# Or use one command
rm .git/hooks/{pre-commit,post-commit,pre-push}

# Verify removed
ls -la .git/hooks/{pre-commit,post-commit,pre-push}
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
# Option 1: Bypass hooks (not recommended)
git push --no-verify

# Option 2: Fix issues first (recommended)
# Fix tests/coverage
git add .
git commit
git push

# Option 3: Disable hook temporarily
mv .git/hooks/pre-push .git/hooks/pre-push.disabled
git push
mv .git/hooks/pre-push.disabled .git/hooks/pre-push
```

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
