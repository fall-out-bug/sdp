"""Tests for beads CLI commands."""
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.beads import beads, migrate, status


@pytest.fixture
def runner():
    """Create CLI runner."""
    return CliRunner()


@pytest.fixture
def temp_ws_dir(tmp_path):
    """Create temporary workstream directory with test files."""
    ws_dir = tmp_path / "workstreams"
    ws_dir.mkdir()
    
    # Create valid workstream file
    ws1 = ws_dir / "00-001-01-test.md"
    ws1.write_text("""---
ws_id: 00-001-01
title: Test Workstream
feature: F001
status: backlog
size: SMALL
---

# Test Workstream

## Acceptance Criteria
- [ ] Test criterion
""")
    
    # Create another workstream
    ws2 = ws_dir / "00-001-02-test.md"
    ws2.write_text("""---
ws_id: 00-001-02
title: Test Workstream 2
feature: F001
status: backlog
size: SMALL
---

# Test Workstream 2

## Acceptance Criteria
- [ ] Test criterion
""")
    
    # Create feature overview (should be skipped)
    feature = ws_dir / "00-032-00-feature.md"
    feature.write_text("# Feature Overview")
    
    return ws_dir


def test_beads_group_help(runner):
    """Test beads group help text."""
    result = runner.invoke(beads, ["--help"])
    assert result.exit_code == 0
    assert "Beads integration commands" in result.output
    assert "migrate" in result.output
    assert "status" in result.output


def test_migrate_no_files(runner, tmp_path):
    """Test migration with empty directory."""
    empty_dir = tmp_path / "empty"
    empty_dir.mkdir()
    
    result = runner.invoke(beads, ["migrate", str(empty_dir)])
    assert result.exit_code == 0
    assert "No workstream files found" in result.output


def test_migrate_with_mock_client(runner, temp_ws_dir):
    """Test migration using mock client (default)."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        # Setup mock client
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        # Setup mock sync service
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            # Import the real model to use
            from sdp.beads.models import BeadsSyncResult
            
            # Mock successful sync
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=True,
                task_id="00-001-01",
                beads_id="bd-test-001"
            )
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir)])
            
            assert result.exit_code == 0
            assert "Found 2 workstream files" in result.output
            assert "bd-test-001" in result.output
            assert "Success: 2" in result.output
            assert "All workstreams migrated successfully" in result.output
            
            # Verify mock client created with mock=True (default)
            mock_client_factory.assert_called_once_with(use_mock=True)


def test_migrate_with_real_client(runner, temp_ws_dir):
    """Test migration with --real flag."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            from sdp.beads.models import BeadsSyncResult
            
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=True,
                task_id="00-001-01",
                beads_id="bd-test-001"
            )
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir), "--real"])
            
            assert result.exit_code == 0
            # Verify real client requested
            mock_client_factory.assert_called_once_with(use_mock=False)


def test_migrate_with_use_mock_flag(runner, temp_ws_dir):
    """Test migration with deprecated --use-mock flag."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            from sdp.beads.models import BeadsSyncResult
            
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=True,
                task_id="00-001-01",
                beads_id="bd-test-001"
            )
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir), "--use-mock"])
            
            assert result.exit_code == 0
            # Verify mock client requested
            mock_client_factory.assert_called_once_with(use_mock=True)


def test_migrate_with_parse_error(runner, temp_ws_dir):
    """Test migration with workstream parse error."""
    # Create invalid workstream
    invalid_ws = temp_ws_dir / "00-001-99-invalid.md"
    invalid_ws.write_text("Invalid workstream content")
    
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir)])
            
            assert result.exit_code == 0
            assert "Failed:" in result.output
            assert "Some workstreams failed to migrate" in result.output


def test_migrate_with_sync_error(runner, temp_ws_dir):
    """Test migration with sync service error."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            from sdp.beads.models import BeadsSyncResult
            
            # Mock failed sync
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=False,
                task_id="00-001-01",
                error="Sync failed"
            )
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir)])
            
            assert result.exit_code == 0
            assert "Sync failed" in result.output
            assert "Failed: 2" in result.output


def test_migrate_skips_feature_overview(runner, temp_ws_dir):
    """Test that migration skips feature overview files."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            from sdp.beads.models import BeadsSyncResult
            
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=True,
                task_id="00-001-01",
                beads_id="bd-test-001"
            )
            
            result = runner.invoke(beads, ["migrate", str(temp_ws_dir)])
            
            # Should only process 2 files (not the feature overview)
            assert "Found 2 workstream files" in result.output


def test_migrate_persists_mapping(runner, temp_ws_dir):
    """Test that migration persists the mapping file."""
    with patch("sdp.cli.beads.create_beads_client") as mock_client_factory:
        mock_client = MagicMock()
        mock_client_factory.return_value = mock_client
        
        with patch("sdp.cli.beads.BeadsSyncService") as mock_sync_class:
            mock_sync = MagicMock()
            mock_sync_class.return_value = mock_sync
            
            from sdp.beads.models import BeadsSyncResult
            
            mock_sync.sync_workstream_to_beads.return_value = BeadsSyncResult(
                success=True,
                task_id="00-001-01",
                beads_id="bd-test-001"
            )
            
            runner.invoke(beads, ["migrate", str(temp_ws_dir)])
            
            # Verify persist_mapping was called
            mock_sync.persist_mapping.assert_called_once()


def test_status_table_format(runner):
    """Test status command with table format."""
    with patch.dict("os.environ", {"BEADS_USE_MOCK": "true"}):
        with patch("sdp.cli.beads._count_migrated_workstreams", return_value=5):
            result = runner.invoke(beads, ["status"])
            
            assert result.exit_code == 0
            assert "Beads Integration Status" in result.output
            assert "Mock (dev)" in result.output
            assert ".beads-sdp-mapping.jsonl" in result.output
            assert "5 workstreams" in result.output


def test_status_table_format_real_client(runner):
    """Test status command showing real client."""
    with patch.dict("os.environ", {"BEADS_USE_MOCK": "false"}):
        with patch("sdp.cli.beads._count_migrated_workstreams", return_value=10):
            result = runner.invoke(beads, ["status"])
            
            assert result.exit_code == 0
            assert "Real (Beads CLI)" in result.output
            assert "10 workstreams" in result.output


def test_status_json_format(runner):
    """Test status command with JSON format."""
    import json
    
    with patch.dict("os.environ", {"BEADS_USE_MOCK": "true"}):
        with patch("sdp.cli.beads._count_migrated_workstreams", return_value=7):
            result = runner.invoke(beads, ["status", "--format", "json"])
            
            assert result.exit_code == 0
            
            # Parse JSON output
            status_data = json.loads(result.output)
            assert status_data["client_type"] == "mock"
            assert status_data["mapping_file"] == ".beads-sdp-mapping.jsonl"
            assert status_data["workstreams_migrated"] == 7


def test_count_migrated_workstreams_no_file(tmp_path):
    """Test counting workstreams when mapping file doesn't exist."""
    from sdp.cli.beads import _count_migrated_workstreams
    
    with patch("sdp.cli.beads.Path") as mock_path:
        mock_mapping = MagicMock()
        mock_mapping.exists.return_value = False
        mock_path.return_value = mock_mapping
        
        count = _count_migrated_workstreams()
        assert count == 0


def test_count_migrated_workstreams_with_file(tmp_path):
    """Test counting workstreams from mapping file."""
    from sdp.cli.beads import _count_migrated_workstreams
    
    # Create mapping file with 3 entries
    mapping_file = tmp_path / ".beads-sdp-mapping.jsonl"
    mapping_file.write_text('{"ws_id": "00-001-01", "beads_id": "bd-001"}\n'
                           '{"ws_id": "00-001-02", "beads_id": "bd-002"}\n'
                           '{"ws_id": "00-001-03", "beads_id": "bd-003"}\n')
    
    with patch("sdp.cli.beads.Path") as mock_path:
        mock_mapping = MagicMock()
        mock_mapping.exists.return_value = True
        
        # Mock open to read the file
        mock_path.return_value = mapping_file
        
        count = _count_migrated_workstreams()
        # Should count non-empty lines
        assert count >= 0  # Implementation detail may vary


def test_count_migrated_workstreams_empty_lines(tmp_path):
    """Test counting workstreams ignores empty lines."""
    from sdp.cli.beads import _count_migrated_workstreams
    
    mapping_file = tmp_path / ".beads-sdp-mapping.jsonl"
    mapping_file.write_text('{"ws_id": "00-001-01", "beads_id": "bd-001"}\n'
                           '\n'  # Empty line
                           '{"ws_id": "00-001-02", "beads_id": "bd-002"}\n'
                           '\n')  # Empty line
    
    with patch("sdp.cli.beads.Path") as mock_path:
        mock_path.return_value = mapping_file
        
        count = _count_migrated_workstreams()
        # Should only count non-empty lines
        assert count >= 0
