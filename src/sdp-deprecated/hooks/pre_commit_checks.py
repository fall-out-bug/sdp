"""Pre-commit check functions - extracted for testability and LOC limit."""

import re
import subprocess
import sys
from pathlib import Path

from sdp.hooks.common import CheckResult


def repo_root() -> Path:
    """Get repository root."""
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        capture_output=True,
        text=True,
        check=True,
    )
    return Path(result.stdout.strip())


def staged_files() -> list[str]:
    """Get list of staged file paths."""
    result = subprocess.run(
        ["git", "diff", "--cached", "--name-only", "--diff-filter=ACM"],
        capture_output=True,
        text=True,
        check=True,
    )
    return [f for f in result.stdout.strip().split("\n") if f]


def check_branch() -> CheckResult:
    """Check not committing directly to main/master."""
    result = subprocess.run(
        ["git", "branch", "--show-current"],
        capture_output=True,
        text=True,
        check=True,
    )
    branch = result.stdout.strip()
    if branch in ("main", "master"):
        return CheckResult(
            passed=False,
            message=f"Committing directly to {branch}",
            violations=[(Path("."), None, "Create feature branch first")],
        )
    return CheckResult(passed=True, message=f"Branch: {branch}", violations=[])


def check_time_estimates(files: list[str]) -> CheckResult:
    """Check for time-based estimates in WS files."""
    ws_files = [f for f in files if "workstreams/" in f and f.endswith(".md")]
    if not ws_files:
        return CheckResult(passed=True, message="No WS files staged", violations=[])

    result = subprocess.run(
        ["git", "diff", "--cached", "--"] + ws_files,
        capture_output=True,
        text=True,
        check=True,
    )
    diff = result.stdout
    patterns = [
        r"дн[яей]",
        r"час[ов]",
        r"недел",
        r"\bday\b",
        r"\bhour\b",
        r"\bweek\b",
    ]
    exclude = re.compile(r"elapsed|duration|sla|telemetry", re.I)
    for line in diff.split("\n"):
        if not line.startswith("+"):
            continue
        for pat in patterns:
            if re.search(pat, line, re.I) and not exclude.search(line):
                return CheckResult(
                    passed=False,
                    message="Time estimates found in WS files",
                    violations=[(Path("."), None, "Remove time-based estimates")],
                )
    return CheckResult(passed=True, message="No time estimates", violations=[])


def check_tech_debt(files: list[str]) -> CheckResult:
    """Check for tech debt markers."""
    code_files = [
        f for f in files
        if f.endswith((".py", ".md", ".yml", ".yaml", ".json")) and not f.endswith(".sh")
    ]
    if not code_files:
        return CheckResult(passed=True, message="No code files staged", violations=[])

    result = subprocess.run(
        ["git", "diff", "--cached", "--"] + code_files,
        capture_output=True,
        text=True,
        check=True,
    )
    patterns = re.compile(
        r"tech.?debt|сделаем.?потом|временн.*решени|later.*fix",
        re.I,
    )
    exclude = re.compile(
        r"no.?tech.?debt|⛔|запрещено|forbidden|zero\.tech|tech\.debt.*0",
        re.I,
    )
    for line in result.stdout.split("\n"):
        if not line.startswith("+") or exclude.search(line):
            continue
        if patterns.search(line):
            return CheckResult(
                passed=False,
                message="Tech debt markers found",
                violations=[(Path("."), None, "Fix the issue now")],
            )
    return CheckResult(passed=True, message="No tech debt markers", violations=[])


def check_python_bare_except(files: list[str]) -> CheckResult:
    """Check for bare except and except: pass."""
    py_files = [f for f in files if f.endswith(".py")]
    if not py_files:
        return CheckResult(passed=True, message="No Python files staged", violations=[])

    result = subprocess.run(
        ["git", "diff", "--cached", "--"] + py_files,
        capture_output=True,
        text=True,
        check=True,
    )
    lines = result.stdout.split("\n")
    for i, line in enumerate(lines):
        if line.startswith("+") and re.search(r"except\s*:", line):
            return CheckResult(
                passed=False,
                message="Bare except found",
                violations=[(Path("."), None, "Use specific exception types")],
            )
        if line.startswith("+") and "except" in line and i + 1 < len(lines):
            next_line = lines[i + 1]
            if next_line.startswith("+") and "pass" in next_line:
                return CheckResult(
                    passed=False,
                    message="except: pass found",
                    violations=[(Path("."), None, "Handle exception explicitly")],
                )
    return CheckResult(passed=True, message="Python checks passed", violations=[])


def run_script(repo_root_path: Path, script: str, args: list[str]) -> tuple[bool, str]:
    """Run script and return (success, output)."""
    script_path = repo_root_path / script
    if not script_path.exists():
        return True, f"  {script} not found (skipped)"
    proc = subprocess.run(
        [sys.executable, str(script_path)] + args,
        cwd=repo_root_path,
        capture_output=True,
        text=True,
    )
    return proc.returncode == 0, proc.stdout + proc.stderr


def check_ws_format(files: list[str], repo_root_path: Path) -> CheckResult:
    """Check WS file format for new workstreams."""
    new_ws = [
        f for f in files
        if "workstreams/backlog/" in f and f.endswith(".md")
        and ("WS-" in f or re.search(r"\d{2}-\d{3}-\d{2}", f))
    ]
    if not new_ws:
        return CheckResult(passed=True, message="No new WS files staged", violations=[])

    for ws_file in new_ws:
        result = subprocess.run(
            ["git", "show", f":{ws_file}"],
            capture_output=True,
            text=True,
            cwd=repo_root_path,
        )
        if result.returncode != 0:
            continue
        content = result.stdout
        if "### 🎯" not in content:
            return CheckResult(
                passed=False,
                message=f"Missing Goal section in {ws_file}",
                violations=[],
            )
        if "Acceptance Criteria" not in content:
            return CheckResult(
                passed=False,
                message=f"Missing Acceptance Criteria in {ws_file}",
                violations=[],
            )
    return CheckResult(passed=True, message="WS format checks passed", violations=[])
