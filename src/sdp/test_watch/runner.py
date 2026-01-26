"""Test runner for watch mode."""

import logging
import re
import subprocess
from pathlib import Path

from ..dashboard.sources.test_runner import TestRunner, TestResults

logger = logging.getLogger(__name__)


class WatchTestRunner(TestRunner):
    """Test runner specifically for watch mode."""

    def __init__(
        self,
        project_dir: str | Path,
        coverage: bool = True,
        pattern: str | None = None,
    ) -> None:
        """Initialize watch test runner.

        Args:
            project_dir: Project root directory
            coverage: Whether to collect coverage
            pattern: Test pattern filter
        """
        super().__init__(project_dir, coverage)
        self._pattern = pattern

    def run(self) -> TestResults:
        """Run tests and return results.

        Returns:
            TestResults with outcome of test run
        """
        tests_dir = self._project / "tests"
        if not tests_dir.exists():
            return TestResults(status="no_tests")

        cmd = ["pytest", "-v", "--tb=short"]

        if self._pattern:
            cmd.extend(["-k", self._pattern])

        if self._coverage:
            cmd.extend(["--cov=", "--cov-report=json"])

        try:
            result = subprocess.run(
                cmd,
                cwd=self._project,
                capture_output=True,
                text=True,
                timeout=300,
            )
            return self._parse_output(result.stdout, result.stderr)
        except subprocess.TimeoutExpired:
            return TestResults(status="error", error_message="Tests timed out")
        except FileNotFoundError:
            return TestResults(status="no_tests")
        except Exception as e:
            return TestResults(status="error", error_message=str(e))

    def run_affected(self, changed_file: str) -> TestResults:
        """Run tests affected by a changed file.

        Args:
            changed_file: Path to changed file

        Returns:
            TestResults with outcome of test run
        """
        # Try to determine which tests to run based on changed file
        changed_path = Path(changed_file)
        module_name = changed_path.stem

        # If a test file changed, run just that file
        if "tests" in changed_path.parts:
            cmd = ["pytest", "-v", "--tb=short", str(changed_path)]
        else:
            # Run tests that might be affected
            # Look for test files with matching names
            test_file = self._project / "tests" / f"test_{module_name}.py"
            if test_file.exists():
                cmd = ["pytest", "-v", "--tb=short", str(test_file)]
            else:
                # Fall back to running all tests
                return self.run()

        if self._coverage:
            cmd.extend(["--cov=", "--cov-report=json"])

        try:
            result = subprocess.run(
                cmd,
                cwd=self._project,
                capture_output=True,
                text=True,
                timeout=300,
            )
            return self._parse_output(result.stdout, result.stderr)
        except subprocess.TimeoutExpired:
            return TestResults(status="error", error_message="Tests timed out")
        except Exception as e:
            return TestResults(status="error", error_message=str(e))
