"""Tests for feature loading utilities."""

import pytest
from pathlib import Path
from unittest.mock import patch, mock_open

from sdp.core.feature.loader import (
    load_feature_from_directory,
    load_feature_from_spec,
)
from sdp.core.workstream import WorkstreamParseError


class TestLoadFeatureFromDirectory:
    """Test loading feature from directory."""

    def test_load_feature_success(self, tmp_path: Path) -> None:
        """Verify successful feature loading."""
        ws_file1 = tmp_path / "00-001-01.md"
        ws_file1.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-01: First Workstream
""")

        ws_file2 = tmp_path / "00-001-02.md"
        ws_file2.write_text("""---
ws_id: 00-001-02
feature: F001
status: backlog
size: MEDIUM
---

## WS-00-001-02: Second Workstream
""")

        feature = load_feature_from_directory("F001", tmp_path, pattern="00-*.md")
        assert feature.feature_id == "F001"
        assert len(feature.workstreams) == 2
        assert feature.workstreams[0].ws_id == "00-001-01"
        assert feature.workstreams[1].ws_id == "00-001-02"

    def test_load_feature_no_files(self, tmp_path: Path) -> None:
        """Verify error when no workstream files found."""
        with pytest.raises(ValueError, match="No workstream files found"):
            load_feature_from_directory("F001", tmp_path)

    def test_load_feature_custom_pattern(self, tmp_path: Path) -> None:
        """Verify custom glob pattern works."""
        ws_file = tmp_path / "WS-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-01: Workstream
""")

        feature = load_feature_from_directory("F001", tmp_path, pattern="WS-*.md")
        assert len(feature.workstreams) == 1

    def test_load_feature_feature_mismatch(self, tmp_path: Path) -> None:
        """Verify error when workstream feature doesn't match."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F002
status: backlog
size: SMALL
---

## WS-00-001-01: Workstream
""")

        with pytest.raises(ValueError, match="has feature F002, expected F001"):
            load_feature_from_directory("F001", tmp_path, pattern="00-*.md")

    def test_load_feature_parse_error(self, tmp_path: Path) -> None:
        """Verify WorkstreamParseError is raised for malformed files."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("Invalid markdown without frontmatter")

        with pytest.raises(WorkstreamParseError, match="Failed to parse"):
            load_feature_from_directory("F001", tmp_path, pattern="00-*.md")

    def test_load_feature_multiple_parse_errors(self, tmp_path: Path) -> None:
        """Verify first parse error is raised."""
        ws_file1 = tmp_path / "00-001-01.md"
        ws_file1.write_text("Invalid markdown")

        ws_file2 = tmp_path / "00-001-02.md"
        ws_file2.write_text("Also invalid")

        with pytest.raises(WorkstreamParseError):
            load_feature_from_directory("F001", tmp_path, pattern="00-*.md")

    def test_load_feature_sorted_files(self, tmp_path: Path) -> None:
        """Verify workstreams are loaded in sorted order."""
        ws_file3 = tmp_path / "00-001-03.md"
        ws_file3.write_text("""---
ws_id: 00-001-03
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-03: Third
""")

        ws_file1 = tmp_path / "00-001-01.md"
        ws_file1.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-01: First
""")

        ws_file2 = tmp_path / "00-001-02.md"
        ws_file2.write_text("""---
ws_id: 00-001-02
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-02: Second
""")

        feature = load_feature_from_directory("F001", tmp_path, pattern="00-*.md")
        assert feature.workstreams[0].ws_id == "00-001-01"
        assert feature.workstreams[1].ws_id == "00-001-02"
        assert feature.workstreams[2].ws_id == "00-001-03"


class TestLoadFeatureFromSpec:
    """Test loading feature from spec file."""

    def test_load_feature_from_spec_success(self, tmp_path: Path) -> None:
        """Verify successful loading from spec file."""
        # Create directory structure matching expected layout
        # spec_file.parent.parent goes up 2 levels: specs -> path -> some
        specs_dir = tmp_path / "some" / "path" / "specs"
        specs_dir.mkdir(parents=True)
        # workstreams_dir should be at spec_file.parent.parent / "workstreams" / "backlog"
        # which is tmp_path / "some" / "workstreams" / "backlog"
        workstreams_dir = tmp_path / "some" / "workstreams" / "backlog"
        workstreams_dir.mkdir(parents=True)

        spec_file = specs_dir / "F001.md"
        spec_file.write_text("# Feature F001")

        ws_file = workstreams_dir / "WS-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F001
status: backlog
size: SMALL
---

## WS-00-001-01: Workstream
""")

        feature = load_feature_from_spec("F001", spec_file)
        assert feature.feature_id == "F001"
        assert len(feature.workstreams) == 1

    def test_load_feature_from_spec_missing_directory(self, tmp_path: Path) -> None:
        """Verify error when workstreams directory doesn't exist."""
        specs_dir = tmp_path / "specs"
        specs_dir.mkdir()
        spec_file = specs_dir / "F001.md"
        spec_file.write_text("# Feature F001")

        with pytest.raises(ValueError, match="Workstreams directory not found"):
            load_feature_from_spec("F001", spec_file)

    def test_load_feature_from_spec_no_workstreams(self, tmp_path: Path) -> None:
        """Verify error when no workstreams found."""
        specs_dir = tmp_path / "some" / "path" / "specs"
        specs_dir.mkdir(parents=True)
        workstreams_dir = tmp_path / "some" / "workstreams" / "backlog"
        workstreams_dir.mkdir(parents=True)

        spec_file = specs_dir / "F001.md"
        spec_file.write_text("# Feature F001")

        with pytest.raises(ValueError, match="No workstream files found"):
            load_feature_from_spec("F001", spec_file)

    def test_load_feature_from_spec_feature_mismatch(self, tmp_path: Path) -> None:
        """Verify error when workstream feature doesn't match."""
        specs_dir = tmp_path / "some" / "path" / "specs"
        specs_dir.mkdir(parents=True)
        workstreams_dir = tmp_path / "some" / "workstreams" / "backlog"
        workstreams_dir.mkdir(parents=True)

        spec_file = specs_dir / "F001.md"
        spec_file.write_text("# Feature F001")

        ws_file = workstreams_dir / "WS-001-01.md"
        ws_file.write_text("""---
ws_id: 00-001-01
feature: F002
status: backlog
size: SMALL
---

## WS-00-001-01: Workstream
""")

        with pytest.raises(ValueError, match="has feature F002, expected F001"):
            load_feature_from_spec("F001", spec_file)

    def test_load_feature_from_spec_parse_error(self, tmp_path: Path) -> None:
        """Verify parse errors are propagated."""
        specs_dir = tmp_path / "some" / "path" / "specs"
        specs_dir.mkdir(parents=True)
        workstreams_dir = tmp_path / "some" / "workstreams" / "backlog"
        workstreams_dir.mkdir(parents=True)

        spec_file = specs_dir / "F001.md"
        spec_file.write_text("# Feature F001")

        ws_file = workstreams_dir / "WS-001-01.md"
        ws_file.write_text("Invalid markdown")

        with pytest.raises(WorkstreamParseError):
            load_feature_from_spec("F001", spec_file)
