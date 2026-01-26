"""Tests for AgentPool."""

import asyncio
import pytest

from sdp.agents.agent_pool import AgentPool, AgentInfo, PoolStats


@pytest.mark.asyncio
async def test_agent_pool_initialization() -> None:
    """Test pool initializes with correct number of agents."""
    pool = AgentPool(max_agents=3)

    stats = pool.get_stats()
    assert stats.total_agents == 3
    assert stats.busy_agents == 0
    assert stats.available_agents == 3


@pytest.mark.asyncio
async def test_agent_pool_acquire_release() -> None:
    """Test acquiring and releasing agents."""
    pool = AgentPool(max_agents=2)

    # Acquire first agent
    agent_id = await pool.acquire()
    assert agent_id in ["agent-1", "agent-2"]

    stats = pool.get_stats()
    assert stats.busy_agents == 1
    assert stats.available_agents == 1

    # Release agent
    await pool.release(agent_id)

    stats = pool.get_stats()
    assert stats.busy_agents == 0
    assert stats.available_agents == 2


@pytest.mark.asyncio
async def test_agent_pool_concurrent_limit() -> None:
    """Test pool respects max_agents limit."""
    pool = AgentPool(max_agents=2)

    # Acquire both agents
    agent1 = await pool.acquire()
    agent2 = await pool.acquire()

    # Third acquire should block
    acquired = False

    async def try_acquire() -> None:
        nonlocal acquired
        await pool.acquire()
        acquired = True

    task = asyncio.create_task(try_acquire())

    # Wait a bit - should not acquire
    await asyncio.sleep(0.1)
    assert acquired is False

    # Release one agent
    await pool.release(agent1)

    # Now third acquire should succeed
    await asyncio.sleep(0.1)
    assert acquired is True

    # Cleanup
    await pool.release(agent2)


@pytest.mark.asyncio
async def test_agent_pool_least_busy_selection() -> None:
    """Test pool assigns least busy agent."""
    pool = AgentPool(max_agents=3)

    # All agents equal, should get agent-1 first
    agent1 = await pool.acquire()
    assert agent1 == "agent-1"

    # Release and acquire again - should prefer agent-1 (now with completed count)
    await pool.release(agent1, success=True)
    agent2 = await pool.acquire()

    # agent-1 has completed_count=1, others have 0
    # So should prefer agent-2 or agent-3
    assert agent2 in ["agent-2", "agent-3"]


@pytest.mark.asyncio
async def test_agent_pool_get_agent_info() -> None:
    """Test getting specific agent info."""
    pool = AgentPool(max_agents=2)

    agent_id = await pool.acquire("test-ws")
    info = pool.get_agent_info(agent_id)

    assert info is not None
    assert info.agent_id == agent_id
    assert info.busy is True
    assert info.current_ws == "test-ws"
    assert info.started_at is not None

    await pool.release(agent_id)

    info = pool.get_agent_info(agent_id)
    assert info.busy is False
    assert info.current_ws is None


@pytest.mark.asyncio
async def test_agent_pool_list_agents() -> None:
    """Test listing all agents."""
    pool = AgentPool(max_agents=2)

    agents = pool.list_agents()
    assert len(agents) == 2
    assert all(isinstance(a, AgentInfo) for a in agents)


@pytest.mark.asyncio
async def test_agent_pool_execute() -> None:
    """Test execute method."""
    pool = AgentPool(max_agents=2)

    executed = []

    async def mock_work():
        executed.append("done")
        return True, None

    success, error = await pool.execute("test-ws", mock_work())

    assert success is True
    assert error is None
    assert executed == ["done"]


@pytest.mark.asyncio
async def test_agent_pool_execute_with_error() -> None:
    """Test execute method handles errors."""
    pool = AgentPool(max_agents=2)

    async def mock_work():
        raise ValueError("Test error")

    success, error = await pool.execute("test-ws", mock_work())

    assert success is False
    assert "Test error" in error


@pytest.mark.asyncio
async def test_agent_pool_context_manager() -> None:
    """Test pool as context manager."""
    async with AgentPool(max_agents=2) as pool:
        agent1 = await pool.acquire()
        agent2 = await pool.acquire()

        assert pool.get_stats().busy_agents == 2

        await pool.release(agent1)
        await pool.release(agent2)

    # After exiting, all agents should be released


@pytest.mark.asyncio
async def test_agent_pool_iterator() -> None:
    """Test AgentPoolIterator."""
    from sdp.agents.agent_pool import AgentPoolIterator

    pool = AgentPool(max_agents=2)

    results = {}

    async def execute_fn(agent_id: str, ws_id: str) -> tuple[bool, str | None]:
        await asyncio.sleep(0.01)
        results[ws_id] = True
        return True, None

    iterator = AgentPoolIterator(
        pool,
        ["ws-001", "ws-002"],
        execute_fn,
    )

    final_results = await iterator.run_all()

    assert len(final_results) == 2
    assert "ws-001" in results
    assert "ws-002" in results


@pytest.mark.asyncio
async def test_agent_pool_iterator_iter_results() -> None:
    """Test AgentPoolIterator iter_results."""
    from sdp.agents.agent_pool import AgentPoolIterator

    pool = AgentPool(max_agents=2)
    results = []

    async def execute_fn(agent_id: str, ws_id: str) -> tuple[bool, str | None]:
        await asyncio.sleep(0.01)
        return True, None

    iterator = AgentPoolIterator(
        pool,
        ["ws-001", "ws-002"],
        execute_fn,
    )

    async for ws_id, success, error in iterator.iter_results():
        results.append((ws_id, success, error))

    assert len(results) == 2


def test_agent_info_elapsed_seconds() -> None:
    """Test AgentInfo elapsed_seconds calculation."""
    from datetime import datetime, timedelta

    info = AgentInfo(agent_id="test")
    assert info.elapsed_seconds == 0.0

    info.started_at = datetime.now() - timedelta(seconds=5)
    # Should be approximately 5 seconds
    assert 4.9 <= info.elapsed_seconds <= 5.1
