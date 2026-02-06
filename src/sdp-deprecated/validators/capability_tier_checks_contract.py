"""Contract-specific validation checks for workstreams.

Provides validation functions for contract sections and verification commands.
"""

import re

from sdp.validators.capability_tier_models import ValidationCheck


def _check_verification_commands(verification_section: str) -> ValidationCheck:
    """Check verification section has concrete commands."""
    # Look for bash code blocks
    bash_blocks = re.findall(r"```bash\n(.*?)```", verification_section, re.DOTALL)

    if not bash_blocks:
        return ValidationCheck(
            name="verification_commands",
            passed=False,
            message="Verification section must contain bash code blocks with commands",
        )

    # Check blocks have actual commands (not just comments)
    has_commands = False
    for block in bash_blocks:
        # Remove comments and empty lines
        lines = [
            line.strip()
            for line in block.split("\n")
            if line.strip() and not line.strip().startswith("#")
        ]
        if lines:
            has_commands = True
            break

    if not has_commands:
        return ValidationCheck(
            name="verification_commands",
            passed=False,
            message="Verification bash blocks must contain actual commands",
        )

    return ValidationCheck(
        name="verification_commands",
        passed=True,
        message=f"Verification has {len(bash_blocks)} bash code block(s) with commands",
    )


def _check_contract_is_read_only(body: str) -> ValidationCheck:
    """Check contract appears to be a specification (not implementation)."""
    # Look for Contract section
    contract_match = re.search(r"#+\s*Contract", body, re.IGNORECASE)
    if not contract_match:
        return ValidationCheck(
            name="contract_read_only",
            passed=True,  # No contract to check
            message="No contract section found",
        )

    # Extract contract section (up to next heading or end)
    contract_start = contract_match.end()
    next_heading = re.search(r"\n#+\s", body[contract_start:])
    contract_end = contract_start + next_heading.start() if next_heading else len(body)
    contract_text = body[contract_start:contract_end]

    # Check for implementation indicators (should NOT be present in contract)
    impl_patterns = [
        r"class\s+\w+:",  # Class definitions
        r"def\s+\w+\([^)]*\):\s*\n\s+[^\s]",  # Non-pass function bodies
    ]

    has_implementation = False
    for pattern in impl_patterns:
        if re.search(pattern, contract_text):
            has_implementation = True
            break

    if has_implementation:
        return ValidationCheck(
            name="contract_read_only",
            passed=False,
            message="Contract section should contain specifications, not implementation",
            details=["Use NotImplementedError placeholders for function bodies"],
        )

    return ValidationCheck(
        name="contract_read_only",
        passed=True,
        message="Contract is read-only (specifications only)",
    )


def _check_no_vague_language(body: str) -> ValidationCheck:
    """Check workstream uses specific, actionable language."""
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
