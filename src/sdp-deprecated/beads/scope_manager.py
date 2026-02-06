"""Scope management for workstream file restrictions.

Manages which files a workstream is allowed to modify by storing
scope information in Beads task metadata.
"""

from typing import List, cast

from sdp.beads.base import BeadsClient
from sdp.beads.exceptions import BeadsClientError


class ScopeManager:
    """Manage workstream file scope via Beads metadata.

    Scope determines which files a workstream can modify. Stored in
    task metadata under 'scope_files' key.
    """

    def __init__(self, client: BeadsClient):
        """Initialize scope manager.

        Args:
            client: Beads client for metadata operations
        """
        self._client = client

    def get_scope(self, ws_id: str) -> List[str]:
        """Get scope files for workstream.

        Args:
            ws_id: Workstream/task ID

        Returns:
            List of file paths in scope. Empty list means unrestricted.

        Raises:
            ValueError: If workstream not found
        """
        task = self._client.get_task(ws_id)
        if not task:
            raise ValueError(f"Workstream not found: {ws_id}")

        scope = task.sdp_metadata.get("scope_files", [])
        return cast(List[str], scope if isinstance(scope, list) else [])

    def set_scope(self, ws_id: str, files: List[str]) -> None:
        """Set scope files for workstream.

        Args:
            ws_id: Workstream/task ID
            files: List of file paths to allow

        Raises:
            ValueError: If workstream not found
        """
        task = self._client.get_task(ws_id)
        if not task:
            raise ValueError(f"Workstream not found: {ws_id}")

        # Update metadata
        metadata = task.sdp_metadata.copy()
        metadata["scope_files"] = files

        try:
            self._client.update_metadata(ws_id, metadata)
        except BeadsClientError as e:
            raise ValueError(f"Failed to update scope for {ws_id}: {e}") from e

    def add_file(self, ws_id: str, file_path: str) -> None:
        """Add file to workstream scope.

        Args:
            ws_id: Workstream/task ID
            file_path: File path to add

        Raises:
            ValueError: If workstream not found or update fails
        """
        scope = self.get_scope(ws_id)
        if file_path not in scope:
            scope.append(file_path)
            self.set_scope(ws_id, scope)

    def remove_file(self, ws_id: str, file_path: str) -> None:
        """Remove file from workstream scope.

        Args:
            ws_id: Workstream/task ID
            file_path: File path to remove

        Raises:
            ValueError: If workstream not found or update fails
        """
        scope = self.get_scope(ws_id)
        if file_path in scope:
            scope.remove(file_path)
            self.set_scope(ws_id, scope)

    def is_in_scope(self, ws_id: str, file_path: str) -> bool:
        """Check if file is in workstream scope.

        Args:
            ws_id: Workstream/task ID
            file_path: File path to check

        Returns:
            True if file is allowed (in scope or scope is empty/unrestricted)

        Raises:
            ValueError: If workstream not found
        """
        scope = self.get_scope(ws_id)

        # Empty scope = all files allowed (unrestricted)
        if not scope:
            return True

        return file_path in scope

    def clear_scope(self, ws_id: str) -> None:
        """Clear scope (make unrestricted).

        Args:
            ws_id: Workstream/task ID

        Raises:
            ValueError: If workstream not found or update fails
        """
        self.set_scope(ws_id, [])
