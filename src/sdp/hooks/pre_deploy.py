"""Pre-deploy hook: E2E tests before deployment."""

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


def main() -> int:  # noqa: C901
    """Run pre-deploy checks. Usage: pre_deploy.py F{XX} [staging|prod]"""
    if len(sys.argv) < 2:
        print("Usage: pre_deploy.py F{XX} [staging|prod]")
        return 1

    feature_id = sys.argv[1]
    environment = sys.argv[2] if len(sys.argv) > 2 else "staging"
    repo_root = _repo_root()

    print(f"🚀 Running pre-deploy checks for {feature_id} ({environment})...")

    # 1. Unit tests
    print("\n=== 1. Unit Tests ===")
    proc = subprocess.run(
        [
            sys.executable, "-m", "pytest", "tests/unit/",
            "-v", "--tb=short", "--cov=src/sdp", "--cov-report=term-missing",
            "--cov-fail-under=70",
        ],
        cwd=repo_root,
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        print("❌ Unit tests failed")
        print(proc.stdout + proc.stderr)
        return 1
    print("✅ Unit tests passed")

    # 2. Integration tests
    print("\n=== 2. Integration Tests ===")
    int_dir = repo_root / "tests" / "integration"
    if int_dir.exists():
        proc = subprocess.run(
            [sys.executable, "-m", "pytest", "tests/integration/", "-v", "--tb=short"],
            cwd=repo_root,
            capture_output=True,
            text=True,
        )
        if proc.returncode != 0:
            print("❌ Integration tests failed")
            print(proc.stdout + proc.stderr)
            return 1
        print("✅ Integration tests passed")
    else:
        print("  Skipped (tests/integration not found)")

    # 3. Type checking
    print("\n=== 3. Type Checking (mypy) ===")
    proc = subprocess.run(
        [sys.executable, "-m", "mypy", "src/sdp/", "--strict"],
        cwd=repo_root,
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        print("⚠️ Type checking issues found")
        print(proc.stdout[:1000])
    else:
        print("✅ Type checking passed")

    # 4. Linting
    print("\n=== 4. Linting (ruff) ===")
    proc = subprocess.run(
        [sys.executable, "-m", "ruff", "check", "src/sdp/"],
        cwd=repo_root,
        capture_output=True,
        text=True,
    )
    if proc.returncode != 0:
        print("⚠️ Linting issues found")
        print(proc.stdout[:1000])
    else:
        print("✅ Linting passed")

    # 5. Production checks
    if environment == "prod":
        print("\n=== 5. Production Readiness ===")
        for py_file in (repo_root / "src").rglob("*.py"):
            if "__pycache__" in str(py_file):
                continue
            content = py_file.read_text()
            if "DEBUG" in content and "= True" in content and "= False" not in content:
                print(f"❌ DEBUG flags found in {py_file}")
                return 1
        print("✅ Production readiness checks passed")

    print("\n✅ All pre-deploy checks passed for " + feature_id)
    print(f"✅ Ready to deploy to {environment}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
