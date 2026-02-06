"""E2E tests for Telegram notification integration.

Tests TelegramNotifier with both mock and real Telegram API,
covering configuration, message sending, formatting, and error handling.
"""

import os
from typing import TYPE_CHECKING
from unittest.mock import MagicMock, patch

import pytest
from requests.exceptions import ConnectionError as RequestsConnectionError
from requests.exceptions import Timeout
from sdp.unified.notifications.mock import MockNotificationProvider
from sdp.unified.notifications.provider import (
    Notification,
    NotificationType,
)
from sdp.unified.notifications.telegram import TelegramConfig, TelegramNotifier

if TYPE_CHECKING:
    pass


class TestTelegramConfigValidation:
    """Test Telegram configuration validation."""

    def test_create_config_with_required_fields(self) -> None:
        """Should create config with bot token and chat ID."""
        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat_id",
        )

        assert config.bot_token == "test_token"
        assert config.chat_id == "test_chat_id"
        assert config.parse_mode == "HTML"  # Default

    def test_create_config_with_all_fields(self) -> None:
        """Should create config with all optional fields."""
        config = TelegramConfig(
            bot_token="test_token",
            chat_id="test_chat_id",
            parse_mode="Markdown",
            disable_notification=True,
            timeout=15,
            max_retries=5,
        )

        assert config.parse_mode == "Markdown"
        assert config.disable_notification is True
        assert config.timeout == 15
        assert config.max_retries == 5


class TestMessageFormatting:
    """Test message formatting for Telegram."""

    def test_format_message_as_html(self) -> None:
        """Should format message with HTML (default)."""
        notifier = TelegramNotifier(
            config=TelegramConfig(
                bot_token="test",
                chat_id="test",
                parse_mode="HTML",
            ),
        )

        notification = Notification(
            type=NotificationType.ERROR,
            message="Error message",
        )

        formatted = notifier._format_message(notification)

        assert "<b>[ERROR]</b>" in formatted
        assert "Error message" in formatted

    def test_format_message_with_emoji(self) -> None:
        """Should add emoji based on notification type."""
        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        test_cases = [
            (NotificationType.INFO, "ℹ️"),
            (NotificationType.SUCCESS, "✅"),
            (NotificationType.WARNING, "⚠️"),
            (NotificationType.ERROR, "🚨"),
        ]

        for notif_type, emoji in test_cases:
            notification = Notification(
                type=notif_type,
                message="Test",
            )
            formatted = notifier._format_message(notification)
            assert emoji in formatted, f"Missing emoji for {notif_type}"

    def test_format_message_with_recipient(self) -> None:
        """Should include recipient when specified."""
        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Test message",
            recipient="user123",
        )

        formatted = notifier._format_message(notification)

        assert "To: user123" in formatted


class TestNotificationSending:
    """Test sending notifications with mock."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_info_notification(self, mock_post: MagicMock) -> None:
        """Should send info notification successfully."""
        mock_post.return_value = MagicMock(
            status_code=200,
            json=lambda: {"ok": True},
        )

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Info message",
        )

        result = notifier.send(notification)

        assert result is True
        assert mock_post.called

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_error_notification(self, mock_post: MagicMock) -> None:
        """Should send error notification successfully."""
        mock_post.return_value = MagicMock(
            status_code=200,
            json=lambda: {"ok": True},
        )

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.ERROR,
            message="Error message",
        )

        result = notifier.send(notification)

        assert result is True
        assert mock_post.called


class TestErrorHandling:
    """Test error handling in TelegramNotifier."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_api_error(self, mock_post: MagicMock) -> None:
        """Should handle Telegram API error."""
        mock_post.return_value = MagicMock(
            status_code=400,
            json=lambda: {"ok": False, "description": "Bad request"},
        )

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_connection_error(self, mock_post: MagicMock) -> None:
        """Should handle connection error."""
        mock_post.side_effect = RequestsConnectionError("Connection error")

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_handles_timeout(self, mock_post: MagicMock) -> None:
        """Should handle request timeout."""
        mock_post.side_effect = Timeout("Request timed out")

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test"),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_retries_on_failure(self, mock_post: MagicMock) -> None:
        """Should retry on transient failures."""
        # Fail twice, then succeed
        mock_post.side_effect = [
            RequestsConnectionError("Network error"),
            RequestsConnectionError("Network error"),
            MagicMock(status_code=200, json=lambda: {"ok": True}),
        ]

        notifier = TelegramNotifier(
            config=TelegramConfig(bot_token="test", chat_id="test", max_retries=3),
        )

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        result = notifier.send(notification)

        assert result is True
        assert mock_post.call_count == 3


class TestMockNotificationProvider:
    """Test MockNotificationProvider as Telegram alternative."""

    def test_mock_provider_sends_notification(self) -> None:
        """Should send notification to mock provider."""
        provider = MockNotificationProvider()

        notification = Notification(
            type=NotificationType.INFO,
            message="Test message",
        )

        result = provider.send(notification)

        assert result is True
        assert len(provider.get_notifications()) == 1

    def test_mock_provider_stores_notifications(self) -> None:
        """Should store all sent notifications."""
        provider = MockNotificationProvider()

        notifications = [
            Notification(type=NotificationType.INFO, message="Info 1"),
            Notification(type=NotificationType.ERROR, message="Error 1"),
            Notification(type=NotificationType.SUCCESS, message="Success 1"),
        ]

        for notif in notifications:
            provider.send(notif)

        stored = provider.get_notifications()
        assert len(stored) == 3
        assert stored[0].message == "Info 1"
        assert stored[1].message == "Error 1"
        assert stored[2].message == "Success 1"

    def test_mock_provider_filters_by_type(self) -> None:
        """Should filter notifications by type."""
        provider = MockNotificationProvider()

        provider.send(Notification(type=NotificationType.INFO, message="Info"))
        provider.send(Notification(type=NotificationType.ERROR, message="Error"))
        provider.send(Notification(type=NotificationType.INFO, message="Info 2"))

        errors = provider.get_by_type(NotificationType.ERROR)
        infos = provider.get_by_type(NotificationType.INFO)

        assert len(errors) == 1
        assert len(infos) == 2

    def test_mock_provider_clears_notifications(self) -> None:
        """Should clear all notifications."""
        provider = MockNotificationProvider()

        provider.send(Notification(type=NotificationType.INFO, message="Test"))
        assert len(provider.get_notifications()) == 1

        provider.clear()
        assert len(provider.get_notifications()) == 0

    def test_mock_provider_counts_by_type(self) -> None:
        """Should count all notifications."""
        provider = MockNotificationProvider()

        provider.send(Notification(type=NotificationType.INFO, message="Info"))
        provider.send(Notification(type=NotificationType.ERROR, message="Error"))
        provider.send(Notification(type=NotificationType.INFO, message="Info 2"))

        # Count returns total count
        assert provider.count() == 3

        # Use get_by_type for type-specific counts
        assert len(provider.get_by_type(NotificationType.INFO)) == 2
        assert len(provider.get_by_type(NotificationType.ERROR)) == 1


class TestRealTelegramIntegration:
    """Integration tests with real Telegram API."""

    @pytest.mark.skipif(
        os.getenv("TELEGRAM_BOT_TOKEN") is None,
        reason="Skipping real Telegram tests (TELEGRAM_BOT_TOKEN not set)"
    )
    def test_real_telegram_send_notification(self) -> None:
        """Test sending notification with real Telegram API."""
        try:
            config = TelegramConfig(
                bot_token=os.getenv("TELEGRAM_BOT_TOKEN", ""),
                chat_id=os.getenv("TELEGRAM_CHAT_ID", ""),
            )
            notifier = TelegramNotifier(config=config)

            notification = Notification(
                type=NotificationType.INFO,
                message="E2E test message from SDP",
            )

            result = notifier.send(notification)

            assert result is True

        except Exception as e:
            pytest.fail(f"Failed to send notification: {e}")

    @pytest.mark.skipif(
        os.getenv("TELEGRAM_BOT_TOKEN") is None,
        reason="Skipping real Telegram tests (TELEGRAM_BOT_TOKEN not set)"
    )
    def test_real_telegram_send_all_types(self) -> None:
        """Test sending all notification types with real Telegram."""
        try:
            config = TelegramConfig(
                bot_token=os.getenv("TELEGRAM_BOT_TOKEN", ""),
                chat_id=os.getenv("TELEGRAM_CHAT_ID", ""),
            )
            notifier = TelegramNotifier(config=config)

            types = [
                NotificationType.INFO,
                NotificationType.SUCCESS,
                NotificationType.WARNING,
                NotificationType.ERROR,
            ]

            for notif_type in types:
                notification = Notification(
                    type=notif_type,
                    message=f"Test {notif_type.value} message",
                )
                result = notifier.send(notification)
                assert result is True, f"Failed to send {notif_type}"

        except Exception as e:
            pytest.fail(f"Failed to send notifications: {e}")


class TestTelegramNotifierFixtures:
    """Test fixtures for Telegram integration."""

    @pytest.fixture
    def telegram_notifier(self) -> TelegramNotifier:
        """Provide Telegram notifier for testing."""
        return TelegramNotifier(
            config=TelegramConfig(
                bot_token="test_token",
                chat_id="test_chat",
            ),
        )

    @pytest.fixture
    def mock_provider(self) -> MockNotificationProvider:
        """Provide mock provider for testing."""
        return MockNotificationProvider()

    def test_fixture_usage(
        self,
        telegram_notifier: TelegramNotifier,
        mock_provider: MockNotificationProvider,
    ) -> None:
        """Test fixture usage in tests."""
        assert telegram_notifier is not None
        assert mock_provider is not None

        # Send to mock provider
        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        result = mock_provider.send(notification)
        assert result is True
