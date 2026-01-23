# Setup Guide

Complete guide to configure GitHub integration for the Spec-Driven Protocol.

## Prerequisites

- GitHub account with repository access
- Git repository: `fall-out-bug/msu_ai_masters`
- Python 3.11+ with Poetry installed
- Git configured with user name and email

## Step 1: Create GitHub Personal Access Token

### 1.1 Navigate to Token Settings

Go to: https://github.com/settings/tokens

### 1.2 Generate New Token

1. Click **"Generate new token (classic)"**
2. Name: `sdp-github-automation`
3. Expiration: Choose expiration (90 days recommended for security)
4. Select scopes:

**Required Scopes:**
- ✅ **`repo`** - Full control of private repositories
  - repo:status
  - repo_deployment
  - public_repo
  - repo:invite
  - security_events
- ✅ **`project`** - Full control of organization projects (new GitHub Projects)
  - read:project
  - write:project
- ✅ **`org:read`** - Read org and team membership, read org projects

### 1.3 Generate and Copy Token

1. Click **"Generate token"**
2. **Copy the token immediately** (you won't see it again!)
3. Token format: `ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

**Security Note:** Treat this token like a password. Never commit it to Git!

## Step 2: Configure Environment Variables

### 2.1 Create .env File

Create `.env` file in project root (`/home/fall_out_bug/msu_ai_masters/.env`):

```bash
# GitHub API Configuration
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_REPO=fall-out-bug/msu_ai_masters
GITHUB_ORG=fall-out-bug

# Optional: Override project routing
# GITHUB_PROJECT=mlsd

# Optional: Telegram notifications (if configured)
# TELEGRAM_BOT_TOKEN=your_bot_token
# TELEGRAM_CHAT_ID=your_chat_id
```

### 2.2 Verify .env is Gitignored

```bash
# Check .gitignore includes .env
grep "^\.env$" .gitignore

# If not, add it
echo ".env" >> .gitignore
```

**Critical:** Never commit `.env` to Git!

### 2.3 Load Environment Variables

```bash
# Option 1: Export manually (temporary, current session only)
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GITHUB_REPO="fall-out-bug/msu_ai_masters"
export GITHUB_ORG="fall-out-bug"

# Option 2: Use direnv (automatic, recommended)
# Install: https://direnv.net/
echo "export GITHUB_TOKEN=ghp_xxx" >> .envrc
echo "export GITHUB_REPO=fall-out-bug/msu_ai_masters" >> .envrc
echo "export GITHUB_ORG=fall-out-bug" >> .envrc
direnv allow

# Option 3: Source .env in shell (manual)
set -a
source .env
set +a
```

### 2.4 Verify Configuration

```bash
# Test environment variables are set
echo $GITHUB_TOKEN | head -c 10  # Should show "ghp_xxxxxx"
echo $GITHUB_REPO                 # Should show "fall-out-bug/msu_ai_masters"
echo $GITHUB_ORG                  # Should show "fall-out-bug"
```

## Step 3: Install Dependencies

### 3.1 Install SDP Package

```bash
cd /home/fall_out_bug/msu_ai_masters/sdp
poetry install --no-interaction
```

### 3.2 Verify Installation

```bash
# Check PyGithub is installed
poetry run python -c "import github; print(github.__version__)"

# Test imports
poetry run python -c "
from sdp.github.project_router import ProjectRouter
from sdp.github.projects_client import ProjectsClient
print('✅ GitHub integration modules loaded')
"
```

## Step 4: Test GitHub Connection

### 4.1 Test API Authentication

```bash
# Test GitHub API connection
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/user \
     | jq '.login'
# Expected: Your GitHub username

# Check rate limit
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/rate_limit \
     | jq '.rate.remaining'
# Expected: ~5000 (for authenticated requests)
```

### 4.2 Test Python GitHub Client

```bash
cd sdp
poetry run python << 'EOF'
import os
from github import Github

token = os.getenv("GITHUB_TOKEN")
repo_name = os.getenv("GITHUB_REPO")

if not token:
    print("❌ GITHUB_TOKEN not set")
    exit(1)

client = Github(token)
repo = client.get_repo(repo_name)

print(f"✅ Connected to: {repo.full_name}")
print(f"✅ Default branch: {repo.default_branch}")
print(f"✅ Open issues: {repo.open_issues_count}")
EOF
```

Expected output:
```
✅ Connected to: fall-out-bug/msu_ai_masters
✅ Default branch: main
✅ Open issues: X
```

## Step 5: Install Git Hooks

### 5.1 Run Hook Installation Script

```bash
# From repository root
cd /home/fall_out_bug/msu_ai_masters
./sdp/hooks/install-hooks.sh
```

This installs:
- `post-commit` - Auto-comment on GitHub issues when committing
- `pre-commit` - Validate WS files and code quality (already installed)
- `commit-msg` - Validate conventional commit format (already installed)

### 5.2 Verify Hooks Installed

```bash
# Check hooks exist
ls -la .git/hooks/post-commit
ls -la .git/hooks/pre-commit
ls -la .git/hooks/commit-msg

# Verify they're executable
test -x .git/hooks/post-commit && echo "✅ post-commit executable"
test -x .git/hooks/pre-commit && echo "✅ pre-commit executable"
test -x .git/hooks/commit-msg && echo "✅ commit-msg executable"
```

### 5.3 Test Post-Commit Hook (Dry Run)

```bash
# Test hook without making a real commit
export GITHUB_DRY_RUN=true
.git/hooks/post-commit
# Should show: "Dry-run mode: would post comment to issue #X"
```

## Step 6: Setup GitHub Projects (Optional)

GitHub Projects are created automatically on first sync, but you can create them manually for better control.

### 6.1 Create Project Manually

1. **Navigate to Projects:**
   - Go to: https://github.com/fall-out-bug/msu_ai_masters/projects
   - Or: https://github.com/orgs/fall-out-bug/projects

2. **Create New Project:**
   - Click **"New project"**
   - Choose **"Board"** template
   - Name: `mlsd` (for ML System Design course)

3. **Configure Board Columns:**
   - Create columns: **Backlog**, **In Progress**, **In Review**, **Done**
   - Or use default columns and rename them

4. **Add Status Field:**
   - Click **"+ New field"**
   - Type: **Single select**
   - Name: **Status**
   - Options: `Backlog`, `In Progress`, `In Review`, `Done`

5. **Repeat for Other Projects:**
   - Create `bdde` project for Big Data & Data Engineering course

### 6.2 Automatic Project Creation

Alternatively, let the system create projects automatically:

```bash
cd sdp
poetry run python << 'EOF'
from github import Github
from sdp.github.projects_client import ProjectsClient
import os

client = Github(os.getenv("GITHUB_TOKEN"))
repo = client.get_repo(os.getenv("GITHUB_REPO"))
projects_client = ProjectsClient(client, repo, os.getenv("GITHUB_REPO"))

# This will auto-create project if missing
project = projects_client.get_or_create_project(
    "mlsd",
    "ML System Design course workstreams"
)
print(f"✅ Project created: {project['title']} (ID: {project['id']})")
EOF
```

## Step 7: Verify Full Setup

### 7.1 Run Test Sync (Dry Run)

```bash
cd /home/fall_out_bug/msu_ai_masters/sdp

# Test sync without making changes
poetry run python -c "
from pathlib import Path
from sdp.github.project_router import ProjectRouter

# Test project routing
ws = Path('tools/hw_checker/docs/workstreams/backlog/WS-160-01.md')
project = ProjectRouter.get_project_for_ws(ws)
print(f'✅ Project routing works: {ws.name} → {project}')
"
```

Expected: `✅ Project routing works: WS-160-01.md → mlsd`

### 7.2 Run Integration Tests

```bash
cd sdp
poetry run pytest tests/integration/test_github_integration.py -v
```

Expected: All 5 integration tests pass

### 7.3 Check Configuration Summary

```bash
# Print configuration summary
echo "GitHub Integration Setup Summary:"
echo "================================="
echo "Token: $(echo $GITHUB_TOKEN | cut -c1-10)... (${#GITHUB_TOKEN} chars)"
echo "Repo: $GITHUB_REPO"
echo "Org: $GITHUB_ORG"
echo "Hooks: $(ls .git/hooks/{pre,post}-commit commit-msg 2>/dev/null | wc -l)/3 installed"
echo "Tests: $(cd sdp && poetry run pytest tests/integration/ -q --co 2>/dev/null | grep test_github | wc -l) integration tests"
echo ""
echo "✅ Setup complete! Ready to use GitHub integration."
```

## Step 8: First Real Sync (Optional)

### 8.1 Sync a Test Workstream

```bash
# Pick an existing WS file
cd /home/fall_out_bug/msu_ai_masters/sdp

# Sync it to GitHub (creates issue if missing)
poetry run python << 'EOF'
from pathlib import Path
from sdp.github.sync_service import SyncService
from github import Github
import os

client = Github(os.getenv("GITHUB_TOKEN"))
service = SyncService(client)

ws_file = Path("../tools/hw_checker/docs/workstreams/completed/WS-160-09-deploy-integration.md")
result = service.sync_workstream(ws_file)

print(f"✅ Synced: {ws_file.name}")
print(f"   Action: {result.action}")
print(f"   Issue: #{result.issue_number}")
EOF
```

### 8.2 Verify Issue Created

1. Go to: https://github.com/fall-out-bug/msu_ai_masters/issues
2. Find the created issue
3. Check:
   - Title matches WS
   - Labels: `workstream`, `feature/F150`, `size/SMALL`, `status/completed`
   - Milestone: Feature F150
   - Description contains WS Goal

## Troubleshooting Setup Issues

### Issue: Token Invalid

**Error:** `BadCredentialsException: 401`

**Solution:**
1. Verify token is copied correctly (no spaces, full token)
2. Check token hasn't expired
3. Regenerate token with correct scopes

### Issue: Repository Not Found

**Error:** `UnknownObjectException: 404`

**Solution:**
1. Verify `GITHUB_REPO` format: `owner/repo`
2. Check token has access to repository
3. For private repos, ensure token has `repo` scope

### Issue: Hooks Not Executing

**Symptom:** No comments on issues after commits

**Solution:**
1. Check hooks are executable: `chmod +x .git/hooks/*`
2. Verify environment variables in shell
3. Test hook manually: `.git/hooks/post-commit`

### Issue: Rate Limit Exceeded Immediately

**Error:** `RateLimitError` on first request

**Solution:**
- Check you're using authenticated requests (token set)
- Unauthenticated: 60 requests/hour
- Authenticated: 5000 requests/hour

## Next Steps

After setup is complete:

1. **Read Usage Guide:** [USAGE.md](USAGE.md)
2. **Try Workflows:** Follow examples in README
3. **Test with /design:** Create a test feature
4. **Monitor Rate Limits:** Check usage periodically

## Configuration Reference

### Complete .env Template

```bash
# === GitHub Integration (Required) ===
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_REPO=fall-out-bug/msu_ai_masters
GITHUB_ORG=fall-out-bug

# === Project Override (Optional) ===
# Override automatic project routing
# GITHUB_PROJECT=mlsd

# === Telegram Notifications (Optional) ===
# For workflow notifications (feature deployments, etc.)
# TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
# TELEGRAM_CHAT_ID=-1001234567890

# === Logging (Optional) ===
# HW_CHECKER_LOG_LEVEL=DEBUG  # For debugging

# === Dry Run (Optional) ===
# GITHUB_DRY_RUN=true  # Test without making changes
```

## Security Best Practices

1. **Never commit tokens:** Always use `.env` and `.gitignore`
2. **Rotate tokens regularly:** Set expiration and regenerate
3. **Use minimal scopes:** Only grant necessary permissions
4. **Separate tokens for CI/CD:** Use different token for GitHub Actions
5. **Revoke compromised tokens immediately:** GitHub Settings → Tokens → Revoke

## Support

- **Troubleshooting:** [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
- **Usage Examples:** [USAGE.md](USAGE.md)
- **Issues:** https://github.com/fall-out-bug/msu_ai_masters/issues
