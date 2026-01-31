"""Scope validation checks for workstreams.

Provides validation functions for workstream size and scope estimates.
"""

from sdp.validators.capability_tier_models import ValidationCheck


def _check_scope_tiny(ws, body: str) -> ValidationCheck:
    """Check T3 workstream has tiny scope (< 50 LOC estimate)."""
    # Check workstream size
    if ws.size.value in ("LARGE", "XLARGE"):
        return ValidationCheck(
            name="scope_tiny",
            passed=False,
            message="T3 workstreams must be TINY scope (< 50 LOC estimate)",
            details=[f"Current size: {ws.size.value}"],
        )

    if ws.size.value != "TINY":
        return ValidationCheck(
            name="scope_tiny",
            passed=False,
            message="T3 workstreams must be TINY scope",
            details=[f"Current size: {ws.size.value}"],
        )

    return ValidationCheck(
        name="scope_tiny",
        passed=True,
        message="Workstream scope is TINY (< 50 LOC)",
    )
