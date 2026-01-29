"""Telegram notification provider and configuration.

Provides TelegramNotifier for sending notifications via Telegram Bot API
with retry logic and error handling.
"""

import logging
from dataclasses import dataclass

import requests

from sdp.unified.notifications.provider import (
    Notification,
    NotificationProvider,
    NotificationType,
)

logger = logging.getLogger(__name__)

TELEGRAM_API_URL = "https://api.telegram.org/bot"


@dataclass
class TelegramConfig:
    """Telegram Bot API configuration.

    Attributes:
        bot_token: Telegram bot token
        chat_id: Target chat ID
        parse_mode: Message parse mode (HTML or Markdown)
        disable_notification: Send silently without notification
        timeout: Request timeout in seconds
        max_retries: Maximum retry attempts
    """
    bot_token: str
    chat_id: str
    parse_mode: str = "HTML"
    disable_notification: bool = False
    timeout: int = 10
    max_retries: int = 3


class TelegramNotifier(NotificationProvider):
    """Send notifications to Telegram via Bot API.

    Provides Telegram Bot API integration with message formatting,
    retry logic, and error handling.
    """

    # Emoji mapping for notification types
    TYPE_EMOJI = {
        NotificationType.INFO: "ℹ️",
        NotificationType.WARNING: "⚠️",
        NotificationType.ERROR: "🚨",
        NotificationType.SUCCESS: "✅",
    }

    def __init__(self, config: TelegramConfig) -> None:
        """Initialize Telegram notifier.

        Args:
            config: Telegram configuration
        """
        self.config = config
        self._api_url = f"{TELEGRAM_API_URL}{config.bot_token}/sendMessage"

    def send(self, notification: Notification) -> bool:
        """Send notification to Telegram.

        Args:
            notification: Notification to send

        Returns:
            True if sent successfully, False otherwise
        """
        # Format message
        message = self._format_message(notification)

        # Prepare payload
        payload = {
            "chat_id": self.config.chat_id,
            "text": message,
            "parse_mode": self.config.parse_mode,
            "disable_notification": self.config.disable_notification,
        }

        # Send with retry logic
        for attempt in range(self.config.max_retries):
            try:
                response = requests.post(
                    self._api_url,
                    json=payload,
                    timeout=self.config.timeout,
                )

                if response.status_code == 200:
                    data = response.json()
                    if data.get("ok"):
                        logger.info(f"Sent Telegram notification: {notification.message[:50]}")
                        return True
                    else:
                        logger.error(
                            f"Telegram API error: {data.get('description')}"
                        )
                        return False
                else:
                    logger.error(
                        f"Telegram HTTP error: {response.status_code}"
                    )
                    # Retry on server errors
                    if response.status_code >= 500 and attempt < self.config.max_retries - 1:
                        continue
                    return False

            except requests.ConnectionError as e:
                logger.error(f"Telegram connection error: {e}")
                if attempt < self.config.max_retries - 1:
                    continue
                return False

            except requests.Timeout as e:
                logger.error(f"Telegram timeout error: {e}")
                if attempt < self.config.max_retries - 1:
                    continue
                return False

            except Exception as e:
                logger.error(f"Unexpected error sending to Telegram: {e}")
                return False

        return False

    def _format_message(self, notification: Notification) -> str:
        """Format notification for Telegram.

        Args:
            notification: Notification to format

        Returns:
            Formatted message string
        """
        # Get emoji for type
        emoji = self.TYPE_EMOJI.get(notification.type, "")

        # Build message parts
        parts = []

        if emoji:
            parts.append(f"{emoji}")

        # Add type label
        parts.append(f"<b>[{notification.type.value.upper()}]</b>")

        # Add message
        parts.append(notification.message)

        # Add recipient if present
        if notification.recipient:
            parts.append(f"\n\nTo: {notification.recipient}")

        return " ".join(parts)
