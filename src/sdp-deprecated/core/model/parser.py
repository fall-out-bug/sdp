"""Model mapping table parser."""

import re

from sdp.core.model.models import ModelProvider


def parse_models_table(section: str) -> list[ModelProvider]:
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
            context_window = parse_context_window(context)

        else:
            # Old format without metrics - use defaults
            context = row_match.group(3).strip()
            tool_use_str = row_match.group(4).strip()
            notes = row_match.group(5).strip()

            cost = 0.0
            availability = 0.99
            context_window = parse_context_window(context)

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


def parse_context_window(context_str: str) -> int:
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
