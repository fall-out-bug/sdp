"""Builder router for model-agnostic workstream execution."""

from sdp.core.builder.model_selector import select_model_for_tier
from sdp.core.builder.policies import HumanEscalationError, RetryPolicy
from sdp.core.model_mapping import ModelProvider, ModelRegistry
from sdp.core.workstream import Workstream


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
        # Get capability_tier from workstream if available, otherwise use default
        tier = getattr(workstream, "capability_tier", None) or self.default_tier
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
        tier = getattr(workstream, "capability_tier", None) or self.default_tier

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
        tier = getattr(workstream, "capability_tier", None) or self.default_tier
        return HumanEscalationError(
            ws_id=workstream.ws_id,
            tier=tier,
            attempts=attempt,
            diagnostics=diagnostics,
        )
