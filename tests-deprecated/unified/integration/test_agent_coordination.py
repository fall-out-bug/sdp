"""Integration tests for agent coordination workflows.

Tests end-to-end workflows involving multiple components:
agent spawning, inter-agent messaging, role loading, and notifications.
"""

from unittest.mock import MagicMock, patch

import pytest
from sdp.unified.agent.bug_report import BugReportFlow, BugSeverity, BugStatus
from sdp.unified.agent.role_loader import RoleLoader
from sdp.unified.agent.role_state import RoleState, RoleStateManager
from sdp.unified.agent.router import Message, SendMessageRouter
from sdp.unified.agent.spawner import AgentConfig, AgentSpawner
from sdp.unified.notifications.mock import MockNotificationProvider
from sdp.unified.notifications.provider import Notification, NotificationType
from sdp.unified.notifications.router import NotificationRouter


class TestAgentSpawningAndMessaging:
    """Test agent spawning with inter-agent messaging."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_agent_and_send_message(self, mock_task):
        """Test spawning agent and sending message."""
        mock_task.return_value = MagicMock(task_id="agent-123")

        # Spawn agent
        spawner = AgentSpawner()
        config = AgentConfig(
            name="test-agent",
            prompt="You are a test agent",
        )
        agent_id = spawner.spawn_agent(config)

        assert agent_id == "agent-123"

        # Send message to spawned agent
        router = SendMessageRouter()
        message = Message(
            sender="orchestrator",
            content="Hello from orchestrator",
            recipient=agent_id,
        )

        result = router.send_message(message)

        assert result.success is True


class TestRoleLoadingAndStateManagement:
    """Test role loading with role state management."""

    def test_load_role_and_activate(self, tmp_path):
        """Test loading role and activating it."""
        # Create role file
        role_file = tmp_path / "planner.md"
        role_file.write_text("""# Planner

You are a planning agent.

**Capabilities:**
- Plan features
- Decompose tasks
""")

        # Load role
        loader = RoleLoader(agents_dir=tmp_path)
        role = loader.load_role("planner")

        assert role is not None
        assert len(role.capabilities) == 2

        # Activate role
        state_mgr = RoleStateManager()
        state = state_mgr.activate_role("planner")

        assert state == RoleState.ACTIVE
        assert state_mgr.is_active("planner")


class TestNotificationRouting:
    """Test notification routing across multiple providers."""

    def test_broadcast_to_multiple_providers(self):
        """Test broadcasting notification to all providers."""
        router = NotificationRouter()

        # Register multiple providers
        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        router.register_provider(provider1)
        router.register_provider(provider2)

        # Broadcast notification
        notification = Notification(
            type=NotificationType.INFO,
            message="Broadcast to all",
        )
        router.broadcast(notification)

        # Both providers should receive
        assert len(provider1.get_notifications()) == 1
        assert len(provider2.get_notifications()) == 1


class TestBugReportWorkflow:
    """Test bug report creation and blocking workflow."""

    def test_bug_blocks_workstream(self):
        """Test bug blocking workstream execution."""
        bug_flow = BugReportFlow()

        # Create blocking bug
        bug = bug_flow.create_report(
            title="Login fails",
            description="Critical issue",
            severity=BugSeverity.P1,
            workstream_id="WS-001",
        )

        # Check blocking
        blocking = bug_flow.get_blocking_workstreams()

        assert "WS-001" in blocking
        assert bug in bug_flow.get_blocking_bugs()

    def test_bug_resolved_unblocks_workstream(self):
        """Test resolving bug unblocks workstream."""
        bug_flow = BugReportFlow()

        # Create blocking bug
        bug = bug_flow.create_report(
            title="Bug",
            description="Issue",
            severity=BugSeverity.P0,
            workstream_id="WS-001",
        )

        # Initially blocking
        assert "WS-001" in bug_flow.get_blocking_workstreams()

        # Mark as resolved
        bug_flow.update_status(bug.id, BugStatus.RESOLVED)

        # Should no longer block (if we filter resolved bugs)
        blocking = [b for b in bug_flow.get_blocking_bugs() if b.status != BugStatus.RESOLVED]
        assert len(blocking) == 0


class TestMultiComponentWorkflow:
    """Test workflows across multiple components."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_agent_lifecycle_with_notifications(self, mock_task):
        """Test complete agent lifecycle with notifications."""
        mock_task.return_value = MagicMock(task_id="agent-lifecycle")

        # Setup notification router
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        # Spawn agent
        spawner = AgentSpawner()
        config = AgentConfig(
            name="lifecycle-agent",
            prompt="Test agent",
        )
        agent_id = spawner.spawn_agent(config)

        # Send notification about spawn
        notification = Notification(
            type=NotificationType.SUCCESS,
            message=f"Agent {agent_id} spawned successfully",
        )
        router.broadcast(notification)

        # Verify notification sent
        notifications = provider.get_notifications()
        assert len(notifications) == 1
        assert "spawned successfully" in notifications[0].message

    def test_role_switch_with_notifications(self, tmp_path):
        """Test role switching with notifications."""
        # Create role files
        for role_name in ["role-a", "role-b"]:
            role_file = tmp_path / f"{role_name}.md"
            role_file.write_text(f"# {role_name.title()}\nTest role")

        # Setup notification router
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        # Load roles
        loader = RoleLoader(agents_dir=tmp_path)
        role_a = loader.load_role("role-a")
        role_b = loader.load_role("role-b")

        # Setup state manager
        state_mgr = RoleStateManager()

        # Activate role-a
        state_mgr.activate_role("role-a")
        assert state_mgr.is_active("role-a")

        # Switch to role-b
        state_mgr.deactivate_role("role-a")
        state_mgr.activate_role("role-b")

        # Verify switch
        assert not state_mgr.is_active("role-a")
        assert state_mgr.is_active("role-b")

        # Send notification about role switch
        notification = Notification(
            type=NotificationType.INFO,
            message="Switched from role-a to role-b",
        )
        router.route(notification)

        assert len(provider.get_notifications()) == 1


class TestErrorHandlingIntegration:
    """Test error handling across components."""

    @patch('sdp.unified.agent.task_wrapper.Task')
    def test_spawn_failure_with_notification(self, mock_task):
        """Test handling spawn failure with notification."""
        # Make Task fail
        mock_task.side_effect = Exception("Spawn failed")

        # Setup notification router
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        # Try to spawn
        spawner = AgentSpawner()
        config = AgentConfig(
            name="failing-agent",
            prompt="Test",
        )

        with pytest.raises(Exception):
            spawner.spawn_agent(config)

        # Send error notification
        notification = Notification(
            type=NotificationType.ERROR,
            message="Failed to spawn agent: Spawn failed",
        )
        router.route(notification)

        assert len(provider.get_notifications()) == 1

    def test_role_load_failure_graceful_handling(self):
        """Test graceful handling of missing role file."""
        loader = RoleLoader()

        # Try to load non-existent role
        role = loader.load_role("nonexistent")

        assert role is None

        # Should not raise exception
        state_mgr = RoleStateManager()
        state_mgr.activate_role("nonexistent")  # Just tracks state

        assert state_mgr.is_active("nonexistent")


class TestCheckpointWorkflow:
    """Test checkpoint save/restore workflows."""

    def test_save_and_restore_agent_state(self):
        """Test saving and restoring agent state."""
        state_mgr = RoleStateManager()

        # Save state
        state_mgr.activate_role("planner")
        state_mgr.activate_role("builder")

        # Create checkpoint data
        checkpoint_data = {
            "active_roles": state_mgr.list_active(),
            "dormant_roles": state_mgr.list_dormant(),
        }

        # Verify checkpoint
        assert "planner" in checkpoint_data["active_roles"]
        assert "builder" in checkpoint_data["active_roles"]

        # Restore state (in new manager)
        new_state_mgr = RoleStateManager()
        for role_name in checkpoint_data["active_roles"]:
            new_state_mgr.activate_role(role_name)

        assert new_state_mgr.is_active("planner")
        assert new_state_mgr.is_active("builder")
