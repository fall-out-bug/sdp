"""Capability tier validation for Contract-Driven WS v2.0.

This module validates workstreams against capability tiers (T0-T3) to ensure
they meet the requirements for model-agnostic execution.
"""

import re
from pathlib import Path

from sdp.core.workstream import WorkstreamParseError, parse_workstream
from sdp.validators.capability_tier_models import CapabilityTier, ValidationCheck, ValidationResult
from sdp.validators.capability_tier_t0_t1 import validate_t0_t1_gates
from sdp.validators.capability_tier_t2_t3 import validate_t2_t3


def validate_workstream_tier(ws_path: Path, tier: str) -> ValidationResult:
    """Validate workstream compatibility with capability tier.

    Args:
        ws_path: Path to WS markdown file
        tier: "T0" | "T1" | "T2" | "T3"

    Returns:
        ValidationResult with passed/failed checks
    """
    if tier != tier.upper():
        raise ValueError(f"Invalid tier: {tier}. Must be uppercase (T0, T1, T2, T3)")

    try:
        tier_enum = CapabilityTier(tier)
    except ValueError as e:
        raise ValueError(f"Invalid tier: {tier}. Must be one of T0, T1, T2, T3") from e

    if not ws_path.exists():
        raise WorkstreamParseError(
            message=f"File not found: {ws_path}",
            file_path=ws_path,
        )

    try:
        ws = parse_workstream(ws_path)
    except WorkstreamParseError as e:
        result = ValidationResult(tier=tier_enum, passed=False)
        result.add_check(
            ValidationCheck(
                name="parse_workstream",
                passed=False,
                message=f"Failed to parse workstream: {e}",
            )
        )
        return result

    content = ws_path.read_text(encoding="utf-8")
    body = _extract_body(content)

    result = ValidationResult(tier=tier_enum, passed=True)

    # T2/T3 specific checks
    if tier_enum in (CapabilityTier.T2, CapabilityTier.T3):
        checks = validate_t2_t3(ws, body, tier_enum)
        for check in checks:
            result.add_check(check)

    # T0/T1 gate checks
    if tier_enum in (CapabilityTier.T0, CapabilityTier.T1):
        checks = validate_t0_t1_gates(ws, body)
        for check in checks:
            result.add_check(check)

    return result


def _extract_body(content: str) -> str:
    """Extract markdown body (without frontmatter)."""
    match = re.match(r"^---\n.*?\n---\n(.*)", content, re.DOTALL)
    return match.group(1) if match else content
