"""Workstream completion verifier - validates WS is actually complete with evidence."""

from pathlib import Path
from typing import Any

from sdp.validators.ws_completion.checkers import (
    extract_coverage_value,
    verify_commands as check_commands,
    verify_coverage as check_coverage,
    verify_output_files as check_output_files,
)
from sdp.validators.ws_completion.models import CheckResult, VerificationResult
from sdp.validators.ws_completion.parser import find_ws_file, parse_ws_file


class WSCompletionVerifier:
    """Verify workstream completion with evidence."""

    def __init__(self, ws_dir: Path = Path("docs/workstreams")):
        """Initialize verifier.

        Args:
            ws_dir: Base directory for workstream files
        """
        self.ws_dir = ws_dir

    def verify_output_files(self, ws_data: dict[str, Any]) -> list[CheckResult]:
        """Check all output files in scope exist.

        Args:
            ws_data: Parsed workstream data

        Returns:
            List of check results for each file
        """
        return check_output_files(ws_data)

    def verify_commands(self, ws_data: dict[str, Any]) -> list[CheckResult]:
        """Run verification commands and check exit codes.

        Args:
            ws_data: Parsed workstream data

        Returns:
            List of check results for each command
        """
        return check_commands(ws_data)

    def verify_coverage(self, ws_data: dict[str, Any]) -> CheckResult | None:
        """Check test coverage meets threshold.

        Args:
            ws_data: Parsed workstream data

        Returns:
            Check result or None if no coverage requirement
        """
        return check_coverage(ws_data)

    def _find_ws_file(self, ws_id: str) -> Path | None:
        """Find WS file by ID (backward compatibility).

        Args:
            ws_id: Workstream ID

        Returns:
            Path to WS file or None
        """
        return find_ws_file(ws_id, self.ws_dir)

    def _parse_coverage_from_output(self, output: str) -> float | None:
        """Parse coverage percentage from pytest output (backward compatibility).

        Args:
            output: pytest stdout

        Returns:
            Coverage percentage or None
        """
        from sdp.validators.ws_completion.checkers import _parse_coverage_from_output

        return _parse_coverage_from_output(output)

    def _extract_coverage(self, check: CheckResult | None) -> float | None:
        """Extract coverage value from check result (backward compatibility).

        Args:
            check: Coverage check result

        Returns:
            Coverage percentage or None
        """
        return extract_coverage_value(check)

    def _parse_ws_file(self, ws_path: Path) -> dict[str, Any]:
        """Parse WS file for verification data (backward compatibility).

        Args:
            ws_path: Path to WS file

        Returns:
            Dict with scope_files, verification_commands, etc.
        """
        return parse_ws_file(ws_path)


    def verify(self, ws_id: str) -> VerificationResult:
        """Run all verification checks.

        Checks:
        1. All scope_files output exist
        2. All Verification commands pass
        3. Coverage meets threshold
        4. AC checkboxes accurate

        Args:
            ws_id: Workstream ID (e.g., "00-032-26")

        Returns:
            VerificationResult with all check results
        """
        checks: list[CheckResult] = []
        missing_files: list[str] = []
        failed_commands: list[str] = []

        # Find WS file
        ws_path = find_ws_file(ws_id, self.ws_dir)
        if not ws_path:
            return VerificationResult(
                ws_id=ws_id,
                passed=False,
                checks=[
                    CheckResult(
                        name="Find WS",
                        passed=False,
                        message=f"Workstream file not found: {ws_id}",
                        evidence=None,
                    )
                ],
                coverage_actual=None,
                missing_files=[],
                failed_commands=[],
            )

        # Parse WS file
        ws_data = parse_ws_file(ws_path)

        # Check 1: Verify output files exist
        file_checks = self.verify_output_files(ws_data)
        checks.extend(file_checks)
        missing_files = [c.message for c in file_checks if not c.passed]

        # Check 2: Run verification commands
        cmd_checks = self.verify_commands(ws_data)
        checks.extend(cmd_checks)
        failed_commands = [c.name for c in cmd_checks if not c.passed]

        # Check 3: Verify coverage
        coverage_check = self.verify_coverage(ws_data)
        if coverage_check:
            checks.append(coverage_check)

        # Determine overall pass/fail
        passed = all(c.passed for c in checks)

        return VerificationResult(
            ws_id=ws_id,
            passed=passed,
            checks=checks,
            coverage_actual=extract_coverage_value(coverage_check) if coverage_check else None,
            missing_files=missing_files,
            failed_commands=failed_commands,
        )
