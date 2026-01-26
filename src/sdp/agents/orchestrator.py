"""Multi-agent orchestrator for workstream execution."""

from __future__ import annotations

import asyncio
import json
import logging
from dataclasses import dataclass, field
from datetime import datetime
from pathlib import Path
from typing import TYPE_CHECKING, TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Callable

from .agent_pool import AgentPool, AgentPoolIterator, PoolStats
from .dependency_graph import DependencyGraph
from ..queue.task_queue import Priority, Task, TaskQueue

logger = logging.getLogger(__name__)


@dataclass
class OrchestratorState:
    """Persistent state for the orchestrator."""

    feature_id: str
    started_at: datetime
    status: str = "running"  # running, completed, failed, cancelled
    results: dict[str, bool] = field(default_factory=dict)
    current_ws: str | None = None
    errors: dict[str, str] = field(default_factory=dict)


@dataclass
class ExecutionResult:
    """Result of executing a feature."""

    feature_id: str
    success: bool
    completed: list[str]
    failed: list[str]
    duration_seconds: float
    errors: dict[str, str]


class Orchestrator:
    """Manages multi-agent workstream execution with dependency resolution.

    Features:
    - Dependency-aware execution order (topological sort)
    - Concurrent agent pool (configurable size)
    - State persistence for crash recovery
    - Progress tracking and callbacks
    - Deadlock detection (circular dependencies)
    """

    def __init__(
        self,
        max_agents: int = 3,
        ws_dir: str | Path = "docs/workstreams",
        state_file: str | Path = ".sdp/orchestrator_state.json",
    ) -> None:
        """Initialize orchestrator.

        Args:
            max_agents: Maximum concurrent agents
            ws_dir: Path to workstreams directory
            state_file: Path to state persistence file
        """
        self._max_agents = max_agents
        self._ws_dir = Path(ws_dir)
        self._state_file = Path(state_file)
        self._state_file.parent.mkdir(parents=True, exist_ok=True)

        self._pool = AgentPool(max_agents)
        self._graph = DependencyGraph(ws_dir)
        # TaskQueue needs a state file path, not a directory
        self._queue = TaskQueue(str(self._ws_dir.parent / ".sdp" / "queue_state.json"))

        self._state: OrchestratorState | None = None
        self._callbacks: list[Callable[[str, bool, str | None], None]] = []

    def on_progress(
        self,
        callback: Callable[[str, bool, str | None], None],
    ) -> None:
        """Register a progress callback.

        Args:
            callback: Function called with (ws_id, success, error) on completion
        """
        self._callbacks.append(callback)

    async def run_feature(
        self,
        feature_id: str,
        execute_fn: Callable[[str, str], tuple[bool, str | None]] | None = None,
    ) -> ExecutionResult:
        """Execute all workstreams for a feature.

        Args:
            feature_id: Feature ID (e.g., "F012")
            execute_fn: Optional async function taking (agent_id, ws_id)
                        and returning (success, error). If None, uses default.

        Returns:
            ExecutionResult with summary
        """
        start_time = datetime.now()

        # Initialize state
        self._state = OrchestratorState(feature_id=feature_id, started_at=start_time)
        self._save_state()

        # Get feature workstreams
        ws_ids = self._get_feature_workstreams(feature_id)
        if not ws_ids:
            logger.warning(f"No workstreams found for feature {feature_id}")
            return self._make_result(feature_id, start_time, {}, {})

        # Build dependency graph
        self._graph.build(ws_ids)

        # Check for circular dependencies
        if self._graph.has_circular_dependency():
            cycle = self._find_cycle()
            raise ValueError(f"Circular dependency detected: {' -> '.join(cycle)}")

        # Execute in dependency order
        results: dict[str, tuple[bool, str | None]] = {}

        if execute_fn is None:
            execute_fn = self._default_execute

        # Get execution order
        execution_order = self._graph.execution_order()
        logger.info(f"Execution order: {' -> '.join(execution_order)}")

        # Execute with level-based parallelization
        # (all workstreams with same dependency level can run in parallel)
        completed: set[str] = set()
        levels = self._group_by_level(execution_order)

        for level, level_ws_ids in levels:
            logger.info(f"Executing level {level}: {level_ws_ids}")

            # Execute level in parallel
            level_results = await self._execute_level(
                level_ws_ids,
                execute_fn,
                completed,
            )

            results.update(level_results)

            # Update completed set
            for ws_id in level_ws_ids:
                if results[ws_id][0]:
                    completed.add(ws_id)
                    self._state.results[ws_id] = True
                else:
                    self._state.errors[ws_id] = results[ws_id][1] or "Unknown error"

            self._save_state()

        # Finalize
        duration = (datetime.now() - start_time).total_seconds()
        self._state.status = "completed"
        self._save_state()

        return self._make_result(
            feature_id,
            start_time,
            {ws_id: success for ws_id, (success, _) in results.items()},
            {ws_id: err for ws_id, (_, err) in results.items() if err},
        )

    async def run_feature_ordered(
        self,
        feature_id: str,
        execute_fn: Callable[[str, str], tuple[bool, str | None]] | None = None,
    ) -> ExecutionResult:
        """Execute workstreams sequentially in dependency order.

        Simpler alternative to run_feature for debugging.

        Args:
            feature_id: Feature ID
            execute_fn: Optional execute function

        Returns:
            ExecutionResult with summary
        """
        start_time = datetime.now()
        self._state = OrchestratorState(feature_id=feature_id, started_at=start_time)

        ws_ids = self._get_feature_workstreams(feature_id)
        if not ws_ids:
            return self._make_result(feature_id, start_time, {}, {})

        self._graph.build(ws_ids)
        execution_order = self._graph.execution_order()

        results: dict[str, tuple[bool, str | None]] = {}

        if execute_fn is None:
            execute_fn = self._default_execute

        for ws_id in execution_order:
            self._state.current_ws = ws_id
            self._save_state()

            # Acquire agent and execute
            agent_id = await self._pool.acquire(ws_id)
            success = False
            try:
                success, error = await execute_fn(agent_id, ws_id)
                results[ws_id] = (success, error)

                self._state.results[ws_id] = success
                if error:
                    self._state.errors[ws_id] = error

                # Notify callbacks
                for callback in self._callbacks:
                    callback(ws_id, success, error)

            finally:
                await self._pool.release(agent_id, success)

        duration = (datetime.now() - start_time).total_seconds()
        self._state.status = "completed"
        self._state.current_ws = None
        self._save_state()

        return self._make_result(
            feature_id,
            start_time,
            {ws_id: success for ws_id, (success, _) in results.items()},
            {ws_id: err for ws_id, (_, err) in results.items() if err},
        )

    async def enqueue_feature(
        self,
        feature_id: str,
        priority: Priority = Priority.NORMAL,
    ) -> list[str]:
        """Enqueue all workstreams for a feature in the task queue.

        Args:
            feature_id: Feature ID
            priority: Task priority

        Returns:
            List of task IDs created
        """
        ws_ids = self._get_feature_workstreams(feature_id)
        task_ids: list[str] = []

        for ws_id in ws_ids:
            task = Task(
                ws_id=ws_id,
                priority=priority,
                metadata={"feature_id": feature_id},
            )
            task_id = self._queue.enqueue(task)
            task_ids.append(task_id)

        logger.info(f"Enqueued {len(task_ids)} tasks for feature {feature_id}")
        return task_ids

    def get_stats(self) -> PoolStats:
        """Get current agent pool statistics.

        Returns:
            PoolStats
        """
        return self._pool.get_stats()

    def get_state(self) -> OrchestratorState | None:
        """Get current orchestrator state.

        Returns:
            Current OrchestratorState or None
        """
        return self._state

    def load_state(self) -> OrchestratorState | None:
        """Load state from disk.

        Returns:
            Loaded state or None if no state file exists
        """
        if not self._state_file.exists():
            return None

        try:
            data = json.loads(self._state_file.read_text())

            # Convert datetime string back
            started_at = datetime.fromisoformat(data["started_at"])

            self._state = OrchestratorState(
                feature_id=data["feature_id"],
                started_at=started_at,
                status=data.get("status", "running"),
                results=data.get("results", {}),
                current_ws=data.get("current_ws"),
                errors=data.get("errors", {}),
            )

            return self._state

        except Exception as e:
            logger.error(f"Failed to load state: {e}")
            return None

    def _get_feature_workstreams(self, feature_id: str) -> list[str]:
        """Get all workstream IDs for a feature.

        Args:
            feature_id: Feature ID

        Returns:
            List of ws_ids
        """
        from ..dashboard.sources.workstream_reader import WorkstreamReader

        reader = WorkstreamReader(self._ws_dir)
        state = reader.read()

        return [
            ws.ws_id
            for ws in state.workstreams.values()
            if ws.feature == feature_id
        ]

    async def _execute_level(
        self,
        ws_ids: list[str],
        execute_fn: Callable[[str, str], tuple[bool, str | None]],
        completed: set[str],
    ) -> dict[str, tuple[bool, str | None]]:
        """Execute a level of workstreams in parallel.

        Args:
            ws_ids: Workstream IDs to execute
            execute_fn: Execution function
            completed: Already completed workstreams

        Returns:
            Dict mapping ws_id to (success, error)
        """
        results: dict[str, tuple[bool, str | None]] = {}

        # Create tasks
        async def run_one(ws_id: str) -> tuple[str, bool, str | None]:
            agent_id = await self._pool.acquire(ws_id)
            try:
                success, error = await execute_fn(agent_id, ws_id)

                # Notify callbacks
                for callback in self._callbacks:
                    callback(ws_id, success, error)

                return ws_id, success, error
            finally:
                await self._pool.release(agent_id, success)

        tasks = [asyncio.create_task(run_one(ws_id)) for ws_id in ws_ids]

        # Wait for all to complete
        for task in tasks:
            ws_id, success, error = await task
            results[ws_id] = (success, error)

        return results

    def _group_by_level(self, execution_order: list[str]) -> list[tuple[int, list[str]]]:
        """Group workstreams by dependency level.

        Workstreams at the same level can run in parallel.

        Args:
            execution_order: List of ws_ids in dependency order

        Returns:
            List of (level, [ws_ids]) tuples
        """
        levels: dict[str, int] = {}

        # Calculate level for each ws_id
        for ws_id in execution_order:
            deps = self._graph.get_dependencies(ws_id)
            if not deps:
                levels[ws_id] = 0
            else:
                # Level is max of dependency levels + 1
                dep_levels = [levels.get(dep, 0) for dep in deps if dep in levels]
                levels[ws_id] = max(dep_levels) + 1 if dep_levels else 0

        # Group by level
        level_groups: dict[int, list[str]] = {}
        for ws_id, level in levels.items():
            if level not in level_groups:
                level_groups[level] = []
            level_groups[level].append(ws_id)

        # Return as sorted list
        return sorted(level_groups.items())

    async def _default_execute(
        self,
        agent_id: str,
        ws_id: str,
    ) -> tuple[bool, str | None]:
        """Default execution function - delegates to AgentExecutor.

        Args:
            agent_id: Agent ID (for logging)
            ws_id: Workstream ID

        Returns:
            Tuple of (success, error)
        """
        import asyncio

        from .executor import AgentExecutor

        def run_sync() -> tuple[bool, str | None]:
            try:
                executor = AgentExecutor(str(self._ws_dir), metrics_file=".sdp/orchestrator_metrics.json")
                success = executor.execute(ws_id)
                return success, None
            except Exception as e:
                return False, str(e)

        return await asyncio.get_event_loop().run_in_executor(None, run_sync)

    def _find_cycle(self) -> list[str]:
        """Find circular dependency cycle for error reporting.

        Returns:
            List of ws_ids forming a cycle
        """
        return self._graph._find_cycle()

    def _save_state(self) -> None:
        """Save current state to disk."""
        if not self._state:
            return

        try:
            data = {
                "feature_id": self._state.feature_id,
                "started_at": self._state.started_at.isoformat(),
                "status": self._state.status,
                "results": self._state.results,
                "current_ws": self._state.current_ws,
                "errors": self._state.errors,
            }
            self._state_file.write_text(json.dumps(data, indent=2))
        except Exception as e:
            logger.error(f"Failed to save state: {e}")

    def _make_result(
        self,
        feature_id: str,
        start_time: datetime,
        results: dict[str, bool],
        errors: dict[str, str],
    ) -> ExecutionResult:
        """Create an ExecutionResult.

        Args:
            feature_id: Feature ID
            start_time: Start time
            results: Dict of ws_id to success
            errors: Dict of ws_id to error message

        Returns:
            ExecutionResult
        """
        duration = (datetime.now() - start_time).total_seconds()

        completed = [ws_id for ws_id, success in results.items() if success]
        failed = [ws_id for ws_id, success in results.items() if not success]

        return ExecutionResult(
            feature_id=feature_id,
            success=len(failed) == 0,
            completed=completed,
            failed=failed,
            duration_seconds=duration,
            errors=errors,
        )

    async def resume(self) -> ExecutionResult | None:
        """Resume a previously interrupted execution.

        Returns:
            ExecutionResult if resuming, None if no state to resume
        """
        state = self.load_state()
        if not state or state.status == "completed":
            return None

        # Calculate remaining workstreams
        completed = set(ws_id for ws_id, success in state.results.items() if success)
        all_ws = self._get_feature_workstreams(state.feature_id)
        remaining = [ws_id for ws_id in all_ws if ws_id not in completed]

        if not remaining:
            # Already done
            return None

        # Resume execution
        start_time = datetime.now()
        self._state = state

        results: dict[str, tuple[bool, str | None]] = {}

        # Execute remaining in order
        execution_order = self._graph.execution_order()
        remaining_order = [ws_id for ws_id in execution_order if ws_id in remaining]

        for ws_id in remaining_order:
            self._state.current_ws = ws_id
            self._save_state()

            agent_id = await self._pool.acquire(ws_id)
            try:
                success, error = await self._default_execute(agent_id, ws_id)
                results[ws_id] = (success, error)

                self._state.results[ws_id] = success
                if error:
                    self._state.errors[ws_id] = error

            finally:
                await self._pool.release(agent_id, success)

        duration = (datetime.now() - start_time).total_seconds()
        self._state.status = "completed"
        self._state.current_ws = None
        self._save_state()

        return self._make_result(
            state.feature_id,
            start_time,
            {ws_id: success for ws_id, (success, _) in results.items()},
            {ws_id: err for ws_id, (_, err) in results.items() if err},
        )
