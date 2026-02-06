"""Agent spawner for multi-agent execution.

Provides AgentSpawner class for spawning autonomous agents via
Claude Code's Task tool with support for background execution.
"""

import logging
from dataclasses import dataclass
from typing import Optional

logger = logging.getLogger(__name__)


@dataclass
class AgentConfig:
    """Configuration for spawning an agent.

    Attributes:
        name: Agent name/identifier
        prompt: Task prompt for the agent
        subagent_type: Optional subagent type (e.g., "builder", "general-purpose")
        model: Optional model to use (e.g., "haiku", "sonnet", "opus")
        description: Optional description for the agent
        run_in_background: Whether to run agent in background mode
    """

    name: str
    prompt: str
    subagent_type: Optional[str] = None
    model: Optional[str] = None
    description: Optional[str] = None
    run_in_background: bool = False


class AgentSpawner:
    """Spawns and manages autonomous agents via Task tool.

    Provides interface for spawning agents with Task tool,
    tracking active agents, and managing agent lifecycle.
    """

    def __init__(self) -> None:
        """Initialize agent spawner."""
        self._active_agents: dict[str, AgentConfig] = {}

    def spawn_agent(self, config: AgentConfig) -> Optional[str]:
        """Spawn an agent via Task tool.

        Args:
            config: Agent configuration

        Returns:
            Agent ID if spawned successfully, None if failed

        Raises:
            Exception: If Task tool invocation fails
        """
        from sdp.unified.agent import task_wrapper

        try:
            # Call Task tool with arguments
            result = task_wrapper.Task(
                description=config.name,
                prompt=config.prompt,
                subagent_type=config.subagent_type,
                model=config.model,
                run_in_background=config.run_in_background,
            )

            # Extract agent ID from result
            agent_id = getattr(result, "task_id", None)

            if agent_id:
                # Track the spawned agent
                self._active_agents[agent_id] = config
                logger.info(f"Spawned agent: {agent_id}")
                return str(agent_id) if agent_id else None

            return None

        except Exception as e:
            logger.error(f"Failed to spawn agent: {e}")
            raise

    def get_active_agents(self) -> list[str]:
        """Get list of active agent IDs.

        Returns:
            List of active agent IDs
        """
        return list(self._active_agents.keys())

    def mark_agent_complete(self, agent_id: str) -> None:
        """Mark an agent as complete and remove from tracking.

        Args:
            agent_id: Agent ID to mark complete
        """
        if agent_id in self._active_agents:
            del self._active_agents[agent_id]
            logger.info(f"Agent completed: {agent_id}")
