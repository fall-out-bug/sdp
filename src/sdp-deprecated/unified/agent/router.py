"""SendMessage router for inter-agent communication.

Routes messages between agents using SendMessage tool with support for
direct messages, broadcasts, and offline agent queuing.
"""

import logging
from dataclasses import dataclass, field
from datetime import datetime, timezone
from typing import Any, Optional

from sdp.unified.agent.send_message import SendMessage

logger = logging.getLogger(__name__)

# Type hint for TeamManager (will be imported when available)
TeamManager = Any


@dataclass
class Message:
    """Message sent between agents.

    Attributes:
        sender: Sender agent ID
        content: Message content
        recipient: Optional recipient agent ID (None for broadcast)
        broadcast: Whether this is a broadcast message
        message_type: Message type ('direct' or 'broadcast')
        timestamp: Message timestamp
    """
    sender: str
    content: str
    recipient: Optional[str] = None
    broadcast: bool = False
    message_type: str = field(init=False)
    timestamp: datetime = field(
        default_factory=lambda: datetime.now(timezone.utc)
    )

    def __post_init__(self) -> None:
        """Set message_type based on broadcast flag."""
        self.message_type = "broadcast" if self.broadcast else "direct"


@dataclass
class DeliveryResult:
    """Result from message delivery attempt.

    Attributes:
        success: Whether delivery succeeded
        error: Optional error message if failed
    """
    success: bool
    error: Optional[str] = None


class SendMessageRouter:
    """Router for inter-agent messages using SendMessage tool."""

    def __init__(self, team_manager: Optional[Any] = None) -> None:
        """Initialize router with optional TeamManager for role lookup."""
        self._team_manager = team_manager
        self._pending_messages: dict[str, list[Message]] = {}

    def send_message(self, message: Message) -> DeliveryResult:
        """Send message to recipient(s). Routes based on broadcast flag."""
        if message.broadcast:
            return self._send_broadcast(message)
        else:
            return self._send_direct(message)

    def _send_direct(self, message: Message) -> DeliveryResult:
        """Send direct message to recipient after resolving name via TeamManager."""
        if not message.recipient:
            return DeliveryResult(
                success=False,
                error="No recipient specified"
            )

        # Look up recipient agent ID via TeamManager if available
        recipient_id = message.recipient
        if self._team_manager and hasattr(self._team_manager, 'find_agent_by_name'):
            resolved_id = self._team_manager.find_agent_by_name(message.recipient)
            if resolved_id:
                recipient_id = resolved_id
            else:
                return DeliveryResult(
                    success=False,
                    error=f"Recipient '{message.recipient}' not found in team"
                )

        # Use SendMessage tool to send message
        result = SendMessage(
            recipient=recipient_id,
            content=message.content,
        )

        if result.success:
            logger.info(
                f"Message sent from {message.sender} to {recipient_id}"
            )
            return DeliveryResult(success=True)
        else:
            logger.error(
                f"Failed to send message from {message.sender} to {recipient_id}: "
                f"{result.error}"
            )
            return DeliveryResult(
                success=False,
                error=result.error or "Delivery failed"
            )

    def _send_broadcast(self, message: Message) -> DeliveryResult:
        """Send broadcast to all team members except sender."""
        if not self._team_manager:
            return DeliveryResult(
                success=False,
                error="No team manager available for broadcast"
            )

        # Get all active agents from TeamManager
        if hasattr(self._team_manager, 'get_active_roles'):
            active_agents = self._team_manager.get_active_roles()
        else:
            active_agents = []

        if not active_agents:
            return DeliveryResult(
                success=False,
                error="No active agents found"
            )

        # Send to all active agents except sender
        sent_count = 0
        for agent_id in active_agents:
            if agent_id != message.sender:
                result = SendMessage(
                    recipient=agent_id,
                    content=message.content,
                )
                if result.success:
                    logger.info(f"Broadcast from {message.sender} to {agent_id}")
                    sent_count += 1

        return DeliveryResult(
            success=sent_count > 0,
            error=f"Sent to {sent_count} recipients"
        )

    def queue_message(self, message: Message) -> None:
        """Queue message for offline recipient."""
        if not message.recipient:
            logger.warning("Cannot queue message without recipient")
            return

        recipient = message.recipient
        if recipient not in self._pending_messages:
            self._pending_messages[recipient] = []
        self._pending_messages[recipient].append(message)
        logger.info(f"Queued message for {recipient}")

    def get_pending_messages(self, recipient: str) -> list[Message]:
        """Get pending messages for recipient."""
        return self._pending_messages.get(recipient, [])

    def deliver_pending(self, recipient: str) -> list[Message]:
        """Deliver all pending messages to recipient and clear queue."""
        pending = self._pending_messages.get(recipient, [])
        if pending:
            self._pending_messages[recipient] = []
            logger.info(f"Delivered {len(pending)} pending messages to {recipient}")
        return pending
