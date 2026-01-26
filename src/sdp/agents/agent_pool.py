"""Agent pool for managing concurrent agent execution."""

from __future__ import annotations

import asyncio
import logging
from dataclasses import dataclass, field
from datetime import datetime
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import AsyncIterator

logger = logging.getLogger(__name__)


@dataclass
class AgentInfo:
    """Information about an agent in the pool."""

    agent_id: str
    busy: bool = False
    current_ws: str | None = None
    started_at: datetime | None = None
    completed_count: int = 0

    @property
    def elapsed_seconds(self) -> float:
        """Time elapsed since agent started current workstream."""
        if not self.started_at:
            return 0.0
        return (datetime.now() - self.started_at).total_seconds()


@dataclass
class PoolStats:
    """Statistics for the agent pool."""

    total_agents: int = 0
    busy_agents: int = 0
    available_agents: int = 0
    total_completed: int = 0
    average_execution_time: float = 0.0


class AgentPool:
    """Manages a pool of concurrent agent workers.

    Provides load balancing by assigning work to the least busy agent.
    Tracks agent statistics and state.
    """

    def __init__(self, max_agents: int = 3) -> None:
        """Initialize agent pool.

        Args:
            max_agents: Maximum number of concurrent agents
        """
        self._max_agents = max_agents
        self._agents: dict[str, AgentInfo] = {}
        self._semaphore = asyncio.Semaphore(max_agents)
        self._lock = asyncio.Lock()
        self._init_agents()

    def _init_agents(self) -> None:
        """Initialize agent instances."""
        for i in range(self._max_agents):
            agent_id = f"agent-{i + 1}"
            self._agents[agent_id] = AgentInfo(agent_id=agent_id)

    async def acquire(self, ws_id: str | None = None) -> str:
        """Acquire an available agent.

        Blocks until an agent is available.

        Args:
            ws_id: Optional workstream ID to assign

        Returns:
            Agent ID that was acquired
        """
        await self._semaphore.acquire()

        async with self._lock:
            # Get least busy agent (prefer idle agents)
            agent_id = self._get_least_busy_agent()
            agent = self._agents[agent_id]

            agent.busy = True
            agent.current_ws = ws_id
            agent.started_at = datetime.now()

            logger.debug(f"Acquired agent {agent_id} for {ws_id}")
            return agent_id

    async def release(self, agent_id: str, success: bool = True) -> None:
        """Release an agent back to the pool.

        Args:
            agent_id: Agent ID to release
            success: Whether the work was completed successfully
        """
        async with self._lock:
            if agent_id in self._agents:
                agent = self._agents[agent_id]
                agent.busy = False
                agent.current_ws = None
                agent.started_at = None
                if success:
                    agent.completed_count += 1

        self._semaphore.release()
        logger.debug(f"Released agent {agent_id}")

    def get_stats(self) -> PoolStats:
        """Get current pool statistics.

        Returns:
            PoolStats with current pool state
        """
        busy_count = sum(1 for a in self._agents.values() if a.busy)
        total_completed = sum(a.completed_count for a in self._agents.values())

        # Calculate average execution time (simplified)
        avg_time = 0.0
        busy_agents = [a for a in self._agents.values() if a.busy]
        if busy_agents:
            avg_time = sum(a.elapsed_seconds for a in busy_agents) / len(busy_agents)

        return PoolStats(
            total_agents=self._max_agents,
            busy_agents=busy_count,
            available_agents=self._max_agents - busy_count,
            total_completed=total_completed,
            average_execution_time=avg_time,
        )

    def get_agent_info(self, agent_id: str) -> AgentInfo | None:
        """Get information about a specific agent.

        Args:
            agent_id: Agent ID to query

        Returns:
            AgentInfo if agent exists, None otherwise
        """
        return self._agents.get(agent_id)

    def list_agents(self) -> list[AgentInfo]:
        """List all agents in the pool.

        Returns:
            List of AgentInfo for all agents
        """
        return list(self._agents.values())

    async def execute(
        self,
        ws_id: str,
        coro,
    ) -> tuple[bool, str | None]:
        """Execute a coroutine with an agent from the pool.

        Acquires an agent, runs the coroutine, releases the agent.

        Args:
            ws_id: Workstream ID being executed
            coro: Async coroutine to execute

        Returns:
            Tuple of (success, error_message)
        """
        agent_id = await self.acquire(ws_id)
        error = None
        success = False

        try:
            result = await coro
            # Assume coroutine returns (success, error) or raises
            if isinstance(result, tuple):
                success, error = result
            else:
                success = bool(result)
        except Exception as e:
            error = str(e)
            success = False
            logger.exception(f"Agent {agent_id} failed on {ws_id}")
        finally:
            await self.release(agent_id, success)

        return success, error

    def _get_least_busy_agent(self) -> str:
        """Get the least busy agent.

        Prefer idle agents, then agents with fewer completed tasks.

        Returns:
            Agent ID
        """
        # Sort by: busy (idle first), then completed count
        sorted_agents = sorted(
            self._agents.values(),
            key=lambda a: (a.busy, a.completed_count),
        )
        return sorted_agents[0].agent_id

    async def __aenter__(self) -> AgentPool:
        """Context manager entry."""
        return self

    async def __aexit__(self, *args) -> None:
        """Context manager exit."""
        # Wait for all agents to finish
        for _ in range(self._max_agents):
            await self._semaphore.acquire()


class AgentPoolIterator:
    """Iterator for executing workstreams with automatic agent management."""

    def __init__(
        self,
        pool: AgentPool,
        ws_ids: list[str],
        execute_fn,
    ) -> None:
        """Initialize iterator.

        Args:
            pool: Agent pool to use
            ws_ids: Workstream IDs to execute
            execute_fn: Async function taking (agent_id, ws_id) and returning (success, error)
        """
        self._pool = pool
        self._ws_ids = ws_ids
        self._execute_fn = execute_fn
        self._results: dict[str, tuple[bool, str | None]] = {}
        self._tasks: dict[str, asyncio.Task] = {}

    async def run_all(self) -> dict[str, tuple[bool, str | None]]:
        """Run all workstreams concurrently.

        Returns:
            Dict mapping ws_id to (success, error) tuples
        """
        # Create tasks for all workstreams
        for ws_id in self._ws_ids:
            self._tasks[ws_id] = asyncio.create_task(self._run_one(ws_id))

        # Wait for all to complete
        await asyncio.gather(*self._tasks.values(), return_exceptions=True)

        # Collect results
        for ws_id, task in self._tasks.items():
            try:
                self._results[ws_id] = task.result()
            except Exception as e:
                self._results[ws_id] = (False, str(e))

        return self._results

    async def _run_one(self, ws_id: str) -> tuple[bool, str | None]:
        """Run a single workstream.

        Args:
            ws_id: Workstream ID

        Returns:
            Tuple of (success, error)
        """
        agent_id = await self._pool.acquire(ws_id)
        try:
            return await self._execute_fn(agent_id, ws_id)
        finally:
            await self._pool.release(agent_id, True)

    async def iter_results(self) -> AsyncIterator[tuple[str, bool, str | None]]:
        """Yield results as they complete.

        Yields:
            Tuples of (ws_id, success, error)
        """
        # Create tasks and track them
        pending_tasks: dict[asyncio.Task, str] = {
            asyncio.create_task(self._run_one(ws_id)): ws_id
            for ws_id in self._ws_ids
        }

        while pending_tasks:
            # Wait for next to complete
            done, _pending_tasks_set = await asyncio.wait(
                pending_tasks.keys(),
                return_when=asyncio.FIRST_COMPLETED,
            )

            for task in done:
                ws_id = pending_tasks[task]
                # Remove from pending
                del pending_tasks[task]
                try:
                    success, error = task.result()
                    yield ws_id, success, error
                except Exception as e:
                    yield ws_id, False, str(e)
