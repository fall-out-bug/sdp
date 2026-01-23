# Usage Guide

Complete guide to using GitHub integration in the Spec-Driven Protocol.

## Overview

GitHub integration provides automatic synchronization between workstream files and GitHub issues, enabling full lifecycle management from planning through deployment.

## Workflows

### Workflow 1: Complete Feature Lifecycle

End-to-end workflow from idea to deployment:

```bash
# === Phase 1: Requirements Gathering ===
/idea "LMS integration for course materials"
# → Interactive dialogue collects requirements
# → Output: docs/drafts/idea-lms-integration.md

# === Phase 2: Design & Planning ===
/design idea-lms-integration
# → Analyzes codebase and creates workstreams
# → Creates WS-060-01.md, WS-060-02.md, WS-060-03.md
# → Creates GitHub issues #123, #124, #125
# → Creates milestone "Feature F60: LMS Integration"
# → Updates WS frontmatter: github_issue: 123

# === Phase 3: Implementation (TDD) ===
/build WS-060-01
# → Red: Write failing test
# → Green: Implement minimum code
# → Refactor: Clean up
# → Tests pass, coverage ≥80%
# → Git commit: "feat(lms): WS-060-01 - implement domain"
# → Post-commit hook comments on issue #123

/build WS-060-02
# → Same TDD cycle
# → Commit comments on issue #124

/build WS-060-03
# → Complete feature implementation

# === Phase 4: Review ===
/codereview F60
# → 17-point quality checklist
# → Verdict: APPROVED or CHANGES_REQUESTED
# → Generates UAT guide

# === Phase 5: Human UAT (10-15 min) ===
# → Follow UAT guide scenarios
# → Sign off if all tests pass

# === Phase 6: Deployment ===
/deploy F60
# → Merges feature to main
# → Creates PR with "Closes #123, #124, #125"
# → Tags v1.1.0
# → Generates release notes
# → Issues auto-close on PR merge
# → Milestone shows 100% complete
```

### Workflow 2: Manual Sync (Ad-hoc)

Sync workstreams to GitHub without using slash commands:

```bash
cd /home/fall_out_bug/msu_ai_masters/sdp

# === Sync Single Workstream (explicit path) ===
poetry run sdp-github sync-ws ../tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# === Sync All Workstreams (backlog/active/completed) ===
poetry run sdp-github sync-all --ws-dir ../tools/hw_checker/docs/workstreams
# Includes: backlog/, active/, completed/ (recursive)

# === Dry-Run (no changes) ===
poetry run sdp-github sync-all --ws-dir ../tools/hw_checker/docs/workstreams --dry-run
```

### Workflow 3: Status Updates

Update workstream status and sync to GitHub:

```bash
# === Update WS Status in Frontmatter ===
# Edit WS file: status: active

# === Sync to GitHub ===
cd sdp
poetry run sdp-github sync-ws ../tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# Status change reflected in:
# - GitHub issue labels (status/in-progress)
# - Project board column (moved to 'In Progress')
```

Status mapping:

| WS Status | GitHub Label | Issue State | Board Column |
|-----------|--------------|-------------|--------------|
| backlog | status/backlog | open | Backlog |
| active | status/in-progress | open | In Progress |
| completed | status/completed | closed | Done |
| blocked | status/blocked | open | Blocked |

### Workflow 4: Multi-Project Management

Route workstreams to different GitHub projects:

```bash
# === Automatic Routing (Path-based) ===
# - tools/hw_checker/ → mlsd
# - courses/mlsd/ → mlsd
# - courses/bdde/ → bdde

# === Manual Project Override ===
# Option 1: Environment variable (session-wide)
export GITHUB_PROJECT=bdde
cd sdp
poetry run python -c "
from pathlib import Path
from sdp.github.project_router import ProjectRouter

ws = Path('../courses/bdde/docs/WS-200-01.md')
project = ProjectRouter.get_project_for_ws(ws)
print(f'Project: {project}')  # → bdde (from environment)
"

# Option 2: Frontmatter override (per-WS)
# Add to WS file frontmatter: github_project: bdde

# === List All Projects ===
poetry run python -c "
from sdp.github.project_router import ProjectRouter
projects = ProjectRouter.get_all_projects()
print('Configured projects:', projects)
"
```

## Git Integration

### Post-Commit Hook (Automatic)

Automatically posts comments to GitHub issues when committing:

```bash
# === Make a Commit with WS ID ===
git commit -m "feat(lms): WS-060-01 - implement domain layer

Add User and Course entities with business logic.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

# === Hook Automatically Posts Comment to Issue #123 ===
# Comment format:
# 🔨 **Commit:** abc1234
# feat(lms): WS-060-01 - implement domain layer
#
# Add User and Course entities with business logic.
#
# [View commit](https://github.com/fall-out-bug/msu_ai_masters/commit/abc1234)
```

### Skip Hook (When Needed)

```bash
# Skip post-commit hook (doesn't post comment)
git commit --no-verify -m "..."

# Or disable temporarily
export SKIP_GITHUB_COMMENT=true
git commit -m "..."
unset SKIP_GITHUB_COMMENT
```

### Dry-Run Hook Testing

```bash
# Test hook without posting
export GITHUB_DRY_RUN=true
.git/hooks/post-commit
# Shows: "Dry-run: would post comment to issue #123"
```

## Commands Reference

### Python API

All operations can be done via Python API:

```python
from github import Github
from pathlib import Path
import os

# === Initialize GitHub Client ===
client = Github(os.getenv("GITHUB_TOKEN"))
repo = client.get_repo(os.getenv("GITHUB_REPO"))

# === Project Router ===
from sdp.github.project_router import ProjectRouter

ws_file = Path("tools/hw_checker/docs/workstreams/backlog/WS-060-01.md")
project = ProjectRouter.get_project_for_ws(ws_file)
# Returns: "mlsd" or "bdde"

all_projects = ProjectRouter.get_all_projects()
# Returns: ["mlsd", "bdde"]

# === Projects Client (GitHub Projects v2 API) ===
from sdp.github.projects_client import ProjectsClient

projects_client = ProjectsClient(client, repo, "fall-out-bug/msu_ai_masters")

# Get or create project
project = projects_client.get_or_create_project("mlsd", "ML System Design")
# Returns: {"id": "PVT_xxx", "number": 1, "title": "mlsd"}

# Add issue to project
item_id = projects_client.add_issue_to_project(project["id"], issue.node_id)

# Update item field (move to column)
projects_client.update_item_field(
    project["id"],
    item_id,
    "Status",
    "In Progress"
)

# === Project Board Sync ===
from sdp.github.project_board_sync import ProjectBoardSync

sync = ProjectBoardSync(client, repo, "fall-out-bug/msu_ai_masters")

# Ensure issue is on board
sync.ensure_issue_on_board(issue, "mlsd")

# Sync WS status to board
ws_file = Path("tools/hw_checker/docs/workstreams/backlog/WS-060-01.md")
sync.sync_ws_status_to_board(ws_file)

# === Sync Service ===
from sdp.github.sync_service import SyncService

service = SyncService(client)

# Sync single workstream
result = service.sync_workstream(ws_file)
# Returns: SyncResult(action="created", issue_number=123, ...)

# === Deploy Integration ===
from sdp.github.deploy_integration import DeployGitHubIntegration

deploy = DeployGitHubIntegration(client)

# Create PR with issue links
pr_result = deploy.create_pr_with_issues(
    base_branch="main",
    head_branch="feature/lms",
    feature_id="F60",
    feature_name="LMS Integration"
)
# Returns: {"pr_url": "...", "description": "Closes #123, #124..."}

# Close milestone if all issues closed
closed = deploy.close_milestone_if_complete("F60")
```

### Error Handling

All operations use custom exceptions with action suggestions:

```python
from sdp.github.exceptions import (
    GitHubSyncError,
    RateLimitError,
    AuthenticationError,
    ProjectNotFoundError,
)

try:
    result = service.sync_workstream(ws_file)
except RateLimitError as e:
    print(f"Rate limit exceeded. Resets at: {e.reset_time}")
    print(f"Action: {e.action}")
    # Wait for reset
except AuthenticationError as e:
    print(f"Authentication failed: {e}")
    print(f"Action: {e.action}")
    # Check GITHUB_TOKEN
except ProjectNotFoundError as e:
    print(f"Project not found: {e}")
    print(f"Action: {e.action}")
    # Create project
except GitHubSyncError as e:
    print(f"Sync error: {e}")
    if e.action:
        print(f"Action: {e.action}")
```

### Retry Logic

Automatic retry with exponential backoff for rate limits:

```python
from sdp.github.retry_logic import retry_on_rate_limit

@retry_on_rate_limit(max_retries=3, base_delay=1.0)
def api_call():
    # Function automatically retries on rate limit (403)
    # Delays: 1s, 2s, 4s
    return client.get_repo("fall-out-bug/msu_ai_masters")

repo = api_call()
```

## Slash Commands Integration

### /design Command

Creates workstreams and GitHub issues:

```bash
/design idea-lms-integration

# Behind the scenes:
# 1. Reads idea draft
# 2. Analyzes codebase
# 3. Creates WS files (WS-060-01.md, WS-060-02.md, ...)
# 4. Creates GitHub milestone "Feature F60: LMS Integration"
# 5. For each WS:
#    - Creates GitHub issue with labels
#    - Adds to milestone
#    - Adds to project board
#    - Updates WS frontmatter: github_issue: 123
# 6. Commits all changes
```

Skip GitHub integration:

```bash
/design idea-lms-integration --skip-github
# Creates WS files only, no GitHub sync
```

### /deploy Command

Creates PR with issue links:

```bash
/deploy F60

# Behind the scenes:
# 1. Generates PR description with "Closes #123, #124, #125"
# 2. Lists all workstreams
# 3. Merges feature branch to main
# 4. Tags version (v1.1.0)
# 5. Generates release notes
# 6. PR merge auto-closes all issues
# 7. Milestone marked complete
```

## Advanced Usage

### Batch Operations

Sync multiple workstreams efficiently:

```python
from pathlib import Path
from sdp.github.sync_service import SyncService
from github import Github
import os

client = Github(os.getenv("GITHUB_TOKEN"))
service = SyncService(client)

# Batch sync all WS for feature
ws_dir = Path("tools/hw_checker/docs/workstreams/backlog")
ws_files = sorted(ws_dir.glob("WS-060-*.md"))

results = []
for ws_file in ws_files:
    try:
        result = service.sync_workstream(ws_file)
        results.append((ws_file.name, result.action, result.issue_number))
        print(f"✅ {ws_file.name}: {result.action} → #{result.issue_number}")
    except Exception as e:
        print(f"❌ {ws_file.name}: {e}")
        results.append((ws_file.name, "error", None))

# Summary
print(f"\nSynced {len([r for r in results if r[1] != 'error'])}/{len(results)} workstreams")
```

### Custom Project Mapping

Add custom project routing:

```python
from sdp.github.project_router import ProjectRouter
from pathlib import Path

# Add custom pattern
ProjectRouter.PATH_PATTERNS["courses/ml_ops"] = "mlops"

# Now this routes to mlops project
ws = Path("courses/ml_ops/docs/WS-300-01.md")
project = ProjectRouter.get_project_for_ws(ws)
# Returns: "mlops"
```

### Monitoring Rate Limits

Check GitHub API rate limit usage:

```python
from github import Github
import os

client = Github(os.getenv("GITHUB_TOKEN"))
rate_limit = client.get_rate_limit()

print(f"Core API:")
print(f"  Remaining: {rate_limit.core.remaining}/{rate_limit.core.limit}")
print(f"  Resets at: {rate_limit.core.reset}")

print(f"GraphQL API:")
print(f"  Remaining: {rate_limit.graphql.remaining}/{rate_limit.graphql.limit}")
```

Rate limits:
- **Authenticated:** 5000 requests/hour (core), 5000 points/hour (GraphQL)
- **Unauthenticated:** 60 requests/hour

## Best Practices

### 1. Always Include WS-ID in Commits

```bash
# ✅ Good: WS-ID in message
git commit -m "feat(lms): WS-060-01 - implement domain"

# ❌ Bad: No WS-ID
git commit -m "feat(lms): implement domain"
```

### 2. Sync Before Deploy

```bash
# Ensure all WS synced before /deploy
cd sdp
poetry run python -c "
from pathlib import Path
from sdp.github.sync_service import SyncService
from github import Github
import os

client = Github(os.getenv('GITHUB_TOKEN'))
service = SyncService(client)

ws_dir = Path('../tools/hw_checker/docs/workstreams/backlog')
for ws_file in ws_dir.glob('WS-060-*.md'):
    service.sync_workstream(ws_file)
"

# Then deploy
/deploy F60
```

### 3. Use Dry-Run for Testing

```bash
# Test sync without changes
export GITHUB_DRY_RUN=true
# ... run sync operations ...
unset GITHUB_DRY_RUN
```

### 4. Monitor Rate Limits

```bash
# Check before large batch operations
curl -H "Authorization: Bearer $GITHUB_TOKEN" \
     https://api.github.com/rate_limit \
     | jq '.rate.remaining'

# If low (<100), wait for reset
```

### 5. Close Milestones After Deployment

```bash
# After /deploy and PR merge
cd sdp
poetry run python -c "
from sdp.github.deploy_integration import DeployGitHubIntegration
from github import Github
import os

client = Github(os.getenv('GITHUB_TOKEN'))
deploy = DeployGitHubIntegration(client)

closed = deploy.close_milestone_if_complete('F60')
if closed:
    print('✅ Milestone closed')
else:
    print('⚠️ Milestone has open issues')
"
```

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for:
- Common errors and solutions
- Debugging techniques
- Rate limit handling
- Hook failures

## See Also

- **Setup Guide:** [SETUP.md](SETUP.md)
- **Overview:** [README.md](README.md)
- **API Reference:** Module docstrings in `sdp/src/sdp/github/`
