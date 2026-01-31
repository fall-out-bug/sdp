"""Pre-commit hook: quality checks on staged files."""

import subprocess
import sys

from sdp.hooks.pre_commit_checks import (
    check_branch,
    check_python_bare_except,
    check_tech_debt,
    check_time_estimates,
    check_ws_format,
    repo_root,
    run_script,
    staged_files,
)


def main() -> int:  # noqa: C901  # pragma: no cover
    """Run pre-commit checks. Returns 0 on success, 1 on failure."""
    repo_root_path = repo_root()
    files = staged_files()

    print("🔍 Pre-commit checks")
    print("====================")

    if not files:
        print("No staged files, skipping checks")
        return 0

    checks = [
        ("Check 0: Branch check", check_branch()),
        ("Check 1: No time estimates", check_time_estimates(files)),
        ("Check 2: No tech debt markers", check_tech_debt(files)),
        ("Check 3: Python code quality", check_python_bare_except(files)),
        ("Check 5: WS file format", check_ws_format(files, repo_root_path)),
    ]

    for name, result in checks:
        print(f"\n{name}")
        print(result.format_terminal())
        if not result.passed:
            return 1

    # Check 1b: Workstreams layout
    print("\nCheck 1b: Workstreams layout")
    if any("workstreams/" in f for f in files):
        ok, out = run_script(repo_root_path, "scripts/check_workstreams_layout.py", [])
        print(out.strip() if out else ("✓" if ok else "❌"))
        if not ok:
            return 1
    else:
        print("  No workstreams files staged")

    # Check 3b: Quality gates
    print("\nCheck 3b: Quality Gates")
    src_py = [f for f in files if "src/" in f and f.endswith(".py")]
    if src_py:
        ok, out = run_script(repo_root_path, "scripts/check_quality_gates.py", ["--staged"])
        if out.strip():
            print(out.strip())
        if not ok:
            return 1
    else:
        print("  No src/ Python files staged")

    # Check 4: Architecture
    print("\nCheck 4: Clean Architecture")
    py_files = [f for f in files if f.endswith(".py")]
    if py_files:
        ok, out = run_script(repo_root_path, "scripts/check_architecture.py", ["--staged"])
        if out.strip():
            print(out.strip())
        if not ok:
            return 1
    else:
        print("  No Python files staged")

    # Check 6: Breaking changes (optional)
    print("\nCheck 6: Breaking changes")
    bc_script = repo_root_path / "tools" / "hw_checker" / "scripts" / "detect_breaking_changes.py"
    if bc_script.exists():
        proc = subprocess.run(
            [sys.executable, str(bc_script), "--staged"],
            cwd=repo_root_path / "tools" / "hw_checker",
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            print(proc.stdout + proc.stderr)
            return 1
        print("✓ No breaking changes")
    else:
        print("  Breaking changes detection script not found (skipping)")

    # Check 7: Test quality (optional)
    print("\nCheck 7: Test quality")
    test_files = [f for f in files if "test_" in f and f.endswith(".py")]
    if test_files and (repo_root_path / "tools" / "hw_checker").exists():
        abs_files = [
            str((repo_root_path / "tools" / "hw_checker" / f).resolve())
            for f in test_files
        ]
        proc = subprocess.run(
            [sys.executable, "scripts/check_test_quality.py", "--strict"] + abs_files,
            cwd=repo_root_path / "tools" / "hw_checker",
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            print(proc.stdout + proc.stderr)
            return 1
        print("✓ Test quality checks passed")
    else:
        msg = "  No test files staged" if not test_files else "  tools/hw_checker not found"
        print(msg)

    print("\n====================")
    print("✅ All pre-commit checks PASSED")
    return 0


if __name__ == "__main__":
    sys.exit(main())
