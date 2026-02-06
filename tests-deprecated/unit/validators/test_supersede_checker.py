"""Tests for supersede_checker validator."""

from pathlib import Path

import pytest

from sdp.validators.supersede_checker import (
    SupersedeChain,
    SupersedeResult,
    SupersedeValidator,
    ValidationReport,
)


@pytest.fixture
def ws_dir(tmp_path: Path) -> Path:
    """Create temp workstream directory structure."""
    for sub in ("backlog", "completed"):
        (tmp_path / sub).mkdir()
    return tmp_path


@pytest.fixture
def validator(ws_dir: Path) -> SupersedeValidator:
    """Create validator with temp ws_dir."""
    return SupersedeValidator(ws_dir=ws_dir)


class TestSupersedeValidator:
    """Tests for SupersedeValidator."""

    def test_supersede_old_not_found_returns_error(
        self, validator: SupersedeValidator, ws_dir: Path
    ) -> None:
        """supersede when old_ws not found returns error."""
        (ws_dir / "completed" / "00-032-02-new.md").write_text("---\nws_id: 00-032-02\n---")
        result = validator.supersede("00-032-99", "00-032-02")
        assert not result.success
        assert "not found" in (result.error or "").lower()

    def test_supersede_new_not_found_returns_error(
        self, validator: SupersedeValidator, ws_dir: Path
    ) -> None:
        """supersede when new_ws not found returns error."""
        (ws_dir / "completed" / "00-032-01-old.md").write_text("---\nws_id: 00-032-01\n---")
        result = validator.supersede("00-032-01", "00-032-99")
        assert not result.success
        assert "replacement" in (result.error or "").lower()

    def test_supersede_success_updates_frontmatter(
        self, validator: SupersedeValidator, ws_dir: Path
    ) -> None:
        """supersede updates old WS frontmatter."""
        (ws_dir / "completed" / "00-032-01-old.md").write_text(
            "---\nws_id: 00-032-01\nstatus: completed\n---\n\n# Old"
        )
        (ws_dir / "completed" / "00-032-02-new.md").write_text(
            "---\nws_id: 00-032-02\nstatus: completed\n---\n\n# New"
        )
        result = validator.supersede("00-032-01", "00-032-02")
        assert result.success
        content = (ws_dir / "completed" / "00-032-01-old.md").read_text()
        assert "superseded" in content
        assert "00-032-02" in content

    def test_trace_chain_final_ws(self, validator: SupersedeValidator, ws_dir: Path) -> None:
        """trace_chain returns final WS when not superseded."""
        (ws_dir / "completed" / "00-032-01.md").write_text(
            "---\nws_id: 00-032-01\nstatus: completed\n---"
        )
        chain = validator.trace_chain("00-032-01")
        assert not chain.has_cycle
        assert chain.final_ws == "00-032-01"

    def test_trace_chain_not_found(self, validator: SupersedeValidator) -> None:
        """trace_chain returns None final_ws when WS not found."""
        chain = validator.trace_chain("99-999-99")
        assert not chain.has_cycle
        assert chain.final_ws is None

    def test_find_orphans_empty_dir(self, validator: SupersedeValidator) -> None:
        """find_orphans returns empty for dir with no superseded."""
        assert validator.find_orphans() == []

    def test_find_orphans_finds_missing_replacement(
        self, validator: SupersedeValidator, ws_dir: Path
    ) -> None:
        """find_orphans finds superseded WS without valid replacement."""
        (ws_dir / "completed" / "00-032-01.md").write_text(
            "---\nws_id: 00-032-01\nstatus: superseded\nsuperseded_by: 00-032-99\n---"
        )
        orphans = validator.find_orphans()
        assert "00-032-01" in orphans

    def test_validate_all_empty(self, validator: SupersedeValidator) -> None:
        """validate_all returns empty report for no superseded."""
        report = validator.validate_all()
        assert report.total_superseded == 0
        assert report.orphans == []
        assert report.cycles == []
