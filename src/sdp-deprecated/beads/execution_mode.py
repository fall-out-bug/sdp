"""
Execution modes and related classes for @oneshot workflow efficiency feature.

DEPRECATED: This module is split into smaller modules for better maintainability.
Import from sdp.beads.execution instead:
- models: ExecutionMode, OneshotResult
- destructive: DestructiveOperations, DestructiveOperationDetector
- audit: AuditLogEntry, AuditLogger

This module remains for backward compatibility.

Defines:
- ExecutionMode enum for different execution modes
- DestructiveOperationDetector for detecting dangerous operations
- AuditLogger for tracking --auto-approve executions
- OneshotResult dataclass for execution results
"""

# Re-export all public APIs for backward compatibility
from .execution import (  # noqa: F401
    AuditLogEntry,
    AuditLogger,
    DestructiveOperationDetector,
    DestructiveOperations,
    ExecutionMode,
    OneshotResult,
)

__all__ = [
    "ExecutionMode",
    "OneshotResult",
    "DestructiveOperations",
    "DestructiveOperationDetector",
    "AuditLogEntry",
    "AuditLogger",
]
