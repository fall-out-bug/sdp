"""
Execution Report Generator for SDP Workstreams.

Auto-generates execution reports on WS completion with:
- Duration tracking
- LOC (lines of code) changes
- Test coverage statistics
- Deviations from plan
"""

from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
from typing import Optional


@dataclass
class ExecutionStats:
    """Statistics collected during workstream execution."""

    duration_minutes: float
    files_changed: int
    loc_added: int
    loc_removed: int
    coverage_pct: float
    tests_passed: int
    tests_failed: int
    deviations: list[str]


class ReportGenerator:
    """
    Generates execution reports for completed workstreams.

    Usage:
        generator = ReportGenerator(ws_id="WS-001-01", ws_path="/path/to/ws.md")
        generator.start_timer()
        # ... execute workstream ...
        stats = generator.collect_stats()
        generator.append_report(stats, executed_by="developer")
    """

    def __init__(self, ws_id: str, ws_path: Optional[Path] = None) -> None:
        """
        Initialize report generator.

        Args:
            ws_id: Workstream ID (e.g., "WS-001-01")
            ws_path: Path to workstream file (auto-detected if None when needed)
        """
        self.ws_id = ws_id
        self._ws_path_provided = ws_path is not None
        self._ws_path = ws_path
        self.start_time: Optional[datetime] = None
        self.end_time: Optional[datetime] = None

    def _find_ws_file(self) -> Path:
        """Auto-detect workstream file from common locations."""
        search_paths = [
            Path("docs/workstreams/in_progress"),
            Path("docs/workstreams/backlog"),
            Path("docs/workstreams/completed"),
        ]

        for base_path in search_paths:
            if base_path.exists():
                candidates = list(base_path.glob(f"{self.ws_id}*.md"))
                if candidates:
                    return candidates[0]

        raise FileNotFoundError(f"Cannot find workstream file for {self.ws_id}")

    @property
    def ws_path(self) -> Optional[Path]:
        """Get workstream path, auto-detecting if needed."""
        if self._ws_path is None and not self._ws_path_provided:
            self._ws_path = self._find_ws_file()
        return self._ws_path

    @ws_path.setter
    def ws_path(self, value: Optional[Path]) -> None:
        """Set workstream path."""
        self._ws_path = value
        if value is not None:
            self._ws_path_provided = True

    def start_timer(self) -> None:
        """Start execution timer."""
        self.start_time = datetime.now()

    def stop_timer(self) -> None:
        """Stop execution timer."""
        self.end_time = datetime.now()

    def get_duration_minutes(self) -> float:
        """
        Get execution duration in minutes.

        Returns:
            Duration in minutes (0 if timer not started/stopped)

        Raises:
            ValueError: If timer wasn't started and stopped
        """
        if not self.start_time or not self.end_time:
            return 0.0

        delta = self.end_time - self.start_time
        return delta.total_seconds() / 60.0

    def collect_stats(
        self,
        files_changed: Optional[list[tuple[str, str, int]]] = None,
        coverage_pct: Optional[float] = None,
        tests_passed: int = 0,
        tests_failed: int = 0,
        deviations: Optional[list[str]] = None,
    ) -> ExecutionStats:
        """
        Collect execution statistics.

        Args:
            files_changed: List of (file_path, action, loc_added) tuples
            coverage_pct: Test coverage percentage
            tests_passed: Number of tests passed
            tests_failed: Number of tests failed
            deviations: List of deviations from plan

        Returns:
            ExecutionStats object with collected data
        """
        duration = self.get_duration_minutes()

        # Calculate LOC changes
        loc_added = sum(loc for _, _, loc in (files_changed or []))
        loc_removed = 0  # Would need git diff to calculate accurately

        return ExecutionStats(
            duration_minutes=duration,
            files_changed=len(files_changed) if files_changed else 0,
            loc_added=loc_added,
            loc_removed=loc_removed,
            coverage_pct=coverage_pct or 0.0,
            tests_passed=tests_passed,
            tests_failed=tests_failed,
            deviations=deviations or [],
        )

    def generate_report_markdown(
        self,
        stats: ExecutionStats,
        executed_by: str,
        commit_hash: Optional[str] = None,
    ) -> str:
        """
        Generate markdown report from statistics.

        Args:
            stats: Execution statistics
            executed_by: Who executed the workstream
            commit_hash: Git commit hash (optional)

        Returns:
            Markdown formatted report
        """
        report_date = datetime.now().strftime("%Y-%m-%d")

        # Build files changed table
        files_table = self._get_files_changed_markdown()

        markdown = f"""

---

## Execution Report

**Executed by:** {executed_by}
**Date:** {report_date}
**Duration:** {stats.duration_minutes:.1f} minutes

### Goal Status
- [ ] AC1-AC3 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
{files_table}

### Statistics
- **Files Changed:** {stats.files_changed}
- **Lines Added:** {stats.loc_added}
- **Lines Removed:** {stats.loc_removed}
- **Test Coverage:** {stats.coverage_pct:.1f}%
- **Tests Passed:** {stats.tests_passed}
- **Tests Failed:** {stats.tests_failed}

### Deviations from Plan
"""
        if stats.deviations:
            for deviation in stats.deviations:
                markdown += f"- {deviation}\n"
        else:
            markdown += "- None (followed plan exactly)\n"

        markdown += f"""
### Commit
{commit_hash or "______"}

---

"""
        return markdown

    def _get_files_changed_markdown(self) -> str:
        """
        Get markdown table of files changed.

        Returns:
            Markdown table rows (empty if no git detected)
        """
        try:
            import subprocess

            # Get last commit's changed files
            result = subprocess.run(
                ["git", "diff", "--stat", "HEAD~1", "HEAD"],
                capture_output=True,
                text=True,
                check=False,
            )

            if result.returncode != 0:
                return "|      |        |     |"

            # Parse git diff --stat output
            lines = []
            for line in result.stdout.strip().split("\n"):
                if line.strip() and not line.startswith(" ") and "|" in line:
                    parts = line.split("|")
                    if len(parts) >= 2:
                        filename = parts[0].strip()
                        stats = parts[1].strip()
                        # Extract +/- counts if available
                        if "+" in stats:
                            loc = stats.split()[0]
                            lines.append(f"| {filename} | modified | {loc} |")

            return "\n".join(lines) if lines else "|      |        |     |"

        except Exception:
            return "|      |        |     |"

    def append_report(
        self,
        stats: ExecutionStats,
        executed_by: str,
        commit_hash: Optional[str] = None,
    ) -> None:
        """
        Append execution report to workstream file.

        Args:
            stats: Execution statistics
            executed_by: Who executed the workstream
            commit_hash: Git commit hash (optional)

        Raises:
            IOError: If cannot write to workstream file
        """
        if not self.ws_path or not self.ws_path.exists():
            raise FileNotFoundError(f"Workstream file not found: {self.ws_path}")

        # Generate report markdown
        report = self.generate_report_markdown(stats, executed_by, commit_hash)

        # Append to file
        with open(self.ws_path, "a") as f:
            f.write(report)

    def generate(
        self,
        executed_by: str,
        files_changed: Optional[list[tuple[str, str, int]]] = None,
        coverage_pct: Optional[float] = None,
        tests_passed: int = 0,
        tests_failed: int = 0,
        deviations: Optional[list[str]] = None,
        commit_hash: Optional[str] = None,
    ) -> str:
        """
        Complete report generation workflow.

        Args:
            executed_by: Who executed the workstream
            files_changed: List of (file_path, action, loc_added) tuples
            coverage_pct: Test coverage percentage
            tests_passed: Number of tests passed
            tests_failed: Number of tests failed
            deviations: List of deviations from plan
            commit_hash: Git commit hash

        Returns:
            Generated report markdown
        """
        self.stop_timer()
        stats = self.collect_stats(
            files_changed=files_changed,
            coverage_pct=coverage_pct,
            tests_passed=tests_passed,
            tests_failed=tests_failed,
            deviations=deviations,
        )
        return self.generate_report_markdown(stats, executed_by, commit_hash)
