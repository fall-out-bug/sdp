"""
Multi-agent oneshot execution using Beads ready detection.

Implements @oneshot skill logic: execute all workstreams for a feature
using multiple agents in parallel with Beads dependency tracking.
"""

from .executor import MultiAgentExecutor

__all__ = ["MultiAgentExecutor"]
