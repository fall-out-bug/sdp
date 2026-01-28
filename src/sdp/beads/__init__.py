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
    MockBeadsClient,
    CLIBeadsClient,
    create_beads_client,
    BeadsClientError,
)
from .models import (
    BeadsTask,
    BeadsTaskCreate,
    BeadsStatus,
    BeadsPriority,
    BeadsDependency,
    BeadsDependencyType,
    BeadsSyncResult,
)
from .sync import BeadsSyncService, BeadsSyncError
from .skills_design import FeatureDecomposer, WorkstreamSpec
from .skills_build import WorkstreamExecutor, ExecutionResult
from .skills_oneshot import MultiAgentExecutor
from .idea_interview import (
    InterviewRound,
    CriticalQuestions,
    AmbiguityDetector,
    IdeaInterviewer,
    InterviewResult,
)
from .execution_mode import (
    ExecutionMode,
    AuditLogger,
    DestructiveOperationDetector,
    DestructiveOperations,
    AuditLogEntry,
)
# Re-export OneshotResult from execution_mode (where it's now defined)
from .execution_mode import OneshotResult

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
