"""T0/T1 gate validation for workstreams.

Validates architectural gate requirements for T0 (Architect) and T1 (Integrator) tiers.
"""

import re
from sdp.core.workstream import Workstream


def _check_no_time_estimates(body: str) -> "ValidationCheck":
    """Check workstream contains no time estimates."""
    from sdp.validators.capability_tier_models import ValidationCheck

    time_patterns = [
        r"\d+\s*(hour|hr|minute|min|day|week|s)",
        r"\d+:\d+",  # HH:MM format
    ]

    found_patterns = []
    for pattern in time_patterns:
        matches = re.findall(pattern, body, re.IGNORECASE)
        if matches:
            found_patterns.extend(matches)

    if found_patterns:
        return ValidationCheck(
            name="no_time_estimates",
            passed=False,
            message="Time estimates found (forbidden in contract-driven development)",
            details=list(set(found_patterns)),
        )

    return ValidationCheck(
        name="no_time_estimates",
        passed=True,
        message="No time estimates found",
    )


def _check_no_vague_language(body: str) -> "ValidationCheck":
    """Check workstream uses specific, actionable language."""
    from sdp.validators.capability_tier_models import ValidationCheck

    vague_patterns = [
        r"\b(appropriate|suitable|reasonable|proper|correct)\b",
        r"\b(etc|etcetera|and so forth)\b",
        r"\b(various|multiple|several|few|some)\b",
        r"\b(maybe|perhaps|possibly|might)\b",
    ]

    found_matches = []
    for pattern in vague_patterns:
        matches = re.findall(pattern, body, re.IGNORECASE)
        if matches:
            found_matches.extend(matches)

    if found_matches:
        return ValidationCheck(
            name="no_vague_language",
            passed=False,
            message="Vague language found (use specific, actionable language)",
            details=list(set(found_matches)),
        )

    return ValidationCheck(
        name="no_vague_language",
        passed=True,
        message="Language is specific and actionable",
    )


def _check_architecture_section_present(body: str) -> "ValidationCheck":
    """Check architecture section exists and has content."""
    from sdp.validators.capability_tier_models import ValidationCheck
    from sdp.validators.capability_tier_extractors import _extract_section

    arch_section = _extract_section(body, "Architecture")

    if not arch_section:
        return ValidationCheck(
            name="architecture_section",
            passed=False,
            message="Architecture section is required for T0/T1",
        )

    return ValidationCheck(
        name="architecture_section",
        passed=True,
        message="Architecture section found",
    )


def _check_design_decisions_documented(body: str) -> "ValidationCheck":
    """Check design decisions are documented with rationale."""
    from sdp.validators.capability_tier_models import ValidationCheck

    # Look for decision patterns like "Decision:", "Rationale:", "Alternative:"
    decision_patterns = [
        r"#+\s*Decision",
        r"#+\s*Rationale",
        r"#+\s*Alternative",
        r"\*\*Decision:\*\*",
        r"\*\*Rationale:\*\*",
    ]

    found = False
    for pattern in decision_patterns:
        if re.search(pattern, body, re.IGNORECASE):
            found = True
            break

    if not found:
        return ValidationCheck(
            name="design_decisions",
            passed=False,
            message="Design decisions should be documented with rationale",
        )

    return ValidationCheck(
        name="design_decisions",
        passed=True,
        message="Design decisions documented",
    )


def validate_t0_t1_gates(ws: Workstream, body: str) -> list["ValidationCheck"]:
    """Validate T0/T1 gate requirements.

    Args:
        ws: Parsed workstream
        body: Workstream body content

    Returns:
        List of validation checks
    """
    checks: list[ValidationCheck] = []

    # Gate 1: No time estimates
    time_estimate_check = _check_no_time_estimates(body)
    checks.append(time_estimate_check)

    # Gate 2: No vague language
    vague_check = _check_no_vague_language(body)
    checks.append(vague_check)

    # Gate 3: Architecture section
    arch_check = _check_architecture_section_present(body)
    checks.append(arch_check)

    # Gate 4: Design decisions documented
    decisions_check = _check_design_decisions_documented(body)
    checks.append(decisions_check)

    return checks
