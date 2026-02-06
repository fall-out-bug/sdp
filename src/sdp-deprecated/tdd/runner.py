"""TDD cycle execution runner."""

import subprocess
from dataclasses import dataclass
from enum import Enum
from pathlib import Path

from sdp.internal.retry import CommandRetryError, RetryConfig, run_with_retry


class Phase(Enum):
    """TDD cycle phases."""
    RED = "red"
    GREEN = "green"
    REFACTOR = "refactor"


@dataclass
class TDDResult:
    """Result of a TDD phase execution."""
    phase: Phase
    success: bool
    output: str
    next_phase: Phase | None = None


class TDDRunner:
    """Manages TDD cycle execution."""

    def __init__(
        self,
        project_dir: str | Path = ".",
        enable_retry: bool = False,
        max_retries: int = 3,
    ):
        """Initialize TDD runner.

        Args:
            project_dir: Project directory for test execution
            enable_retry: Enable automatic retry on transient failures
                (default: False - TDD failures should not be retried)
            max_retries: Maximum retry attempts per phase
        """
        self._project = Path(project_dir)
        self._enable_retry = enable_retry
        self._max_retries = max_retries

    def _run_command(
        self,
        command: list[str],
    ) -> subprocess.CompletedProcess[str]:
        """Run command with optional retry logic.

        Retry only applies to transient failures (e.g., subprocess errors,
        file locks, network issues). Test failures are NOT retried as they
        indicate real issues that need fixing.

        Args:
            command: Command to execute

        Returns:
            Completed process result
        """
        if self._enable_retry:
            try:
                # Use retry with exponential backoff
                result = run_with_retry(
                    command,
                    cwd=self._project,
                    config=RetryConfig(max_retries=self._max_retries),
                )
                # Convert to subprocess.CompletedProcess-like object
                return subprocess.CompletedProcess(
                    args=command,
                    returncode=result.exit_code,
                    stdout=result.stdout,
                    stderr=result.stderr,
                )
            except CommandRetryError as e:
                # All retries failed - return last attempt result
                return subprocess.CompletedProcess(
                    args=command,
                    returncode=e.last_exit_code,
                    stdout=e.last_output,
                    stderr="",
                )
        else:
            # Direct execution without retry
            return subprocess.run(
                command,
                cwd=self._project,
                capture_output=True,
                text=True,
            )

    def red_phase(self, test_path: str) -> TDDResult:
        """Run test expecting failure.

        Args:
            test_path: Path to test file or directory

        Returns:
            TDDResult with success=True if test failed as expected
        """
        result = self._run_command(["pytest", test_path, "-v"])

        # In RED phase, we EXPECT failure
        success = result.returncode != 0
        return TDDResult(
            phase=Phase.RED,
            success=success,
            output=result.stdout + result.stderr,
            next_phase=Phase.GREEN if success else None
        )

    def green_phase(self, test_path: str) -> TDDResult:
        """Run tests expecting success.

        Args:
            test_path: Path to test file or directory

        Returns:
            TDDResult with success=True if tests passed
        """
        result = self._run_command(["pytest", test_path, "-v"])

        success = result.returncode == 0
        return TDDResult(
            phase=Phase.GREEN,
            success=success,
            output=result.stdout + result.stderr,
            next_phase=Phase.REFACTOR if success else Phase.GREEN
        )

    def refactor_phase(self, test_path: str) -> TDDResult:
        """Run tests after refactoring.

        Args:
            test_path: Path to test file or directory

        Returns:
            TDDResult with success=True if tests still pass
        """
        result = self._run_command(["pytest", test_path, "-v"])

        success = result.returncode == 0
        return TDDResult(
            phase=Phase.REFACTOR,
            success=success,
            output=result.stdout + result.stderr,
            next_phase=None  # Cycle complete
        )
