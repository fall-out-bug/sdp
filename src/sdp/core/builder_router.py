"""Builder router for model-agnostic workstream execution.

Routes /build commands to appropriate model providers based on capability tier
and applies retry/escalation policies.
"""

from dataclasses import dataclass
from typing import Optional

from sdp.core.model_mapping import ModelProvider, ModelRegistry
from sdp.core.workstream import Workstream


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


class HumanEscalationError(Exception):
    """Error raised when build should be escalated to human."""

    def __init__(
        self,
        ws_id: str,
        tier: Optional[str],
        attempts: int,
        diagnostics: str,
    ) -> None:
        """Initialize escalation error.

        Args:
            ws_id: Workstream ID
            tier: Capability tier
            attempts: Number of failed attempts
            diagnostics: Diagnostic information for human
        """
        self.ws_id = ws_id
        self.tier = tier
        self.attempts = attempts
        self.diagnostics = diagnostics
        super().__init__(
            f"Build failed after {attempts} attempts for {ws_id} (tier: {tier}). "
            f"Escalating to human. Diagnostics: {diagnostics}"
        )


class BuilderRouter:
    """Router for selecting models and applying retry policies."""

    def __init__(
        self,
        registry: ModelRegistry,
        default_tier: str = "T2",
    ) -> None:
        """Initialize builder router.

        Args:
            registry: Model registry with tier → models mapping
            default_tier: Default tier to use if workstream has no capability_tier
        """
        self.registry = registry
        self.default_tier = default_tier

    def select_model(self, workstream: Workstream) -> ModelProvider:
        """Select model provider for workstream based on capability tier.

        Args:
            workstream: Workstream to build

        Returns:
            Selected ModelProvider instance

        Raises:
            ValueError: If tier is invalid or no models available
        """
        tier = workstream.capability_tier or self.default_tier
        return select_model_for_tier(tier, self.registry)

    def get_retry_policy(self, workstream: Workstream) -> RetryPolicy:
        """Get retry policy for workstream based on tier.

        For T2/T3: 3 attempts → escalate to human (policy D1).
        For T0/T1: No retry limit (can escalate to stronger model).

        Args:
            workstream: Workstream to build

        Returns:
            RetryPolicy instance
        """
        tier = workstream.capability_tier or self.default_tier

        # Policy D1: T2/T3 get 3 attempts → human escalation
        if tier in ("T2", "T3"):
            return RetryPolicy(max_attempts=3)

        # T0/T1: No hard limit (can escalate to stronger model)
        return RetryPolicy(max_attempts=10)

    def should_escalate_to_human(
        self,
        workstream: Workstream,
        attempt: int,
        failed: bool,
    ) -> bool:
        """Check if build should be escalated to human.

        Args:
            workstream: Workstream being built
            attempt: Current attempt number
            failed: Whether attempt failed

        Returns:
            True if should escalate to human
        """
        policy = self.get_retry_policy(workstream)
        return policy.should_escalate(attempt, failed)

    def create_escalation_error(
        self,
        workstream: Workstream,
        attempt: int,
        diagnostics: str,
    ) -> HumanEscalationError:
        """Create human escalation error with diagnostics.

        Args:
            workstream: Workstream that failed
            attempt: Number of failed attempts
            diagnostics: Diagnostic information

        Returns:
            HumanEscalationError instance
        """
        return HumanEscalationError(
            ws_id=workstream.ws_id,
            tier=workstream.capability_tier,
            attempts=attempt,
            diagnostics=diagnostics,
        )


def select_model_for_tier(tier: str, registry: ModelRegistry) -> ModelProvider:
    """Select model provider for given capability tier.

    Currently selects first available model. Future: add cost/availability logic.

    Args:
        tier: Capability tier (T0, T1, T2, T3)
        registry: Model registry

    Returns:
        Selected ModelProvider instance

    Raises:
        ValueError: If tier is invalid or no models available
    """
    models = registry.get_models_for_tier(tier)

    if not models:
        raise ValueError(f"No models available for tier {tier}")

    # For now, select first model (primary choice)
    # Future: add cost/availability/context size selection logic
    return models[0]
