"""Command retry with exponential backoff.

This module provides retry logic for subprocess commands that may fail
transiently. It implements exponential backoff and integrates with telemetry.
"""

import subprocess
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Callable, Sequence


@dataclass
class RetryConfig:
    """Configuration for command retry behavior.

    Attributes:
        max_retries: Maximum number of retry attempts (default: 3)
        base_delay: Base delay in seconds before first retry (default: 1.0)
        backoff_factor: Multiplier for exponential backoff (default: 2.0)
    """

    max_retries: int = 3
    base_delay: float = 1.0
    backoff_factor: float = 2.0


@dataclass
class RetryResult:
    """Result of a command execution with retry.

    Attributes:
        success: Whether command succeeded
        exit_code: Process exit code
        stdout: Standard output as string
        stderr: Standard error as string
        attempts: Number of execution attempts
        total_delay: Total time spent in delays
    """

    success: bool
    exit_code: int
    stdout: str
    stderr: str
    attempts: int
    total_delay: float


class CommandRetryError(Exception):
    """Exception raised when command fails after all retries.

    Attributes:
        command: Command that failed
        last_exit_code: Last exit code received
        attempts: Total number of attempts
        total_delay: Total delay time
        last_output: Last output (stdout + stderr)
    """

    def __init__(
        self,
        command: Sequence[str],
        last_exit_code: int,
        attempts: int,
        total_delay: float,
        last_output: str = "",
    ) -> None:
        """Initialize retry error.

        Args:
            command: Command that failed
            last_exit_code: Last exit code received
            attempts: Total number of attempts
            total_delay: Total delay time
            last_output: Last output
        """
        self.command = command
        self.last_exit_code = last_exit_code
        self.attempts = attempts
        self.total_delay = total_delay
        self.last_output = last_output
        super().__init__(f"Command {list(command)} failed after {attempts} attempts")


def run_with_retry(
    command: Sequence[str],
    config: RetryConfig | None = None,
    cwd: str | Path | None = None,
    timeout: float | None = None,
    on_retry: Callable[[RetryResult], None] | None = None,
) -> RetryResult:
    """Run command with exponential backoff retry.

    Args:
        command: Command to execute (list of strings)
        config: Retry configuration (uses default if None)
        cwd: Working directory for command
        timeout: Timeout in seconds for each attempt
        on_retry: Optional callback after each attempt (for telemetry)

    Returns:
        RetryResult with execution details

    Raises:
        CommandRetryError: If command fails after all retries
        subprocess.SubprocessError: For subprocess exceptions (not caught)

    Examples:
        >>> result = run_with_retry(["pytest", "tests/"])
        >>> if result.success:
        ...     print(f"Passed after {result.attempts} attempts")

        >>> # Custom retry config
        >>> config = RetryConfig(max_retries=5, base_delay=2.0)
        >>> result = run_with_retry(["npm", "test"], config=config)

        >>> # With telemetry callback
        >>> def log_telemetry(result: RetryResult) -> None:
        ...     print(f"Attempt {result.attempts}: success={result.success}")
        >>> result = run_with_retry(["make", "test"], on_retry=log_telemetry)
    """
    if config is None:
        config = RetryConfig()

    total_delay = 0.0
    last_result: RetryResult | None = None

    for attempt in range(config.max_retries + 1):
        try:
            # Execute command
            proc_result = subprocess.run(
                command,
                cwd=cwd,
                timeout=timeout,
                capture_output=True,
                text=True,
            )

            # Create result for this attempt
            result = RetryResult(
                success=proc_result.returncode == 0,
                exit_code=proc_result.returncode,
                stdout=proc_result.stdout,
                stderr=proc_result.stderr,
                attempts=attempt + 1,
                total_delay=total_delay,
            )

            # Call telemetry callback if provided
            if on_retry:
                on_retry(result)

            # If successful, return result
            if result.success:
                return result

            # Store for error reporting
            last_result = result

            # Calculate delay if not the last attempt
            if attempt < config.max_retries:
                delay = config.base_delay * (config.backoff_factor**attempt)
                time.sleep(delay)
                total_delay += delay

        except subprocess.SubprocessError:
            # Re-raise subprocess exceptions (FileNotFoundError, etc.)
            raise

    # All retries failed
    if last_result:
        raise CommandRetryError(
            command=command,
            last_exit_code=last_result.exit_code,
            attempts=last_result.attempts,
            total_delay=last_result.total_delay,
            last_output=last_result.stdout + last_result.stderr,
        )

    # Should not reach here, but satisfy type checker
    raise RuntimeError("Retry loop completed without result")
