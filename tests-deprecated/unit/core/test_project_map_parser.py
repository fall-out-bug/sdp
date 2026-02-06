"""Tests for project map parsing."""

import pytest
from pathlib import Path

from sdp.core.project_map_parser import parse_project_map
from sdp.core.project_map_types import ProjectMapParseError


class TestParseProjectMap:
    """Test project map parsing."""

    def test_parse_valid_project_map(self) -> None:
        """Verify parse_project_map loads valid file."""
        project_map_file = (
            Path(__file__).parent.parent.parent.parent / "docs" / "PROJECT_MAP.md"
        )

        result = parse_project_map(project_map_file)

        assert result is not None
        assert result.project_name == "SDP"
        assert result.file_path == project_map_file

    def test_parse_raises_file_not_found(self) -> None:
        """Verify parse_project_map raises when file missing."""
        with pytest.raises(FileNotFoundError, match="Project map file not found"):
            parse_project_map(Path("/nonexistent/PROJECT_MAP.md"))

    def test_parse_extracts_decisions(self) -> None:
        """Verify parse_project_map extracts decisions if present."""
        project_map_file = (
            Path(__file__).parent.parent.parent.parent / "docs" / "PROJECT_MAP.md"
        )

        result = parse_project_map(project_map_file)

        assert hasattr(result, "decisions")
        assert isinstance(result.decisions, list)
