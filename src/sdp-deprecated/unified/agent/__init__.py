"""Agent runtime module for multi-agent execution.

Provides agent spawning, management, and coordination capabilities
using Claude Code's Task tool.
"""

from sdp.unified.agent.spawner import AgentConfig, AgentSpawner

__all__ = [
    "AgentConfig",
    "AgentSpawner",
]
