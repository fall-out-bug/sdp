"""Tests for core exception SDPError compliance (F031)."""

import json
import pytest
from pathlib import Path

from sdp.core.workstream import WorkstreamParseError
from sdp.core.feature import CircularDependencyError, MissingDependencyError
from sdp.core.model_mapping import ModelMappingError
from sdp.errors import SDPError


class TestWorkstreamParseError:
    """Test WorkstreamParseError SDPError compliance."""

    def test_includes_error_category(self) -> None:
        """AC1: Verify WorkstreamParseError includes VALIDATION category."""
        error = WorkstreamParseError("Invalid ws_id")

        assert error.category.value == "validation"
        assert "VALIDATION" in error.format_terminal()

    def test_includes_remediation_steps(self) -> None:
        """AC4: Verify error includes actionable remediation steps."""
        error = WorkstreamParseError("Invalid ws_id")

        formatted = error.format_terminal()
        assert "💡" in formatted or "Remediation" in formatted
        assert "1. Check WS ID format" in formatted
        assert "PP-FFF-SS" in formatted

    def test_includes_docs_url(self) -> None:
        """Verify error includes documentation link."""
        error = WorkstreamParseError("Invalid ws_id")

        formatted = error.format_terminal()
        assert "Docs" in formatted or "Documentation" in formatted
        assert "sdp.dev" in formatted

    def test_includes_context_when_provided(self) -> None:
        """Verify error includes context when available."""
        error = WorkstreamParseError(
            message="Invalid ws_id",
            file_path=Path("/path/to/ws.md"),
        )

        formatted = error.format_terminal()
        assert "Context" in formatted
        assert "/path/to/ws.md" in formatted

    def test_formats_json_for_machine_parsing(self) -> None:
        """AC5/AC6: Verify error formats as JSON for CI/CD systems."""
        error = WorkstreamParseError("Invalid ws_id")

        json_str = error.format_json()
        parsed = json.loads(json_str)

        assert parsed["category"] == "validation"
        assert parsed["message"] == "Invalid ws_id"
        assert "remediation" in parsed
        assert "docs_url" in parsed


class TestCircularDependencyError:
    """Test CircularDependencyError SDPError compliance."""

    def test_shows_cycle_in_error_message(self) -> None:
        """AC2: Verify CircularDependencyError shows the dependency cycle."""
        error = CircularDependencyError(
            ws_id="00-001-01",
            cycle=["00-001-02", "00-001-03", "00-001-01"],
        )

        formatted = error.format_terminal()
        assert "00-001-02" in formatted
        assert "00-001-03" in formatted
        assert "00-001-01" in formatted

    def test_suggests_breaking_cycle(self) -> None:
        """Verify remediation suggests breaking the cycle."""
        error = CircularDependencyError(
            ws_id="00-001-01",
            cycle=["00-001-02", "00-001-03", "00-001-01"],
        )

        formatted = error.format_terminal()
        assert "Break the cycle" in formatted or "cycle" in formatted
        assert "dependency" in formatted


class TestMissingDependencyError:
    """Test MissingDependencyError SDPError compliance."""

    def test_shows_available_alternatives(self) -> None:
        """AC3: Verify MissingDependencyError shows available workstreams."""
        error = MissingDependencyError(
            ws_id="00-001-01",
            missing_dep="00-001-02",
            available_workstreams=["00-001-03", "00-001-04"],
        )

        formatted = error.format_terminal()
        assert "00-001-03" in formatted
        assert "00-001-04" in formatted


class TestModelMappingError:
    """Test ModelMappingError SDPError compliance."""

    def test_inherits_from_sdperror(self) -> None:
        """AC6: Verify ModelMappingError inherits from SDPError."""
        error = ModelMappingError(
            message="File not found",
            mapping_file=Path("/path/to/mapping.md"),
        )

        assert isinstance(error, SDPError)
        assert error.category.value == "configuration"
