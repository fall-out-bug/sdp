"""Quality-related errors.

Covers: Beads integration, coverage violations, quality gate failures,
workstream validation.
"""

from .base import ErrorCategory, SDPError


class BeadsNotFoundError(SDPError):
    """Beads task data not found."""

    def __init__(
        self,
        task_id: str,
        search_paths: list[str] | None = None,
    ) -> None:
        """Initialize BeadsNotFoundError.

        Args:
            task_id: The task ID that was not found
            search_paths: Paths that were searched
        """
        super().__init__(
            category=ErrorCategory.BEADS,
            message=f"Beads task '{task_id}' not found",
            remediation=(
                "1. Check the task ID is correct\n"
                "2. Ensure Beads server is running: beads status\n"
                "3. Verify task exists: beads tasks\n"
                "4. Check BEADS_HOME environment variable"
            ),
            docs_url="https://docs.sdp.dev/troubleshooting#beads-not-found",
            context={"task_id": task_id, "search_paths": search_paths},
        )


class CoverageTooLowError(SDPError):
    """Test coverage below required threshold."""

    def __init__(
        self,
        coverage_pct: float,
        required_pct: float,
        module: str,
        missing_files: list[str] | None = None,
    ) -> None:
        """Initialize CoverageTooLowError.

        Args:
            coverage_pct: Actual coverage percentage
            required_pct: Required coverage percentage
            module: Module being tested
            missing_files: Files with missing coverage
        """
        super().__init__(
            category=ErrorCategory.COVERAGE,
            message=(
                f"Coverage {coverage_pct:.1f}% is below required {required_pct:.1f}% "
                f"for module '{module}'"
            ),
            remediation=(
                f"1. Add tests for uncovered code\n"
                f"2. Run: pytest --cov={module} --cov-report=term-missing\n"
                f"3. Focus on files with 0% coverage\n"
                f"4. Target: ≥{required_pct:.0f}% coverage"
            ),
            docs_url="https://docs.sdp.dev/quality-gates#coverage",
            context={
                "actual_coverage": f"{coverage_pct:.1f}%",
                "required_coverage": f"{required_pct:.1f}%",
                "module": module,
                "missing_files": missing_files or [],
            },
        )


class QualityGateViolationError(SDPError):
    """Quality gate check failed."""

    def __init__(
        self,
        gate_name: str,
        violations: list[str],
        severity: str = "error",
    ) -> None:
        """Initialize QualityGateViolationError.

        Args:
            gate_name: Name of the failed gate
            violations: List of violation messages
            severity: 'error', 'warning', or 'critical'
        """
        super().__init__(
            category=ErrorCategory.VALIDATION,
            message=f"Quality gate '{gate_name}' failed with {len(violations)} violation(s)",
            remediation=(
                f"1. Fix all violations listed below\n"
                f"2. Run quality gate check again\n"
                f"3. See docs for gate requirements: {gate_name}\n"
                f"4. Severity: {severity}"
            ),
            docs_url="https://docs.sdp.dev/quality-gates",
            context={
                "gate": gate_name,
                "severity": severity,
                "violations": violations,
            },
        )


class WorkstreamValidationError(SDPError):
    """Workstream file validation failed."""

    def __init__(
        self,
        ws_id: str,
        errors: list[str],
        file_path: str | None = None,
    ) -> None:
        """Initialize WorkstreamValidationError.

        Args:
            ws_id: Workstream ID
            errors: List of validation error messages
            file_path: Path to workstream file
        """
        super().__init__(
            category=ErrorCategory.VALIDATION,
            message=f"Workstream validation failed: {ws_id}",
            remediation=(
                "Fix workstream structure:\n"
                + "\n".join(f"  • {e}" for e in errors[:5])
                + ("\n  ..." if len(errors) > 5 else "")
            ),
            docs_url="https://docs.sdp.dev/workstreams#template",
            context={
                "ws_id": ws_id,
                "file_path": file_path,
                "errors": errors,
            },
        )
