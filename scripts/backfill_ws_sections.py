"""Backfill missing WS sections for legacy files."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import re


@dataclass
class BackfillResult:
    """Result for a single file."""

    path: Path
    updated: bool


GOAL_RE = re.compile(r"^###\s+🎯\s*(Goal|Цель)", re.MULTILINE)
AC_RE = re.compile(r"^###\s+Acceptance Criteria", re.MULTILINE)


GOAL_SECTION = """### 🎯 Goal

Ensure this legacy workstream is tracked with required metadata while
preserving the original content.
"""

AC_SECTION = """### Acceptance Criteria

- [ ] AC1: Frontmatter fields are present and valid
- [ ] AC2: Original content is preserved without loss
"""


def _has_frontmatter(content: str) -> bool:
    return bool(re.search(r"^---\n(.*?)\n---", content, re.DOTALL | re.MULTILINE))


def backfill_sections(path: Path) -> BackfillResult:
    content = path.read_text(encoding="utf-8")
    if not _has_frontmatter(content):
        return BackfillResult(path=path, updated=False)

    updated = False
    additions: list[str] = []

    if not GOAL_RE.search(content):
        additions.append(GOAL_SECTION)
        updated = True

    if not AC_RE.search(content):
        additions.append(AC_SECTION)
        updated = True

    if updated:
        content = content.rstrip() + "\n\n" + "\n".join(additions) + "\n"
        path.write_text(content, encoding="utf-8")

    return BackfillResult(path=path, updated=updated)


def backfill_workspace(root: Path) -> list[BackfillResult]:
    workstreams_dir = root / "tools" / "hw_checker" / "docs" / "workstreams"
    results: list[BackfillResult] = []
    for path in sorted(workstreams_dir.rglob("*.md")):
        results.append(backfill_sections(path))
    return results


def main() -> None:
    repo_root = Path(__file__).resolve().parents[2]
    results = backfill_workspace(repo_root)
    updated = [result for result in results if result.updated]
    print(f"Updated {len(updated)} files with missing sections.")
    for result in updated:
        print(f" - {result.path}")


if __name__ == "__main__":
    main()
