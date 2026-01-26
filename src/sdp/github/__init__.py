"""GitHub API integration for SDP workstream automation.

This module provides:
- GitHubConfig: Configuration loading from .env
- GitHubClient: API client wrapper around PyGithub
- IssueSync: Workstream to GitHub issue synchronization
- LabelManager: Label creation and management
- MilestoneManager: Feature milestone management
- ConflictResolver: Conflict detection for sync
- EnhancedSyncService: Sync with conflict resolution
"""

from sdp.github.client import GitHubClient
from sdp.github.config import GitHubConfig
from sdp.github.conflict_resolver import Conflict, ConflictResolver
from sdp.github.sync_enhanced import EnhancedSyncService

__all__ = [
    "GitHubConfig",
    "GitHubClient",
    "Conflict",
    "ConflictResolver",
    "EnhancedSyncService",
]
