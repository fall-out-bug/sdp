"""Capability tier validation for Contract-Driven WS v2.0.

This module validates workstreams against capability tiers (T0-T3) to ensure
they meet the requirements for model-agnostic execution.
"""

import re
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Optional

from sdp.core.workstream import Workstream, WorkstreamParseError, WorkstreamSize, parse_workstream


class CapabilityTier(str, Enum):
    """Capability tier levels."""

    T0 = "T0"  # Architect (Contract writer)
    T1 = "T1"  # Integrator (Complex build)
    T2 = "T2"  # Implementer (Contract-driven build)
    T3 = "T3"  # Autocomplete (Micro build)


@dataclass
class ValidationCheck:
    """Single validation check result."""

    name: str
    passed: bool
    message: str = ""
    details: list[str] = field(default_factory=list)


@dataclass
class ValidationResult:
    """Result of workstream tier validation."""

    tier: CapabilityTier
    passed: bool
    checks: list[ValidationCheck] = field(default_factory=list)

    def add_check(self, check: ValidationCheck) -> None:
        """Add a validation check result."""
        self.checks.append(check)
        if not check.passed:
            self.passed = False


def validate_workstream_tier(ws_path: Path, tier: str) -> ValidationResult:
    """Validate workstream compatibility with capability tier.

    Args:
        ws_path: Path to WS markdown file
        tier: "T0" | "T1" | "T2" | "T3"

    Returns:
        ValidationResult with passed/failed checks
    """
    try:
        tier_enum = CapabilityTier(tier.upper())
    except ValueError as e:
        raise ValueError(f"Invalid tier: {tier}. Must be one of T0, T1, T2, T3") from e

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
        checks = _validate_t2_t3(ws, body, tier_enum)
        for check in checks:
            result.add_check(check)

    # T0/T1 gate checks
    if tier_enum in (CapabilityTier.T0, CapabilityTier.T1):
        checks = _validate_t0_t1_gates(ws, body)
        for check in checks:
            result.add_check(check)

    return result


def _extract_body(content: str) -> str:
    """Extract markdown body (without frontmatter)."""
    match = re.match(r"^---\n.*?\n---\n(.*)", content, re.DOTALL)
    return match.group(1) if match else content


def _extract_section(body: str, section_name: str) -> str:
    """Extract content of a section by name (case-insensitive)."""
    # Match ### Section or #### Subsection
    heading_pattern = rf"^###+ .*{section_name}.*$"
    heading_match = re.search(heading_pattern, body, re.MULTILINE | re.IGNORECASE)

    if not heading_match:
        return ""

    start_pos = heading_match.end() + 1
    if start_pos >= len(body):
        return ""

    remaining = body[start_pos:]
    # Find next heading at same or higher level
    next_heading = re.search(r"^###", remaining, re.MULTILINE)

    end_pos = next_heading.start() if next_heading else len(remaining)
    content = remaining[:end_pos]
    return content.strip()


def _extract_code_block(section: str, language: str = "python") -> Optional[str]:
    """Extract code block from section."""
    pattern = rf"```{language}\n(.*?)```"
    match = re.search(pattern, section, re.DOTALL)
    return match.group(1).strip() if match else None


def _validate_t2_t3(ws: Workstream, body: str, tier: CapabilityTier) -> list[ValidationCheck]:
    """Validate T2/T3 specific requirements."""
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


def _validate_t0_t1_gates(ws: Workstream, body: str) -> list[ValidationCheck]:
    """Validate T0/T1 gate requirements."""
    checks: list[ValidationCheck] = []

    # Gate 1: No time estimates
    time_estimate_check = _check_no_time_estimates(body)
    checks.append(time_estimate_check)

    # Gate 2: Scope <= MEDIUM
    scope_check = _check_scope_not_large(ws)
    checks.append(scope_check)

    return checks


def _check_signatures_complete(interface_code: str) -> ValidationCheck:
    """Check that all functions have complete type hints."""
    # Look for function definitions
    func_pattern = r"def\s+(\w+)\s*\([^)]*\)\s*(->\s*[^:]+)?:"
    functions = re.findall(func_pattern, interface_code)

    if not functions:
        return ValidationCheck(
            name="signatures_complete",
            passed=False,
            message="No function signatures found in Interface section",
        )

    issues: list[str] = []
    # Check each function for type hints
    for func_match in re.finditer(r"def\s+(\w+)\s*\(([^)]*)\)\s*(->\s*[^:]+)?:", interface_code):
        params = func_match.group(2)
        return_type = func_match.group(3)

        # Check return type
        if not return_type:
            issues.append(f"Function {func_match.group(1)} missing return type hint")

        # Check parameters (basic check - should have type hints)
        if params and params.strip() != "self":
            # Simple check: parameters should have colons (type hints)
            param_list = [p.strip() for p in params.split(",")]
            for param in param_list:
                if param and ":" not in param and param != "self" and param != "*args" and param != "**kwargs":
                    issues.append(f"Parameter '{param}' in {func_match.group(1)} missing type hint")

    if issues:
        return ValidationCheck(
            name="signatures_complete",
            passed=False,
            message="Some functions missing type hints",
            details=issues[:5],  # Limit details
        )

    return ValidationCheck(
        name="signatures_complete",
        passed=True,
        message=f"All {len(functions)} functions have complete type hints",
    )


def _check_tests_complete(tests_code: str) -> ValidationCheck:
    """Check that tests are complete and executable."""
    # Look for test functions
    test_pattern = r"def\s+test_\w+"
    tests = re.findall(test_pattern, tests_code)

    if not tests:
        return ValidationCheck(
            name="tests_complete",
            passed=False,
            message="No test functions found in Tests section",
        )

    # Check for assertions
    assert_count = len(re.findall(r"\bassert\b", tests_code))
    if assert_count == 0:
        return ValidationCheck(
            name="tests_complete",
            passed=False,
            message="Tests section contains no assertions",
        )

    return ValidationCheck(
        name="tests_complete",
        passed=True,
        message=f"Found {len(tests)} test functions with {assert_count} assertions",
    )


def _check_placeholders_present(interface_code: str) -> ValidationCheck:
    """Check that functions have NotImplementedError placeholders."""
    # Look for function bodies
    func_pattern = r"def\s+\w+\s*\([^)]*\)\s*(->\s*[^:]+)?:\s*\n\s*(.*?)(?=\n\s*def|\n\s*class|\Z)"
    matches = list(re.finditer(func_pattern, interface_code, re.DOTALL))

    if not matches:
        return ValidationCheck(
            name="placeholders_present",
            passed=False,
            message="No function bodies found",
        )

    issues: list[str] = []
    for match in matches:
        body = match.group(2).strip()
        # Check if body contains NotImplementedError or pass or ...
        if not re.search(r"(NotImplementedError|pass|\.\.\.)", body, re.IGNORECASE):
            func_name_match = re.search(r"def\s+(\w+)", match.group(0))
            func_name = func_name_match.group(1) if func_name_match else "unknown"
            issues.append(f"Function {func_name} missing NotImplementedError/pass placeholder")

    if issues:
        return ValidationCheck(
            name="placeholders_present",
            passed=False,
            message="Some functions missing placeholders",
            details=issues[:5],
        )

    return ValidationCheck(
        name="placeholders_present",
        passed=True,
        message="All functions have placeholders",
    )


def _check_no_vague_language(body: str) -> ValidationCheck:
    """Check for vague language patterns that are forbidden for T2/T3."""
    vague_patterns = [
        (r"по\s+смыслу", "by meaning"),
        (r"как\s+считаешь\s+правильно", "as you think right"),
        (r"на\s+сво[ёе]\s+усмотрение", "at your discretion"),
        (r"выбери\s+подход", "choose approach"),
        (r"определи\s+архитектуру", "define architecture"),
        (r"реши\s+где\s+лучше", "decide where better"),
        (r"смотри\s+в\s+файле", "see in file"),  # Without quote context
        (r"убедись\s+что\s+(ок|красиво|правильно)", "make sure it's ok/beautiful/right"),
        (r"реализуй\s+по\s+смыслу", "implement by meaning"),
        (r"на\s+тво[ёе]\s+усмотрение", "at your discretion"),
    ]

    found: list[str] = []
    for pattern, description in vague_patterns:
        matches = re.finditer(pattern, body, re.IGNORECASE)
        for match in matches:
            # Get context (20 chars before and after)
            start = max(0, match.start() - 20)
            end = min(len(body), match.end() + 20)
            context = body[start:end].replace("\n", " ")
            found.append(f"{description}: ...{context}...")

    if found:
        return ValidationCheck(
            name="no_vague_language",
            passed=False,
            message=f"Found {len(found)} vague language patterns",
            details=found[:5],  # Limit to 5 examples
        )

    return ValidationCheck(
        name="no_vague_language",
        passed=True,
        message="No vague language patterns found",
    )


def _check_verification_commands(verification_section: str) -> ValidationCheck:
    """Check that verification section contains bash commands."""
    # Look for code blocks with bash/shell
    bash_blocks = re.findall(r"```(?:bash|sh|shell)\n(.*?)```", verification_section, re.DOTALL)

    if not bash_blocks:
        return ValidationCheck(
            name="verification_commands",
            passed=False,
            message="Verification section must contain bash/shell code block",
        )

    # Check for commands that can exit 0
    commands = []
    for block in bash_blocks:
        lines = [line.strip() for line in block.split("\n") if line.strip() and not line.strip().startswith("#")]
        commands.extend(lines)

    if not commands:
        return ValidationCheck(
            name="verification_commands",
            passed=False,
            message="Verification section contains no executable commands",
        )

    # Check for common test/lint commands
    has_test_command = any(
        re.search(r"(pytest|test|check|lint|mypy|ruff)", cmd, re.IGNORECASE) for cmd in commands
    )

    if not has_test_command:
        return ValidationCheck(
            name="verification_commands",
            passed=True,  # Not required, but recommended
            message=f"Found {len(commands)} commands (test/lint commands recommended)",
        )

    return ValidationCheck(
        name="verification_commands",
        passed=True,
        message=f"Found {len(commands)} verification commands including test/lint",
    )


def _check_contract_is_read_only(body: str) -> ValidationCheck:
    """Check that contract sections are marked as read-only."""
    # Look for explicit read-only markers
    read_only_patterns = [
        r"DO\s+NOT\s+MODIFY",
        r"read-only",
        r"readonly",
        r"не\s+менять",
        r"запрещено\s+изменять",
    ]

    contract_section = _extract_section(body, "Contract")
    if not contract_section:
        return ValidationCheck(
            name="contract_read_only",
            passed=False,
            message="Contract section not found",
        )

    # Check Interface section
    interface_section = _extract_section(contract_section, "Interface")
    has_interface_marker = False
    if interface_section:
        for pattern in read_only_patterns:
            if re.search(pattern, interface_section, re.IGNORECASE):
                has_interface_marker = True
                break

    # Check Tests section
    tests_section = _extract_section(contract_section, "Tests")
    has_tests_marker = False
    if tests_section:
        for pattern in read_only_patterns:
            if re.search(pattern, tests_section, re.IGNORECASE):
                has_tests_marker = True
                break

    issues: list[str] = []
    if interface_section and not has_interface_marker:
        issues.append("Interface section missing read-only marker")
    if tests_section and not has_tests_marker:
        issues.append("Tests section missing read-only marker")

    if issues:
        return ValidationCheck(
            name="contract_read_only",
            passed=False,
            message="Contract sections should be marked as read-only",
            details=issues,
        )

    return ValidationCheck(
        name="contract_read_only",
        passed=True,
        message="Contract sections properly marked as read-only",
    )


def _check_scope_tiny(ws: Workstream, body: str) -> ValidationCheck:
    """Check that T3 WS has TINY scope (< 50 LOC)."""
    # Check size field
    if ws.size != WorkstreamSize.SMALL:
        return ValidationCheck(
            name="scope_tiny",
            passed=False,
            message=f"T3 requires SMALL size, got {ws.size.value}",
        )

    # Check for scope estimates in body
    scope_section = _extract_section(body, "Scope Estimate")
    if scope_section:
        # Look for LOC estimates
        loc_pattern = r"(\d+)\s*(?:LOC|lines|строк)"
        loc_matches = re.findall(loc_pattern, scope_section, re.IGNORECASE)
        for loc_str in loc_matches:
            try:
                loc = int(loc_str)
                if loc >= 50:
                    return ValidationCheck(
                        name="scope_tiny",
                        passed=False,
                        message=f"T3 requires < 50 LOC, found estimate: {loc}",
                    )
            except ValueError:
                pass

    return ValidationCheck(
        name="scope_tiny",
        passed=True,
        message="Scope is TINY (< 50 LOC)",
    )


def _check_no_time_estimates(body: str) -> ValidationCheck:
    """Check for time estimates (forbidden for T0/T1)."""
    time_patterns = [
        r"\d+\s*(?:час|часа|часов|hour|hours)",
        r"\d+\s*(?:день|дня|дней|day|days)",
        r"\d+\s*(?:неделя|недели|недель|week|weeks)",
        r"\d+\s*(?:минут|минуты|minute|minutes)",
    ]

    found: list[str] = []
    for pattern in time_patterns:
        matches = re.finditer(pattern, body, re.IGNORECASE)
        for match in matches:
            start = max(0, match.start() - 20)
            end = min(len(body), match.end() + 20)
            context = body[start:end].replace("\n", " ")
            found.append(f"...{context}...")

    if found:
        return ValidationCheck(
            name="no_time_estimates",
            passed=False,
            message=f"Found {len(found)} time estimates (forbidden)",
            details=found[:5],
        )

    return ValidationCheck(
        name="no_time_estimates",
        passed=True,
        message="No time estimates found",
    )


def _check_scope_not_large(ws: Workstream) -> ValidationCheck:
    """Check that scope is not LARGE (forbidden for T0/T1)."""
    if ws.size == WorkstreamSize.LARGE:
        return ValidationCheck(
            name="scope_not_large",
            passed=False,
            message="LARGE scope is forbidden for T0/T1 (must be SMALL or MEDIUM)",
        )

    return ValidationCheck(
        name="scope_not_large",
        passed=True,
        message=f"Scope is {ws.size.value} (allowed)",
    )
