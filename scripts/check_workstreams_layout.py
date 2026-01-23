"""Validate workstreams directory contains WS-only files."""

from __future__ import annotations

import sys
from pathlib import Path


ALLOWED_NAMES = {"INDEX.md", "TEMPLATE.md"}


def find_invalid_files(workstreams_dir: Path) -> list[Path]:
    invalid: list[Path] = []
    for path in workstreams_dir.rglob("*.md"):
        if path.name in ALLOWED_NAMES:
            continue
        if path.name.startswith("WS-"):
            continue
        invalid.append(path)
    return invalid


def main() -> None:
    repo_root = Path(__file__).resolve().parents[2]
    workstreams_dir = repo_root / "tools" / "hw_checker" / "docs" / "workstreams"
    if not workstreams_dir.exists():
        print("Workstreams directory not found.")
        sys.exit(0)

    invalid = find_invalid_files(workstreams_dir)
    if not invalid:
        print("✓ Workstreams layout valid (WS-only).")
        sys.exit(0)

    print("✗ Non-WS files found in workstreams directory:")
    for path in invalid:
        print(f" - {path}")
    print("Move non-WS files to docs/reviews or docs/archive.")
    sys.exit(1)


if __name__ == "__main__":
    main()
