"""
SDP + Beads Integration.

This module provides integration between SDP (Spec-Driven Protocol) and Beads
(git-backed issue tracker for AI agents).

Key components:
- BeadsClient: Interface for interacting with Beads
- MockBeadsClient: In-memory mock for development/testing
- CLIBeadsClient: Real Beads via CLI subprocess
- BeadsSyncService: Bidirectional sync between SDP workstreams and Beads tasks
- FeatureDecomposer: Decompose features into workstreams

Usage:
    from sdp.beads import create_beads_client, BeadsSyncService, FeatureDecomposer

    # Create client (mock for development, real for production)
    client = create_beads_client(use_mock=True)  # Mock
    client = create_beads_client()  # Real (requires Beads installed)

    # Create sync service
    sync = BeadsSyncService(client)

    # Decompose feature
    decomposer = FeatureDecomposer(client)
    ws_ids = decomposer.decompose("bd-0001")
"""

from .client import (
    BeadsClient,
    BeadsClientError,
    CLIBeadsClient,
    MockBeadsClient,
    create_beads_client,
)

# Re-export OneshotResult from execution_mode (where it's now defined)
from .execution_mode import (
    AuditLogEntry,
    AuditLogger,
    DestructiveOperationDetector,
    DestructiveOperations,
    ExecutionMode,
    OneshotResult,
)
from .idea_interview import (
    AmbiguityDetector,
    CriticalQuestions,
    IdeaInterviewer,
    InterviewResult,
    InterviewRound,
)
from .models import (
    BeadsDependency,
    BeadsDependencyType,
    BeadsPriority,
    BeadsStatus,
    BeadsSyncResult,
    BeadsTask,
    BeadsTaskCreate,
)
from .skills_build import ExecutionResult, WorkstreamExecutor
from .skills_design import FeatureDecomposer, WorkstreamSpec
from .skills_oneshot import MultiAgentExecutor
from .sync import BeadsSyncError, BeadsSyncService

__all__ = [
    # Client
    "BeadsClient",
    "MockBeadsClient",
    "CLIBeadsClient",
    "create_beads_client",
    "BeadsClientError",
    # Models
    "BeadsTask",
    "BeadsTaskCreate",
    "BeadsStatus",
    "BeadsPriority",
    "BeadsDependency",
    "BeadsDependencyType",
    "BeadsSyncResult",
    # Services
    "BeadsSyncService",
    "BeadsSyncError",
    "FeatureDecomposer",
    "WorkstreamSpec",
    "WorkstreamExecutor",
    "ExecutionResult",
    "MultiAgentExecutor",
    # Idea interview (F014)
    "InterviewRound",
    "CriticalQuestions",
    "AmbiguityDetector",
    "IdeaInterviewer",
    "InterviewResult",
    # Execution modes (F014)
    "ExecutionMode",
    "AuditLogger",
    "DestructiveOperationDetector",
    "DestructiveOperations",
    "AuditLogEntry",
    "OneshotResult",
]

__version__ = "0.4.0"
