---
ws_id: 00-034-06
feature: F034
status: completed
complexity: MEDIUM
project_id: "00"
---

# Workstream: Add `sdp status` Command

**ID:** 00-034-06  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** MEDIUM (~400 LOC)

---

## Goal

Добавить команду `sdp status`, показывающую текущий контекст: активные workstreams, blockers, next actions.

---

## Context

**Текущая проблема:**

Пользователи не имеют быстрого способа узнать:
- Какой WS сейчас активен?
- Какие WS заблокированы?
- Что делать дальше?

Приходится вручную проверять:
```bash
ls docs/workstreams/in_progress/
ls docs/workstreams/backlog/
sdp guard status  # отдельная команда
bd ready          # ещё одна команда
```

**Решение:** Единая команда `sdp status` с агрегированной информацией.

---

## Scope

### In Scope
- ✅ `sdp status` — human-readable output
- ✅ `sdp status --json` — machine-readable output
- ✅ Show active WS, blockers, ready tasks
- ✅ Show guard status (if active)
- ✅ Show Beads sync status (if available)

### Out of Scope
- ❌ Interactive mode
- ❌ Web dashboard
- ❌ Notifications

---

## Dependencies

**Depends On:**
- None (can start immediately)

**Blocks:**
- None

---

## Acceptance Criteria

- [ ] `sdp status` shows current state
- [ ] `sdp status --json` returns valid JSON
- [ ] Works without Beads installed
- [ ] Shows helpful "next steps" suggestions
- [ ] Handles empty state gracefully

---

## Implementation Plan

### Task 1: Define Status Model

```python
# src/sdp/cli/status/models.py
from dataclasses import dataclass, field
from typing import Optional


@dataclass
class WorkstreamSummary:
    """Summary of a workstream for status display."""
    id: str
    title: str
    status: str
    scope: str
    blockers: list[str] = field(default_factory=list)


@dataclass
class GuardStatus:
    """Guard state summary."""
    active: bool
    workstream_id: Optional[str] = None
    allowed_files: list[str] = field(default_factory=list)


@dataclass
class BeadsStatus:
    """Beads integration status."""
    available: bool
    synced: bool
    ready_tasks: list[str] = field(default_factory=list)
    last_sync: Optional[str] = None


@dataclass
class ProjectStatus:
    """Complete project status."""
    # Workstreams
    in_progress: list[WorkstreamSummary]
    blocked: list[WorkstreamSummary]
    ready: list[WorkstreamSummary]
    
    # Integrations
    guard: GuardStatus
    beads: BeadsStatus
    
    # Suggestions
    next_actions: list[str]
```

### Task 2: Create Status Collector

```python
# src/sdp/cli/status/collector.py
from pathlib import Path
from typing import Optional

from .models import ProjectStatus, WorkstreamSummary, GuardStatus, BeadsStatus


class StatusCollector:
    """Collects project status from various sources."""
    
    def __init__(self, root: Path):
        self.root = root
        self.ws_dir = root / "docs" / "workstreams"
    
    def collect(self) -> ProjectStatus:
        """Collect complete project status."""
        return ProjectStatus(
            in_progress=self._collect_in_progress(),
            blocked=self._collect_blocked(),
            ready=self._collect_ready(),
            guard=self._collect_guard_status(),
            beads=self._collect_beads_status(),
            next_actions=self._suggest_actions(),
        )
    
    def _collect_in_progress(self) -> list[WorkstreamSummary]:
        """Find workstreams in progress."""
        ws_list = []
        in_progress_dir = self.ws_dir / "in_progress"
        if in_progress_dir.exists():
            for ws_file in in_progress_dir.glob("*.md"):
                ws = self._parse_ws_file(ws_file)
                if ws:
                    ws_list.append(ws)
        return ws_list
    
    def _collect_blocked(self) -> list[WorkstreamSummary]:
        """Find blocked workstreams."""
        # Check for workstreams with unmet dependencies
        pass
    
    def _collect_ready(self) -> list[WorkstreamSummary]:
        """Find ready workstreams (no blockers)."""
        pass
    
    def _collect_guard_status(self) -> GuardStatus:
        """Get guard status."""
        state_file = self.root / ".sdp" / "state.json"
        if state_file.exists():
            # Parse and return
            pass
        return GuardStatus(active=False)
    
    def _collect_beads_status(self) -> BeadsStatus:
        """Get Beads integration status."""
        beads_dir = self.root / ".beads"
        if not beads_dir.exists():
            return BeadsStatus(available=False, synced=False)
        # Check sync status, ready tasks
        pass
    
    def _suggest_actions(self) -> list[str]:
        """Suggest next actions based on state."""
        actions = []
        # Logic based on current state
        return actions
```

### Task 3: Create CLI Command

```python
# src/sdp/cli/status/command.py
import json
from pathlib import Path

import typer

from .collector import StatusCollector
from .formatter import format_status_human, format_status_json


app = typer.Typer()


@app.command()
def status(
    json_output: bool = typer.Option(False, "--json", help="Output as JSON"),
    verbose: bool = typer.Option(False, "-v", "--verbose", help="Show more details"),
):
    """Show current project status."""
    root = Path.cwd()
    collector = StatusCollector(root)
    status = collector.collect()
    
    if json_output:
        typer.echo(format_status_json(status))
    else:
        typer.echo(format_status_human(status, verbose=verbose))
```

### Task 4: Create Output Formatter

```python
# src/sdp/cli/status/formatter.py
import json
from dataclasses import asdict

from rich.console import Console
from rich.table import Table
from rich.panel import Panel

from .models import ProjectStatus


def format_status_human(status: ProjectStatus, verbose: bool = False) -> str:
    """Format status for human reading."""
    console = Console(record=True)
    
    # Header
    console.print(Panel("[bold]SDP Project Status[/bold]"))
    
    # In Progress
    if status.in_progress:
        console.print("\n[yellow]⏳ In Progress[/yellow]")
        for ws in status.in_progress:
            console.print(f"  • {ws.id}: {ws.title}")
    else:
        console.print("\n[dim]No workstreams in progress[/dim]")
    
    # Ready
    if status.ready:
        console.print("\n[green]✅ Ready to Start[/green]")
        for ws in status.ready:
            console.print(f"  • {ws.id}: {ws.title}")
    
    # Blocked
    if status.blocked:
        console.print("\n[red]🚫 Blocked[/red]")
        for ws in status.blocked:
            blockers = ", ".join(ws.blockers) if ws.blockers else "unknown"
            console.print(f"  • {ws.id}: {ws.title} (by: {blockers})")
    
    # Guard Status
    console.print("\n[blue]🛡️ Guard[/blue]")
    if status.guard.active:
        console.print(f"  Active: {status.guard.workstream_id}")
    else:
        console.print("  [dim]Inactive[/dim]")
    
    # Beads Status
    if status.beads.available:
        console.print("\n[magenta]📿 Beads[/magenta]")
        sync_status = "✅ Synced" if status.beads.synced else "⚠️ Needs sync"
        console.print(f"  {sync_status}")
        if status.beads.ready_tasks:
            console.print(f"  Ready: {len(status.beads.ready_tasks)} tasks")
    
    # Next Actions
    if status.next_actions:
        console.print("\n[bold cyan]💡 Suggested Actions[/bold cyan]")
        for action in status.next_actions:
            console.print(f"  → {action}")
    
    return console.export_text()


def format_status_json(status: ProjectStatus) -> str:
    """Format status as JSON."""
    return json.dumps(asdict(status), indent=2)
```

### Task 5: Register Command

```python
# src/sdp/cli/main.py
from .status.command import app as status_app

# Add to main CLI
app.add_typer(status_app, name="status")
```

### Task 6: Add Tests

```python
# tests/unit/cli/test_status.py
import pytest
from pathlib import Path
from sdp.cli.status.collector import StatusCollector
from sdp.cli.status.models import ProjectStatus


def test_status_collector_empty_project(tmp_path):
    """Test status collection on empty project."""
    collector = StatusCollector(tmp_path)
    status = collector.collect()
    
    assert status.in_progress == []
    assert status.ready == []
    assert status.blocked == []
    assert status.guard.active is False


def test_status_collector_with_workstreams(tmp_path):
    """Test status collection with workstreams."""
    # Setup workstream files
    ws_dir = tmp_path / "docs" / "workstreams" / "in_progress"
    ws_dir.mkdir(parents=True)
    (ws_dir / "00-034-01.md").write_text("---\nws_id: 00-034-01\n---\n# Test")
    
    collector = StatusCollector(tmp_path)
    status = collector.collect()
    
    assert len(status.in_progress) == 1
    assert status.in_progress[0].id == "00-034-01"


def test_status_json_output(tmp_path):
    """Test JSON output format."""
    from sdp.cli.status.formatter import format_status_json
    
    status = ProjectStatus(
        in_progress=[],
        blocked=[],
        ready=[],
        guard=GuardStatus(active=False),
        beads=BeadsStatus(available=False, synced=False),
        next_actions=["Run @build to start"],
    )
    
    output = format_status_json(status)
    import json
    parsed = json.loads(output)
    
    assert "next_actions" in parsed
    assert parsed["next_actions"] == ["Run @build to start"]
```

---

## Expected Output

### Human-Readable

```
╭──────────────────────────────────╮
│       SDP Project Status         │
╰──────────────────────────────────╯

⏳ In Progress
  • 00-034-01: Split Large Files (Phase 1: Core)

✅ Ready to Start
  • 00-034-04: Documentation Consistency

🚫 Blocked
  • 00-034-02: Split Large Files (Phase 2) (by: 00-034-01)
  • 00-034-03: Increase Test Coverage (by: 00-034-02)

🛡️ Guard
  Active: 00-034-01

📿 Beads
  ✅ Synced
  Ready: 2 tasks

💡 Suggested Actions
  → Complete 00-034-01 to unblock 00-034-02
  → Run `bd sync` after completing current WS
```

### JSON

```json
{
  "in_progress": [
    {"id": "00-034-01", "title": "Split Large Files", "status": "in_progress"}
  ],
  "blocked": [
    {"id": "00-034-02", "title": "Split Large Files (Phase 2)", "blockers": ["00-034-01"]}
  ],
  "ready": [
    {"id": "00-034-04", "title": "Documentation Consistency"}
  ],
  "guard": {"active": true, "workstream_id": "00-034-01"},
  "beads": {"available": true, "synced": true, "ready_tasks": ["bd-0001", "bd-0002"]},
  "next_actions": ["Complete 00-034-01 to unblock 00-034-02"]
}
```

---

## DO / DON'T

### CLI Design

**✅ DO:**
- Use Rich for formatted output
- Provide JSON option for scripting
- Handle missing files gracefully
- Show helpful suggestions

**❌ DON'T:**
- Fail on missing optional integrations (Beads)
- Show too much detail by default
- Use colors that don't work on all terminals

---

## Files to Create

- [ ] `src/sdp/cli/status/__init__.py`
- [ ] `src/sdp/cli/status/models.py`
- [ ] `src/sdp/cli/status/collector.py`
- [ ] `src/sdp/cli/status/formatter.py`
- [ ] `src/sdp/cli/status/command.py`
- [ ] `tests/unit/cli/test_status.py`

## Files to Modify

- [ ] `src/sdp/cli/main.py` — register status command

---

## Test Plan

### Unit Tests
- [ ] StatusCollector handles empty project
- [ ] StatusCollector finds in_progress WS
- [ ] StatusCollector detects blocked WS
- [ ] JSON output is valid
- [ ] Human output is readable

### Integration Tests
- [ ] `sdp status` runs without error
- [ ] `sdp status --json` returns valid JSON
- [ ] Works with and without Beads

---

**Version:** 1.0  
**Created:** 2026-01-31
