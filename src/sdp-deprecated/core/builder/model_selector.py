"""Model selection utilities for builder routing."""

from typing import Optional

from sdp.core.model import ModelProvider, ModelRegistry

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
