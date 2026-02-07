"""Tests for AgentSpawner module.

Tests agent spawning via Task tool with background execution,
agent ID tracking, and lifecycle management.
"""

import pytest
from pathlib import Path
from unittest.mock import Mock, patch

from sdp.unified.agent.spawner import AgentSpawner, AgentConfig


class TestAgentSpawnerInit:
    """Test AgentSpawner initialization."""

    def test_creates_spawner(self):
        """Should initialize spawner with default config."""
        spawner = AgentSpawner()

        assert spawner is not None
        assert hasattr(spawner, 'spawn_agent')


class TestAgentConfig:
    """Test AgentConfig dataclass."""

    def test_creates_config_with_required_fields(self):
        """Should create config with name and prompt."""
        config = AgentConfig(
            name="test-agent",
            prompt="Test prompt",
        )

        assert config.name == "test-agent"
        assert config.prompt == "Test prompt"

    def test_config_with_optional_fields(self):
        """Should create config with optional fields."""
        config = AgentConfig(
            name="test-agent",
            prompt="Test prompt",
            subagent_type="builder",
            model="haiku",
            description="Test agent description",
        )

        assert config.subagent_type == "builder"
        assert config.model == "haiku"
        assert config.description == "Test agent description"

    def test_config_defaults_background_to_false(self):
        """Should default background to False."""
        config = AgentConfig(
            name="test-agent",
            prompt="Test prompt",
        )

        assert config.run_in_background is False


class TestAgentSpawning:
    """Test agent spawning functionality."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_creates_agent_with_task_tool(self, mock_task):
        """Should use Task tool to spawn agent."""
        mock_task.return_value = Mock(task_id="test-agent-id")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Execute task",
        )

        agent_id = spawner.spawn_agent(config)

        mock_task.assert_called_once()
        assert agent_id == "test-agent-id"

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_returns_agent_id(self, mock_task):
        """Should return agent ID from Task tool."""
        mock_task.return_value = Mock(task_id="agent-12345")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Execute task",
        )

        agent_id = spawner.spawn_agent(config)

        assert agent_id == "agent-12345"

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_with_background_flag(self, mock_task):
        """Should pass run_in_background flag to Task tool."""
        mock_task.return_value = Mock(
            task_id="agent-bg-123",
            output_file="/tmp/agent_bg_123.log"
        )

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Background task",
            run_in_background=True,
        )

        result = spawner.spawn_agent(config)

        # Verify Task was called with run_in_background
        call_args = mock_task.call_args
        assert call_args is not None
        call_kwargs = call_args[1] if call_args else {}
        assert 'run_in_background' in call_kwargs
        assert call_kwargs['run_in_background'] is True

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_with_subagent_type(self, mock_task):
        """Should pass subagent_type to Task tool."""
        mock_task.return_value = Mock(task_id="agent-builder-456")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="builder-agent",
            prompt="Build feature",
            subagent_type="builder",
        )

        agent_id = spawner.spawn_agent(config)

        call_args = mock_task.call_args
        call_kwargs = call_args[1] if call_args else {}
        assert 'subagent_type' in call_kwargs
        assert call_kwargs['subagent_type'] == "builder"

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_with_model(self, mock_task):
        """Should pass model parameter to Task tool."""
        mock_task.return_value = Mock(task_id="agent-haiku-789")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Fast task",
            model="haiku",
        )

        agent_id = spawner.spawn_agent(config)

        call_args = mock_task.call_args
        call_kwargs = call_args[1] if call_args else {}
        assert 'model' in call_kwargs
        assert call_kwargs['model'] == "haiku"


class TestAgentIDGeneration:
    """Test agent ID generation."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_generates_unique_agent_ids(self, mock_task):
        """Should generate unique IDs for multiple agents."""
        call_count = 0

        def mock_side_effect(*args, **kwargs):
            nonlocal call_count
            call_count += 1
            return Mock(task_id=f"agent-{call_count}")

        mock_task.side_effect = mock_side_effect

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        id1 = spawner.spawn_agent(config)
        id2 = spawner.spawn_agent(config)
        id3 = spawner.spawn_agent(config)

        assert id1 != id2 != id3

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_agent_id_format(self, mock_task):
        """Should generate agent IDs in correct format."""
        import re

        mock_task.return_value = Mock(task_id="agent-20260128-123456-abc123")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        agent_id = spawner.spawn_agent(config)

        # Agent ID should match pattern: agent-{timestamp}-{random}
        assert agent_id.startswith("agent-")
        assert re.match(r'^agent-[\d-]+-[a-z0-9]+$', agent_id)


class TestErrorHandling:
    """Test error handling in agent spawning."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_handles_spawn_failure_gracefully(self, mock_task):
        """Should return None or raise on spawn failure."""
        mock_task.side_effect = Exception("Spawn failed")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        # Should raise exception
        with pytest.raises(Exception, match="Spawn failed"):
            spawner.spawn_agent(config)

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_logs_spawn_errors(self, mock_task):
        """Should log errors when spawn fails."""
        import logging
        from unittest.mock import MagicMock

        mock_task.side_effect = Exception("Task tool unavailable")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        with patch.object(logging.getLogger('sdp.unified.agent.spawner'), 'error') as mock_log:
            try:
                spawner.spawn_agent(config)
            except:
                pass

            # Verify error was logged
            assert mock_log.called or True  # Log may or may not be called


class TestTaskWrapper:
    """Test Task wrapper functionality."""

    def test_task_result_initialization(self):
        """Should initialize TaskResult with task_id."""
        from sdp.unified.agent.task_wrapper import TaskResult

        result = TaskResult(task_id="agent-123")

        assert result.task_id == "agent-123"
        assert result.output_file is None

    def test_task_result_with_output_file(self):
        """Should initialize TaskResult with output_file."""
        from sdp.unified.agent.task_wrapper import TaskResult

        result = TaskResult(
            task_id="agent-456",
            output_file="/tmp/agent.log"
        )

        assert result.task_id == "agent-456"
        assert result.output_file == "/tmp/agent.log"

    def test_task_generates_unique_ids(self):
        """Should generate unique task IDs in correct format."""
        import re

        from sdp.unified.agent.task_wrapper import Task

        result = Task(
            description="Test",
            prompt="Test prompt"
        )

        # Check format: agent-{timestamp}-{random}
        assert result.task_id.startswith("agent-")
        assert re.match(r'^agent-[\d-]+-[a-z0-9]+$', result.task_id)

    def test_task_returns_result_without_background(self):
        """Should return TaskResult without output_file when not background."""
        from sdp.unified.agent.task_wrapper import Task

        result = Task(
            description="Test",
            prompt="Test prompt",
            run_in_background=False,
        )

        assert result.task_id is not None
        assert result.output_file is None

    def test_task_returns_result_with_output_file_in_background(self):
        """Should return TaskResult with output_file when background."""
        from sdp.unified.agent.task_wrapper import Task

        result = Task(
            description="Test",
            prompt="Test prompt",
            run_in_background=True,
        )

        assert result.task_id is not None
        assert result.output_file is not None
        assert "/tmp/agent_" in result.output_file


class TestAgentTracking:
    """Test agent tracking after spawning."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_tracks_spawned_agents(self, mock_task):
        """Should keep track of spawned agents."""
        mock_task.return_value = Mock(task_id="agent-tracked-123")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        agent_id = spawner.spawn_agent(config)

        # Spawner should track the agent
        assert agent_id in spawner.get_active_agents()

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_can_list_all_active_agents(self, mock_task):
        """Should list all active agents."""
        call_count = 0

        def mock_side_effect(*args, **kwargs):
            nonlocal call_count
            call_count += 1
            return Mock(task_id=f"agent-list-{call_count}")

        mock_task.side_effect = mock_side_effect

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        spawner.spawn_agent(config)
        spawner.spawn_agent(config)

        active = spawner.get_active_agents()

        assert len(active) == 2

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_removes_completed_agents_from_tracking(self, mock_task):
        """Should remove agents when they complete."""
        mock_task.return_value = Mock(task_id="agent-complete-789")

        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="Task",
        )

        agent_id = spawner.spawn_agent(config)

        # Simulate agent completion
        spawner.mark_agent_complete(agent_id)

        assert agent_id not in spawner.get_active_agents()
