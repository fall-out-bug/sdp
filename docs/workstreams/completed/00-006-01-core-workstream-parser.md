---
ws_id: 00-190-01
project_id: 00
feature: F006
status: completed
size: MEDIUM
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-190-01: Core Workstream Parser

### 🎯 Goal

**What must WORK after this WS is complete:**
- Standalone `sdp-protocol` package parses WS markdown files
- Parser extracts frontmatter, goal, AC, steps, code blocks
- Works independently of hw_checker
- 100% test coverage for parser

**Acceptance Criteria:**
- [ ] AC1: `sdp/core/workstream.py` parses WS files to dataclass
- [ ] AC2: Frontmatter extraction (ws_id, feature, status, size, etc.)
- [ ] AC3: Content extraction (goal, context, steps, code blocks)
- [ ] AC4: Validation (required fields, format checks)
- [ ] AC5: Unit tests with 100% coverage

**WS is NOT complete until Goal achieved (all AC checked).**

---

### Context

SDP Universal needs project-agnostic core. First step: extract WS parsing from current tightly-coupled implementation to standalone module.

**Current state:** WS parsing scattered across multiple files
**Target state:** Single `workstream.py` module with clean API

---

### Dependencies

None (first WS in F190)

---

### Input Files

**Existing (reference):**
- `sdp/src/sdp/github/sync_service.py` - has some WS parsing
- `tools/hw_checker/docs/workstreams/TEMPLATE.md` - WS format spec

**Created:**
- `sdp/src/sdp/core/__init__.py`
- `sdp/src/sdp/core/workstream.py`
- `sdp/tests/unit/core/test_workstream.py`

---

### Steps

#### 1. Create Core Module Structure

```bash
mkdir -p sdp/src/sdp/core
touch sdp/src/sdp/core/__init__.py
```

#### 2. Implement Workstream Dataclass

```python
# sdp/src/sdp/core/workstream.py
"""Workstream parsing and validation."""

from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Optional
import re
import yaml


class WorkstreamStatus(Enum):
    """Workstream lifecycle status."""
    BACKLOG = "backlog"
    ACTIVE = "active"
    COMPLETED = "completed"
    BLOCKED = "blocked"


class WorkstreamSize(Enum):
    """Workstream scope size."""
    SMALL = "SMALL"
    MEDIUM = "MEDIUM"
    LARGE = "LARGE"


@dataclass
class AcceptanceCriterion:
    """Single acceptance criterion."""
    id: str  # AC1, AC2, etc.
    description: str
    checked: bool = False


@dataclass
class Workstream:
    """Parsed workstream specification."""

    # Frontmatter (required)
    ws_id: str
    feature: str
    status: WorkstreamStatus
    size: WorkstreamSize

    # Frontmatter (optional)
    github_issue: Optional[int] = None
    assignee: Optional[str] = None
    started: Optional[str] = None
    completed: Optional[str] = None
    blocked_reason: Optional[str] = None

    # Content
    title: str = ""
    goal: str = ""
    acceptance_criteria: list[AcceptanceCriterion] = field(default_factory=list)
    context: str = ""
    dependencies: list[str] = field(default_factory=list)
    steps: list[str] = field(default_factory=list)
    code_blocks: list[str] = field(default_factory=list)

    # Source
    file_path: Optional[Path] = None


class WorkstreamParseError(Exception):
    """Error parsing workstream file."""
    pass


def parse_workstream(file_path: Path) -> Workstream:
    """Parse workstream markdown file.

    Args:
        file_path: Path to WS markdown file

    Returns:
        Parsed Workstream instance

    Raises:
        WorkstreamParseError: If parsing fails
    """
    content = file_path.read_text(encoding="utf-8")

    # Extract frontmatter
    frontmatter = _extract_frontmatter(content)

    # Extract content sections
    body = _strip_frontmatter(content)

    return Workstream(
        ws_id=frontmatter["ws_id"],
        feature=frontmatter["feature"],
        status=WorkstreamStatus(frontmatter["status"]),
        size=WorkstreamSize(frontmatter["size"]),
        github_issue=frontmatter.get("github_issue"),
        assignee=frontmatter.get("assignee"),
        started=frontmatter.get("started"),
        completed=frontmatter.get("completed"),
        blocked_reason=frontmatter.get("blocked_reason"),
        title=_extract_title(body),
        goal=_extract_section(body, "Goal"),
        acceptance_criteria=_extract_acceptance_criteria(body),
        context=_extract_section(body, "Context"),
        dependencies=_extract_dependencies(body),
        steps=_extract_steps(body),
        code_blocks=_extract_code_blocks(body),
        file_path=file_path,
    )


def _extract_frontmatter(content: str) -> dict:
    """Extract YAML frontmatter from markdown."""
    match = re.match(r"^---\n(.+?)\n---", content, re.DOTALL)
    if not match:
        raise WorkstreamParseError("No frontmatter found")

    try:
        data = yaml.safe_load(match.group(1))
    except yaml.YAMLError as e:
        raise WorkstreamParseError(f"Invalid YAML: {e}") from e

    # Validate required fields
    required = ["ws_id", "feature", "status", "size"]
    missing = [f for f in required if f not in data]
    if missing:
        raise WorkstreamParseError(f"Missing required fields: {missing}")

    return data


def _strip_frontmatter(content: str) -> str:
    """Remove frontmatter from markdown content."""
    return re.sub(r"^---\n.+?\n---\n*", "", content, flags=re.DOTALL)


def _extract_title(body: str) -> str:
    """Extract WS title from ## heading."""
    match = re.search(r"^## (.+)$", body, re.MULTILINE)
    return match.group(1) if match else ""


def _extract_section(body: str, section_name: str) -> str:
    """Extract content of a ### section."""
    pattern = rf"### .*{section_name}.*\n(.*?)(?=\n### |\n---|\Z)"
    match = re.search(pattern, body, re.DOTALL | re.IGNORECASE)
    return match.group(1).strip() if match else ""


def _extract_acceptance_criteria(body: str) -> list[AcceptanceCriterion]:
    """Extract acceptance criteria from Goal section."""
    goal_section = _extract_section(body, "Goal")
    criteria = []

    # Match: - [ ] AC1: description or - [x] AC1: description
    pattern = r"- \[([ x])\] (AC\d+): (.+)"
    for match in re.finditer(pattern, goal_section):
        checked = match.group(1) == "x"
        ac_id = match.group(2)
        description = match.group(3).strip()
        criteria.append(AcceptanceCriterion(ac_id, description, checked))

    return criteria


def _extract_dependencies(body: str) -> list[str]:
    """Extract WS dependencies."""
    dep_section = _extract_section(body, "Dependenc")
    if not dep_section or "none" in dep_section.lower():
        return []

    # Match WS-XXX-YY patterns
    return re.findall(r"WS-\d{3}-\d{2}", dep_section)


def _extract_steps(body: str) -> list[str]:
    """Extract numbered steps."""
    steps_section = _extract_section(body, "Steps")
    steps = []

    # Match: 1. Step description or #### 1. Step
    pattern = r"(?:^|\n)(?:####\s*)?\d+\.\s+(.+?)(?=\n(?:####\s*)?\d+\.|\n###|\Z)"
    for match in re.finditer(pattern, steps_section, re.DOTALL):
        steps.append(match.group(1).strip())

    return steps


def _extract_code_blocks(body: str) -> list[str]:
    """Extract fenced code blocks."""
    return re.findall(r"```[\w]*\n(.+?)```", body, re.DOTALL)
```

#### 3. Write Unit Tests

```python
# sdp/tests/unit/core/test_workstream.py
"""Tests for workstream parser."""

import pytest
from pathlib import Path
from sdp.core.workstream import (
    parse_workstream,
    Workstream,
    WorkstreamStatus,
    WorkstreamSize,
    WorkstreamParseError,
    AcceptanceCriterion,
)


@pytest.fixture
def sample_ws_content():
    return '''---
ws_id: 00-190-01
feature: F006
status: backlog
size: MEDIUM
github_issue: 123
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-190-01: Test Workstream

### 🎯 Goal

**What must WORK:**
- Feature X works

**Acceptance Criteria:**
- [ ] AC1: First criterion
- [x] AC2: Second criterion (done)
- [ ] AC3: Third criterion

---

### Context

This is the context.

### Dependencies

00--01, 00--02

### Steps

1. First step
2. Second step
3. Third step

### Code

```python
def example():
    pass
```
'''


@pytest.fixture
def ws_file(tmp_path, sample_ws_content):
    ws_path = tmp_path / "00--01-test.md"
    ws_path.write_text(sample_ws_content)
    return ws_path


class TestParseWorkstream:
    def test_parses_frontmatter(self, ws_file):
        ws = parse_workstream(ws_file)

        assert ws.ws_id == "00--01"
        assert ws.feature == "F190"
        assert ws.status == WorkstreamStatus.BACKLOG
        assert ws.size == WorkstreamSize.MEDIUM
        assert ws.github_issue == 123

    def test_parses_title(self, ws_file):
        ws = parse_workstream(ws_file)
        assert "Test Workstream" in ws.title

    def test_parses_acceptance_criteria(self, ws_file):
        ws = parse_workstream(ws_file)

        assert len(ws.acceptance_criteria) == 3
        assert ws.acceptance_criteria[0].id == "AC1"
        assert ws.acceptance_criteria[0].checked is False
        assert ws.acceptance_criteria[1].id == "AC2"
        assert ws.acceptance_criteria[1].checked is True

    def test_parses_dependencies(self, ws_file):
        ws = parse_workstream(ws_file)

        assert "00--01" in ws.dependencies
        assert "00--02" in ws.dependencies

    def test_parses_steps(self, ws_file):
        ws = parse_workstream(ws_file)

        assert len(ws.steps) >= 3

    def test_parses_code_blocks(self, ws_file):
        ws = parse_workstream(ws_file)

        assert len(ws.code_blocks) >= 1
        assert "def example" in ws.code_blocks[0]

    def test_missing_frontmatter_raises(self, tmp_path):
        bad_file = tmp_path / "bad.md"
        bad_file.write_text("# No frontmatter")

        with pytest.raises(WorkstreamParseError, match="No frontmatter"):
            parse_workstream(bad_file)

    def test_missing_required_field_raises(self, tmp_path):
        bad_file = tmp_path / "bad.md"
        bad_file.write_text('''---
ws_id: 00-190-01
---
# Missing fields
''')

        with pytest.raises(WorkstreamParseError, match="Missing required"):
            parse_workstream(bad_file)
```

---

### Expected Result

**Created:**
```
sdp/src/sdp/core/
├── __init__.py
└── workstream.py    (~200 LOC)

sdp/tests/unit/core/
└── test_workstream.py    (~150 LOC)
```

---

### Scope Estimate

- **Files:** 3 created
- **Lines:** ~350 (code: ~200, tests: ~150)
- **Tokens:** ~3000
- **Size:** MEDIUM

---

### Completion Criteria

```bash
# Tests pass
pytest sdp/tests/unit/core/test_workstream.py -v

# Coverage
pytest --cov=sdp.core.workstream --cov-fail-under=90

# Linters
ruff check sdp/src/sdp/core/
mypy sdp/src/sdp/core/ --strict

# No TODO/FIXME
grep -rn "TODO\|FIXME" sdp/src/sdp/core/
```

---

### Constraints

- NOT importing hw_checker modules
- NOT depending on GitHub integration
- NOT writing to files (read-only parser)
- USING dataclasses (not Pydantic yet)
