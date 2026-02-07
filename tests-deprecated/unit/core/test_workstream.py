"""Tests for workstream parsing and validation."""

import pytest
from pathlib import Path

from sdp.core.workstream.parser import WorkstreamParseError, parse_workstream
from sdp.domain.workstream import (
    WorkstreamID,
    WorkstreamStatus,
    WorkstreamSize,
    Workstream,
)


class TestWorkstreamID:
    """Test WorkstreamID parsing."""

    def test_parse_pp_fff_ss_format(self) -> None:
        """Verify parsing PP-FFF-SS format."""
        ws_id = WorkstreamID.parse("00-001-01")

        assert ws_id.project_id == 0
        assert ws_id.feature_id == 1
        assert ws_id.sequence == 1
        assert str(ws_id) == "00-001-01"

    def test_parse_legacy_ws_format(self) -> None:
        """Verify parsing legacy WS-FFF-SS format."""
        ws_id = WorkstreamID.parse("WS-500-01")

        assert ws_id.project_id == 0
        assert ws_id.feature_id == 500
        assert ws_id.sequence == 1
        assert str(ws_id) == "00-500-01"

    def test_parse_raises_invalid_format(self) -> None:
        """Verify parse raises for invalid format."""
        with pytest.raises(ValueError, match="Invalid WS ID format"):
            WorkstreamID.parse("invalid")

    def test_is_sdp_property(self) -> None:
        """Verify is_sdp for project 00."""
        ws_id = WorkstreamID.parse("00-001-01")
        assert ws_id.is_sdp is True

    def test_validate_project_id(self) -> None:
        """Verify validate_project_id accepts valid IDs."""
        ws_id = WorkstreamID.parse("00-001-01")
        ws_id.validate_project_id()

    def test_validate_project_id_raises_invalid(self) -> None:
        """Verify validate_project_id raises for invalid ID."""
        ws_id = WorkstreamID(project_id=99, feature_id=1, sequence=1)

        with pytest.raises(ValueError, match="Invalid project_id"):
            ws_id.validate_project_id()


class TestWorkstreamStatus:
    """Test WorkstreamStatus enum."""

    def test_status_values(self) -> None:
        """Verify status enum values."""
        assert WorkstreamStatus.BACKLOG.value == "backlog"
        assert WorkstreamStatus.ACTIVE.value == "active"
        assert WorkstreamStatus.COMPLETED.value == "completed"
        assert WorkstreamStatus.BLOCKED.value == "blocked"


class TestWorkstreamSize:
    """Test WorkstreamSize enum."""

    def test_size_values(self) -> None:
        """Verify size enum values."""
        assert WorkstreamSize.SMALL.value == "SMALL"
        assert WorkstreamSize.MEDIUM.value == "MEDIUM"
        assert WorkstreamSize.LARGE.value == "LARGE"


class TestParseWorkstream:
    """Test workstream file parsing."""

    def test_parse_valid_workstream(self, tmp_path: Path) -> None:
        """Verify parser extracts all fields correctly."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
project_id: 00
---

## WS-00-001-01: Test Workstream

### Goal

Test goal

### Context

Test context
""")

        ws = parse_workstream(ws_file)

        assert ws.ws_id == "00-001-01"
        assert ws.feature == "F001"
        assert ws.status == WorkstreamStatus.BACKLOG
        assert ws.size == WorkstreamSize.SMALL
        assert "Test goal" in ws.goal
        assert ws.file_path == ws_file

    def test_parse_raises_no_frontmatter(self, tmp_path: Path) -> None:
        """Verify parser raises when no frontmatter."""
        ws_file = tmp_path / "bad.md"
        ws_file.write_text("No frontmatter here")

        with pytest.raises(WorkstreamParseError, match="No frontmatter"):
            parse_workstream(ws_file)

    def test_parse_raises_invalid_ws_id(self, tmp_path: Path) -> None:
        """Verify parser raises for invalid ws_id format."""
        ws_file = tmp_path / "bad.md"
        ws_file.write_text("""---
ws_id: invalid-id
feature: F001
status: backlog
size: SMALL
---

## Bad
""")

        with pytest.raises(WorkstreamParseError, match="Invalid ws_id"):
            parse_workstream(ws_file)

    def test_parse_raises_missing_required_fields(self, tmp_path: Path) -> None:
        """Verify parser raises for missing required fields."""
        ws_file = tmp_path / "bad.md"
        ws_file.write_text("""---
ws_id: 00-001-01
---

## Bad
""")

        with pytest.raises(WorkstreamParseError, match="Missing required"):
            parse_workstream(ws_file)

    def test_parse_raises_invalid_status(self, tmp_path: Path) -> None:
        """Verify parser raises for invalid status."""
        ws_file = tmp_path / "bad.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: invalid_status
size: SMALL
---

## Bad
""")

        with pytest.raises(WorkstreamParseError, match="Invalid status"):
            parse_workstream(ws_file)

    def test_parse_raises_invalid_size(self, tmp_path: Path) -> None:
        """Verify parser raises for invalid size."""
        ws_file = tmp_path / "bad.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: INVALID_SIZE
---

## Bad
""")

        with pytest.raises(WorkstreamParseError, match="Invalid size"):
            parse_workstream(ws_file)

    def test_parse_with_depends_on_list(self, tmp_path: Path) -> None:
        """Verify parser merges depends_on list from frontmatter."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
depends_on:
  - 00-002-01
  - 00-003-01
---

## WS-00-001-01: Test Workstream

### Dependencies

00-004-01
""")

        ws = parse_workstream(ws_file)
        assert "00-002-01" in ws.dependencies
        assert "00-003-01" in ws.dependencies
        assert "00-004-01" in ws.dependencies

    def test_parse_with_depends_on_string(self, tmp_path: Path) -> None:
        """Verify parser merges depends_on string from frontmatter."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
depends_on: 00-002-01
---

## WS-00-001-01: Test Workstream

### Dependencies

00-003-01
""")

        ws = parse_workstream(ws_file)
        assert "00-002-01" in ws.dependencies
        assert "00-003-01" in ws.dependencies

    def test_parse_with_depends_on_duplicate_prevention(self, tmp_path: Path) -> None:
        """Verify parser prevents duplicate dependencies."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
depends_on: 00-002-01
---

## WS-00-001-01: Test Workstream

### Dependencies

00-002-01
""")

        ws = parse_workstream(ws_file)
        # Should only appear once
        assert ws.dependencies.count("00-002-01") == 1

    def test_parse_with_empty_depends_on(self, tmp_path: Path) -> None:
        """Verify parser handles empty depends_on values."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
depends_on: ""
---

## WS-00-001-01: Test Workstream
""")

        ws = parse_workstream(ws_file)
        assert ws.dependencies == []

    def test_parse_with_github_issue(self, tmp_path: Path) -> None:
        """Verify parser extracts github_issue."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
github_issue: 123
---

## WS-00-001-01: Test Workstream
""")

        ws = parse_workstream(ws_file)
        assert ws.github_issue == 123

    def test_parse_with_assignee(self, tmp_path: Path) -> None:
        """Verify parser extracts assignee."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
assignee: john@example.com
---

## WS-00-001-01: Test Workstream
""")

        ws = parse_workstream(ws_file)
        assert ws.assignee == "john@example.com"

    def test_parse_file_read_error(self, tmp_path: Path) -> None:
        """Verify parser handles file read errors."""
        ws_file = tmp_path / "00-001-01.md"
        # Create a directory with the same name to cause read error
        ws_file.mkdir()

        with pytest.raises((OSError, IsADirectoryError)):
            parse_workstream(ws_file)
