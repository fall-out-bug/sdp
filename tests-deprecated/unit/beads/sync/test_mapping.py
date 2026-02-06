"""Unit tests for beads/sync/mapping.py

Tests for MappingManager and resolve_ws_id_to_beads_id to increase coverage.
"""

import json
from pathlib import Path

import pytest

from sdp.beads.sync.mapping import BeadsSyncError, MappingManager, resolve_ws_id_to_beads_id


class TestResolveWsIdToBeadsId:
    """Test resolve_ws_id_to_beads_id function."""

    def test_resolve_with_empty_lines(self, tmp_path: Path) -> None:
        """Handle empty lines in mapping file."""
        mapping = tmp_path / "mapping.jsonl"
        mapping.write_text('\n{"sdp_id": "00-001-01", "beads_id": "bd-abc"}\n\n')
        
        assert resolve_ws_id_to_beads_id("00-001-01", mapping) == "bd-abc"

    def test_resolve_json_decode_error(self, tmp_path: Path) -> None:
        """Handle malformed JSON gracefully."""
        mapping = tmp_path / "mapping.jsonl"
        mapping.write_text("not json\n")
        
        assert resolve_ws_id_to_beads_id("00-001-01", mapping) is None

    def test_resolve_io_error(self, tmp_path: Path) -> None:
        """Handle IO errors gracefully."""
        # Directory instead of file
        mapping = tmp_path / "baddir"
        mapping.mkdir()
        
        assert resolve_ws_id_to_beads_id("00-001-01", mapping) is None

    def test_resolve_non_string_beads_id(self, tmp_path: Path) -> None:
        """Handle non-string beads_id values."""
        mapping = tmp_path / "mapping.jsonl"
        mapping.write_text('{"sdp_id": "00-001-01", "beads_id": 123}\n')
        
        assert resolve_ws_id_to_beads_id("00-001-01", mapping) is None


class TestMappingManager:
    """Test MappingManager class."""

    def test_load_nonexistent_file(self, tmp_path: Path) -> None:
        """Load handles nonexistent file gracefully."""
        mapping_file = tmp_path / "nonexistent.jsonl"
        manager = MappingManager(mapping_file)
        
        manager.load()  # Should not raise
        
        assert manager.get_beads_id("00-001-01") is None

    def test_load_existing_mappings(self, tmp_path: Path) -> None:
        """Load existing mappings from file."""
        mapping_file = tmp_path / "mapping.jsonl"
        mapping_file.write_text(
            '{"sdp_id": "00-001-01", "beads_id": "bd-abc"}\n'
            '{"sdp_id": "00-001-02", "beads_id": "bd-def"}\n'
        )
        
        manager = MappingManager(mapping_file)
        manager.load()
        
        assert manager.get_beads_id("00-001-01") == "bd-abc"
        assert manager.get_beads_id("00-001-02") == "bd-def"
        assert manager.get_sdp_id("bd-abc") == "00-001-01"
        assert manager.get_sdp_id("bd-def") == "00-001-02"

    def test_load_malformed_json_raises_error(self, tmp_path: Path) -> None:
        """Load raises BeadsSyncError on malformed JSON."""
        mapping_file = tmp_path / "bad.jsonl"
        mapping_file.write_text("not json")
        
        manager = MappingManager(mapping_file)
        
        with pytest.raises(BeadsSyncError, match="Failed to load mapping file"):
            manager.load()

    def test_load_skips_empty_lines(self, tmp_path: Path) -> None:
        """Load skips empty lines."""
        mapping_file = tmp_path / "mapping.jsonl"
        mapping_file.write_text(
            '{"sdp_id": "00-001-01", "beads_id": "bd-abc"}\n'
            '\n'
            '{"sdp_id": "00-001-02", "beads_id": "bd-def"}\n'
        )
        
        manager = MappingManager(mapping_file)
        manager.load()
        
        assert len(manager._mapping) == 2

    def test_load_skips_incomplete_entries(self, tmp_path: Path) -> None:
        """Load skips entries missing sdp_id or beads_id."""
        mapping_file = tmp_path / "mapping.jsonl"
        mapping_file.write_text(
            '{"sdp_id": "00-001-01"}\n'  # Missing beads_id
            '{"beads_id": "bd-def"}\n'    # Missing sdp_id
            '{"sdp_id": "00-001-03", "beads_id": "bd-ghi"}\n'
        )
        
        manager = MappingManager(mapping_file)
        manager.load()
        
        assert len(manager._mapping) == 1
        assert manager.get_beads_id("00-001-03") == "bd-ghi"

    def test_save_creates_file(self, tmp_path: Path) -> None:
        """Save creates mapping file with current state."""
        mapping_file = tmp_path / "mapping.jsonl"
        manager = MappingManager(mapping_file)
        
        manager.add_mapping("00-001-01", "bd-abc")
        manager.add_mapping("00-001-02", "bd-def")
        manager.save()
        
        assert mapping_file.exists()
        lines = mapping_file.read_text().strip().split("\n")
        assert len(lines) == 2
        
        entries = [json.loads(line) for line in lines]
        assert any(e["sdp_id"] == "00-001-01" and e["beads_id"] == "bd-abc" for e in entries)

    def test_save_overwrites_existing_file(self, tmp_path: Path) -> None:
        """Save overwrites existing file."""
        mapping_file = tmp_path / "mapping.jsonl"
        mapping_file.write_text('{"old": "data"}\n')
        
        manager = MappingManager(mapping_file)
        manager.add_mapping("00-001-01", "bd-abc")
        manager.save()
        
        content = mapping_file.read_text()
        assert "old" not in content
        assert "00-001-01" in content

    def test_save_io_error_raises(self, tmp_path: Path, monkeypatch) -> None:
        """Save raises BeadsSyncError on IO failure."""
        mapping_file = tmp_path / "mapping.jsonl"
        manager = MappingManager(mapping_file)
        manager.add_mapping("00-001-01", "bd-abc")
        
        def mock_open(*args, **kwargs):
            raise IOError("Permission denied")
        
        monkeypatch.setattr("builtins.open", mock_open)
        
        with pytest.raises(BeadsSyncError, match="Failed to save mapping file"):
            manager.save()

    def test_get_beads_id_not_found(self, tmp_path: Path) -> None:
        """get_beads_id returns None for unknown sdp_id."""
        manager = MappingManager(tmp_path / "mapping.jsonl")
        
        assert manager.get_beads_id("00-999-99") is None

    def test_get_sdp_id_not_found(self, tmp_path: Path) -> None:
        """get_sdp_id returns None for unknown beads_id."""
        manager = MappingManager(tmp_path / "mapping.jsonl")
        
        assert manager.get_sdp_id("bd-unknown") is None

    def test_add_mapping_bidirectional(self, tmp_path: Path) -> None:
        """add_mapping creates bidirectional mapping."""
        manager = MappingManager(tmp_path / "mapping.jsonl")
        
        manager.add_mapping("00-001-01", "bd-abc")
        
        assert manager.get_beads_id("00-001-01") == "bd-abc"
        assert manager.get_sdp_id("bd-abc") == "00-001-01"

    def test_add_mapping_updates_existing(self, tmp_path: Path) -> None:
        """add_mapping updates existing mapping."""
        manager = MappingManager(tmp_path / "mapping.jsonl")
        
        manager.add_mapping("00-001-01", "bd-abc")
        manager.add_mapping("00-001-01", "bd-xyz")  # Update
        
        assert manager.get_beads_id("00-001-01") == "bd-xyz"
        assert manager.get_sdp_id("bd-xyz") == "00-001-01"
        # Note: Old reverse mapping persists (bd-abc still points to 00-001-01)
