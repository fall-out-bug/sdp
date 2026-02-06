"""GitHub integration exceptions."""

import datetime


class GitHubSyncError(Exception):
    """Base exception for GitHub sync errors."""

    def __init__(self, message: str, action: str | None = None) -> None:
        """Initialize exception.

        Args:
            message: Error message
            action: User action suggestion

        """
        super().__init__(message)
        self.action = action


class RateLimitError(GitHubSyncError):
    """GitHub API rate limit exceeded."""

    def __init__(self, reset_time: int) -> None:
        """Initialize rate limit error.

        Args:
            reset_time: Unix timestamp when limit resets

        """
        reset_dt = datetime.datetime.fromtimestamp(reset_time)
        message = f"GitHub API rate limit exceeded. Resets at {reset_dt}"
        action = "Wait for rate limit reset or use different token"
        super().__init__(message, action)
        self.reset_time = reset_time


class AuthenticationError(GitHubSyncError):
    """GitHub authentication failed."""

    def __init__(self) -> None:
        """Initialize authentication error."""
        message = "GitHub authentication failed. Check GITHUB_TOKEN"
        action = (
            "Verify GITHUB_TOKEN in .env has correct permissions "
            "(repo, project)"
        )
        super().__init__(message, action)


class ProjectNotFoundError(GitHubSyncError):
    """GitHub project not found."""

    def __init__(self, project_name: str) -> None:
        """Initialize project not found error.

        Args:
            project_name: Name of the missing project

        """
        message = f"GitHub project '{project_name}' not found"
        action = (
            f"Create project '{project_name}' in GitHub or run: "
            "hwc github setup"
        )
        super().__init__(message, action)
