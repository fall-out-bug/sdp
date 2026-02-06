"""Pre-push hook: regression tests before pushing."""

import os
import subprocess
import sys
from pathlib import Path


def _repo_root() -> Path:
    """Get repository root."""
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        capture_output=True,
        text=True,
        check=True,
    )
    return Path(result.stdout.strip())


def _files_to_push() -> list[str]:
    """Get list of files to be pushed."""
    result = subprocess.run(
        ["git", "diff", "--name-only", "HEAD", "@{u}"],
        capture_output=True,
        text=True,
        cwd=_repo_root(),
    )
    if result.returncode != 0 or not result.stdout.strip():
        result = subprocess.run(
            ["git", "diff", "--name-only", "HEAD~1", "HEAD"],
            capture_output=True,
            text=True,
            cwd=_repo_root(),
        )
    return [f for f in result.stdout.strip().split("\n") if f] if result.returncode == 0 else []


def main() -> int:  # noqa: C901  # pragma: no cover
    """Run pre-push checks. Returns 0 on success, 1 on failure."""
    hard_push = os.environ.get("SDP_HARD_PUSH", "0") == "1"
    repo_root = _repo_root()

    print("🔍 Running pre-push checks...")
    print("")
    if hard_push:
        print("🔒 HARD PUSH MODE enabled (SDP_HARD_PUSH=1)")
        print("   Test failures will BLOCK push")
    else:
        print("⚠️  Soft mode (set SDP_HARD_PUSH=1 to enable hard blocking)")
        print("   Test failures will WARN but not block push")
    print("")

    files = _files_to_push()
    py_files = [f for f in files if f.endswith(".py")]

    if not py_files:
        print("No Python files to push, skipping tests")
        print("")
        print("✅ Pre-push checks complete")
        return 0

    has_failures = False

    # Run regression tests
    print("1. Running regression tests...")
    tests_dir = repo_root / "tests"
    if tests_dir.exists():
        proc = subprocess.run(
            [sys.executable, "-m", "pytest", "tests/", "-x", "-q", "--tb=no"],
            cwd=repo_root,
            capture_output=True,
            text=True,
        )
        if proc.returncode == 0:
            print("✓ Regression tests passed")
        else:
            print("❌ Regression tests failed")
            print(proc.stdout + proc.stderr)
            if hard_push:
                return 1
            has_failures = True
    else:
        print("⚠️  tests/ directory not found, skipping tests")

    # Check coverage
    cov_file = repo_root / ".coverage"
    if cov_file.exists():
        print("")
        print("2. Checking coverage...")
        proc = subprocess.run(
            [sys.executable, "-m", "coverage", "report", "--show-missing"],
            cwd=repo_root,
            capture_output=True,
            text=True,
        )
        if proc.returncode == 0 and proc.stdout:
            import re as re_mod
            match = re_mod.search(r"(\d+)%", proc.stdout)
            if match:
                pct = int(match.group(1))
                if pct < 80:
                    print(f"❌ Coverage is below 80% (currently: {pct}%)")
                    if hard_push:
                        return 1
                    has_failures = True
                else:
                    print(f"✓ Coverage is {pct}% (≥ 80%)")
        else:
            print("  Coverage report not available")

    print("")
    if hard_push and has_failures:
        print("🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)")
        print("To bypass: git push --no-verify")
        return 1
    if has_failures:
        print("⚠️  Pre-push checks complete (WARNING mode - failures detected)")
    else:
        print("✅ Pre-push checks complete")
    return 0


if __name__ == "__main__":
    sys.exit(main())
