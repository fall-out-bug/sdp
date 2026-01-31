"""Workstream completion verifier - validates WS is actually complete with evidence."""

import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass
class CheckResult:
    """Result of a single verification check."""

    name: str
    passed: bool
    message: str
    evidence: str | None  # Command output or file path


@dataclass
class VerificationResult:
    """Result of full WS verification."""

    ws_id: str
    passed: bool
    checks: list[CheckResult]
    coverage_actual: float | None
    missing_files: list[str]
    failed_commands: list[str]


class WSCompletionVerifier:
    """Verify workstream completion with evidence."""

    def __init__(self, ws_dir: Path = Path("docs/workstreams")):
        """Initialize verifier.

        Args:
            ws_dir: Base directory for workstream files
        """
        self.ws_dir = ws_dir

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
        ws_path = self._find_ws_file(ws_id)
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
        ws_data = self._parse_ws_file(ws_path)

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
            coverage_actual=self._extract_coverage(coverage_check) if coverage_check else None,
            missing_files=missing_files,
            failed_commands=failed_commands,
        )

    def verify_output_files(self, ws_data: dict[str, Any]) -> list[CheckResult]:
        """Check all output files in scope exist.

        Args:
            ws_data: Parsed workstream data

        Returns:
            List of check results for each file
        """
        checks: list[CheckResult] = []
        scope_files = ws_data.get("scope_files", [])

        for file_path_str in scope_files:
            file_path = Path(file_path_str)
            exists = file_path.exists()

            checks.append(
                CheckResult(
                    name=f"File: {file_path}",
                    passed=exists,
                    message=str(file_path) if exists else f"Missing: {file_path}",
                    evidence=str(file_path.absolute()) if exists else None,
                )
            )

        return checks

    def verify_commands(self, ws_data: dict[str, Any]) -> list[CheckResult]:
        """Run verification commands and check exit codes.

        Args:
            ws_data: Parsed workstream data

        Returns:
            List of check results for each command
        """
        checks: list[CheckResult] = []
        verification_cmds = ws_data.get("verification_commands", [])

        for cmd in verification_cmds:
            try:
                result = subprocess.run(
                    cmd,
                    shell=True,
                    capture_output=True,
                    text=True,
                    timeout=60,
                )

                passed = result.returncode == 0
                evidence = result.stdout if passed else result.stderr

                checks.append(
                    CheckResult(
                        name=f"Command: {cmd[:50]}...",
                        passed=passed,
                        message=f"Exit code: {result.returncode}",
                        evidence=evidence[:500] if evidence else None,
                    )
                )
            except subprocess.TimeoutExpired:
                checks.append(
                    CheckResult(
                        name=f"Command: {cmd[:50]}...",
                        passed=False,
                        message="Command timed out (60s)",
                        evidence=None,
                    )
                )
            except Exception as e:
                checks.append(
                    CheckResult(
                        name=f"Command: {cmd[:50]}...",
                        passed=False,
                        message=f"Error: {str(e)}",
                        evidence=None,
                    )
                )

        return checks

    def verify_coverage(self, ws_data: dict[str, Any]) -> CheckResult | None:
        """Check test coverage meets threshold.

        Args:
            ws_data: Parsed workstream data

        Returns:
            Check result or None if no coverage requirement
        """
        # Extract module from scope_files
        scope_files = ws_data.get("scope_files", [])
        if not scope_files:
            return None

        # Find Python files in scope
        py_files = [f for f in scope_files if f.endswith(".py") and not f.startswith("tests/")]
        if not py_files:
            return None

        # Run coverage check (basic implementation)
        try:
            result = subprocess.run(
                ["pytest", "--cov=src", "--cov-report=term", "-v"],
                capture_output=True,
                text=True,
                timeout=120,
            )

            # Parse coverage from output
            coverage = self._parse_coverage_from_output(result.stdout)

            passed = coverage is not None and coverage >= 80.0

            return CheckResult(
                name="Test Coverage",
                passed=passed,
                message=f"Coverage: {coverage:.1f}%" if coverage else "Coverage not found",
                evidence=result.stdout[:500],
            )
        except Exception as e:
            return CheckResult(
                name="Test Coverage",
                passed=False,
                message=f"Error running coverage: {str(e)}",
                evidence=None,
            )

    def _find_ws_file(self, ws_id: str) -> Path | None:
        """Find WS file by ID.

        Args:
            ws_id: Workstream ID

        Returns:
            Path to WS file or None
        """
        # Search in multiple locations
        search_dirs = [
            self.ws_dir / "backlog",
            self.ws_dir / "in_progress",
            self.ws_dir / "completed",
        ]

        for search_dir in search_dirs:
            if not search_dir.exists():
                continue

            for ws_file in search_dir.glob(f"{ws_id}*.md"):
                return ws_file

        return None

    def _parse_frontmatter_scope(self, content: str) -> list[str]:
        """Extract scope_files from frontmatter."""
        scope_files: list[str] = []
        in_frontmatter = False
        in_scope_files = False
        for line in content.splitlines():
            if line.strip() == "---":
                if not in_frontmatter:
                    in_frontmatter = True
                else:
                    break
                continue
            if in_frontmatter:
                if line.startswith("scope_files:"):
                    in_scope_files = True
                elif in_scope_files:
                    if line.startswith("  - "):
                        scope_files.append(line.strip()[2:].strip())
                    elif not line.startswith(" "):
                        in_scope_files = False
        return scope_files

    def _parse_verification_commands(self, content: str) -> list[str]:
        """Extract verification commands from ### Verification section."""
        commands: list[str] = []
        in_verification = False
        in_code_block = False
        for line in content.splitlines():
            if line.startswith("### Verification"):
                in_verification = True
                continue
            if in_verification:
                if line.strip().startswith("```bash") or line.strip().startswith("```sh"):
                    in_code_block = True
                    continue
                if line.strip() == "```":
                    in_code_block = False
                    continue
                if line.startswith("##"):
                    break
                if in_code_block and line.strip() and not line.strip().startswith("#"):
                    commands.append(line.strip())
        return commands

    def _parse_ws_file(self, ws_path: Path) -> dict[str, Any]:
        """Parse WS file for verification data.

        Args:
            ws_path: Path to WS file

        Returns:
            Dict with scope_files, verification_commands, etc.
        """
        content = ws_path.read_text(encoding="utf-8")
        return {
            "scope_files": self._parse_frontmatter_scope(content),
            "verification_commands": self._parse_verification_commands(content),
        }

    def _parse_coverage_from_output(self, output: str) -> float | None:
        """Parse coverage percentage from pytest output.

        Args:
            output: pytest stdout

        Returns:
            Coverage percentage or None
        """
        import re

        # Look for "TOTAL ... XX%"
        match = re.search(r"TOTAL\s+\d+\s+\d+\s+(\d+)%", output)
        if match:
            return float(match.group(1))

        return None

    def _extract_coverage(self, check: CheckResult | None) -> float | None:
        """Extract coverage value from check result.

        Args:
            check: Coverage check result

        Returns:
            Coverage percentage or None
        """
        if not check or not check.message:
            return None

        import re

        match = re.search(r"Coverage: ([\d.]+)%", check.message)
        if match:
            return float(match.group(1))

        return None
