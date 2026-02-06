"""Retry policy and result models for builder routing."""

from dataclasses import dataclass
from typing import Optional

from sdp.errors import ErrorCategory, SDPError


@dataclass
class RetryPolicy:
    """Retry policy configuration."""

    max_attempts: int = 3

    def should_retry(self, attempt: int) -> bool:
        """Check if retry is allowed for given attempt number.

        Args:
            attempt: Current attempt number (1-indexed)

        Returns:
            True if retry allowed, False otherwise
        """
        return attempt <= self.max_attempts

    def should_escalate(self, attempt: int, failed: bool) -> bool:
        """Check if should escalate to human after failed attempts.

        Args:
            attempt: Current attempt number (1-indexed)
            failed: Whether the attempt failed

        Returns:
            True if should escalate to human
        """
        return failed and attempt > self.max_attempts


@dataclass
class BuildResult:
    """Result of build execution."""

    success: bool
    attempt: int
    output: str = ""
    error: Optional[str] = None
    diagnostics: Optional[str] = None


class HumanEscalationError(SDPError):
    """Error raised when build should be escalated to human."""

    def __init__(
        self,
        ws_id: str,
        tier: Optional[str],
        attempts: int,
        diagnostics: str,
    ) -> None:
        """Initialize escalation error."""
        super().__init__(
            category=ErrorCategory.BUILD,
            message=(
                f"Build failed after {attempts} attempts for {ws_id} (tier: {tier}). "
                "Escalating to human."
            ),
            remediation=(
                "1. Review diagnostics below for failure cause\n"
                "2. Fix the underlying issue (tests, implementation)\n"
                "3. Retry build manually: sdp build {ws_id}\n"
                "4. Escalate to human if issue persists"
            ),
            docs_url="https://docs.sdp.dev/build#escalation",
            context={
                "ws_id": ws_id,
                "tier": tier,
                "attempts": attempts,
                "diagnostics": diagnostics,
            },
        )
