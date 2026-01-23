# Troubleshooting Guide

Solutions to common issues with GitHub integration.

## Common Issues

### Authentication Error

**Error:**
```
GitHubSyncError: GitHub authentication failed
Action: Verify GITHUB_TOKEN in .env has correct permissions
```

**Symptoms:**
- 401 Unauthorized responses
- "Bad credentials" messages
- Cannot access repository

**Solutions:**

1. **Check Token is Set:**
   ```bash
   echo $GITHUB_TOKEN | head -c 10
   # Should show: ghp_xxxxxx
   ```

2. **Verify Token is Valid:**
   ```bash
   curl -H "Authorization: Bearer $GITHUB_TOKEN" \
        https://api.github.com/user
   # Should return your GitHub user info
   ```

3. **Check Token Scopes:**
   - Go to: https://github.com/settings/tokens
   - Find your token
   - Verify scopes: `repo`, `project`, `org:read`

4. **Regenerate Token:**
   - Delete old token
   - Create new token with correct scopes
   - Update `.env` file

5. **Check .env is Loaded:**
   ```bash
   # Source .env
   set -a; source .env; set +a

   # Verify
   env | grep GITHUB
   ```

---

### Rate Limit Exceeded

**Error:**
```
RateLimitError: GitHub API rate limit exceeded. Resets at 2026-01-18 15:30:00
Action: Wait for rate limit reset or use different token
```

**Symptoms:**
- 403 Forbidden responses with "rate limit" message
- Sync operations fail after many requests
- Error shows reset time

**Solutions:**

1. **Check Current Rate Limit:**
   ```bash
   curl -H "Authorization: Bearer $GITHUB_TOKEN" \
        https://api.github.com/rate_limit | jq
   ```

2. **Wait for Reset:**
   - Note the reset time in error message
   - Wait until that time
   - Rate limits reset hourly

3. **Use Different Token:**
   - Create second Personal Access Token
   - Use for different operations
   - Each token has separate rate limit

4. **Reduce Request Frequency:**
   - Batch sync operations
   - Use dry-run mode for testing
   - Cache project/repo data

5. **Automatic Retry:**
   - Integration has built-in exponential backoff
   - Automatically retries on rate limit
   - Max 3 retries with delays: 1s, 2s, 4s

**Rate Limits (telemetry):**
- **Authenticated:** 5000 requests/hour
- **Unauthenticated:** 60 requests/hour
- **GraphQL:** 5000 points/hour

---

### Project Not Found

**Error:**
```
ProjectNotFoundError: GitHub project 'mlsd' not found
Action: Create project 'mlsd' in GitHub or run auto-setup
```

**Symptoms:**
- "Project not found" error
- Issues not added to board
- Project operations fail

**Solutions:**

1. **Check Project Exists:**
   - Go to: https://github.com/fall-out-bug/msu_ai_masters/projects
   - Verify project "mlsd" (or "bdde") exists

2. **Create Project Manually:**
   ```
   1. Go to repository projects
   2. Click "New project"
   3. Choose "Board" template
   4. Name: "mlsd"
   5. Add columns: Backlog, In Progress, In Review, Done
   ```

3. **Auto-Create Project:**
   ```python
   from github import Github
   from sdp.github.projects_client import ProjectsClient
   import os

   client = Github(os.getenv("GITHUB_TOKEN"))
   repo = client.get_repo(os.getenv("GITHUB_REPO"))
   projects = ProjectsClient(client, repo, os.getenv("GITHUB_REPO"))

   # Auto-creates if missing
   project = projects.get_or_create_project("mlsd", "ML System Design")
   print(f"Created: {project['title']}")
   ```

4. **Check Token Has Project Scope:**
   - Token needs `project` scope
   - Regenerate if missing

---

### Issue Not Found

**Error:**
```
github.GithubException.UnknownObjectException: 404
{"message": "Not Found", "documentation_url": "..."}
```

**Symptoms:**
- Cannot update existing issue
- WS has github_issue number but sync fails
- 404 Not Found error

**Solutions:**

1. **Check Issue Number in Frontmatter:**
   ```bash
   head -n 10 tools/hw_checker/docs/workstreams/backlog/WS-060-01.md
   # Look for: github_issue: 123
   ```

2. **Verify Issue Exists:**
   - Go to: https://github.com/fall-out-bug/msu_ai_masters/issues/123
   - Check issue hasn't been deleted

3. **Recreate Issue:**
   ```bash
   # Edit WS frontmatter: github_issue: null

   # Re-sync (creates new issue)
   cd sdp
   poetry run sdp-github sync-ws ../tools/hw_checker/docs/workstreams/backlog/WS-060-01.md
   ```

4. **Check Repository Access:**
   - Verify token has access to repository
   - For private repos, ensure `repo` scope

---

### Post-Commit Hook Fails Silently

**Symptom:** No comments appear on GitHub issues after commits

**Symptoms:**
- Commits don't trigger issue comments
- Hook doesn't execute
- No errors shown

**Solutions:**

1. **Check Hook is Installed:**
   ```bash
   ls -la .git/hooks/post-commit
   # Should exist and be executable
   ```

2. **Verify Hook Permissions:**
   ```bash
   chmod +x .git/hooks/post-commit
   ```

3. **Check Environment Variables:**
   ```bash
   # Hook needs these set
   echo $GITHUB_TOKEN
   echo $GITHUB_REPO

   # If empty, source .env before committing
   set -a; source .env; set +a
   ```

4. **Test Hook Manually:**
   ```bash
   # Dry-run test
   export GITHUB_DRY_RUN=true
   .git/hooks/post-commit

   # Should show: "Dry-run: would post comment to issue #X"
   ```

5. **Check Hook Logs:**
   ```bash
   # Enable debug mode
   export HW_CHECKER_LOG_LEVEL=DEBUG
   git commit -m "test: check hook"

   # Check for errors in output
   ```

6. **Re-install Hooks:**
   ```bash
   ./sdp/hooks/install-hooks.sh
   ```

7. **Check WS-ID in Commit Message:**
   ```bash
   # ✅ Good: Contains WS-ID
   git commit -m "feat(lms): WS-060-01 - implement domain"

   # ❌ Bad: No WS-ID (hook won't find issue)
   git commit -m "feat(lms): implement domain"
   ```

---

### Dry-Run Shows Plan But Doesn't Execute

**Expected behavior!** Dry-run mode never makes actual changes.

**Symptoms:**
- Sync shows "would create issue" but nothing happens
- Commands don't create issues/projects
- No errors, but no results

**Solution:**

Remove `--dry-run` flag or `GITHUB_DRY_RUN` environment variable:

```bash
# Remove environment variable
unset GITHUB_DRY_RUN

# Run actual sync (without dry-run)
cd sdp
poetry run python -c "
from pathlib import Path
from sdp.github.sync_service import SyncService
from github import Github
import os

client = Github(os.getenv('GITHUB_TOKEN'))
service = SyncService(client)

ws_file = Path('../tools/hw_checker/docs/workstreams/backlog/WS-060-01.md')
result = service.sync_workstream(ws_file)  # No dry-run flag
print(f'Created issue: #{result.issue_number}')
"
```

---

### GraphQL API Errors

**Error:**
```
github.GithubException.GithubException: 502
{"message": "Server Error"}
```

**Symptoms:**
- Projects v2 operations fail
- GraphQL queries timeout
- 502/503 errors

**Solutions:**

1. **Check GitHub Status:**
   - Go to: https://www.githubstatus.com/
   - Look for API incidents

2. **Retry Operation:**
   - Built-in retry logic handles transient errors
   - Wait a few minutes and try again

3. **Use REST API Fallback:**
   - Some operations have REST alternatives
   - Projects v2 requires GraphQL (no fallback)

4. **Reduce Payload Size:**
   - Query fewer fields
   - Batch operations smaller

---

### Milestone Not Closing

**Symptom:** Milestone stays open after all issues closed

**Symptoms:**
- All issues closed but milestone open
- `/deploy` doesn't close milestone
- Manual check needed

**Solutions:**

1. **Check All Issues Closed:**
   ```bash
   # Go to milestone page
   # https://github.com/fall-out-bug/msu_ai_masters/milestone/X

   # Verify: "0 open, N closed"
   ```

2. **Manually Close Milestone:**
   ```python
   from sdp.github.deploy_integration import DeployGitHubIntegration
   from github import Github
   import os

   client = Github(os.getenv("GITHUB_TOKEN"))
   deploy = DeployGitHubIntegration(client)

   closed = deploy.close_milestone_if_complete("F60")
   if not closed:
       print("Some issues still open")
   ```

3. **Force Close:**
   - Go to milestone page on GitHub
   - Click "Close milestone" button

---

### Wrong Project Assignment

**Symptom:** WS synced to wrong project (mlsd vs bdde)

**Symptoms:**
- Issue appears in wrong project board
- Routing incorrect

**Solutions:**

1. **Check Path-Based Routing:**
   ```python
   from pathlib import Path
   from sdp.github.project_router import ProjectRouter

   ws = Path("courses/bdde/docs/WS-200-01.md")
   project = ProjectRouter.get_project_for_ws(ws)
   print(f"Routes to: {project}")
   # Should be "bdde" for courses/bdde/
   ```

2. **Override in Frontmatter:**
   ```yaml
   ---
   ws_id: WS-200-01
   github_project: bdde
   ---
   ```

3. **Override via Environment:**
   ```bash
   export GITHUB_PROJECT=bdde
   # Sync operations now use bdde project
   ```

4. **Add Custom Path Pattern:**
   ```python
   from sdp.github.project_router import ProjectRouter

   # Add custom pattern
   ProjectRouter.PATH_PATTERNS["new/path"] = "new_project"
   ```

---

## Debugging

### Enable Debug Logging

```bash
export HW_CHECKER_LOG_LEVEL=DEBUG
cd sdp
poetry run python -c "
from pathlib import Path
from sdp.github.sync_service import SyncService
from github import Github
import os

client = Github(os.getenv('GITHUB_TOKEN'))
service = SyncService(client)

ws_file = Path('../tools/hw_checker/docs/workstreams/backlog/WS-060-01.md')
result = service.sync_workstream(ws_file)
"
# Verbose output with API calls, responses, etc.
```

### Test GitHub Connection

```bash
# Test authentication
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/user \
     | jq '.login'

# Check rate limit
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/rate_limit \
     | jq '.rate'

# Test repository access
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/repos/fall-out-bug/msu_ai_masters \
     | jq '.full_name'
```

### Inspect WS Frontmatter

```bash
# Check frontmatter format
head -n 15 tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# Should have:
# ---
# ws_id: WS-060-01
# feature: F60
# status: backlog
# github_issue: 123 (or null)
# ---
```

### Validate Token Scopes

```bash
# Get token info (shows scopes)
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/user \
     -I \
     | grep -i "x-oauth-scopes"

# Should include: repo, project, org:read
```

### Test Import Paths

```bash
cd sdp
poetry run python -c "
# Test all imports work
from sdp.github.project_router import ProjectRouter
from sdp.github.projects_client import ProjectsClient
from sdp.github.project_board_sync import ProjectBoardSync
from sdp.github.sync_service import SyncService
from sdp.github.deploy_integration import DeployGitHubIntegration
from sdp.github.exceptions import GitHubSyncError
from sdp.github.retry_logic import retry_on_rate_limit

print('✅ All imports successful')
"
```

### Run Integration Tests

```bash
cd sdp

# Run all integration tests
poetry run pytest tests/integration/test_github_integration.py -v

# Run specific test
poetry run pytest tests/integration/test_github_integration.py::test_project_router_integration -v

# Run with coverage
poetry run pytest tests/integration/ --cov=sdp.github --cov-report=term-missing
```

---

## Error Code Reference

| Error Code | Meaning | Common Cause |
|------------|---------|--------------|
| 401 | Unauthorized | Invalid/expired token |
| 403 | Forbidden | Rate limit or insufficient scopes |
| 404 | Not Found | Issue/project/repo doesn't exist |
| 422 | Unprocessable Entity | Invalid input data |
| 502/503 | Server Error | GitHub API down/overloaded |

---

## Getting Help

### Before Asking for Help

1. **Check Logs:**
   - Enable `HW_CHECKER_LOG_LEVEL=DEBUG`
   - Look for error details

2. **Verify Setup:**
   - Token set correctly
   - Repository accessible
   - Hooks installed

3. **Test Connection:**
   - Run `curl` tests above
   - Verify rate limit not exceeded

4. **Run Tests:**
   - `pytest tests/integration/`
   - Check for failures

### How to Report Issues

Include in bug report:

1. **Error Message:**
   ```
   Exact error text with stack trace
   ```

2. **Environment:**
   ```bash
   python --version
   poetry --version
   git --version
   ```

3. **Configuration:**
   ```bash
   # Don't include actual token!
   echo "GITHUB_REPO: $GITHUB_REPO"
   echo "Token length: ${#GITHUB_TOKEN} chars"
   ```

4. **Steps to Reproduce:**
   ```bash
   # Exact commands that cause error
   cd sdp
   poetry run python -c "..."
   ```

5. **Expected vs Actual:**
   - What you expected to happen
   - What actually happened

### Resources

- **GitHub API Status:** https://www.githubstatus.com/
- **GitHub API Docs:** https://docs.github.com/en/rest
- **Projects v2 Docs:** https://docs.github.com/en/issues/planning-and-tracking-with-projects
- **Rate Limits:** https://docs.github.com/en/rest/rate-limit
- **Token Scopes:** https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps

### Contact

- **Issues:** https://github.com/fall-out-bug/msu_ai_masters/issues
- **Documentation:** [README.md](README.md), [SETUP.md](SETUP.md), [USAGE.md](USAGE.md)
