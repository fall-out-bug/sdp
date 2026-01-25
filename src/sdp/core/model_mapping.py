"""Model mapping registry loader.

Loads capability tier → model mappings from sdp/docs/model-mapping.md.
"""

import re
from dataclasses import dataclass
from pathlib import Path


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


class ModelMappingError(Exception):
    """Error loading model mapping registry."""

    pass


def load_model_registry(mapping_file: Path) -> ModelRegistry:
    """Load model mapping registry from markdown file.

    Args:
        mapping_file: Path to model-mapping.md file

    Returns:
        ModelRegistry instance with tier → models mapping

    Raises:
        ModelMappingError: If file not found or cannot be parsed
    """
    if not mapping_file.exists():
        raise ModelMappingError(f"Model mapping file not found: {mapping_file}")

    content = mapping_file.read_text(encoding="utf-8")
    tiers: dict[str, list[ModelProvider]] = {}

    # Find all tier sections (### T0 —, ### T1 —, etc.)
    # Match: ### T0 — Architect or ### T0 — Architect (Contract Writer)
    tier_pattern = r"^### (T\d+) —"
    tier_matches = list(re.finditer(tier_pattern, content, re.MULTILINE))

    for i, tier_match in enumerate(tier_matches):
        tier = tier_match.group(1)
        start_pos = tier_match.end()

        # Find end of this tier section (next tier or end of file)
        if i + 1 < len(tier_matches):
            end_pos = tier_matches[i + 1].start()
        else:
            end_pos = len(content)

        tier_section = content[start_pos:end_pos]

        # Extract models from table
        models = _parse_models_table(tier_section)
        tiers[tier] = models

    return ModelRegistry(tiers=tiers)


def _parse_models_table(section: str) -> list[ModelProvider]:
    """Parse models from markdown table in tier section.

    Args:
        section: Markdown section content for one tier

    Returns:
        List of ModelProvider instances
    """
    models: list[ModelProvider] = []

    # Find table (starts with | Provider | Model | ...)
    # Match header row and separator, then capture all data rows
    table_pattern = r"\| Provider \| Model \|.*?\n\|[-\s|]+\n((?:\|.*?\n)+)"
    table_match = re.search(table_pattern, section, re.MULTILINE)

    if not table_match:
        return models

    table_rows = table_match.group(1)

    # Parse each row - try 8 columns first (new format), then 5 columns (old format)
    # New format: Provider | Model | Cost ($/1M) | Availability | Context |
    #              Tool Use | Notes
    # Old format: Provider | Model | Context | Tool Use | Notes
    row_pattern_new = (
        r"\| ([^|]+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \| "
        r"([^^|]+) \| ([^|]+) \| ([^|]+) \|"
    )
    row_pattern_old = r"\| ([^|]+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \|"

    # Try new format first
    matches = list(re.finditer(row_pattern_new, table_rows))
    use_new_format = len(matches) > 0

    if not use_new_format:
        # Fall back to old format
        matches = list(re.finditer(row_pattern_old, table_rows))

    for row_match in matches:
        provider = row_match.group(1).strip()
        model = row_match.group(2).strip()

        # Skip empty rows
        if not provider or not model:
            continue

        if use_new_format:
            # New format with metrics
            cost_str = row_match.group(3).strip()
            availability_str = row_match.group(4).strip()
            context = row_match.group(5).strip()
            tool_use_str = row_match.group(6).strip()
            notes = row_match.group(7).strip()

            # Parse cost (e.g., "3.00", "0.25", "0.00")
            try:
                cost = float(cost_str.replace("$", "").strip())
            except ValueError:
                cost = 0.0

            # Parse availability (e.g., "99.9%", "99.0%", "100%")
            try:
                availability = float(availability_str.replace("%", "").strip()) / 100.0
            except ValueError:
                availability = 0.99

            # Parse context window (e.g., "200K", "128K", "1M+", "2K")
            context_window = _parse_context_window(context)

        else:
            # Old format without metrics - use defaults
            context = row_match.group(3).strip()
            tool_use_str = row_match.group(4).strip()
            notes = row_match.group(5).strip()

            cost = 0.0
            availability = 0.99
            context_window = _parse_context_window(context)

        # Parse tool_use (✅ Full, ⚠️ Limited, ❌ None)
        tool_use = "✅" in tool_use_str or "Full" in tool_use_str

        models.append(
            ModelProvider(
                provider=provider,
                model=model,
                context=context,
                tool_use=tool_use,
                cost_per_1m_tokens=cost,
                availability_pct=availability,
                context_window=context_window,
                notes=notes,
            )
        )

    return models


def _parse_context_window(context_str: str) -> int:
    """Parse context window string to integer tokens.

    Args:
        context_str: Context string like "200K", "128K", "1M+", "2K"

    Returns:
        Context window in tokens
    """
    context_str = context_str.strip().upper()

    # Handle "M+" suffix (e.g., "1M+" = 1,000,000)
    if "M+" in context_str:
        return int(float(context_str.replace("M+", "").replace("M", "")) * 1_000_000)
    elif "M" in context_str:
        return int(float(context_str.replace("M", "")) * 1_000_000)

    # Handle "K" suffix (e.g., "200K" = 200,000)
    if "K" in context_str:
        return int(float(context_str.replace("K", "")) * 1_000)

    # Assume plain number
    try:
        return int(context_str)
    except ValueError:
        return 128_000  # Default fallback
