"""NotificationProvider interface and implementations.

Provides abstract interface for notification delivery with console
and in-memory implementations for production and testing.
"""

import logging
import sys
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime, timezone
from enum import Enum
from typing import IO, Optional

logger = logging.getLogger(__name__)


class NotificationType(Enum):
    """Notification type levels."""

    INFO = "info"
    WARNING = "warning"
    ERROR = "error"
    SUCCESS = "success"


@dataclass
class Notification:
    """Notification message.

    Attributes:
        type: Notification type
        message: Notification message
        recipient: Optional recipient (agent ID or user)
        timestamp: Notification timestamp
    """
    type: NotificationType
    message: str
    recipient: Optional[str] = None
    timestamp: datetime = field(
        default_factory=lambda: datetime.now(timezone.utc)
    )


class NotificationProvider(ABC):
    """Abstract notification provider interface.

    Defines interface for sending notifications to various
    destinations (console, in-memory, external services).
    """

    @abstractmethod
    def send(self, notification: Notification) -> bool:
        """Send notification.

        Args:
            notification: Notification to send

        Returns:
            True if sent successfully, False otherwise
        """
        pass


class ConsoleNotificationProvider(NotificationProvider):
    """Console notification provider.

    Outputs notifications to stdout/stderr with formatting
    and timestamps.
    """

    def __init__(
        self,
        stdout: Optional[IO[str]] = None,
        stderr: Optional[IO[str]] = None,
    ) -> None:
        """Initialize console provider with stdout/stderr streams."""
        self._stdout = stdout or sys.stdout
        self._stderr = stderr or sys.stderr

    def send(self, notification: Notification) -> bool:
        """Send notification to console.

        Args:
            notification: Notification to send

        Returns:
            True if sent successfully
        """
        try:
            # Format timestamp
            ts = notification.timestamp.strftime("%Y-%m-%d %H:%M:%S")

            # Build output
            parts = [f"[{ts}]", f"[{notification.type.value.upper()}]"]

            if notification.recipient:
                parts.append(f"[To: {notification.recipient}]")

            parts.append(notification.message)

            output = " ".join(parts)

            # Route to stderr for errors, stdout for others
            if notification.type == NotificationType.ERROR:
                self._stderr.write(output + "\n")
            else:
                self._stdout.write(output + "\n")

            return True
        except Exception as e:
            logger.error(f"Failed to send notification: {e}")
            return False


class InMemoryNotificationProvider(NotificationProvider):
    """In-memory notification provider for testing.

    Stores notifications in memory for retrieval and inspection.
    """

    def __init__(self) -> None:
        """Initialize in-memory provider."""
        self._notifications: list[Notification] = []

    def send(self, notification: Notification) -> bool:
        """Store notification in memory.

        Args:
            notification: Notification to store

        Returns:
            True if stored successfully
        """
        self._notifications.append(notification)
        return True

    def get_notifications(self) -> list[Notification]:
        """Get all stored notifications.

        Returns:
            List of all notifications
        """
        return self._notifications.copy()

    def get_notifications_by_type(self, notification_type: NotificationType) -> list[Notification]:
        """Get notifications by type.

        Args:
            notification_type: Type to filter by

        Returns:
            List of matching notifications
        """
        return [
            n for n in self._notifications
            if n.type == notification_type
        ]

    def get_notifications_by_recipient(self, recipient: str) -> list[Notification]:
        """Get notifications by recipient.

        Args:
            recipient: Recipient to filter by

        Returns:
            List of matching notifications
        """
        return [
            n for n in self._notifications
            if n.recipient == recipient
        ]

    def clear(self) -> None:
        """Clear all stored notifications."""
        self._notifications.clear()

    def count(self) -> int:
        """Get count of stored notifications.

        Returns:
            Number of notifications
        """
        return len(self._notifications)
