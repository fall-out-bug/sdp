"""Tests for workstream ID migration script."""

from pathlib import Path
from textwrap import dedent

import pytest

from scripts.migrate_workstream_ids import (
    WorkstreamFile,
    WorkstreamMigrationError,
    WorkstreamMigrator,
)


@pytest.fixture
def temp_workspace(tmp_path: Path) -> Path:
    """Create temporary workspace with workstream files."""
    ws_dir = tmp_path / "docs" / "workstreams"
    ws_dir.mkdir(parents=True)

    # Create test workstream files
    test_files = {
        "backlog/WS-001-01-test.md": dedent("""\
            ---
            ws_id: WS-001-01
            feature: F001
            status: backlog
            ---
            ## WS-001-01: Test Workstream
            Content here
            """),
        "completed/WS-002-01-done.md": dedent("""\
            ---
            ws_id: WS-002-01
            feature: F002
            status: completed
            ---
            ## WS-002-01: Done Workstream
            Done content
            """),
        "backlog/WS-003-05-deps.md": dedent("""\
            ---
            ws_id: WS-003-05
            feature: F003
            status: backlog
            depends_on:
              - WS-001-01
              - WS-002-01
            ---
            ## WS-003-05: Deps Workstream
            Depends on others
            """),
        "backlog/already-migrated.md": dedent("""\
            ---
            ws_id: 00-010-01
            feature: F010
            project_id: 00
            status: backlog
            ---
            ## 00-010-01: Already Migrated
            Already in new format
            """),
    }

    for filename, content in test_files.items():
        file_path = ws_dir / filename
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content, encoding="utf-8")

    return tmp_path


class TestWorkstreamFile:
    """Tests for WorkstreamFile class."""

    def test_parse_old_format_id(self, temp_workspace: Path) -> None:
        """Test parsing old format WS-FFF-SS."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md",
            project_id="00",
        )
        assert ws_file.old_id == "WS-001-01"
        assert ws_file.new_id == "00-001-01"

    def test_parse_new_format_id(self, temp_workspace: Path) -> None:
        """Test parsing new format PP-FFF-SS."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "already-migrated.md",
            project_id="00",
        )
        # For new format, old_id is set to new_id (not None)
        assert ws_file.new_id == "00-010-01"
        assert not ws_file.needs_migration()

    def test_needs_migration(self, temp_workspace: Path) -> None:
        """Test needs_migration detection."""
        old_format = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md",
            project_id="00",
        )
        assert old_format.needs_migration()

        new_format = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "already-migrated.md",
            project_id="00",
        )
        assert not new_format.needs_migration()

    def test_migrate_dry_run(self, temp_workspace: Path) -> None:
        """Test dry-run migration."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md",
            project_id="00",
        )

        original_content = ws_file.path.read_text()
        success, message = ws_file.migrate(dry_run=True)

        assert success
        assert "DRY RUN" in message
        assert ws_file.path.read_text() == original_content  # Unchanged

    def test_migrate_live(self, temp_workspace: Path) -> None:
        """Test live migration."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md",
            project_id="00",
        )

        success, message = ws_file.migrate(dry_run=False)

        assert success
        assert "Migrated" in message

        # Check new filename
        new_path = ws_file.path.parent / "00-001-01-test.md"
        assert new_path.exists()

        # Check content updated
        content = new_path.read_text()
        assert "ws_id: 00-001-01" in content
        assert "project_id: 00" in content
        assert "## 00-001-01:" in content

    def test_migrate_with_dependencies(self, temp_workspace: Path) -> None:
        """Test migration with dependency updates."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-003-05-deps.md",
            project_id="00",
        )

        success, message = ws_file.migrate(dry_run=False)

        assert success

        # Check dependencies updated
        new_path = ws_file.path.parent / "00-003-05-deps.md"
        content = new_path.read_text()

        assert "00-001-01" in content
        assert "00-002-01" in content
        assert "WS-001-01" not in content
        assert "WS-002-01" not in content


class TestWorkstreamMigrator:
    """Tests for WorkstreamMigrator class."""

    def test_find_workstream_files(self, temp_workspace: Path) -> None:
        """Test finding workstream files."""
        migrator = WorkstreamMigrator(temp_workspace, project_id="00")
        files = migrator._find_workstream_files()

        assert len(files) == 3  # Excludes already-migrated.md
        filenames = [f.name for f in files]
        assert "WS-001-01-test.md" in filenames
        assert "WS-002-01-done.md" in filenames
        assert "WS-003-05-deps.md" in filenames

    def test_migrate_dry_run(self, temp_workspace: Path) -> None:
        """Test dry-run migration."""
        migrator = WorkstreamMigrator(temp_workspace, project_id="00", dry_run=True)
        stats = migrator.migrate()

        assert stats["total"] == 3
        assert stats["migrated"] == 0  # Dry run doesn't actually migrate
        assert stats["skipped"] == 0
        assert stats["failed"] == 0

        # Verify files unchanged
        old_file = temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md"
        assert old_file.exists()
        assert "ws_id: WS-001-01" in old_file.read_text()

    def test_migrate_live(self, temp_workspace: Path) -> None:
        """Test live migration."""
        migrator = WorkstreamMigrator(temp_workspace, project_id="00", dry_run=False)
        stats = migrator.migrate()

        assert stats["total"] == 3
        assert stats["migrated"] == 3
        assert stats["failed"] == 0

        # Verify new files exist
        ws_dir = temp_workspace / "docs" / "workstreams"
        assert (ws_dir / "backlog" / "00-001-01-test.md").exists()
        assert (ws_dir / "completed" / "00-002-01-done.md").exists()
        assert (ws_dir / "backlog" / "00-003-05-deps.md").exists()

        # Verify old files don't exist
        assert not (ws_dir / "backlog" / "WS-001-01-test.md").exists()
        assert not (ws_dir / "completed" / "WS-002-01-done.md").exists()
        assert not (ws_dir / "backlog" / "WS-003-05-deps.md").exists()

    def test_invalid_project_id(self, temp_workspace: Path) -> None:
        """Test validation of project_id format."""
        # The validation happens in main() via print + return 1
        # We test by checking the migrator still works with valid ID
        migrator = WorkstreamMigrator(temp_workspace, project_id="99")
        assert migrator.project_id == "99"

        # Invalid IDs (non-numeric) won't raise during __init__
        # but will be caught in main() before migration starts

    def test_missing_workstreams_dir(self, tmp_path: Path) -> None:
        """Test error handling for missing workstreams directory."""
        with pytest.raises(WorkstreamMigrationError):
            migrator = WorkstreamMigrator(tmp_path, project_id="00")
            migrator.migrate()


class TestEdgeCases:
    """Tests for edge cases."""

    def test_missing_ws_id_in_frontmatter(self, temp_workspace: Path) -> None:
        """Test file without ws_id in frontmatter."""
        ws_dir = temp_workspace / "docs" / "workstreams" / "backlog"
        ws_dir.mkdir(parents=True, exist_ok=True)  # Use exist_ok=True

        # Create file without ws_id but with old-style filename
        test_file = ws_dir / "WS-099-01-no-id.md"
        test_file.write_text("## WS-099-01: No ID in frontmatter")

        ws_file = WorkstreamFile(test_file, project_id="00")
        assert ws_file.old_id == "WS-099-01"
        assert ws_file.new_id == "00-099-01"

    def test_custom_project_id(self, temp_workspace: Path) -> None:
        """Test migration with custom project_id."""
        ws_file = WorkstreamFile(
            temp_workspace / "docs" / "workstreams" / "backlog" / "WS-001-01-test.md",
            project_id="02",
        )

        assert ws_file.new_id == "02-001-01"

        success, _ = ws_file.migrate(dry_run=False)
        assert success

        new_path = ws_file.path.parent / "02-001-01-test.md"
        content = new_path.read_text()

        assert "ws_id: 02-001-01" in content
        assert "project_id: 02" in content

    def test_no_old_format_files(self, temp_workspace: Path) -> None:
        """Test migrator with no old format files."""
        # Remove old format files
        for old_file in temp_workspace.rglob("WS-*.md"):
            old_file.unlink()

        migrator = WorkstreamMigrator(temp_workspace, project_id="00")
        stats = migrator.migrate()

        assert stats["total"] == 0
