"""Mock notification provider for testing.

Provides in-memory notification storage for testing without
external dependencies.
"""

import logging
from typing import List

from sdp.unified.notifications.provider import (
    Notification,
    NotificationProvider,
    NotificationType,
)

logger = logging.getLogger(__name__)


class MockNotificationProvider(NotificationProvider):
    """Mock notification provider for testing.

    Stores notifications in memory for inspection during tests.
    """

    def __init__(self) -> None:
        """Initialize mock provider."""
        self._notifications: List[Notification] = []

    def send(self, notification: Notification) -> bool:
        """Store notification in memory.

        Args:
            notification: Notification to store

        Returns:
            Always returns True
        """
        self._notifications.append(notification)
        return True

    def get_notifications(self) -> List[Notification]:
        """Get all stored notifications.

        Returns:
            List of stored notifications
        """
        return self._notifications.copy()

    def clear(self) -> None:
        """Clear all stored notifications."""
        self._notifications.clear()

    def reset(self) -> None:
        """Reset provider to initial state."""
        self._notifications.clear()

    def count(self) -> int:
        """Get count of stored notifications.

        Returns:
            Number of notifications
        """
        return len(self._notifications)

    def get_by_type(self, notification_type: NotificationType) -> List[Notification]:
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

    def get_by_recipient(self, recipient: str) -> List[Notification]:
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
