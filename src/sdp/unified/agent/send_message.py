"""SendMessage tool wrapper for inter-agent communication.

Provides wrapper interface for Claude Code's SendMessage tool to enable
mocking and testing of message routing functionality.
"""

from typing import Optional


class SendMessageResult:
    """Result from SendMessage tool invocation."""

    def __init__(
        self,
        success: bool,
        error: Optional[str] = None,
    ) -> None:
        """Initialize SendMessage result.

        Args:
            success: Whether message was sent successfully
            error: Optional error message if failed
        """
        self.success = success
        self.error = error


def SendMessage(  # noqa: N802
    recipient: str,
    content: str,
) -> SendMessageResult:
    """Invoke SendMessage tool to send message to agent.

    Wrapper function for actual SendMessage tool. In tests, mocked to
    return SendMessageResult objects.

    Args:
        recipient: Recipient agent ID or name
        content: Message content

    Returns:
        SendMessageResult with success status
    """
    # Placeholder - in production calls actual SendMessage tool
    return SendMessageResult(success=True)
