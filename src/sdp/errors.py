"""Comprehensive error framework for SDP.

This module provides structured error types with contextual information,
remediation guidance, and documentation links for common SDP errors.
"""

from dataclasses import dataclass
from enum import Enum
from typing import Any


class ErrorCategory(Enum):
    """Categories of SDP errors for better organization."""

    VALIDATION = "validation"
    BUILD = "build"
    TEST = "test"
    CONFIGURATION = "configuration"
    DEPENDENCY = "dependency"
    HOOK = "hook"
    ARTIFACT = "artifact"
    BEADS = "beads"
    COVERAGE = "coverage"


@dataclass(frozen=True)
class SDPError(Exception):
    """Base error class for all SDP errors.

    Provides structured error information including:
    - Category: Error classification
    - Message: Human-readable error description
    - Remediation: How to fix the error
    - Docs URL: Link to documentation (optional)
    - Context: Additional error-specific data (optional)

    Example:
        raise SDPError(
            category=ErrorCategory.VALIDATION,
            message="Workstream file not found",
            remediation="Check the WS-ID and ensure file exists in workstreams/backlog/",
            docs_url="https://docs.sdp.dev/troubleshooting#ws-not-found",
            context={"ws_id": "WS-001-01", "search_paths": [...]}
        )
    """

    category: ErrorCategory
    message: str
    remediation: str
    docs_url: str | None = None
    context: dict[str, Any] | None = None

    def __str__(self) -> str:
        """Format error for display."""
        lines = [
            f"❌ {self.category.value.upper()} Error",
            f"   {self.message}",
            "",
            "   💡 Remediation:",
            f"   {self.remediation}",
        ]

        if self.docs_url:
            lines.extend([
                "",
                "   📚 Documentation:",
                f"   {self.docs_url}",
            ])

        if self.context:
            lines.extend([
                "",
                "   🔍 Context:",
            ])
            for key, value in self.context.items():
                lines.append(f"   {key}: {value}")

        return "\n".join(lines)


# =============================================================================
# Predefined Error Types
# =============================================================================


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
            message=f"Workstream '{ws_id}' validation failed with {len(errors)} error(s)",
            remediation=(
                "1. Fix all validation errors listed below\n"
                "2. Ensure required sections exist (Goal, Acceptance Criteria)\n"
                "3. Check frontmatter format\n"
                "4. Validate against template: docs/workstreams/template.md"
            ),
            docs_url="https://docs.sdp.dev/workstreams#validation",
            context={
                "ws_id": ws_id,
                "file_path": file_path,
                "errors": errors,
            },
        )


class ConfigurationError(SDPError):
    """SDP configuration error."""

    def __init__(
        self,
        config_file: str,
        errors: list[str],
        missing_keys: list[str] | None = None,
    ) -> None:
        """Initialize ConfigurationError.

        Args:
            config_file: Configuration file path
            errors: List of configuration errors
            missing_keys: Missing required configuration keys
        """
        super().__init__(
            category=ErrorCategory.CONFIGURATION,
            message=f"Configuration error in '{config_file}': {len(errors)} issue(s)",
            remediation=(
                "1. Review configuration file format\n"
                "2. Add missing required keys\n"
                "3. Check syntax (TOML/YAML/JSON)\n"
                "4. See docs for config schema"
            ),
            docs_url="https://docs.sdp.dev/configuration",
            context={
                "config_file": config_file,
                "errors": errors,
                "missing_keys": missing_keys or [],
            },
        )


class DependencyNotFoundError(SDPError):
    """Required dependency not found."""

    def __init__(
        self,
        dependency: str,
        ws_id: str | None = None,
        available_ws: list[str] | None = None,
    ) -> None:
        """Initialize DependencyNotFoundError.

        Args:
            dependency: Missing dependency identifier
            ws_id: Workstream that depends on this
            available_ws: List of available workstreams
        """
        ws_context = f" for workstream '{ws_id}'" if ws_id else ""
        super().__init__(
            category=ErrorCategory.DEPENDENCY,
            message=f"Dependency '{dependency}' not found{ws_context}",
            remediation=(
                f"1. Complete dependency workstream first: {dependency}\n"
                f"2. Check INDEX.md for workstream status\n"
                f"3. Verify dependency ID is correct\n"
                f"4. Update workstream frontmatter if needed"
            ),
            docs_url="https://docs.sdp.dev/workstreams#dependencies",
            context={
                "dependency": dependency,
                "ws_id": ws_id,
                "available_workstreams": available_ws or [],
            },
        )


class HookExecutionError(SDPError):
    """Git hook or build hook execution failed."""

    def __init__(
        self,
        hook_name: str,
        stage: str,
        output: str,
        exit_code: int,
    ) -> None:
        """Initialize HookExecutionError.

        Args:
            hook_name: Name of the hook that failed
            stage: pre-commit, post-build, etc.
            output: Hook output/stderr
            exit_code: Hook exit code
        """
        super().__init__(
            category=ErrorCategory.HOOK,
            message=f"Hook '{hook_name}' failed during {stage} (exit code: {exit_code})",
            remediation=(
                "1. Review hook output above for specific errors\n"
                "2. Fix the issue that caused hook to fail\n"
                "3. Test hook manually: hooks/{hook_name}.sh\n"
                "4. Bypass with SKIP_CHECK=1 if needed (not recommended)"
            ),
            docs_url="https://docs.sdp.dev/hooks#troubleshooting",
            context={
                "hook": hook_name,
                "stage": stage,
                "exit_code": exit_code,
                "output": output[:500],  # Truncate long output
            },
        )


class TestFailureError(SDPError):
    """Test execution failed."""

    def __init__(
        self,
        test_command: str,
        failed_tests: list[str],
        total_tests: int,
        passed_tests: int,
    ) -> None:
        """Initialize TestFailureError.

        Args:
            test_command: Command that was run
            failed_tests: List of failed test names
            total_tests: Total number of tests
            passed_tests: Number of passed tests
        """
        super().__init__(
            category=ErrorCategory.TEST,
            message=(
                f"Tests failed: {passed_tests}/{total_tests} passed, "
                f"{len(failed_tests)} failed"
            ),
            remediation=(
                "1. Run tests with verbose output: pytest -v\n"
                "2. Fix failing tests one by one\n"
                "3. Check for regression: git diff\n"
                "4. Ensure all acceptance criteria are tested"
            ),
            docs_url="https://docs.sdp.dev/testing#debugging",
            context={
                "command": test_command,
                "passed": passed_tests,
                "failed": len(failed_tests),
                "total": total_tests,
                "failed_tests": failed_tests[:10],  # Show first 10
            },
        )


class BuildValidationError(SDPError):
    """Build validation check failed."""

    def __init__(
        self,
        ws_id: str,
        stage: str,
        check_name: str,
        details: str,
    ) -> None:
        """Initialize BuildValidationError.

        Args:
            ws_id: Workstream being built
            stage: pre-build or post-build
            check_name: Name of the failed check
            details: Check failure details
        """
        super().__init__(
            category=ErrorCategory.BUILD,
            message=f"Build validation failed: {check_name} ({stage})",
            remediation=(
                f"1. Fix the issue: {details}\n"
                f"2. Re-run the check: hooks/{stage}.sh {ws_id}\n"
                f"3. Verify all pre-build checks pass before starting\n"
                f"4. Verify all post-build checks pass before committing"
            ),
            docs_url="https://docs.sdp.dev/building#validation",
            context={
                "ws_id": ws_id,
                "stage": stage,
                "check": check_name,
                "details": details,
            },
        )


class ArtifactValidationError(SDPError):
    """Workstream artifact validation failed."""

    def __init__(
        self,
        artifact_type: str,
        artifact_path: str,
        errors: list[str],
    ) -> None:
        """Initialize ArtifactValidationError.

        Args:
            artifact_type: Type of artifact (code, test, docs)
            artifact_path: Path to the artifact
            errors: List of validation errors
        """
        super().__init__(
            category=ErrorCategory.ARTIFACT,
            message=(
                f"Artifact validation failed for {artifact_type}: "
                f"{len(errors)} error(s)"
            ),
            remediation=(
                "1. Review artifact quality requirements\n"
                "2. Fix all validation errors listed below\n"
                "3. Ensure artifact meets SDP standards\n"
                "4. See docs for artifact checklist"
            ),
            docs_url="https://docs.sdp.dev/artifacts#validation",
            context={
                "type": artifact_type,
                "path": artifact_path,
                "errors": errors,
            },
        )


# =============================================================================
# Error Formatting Utilities
# =============================================================================


def format_error_for_terminal(error: Exception) -> str:
    """Format any exception for terminal display.

    Args:
        error: Exception to format

    Returns:
        Formatted error message
    """
    if isinstance(error, SDPError):
        return str(error)

    # Standard Python exception
    return f"❌ Error: {type(error).__name__}\n   {error}"


def format_error_for_json(error: Exception) -> dict[str, Any]:
    """Format exception as JSON-serializable dict.

    Args:
        error: Exception to format

    Returns:
        Dictionary with error details
    """
    if isinstance(error, SDPError):
        return {
            "category": error.category.value,
            "message": error.message,
            "remediation": error.remediation,
            "docs_url": error.docs_url,
            "context": error.context,
            "type": type(error).__name__,
        }

    return {
        "type": type(error).__name__,
        "message": str(error),
        "category": "unknown",
        "remediation": "Review error message and stack trace",
    }
