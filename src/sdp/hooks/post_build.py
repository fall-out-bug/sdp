"""Post-build hook: quality checks after workstream execution."""

import os
import subprocess
import sys
from pathlib import Path

from sdp.hooks.common import find_project_root, find_workstream_dir


def _repo_root() -> Path:
    """Get repository root."""
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        capture_output=True,
        text=True,
        check=True,
    )
    return Path(result.stdout.strip())


def _project_root(repo_root: Path) -> tuple[Path, str]:
    """Return (work_dir, src_prefix). Project-agnostic: SDP or hw_checker."""
    hw = repo_root / "tools" / "hw_checker"
    if hw.exists():
        return hw, "src/hw_checker"
    return repo_root, "src/sdp"


def main() -> int:  # noqa: C901  # pragma: no cover
    """Run post-build checks. Usage: post_build.py WS-ID [module_path]"""
    if len(sys.argv) < 2:
        print("❌ Usage: post_build.py WS-ID [module_path]")
        return 1

    ws_id = sys.argv[1]
    module = sys.argv[2] if len(sys.argv) > 2 else ""
    repo_root = _repo_root()
    work_dir, src_prefix = _project_root(repo_root)

    print(f"🔍 Post-build checks for {ws_id}")
    print("================================")

    if (repo_root / "quality-gate.toml").exists():
        print("✓ Using quality-gate.toml")
    else:
        print("⚠️ quality-gate.toml not found, using defaults")

    lint_path = f"{src_prefix}/{module}" if module else src_prefix
    full_lint = work_dir / lint_path

    # Check 1: Regression tests
    print("\nCheck 1: Regression tests")
    tests_dir = work_dir / "tests" / "unit"
    if not tests_dir.exists():
        tests_dir = work_dir / "tests"
    if tests_dir.exists():
        proc = subprocess.run(
            [sys.executable, "-m", "pytest", str(tests_dir), "-q", "--tb=no"],
            cwd=work_dir,
            capture_output=True,
            text=True,
        )
        if proc.returncode == 0:
            print("✓ Regression tests passed")
        else:
            print("❌ Regression tests failed")
            print(proc.stdout + proc.stderr)
            return 1
    else:
        print("  Skipped (tests/ not found)")

    # Check 2: Ruff
    print("\nCheck 2: Linters")
    ruff_path = str(full_lint)
    proc = subprocess.run(
        [sys.executable, "-m", "ruff", "check", ruff_path, "--quiet"],
        cwd=work_dir,
        capture_output=True,
        text=True,
    )
    if proc.returncode == 0:
        print("✓ Ruff: no issues")
    else:
        print("❌ Ruff found issues")
        print(proc.stdout + proc.stderr)
        return 1

    # Mypy
    proc = subprocess.run(
        [sys.executable, "-m", "mypy", ruff_path, "--ignore-missing-imports"],
        cwd=work_dir,
        capture_output=True,
        text=True,
    )
    if proc.returncode == 0:
        print("✓ Mypy: no issues")
    else:
        print("❌ Mypy found type errors")
        print(proc.stdout + proc.stderr)
        return 1

    # Check 3: No deferral markers (TODO, FIXME, HACK, XXX)
    print("\nCheck 3: No deferral markers")
    todo_path = full_lint if full_lint.is_dir() else full_lint.parent
    proc = subprocess.run(
        ["grep", "-rn", "--exclude-dir=__pycache__", "--exclude=post_build.py",
         "TODO\\|FIXME\\|HACK\\|XXX", str(todo_path)],
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:  # grep returns 1 when no match
        print("✓ No deferral markers")
    else:
        print("❌ Found deferral markers (TODO/FIXME/HACK/XXX)")
        print(proc.stdout[:500])
        return 1

    # Check 4: File sizes
    print("\nCheck 4: File sizes (< 200 LOC)")
    large: list[str] = []
    for py_file in Path(todo_path).rglob("*.py"):
        if "__pycache__" in str(py_file):
            continue
        lines = len(py_file.read_text().splitlines())
        if lines > 200:
            large.append(f"{py_file} ({lines} lines)")
    if not large:
        print("✓ All files < 200 LOC")
    else:
        print("❌ Large files found")
        for f in large[:5]:
            print(f"  {f}")
        return 1

    # Check 8: Git commit with WS-ID
    print("\nCheck 8: Git commit")
    if os.environ.get("SKIP_COMMIT_CHECK") == "1":
        print("⚠️ Skipped (SKIP_COMMIT_CHECK=1)")
    else:
        proc = subprocess.run(
            ["git", "log", "-1", "--oneline"],
            cwd=repo_root,
            capture_output=True,
            text=True,
        )
        last = proc.stdout.strip() if proc.returncode == 0 else ""
        if ws_id in last:
            print(f"✓ Commit found: {last}")
        else:
            print(f"❌ No commit with WS-ID '{ws_id}' found")
            print(f"   Last commit: {last}")
            return 1

    # Check 9: Execution Report
    print("\nCheck 9: Execution Report")
    try:
        project_root = find_project_root(repo_root)
        ws_dir = find_workstream_dir(project_root)
        ws_file = next(ws_dir.rglob(f"{ws_id}*.md"), None)
    except RuntimeError:
        ws_dirs = [
            repo_root / "docs" / "workstreams",
            repo_root / "tools" / "hw_checker" / "docs" / "workstreams",
        ]
        ws_file = None
        for d in ws_dirs:
            if d.exists():
                ws_file = next(d.rglob(f"{ws_id}*.md"), None)
                if ws_file:
                    break
    if ws_file and ws_file.exists():
        content = ws_file.read_text()
        if "Execution Report" in content or "### Execution Report" in content:
            print(f"✓ Execution Report found in {ws_file.name}")
        else:
            print("❌ Execution Report NOT found in WS file")
            return 1
    else:
        print("⚠️ WS file not found")

    print("\n================================")
    print("✅ Post-build checks PASSED")
    return 0


if __name__ == "__main__":
    sys.exit(main())
