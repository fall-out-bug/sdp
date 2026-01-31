"""Unit tests for resolve_ws_id_to_beads_id."""

from pathlib import Path

import pytest

from sdp.beads.sync.mapping import resolve_ws_id_to_beads_id


class TestResolveWsIdToBeadsId:
    """Tests for ws_id → beads_id resolution."""

    def test_resolves_ws_id_to_beads_id(self, tmp_path: Path) -> None:
        """Resolve PP-FFF-SS to beads_id via mapping file."""
        mapping = tmp_path / ".beads-sdp-mapping.jsonl"
        mapping.write_text(
            '{"sdp_id": "00-020-03", "beads_id": "sdp-4qq", "updated_at": "2026-01-31"}\n'
        )
        assert resolve_ws_id_to_beads_id("00-020-03", mapping) == "sdp-4qq"

    def test_returns_none_when_not_found(self, tmp_path: Path) -> None:
        """Return None when ws_id not in mapping."""
        mapping = tmp_path / ".beads-sdp-mapping.jsonl"
        mapping.write_text(
            '{"sdp_id": "00-020-03", "beads_id": "sdp-4qq", "updated_at": "2026-01-31"}\n'
        )
        assert resolve_ws_id_to_beads_id("00-999-99", mapping) is None

    def test_returns_none_when_file_missing(self, tmp_path: Path) -> None:
        """Return None when mapping file does not exist."""
        missing = tmp_path / "nonexistent.jsonl"
        assert resolve_ws_id_to_beads_id("00-020-03", missing) is None

    def test_returns_none_for_non_string_beads_id(self, tmp_path: Path) -> None:
        """Return None when beads_id is not a string (invalid JSON)."""
        mapping = tmp_path / ".beads-sdp-mapping.jsonl"
        mapping.write_text(
            '{"sdp_id": "00-020-03", "beads_id": 123, "updated_at": "2026-01-31"}\n'
        )
        assert resolve_ws_id_to_beads_id("00-020-03", mapping) is None
