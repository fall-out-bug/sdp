"""
Multi-agent oneshot execution using Beads ready detection.

DEPRECATED: This module is split into smaller modules for better maintainability.
Import from sdp.beads.oneshot instead:
- executor: MultiAgentExecutor

This module remains for backward compatibility.

Implements @oneshot skill logic: execute all workstreams for a feature
using multiple agents in parallel with Beads dependency tracking.

Enhanced with workflow efficiency modes (F014):
- Standard mode (PR required)
- Auto-approve mode (skip PR)
- Sandbox mode (skip PR, sandbox only)
- Dry-run mode (preview changes)
"""

# Re-export all public APIs for backward compatibility
from .oneshot.executor import MultiAgentExecutor  # noqa: F401

__all__ = ["MultiAgentExecutor"]
