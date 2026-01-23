"""Update workstream file frontmatter."""

import re
from pathlib import Path
from typing import Optional


class FrontmatterUpdater:
    """Update YAML frontmatter in workstream markdown files.

    Handles safe updates to WS file frontmatter, specifically:
    - github_issue: null → github_issue: N
    - status: backlog → status: active
    """

    @staticmethod
    def update_github_issue(ws_file: Path, issue_number: int) -> None:
        """Update github_issue field in frontmatter.

        Replaces `github_issue: null` or existing number with new issue number.

        Args:
            ws_file: Path to WS markdown file
            issue_number: GitHub issue number to set

        Raises:
            FileNotFoundError: If file does not exist
            IOError: If file cannot be read/written
        """
        content = ws_file.read_text(encoding="utf-8")

        # Replace github_issue: null with github_issue: N
        # Also handles existing issue numbers (e.g., github_issue: 123)
        updated = re.sub(
            r"github_issue:\s*(null|\d+)",
            f"github_issue: {issue_number}",
            content,
        )

        ws_file.write_text(updated, encoding="utf-8")

    @staticmethod
    def update_status(ws_file: Path, new_status: str) -> None:
        """Update status field in frontmatter.

        Changes WS status from current value to new_status (e.g., backlog → active).

        Args:
            ws_file: Path to WS markdown file
            new_status: New status (backlog, active, completed, blocked)

        Raises:
            FileNotFoundError: If file does not exist
            IOError: If file cannot be read/written
        """
        content = ws_file.read_text(encoding="utf-8")

        # Replace status: xxx with status: new_status
        updated = re.sub(
            r"status:\s*\w+",
            f"status: {new_status}",
            content,
        )

        ws_file.write_text(updated, encoding="utf-8")

    @staticmethod
    def get_github_issue(ws_file: Path) -> Optional[int]:
        """Extract github_issue number from frontmatter.

        Args:
            ws_file: Path to WS markdown file

        Returns:
            Issue number if found, None if not set or invalid

        Raises:
            FileNotFoundError: If file does not exist
            IOError: If file cannot be read
        """
        content = ws_file.read_text(encoding="utf-8")

        match = re.search(r"github_issue:\s*(\d+)", content)
        if match:
            return int(match.group(1))
        return None
