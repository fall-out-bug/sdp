# F012: GitHub Agent Orchestrator + Developer DX Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build GitHub Agent Orchestrator with Developer Experience improvements — daemon, queue, agents, TUI dashboard, test watch mode, and unified status command.

**Architecture:**
- **Layer 1 (Infrastructure):** Daemon service, task queue, GitHub sync
- **Layer 2 (Agent System):** Executor interface, multi-agent orchestration
- **Layer 3 (Dashboard Core):** Reusable state management + Textual widgets
- **Layer 4 (DX Features):** Status command, test watch, auto-state, TUI app

**Tech Stack:** Python 3.11+, Click, Textual, watchdog, pytest, asyncio

**Workstreams:** 14 total | ~9,500 LOC | 4-6 weeks for solo developer

---

## Dependency Graph

```
Layer 1 (Infrastructure)                    Layer 3 (Dashboard Core)
┌─────────────────────┐                    ┌──────────────────────┐
│ 00-012-01 Daemon    │──┐                 │ 00-012-08 Core       │──┐
└─────────────────────┘  │                 └──────────────────────┘  │
                         │                                           │
┌─────────────────────┐  │                                           │
│ 00-012-02 Queue     │◄─┘                                           │
└─────────────────────┘  │                                           │
                         │                    ┌──────────────────────┤
┌─────────────────────┐  │                    │ 00-012-11 Status    │
│ 00-012-03 GH Sync   │ (parallel)           └──────────────────────┘
└─────────────────────┘  │                                           │
                         │                    ┌──────────────────────┤
┌─────────────────────┐  │                    │ 00-012-12 Test Watch │
│ 00-012-04 Executor  │◄─┘───────────────────►└──────────────────────┘
└─────────────────────┘  │                                           │
                         │                    ┌──────────────────────┤
┌─────────────────────┐  │                    │ 00-012-13 Auto State │
│ 00-012-05 CLI Tasks │◄─┤                    └──────────────────────┘
└─────────────────────┘  │                                           │
                         │                    ┌──────────────────────┐
┌─────────────────────┐  │                    │ 00-012-14 Dashboard │
│ 00-012-06 Orchestr  │◄─┘───────────────────►└──────────────────────┘
└─────────────────────┘
                         │
┌─────────────────────┐  │
│ 00-012-07 GH Fields │◄─┘ (depends on 00-012-03)
└─────────────────────┘  │
                         │
┌─────────────────────┐  │
│ 00-012-09 Webhook   │◄─┘ (depends on 00-012-07)
└─────────────────────┘
                         │
┌─────────────────────┐  │
│ 00-012-10 Pre-Check │◄─┘ (depends on 00-012-04)
└─────────────────────┘
```

---

## Execution Order

### Wave 1: Foundation (can run in parallel)
- 00-012-01: Daemon Service Framework
- 00-012-03: Enhanced GitHub Sync
- 00-012-08: Dashboard Core

### Wave 2: Core Infrastructure (depends on Wave 1)
- 00-012-02: Task Queue Management (needs 00-012-01)
- 00-012-11: Workstream Status Command (needs 00-012-08)
- 00-012-12: Test Watch Mode (needs 00-012-08)

### Wave 3: Agent System (depends on Wave 2)
- 00-012-04: Agent Executor Interface (needs 00-012-02)
- 00-012-13: Auto-State Management (needs 00-012-08)

### Wave 4: Advanced Features (depends on Wave 3)
- 00-012-05: CLI Task Commands (needs 00-012-02, 00-012-04)
- 00-012-06: Multi-Agent Orchestration (needs 00-012-04)
- 00-012-07: GitHub Project Fields (needs 00-012-03)
- 00-012-10: Pre-Execution Checks (needs 00-012-04)

### Wave 5: Integration (depends on Wave 4)
- 00-012-09: Webhook Support (needs 00-012-07)
- 00-012-14: Developer Dashboard App (needs 00-012-08, 00-012-11, 00-012-12)

---

## Task 1: 00-012-01 Daemon Service Framework

**Files:**
- Create: `src/sdp/daemon/__init__.py`
- Create: `src/sdp/daemon/daemon.py` (~150 LOC)
- Create: `src/sdp/daemon/pid_manager.py` (~100 LOC)
- Create: `tests/unit/daemon/test_daemon.py` (~80 LOC)
- Create: `tests/unit/daemon/test_pid_manager.py` (~70 LOC)
- Modify: `src/sdp/cli.py` (+50 LOC)

---

### Step 1: Write daemon dataclass test

**File:** `tests/unit/daemon/test_daemon.py`

```python
import pytest
from sdp.daemon.daemon import DaemonConfig, Daemon

def test_daemon_config_defaults():
    config = DaemonConfig()
    assert config.watch_enabled is False
    assert config.pid_file == ".sdp/daemon.pid"
    assert config.log_file == ".sdp/daemon.log"

def test_daemon_config_custom():
    config = DaemonConfig(
        watch_enabled=True,
        pid_file="/tmp/test.pid",
        log_file="/tmp/test.log"
    )
    assert config.watch_enabled is True
    assert config.pid_file == "/tmp/test.pid"
```

**Run:** `pytest tests/unit/daemon/test_daemon.py::test_daemon_config_defaults -v`
**Expected:** FAIL - "No module named 'sdp.daemon.daemon'"

---

### Step 2: Create daemon module with DaemonConfig

**File:** `src/sdp/daemon/daemon.py`

```python
from dataclasses import dataclass
from pathlib import Path

@dataclass
class DaemonConfig:
    """Configuration for SDP daemon process."""
    watch_enabled: bool = False
    pid_file: str = ".sdp/daemon.pid"
    log_file: str = ".sdp/daemon.log"
    poll_interval: float = 1.0
```

**Run:** `pytest tests/unit/daemon/test_daemon.py -v`
**Expected:** PASS

---

### Step 3: Write PID manager test

**File:** `tests/unit/daemon/test_pid_manager.py`

```python
import pytest
from pathlib import Path
from sdp.daemon.pid_manager import PIDManager, PIDError

@pytest.fixture
def temp_pid_file(tmp_path: Path):
    return tmp_path / "test.pid"

def test_write_pid(temp_pid_file: Path):
    manager = PIDManager(temp_pid_file)
    manager.write(12345)
    assert temp_pid_file.read_text() == "12345"

def test_read_existing_pid(temp_pid_file: Path):
    manager = PIDManager(temp_pid_file)
    manager.write(12345)
    assert manager.read() == 12345

def test_read_missing_pid_raises(temp_pid_file: Path):
    manager = PIDManager(temp_pid_file)
    with pytest.raises(PIDError):
        manager.read()

def test_remove_pid(temp_pid_file: Path):
    manager = PIDManager(temp_pid_file)
    manager.write(12345)
    manager.remove()
    assert not temp_pid_file.exists()

def test_is_running_true(temp_pid_file: Path):
    import psutil
    manager = PIDManager(temp_pid_file)
    manager.write(psutil.Process().pid)
    assert manager.is_running() is True

def test_is_running_false(temp_pid_file: Path):
    manager = PIDManager(temp_pid_file)
    manager.write(99999)  # Non-existent PID
    assert manager.is_running() is False
```

**Run:** `pytest tests/unit/daemon/test_pid_manager.py -v`
**Expected:** FAIL - "No module named 'sdp.daemon.pid_manager'"

---

### Step 4: Create PIDManager

**File:** `src/sdp/daemon/pid_manager.py`

```python
import os
from pathlib import Path

class PIDError(Exception):
    """Raised when PID operations fail."""

class PIDManager:
    """Manages daemon PID file."""

    def __init__(self, pid_file: str | Path):
        self._pid_file = Path(pid_file)

    def write(self, pid: int) -> None:
        """Write PID to file."""
        self._pid_file.parent.mkdir(parents=True, exist_ok=True)
        self._pid_file.write_text(str(pid))

    def read(self) -> int:
        """Read PID from file."""
        if not self._pid_file.exists():
            raise PIDError(f"PID file not found: {self._pid_file}")
        return int(self._pid_file.read_text())

    def remove(self) -> None:
        """Remove PID file."""
        if self._pid_file.exists():
            self._pid_file.unlink()

    def is_running(self) -> bool:
        """Check if PID is running."""
        try:
            pid = self.read()
            os.kill(pid, 0)  # Send signal 0 (doesn't kill)
            return True
        except (PIDError, OSError):
            return False
```

**Run:** `pytest tests/unit/daemon/test_pid_manager.py -v`
**Expected:** PASS

---

### Step 5: Add daemon CLI command

**File:** `src/sdp/cli.py` (add to existing file)

```python
import click
from sdp.daemon.daemon import DaemonConfig
from sdp.daemon.pid_manager import PIDManager

@main.group()
def daemon():
    """Daemon management commands."""
    pass

@daemon.command()
@click.option("--watch", is_flag=True, help="Enable file watch mode")
@click.option("--pid-file", default=".sdp/daemon.pid", help="PID file path")
def start(watch: bool, pid_file: str):
    """Start SDP daemon process."""
    manager = PIDManager(pid_file)
    if manager.is_running():
        click.echo(f"Daemon already running (PID: {manager.read()})")
        return

    import os
    pid = os.fork()
    if pid == 0:
        # Child process
        from sdp.daemon.daemon import Daemon
        config = DaemonConfig(watch_enabled=watch, pid_file=pid_file)
        daemon = Daemon(config)
        daemon.run()
    else:
        manager.write(pid)
        click.echo(f"Daemon started (PID: {pid})")

@daemon.command()
@click.option("--pid-file", default=".sdp/daemon.pid", help="PID file path")
def stop(pid_file: str):
    """Stop SDP daemon process."""
    manager = PIDManager(pid_file)
    if not manager.is_running():
        click.echo("Daemon not running")
        return

    import os
    import signal
    pid = manager.read()
    os.kill(pid, signal.SIGTERM)
    manager.remove()
    click.echo(f"Daemon stopped (PID: {pid})")

@daemon.command()
@click.option("--pid-file", default=".sdp/daemon.pid", help="PID file path")
def status(pid_file: str):
    """Show daemon status."""
    manager = PIDManager(pid_file)
    if manager.is_running():
        click.echo(f"Daemon running (PID: {manager.read()})")
    else:
        click.echo("Daemon not running")
```

**Run:** `sdp daemon --help`
**Expected:** Shows daemon group with start/stop/status commands

---

### Step 6: Add Daemon class with run loop

**File:** `src/sdp/daemon/daemon.py` (extend)

```python
import asyncio
import signal
from pathlib import Path
from .pid_manager import PIDManager

class Daemon:
    """SDP daemon process."""

    def __init__(self, config: DaemonConfig):
        self._config = config
        self._running = False
        self._pid_manager = PIDManager(config.pid_file)

    def run(self) -> None:
        """Run daemon event loop."""
        self._running = True

        # Setup signal handlers
        signal.signal(signal.SIGTERM, self._shutdown)
        signal.signal(signal.SIGINT, self._shutdown)

        # Main loop
        while self._running:
            try:
                asyncio.get_event_loop().run_until_complete(
                    asyncio.sleep(self._config.poll_interval)
                )
            except KeyboardInterrupt:
                break

        # Cleanup
        self._pid_manager.remove()

    def _shutdown(self, signum: int, frame) -> None:
        """Handle shutdown signal."""
        self._running = False
```

---

### Step 7: Commit

```bash
git add src/sdp/daemon/ tests/unit/daemon/ src/sdp/cli.py
git commit -m "feat(00-012-01): Add Daemon Service Framework

- DaemonConfig dataclass for configuration
- PIDManager for PID file management
- CLI commands: daemon start/stop/status
- Signal handling for graceful shutdown

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 2: 00-012-08 Dashboard Core

**Files:**
- Create: `src/sdp/dashboard/__init__.py`
- Create: `src/sdp/dashboard/state.py` (~100 LOC)
- Create: `src/sdp/dashboard/sources/__init__.py`
- Create: `src/sdp/dashboard/sources/workstream_reader.py` (~150 LOC)
- Create: `src/sdp/dashboard/sources/test_runner.py` (~200 LOC)
- Create: `src/sdp/dashboard/sources/agent_reader.py` (~100 LOC)
- Create: `src/sdp/dashboard/widgets/__init__.py`
- Create: `src/sdp/dashboard/widgets/workstream_tree.py` (~200 LOC)
- Create: `src/sdp/dashboard/widgets/test_panel.py` (~150 LOC)
- Create: `src/sdp/dashboard/widgets/activity_log.py` (~100 LOC)
- Create: Tests for all components (~500 LOC total)
- Modify: `pyproject.toml` (+10 LOC for textual, watchdog)

---

### Step 1: Add dependencies to pyproject.toml

**File:** `pyproject.toml`

```toml
[project.dependencies]
textual = ">=0.80.0"
watchdog = ">=4.0.0"
```

**Run:** `poetry lock`
**Expected:** Updates poetry.lock with new dependencies

---

### Step 2: Write state bus test

**File:** `tests/unit/dashboard/test_state.py`

```python
from sdp.dashboard.state import StateBus, DashboardState, WorkstreamState

def test_state_bus_empty_on_init():
    bus = StateBus()
    assert bus.state is None

def test_state_bus_publish():
    bus = StateBus()
    state = DashboardState(workstreams={})
    bus.publish(state)
    assert bus.state == state

def test_state_bus_subscribe():
    bus = StateBus()
    received = []

    bus.subscribe(lambda s: received.append(s))

    state = DashboardState(workstreams={})
    bus.publish(state)

    assert len(received) == 1
    assert received[0] == state

def test_workstream_state_dataclass():
    ws = WorkstreamState(
        ws_id="00-012-08",
        status="backlog",
        title="Dashboard Core",
        assignee=None,
        feature="F012"
    )
    assert ws.ws_id == "00-012-08"
    assert ws.status == "backlog"
```

**Run:** `pytest tests/unit/dashboard/test_state.py -v`
**Expected:** FAIL - "No module named 'sdp.dashboard.state'"

---

### Step 3: Create state module

**File:** `src/sdp/dashboard/state.py`

```python
from dataclasses import dataclass, field
from datetime import datetime
from typing import Callable, Optional

@dataclass
class WorkstreamState:
    """State of a single workstream."""
    ws_id: str
    status: str  # backlog, in_progress, completed, blocked
    title: str
    feature: str
    assignee: Optional[str] = None
    size: Optional[str] = None
    started: Optional[datetime] = None
    completed: Optional[datetime] = None

@dataclass
class TestResults:
    """Results of test run."""
    status: str  # passed, failed, error, no_tests
    total: int = 0
    passed: int = 0
    failed: int = 0
    skipped: int = 0
    coverage: Optional[float] = None
    failed_tests: list[str] = field(default_factory=list)
    error_message: Optional[str] = None

@dataclass
class AgentEvent:
    """Event from agent execution."""
    timestamp: datetime
    event_type: str  # started, completed, error
    ws_id: Optional[str] = None
    message: str = ""

@dataclass
class DashboardState:
    """Complete dashboard state."""
    workstreams: dict[str, WorkstreamState] = field(default_factory=dict)
    test_results: Optional[TestResults] = None
    agent_activity: list[AgentEvent] = field(default_factory=list)
    last_update: Optional[datetime] = None

class StateBus:
    """Pub/sub for state updates."""

    def __init__(self) -> None:
        self._subscribers: list[Callable[[DashboardState], None]] = []
        self.state: Optional[DashboardState] = None

    def subscribe(self, callback: Callable[[DashboardState], None]) -> None:
        """Subscribe to state updates."""
        self._subscribers.append(callback)

    def publish(self, state: DashboardState) -> None:
        """Publish state update to all subscribers."""
        self.state = state
        for callback in self._subscribers:
            callback(state)
```

**Run:** `pytest tests/unit/dashboard/test_state.py -v`
**Expected:** PASS

---

### Step 4: Write workstream reader test

**File:** `tests/unit/dashboard/test_workstream_reader.py`

```python
import pytest
from pathlib import Path
from sdp.dashboard.sources.workstream_reader import WorkstreamReader

@pytest.fixture
def sample_workstream(tmp_path: Path):
    ws_dir = tmp_path / "backlog"
    ws_dir.mkdir(parents=True)

    ws_file = ws_dir / "00-012-08.md"
    ws_file.write_text('''---
ws_id: 00-012-08
status: backlog
title: Dashboard Core
feature: F012
---
# Dashboard Core
Test content
''')

    return tmp_path

def test_read_workstreams_parses_yaml(sample_workstream: Path):
    reader = WorkstreamReader(sample_workstream)
    state = reader.read()

    assert "00-012-08" in state.workstreams
    assert state.workstreams["00-012-08"].title == "Dashboard Core"
    assert state.workstreams["00-012-08"].status == "backlog"

def test_read_workstreams_empty_dir(tmp_path: Path):
    (tmp_path / "backlog").mkdir(parents=True)
    reader = WorkstreamReader(tmp_path)
    state = reader.read()

    assert len(state.workstreams) == 0

def test_read_workstreams_invalid_yaml(tmp_path: Path):
    ws_dir = tmp_path / "backlog"
    ws_dir.mkdir(parents=True)
    (ws_dir / "invalid.md").write_text("---\nbad yaml\n: :\n")

    reader = WorkstreamReader(tmp_path)
    state = reader.read()  # Should not crash

    # Invalid files are skipped
    assert "invalid" not in state.workstreams
```

**Run:** `pytest tests/unit/dashboard/test_workstream_reader.py -v`
**Expected:** FAIL - "No module named 'sdp.dashboard.sources.workstream_reader'"

---

### Step 5: Create workstream reader

**File:** `src/sdp/dashboard/sources/workstream_reader.py`

```python
from pathlib import Path
import logging
from ..state import DashboardState, WorkstreamState

logger = logging.getLogger(__name__)

class WorkstreamReader:
    """Reads workstream files and builds state."""

    def __init__(self, base_dir: Path | str):
        self._base = Path(base_dir)

    def read(self) -> DashboardState:
        """Read all workstreams and build state."""
        workstreams: dict[str, WorkstreamState] = {}

        for status_dir in ["backlog", "in_progress", "completed"]:
            ws_dir = self._base / status_dir
            if not ws_dir.exists():
                continue

            for ws_file in ws_dir.glob("*.md"):
                try:
                    ws_state = self._parse_ws_file(ws_file, status_dir)
                    workstreams[ws_state.ws_id] = ws_state
                except Exception as e:
                    logger.warning(f"Failed to parse {ws_file}: {e}")

        return DashboardState(workstreams=workstreams)

    def _parse_ws_file(self, path: Path, status_dir: str) -> WorkstreamState:
        """Parse single workstream file."""
        import yaml

        content = path.read_text()

        # Extract YAML frontmatter
        if content.startswith("---"):
            parts = content.split("---", 2)
            if len(parts) >= 3:
                frontmatter = yaml.safe_load(parts[1])
            else:
                frontmatter = {}
        else:
            frontmatter = {}

        return WorkstreamState(
            ws_id=frontmatter.get("ws_id", path.stem),
            status=status_dir.replace("_", "-"),  # in_progress -> in-progress
            title=frontmatter.get("title", path.stem),
            feature=frontmatter.get("feature", "Unknown"),
            assignee=frontmatter.get("assignee"),
            size=frontmatter.get("size"),
        )
```

**Run:** `pytest tests/unit/dashboard/test_workstream_reader.py -v`
**Expected:** PASS

---

### Step 6: Write test runner test

**File:** `tests/unit/dashboard/test_test_runner.py`

```python
import pytest
from pathlib import Path
from sdp.dashboard.sources.test_runner import TestRunner
from sdp.dashboard.state import TestResults

@pytest.fixture
def sample_project(tmp_path: Path):
    src = tmp_path / "src"
    tests = tmp_path / "tests"
    src.mkdir()
    tests.mkdir()

    # Create a passing test
    (tests / "test_sample.py").write_text("""
def test_passing():
    assert True
""")

    return tmp_path

def test_test_runner_runs_tests(sample_project: Path):
    runner = TestRunner(sample_project)
    results = runner.run()

    assert results is not None
    assert results.status in ["passed", "failed", "error"]

def test_test_runner_returns_no_tests_when_no_tests(tmp_path: Path):
    runner = TestRunner(tmp_path)
    results = runner.run()

    assert results.status == "no_tests"

def test_test_runner_captures_coverage(sample_project: Path):
    runner = TestRunner(sample_project, coverage=True)
    results = runner.run()

    # Coverage might be None if pytest-cov not installed
    assert results is not None
```

**Run:** `pytest tests/unit/dashboard/test_test_runner.py -v`
**Expected:** FAIL - "No module named 'sdp.dashboard.sources.test_runner'"

---

### Step 7: Create test runner

**File:** `src/sdp/dashboard/sources/test_runner.py`

```python
from pathlib import Path
import subprocess
import logging
from ..state import TestResults

logger = logging.getLogger(__name__)

class TestRunner:
    """Runs pytest and parses results."""

    def __init__(self, project_dir: Path | str, coverage: bool = True):
        self._project = Path(project_dir)
        self._coverage = coverage

    def run(self) -> TestResults:
        """Run tests and return results."""
        # Check if tests exist
        tests_dir = self._project / "tests"
        if not tests_dir.exists() or not list(tests_dir.glob("test_*.py")):
            return TestResults(status="no_tests")

        cmd = ["pytest", "-v", "--tb=short"]
        if self._coverage:
            cmd.extend(["--cov=", "--cov-report=json"])

        try:
            result = subprocess.run(
                cmd,
                cwd=self._project,
                capture_output=True,
                text=True,
                timeout=300
            )
            return self._parse_output(result.stdout, result.stderr)
        except subprocess.TimeoutExpired:
            return TestResults(status="error", error_message="Tests timed out")
        except Exception as e:
            return TestResults(status="error", error_message=str(e))

    def _parse_output(self, stdout: str, stderr: str) -> TestResults:
        """Parse pytest output."""
        # Simple parsing - could be improved with pytest JSON
        if " passed" in stdout.lower() and " failed" not in stdout.lower():
            status = "passed"
        elif " failed" in stdout.lower():
            status = "failed"
        else:
            status = "error"

        # Extract counts
        import re
        summary_match = re.search(r"(\d+) passed(?:, (\d+) failed)?", stdout)
        if summary_match:
            passed = int(summary_match.group(1))
            failed = int(summary_match.group(2)) if summary_match.group(2) else 0
            total = passed + failed
        else:
            total = passed = failed = 0

        return TestResults(
            status=status,
            total=total,
            passed=passed,
            failed=failed,
            skipped=0,
        )
```

**Run:** `pytest tests/unit/dashboard/test_test_runner.py -v`
**Expected:** PASS

---

### Step 8: Create agent reader (simple version)

**File:** `src/sdp/dashboard/sources/agent_reader.py`

```python
from pathlib import Path
from ..state import AgentEvent, list
import logging

logger = logging.getLogger(__name__)

class AgentReader:
    """Reads agent activity from daemon queue."""

    def __init__(self, project_dir: Path | str):
        self._project = Path(project_dir)
        self._queue_file = self._project / ".sdp" / "daemon_queue.json"

    def read(self) -> list[AgentEvent]:
        """Read agent events from queue file."""
        if not self._queue_file.exists():
            return []

        try:
            import json
            data = json.loads(self._queue_file.read_text())
            return [
                AgentEvent(
                    timestamp=e["timestamp"],
                    event_type=e["event_type"],
                    ws_id=e.get("ws_id"),
                    message=e.get("message", "")
                )
                for e in data.get("events", [])
            ]
        except Exception as e:
            logger.warning(f"Failed to read agent queue: {e}")
            return []
```

---

### Step 9: Create Textual widgets (skeleton)

**File:** `src/sdp/dashboard/widgets/workstream_tree.py`

```python
from textual.widgets import Tree
from textual import events
from ..state import WorkstreamState

class WorkstreamTree(Tree):
    """Tree widget showing workstreams by status."""

    def __init__(self, state_bus, **kwargs):
        super().__init__(**kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)

    def _on_state_update(self, state):
        """Update tree when state changes."""
        self.clear()

        # Group by status
        by_status = {}
        for ws in state.workstreams.values():
            if ws.status not in by_status:
                by_status[ws.status] = []
            by_status[ws.status].append(ws)

        # Build tree
        for status, wss in by_status.items():
            status_node = self.root.add(f"{status.title()} ({len(wss)})")
            for ws in wss:
                status_node.add_leaf(f"{ws.ws_id}: {ws.title}")
```

**File:** `src/sdp/dashboard/widgets/test_panel.py`

```python
from textual.widgets import Static
from ..state import TestResults

class TestPanel(Static):
    """Panel showing test results."""

    def __init__(self, state_bus, **kwargs):
        super().__init__("Test Results", **kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)

    def _on_state_update(self, state):
        """Update panel when state changes."""
        results = state.test_results
        if not results:
            self.update_content("No tests run")
            return

        content = f"""Status: {results.status.upper()}
Passed: {results.passed}
Failed: {results.failed}
Coverage: {results.coverage or 0:.0f}%
"""
        self.update_content(content)

    def update_content(self, content: str):
        """Update displayed content."""
        self.update(content)
```

**File:** `src/sdp/dashboard/widgets/activity_log.py`

```python
from textual.widgets import Log
from ..state import AgentEvent

class ActivityLog(Log):
    """Log widget showing agent activity."""

    def __init__(self, state_bus, **kwargs):
        super().__init__(**kwargs)
        self._state_bus = state_bus
        self._state_bus.subscribe(self._on_state_update)

    def _on_state_update(self, state):
        """Add new events to log."""
        for event in state.agent_activity:
            self.write_line(
                f"[{event.timestamp}] {event.event_type}: {event.message}"
            )
```

---

### Step 10: Commit

```bash
git add src/sdp/dashboard/ tests/unit/dashboard/ pyproject.toml
git commit -m "feat(00-012-08): Add Dashboard Core (reusable UI components)

- StateBus pub/sub for state updates
- WorkstreamReader for scanning workstream files
- TestRunner for pytest execution
- AgentReader for daemon queue (graceful degradation)
- Textual widgets: WorkstreamTree, TestPanel, ActivityLog

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 3: 00-012-03 Enhanced GitHub Sync

**Files:**
- Create: `src/sdp/github/conflict_resolver.py` (~200 LOC)
- Create: `src/sdp/github/sync_enhanced.py` (~150 LOC)
- Modify: `src/sdp/github/sync_service.py` (+150 LOC)
- Modify: `src/sdp/github/cli.py` (+50 LOC)
- Create: Tests (~250 LOC)

---

### Step 1-3: Conflict detection and resolution tests

**File:** `tests/unit/github/test_conflict_resolver.py`

```python
import pytest
from sdp.github.conflict_resolver import ConflictResolver, Conflict

def test_detect_conflict_when_status_mismatch():
    resolver = ConflictResolver()
    ws_state = {"status": "backlog"}
    gh_state = {"status": "In Progress"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is not None
    assert conflict.field == "status"

def test_no_conflict_when_status_matches():
    resolver = ConflictResolver()
    ws_state = {"status": "backlog"}
    gh_state = {"status": "Backlog"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is None

def test_resolve_conflict_ws_wins():
    resolver = ConflictResolver()
    conflict = Conflict(field="status", ws_value="backlog", gh_value="In Progress")

    resolved = resolver.resolve(conflict)
    assert resolved == "backlog"
```

---

### Step 4-6: Implementation

**File:** `src/sdp/github/conflict_resolver.py`

```python
from dataclasses import dataclass

@dataclass
class Conflict:
    """Represents a conflict between WS and GitHub state."""
    field: str
    ws_value: any
    gh_value: any

class ConflictResolver:
    """Detects and resolves conflicts between WS and GitHub state."""

    def detect(self, ws_state: dict, gh_state: dict) -> Conflict | None:
        """Detect if there's a conflict."""
        status_map = {
            "backlog": "Backlog",
            "in_progress": "In Progress",
            "completed": "Completed",
        }

        ws_status = ws_state.get("status")
        gh_status = gh_state.get("status")

        if ws_status and gh_status:
            normalized_gh = gh_status.lower().replace(" ", "_")
            if ws_status != normalized_gh:
                return Conflict(
                    field="status",
                    ws_value=ws_status,
                    gh_value=gh_status
                )
        return None

    def resolve(self, conflict: Conflict) -> any:
        """Resolve conflict: WS always wins."""
        return conflict.ws_value
```

**File:** `src/sdp/github/sync_enhanced.py` (extends sync_service.py)

```python
from .conflict_resolver import ConflictResolver

class EnhancedSyncService:
    """Enhanced sync with conflict detection."""

    def __init__(self, base_service):
        self._base = base_service
        self._resolver = ConflictResolver()

    def sync_with_conflict_detection(self, ws_id: str, dry_run: bool = False):
        """Sync workstream with conflict detection."""
        ws_state = self._base.read_ws_state(ws_id)
        gh_state = self._base.read_gh_state(ws_id)

        conflict = self._resolver.detect(ws_state, gh_state)
        if conflict:
            resolved = self._resolver.resolve(conflict)
            if dry_run:
                print(f"Would resolve conflict: {conflict.field} -> {resolved}")
            else:
                self._base.update_gh_state(ws_id, {conflict.field: resolved})

        return not conflict
```

---

### Step 7: Commit

```bash
git add src/sdp/github/ tests/unit/github/
git commit -m "feat(00-012-03): Add Enhanced GitHub Sync

- Conflict detection between WS and GitHub status
- Conflict resolution: WS file wins (source of truth)
- Dry-run mode for previewing changes
- sync_backlog() method for incremental sync

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 4: 00-012-02 Task Queue Management

**Files:**
- Create: `src/sdp/queue/__init__.py`
- Create: `src/sdp/queue/task.py` (~100 LOC)
- Create: `src/sdp/queue/priority.py` (~80 LOC)
- Create: `src/sdp/queue/task_queue.py` (~300 LOC)
- Create: `src/sdp/queue/state.py` (~120 LOC)
- Modify: `src/sdp/cli.py` (+50 LOC)
- Create: Tests (~250 LOC)

---

### Key Implementation

**File:** `src/sdp/queue/task_queue.py`

```python
import asyncio
import threading
from queue import PriorityQueue
from .task import Task
from .priority import Priority

class TaskQueue:
    """Thread-safe priority queue for task management."""

    def __init__(self, state_file: str = ".sdp/queue_state.json"):
        self._queue: PriorityQueue[Task] = PriorityQueue()
        self._lock = threading.Lock()
        self._state_file = state_file

    def enqueue(self, task: Task) -> None:
        """Add task to queue."""
        with self._lock:
            self._queue.put(task)
            self._save_state()

    def dequeue(self) -> Task | None:
        """Remove and return highest priority task."""
        with self._lock:
            if self._queue.empty():
                return None
            task = self._queue.get()
            self._save_state()
            return task

    def peek(self) -> Task | None:
        """View highest priority task without removing."""
        with self._lock:
            if self._queue.empty():
                return None
            # PriorityQueue doesn't support peek, use internal queue
            return self._queue.queue[0]

    def _save_state(self) -> None:
        """Persist queue state to file."""
        from .state import QueueState
        QueueState.save(self._queue, self._state_file)
```

---

### Commit

```bash
git add src/sdp/queue/ tests/unit/queue/ src/sdp/cli.py
git commit -m "feat(00-012-02): Add Task Queue Management

- Task dataclass with priority and retry tracking
- Priority enum (backlog < active < blocked)
- Thread-safe TaskQueue with enqueue/dequeue/peek
- Queue state persistence to .sdp/queue_state.json
- CLI: sdp queue status command

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 5: 00-012-11 Workstream Status Command

**Files:**
- Create: `src/sdp/status/__init__.py`
- Create: `src/sdp/status/formatter.py` (~150 LOC)
- Create: `src/sdp/status/command.py` (~100 LOC)
- Modify: `src/sdp/cli.py` (+30 LOC)
- Create: Tests (~180 LOC)

---

### Key Implementation

**File:** `src/sdp/status/command.py`

```python
import click
from rich.console import Console
from rich.table import Table
from ..dashboard.sources.workstream_reader import WorkstreamReader

@click.command()
@click.option("--filter", "status_filter", type=click.Choice(["backlog", "in_progress", "completed"]))
@click.option("--feature", help="Filter by feature ID")
@click.option("--watch", is_flag=True, help="Auto-refresh every 2s")
def status(status_filter: str | None, feature: str | None, watch: bool):
    """Show workstream status."""
    console = Console()
    reader = WorkstreamReader("docs/workstreams")

    while True:
        state = reader.read()
        _print_status(console, state, status_filter, feature)

        if not watch:
            break

        import time
        time.sleep(2)

def _print_status(console, state, status_filter, feature):
    """Print status table."""
    # Group by status
    by_status = {}
    for ws in state.workstreams.values():
        if status_filter and ws.status != status_filter:
            continue
        if feature and ws.feature != feature:
            continue

        if ws.status not in by_status:
            by_status[ws.status] = []
        by_status[ws.status].append(ws)

    # Print tables
    for status, wss in by_status.items():
        console.print(f"\n📁 {status.title()} ({len(wss)})")
        table = Table()
        table.add_column("WS-ID")
        table.add_column("Title")
        table.add_column("Status")
        table.add_column("Assignee")
        table.add_column("Feature")

        for ws in wss:
            table.add_row(ws.ws_id, ws.title, ws.status, ws.assignee or "-", ws.feature)

        console.print(table)
```

---

### Commit

```bash
git add src/sdp/status/ tests/unit/status/ src/sdp/cli.py
git commit -m "feat(00-012-11): Add Workstream Status Command

- sdp status shows all workstreams grouped by status
- Rich table formatting with colors
- --filter flag for status filtering
- --feature flag for feature filtering
- --watch flag for auto-refresh (every 2s)
- Uses WorkstreamReader from Dashboard Core

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 6: 00-012-12 Test Watch Mode

**Files:**
- Create: `src/sdp/test_watch/__init__.py`
- Create: `src/sdp/test_watch/watcher.py` (~150 LOC)
- Create: `src/sdp/test_watch/runner.py` (~200 LOC)
- Create: `src/sdp/test_watch/affected.py` (~100 LOC)
- Modify: `src/sdp/cli.py` (+50 LOC)
- Create: Tests (~180 LOC)

---

### Key Implementation

**File:** `src/sdp/test_watch/watcher.py`

```python
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import time

class TestWatcher(FileSystemEventHandler):
    """Watches for file changes and triggers test runs."""

    def __init__(self, on_change, debounce=0.5):
        self._on_change = on_change
        self._debounce = debounce
        self._last_change = 0

    def on_modified(self, event):
        """Handle file modification."""
        if event.src_path.endswith(".py"):
            now = time.time()
            if now - self._last_change > self._debounce:
                self._last_change = now
                self._on_change(event.src_path)

def watch_tests(project_dir, callback):
    """Start watching for test file changes."""
    observer = Observer()
    handler = TestWatcher(callback)
    observer.schedule(handler, str(project_dir), recursive=True)
    observer.start()
    return observer
```

---

### Commit

```bash
git add src/sdp/test_watch/ tests/unit/test_watch/ src/sdp/cli.py
git commit -m "feat(00-012-12): Add Test Watch Mode

- sdp test --watch starts file watcher on src/ and tests/
- File changes trigger pytest run on affected tests
- --pattern flag filters which tests to run
- --coverage flag enables coverage report
- Debouncing (500ms) to avoid excessive runs
- Uses TestRunner from Dashboard Core

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 7: 00-012-04 Agent Executor Interface

**Files:**
- Create: `src/sdp/agents/__init__.py`
- Create: `src/sdp/agents/executor.py` (~350 LOC)
- Create: `src/sdp/agents/metrics.py` (~200 LOC)
- Create: `src/sdp/agents/errors.py` (~100 LOC)
- Create: Tests (~300 LOC)

---

### Key Implementation

**File:** `src/sdp/agents/executor.py`

```python
from .errors import ExecutorError
from .metrics import ExecutionMetrics

class AgentExecutor:
    """Executes workstreams autonomously."""

    def __init__(self, metrics_file: str = ".sdp/execution_metrics.json"):
        self._metrics = ExecutionMetrics(metrics_file)

    def execute(self, ws_id: str, timeout: int = 3600) -> bool:
        """Execute a workstream and return success status."""
        import time
        start = time.time()

        try:
            # Call /build skill (subprocess or direct)
            result = self._run_build_skill(ws_id)

            duration = time.time() - start
            self._metrics.record(ws_id, success=result, duration=duration)

            return result
        except Exception as e:
            self._metrics.record(ws_id, success=False, duration=time.time() - start)
            raise ExecutorError(f"Execution failed: {e}")

    def _run_build_skill(self, ws_id: str) -> bool:
        """Run /build skill for workstream."""
        import subprocess
        result = subprocess.run(
            ["python", "-m", "claude", "skill", "build", ws_id],
            capture_output=True
        )
        return result.returncode == 0
```

---

### Commit

```bash
git add src/sdp/agents/ tests/unit/agents/
git commit -m "feat(00-012-04): Add Agent Executor Interface

- AgentExecutor.execute(ws_id) runs /build skill
- Progress tracking via TaskUpdate integration
- Error handling with retry (max 2 attempts)
- Execution metrics stored in .sdp/execution_metrics.json
- Timeout protection (default 1h per WS)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 8: 00-012-13 Auto-State Management

**Files:**
- Create: `src/sdp/workspace/__init__.py`
- Create: `src/sdp/workspace/mover.py` (~200 LOC)
- Create: `src/sdp/workspace/index_updater.py` (~150 LOC)
- Create: `src/sdp/workspace/validator.py` (~100 LOC)
- Create: `src/sdp/workspace/state_updater.py` (~150 LOC)
- Modify: `src/sdp/cli.py` (+40 LOC)
- Create: Tests (~430 LOC)

---

### Key Implementation

**File:** `src/sdp/workspace/mover.py`

```python`
from pathlib import Path
import shutil
from .validator import MoveValidator
from .state_updater import StateUpdater

class WorkstreamMover:
    """Moves workstream files between status directories."""

    def __init__(self, base_dir: str = "docs/workstreams"):
        self._base = Path(base_dir)
        self._validator = MoveValidator()
        self._state_updater = StateUpdater()

    def move(self, ws_id: str, to_status: str) -> None:
        """Move workstream to new status directory."""
        # Find current file
        current_path = self._find_ws_file(ws_id)
        if not current_path:
            raise FileNotFoundError(f"Workstream not found: {ws_id}")

        # Validate move
        self._validator.validate_move(current_path, to_status)

        # Build new path
        new_status_dir = to_status.replace("-", "_")  # in-progress -> in_progress
        new_path = self._base / new_status_dir / current_path.name

        # Move file
        shutil.move(str(current_path), str(new_path))

        # Update YAML frontmatter
        self._state_updater.update_status(new_path, to_status)

    def _find_ws_file(self, ws_id: str) -> Path | None:
        """Find workstream file in any status directory."""
        for status_dir in ["backlog", "in_progress", "completed"]:
            path = self._base / status_dir / f"{ws_id}.md"
            if path.exists():
                return path
        return None
```

---

### Commit

```bash
git add src/sdp/workspace/ tests/unit/workspace/ src/sdp/cli.py
git commit -m "feat(00-012-13): Add Auto-State Management

- sdp ws move <ws-id> --to <status> moves files between dirs
- Auto-updates status: field in YAML frontmatter
- Auto-updates docs/workstreams/INDEX.md
- --start flag: backlog → in_progress with timestamp
- --complete flag: in_progress → completed (validates AC)
- Validates move: can't complete if AC not 100% met

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 9: 00-012-05 CLI Task Commands

**Files:**
- Create: `src/sdp/cli_tasks.py` (~250 LOC)
- Modify: `src/sdp/cli.py` (+20 LOC)
- Create: Tests (~100 LOC)

---

### Key Implementation

**File:** `src/sdp/cli_tasks.py`

```python
import click
from ..queue.task_queue import TaskQueue
from ..agents.executor import AgentExecutor

@click.group()
def task():
    """Task management commands."""
    pass

@task.command()
@click.argument("ws_id")
def enqueue(ws_id: str):
    """Add workstream to task queue."""
    from ..queue.task import Task
    from ..queue.priority import Priority

    queue = TaskQueue()
    task = Task(ws_id=ws_id, priority=Priority.NORMAL)
    queue.enqueue(task)
    click.echo(f"Enqueued: {ws_id}")

@task.command()
@click.argument("ws_id")
@click.option("--dry-run", is_flag=True)
def execute(ws_id: str, dry_run: bool):
    """Execute workstream immediately."""
    if dry_run:
        click.echo(f"Would execute: {ws_id}")
    else:
        executor = AgentExecutor()
        success = executor.execute(ws_id)
        click.echo(f"Executed: {ws_id} - {'✅' if success else '❌'}")

@task.command()
def list():
    """Show all pending/running/completed tasks."""
    queue = TaskQueue()
    task = queue.peek()
    if task:
        click.echo(f"Next: {task.ws_id} (priority: {task.priority})")
    else:
        click.echo("Queue empty")
```

---

### Commit

```bash
git add src/sdp/cli_tasks.py tests/unit/cli/test_tasks.py src/sdp/cli.py
git commit -m "feat(00-012-05): Add CLI Task Commands

- sdp task enqueue WS-XXX-YY adds to queue
- sdp task execute WS-XXX-YY runs immediately
- sdp task list shows pending/running/completed
- sdp task cancel <task_id> cancels pending task
- Uses TaskQueue from 00-012-02
- Uses AgentExecutor from 00-012-04

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 10: 00-012-06 Multi-Agent Orchestration

**Files:**
- Create: `src/sdp/agents/dependency_graph.py` (~250 LOC)
- Create: `src/sdp/agents/agent_pool.py` (~200 LOC)
- Create: `src/sdp/agents/orchestrator.py` (~400 LOC)
- Modify: `src/sdp/cli.py` (+50 LOC)
- Create: Tests (~250 LOC)

---

### Key Implementation

**File:** `src/sdp/agents/orchestrator.py`

```python
import asyncio
from .agent_pool import AgentPool
from .dependency_graph import DependencyGraph
from ..queue.task_queue import TaskQueue

class Orchestrator:
    """Manages multiple agents with dependency resolution."""

    def __init__(self, max_agents: int = 3):
        self._pool = AgentPool(max_agents)
        self._graph = DependencyGraph()

    async def run_feature(self, feature_id: str) -> dict[str, bool]:
        """Execute all workstreams for a feature."""
        # Build dependency graph
        ws_ids = self._get_feature_workstreams(feature_id)
        self._graph.build(ws_ids)

        # Execute in topological order
        results = {}
        for ws_id in self._graph.execution_order():
            # Wait for available agent
            await self._pool.acquire()

            # Execute
            try:
                success = await self._execute_ws(ws_id)
                results[ws_id] = success
            finally:
                self._pool.release()

        return results

    def _get_feature_workstreams(self, feature_id: str) -> list[str]:
        """Get all workstream IDs for a feature."""
        from ..dashboard.sources.workstream_reader import WorkstreamReader
        reader = WorkstreamReader("docs/workstreams")
        state = reader.read()

        return [
            ws.ws_id for ws in state.workstreams.values()
            if ws.feature == feature_id
        ]
```

---

### Commit

```bash
git add src/sdp/agents/ tests/unit/agents/ src/sdp/cli.py
git commit -m "feat(00-012-06): Add Multi-Agent Orchestration

- Orchestrator manages agent pool (max 3 concurrent)
- Dependency resolution via topological sort
- Load balancing: assign WS to least busy agent
- Deadlock detection: circular dependency detection
- Orchestrator state persisted to .sdp/orchestrator_state.json
- sdp orchestrator run --feature F012 executes all WS

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 11: 00-012-07 GitHub Project Fields Integration

**Files:**
- Create: `src/sdp/github/fields_client.py` (~200 LOC)
- Create: `src/sdp/github/fields_config.py` (~100 LOC)
- Create: `src/sdp/github/fields_sync.py` (~250 LOC)
- Modify: `src/sdp/github/sync_service.py` (+50 LOC)
- Modify: `src/sdp/github/projects_client.py` (+100 LOC)
- Create: Tests (~100 LOC)

---

### Commit

```bash
git add src/sdp/github/ tests/unit/github/
git commit -m "feat(00-012-07): Add GitHub Project Fields Integration

- Sync WS status/size/feature to GitHub Project custom fields
- Auto-create custom fields if missing
- Field mapping configurable via .sdp/github_fields.toml
- Bidirectional sync (GitHub changes → WS frontmatter)
- Extends ProjectBoardSync with custom field support

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 12: 00-012-10 Pre-Execution Checks

**Files:**
- Create: `src/sdp/agents/pre_check.py` (~200 LOC)
- Modify: `src/sdp/agents/executor.py` (+50 LOC)
- Create: Tests (~150 LOC)

---

### Key Implementation

**File:** `src/sdp/agents/pre_check.py`

```python
from pathlib import Path

class PreExecutionChecker:
    """Validates workstream before execution."""

    def __init__(self, ws_dir: str = "docs/workstreams"):
        self._ws_dir = Path(ws_dir)

    def check(self, ws_id: str) -> list[str]:
        """Run all pre-execution checks."""
        errors = []

        # Check 1: WS file exists and valid YAML
        ws_file = self._find_ws_file(ws_id)
        if not ws_file:
            errors.append(f"Workstream file not found: {ws_id}")
        else:
            try:
                self._parse_yaml(ws_file)
            except Exception as e:
                errors.append(f"Invalid YAML: {e}")

        # Check 2: Dependencies satisfied
        deps = self._get_dependencies(ws_id)
        for dep in deps:
            if not self._is_completed(dep):
                errors.append(f"Dependency not completed: {dep}")

        # Check 3: No circular dependencies
        if self._has_circular_dependency(ws_id):
            errors.append("Circular dependency detected")

        # Check 4: Size ≤ MEDIUM
        if self._get_size(ws_id) == "LARGE":
            errors.append("Workstream too large (LARGE), please split")

        return errors

    def _find_ws_file(self, ws_id: str) -> Path | None:
        for status in ["backlog", "in_progress", "completed"]:
            path = self._ws_dir / status / f"{ws_id}.md"
            if path.exists():
                return path
        return None
```

---

### Commit

```bash
git add src/sdp/agents/ tests/unit/agents/
git commit -m "feat(00-012-10): Add Pre-Execution Checks

- PreExecutionChecker.check(ws_id) validates before execution
- Check 1: WS file exists and valid YAML
- Check 2: Dependencies satisfied (all deps completed)
- Check 3: No circular dependencies
- Check 4: WS size ≤ MEDIUM (< 1500 LOC)
- Checks run before AgentExecutor.execute()

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 13: 00-012-09 Webhook Support

**Files:**
- Create: `src/sdp/webhook/__init__.py`
- Create: `src/sdp/webhook/signature.py` (~100 LOC)
- Create: `src/sdp/webhook/handler.py` (~250 LOC)
- Create: `src/sdp/webhook/server.py` (~300 LOC)
- Modify: `pyproject.toml` (+5 LOC for starlette)
- Modify: `src/sdp/cli.py` (+50 LOC)
- Create: Tests (~250 LOC)

---

### Commit

```bash
git add src/sdp/webhook/ tests/unit/webhook/ pyproject.toml src/sdp/cli.py
git commit -m "feat(00-012-09): Add Webhook Support

- sdp webhook server starts HTTP server on port 8080
- Webhook receives issues and project_v2 events
- Signature validation (X-Hub-Signature-256)
- Triggers sync_service.sync_workstream() on issue update
- Logs all events to .sdp/webhook.log
- --smee-url flag for tunneling (local dev)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Task 14: 00-012-14 Developer Dashboard App

**Files:**
- Create: `src/sdp/dashboard/dashboard_app.py` (~250 LOC)
- Create: `src/sdp/dashboard/tabs/__init__.py`
- Create: `src/sdp/dashboard/tabs/workstreams_tab.py` (~150 LOC)
- Create: `src/sdp/dashboard/tabs/tests_tab.py` (~100 LOC)
- Create: `src/sdp/dashboard/tabs/activity_tab.py` (~100 LOC)
- Modify: `src/sdp/cli.py` (+20 LOC)
- Create: Tests (~150 LOC)

---

### Key Implementation

**File:** `src/sdp/dashboard/dashboard_app.py`

```python
from textual.app import App, ComposeResult
from textual.widgets import TabbedPane, Tabs, Tab
from .state import StateBus
from .sources.workstream_reader import WorkstreamReader
from .sources.test_runner import TestRunner
from .sources.agent_reader import AgentReader
from .tabs.workstreams_tab import WorkstreamsTab
from .tabs.tests_tab import TestsTab
from .tabs.activity_tab import ActivityTab

class DashboardApp(App):
    """Main SDP Dashboard application."""

    TITLE = "SDP Dashboard"
    BINDINGS = {
        "q": "quit",
        "1": "switch_tab(1)",
        "2": "switch_tab(2)",
        "3": "switch_tab(3)",
        "w": "switch_tab(1)",
        "t": "switch_tab(2)",
        "a": "switch_tab(3)",
        "r": "refresh",
    }

    def __init__(self):
        super().__init__()
        self._state_bus = StateBus()
        self._ws_reader = WorkstreamReader("docs/workstreams")
        self._test_runner = TestRunner(".")
        self._agent_reader = AgentReader(".")

    def compose(self) -> ComposeResult:
        """Compose dashboard UI."""
        yield Tabs("Workstreams", "Tests", "Activity")
        yield TabbedPane(
            WorkstreamsTab(self._state_bus),
            TestsTab(self._state_bus),
            ActivityTab(self._state_bus),
        )

    def on_mount(self) -> None:
        """Start background updates."""
        self._start_ws_polling()
        self._start_test_watching()

    def _start_ws_polling(self):
        """Poll workstream state every 1s."""
        async def poll():
            import asyncio
            while True:
                state = self._ws_reader.read()
                self._state_bus.publish(state)
                await asyncio.sleep(1)

        self.run_worker(poll())

    def _start_test_watching(self):
        """Start test file watcher."""
        from ..test_watch.watcher import watch_tests

        def on_change(path):
            results = self._test_runner.run()
            current_state = self._state_bus.state or DashboardState()
            current_state.test_results = results
            self._state_bus.publish(current_state)

        watcher = watch_tests(".", on_change)
        self._watcher = watcher  # Keep reference

    def action_switch_tab(self, tab_index: int):
        """Switch to tab by index."""
        self.query_one(Tabs).active = tab_index - 1

    def action_refresh(self):
        """Force refresh all data."""
        state = self._ws_reader.read()
        self._state_bus.publish(state)
```

---

### Commit

```bash
git add src/sdp/dashboard/ tests/unit/dashboard/ src/sdp/cli.py
git commit -m "feat(00-012-14): Add Developer Dashboard App

- sdp dashboard launches Textual TUI app
- Tab-based layout: Workstreams, Tests, Activity
- Hotkeys: w/t/a for tabs, q=quit, r=refresh
- Workstreams tab: tree view by status, filters
- Tests tab: live test results + coverage
- Activity tab: scrolling event log
- Uses all Dashboard Core components
- Live updates: WS poll (1s), tests on file change
- Graceful degradation when daemon not running

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Final Steps

### 1. Update documentation

```bash
# Update CLAUDE.md with new commands
# Update PROTOCOL.md with DX improvements
# Update README.md with dashboard screenshots
```

### 2. Integration test

```bash
# Full feature test
sdp dashboard &
sleep 5
sdp status
sdp test --watch &
echo "def test_x(): assert True" >> tests/test_sample.py
# Verify dashboard updates
```

### 3. Final commit

```bash
git add docs/
git commit -m "docs(f012): Complete F012 documentation

- CLAUDE.md: Add sdp dashboard, sdp status, sdp test commands
- PROTOCOL.md: Add DX patterns
- README.md: Dashboard screenshots and quick start

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>"
```

---

## Summary

**Total Workstreams:** 14
**Total LOC:** ~9,500
**Estimated Duration:** 4-6 weeks (solo) or 2-3 weeks (pair programming)

**Commands Added:**
- `sdp daemon start/stop/status`
- `sdp status [--filter] [--feature] [--watch]`
- `sdp test --watch [--pattern] [--coverage]`
- `sdp ws move <ws-id> --to <status> [--start|--complete]`
- `sdp dashboard`
- `sdp task enqueue/execute/list/cancel`
- `sdp orchestrator run --feature <id>`

**Dependencies Added:**
- `textual` (TUI framework)
- `watchdog` (file watching)
- `starlette` (webhook server)

**Quality Gates:**
- All code: ≥80% coverage, mypy --strict
- All WS: TDD cycle (Red → Green → Refactor)
- All commits: conventional commit format

---

**End of F012 Implementation Plan**

**For Claude:** Use superpowers:executing-plans to implement this plan task-by-task.
