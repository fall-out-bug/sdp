"""Tests for TelegramNotifier and MockNotificationProvider.

Tests Telegram notification delivery with HTTP client integration
and mock provider for testing without external dependencies.
"""

import pytest
from unittest.mock import Mock, patch

from sdp.unified.notifications.telegram import (
    TelegramNotifier,
    TelegramConfig,
)
from sdp.unified.notifications.provider import (
    NotificationProvider,
    Notification,
    NotificationType,
)


class TestTelegramConfig:
    """Test TelegramConfig dataclass."""

    def test_creates_config_with_required_fields(self):
        """Should create config with bot token and chat ID."""
        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )

        assert config.bot_token == "test_token"
        assert config.chat_id == "test_chat"

    def test_creates_config_with_optional_fields(self):
        """Should create config with optional fields."""
        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
            parse_mode="Markdown",
            disable_notification=True,
        )

        assert config.parse_mode == "Markdown"
        assert config.disable_notification is True

    def test_config_defaults(self):
        """Should have correct default values."""
        config = TelegramConfig(
            bot_token="token",
            chat_id="chat",
        )

        assert config.parse_mode == "HTML"
        assert config.disable_notification is False


class TestTelegramNotifierInit:
    """Test TelegramNotifier initialization."""

    def test_creates_notifier(self):
        """Should initialize notifier."""
        config = TelegramConfig(
            bot_token="token",
            chat_id="chat",
        )
        notifier = TelegramNotifier(config)

        assert notifier is not None
        assert hasattr(notifier, 'send')

    def test_stores_config(self):
        """Should store configuration."""
        config = TelegramConfig(
            bot_token="token",
            chat_id="chat",
        )
        notifier = TelegramNotifier(config)

        assert notifier.config is config

    def test_is_notification_provider(self):
        """Should be NotificationProvider subclass."""
        config = TelegramConfig(
            bot_token="token",
            chat_id="chat",
        )
        notifier = TelegramNotifier(config)

        assert isinstance(notifier, NotificationProvider)


class TestTelegramNotificationSending:
    """Test Telegram notification sending."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_sends_info_notification(self, mock_post):
        """Should send info notification to Telegram."""
        mock_post.return_value = Mock(status_code=200, json=lambda: {"ok": True})

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="Info message",
        )
        result = notifier.send(notification)

        assert result is True
        mock_post.assert_called_once()

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_sends_error_notification(self, mock_post):
        """Should send error notification to Telegram."""
        mock_post.return_value = Mock(status_code=200, json=lambda: {"ok": True})

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.ERROR,
            message="Error message",
        )
        result = notifier.send(notification)

        assert result is True

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_formats_message_with_type_emoji(self, mock_post):
        """Should format message with type emoji."""
        mock_post.return_value = Mock(status_code=200, json=lambda: {"ok": True})

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.ERROR,
            message="Test",
        )
        notifier.send(notification)

        # Check call arguments
        call_args = mock_post.call_args
        message_text = call_args[1]["json"]["text"]
        assert "🚨" in message_text or "ERROR" in message_text

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_api_error(self, mock_post):
        """Should handle Telegram API errors."""
        mock_post.return_value = Mock(
            status_code=400,
            json=lambda: {"ok": False, "description": "Bad Request"}
        )

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_retries_on_failure(self, mock_post):
        """Should retry on transient failures."""
        # First call fails, second succeeds
        mock_post.side_effect = [
            Mock(status_code=500, json=lambda: {"ok": False}),
            Mock(status_code=200, json=lambda: {"ok": True}),
        ]

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = notifier.send(notification)

        assert result is True
        assert mock_post.call_count == 2


class TestMessageFormatting:
    """Test Telegram message formatting."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_formats_for_markdown_mode(self, mock_post):
        """Should format message for Markdown mode."""
        mock_post.return_value = Mock(status_code=200, json=lambda: {"ok": True})

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
            parse_mode="Markdown",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="*bold* message",
        )
        notifier.send(notification)

        call_args = mock_post.call_args
        json_data = call_args[1]["json"]
        assert json_data["parse_mode"] == "Markdown"

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_formats_for_html_mode(self, mock_post):
        """Should format message for HTML mode."""
        mock_post.return_value = Mock(status_code=200, json=lambda: {"ok": True})

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
            parse_mode="HTML",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.WARNING,
            message="<b>warning</b>",
        )
        notifier.send(notification)

        call_args = mock_post.call_args
        json_data = call_args[1]["json"]
        assert json_data["parse_mode"] == "HTML"


class TestMockNotificationProvider:
    """Test MockNotificationProvider for testing."""

    def test_creates_mock_provider(self):
        """Should create mock provider."""
        from sdp.unified.notifications.mock import MockNotificationProvider

        provider = MockNotificationProvider()

        assert provider is not None
        assert isinstance(provider, NotificationProvider)

    def test_mock_send_returns_true(self):
        """Should always return True from send."""
        from sdp.unified.notifications.mock import MockNotificationProvider

        provider = MockNotificationProvider()

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = provider.send(notification)

        assert result is True

    def test_mock_stores_notifications(self):
        """Should store sent notifications."""
        from sdp.unified.notifications.mock import MockNotificationProvider

        provider = MockNotificationProvider()

        notification1 = Notification(NotificationType.INFO, "Msg 1")
        notification2 = Notification(NotificationType.ERROR, "Msg 2")

        provider.send(notification1)
        provider.send(notification2)

        notifications = provider.get_notifications()

        assert len(notifications) == 2
        assert notifications[0].message == "Msg 1"

    def test_mock_clears_notifications(self):
        """Should clear stored notifications."""
        from sdp.unified.notifications.mock import MockNotificationProvider

        provider = MockNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "Test"))
        provider.clear()

        assert len(provider.get_notifications()) == 0

    def test_mock_resets(self):
        """Should reset all state."""
        from sdp.unified.notifications.mock import MockNotificationProvider

        provider = MockNotificationProvider()

        provider.send(Notification(NotificationType.INFO, "Test"))
        provider.reset()

        assert len(provider.get_notifications()) == 0


class TestErrorHandling:
    """Test error handling."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_network_error(self, mock_post):
        """Should handle network errors gracefully."""
        import requests

        mock_post.side_effect = requests.ConnectionError("Network error")

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_timeout(self, mock_post):
        """Should handle timeout errors."""
        import requests

        mock_post.side_effect = requests.Timeout("Request timed out")

        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat",
        )
        notifier = TelegramNotifier(config)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = notifier.send(notification)

        assert result is False
