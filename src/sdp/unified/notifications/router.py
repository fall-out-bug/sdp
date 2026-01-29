"""NotificationRouter for routing notifications to providers.

Routes notifications to registered providers with support for
filtering by type, recipient, and broadcasting.
"""

import logging
from dataclasses import dataclass
from typing import List, Optional

from sdp.unified.notifications.provider import (
    Notification,
    NotificationProvider,
    NotificationType,
)

logger = logging.getLogger(__name__)


@dataclass
class ProviderRegistration:
    """Provider registration with filters.

    Attributes:
        provider: Registered provider
        recipient: Optional recipient filter
        types: Optional list of notification types to handle
    """
    provider: NotificationProvider
    recipient: Optional[str] = None
    types: Optional[List[NotificationType]] = None


class NotificationRouter:
    """Router for notification delivery to providers.

    Manages provider registration and routes notifications based
    on type and recipient filters.
    """

    def __init__(self, providers: Optional[List[NotificationProvider]] = None) -> None:
        """Initialize notification router.

        Args:
            providers: Optional list of initial providers
        """
        self._providers: List[ProviderRegistration] = []

        if providers:
            for provider in providers:
                self.register_provider(provider)

    def register_provider(
        self,
        provider: NotificationProvider,
        recipient: Optional[str] = None,
        types: Optional[List[NotificationType]] = None,
    ) -> None:
        """Register notification provider.

        Args:
            provider: Provider to register
            recipient: Optional recipient filter
            types: Optional notification type filter
        """
        # Check if already registered
        for registration in self._providers:
            if registration.provider is provider:
                logger.warning("Provider already registered, skipping")
                return

        registration = ProviderRegistration(
            provider=provider,
            recipient=recipient,
            types=types,
        )
        self._providers.append(registration)
        logger.info(f"Registered provider: {provider.__class__.__name__}")

    def remove_provider(self, provider: NotificationProvider) -> None:
        """Remove registered provider.

        Args:
            provider: Provider to remove
        """
        self._providers = [
            r for r in self._providers
            if r.provider is not provider
        ]
        logger.info(f"Removed provider: {provider.__class__.__name__}")

    def list_providers(self) -> List[NotificationProvider]:
        """List all registered providers.

        Returns:
            List of registered providers
        """
        return [r.provider for r in self._providers]

    def route(self, notification: Notification) -> bool:
        """Route notification to matching providers.

        Routes to providers that match notification type and recipient.

        Args:
            notification: Notification to route

        Returns:
            True if delivered to at least one provider
        """
        delivered = False

        for registration in self._providers:
            # Check type filter
            if registration.types and notification.type not in registration.types:
                continue

            # Check recipient filter
            if (
                registration.recipient is not None
                and notification.recipient != registration.recipient
            ):
                continue

            # Deliver to provider
            try:
                success = registration.provider.send(notification)
                if success:
                    delivered = True
            except Exception as e:
                logger.error(
                    f"Failed to deliver notification to "
                    f"{registration.provider.__class__.__name__}: {e}"
                )

        return delivered

    def broadcast(self, notification: Notification) -> bool:
        """Broadcast notification to all providers.

        Ignores type and recipient filters, sends to all providers.

        Args:
            notification: Notification to broadcast

        Returns:
            True if delivered to at least one provider
        """
        delivered = False

        for registration in self._providers:
            try:
                success = registration.provider.send(notification)
                if success:
                    delivered = True
            except Exception as e:
                logger.error(
                    f"Failed to broadcast to "
                    f"{registration.provider.__class__.__name__}: {e}"
                )

        return delivered
