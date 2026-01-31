"""Builder router for model-agnostic workstream execution.

Routes /build commands to appropriate model providers based on capability tier
and applies retry/escalation policies.
"""

from dataclasses import dataclass
from typing import Optional

from sdp.core.model_mapping import ModelProvider, ModelRegistry
from sdp.core.workstream import Workstream
from sdp.errors import ErrorCategory, SDPError

# Default weights for model selection
DEFAULT_WEIGHTS = {
    "cost": 0.4,
    "availability": 0.3,
    "context": 0.3,
}


def select_model_weighted(
    models: list[ModelProvider],
    weights: dict[str, float],
    required_context: int = 0,
) -> ModelProvider:
    """Select model using weighted scoring.

    Args:
        models: Available models for tier
        weights: Weight dict (cost, availability, context)
        required_context: Minimum context window required

    Returns:
        Best model according to weighted score

    Raises:
        ValueError: If no models meet requirements
    """
    # Filter by context requirements
    if required_context > 0:
        candidates = [m for m in models if m.context_window >= required_context]
    else:
        candidates = models

    if not candidates:
        raise ValueError(f"No models with context >= {required_context}")

    # Normalize metrics (0-1 scale, higher is better)
    max_cost = max(m.cost_per_1m_tokens for m in candidates) if candidates else 1.0
    if max_cost == 0:
        max_cost = 1.0  # Avoid division by zero

    scores = []
    for model in candidates:
        # Cost score: lower cost = higher score (inverted)
        if max_cost > 0:
            cost_score = 1 - (model.cost_per_1m_tokens / max_cost)
        else:
            cost_score = 1.0

        # Availability score: already 0-1
        avail_score = model.availability_pct

        # Context score: normalize to 200K max
        context_score = min(model.context_window / 200_000, 1.0)

        # Weighted score
        weighted_score = (
            weights["cost"] * cost_score +
            weights["availability"] * avail_score +
            weights["context"] * context_score
        )
        scores.append((weighted_score, model))

    # Return highest score
    scores.sort(key=lambda x: x[0], reverse=True)
    return scores[0][1]


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


def select_model_for_tier(
    tier: str,
    registry: ModelRegistry,
    required_context: int = 0,
    weights: Optional[dict[str, float]] = None,
) -> ModelProvider:
    """Select model provider for given capability tier.

    Uses weighted selection based on cost, availability, context.

    Args:
        tier: Capability tier (T0, T1, T2, T3)
        registry: Model registry
        required_context: Minimum context tokens needed
        weights: Custom weights (defaults to DEFAULT_WEIGHTS)

    Returns:
        Selected ModelProvider instance

    Raises:
        ValueError: If tier is invalid or no models available
    """
    models = registry.get_models_for_tier(tier)

    if not models:
        raise ValueError(f"No models available for tier {tier}")

    return select_model_weighted(
        models,
        weights or DEFAULT_WEIGHTS,
        required_context
    )
