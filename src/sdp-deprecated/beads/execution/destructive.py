"""
Destructive operations detection.

Identifies potentially dangerous operations that require user confirmation.
"""

from dataclasses import dataclass
from typing import List, Optional


@dataclass
class DestructiveOperations:
    """Result of destructive operations check."""

    has_destructive_operations: bool
    operation_types: List[str]
    files_affected: List[str]
    details: Optional[str] = None


class DestructiveOperationDetector:
    """Detect destructive operations that require user confirmation."""

    # Patterns that indicate destructive operations
    DESTRUCTIVE_PATTERNS = {
        "database_migration": ["migration", "migrate", "schema", "alembic"],
        "file_deletion": ["delete", "remove", "rm"],
        "data_loss": ["drop", "truncate", "wipe"],
    }

    def check_operations(
        self,
        files_to_create: List[str],
        files_to_modify: List[str],
        files_to_delete: List[str],
    ) -> DestructiveOperations:
        """Check if operations are destructive.

        Args:
            files_to_create: List of file paths to be created
            files_to_modify: List of file paths to be modified
            files_to_delete: List of file paths to be deleted

        Returns:
            DestructiveOperations with check result
        """
        operation_types = []
        files_affected = []

        # Check file deletions (always destructive)
        if files_to_delete:
            operation_types.append("file_deletion")
            files_affected.extend(files_to_delete)

        # Check for database migrations in created files
        for file_path in files_to_create:
            if self._is_database_migration(file_path):
                operation_types.append("database_migration")
                files_affected.append(file_path)

        # Check for destructive patterns in modified files
        for file_path in files_to_modify:
            if self._has_destructive_changes(file_path):
                operation_types.append("data_loss")
                files_affected.append(file_path)

        return DestructiveOperations(
            has_destructive_operations=len(operation_types) > 0,
            operation_types=list(set(operation_types)),  # Deduplicate
            files_affected=list(set(files_affected)),  # Deduplicate
            details=f"Found {len(operation_types)} destructive operation types",
        )

    def _is_database_migration(self, file_path: str) -> bool:
        """Check if file is a database migration."""
        path_lower = file_path.lower()
        return any(
            pattern in path_lower
            for patterns in self.DESTRUCTIVE_PATTERNS.values()
            for pattern in patterns
            if pattern in path_lower
        )

    def _has_destructive_changes(self, file_path: str) -> bool:
        """Check if file modification involves destructive changes."""
        # Simplified check - in real implementation would read file content
        path_lower = file_path.lower()
        return any(
            pattern in path_lower
            for pattern in self.DESTRUCTIVE_PATTERNS.get("data_loss", [])
        )
