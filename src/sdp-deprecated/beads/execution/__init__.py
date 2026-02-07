"""
Execution modes and related classes for @oneshot workflow efficiency.

This module is split into smaller modules for better maintainability.
Import from sub-modules:
- models: ExecutionMode, OneshotResult
- destructive: DestructiveOperations, DestructiveOperationDetector
- audit: AuditLogEntry, AuditLogger

This module remains for backward compatibility.
"""

from .audit import AuditLogEntry, AuditLogger  # noqa: F401
from .destructive import DestructiveOperationDetector, DestructiveOperations  # noqa: F401
from .models import ExecutionMode, OneshotResult  # noqa: F401

__all__ = [
    "ExecutionMode",
    "OneshotResult",
    "DestructiveOperations",
    "DestructiveOperationDetector",
    "AuditLogEntry",
    "AuditLogger",
]
