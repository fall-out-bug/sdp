"""T2/T3 contract validation for workstreams.

Validates contract requirements for T2 (Implementer) and T3 (Autocomplete) tiers.
"""

from sdp.core.workstream import Workstream
from sdp.validators.capability_tier_models import CapabilityTier, ValidationCheck


def validate_t2_t3(
    ws: Workstream,
    body: str,
    tier: CapabilityTier,
) -> list[ValidationCheck]:
    """Validate T2/T3 specific requirements.

    Args:
        ws: Parsed workstream
        body: Workstream body content
        tier: Capability tier (T2 or T3)

    Returns:
        List of validation checks
    """
    from sdp.validators.capability_tier_checks import (
        _check_contract_is_read_only,
        _check_no_vague_language,
        _check_placeholders_present,
        _check_scope_tiny,
        _check_signatures_complete,
        _check_tests_complete,
        _check_verification_commands,
    )
    from sdp.validators.capability_tier_extractors import (
        _extract_code_block,
        _extract_section,
    )

    checks: list[ValidationCheck] = []

    # Check 1: Contract section exists
    contract_section = _extract_section(body, "Contract")
    if not contract_section:
        checks.append(
            ValidationCheck(
                name="contract_section_exists",
                passed=False,
                message="Contract section is required for T2/T3",
            )
        )
        return checks  # Can't continue without contract
    else:
        checks.append(
            ValidationCheck(
                name="contract_section_exists",
                passed=True,
                message="Contract section found",
            )
        )

    # Check 2: Interface section exists with full signatures
    interface_section = _extract_section(contract_section, "Interface")
    if not interface_section:
        checks.append(
            ValidationCheck(
                name="interface_section_exists",
                passed=False,
                message="Interface section is required in Contract for T2/T3",
            )
        )
    else:
        interface_code = _extract_code_block(interface_section, "python")
        if not interface_code:
            checks.append(
                ValidationCheck(
                    name="interface_code_block",
                    passed=False,
                    message="Interface section must contain Python code block",
                )
            )
        else:
            sig_check = _check_signatures_complete(interface_code)
            checks.append(sig_check)

    # Check 3: Tests section exists with full tests
    tests_section = _extract_section(contract_section, "Tests")
    if not tests_section:
        checks.append(
            ValidationCheck(
                name="tests_section_exists",
                passed=False,
                message="Tests section is required in Contract for T2/T3",
            )
        )
    else:
        tests_code = _extract_code_block(tests_section, "python")
        if not tests_code:
            checks.append(
                ValidationCheck(
                    name="tests_code_block",
                    passed=False,
                    message="Tests section must contain Python code block",
                )
            )
        else:
            tests_check = _check_tests_complete(tests_code)
            checks.append(tests_check)

    # Check 4: NotImplementedError placeholders
    if interface_section:
        interface_code = _extract_code_block(interface_section, "python")
        if interface_code:
            placeholder_check = _check_placeholders_present(interface_code)
            checks.append(placeholder_check)

    # Check 5: No vague language
    vague_check = _check_no_vague_language(body)
    checks.append(vague_check)

    # Check 6: Verification commands exist
    verification_section = _extract_section(body, "Verification")
    if not verification_section:
        checks.append(
            ValidationCheck(
                name="verification_section_exists",
                passed=False,
                message="Verification section is required for T2/T3",
            )
        )
    else:
        verif_check = _check_verification_commands(verification_section)
        checks.append(verif_check)

    # Check 7: Contract is read-only (no permission to modify)
    read_only_check = _check_contract_is_read_only(body)
    checks.append(read_only_check)

    # Check 8: T3 specific - scope must be TINY (< 50 LOC estimate)
    if tier == CapabilityTier.T3:
        scope_check = _check_scope_tiny(ws, body)
        checks.append(scope_check)

    return checks
