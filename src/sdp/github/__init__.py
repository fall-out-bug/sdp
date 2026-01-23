"""GitHub API integration for SDP workstream automation.

This module provides:
- GitHubConfig: Configuration loading from .env
- GitHubClient: API client wrapper around PyGithub
- IssueSync: Workstream to GitHub issue synchronization
- LabelManager: Label creation and management
- MilestoneManager: Feature milestone management
"""

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig

__all__ = [
    "GitHubConfig",
    "GitHubClient",
]
