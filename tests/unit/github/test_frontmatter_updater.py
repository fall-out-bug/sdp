"""Tests for frontmatter updater."""

from pathlib import Path

import pytest

from sdp.github.frontmatter_updater import FrontmatterUpdater


class TestFrontmatterUpdaterUpdateGitHubIssue:
    """Test update_github_issue method."""

    def test_update_github_issue_replaces_null(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue replaces null with issue number."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "github_issue: null\n"
            "---\n"
            "# Content"
        )

        FrontmatterUpdater.update_github_issue(ws_file, 123)

        content = ws_file.read_text()
        assert "github_issue: 123" in content
        assert "github_issue: null" not in content

    def test_update_github_issue_replaces_existing_number(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue replaces existing issue number."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "github_issue: 456\n"
            "---\n"
            "# Content"
        )

        FrontmatterUpdater.update_github_issue(ws_file, 789)

        content = ws_file.read_text()
        assert "github_issue: 789" in content
        assert "github_issue: 456" not in content

    def test_update_github_issue_preserves_other_fields(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue preserves other frontmatter fields."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "feature: F001\n"
            "github_issue: null\n"
            "status: backlog\n"
            "---\n"
            "# Content"
        )

        FrontmatterUpdater.update_github_issue(ws_file, 999)

        content = ws_file.read_text()
        assert "github_issue: 999" in content
        assert "ws_id: 00-001-01" in content
        assert "feature: F001" in content
        assert "status: backlog" in content

    def test_update_github_issue_handles_whitespace(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue handles various whitespace patterns."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "github_issue:   null  \n"
            "---\n"
        )

        FrontmatterUpdater.update_github_issue(ws_file, 111)

        content = ws_file.read_text()
        assert "github_issue: 111" in content

    def test_update_github_issue_raises_file_not_found(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue raises FileNotFoundError for missing file."""
        ws_file = tmp_path / "nonexistent.md"

        with pytest.raises(FileNotFoundError):
            FrontmatterUpdater.update_github_issue(ws_file, 123)

    def test_update_github_issue_handles_multiple_occurrences(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue updates all occurrences."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "github_issue: null\n"
            "---\n"
            "# Content\n"
            "github_issue: null\n"
        )

        FrontmatterUpdater.update_github_issue(ws_file, 222)

        content = ws_file.read_text()
        # Should update both occurrences
        assert content.count("github_issue: 222") == 2


class TestFrontmatterUpdaterUpdateStatus:
    """Test update_status method."""

    def test_update_status_changes_status(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status changes status field."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "status: backlog\n"
            "---\n"
            "# Content"
        )

        FrontmatterUpdater.update_status(ws_file, "active")

        content = ws_file.read_text()
        assert "status: active" in content
        assert "status: backlog" not in content

    def test_update_status_preserves_other_fields(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status preserves other frontmatter fields."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "feature: F001\n"
            "status: backlog\n"
            "github_issue: 123\n"
            "---\n"
            "# Content"
        )

        FrontmatterUpdater.update_status(ws_file, "completed")

        content = ws_file.read_text()
        assert "status: completed" in content
        assert "ws_id: 00-001-01" in content
        assert "feature: F001" in content
        assert "github_issue: 123" in content

    def test_update_status_handles_whitespace(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status handles whitespace variations."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "status:   backlog  \n"
            "---\n"
        )

        FrontmatterUpdater.update_status(ws_file, "active")

        content = ws_file.read_text()
        assert "status: active" in content

    def test_update_status_raises_file_not_found(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status raises FileNotFoundError for missing file."""
        ws_file = tmp_path / "nonexistent.md"

        with pytest.raises(FileNotFoundError):
            FrontmatterUpdater.update_status(ws_file, "active")

    def test_update_status_handles_all_status_values(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status handles all valid status values."""
        statuses = ["backlog", "active", "completed", "blocked"]

        for old_status in statuses:
            for new_status in statuses:
                if old_status != new_status:
                    ws_file = tmp_path / f"test-{old_status}-{new_status}.md"
                    ws_file.write_text(
                        f"---\n"
                        f"status: {old_status}\n"
                        f"---\n"
                    )

                    FrontmatterUpdater.update_status(ws_file, new_status)

                    content = ws_file.read_text()
                    assert f"status: {new_status}" in content
                    assert f"status: {old_status}" not in content


class TestFrontmatterUpdaterGetGitHubIssue:
    """Test get_github_issue method."""

    def test_get_github_issue_returns_number(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue returns issue number when present."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "github_issue: 456\n"
            "---\n"
            "# Content"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result == 456

    def test_get_github_issue_returns_none_when_null(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue returns None when set to null."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "github_issue: null\n"
            "---\n"
            "# Content"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result is None

    def test_get_github_issue_returns_none_when_missing(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue returns None when field missing."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "ws_id: 00-001-01\n"
            "---\n"
            "# Content"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result is None

    def test_get_github_issue_handles_whitespace(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue handles whitespace around number."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "github_issue:   789   \n"
            "---\n"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result == 789

    def test_get_github_issue_raises_file_not_found(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue raises FileNotFoundError for missing file."""
        ws_file = tmp_path / "nonexistent.md"

        with pytest.raises(FileNotFoundError):
            FrontmatterUpdater.get_github_issue(ws_file)

    def test_get_github_issue_handles_multiple_occurrences(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue returns first occurrence."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "github_issue: 111\n"
            "---\n"
            "# Content\n"
            "github_issue: 222\n"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result == 111

    def test_get_github_issue_handles_invalid_format(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue returns None for invalid format."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text(
            "---\n"
            "github_issue: abc\n"
            "---\n"
        )

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result is None


class TestFrontmatterUpdaterErrorHandling:
    """Test error handling."""

    def test_update_github_issue_handles_io_error(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue handles IO errors."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("---\n---\n")

        # Make file read-only to simulate IO error on write
        ws_file.chmod(0o444)

        try:
            with pytest.raises((PermissionError, IOError)):
                FrontmatterUpdater.update_github_issue(ws_file, 123)
        finally:
            # Restore permissions for cleanup
            ws_file.chmod(0o644)

    def test_update_status_handles_io_error(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status handles IO errors."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("---\nstatus: backlog\n---\n")

        # Make file read-only to simulate IO error on write
        ws_file.chmod(0o444)

        try:
            with pytest.raises((PermissionError, IOError)):
                FrontmatterUpdater.update_status(ws_file, "active")
        finally:
            # Restore permissions for cleanup
            ws_file.chmod(0o644)

    def test_get_github_issue_handles_encoding_error(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue handles encoding errors."""
        ws_file = tmp_path / "00-001-01.md"
        # Write binary data that can't be decoded as UTF-8
        ws_file.write_bytes(b"\xff\xfe\x00\x00")

        with pytest.raises(UnicodeDecodeError):
            FrontmatterUpdater.get_github_issue(ws_file)


class TestFrontmatterUpdaterEdgeCases:
    """Test edge cases."""

    def test_update_github_issue_with_empty_file(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue handles empty file."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("")

        # Should not raise, but may not update anything
        FrontmatterUpdater.update_github_issue(ws_file, 123)

        content = ws_file.read_text()
        # Empty file remains empty
        assert content == ""

    def test_update_status_with_empty_file(
        self, tmp_path: Path
    ) -> None:
        """Verify update_status handles empty file."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("")

        # Should not raise, but may not update anything
        FrontmatterUpdater.update_status(ws_file, "active")

        content = ws_file.read_text()
        # Empty file remains empty
        assert content == ""

    def test_get_github_issue_with_empty_file(
        self, tmp_path: Path
    ) -> None:
        """Verify get_github_issue handles empty file."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("")

        result = FrontmatterUpdater.get_github_issue(ws_file)

        assert result is None

    def test_update_github_issue_with_no_frontmatter(
        self, tmp_path: Path
    ) -> None:
        """Verify update_github_issue handles file without frontmatter."""
        ws_file = tmp_path / "00-001-01.md"
        ws_file.write_text("# Just content, no frontmatter")

        # Should not raise, but may not update anything
        FrontmatterUpdater.update_github_issue(ws_file, 123)

        content = ws_file.read_text()
        # File content unchanged
        assert "# Just content, no frontmatter" in content
