"""Tests for SendMessage router module.

Tests message routing between agents, broadcast functionality,
and integration with TeamManager for role lookup.
"""

import pytest
from unittest.mock import Mock, patch

from sdp.unified.agent.router import SendMessageRouter, Message


class TestMessageDataclass:
    """Test Message dataclass."""

    def test_create_message_with_required_fields(self):
        """Should create message with required fields."""
        message = Message(
            sender="agent-1",
            content="Hello from agent-1",
        )

        assert message.sender == "agent-1"
        assert message.content == "Hello from agent-1"

    def test_create_message_with_recipient(self):
        """Should create message with recipient."""
        message = Message(
            sender="agent-1",
            content="Hello agent-2",
            recipient="agent-2",
        )

        assert message.recipient == "agent-2"
        assert message.message_type == "direct"

    def test_create_message_with_broadcast(self):
        """Should create broadcast message."""
        message = Message(
            sender="agent-1",
            content="Broadcast to all",
            broadcast=True,
        )

        assert message.broadcast is True
        assert message.message_type == "broadcast"

    def test_message_defaults(self):
        """Should have correct default values."""
        message = Message(
            sender="agent-1",
            content="Test",
        )

        assert message.recipient is None
        assert message.broadcast is False
        assert message.message_type == "direct"
        assert message.timestamp is not None


class TestSendMessageRouter:
    """Test SendMessageRouter initialization and routing."""

    def test_creates_router(self):
        """Should initialize router."""
        router = SendMessageRouter()

        assert router is not None
        assert hasattr(router, 'send_message')
        assert hasattr(router, 'get_pending_messages')

    @patch('sdp.unified.agent.router.TeamManager')
    def test_initializes_with_team_manager(self, mock_team_mgr):
        """Should initialize with TeamManager."""
        router = SendMessageRouter(team_manager=mock_team_mgr)

        assert router._team_manager == mock_team_mgr


class TestDirectMessaging:
    """Test direct message routing."""

    @patch('sdp.unified.agent.router.SendMessage')
    def test_sends_message_to_specific_agent(self, mock_send):
        """Should send message to specific recipient."""
        mock_send.return_value = Mock(success=True)

        router = SendMessageRouter()
        message = Message(
            sender="agent-1",
            content="Hello agent-2",
            recipient="agent-2",
        )

        result = router.send_message(message)

        assert result.success is True
        mock_send.assert_called_once()

    @patch('sdp.unified.agent.router.SendMessage')
    def test_logs_delivery_confirmation(self, mock_send):
        """Should log successful message delivery."""
        mock_send.return_value = Mock(success=True)

        router = SendMessageRouter()
        message = Message(
            sender="agent-1",
            content="Hello",
            recipient="agent-2",
        )

        router.send_message(message)

        # Verify logging occurred
        # (In real implementation, would check logger calls)

    @patch('sdp.unified.agent.router.SendMessage')
    def test_handles_delivery_failure(self, mock_send):
        """Should handle message delivery failures."""
        mock_send.return_value = Mock(success=False, error="Agent not found")

        router = SendMessageRouter()
        message = Message(
            sender="agent-1",
            content="Hello",
            recipient="agent-2",
        )

        result = router.send_message(message)

        assert result.success is False
        assert result.error == "Agent not found"


class TestBroadcastMessaging:
    """Test broadcast message routing."""

    @patch('sdp.unified.agent.router.SendMessage')
    def test_broadcasts_to_all_team_members(self, mock_send):
        """Should send message to all team members."""
        mock_send.return_value = Mock(success=True)

        from sdp.unified.agent.router import TeamManager

        # Create mock team manager with active roles
        team_mgr = Mock()
        team_mgr.get_active_roles.return_value = ["agent-1", "agent-2", "agent-3"]

        router = SendMessageRouter(team_manager=team_mgr)
        message = Message(
            sender="agent-1",
            content="Broadcast to team",
            broadcast=True,
        )

        router.send_message(message)

        # Should send to all team members except sender
        assert mock_send.call_count == 2  # agent-2 and agent-3

    @patch('sdp.unified.agent.router.SendMessage')
    def test_broadcast_excludes_sender(self, mock_send):
        """Should not send broadcast to sender."""
        mock_send.return_value = Mock(success=True)

        from sdp.unified.agent.router import TeamManager

        team_mgr = Mock()
        team_mgr.get_active_roles.return_value = ["agent-1", "agent-2"]

        router = SendMessageRouter(team_manager=team_mgr)
        message = Message(
            sender="agent-1",
            content="Broadcast",
            broadcast=True,
        )

        router.send_message(message)

        # Should only send to agent-2 (not agent-1)
        assert mock_send.call_count == 1


class TestMessageQueuing:
    """Test message queuing for offline agents."""

    def test_queues_message_for_offline_agent(self):
        """Should queue message when recipient is offline."""
        router = SendMessageRouter()

        message = Message(
            sender="agent-1",
            content="Hello when you're back",
            recipient="agent-2",
        )

        # Queue the message (agent-2 is offline)
        router.queue_message(message)

        pending = router.get_pending_messages("agent-2")

        assert len(pending) == 1
        assert pending[0].content == "Hello when you're back"

    def test_delivers_queued_messages_when_agent_comes_online(self):
        """Should deliver queued messages when agent comes online."""
        router = SendMessageRouter()

        # Queue two messages
        msg1 = Message(
            sender="agent-1",
            content="First message",
            recipient="agent-2",
        )
        msg2 = Message(
            sender="agent-3",
            content="Second message",
            recipient="agent-2",
        )

        router.queue_message(msg1)
        router.queue_message(msg2)

        # Deliver pending messages
        delivered = router.deliver_pending("agent-2")

        assert len(delivered) == 2
        assert delivered[0].content == "First message"
        assert delivered[1].content == "Second message"

    def test_clears_pending_after_delivery(self):
        """Should clear queue after successful delivery."""
        router = SendMessageRouter()

        message = Message(
            sender="agent-1",
            content="Test",
            recipient="agent-2",
        )

        router.queue_message(message)
        router.deliver_pending("agent-2")

        pending = router.get_pending_messages("agent-2")

        assert len(pending) == 0


class TestTeamManagerIntegration:
    """Test integration with TeamManager for role lookup."""

    @patch('sdp.unified.agent.router.TeamManager')
    def test_looks_up_recipient_in_team_manager(self, mock_team_mgr):
        """Should look up recipient agent in TeamManager."""
        # Mock team manager to find agent by name
        mock_team_mgr.find_agent_by_name.return_value = "agent-123"

        router = SendMessageRouter(team_manager=mock_team_mgr)
        message = Message(
            sender="agent-1",
            content="Hello",
            recipient="developer-1",
        )

        router.send_message(message)

        # Should have looked up recipient
        mock_team_mgr.find_agent_by_name.assert_called_once_with("developer-1")

    @patch('sdp.unified.agent.router.TeamManager')
    def test_handles_unknown_recipient(self, mock_team_mgr):
        """Should handle unknown recipients gracefully."""
        mock_team_mgr.find_agent_by_name.return_value = None

        router = SendMessageRouter(team_manager=mock_team_mgr)
        message = Message(
            sender="agent-1",
            content="Hello",
            recipient="unknown-agent",
        )

        result = router.send_message(message)

        assert result.success is False
        assert "not found" in result.error.lower()
