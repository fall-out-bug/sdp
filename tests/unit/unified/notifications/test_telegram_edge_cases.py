"""Tests for Telegram notifier edge cases."""

import pytest
from unittest.mock import Mock, patch
from requests.exceptions import ConnectionError, Timeout

from sdp.unified.notifications.telegram import (
    TelegramNotifier,
    TelegramConfig,
)
from sdp.unified.notifications.provider import Notification, NotificationType


@pytest.fixture
def config():
    """Create Telegram config."""
    return TelegramConfig(
        bot_token="test_token",
        chat_id="123456",
        timeout=5,
        max_retries=3,
    )


@pytest.fixture
def notifier(config):
    """Create Telegram notifier."""
    return TelegramNotifier(config)


@pytest.fixture
def notification():
    """Create test notification."""
    return Notification(
        type=NotificationType.INFO,
        message="Test message",
    )


class TestTelegramNotifier:
    """Tests for TelegramNotifier."""

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_success(self, mock_post, notifier, notification):
        """Test sends notification successfully."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"ok": True}
        mock_post.return_value = mock_response

        result = notifier.send(notification)

        assert result is True
        mock_post.assert_called_once()

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_api_error(self, mock_post, notifier, notification):
        """Test handles Telegram API error."""
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "ok": False,
            "description": "Bad request"
        }
        mock_post.return_value = mock_response

        result = notifier.send(notification)

        assert result is False

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_http_error_no_retry(self, mock_post, notifier, notification):
        """Test handles HTTP client error (no retry)."""
        mock_response = Mock()
        mock_response.status_code = 400
        mock_post.return_value = mock_response

        result = notifier.send(notification)

        assert result is False
        assert mock_post.call_count == 1  # No retry

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_http_error_with_retry(self, mock_post, notifier, notification):
        """Test retries on HTTP server error."""
        mock_response = Mock()
        mock_response.status_code = 500
        mock_post.return_value = mock_response

        result = notifier.send(notification)

        assert result is False
        assert mock_post.call_count == 3  # max_retries

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_connection_error_retries(self, mock_post, notifier, notification):
        """Test retries on connection error."""
        mock_post.side_effect = ConnectionError("Network unreachable")

        result = notifier.send(notification)

        assert result is False
        assert mock_post.call_count == 3  # max_retries

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_timeout_error_retries(self, mock_post, notifier, notification):
        """Test retries on timeout error."""
        mock_post.side_effect = Timeout("Request timeout")

        result = notifier.send(notification)

        assert result is False
        assert mock_post.call_count == 3  # max_retries

    @patch('sdp.unified.notifications.telegram.requests.post')
    def test_send_unexpected_exception(self, mock_post, notifier, notification):
        """Test handles unexpected exception."""
        mock_post.side_effect = RuntimeError("Unexpected error")

        result = notifier.send(notification)

        assert result is False
        assert mock_post.call_count == 1  # No retry on unexpected error

    def test_format_message_with_emoji(self, notifier):
        """Test formats message with emoji."""
        notification = Notification(
            type=NotificationType.SUCCESS,
            message="Build completed",
        )

        formatted = notifier._format_message(notification)

        assert "✅" in formatted
        assert "[SUCCESS]" in formatted
        assert "Build completed" in formatted

    def test_format_message_with_recipient(self, notifier):
        """Test formats message with recipient."""
        notification = Notification(
            type=NotificationType.WARNING,
            message="Approval needed",
            recipient="@developer",
        )

        formatted = notifier._format_message(notification)

        assert "⚠️" in formatted
        assert "[WARNING]" in formatted
        assert "Approval needed" in formatted
        assert "To: @developer" in formatted

    def test_format_message_unknown_type(self, notifier):
        """Test formats message with unknown type (no emoji)."""
        notification = Notification(
            type=NotificationType.INFO,
            message="Info message",
        )

        formatted = notifier._format_message(notification)

        assert "ℹ️" in formatted
        assert "[INFO]" in formatted
