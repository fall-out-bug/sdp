"""Model mapping registry loader."""

import re
from pathlib import Path

from sdp.core.model.models import ModelMappingError, ModelProvider, ModelRegistry
from sdp.core.model.parser import parse_models_table


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
        raise ModelMappingError(
            message=f"Model mapping file not found: {mapping_file}",
            mapping_file=mapping_file,
        )

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
        models = parse_models_table(tier_section)
        tiers[tier] = models

    return ModelRegistry(tiers=tiers)
