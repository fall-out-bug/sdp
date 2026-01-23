"""Map workstream status to GitHub issue status."""

from typing import Optional


class StatusMapper:
    """Bidirectional mapping between WS and GitHub status.

    Maps WS status (backlog, active, completed, blocked) to GitHub
    issue labels and states, enabling bidirectional synchronization.

    Attributes:
        WS_TO_GITHUB: Mapping of WS status to (label, state) tuples
        GITHUB_TO_WS: Mapping of GitHub status label to WS status
    """

    # WS status → GitHub (label, state)
    WS_TO_GITHUB: dict[str, tuple[str, str]] = {
        "backlog": ("status/backlog", "open"),
        "active": ("status/in-progress", "open"),
        "completed": ("status/completed", "closed"),
        "blocked": ("status/blocked", "open"),
    }

    # GitHub label → WS status
    GITHUB_TO_WS: dict[str, str] = {
        "status/backlog": "backlog",
        "status/in-progress": "active",
        "status/completed": "completed",
        "status/blocked": "blocked",
    }

    @classmethod
    def ws_to_github_label(cls, ws_status: str) -> str:
        """Convert WS status to GitHub label.

        Args:
            ws_status: WS status (backlog, active, completed, blocked)

        Returns:
            GitHub label (e.g., "status/in-progress")

        Raises:
            ValueError: If status unknown
        """
        if ws_status not in cls.WS_TO_GITHUB:
            raise ValueError(f"Unknown WS status: {ws_status}")
        return cls.WS_TO_GITHUB[ws_status][0]

    @classmethod
    def ws_to_github_state(cls, ws_status: str) -> str:
        """Convert WS status to GitHub issue state.

        Args:
            ws_status: WS status (backlog, active, completed, blocked)

        Returns:
            "open" or "closed"

        Raises:
            ValueError: If status unknown
        """
        if ws_status not in cls.WS_TO_GITHUB:
            raise ValueError(f"Unknown WS status: {ws_status}")
        return cls.WS_TO_GITHUB[ws_status][1]

    @classmethod
    def github_label_to_ws(cls, label: str) -> Optional[str]:
        """Convert GitHub status label to WS status.

        Args:
            label: GitHub label (e.g., "status/in-progress")

        Returns:
            WS status or None if not a status label
        """
        return cls.GITHUB_TO_WS.get(label)

    @classmethod
    def detect_conflict(
        cls, ws_status: str, github_label: str, github_state: str
    ) -> bool:
        """Detect if WS and GitHub status are inconsistent.

        WS file is source of truth. Conflict occurs when GitHub status
        differs from what WS status should map to.

        Args:
            ws_status: Status from WS file (backlog, active, etc)
            github_label: Status label from GitHub issue
            github_state: Issue state (open/closed)

        Returns:
            True if conflict detected (GitHub differs from expected)

        Raises:
            ValueError: If ws_status unknown
        """
        expected_label = cls.ws_to_github_label(ws_status)
        expected_state = cls.ws_to_github_state(ws_status)

        return github_label != expected_label or github_state != expected_state
