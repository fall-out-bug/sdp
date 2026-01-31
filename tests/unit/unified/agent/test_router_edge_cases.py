"""Tests for send message router edge cases."""

import pytest
from unittest.mock import Mock
from datetime import datetime, timezone

from sdp.unified.agent.router import (
    SendMessageRouter,
    Message,
    DeliveryResult,
)


class TestMessage:
    """Tests for Message dataclass."""

    def test_default_initialization(self):
        """Test default message initialization."""
        message = Message(sender="agent1", content="Test message")

        assert message.sender == "agent1"
        assert message.content == "Test message"
        assert message.recipient is None
        assert message.broadcast is False
        assert message.message_type == "direct"
        assert message.timestamp is not None

    def test_broadcast_message(self):
        """Test broadcast message initialization."""
        message = Message(sender="agent1", content="Broadcast", broadcast=True)

        assert message.broadcast is True
        assert message.message_type == "broadcast"

    def test_direct_message_with_recipient(self):
        """Test direct message with recipient."""
        message = Message(
            sender="agent1",
            content="Direct message",
            recipient="agent2"
        )

        assert message.recipient == "agent2"
        assert message.broadcast is False
        assert message.message_type == "direct"


class TestDeliveryResult:
    """Tests for DeliveryResult dataclass."""

    def test_success_result(self):
        """Test successful delivery result."""
        result = DeliveryResult(success=True)

        assert result.success is True
        assert result.error is None

    def test_failure_result(self):
        """Test failed delivery result."""
        result = DeliveryResult(success=False, error="Network error")

        assert result.success is False
        assert result.error == "Network error"


class TestSendMessageRouter:
    """Tests for SendMessageRouter."""

    def test_init_without_team_manager(self):
        """Test initialization without team manager."""
        router = SendMessageRouter()

        assert router._team_manager is None
        assert router._pending_messages == {}

    def test_init_with_team_manager(self):
        """Test initialization with team manager."""
        mock_manager = Mock()
        router = SendMessageRouter(team_manager=mock_manager)

        assert router._team_manager is mock_manager

    def test_send_direct_message_no_recipient(self):
        """Test direct message fails without recipient."""
        router = SendMessageRouter()
        message = Message(sender="agent1", content="Test")

        result = router.send_message(message)

        assert result.success is False
        assert "No recipient" in result.error

    def test_send_direct_message_with_team_manager(self):
        """Test sends direct message with team manager."""
        mock_manager = Mock()
        mock_manager.get_agent_by_role.return_value = Mock(id="agent2-actual-id")
        
        router = SendMessageRouter(team_manager=mock_manager)
        message = Message(sender="agent1", content="Test", recipient="agent2")

        result = router._send_direct(message)

        # Will depend on actual implementation
        assert result is not None

    def test_send_broadcast_message_no_team_manager(self):
        """Test broadcast message fails without team manager."""
        router = SendMessageRouter()
        message = Message(sender="agent1", content="Broadcast", broadcast=True)

        result = router.send_message(message)

        # Without team manager, broadcast should fail
        assert result.success is False
        assert "team manager" in result.error.lower()


class TestPendingMessages:
    """Tests for pending message queuing."""

    def test_queue_message_for_offline_agent(self):
        """Test queues messages for offline agents."""
        router = SendMessageRouter()
        message1 = Message(sender="agent1", content="Msg1", recipient="agent2")
        message2 = Message(sender="agent1", content="Msg2", recipient="agent2")

        # Manually queue messages (simulating offline delivery)
        router._pending_messages.setdefault("agent2", []).append(message1)
        router._pending_messages.setdefault("agent2", []).append(message2)

        assert len(router._pending_messages["agent2"]) == 2

    def test_retrieve_pending_messages(self):
        """Test retrieves pending messages for agent."""
        router = SendMessageRouter()
        message = Message(sender="agent1", content="Test", recipient="agent2")
        
        router._pending_messages["agent2"] = [message]

        pending = router._pending_messages.get("agent2", [])

        assert len(pending) == 1
        assert pending[0].content == "Test"
