"""
Checkpoint database schema and management.

DEPRECATED: This module is split into smaller modules for better maintainability.
Import from the new modules instead:
- models: Checkpoint, CheckpointStatus
- storage: CheckpointDatabase
- serialization: row_to_checkpoint, checkpoint_to_insert_params, checkpoint_to_update_params
- schema_manager: SchemaManager

This module remains for backward compatibility.
"""

# Re-export all public APIs for backward compatibility
from .models import Checkpoint, CheckpointStatus  # noqa: F401
from .schema_manager import SchemaManager  # noqa: F401
from .serialization import (  # noqa: F401
    checkpoint_to_insert_params,
    checkpoint_to_update_params,
    row_to_checkpoint,
)
from .storage import CheckpointDatabase  # noqa: F401

__all__ = [
    "Checkpoint",
    "CheckpointStatus",
    "CheckpointDatabase",
    "SchemaManager",
    "row_to_checkpoint",
    "checkpoint_to_insert_params",
    "checkpoint_to_update_params",
]
