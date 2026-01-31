"""TDD cycle execution runner."""

import subprocess
from dataclasses import dataclass
from enum import Enum
from pathlib import Path


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

    def __init__(self, project_dir: str | Path = "."):
        self._project = Path(project_dir)

    def red_phase(self, test_path: str) -> TDDResult:
        """Run test expecting failure.

        Args:
            test_path: Path to test file or directory

        Returns:
            TDDResult with success=True if test failed as expected
        """
        result = subprocess.run(
            ["pytest", test_path, "-v"],
            cwd=self._project,
            capture_output=True,
            text=True
        )

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
        result = subprocess.run(
            ["pytest", test_path, "-v"],
            cwd=self._project,
            capture_output=True,
            text=True
        )

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
        result = subprocess.run(
            ["pytest", test_path, "-v"],
            cwd=self._project,
            capture_output=True,
            text=True
        )

        success = result.returncode == 0
        return TDDResult(
            phase=Phase.REFACTOR,
            success=success,
            output=result.stdout + result.stderr,
            next_phase=None  # Cycle complete
        )
