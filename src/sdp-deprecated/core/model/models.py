"""Model provider and registry data models."""

from dataclasses import dataclass
from pathlib import Path

from sdp.errors import ErrorCategory, SDPError


@dataclass
class ModelProvider:
    """Single model provider entry with performance metrics."""

    provider: str
    model: str
    context: str
    tool_use: bool
    cost_per_1m_tokens: float = 0.0
    availability_pct: float = 0.99
    context_window: int = 128000
    notes: str = ""


@dataclass
class ModelRegistry:
    """Registry of models by capability tier."""

    tiers: dict[str, list[ModelProvider]]

    def get_models_for_tier(self, tier: str) -> list[ModelProvider]:
        """Get list of models for given tier.

        Args:
            tier: Capability tier (T0, T1, T2, T3)

        Returns:
            List of ModelProvider instances

        Raises:
            ValueError: If tier is invalid
        """
        tier_upper = tier.upper()
        if tier_upper not in self.tiers:
            raise ValueError(f"Invalid tier: {tier}. Must be one of T0, T1, T2, T3")
        return self.tiers[tier_upper]


class ModelMappingError(SDPError):
    """Error loading model mapping registry."""

    def __init__(self, message: str, mapping_file: Path | None = None) -> None:
        super().__init__(
            category=ErrorCategory.CONFIGURATION,
            message=message,
            remediation=(
                "1. Ensure model-mapping.md exists in docs/\n"
                "2. Check file format (markdown with tier sections)\n"
                "3. Verify YAML/table syntax in tier sections\n"
                "4. See docs/reference/model-mapping.md for format"
            ),
            docs_url="https://docs.sdp.dev/model-mapping",
            context={"mapping_file": str(mapping_file)} if mapping_file else None,
        )
