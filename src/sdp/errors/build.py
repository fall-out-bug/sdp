"""Build and test errors.

Covers: Test failures, build validation errors, artifact issues.
"""

from .base import ErrorCategory, SDPError


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
