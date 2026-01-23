# GitHub Integration

Автоматическая синхронизация Workstreams и Issues с GitHub Projects.

## Setup

### 1. Install GitHub CLI

```bash
# macOS
brew install gh

# Linux
sudo apt install gh

# Authenticate
gh auth login
```

### 2. Configure Repository

```bash
# Set default repository
gh repo set-default your-org/your-repo

# Test
gh issue list
```

### 3. Create GitHub Project

```bash
# Create project board
gh project create \
  --owner "your-org" \
  --title "AI Workflow Automation"

# Add custom fields
gh project field-create \
  --owner "your-org" \
  --project "AI Workflow Automation" \
  --name "WS ID" \
  --data-type TEXT

gh project field-create \
  --name "Feature" \
  --data-type TEXT

gh project field-create \
  --name "Size" \
  --data-type SINGLE_SELECT \
  --single-select-options "SMALL,MEDIUM,LARGE"

gh project field-create \
  --name "Priority" \
  --data-type SINGLE_SELECT \
  --single-select-options "P0,P1,P2,P3"
```

---

## Workflow Integration

### /design Command

Creates GitHub issues for:
- Feature (meta-issue with epic label)
- All Workstreams (linked to feature)

```bash
/design idea-lms-integration

# Creates:
# - Issue #123: [F60] LMS Integration (feature, epic)
# - Issue #124: WS-060-01: Domain entities (workstream)
# - Issue #125: WS-060-02: Repository layer (workstream)
# - ...
```

### /build Command

Updates GitHub issue status:
- Start: `status:backlog` → `status:in-progress`
- Complete: Close issue with comment
- Fail: Add `status:blocked` label

### /issue Command

Creates bug issue:
- Priority label (P0/P1/P2/P3)
- Auto-routing (hotfix/bugfix/backlog)
- Links to affected WS

### /hotfix Command

Creates hotfix issue:
- P0 CRITICAL label
- Closes after deployment
- Links to PR

### /bugfix Command

Creates bugfix issue:
- P1/P2 label
- Closes after merge

---

## GitHub Issue Lifecycle

```
[/design]
    ↓
Create Issue (#123)
status: backlog
    ↓
[/build WS-ID]
    ↓
Update: status:in-progress
    ↓
[Build complete]
    ↓
Close issue
Add comment with metrics
    ↓
[/codereview]
    ↓
Verify all issues closed
```

---

## Issue Templates

Configured in `.github/ISSUE_TEMPLATE/`:

### workstream.yml
```yaml
name: Workstream
description: Workstream created by /design
labels: ["workstream", "auto-generated"]
```

### bug_report.yml
```yaml
name: Bug Report
description: Issue created by /issue
labels: ["bug", "triage"]
```

---

## Project Views

### By Feature
```
Group by: Feature
Sort by: Status
Filter: label:workstream
```

### By Priority (Bugs)
```
Group by: Priority
Sort by: Created
Filter: label:bug
```

### Sprint Board
```
Columns: Backlog | In Progress | Review | Done
Filter: status:open
```

---

## Automation Workflows

See `.github/workflows/`:

### sync-workstreams.yml
- Triggers on WS file changes
- Creates missing GitHub issues
- Updates frontmatter

### update-dashboard.yml
- Runs every 6 hours
- Generates metrics dashboard
- Updates `docs/dashboard/GITHUB_DASHBOARD.md`

---

## GitHub CLI Cheatsheet

```bash
# List issues
gh issue list --label "workstream"
gh issue list --label "P0"

# View issue
gh issue view 123

# Update issue
gh issue edit 123 --add-label "in-progress"
gh issue close 123 --comment "Completed by /build"

# Create PR
gh pr create --title "feat: ..." --body "..."

# Project operations
gh project list
gh project item-list 1
```

---

## Troubleshooting

### "gh: command not found"
```bash
# Install GitHub CLI
brew install gh  # macOS
sudo apt install gh  # Linux
```

### "Permission denied"
```bash
# Re-authenticate
gh auth logout
gh auth login
```

### "Project not found"
```bash
# List projects
gh project list --owner your-org

# Set project number in commands
--project 1
```

---

## Benefits

1. **Visual Tracking:** Kanban board shows progress
2. **Team Visibility:** Everyone sees what's happening
3. **Automated Updates:** Statuses sync automatically
4. **Metrics:** GitHub Insights provides analytics
5. **Integration:** Links between WS, issues, PRs, commits
6. **Notifications:** GitHub notifies on events
7. **Historical Record:** Complete audit trail
