# GitHub Issues Automation

Automatic synchronization between SDP workstreams and GitHub issues.

## Features

✅ **Automatic Issue Creation** - `/design` creates GitHub issues for all workstreams
✅ **Bidirectional Sync** - WS file status ↔ GitHub issue status
✅ **Milestones** - Features tracked as GitHub milestones
✅ **Git Integration** - Commits auto-comment on issues
✅ **PR Automation** - Auto-generate PR descriptions with "Closes #N"
✅ **Project Boards** - Kanban visualization (Backlog → In Progress → Done)
✅ **Multi-project** - Separate boards for mlsd, bdde courses
✅ **Error Handling** - Exponential backoff retry logic for rate limits
✅ **Test Coverage** - 189 tests, 100% coverage

## Quick Start

### 1. Setup (one-time)

```bash
# Install git hooks
./sdp/hooks/install-hooks.sh

# Configure environment variables
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_REPO="fall-out-bug/msu_ai_masters"
export GITHUB_ORG="fall-out-bug"
```

### 2. Design workflow

```bash
/design idea-my-feature
# → Creates WS files + GitHub issues + milestone
```

### 3. Manual sync (optional)

```bash
cd sdp

# Sync single workstream (explicit path)
poetry run sdp-github sync-ws ../tools/hw_checker/docs/workstreams/backlog/WS-160-01.md

# Sync all workstreams (backlog/active/completed)
poetry run sdp-github sync-all --ws-dir ../tools/hw_checker/docs/workstreams
# Includes: completed/ subfolders (recursive)
```

### 4. Deploy workflow

```bash
/deploy F160
# → Creates PR with auto-generated description
```

## Documentation

- **[Setup Guide](SETUP.md)** - Initial configuration and token creation
- **[Usage Guide](USAGE.md)** - Commands and workflows
- **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions

## Architecture

```
WS File (markdown)
    ↓ sync
GitHub Issue (with labels, milestone)
    ↓ add to
GitHub Project Board (Kanban)
    ↓ commit
Git Hook → Comment on Issue
    ↓ deploy
PR with "Closes #N" → Auto-close issues
```

### Components

| Component | Purpose | Location |
|-----------|---------|----------|
| **ProjectRouter** | Route WS to correct project (mlsd/bdde) | `sdp/src/sdp/github/project_router.py` |
| **ProjectsClient** | GitHub Projects v2 (GraphQL) API | `sdp/src/sdp/github/projects_client.py` |
| **ProjectBoardSync** | Kanban board automation | `sdp/src/sdp/github/project_board_sync.py` |
| **DeployIntegration** | PR creation with issue linking | `sdp/src/sdp/github/deploy_integration.py` |
| **Exceptions** | Custom error handling | `sdp/src/sdp/github/exceptions.py` |
| **RetryLogic** | Exponential backoff decorators | `sdp/src/sdp/github/retry_logic.py` |

## Workflows

### Complete Feature Workflow

```bash
# 1. Gather requirements
/idea "new feature description"
# → Interactive dialogue
# → docs/drafts/idea-{slug}.md

# 2. Design & create issues
/design idea-{slug}
# → WS-160-01.md, WS-160-02.md, ...
# → GitHub issues #123, #124, ...
# → Milestone "Feature F160: {name}"
# → WS frontmatter updated: github_issue: 123

# 3. Implement each workstream
/build WS-160-01
# → TDD implementation
# → Tests pass
# → Git commit

# Post-commit hook automatically comments on issue #123

# 4. Review all workstreams
/codereview F160
# → 17-point quality checklist
# → APPROVED/CHANGES_REQUESTED

# 5. Human UAT (10-15 min)
# → Follow UAT guide
# → Smoke tests + scenarios

# 6. Deploy to main
/deploy F160
# → PR with "Closes #123, #124, ..."
# → Merge → issues auto-close
```

### Status Synchronization

WS status automatically synced to GitHub:

| WS Status | GitHub Label | Issue State | Board Column |
|-----------|--------------|-------------|--------------|
| backlog | status/backlog | open | Backlog |
| active | status/in-progress | open | In Progress |
| completed | status/completed | closed | Done |
| blocked | status/blocked | open | Blocked |

### Multi-Project Routing

Issues route to projects based on WS ID and file path:

**SDP-related WS** (routed to "MSU AI Masters"):
- WS-190-* (SDP Universal Core)
- WS-191-* (Superpowers Techniques)
- WS-192-* (Multi-Platform)
- WS-193-* (Extension System)
- WS-410+ (SDP Protocol improvements)

**Path-based routing:**
| File Path | Project |
|-----------|---------|
| `sdp/` | MSU AI Masters |
| `tools/hw_checker/` | MSU AI Masters (all hw_checker WS) |
| `courses/mlsd/` | mlsd |
| `courses/bdde/` | bdde |

**Override with:**
- Environment: `export GITHUB_PROJECT=bdde`
- Frontmatter: `github_project: bdde`

## Testing

```bash
cd sdp

# Run all tests
poetry run pytest tests/ -v
# 189 tests, 1.68s execution

# Run integration tests
poetry run pytest tests/integration/test_github_integration.py -v
# 5 integration tests

# Check coverage
poetry run pytest tests/ --cov=sdp.github --cov-report=html
# 100% coverage
```

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `GITHUB_TOKEN` | GitHub API token (repo, project scopes) | Yes | None |
| `GITHUB_REPO` | Repository (format: "org/repo") | Yes | None |
| `GITHUB_ORG` | Organization name | Yes | None |
| `GITHUB_PROJECT` | Override project routing | No | Path-based |

### Token Permissions

GitHub token needs:
- `repo` scope (read/write issues, PRs)
- `project` scope (read/write projects)
- `org:read` scope (read organization projects)

Create token: https://github.com/settings/tokens/new

## Error Handling

### Custom Exceptions

```python
from sdp.github.exceptions import (
    GitHubSyncError,        # Base exception
    RateLimitError,         # Rate limit exceeded
    AuthenticationError,    # Invalid token
    ProjectNotFoundError,   # Project missing
)
```

### Retry Logic

Exponential backoff for rate limits:

```python
from sdp.github.retry_logic import retry_on_rate_limit

@retry_on_rate_limit(max_retries=3, base_delay=1.0)
def api_call():
    # Automatically retries on 403 rate limit
    # Delays: 1s, 2s, 4s
    pass
```

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for:
- Authentication errors
- Rate limit issues
- Project not found
- Hook failures
- Debugging tips

## API Reference

All public classes and methods are fully documented with Google-style docstrings:

```python
from sdp.github.project_router import ProjectRouter

# Get project for workstream file
project = ProjectRouter.get_project_for_ws(ws_file)
# Returns: "mlsd" or "bdde"

# Get all configured projects
projects = ProjectRouter.get_all_projects()
# Returns: ["mlsd", "bdde"]
```

See module docstrings for detailed API documentation:
- `sdp/src/sdp/github/project_router.py`
- `sdp/src/sdp/github/projects_client.py`
- `sdp/src/sdp/github/project_board_sync.py`

## Feature Status

**Feature F150: GitHub Integration**
- Status: ✅ Complete
- Review: ✅ APPROVED
- Tests: 189/189 passing
- Coverage: 100%
- Version: v1.0.0

See:
- Code Review: `tools/hw_checker/docs/reviews/F150-REVIEW.md`
- Release Notes: `docs/releases/v1.0.0-F150.md`
- UAT Guide: `tools/hw_checker/docs/uat/F150-uat-guide.md`

## Contributing

When contributing to GitHub integration:
1. Follow TDD (Red-Green-Refactor)
2. Maintain 100% test coverage
3. Update documentation
4. Follow Clean Architecture (domain → application → infrastructure)
5. Use type hints (mypy --strict)

See: `CONTRIBUTING.md` for general contribution guidelines

## License

Part of the Spec-Driven Protocol (SDP) v0.3.0 project.
