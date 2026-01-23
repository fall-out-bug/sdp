"""Model mapping registry loader.

Loads capability tier → model mappings from sdp/docs/model-mapping.md.
"""

import re
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ModelProvider:
    """Single model provider entry."""

    provider: str
    model: str
    context: str
    tool_use: bool
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

    # Parse each row (5 columns: Provider | Model | Context | Tool Use | Notes)
    row_pattern = r"\| ([^|]+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \|"
    for row_match in re.finditer(row_pattern, table_rows):
        provider = row_match.group(1).strip()
        model = row_match.group(2).strip()
        context = row_match.group(3).strip()
        tool_use_str = row_match.group(4).strip()
        notes = row_match.group(5).strip()

        # Skip empty rows
        if not provider or not model:
            continue

        # Parse tool_use (✅ Full, ⚠️ Limited, ❌ None)
        tool_use = "✅" in tool_use_str or "Full" in tool_use_str

        models.append(
            ModelProvider(
                provider=provider,
                model=model,
                context=context,
                tool_use=tool_use,
                notes=notes,
            )
        )

    return models
