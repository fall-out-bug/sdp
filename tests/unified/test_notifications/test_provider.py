"""Tests for NotificationProvider module.

Tests notification provider interface, console and in-memory
implementations, and notification delivery functionality.
"""

import pytest
from abc import ABC
from datetime import datetime
from io import StringIO
from unittest.mock import Mock

from sdp.unified.notifications.provider import (
    NotificationProvider,
    NotificationType,
    Notification,
    ConsoleNotificationProvider,
    InMemoryNotificationProvider,
)


class TestNotificationTypeEnum:
    """Test NotificationType enum."""

    def test_has_info_type(self):
        """Should have INFO type."""
        assert NotificationType.INFO is not None
        assert NotificationType.INFO.value == "info"

    def test_has_warning_type(self):
        """Should have WARNING type."""
        assert NotificationType.WARNING is not None

    def test_has_error_type(self):
        """Should have ERROR type."""
        assert NotificationType.ERROR is not None

    def test_has_success_type(self):
        """Should have SUCCESS type."""
        assert NotificationType.SUCCESS is not None


class TestNotificationDataclass:
    """Test Notification dataclass."""

    def test_create_notification_with_required_fields(self):
        """Should create notification with required fields."""
        notification = Notification(
            type=NotificationType.INFO,
            message="Test message",
        )

        assert notification.type == NotificationType.INFO
        assert notification.message == "Test message"
        assert notification.timestamp is not None

    def test_create_notification_with_recipient(self):
        """Should create notification with recipient."""
        notification = Notification(
            type=NotificationType.ERROR,
            message="Error occurred",
            recipient="agent-123",
        )

        assert notification.recipient == "agent-123"

    def test_notification_has_timestamp(self):
        """Should have automatic timestamp."""
        from datetime import timezone

        before = datetime.now(timezone.utc)
        notification = Notification(
            type=NotificationType.SUCCESS,
            message="Success",
        )
        after = datetime.now(timezone.utc)

        assert notification.timestamp >= before
        assert notification.timestamp <= after


class TestNotificationProviderInterface:
    """Test NotificationProvider abstract interface."""

    def test_is_abstract_base_class(self):
        """Should be abstract base class."""
        assert issubclass(NotificationProvider, ABC)

    def test_requires_send_method(self):
        """Should require send method implementation."""

        class IncompleteProvider(NotificationProvider):
            pass

        with pytest.raises(TypeError):
            IncompleteProvider()

    def test_can_create_concrete_implementation(self):
        """Should allow concrete implementation."""

        class TestProvider(NotificationProvider):
            def send(self, notification: Notification) -> bool:
                return True

        provider = TestProvider()
        assert provider is not None


class TestConsoleNotificationProvider:
    """Test ConsoleNotificationProvider implementation."""

    def test_creates_provider(self):
        """Should initialize console provider."""
        provider = ConsoleNotificationProvider()

        assert provider is not None
        assert hasattr(provider, 'send')

    def test_sends_info_notification(self, capsys):
        """Should send info notification to console."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="Info message",
        )

        result = provider.send(notification)

        assert result is True
        captured = capsys.readouterr()
        assert "Info message" in captured.out

    def test_sends_warning_notification(self, capsys):
        """Should send warning notification to console."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.WARNING,
            message="Warning message",
        )

        provider.send(notification)

        captured = capsys.readouterr()
        assert "Warning" in captured.out

    def test_sends_error_notification(self, capsys):
        """Should send error notification to console."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.ERROR,
            message="Error message",
        )

        provider.send(notification)

        captured = capsys.readouterr()
        # Errors go to stderr, not stdout
        assert "ERROR" in captured.err

    def test_sends_success_notification(self, capsys):
        """Should send success notification to console."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.SUCCESS,
            message="Success message",
        )

        provider.send(notification)

        captured = capsys.readouterr()
        assert "Success" in captured.out

    def test_includes_recipient_in_output(self, capsys):
        """Should include recipient in console output."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="Message",
            recipient="agent-123",
        )

        provider.send(notification)

        captured = capsys.readouterr()
        assert "agent-123" in captured.out

    def test_formats_notification_with_timestamp(self, capsys):
        """Should format notification with timestamp."""
        provider = ConsoleNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="Message",
        )

        provider.send(notification)

        captured = capsys.readouterr()
        # Should have timestamp
        assert notification.timestamp.strftime("%Y-%m-%d") in captured.out


class TestInMemoryNotificationProvider:
    """Test InMemoryNotificationProvider implementation."""

    def test_creates_provider(self):
        """Should initialize in-memory provider."""
        provider = InMemoryNotificationProvider()

        assert provider is not None
        assert hasattr(provider, 'send')
        assert hasattr(provider, 'get_notifications')

    def test_sends_notification(self):
        """Should send notification to memory."""
        provider = InMemoryNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        result = provider.send(notification)

        assert result is True

    def test_stores_notifications(self):
        """Should store sent notifications."""
        provider = InMemoryNotificationProvider()
        notification = Notification(
            type=NotificationType.WARNING,
            message="Warning",
        )

        provider.send(notification)
        notifications = provider.get_notifications()

        assert len(notifications) == 1
        assert notifications[0].message == "Warning"

    def test_gets_notifications_by_type(self):
        """Should filter notifications by type."""
        provider = InMemoryNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "Info"))
        provider.send(Notification(NotificationType.ERROR, "Error"))
        provider.send(Notification(NotificationType.INFO, "Info 2"))

        info_notifications = provider.get_notifications_by_type(NotificationType.INFO)

        assert len(info_notifications) == 2

    def test_gets_notifications_by_recipient(self):
        """Should filter notifications by recipient."""
        provider = InMemoryNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "Msg", "agent-1"))
        provider.send(Notification(NotificationType.INFO, "Msg", "agent-2"))
        provider.send(Notification(NotificationType.INFO, "Msg", "agent-1"))

        agent1_notifications = provider.get_notifications_by_recipient("agent-1")

        assert len(agent1_notifications) == 2

    def test_clears_notifications(self):
        """Should clear all notifications."""
        provider = InMemoryNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "Msg"))
        provider.clear()

        notifications = provider.get_notifications()

        assert len(notifications) == 0

    def test_counts_notifications(self):
        """Should count notifications."""
        provider = InMemoryNotificationProvider()

        assert provider.count() == 0

        provider.send(Notification(NotificationType.INFO, "Msg"))
        provider.send(Notification(NotificationType.WARNING, "Msg"))

        assert provider.count() == 2


class TestNotificationDelivery:
    """Test notification delivery functionality."""

    def test_handles_delivery_failure_in_console(self):
        """Should handle console delivery failure gracefully."""
        # This is more about ensuring no exceptions are raised
        provider = ConsoleNotificationProvider()
        notification = Notification(NotificationType.INFO, "Test")

        # Should not raise exception
        result = provider.send(notification)
        assert result is True

    def test_handles_multiple_notifications(self):
        """Should handle multiple notifications in sequence."""
        provider = InMemoryNotificationProvider()

        for i in range(5):
            notification = Notification(
                type=NotificationType.INFO,
                message=f"Message {i}",
            )
            provider.send(notification)

        assert provider.count() == 5

    def test_preserves_notification_order(self):
        """Should preserve notification order."""
        provider = InMemoryNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "First"))
        provider.send(Notification(NotificationType.WARNING, "Second"))
        provider.send(Notification(NotificationType.ERROR, "Third"))

        notifications = provider.get_notifications()

        assert notifications[0].message == "First"
        assert notifications[1].message == "Second"
        assert notifications[2].message == "Third"


class TestErrorHandling:
    """Test error handling."""

    def test_handles_empty_message(self):
        """Should handle empty message gracefully."""
        provider = InMemoryNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="",
        )

        result = provider.send(notification)

        assert result is True
        assert provider.count() == 1

    def test_handles_none_recipient(self):
        """Should handle None recipient gracefully."""
        provider = InMemoryNotificationProvider()
        notification = Notification(
            type=NotificationType.INFO,
            message="Message",
            recipient=None,
        )

        result = provider.send(notification)

        assert result is True

    def test_handles_special_characters(self):
        """Should handle special characters in message."""
        provider = InMemoryNotificationProvider()
        notification = Notification(
            type=NotificationType.ERROR,
            message="Error: <>&\"'\\n\\t",
        )

        result = provider.send(notification)

        assert result is True
        assert provider.count() == 1
