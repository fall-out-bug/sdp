"""Tests for checkpoint file manager edge cases."""

import json
import pytest
from pathlib import Path
from unittest.mock import Mock, patch, mock_open

from sdp.unified.orchestrator.checkpoint import CheckpointFileManager


@pytest.fixture
def manager(tmp_path):
    """Create checkpoint manager with temp directory."""
    return CheckpointFileManager(base_path=str(tmp_path / ".oneshot"))


class TestCheckpointFileManager:
    """Tests for CheckpointFileManager."""

    def test_init_creates_directory(self, tmp_path):
        """Test initialization creates checkpoint directory."""
        base_path = tmp_path / ".oneshot"
        assert not base_path.exists()

        manager = CheckpointFileManager(base_path=str(base_path))

        assert base_path.exists()
        assert base_path.is_dir()

    def test_save_checkpoint_success(self, manager, tmp_path):
        """Test saves checkpoint data successfully."""
        checkpoint_data = {
            "feature_id": "F001",
            "state": "in_progress",
            "workstreams": ["00-001-01"],
        }

        manager.save_checkpoint("F001", checkpoint_data)

        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        assert checkpoint_file.exists()
        
        loaded = json.loads(checkpoint_file.read_text())
        assert loaded == checkpoint_data

    def test_save_checkpoint_io_error(self, manager):
        """Test handles IO error during save."""
        checkpoint_data = {"feature_id": "F001"}

        with patch("builtins.open", side_effect=OSError("Permission denied")):
            with pytest.raises(OSError, match="Permission denied"):
                manager.save_checkpoint("F001", checkpoint_data)

    def test_load_checkpoint_success(self, manager, tmp_path):
        """Test loads checkpoint data successfully."""
        checkpoint_data = {
            "feature_id": "F001",
            "state": "in_progress",
        }
        
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text(json.dumps(checkpoint_data))

        loaded = manager.load_checkpoint("F001")

        assert loaded == checkpoint_data

    def test_load_checkpoint_not_found(self, manager):
        """Test returns None when checkpoint file not found."""
        result = manager.load_checkpoint("F999")

        assert result is None

    def test_load_checkpoint_json_error(self, manager, tmp_path):
        """Test handles JSON decode error."""
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text("invalid json {")

        result = manager.load_checkpoint("F001")

        assert result is None

    def test_load_checkpoint_io_error(self, manager, tmp_path):
        """Test handles IO error during load."""
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text('{"test": "data"}')

        with patch("builtins.open", side_effect=OSError("Read error")):
            result = manager.load_checkpoint("F001")

        assert result is None

    def test_delete_checkpoint_success(self, manager, tmp_path):
        """Test deletes checkpoint file successfully."""
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text('{"test": "data"}')
        assert checkpoint_file.exists()

        manager.delete_checkpoint("F001")

        assert not checkpoint_file.exists()

    def test_delete_checkpoint_not_found(self, manager):
        """Test handles deletion when file doesn't exist."""
        # Should not raise
        manager.delete_checkpoint("F999")

    def test_delete_checkpoint_io_error(self, manager, tmp_path):
        """Test handles IO error during deletion."""
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text('{"test": "data"}')

        with patch.object(Path, 'unlink', side_effect=OSError("Permission denied")):
            # Should not raise, just log error
            manager.delete_checkpoint("F001")

    def test_checkpoint_exists_true(self, manager, tmp_path):
        """Test returns True when checkpoint exists."""
        checkpoint_file = tmp_path / ".oneshot" / "F001-checkpoint.json"
        checkpoint_file.write_text('{"test": "data"}')

        result = manager.checkpoint_exists("F001")

        assert result is True

    def test_checkpoint_exists_false(self, manager):
        """Test returns False when checkpoint doesn't exist."""
        result = manager.checkpoint_exists("F999")

        assert result is False
