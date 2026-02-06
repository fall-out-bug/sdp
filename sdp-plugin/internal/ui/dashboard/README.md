# TUI Dashboard (`sdp status`)

Rich terminal UI (TUI) dashboard for real-time project status visibility.

## Features

- **Workstreams Tab**: View all tasks grouped by status (Open, In Progress, Completed, Blocked)
- **Ideas Tab**: Browse draft ideas from `docs/drafts/`
- **Tests Tab**: Monitor test coverage and quality gates
- **Activity Tab**: View recent activity (placeholder)

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `1` | Switch to Workstreams tab |
| `2` | Switch to Ideas tab |
| `3` | Switch to Tests tab |
| `4` | Switch to Activity tab |
| `r` | Refresh all data |
| `q` | Quit dashboard |

## Color Coding

### Status Colors
- **Green** (76): Open tasks
- **Yellow** (226): In Progress
- **Blue** (69): Completed
- **Red** (196): Blocked

### Priority Colors
- **Red** (196): P0 - Critical
- **Orange** (208): P1 - High
- **Yellow** (226): P2 - Medium
- **Gray** (faint): P3 - Low

## Usage

```bash
# Start the dashboard
sdp status

# The dashboard will auto-refresh every 2 seconds
# Press 'r' to force refresh
# Press 'q' to quit
```

## Data Sources

The dashboard fetches data from multiple sources:

1. **Beads CLI** (if available):
   - Uses `bd ready` to fetch open tasks
   - Shows task ID, title, status, priority

2. **docs/workstreams/** (fallback):
   - Parses workstream markdown files
   - Extracts ID, goal, status, size from frontmatter

3. **docs/drafts/**:
   - Lists all markdown files
   - Shows title, path, modification time

4. **Test results** (placeholder):
   - Coverage: N/A (TODO: read from coverage files)
   - Quality Gates: N/A (TODO: run quality checks)

## Architecture

```
cmd/sdp/
├── status.go          # CLI command entry point

internal/ui/dashboard/
├── app.go             # Bubbletea app (Model, Update, View)
├── state.go           # State structures
├── styles.go          # Lipgloss styles for colors
├── data.go            # Data fetching from Beads, docs
└── app_test.go        # Unit tests
```

## Tech Stack

- **bubbletea**: TUI framework (Model-Update-View)
- **lipgloss**: Styling and layout
- **fsnotify**: File watcher (for auto-refresh)

## Future Enhancements

- [ ] Implement Activity tab with real activity log
- [ ] Add test coverage parsing from coverage files
- [ ] Add quality gate status from `sdp quality check`
- [ ] Implement file watcher for auto-refresh
- [ ] Add filters (by status, priority, assignee)
- [ ] Add sorting options
- [ ] Show workstream details on select
- [ ] Navigate to workstream file with Enter key
