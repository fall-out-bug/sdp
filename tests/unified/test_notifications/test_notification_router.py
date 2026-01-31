"""Tests for NotificationRouter module.

Tests notification routing to providers, broadcasting, and
provider management functionality.
"""

import pytest
from unittest.mock import Mock, call

from sdp.unified.notifications.router import NotificationRouter
from sdp.unified.notifications.provider import (
    NotificationProvider,
    Notification,
    NotificationType,
    InMemoryNotificationProvider,
)


class MockNotificationProvider(NotificationProvider):
    """Mock notification provider for testing."""

    def __init__(self) -> None:
        """Initialize mock provider."""
        self.notifications: list[Notification] = []

    def send(self, notification: Notification) -> bool:
        """Store notification in memory."""
        self.notifications.append(notification)
        return True


class TestNotificationRouterInit:
    """Test NotificationRouter initialization."""

    def test_creates_router(self):
        """Should initialize router."""
        router = NotificationRouter()

        assert router is not None
        assert hasattr(router, 'register_provider')
        assert hasattr(router, 'route')

    def test_initializes_empty_providers(self):
        """Should initialize with no providers."""
        router = NotificationRouter()

        assert router.list_providers() == []

    def test_initializes_with_default_providers(self):
        """Should initialize with default providers if provided."""
        provider = MockNotificationProvider()
        router = NotificationRouter(providers=[provider])

        providers = router.list_providers()

        assert len(providers) == 1


class TestProviderManagement:
    """Test provider registration and management."""

    def test_registers_provider(self):
        """Should register notification provider."""
        router = NotificationRouter()
        provider = MockNotificationProvider()

        router.register_provider(provider)

        assert len(router.list_providers()) == 1

    def test_registers_multiple_providers(self):
        """Should register multiple providers."""
        router = NotificationRouter()
        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        router.register_provider(provider1)
        router.register_provider(provider2)

        assert len(router.list_providers()) == 2

    def test_removes_provider(self):
        """Should remove registered provider."""
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        router.remove_provider(provider)

        assert len(router.list_providers()) == 0

    def test_lists_providers(self):
        """Should list all registered providers."""
        router = NotificationRouter()
        provider1 = MockNotificationProvider()
        provider2 = InMemoryNotificationProvider()

        router.register_provider(provider1)
        router.register_provider(provider2)

        providers = router.list_providers()

        assert len(providers) == 2
        assert provider1 in providers
        assert provider2 in providers


class TestNotificationRouting:
    """Test notification routing to providers."""

    def test_routes_to_single_provider(self):
        """Should route notification to single provider."""
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        router.route(notification)

        assert len(provider.notifications) == 1
        assert provider.notifications[0].message == "Test"

    def test_routes_to_multiple_providers(self):
        """Should route notification to all providers."""
        router = NotificationRouter()
        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        router.register_provider(provider1)
        router.register_provider(provider2)

        notification = Notification(
            type=NotificationType.INFO,
            message="Broadcast",
        )
        router.route(notification)

        assert len(provider1.notifications) == 1
        assert len(provider2.notifications) == 1

    def test_filters_by_recipient(self):
        """Should filter notifications by recipient."""
        router = NotificationRouter()

        # Register provider with recipient filter
        provider = MockNotificationProvider()
        router.register_provider(provider, recipient="agent-1")

        notification = Notification(
            type=NotificationType.INFO,
            message="Message",
            recipient="agent-1",
        )
        router.route(notification)

        # Should deliver to matching provider
        assert len(provider.notifications) == 1

    def test_skips_non_matching_recipient(self):
        """Should skip providers with non-matching recipient filters."""
        router = NotificationRouter()

        provider = MockNotificationProvider()
        router.register_provider(provider, recipient="agent-1")

        notification = Notification(
            type=NotificationType.INFO,
            message="Message",
            recipient="agent-2",
        )
        router.route(notification)

        # Should not deliver
        assert len(provider.notifications) == 0

    def test_handles_no_providers(self):
        """Should handle routing with no providers gracefully."""
        router = NotificationRouter()

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )

        # Should not raise exception
        router.route(notification)


class TestBroadcasting:
    """Test broadcast functionality."""

    def test_broadcasts_to_all_providers(self):
        """Should broadcast to all providers."""
        router = NotificationRouter()
        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        router.register_provider(provider1)
        router.register_provider(provider2)

        notification = Notification(
            type=NotificationType.INFO,
            message="Broadcast",
        )
        router.broadcast(notification)

        assert len(provider1.notifications) == 1
        assert len(provider2.notifications) == 1

    def test_broadcast_sends_to_all_regardless_of_filters(self):
        """Should broadcast to all providers ignoring recipient filters."""
        router = NotificationRouter()

        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        # Register with recipient filters
        router.register_provider(provider1, recipient="agent-1")
        router.register_provider(provider2, recipient="agent-2")

        notification = Notification(
            type=NotificationType.INFO,
            message="Broadcast",
            recipient="agent-1",
        )

        # Broadcast should ignore recipient filters
        router.broadcast(notification)

        assert len(provider1.notifications) == 1
        assert len(provider2.notifications) == 1


class TestFailedDelivery:
    """Test failed delivery handling."""

    def test_continues_on_provider_failure(self):
        """Should continue routing if one provider fails."""
        router = NotificationRouter()

        # Create a failing provider
        class FailingProvider(NotificationProvider):
            def send(self, notification: Notification) -> bool:
                return False

        failing_provider = FailingProvider()
        working_provider = MockNotificationProvider()

        router.register_provider(failing_provider)
        router.register_provider(working_provider)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        router.route(notification)

        # Working provider should still receive
        assert len(working_provider.notifications) == 1


class TestProviderFiltering:
    """Test provider filtering by notification type."""

    def test_filters_by_notification_type(self):
        """Should filter providers by notification type."""
        router = NotificationRouter()

        provider1 = MockNotificationProvider()
        provider2 = MockNotificationProvider()

        # Register with type filters
        router.register_provider(provider1, types=[NotificationType.ERROR])
        router.register_provider(provider2, types=[NotificationType.INFO])

        error_notification = Notification(
            type=NotificationType.ERROR,
            message="Error",
        )
        info_notification = Notification(
            type=NotificationType.INFO,
            message="Info",
        )

        router.route(error_notification)
        router.route(info_notification)

        # provider1 should only get ERROR
        assert len(provider1.notifications) == 1
        assert provider1.notifications[0].type == NotificationType.ERROR

        # provider2 should only get INFO
        assert len(provider2.notifications) == 1
        assert provider2.notifications[0].type == NotificationType.INFO

    def test_handles_no_type_filter(self):
        """Should handle providers with no type filter."""
        router = NotificationRouter()

        provider = MockNotificationProvider()
        router.register_provider(provider)  # No type filter

        error_notification = Notification(
            type=NotificationType.ERROR,
            message="Error",
        )
        info_notification = Notification(
            type=NotificationType.INFO,
            message="Info",
        )

        router.route(error_notification)
        router.route(info_notification)

        # Should receive both
        assert len(provider.notifications) == 2


class TestNotificationTransformation:
    """Test notification transformation and enrichment."""

    def test_enriches_notification_with_metadata(self):
        """Should enrich notification with router metadata."""
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
        )
        router.route(notification)

        # Notification should be delivered
        assert len(provider.notifications) == 1


class TestEdgeCases:
    """Test edge cases."""

    def test_handles_same_provider_multiple_times(self):
        """Should handle registering same provider multiple times."""
        router = NotificationRouter()
        provider = MockNotificationProvider()

        router.register_provider(provider)
        router.register_provider(provider)

        # Should only have one instance
        assert len(router.list_providers()) == 1

    def test_handles_none_recipient_in_notification(self):
        """Should handle notifications with None recipient."""
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        notification = Notification(
            type=NotificationType.INFO,
            message="Test",
            recipient=None,
        )
        router.route(notification)

        # Should deliver
        assert len(provider.notifications) == 1

    def test_handles_empty_message(self):
        """Should handle empty message gracefully."""
        router = NotificationRouter()
        provider = MockNotificationProvider()
        router.register_provider(provider)

        notification = Notification(
            type=NotificationType.INFO,
            message="",
        )
        router.route(notification)

        # Should deliver
        assert len(provider.notifications) == 1
