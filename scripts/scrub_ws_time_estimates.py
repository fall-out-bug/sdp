"""Remove time-based estimates from workstream documents."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import re


@dataclass
class ScrubResult:
    """Result for a single file."""

    path: Path
    updated: bool


TIME_RANGE_RE = re.compile(
    r"\b\d+(?:\.\d+)?\s*(?:-\s*\d+(?:\.\d+)?)?\s*hours?\b",
    re.IGNORECASE,
)
PER_HOUR_RE = re.compile(r"\bper hour\b", re.IGNORECASE)
LAST_HOUR_RE = re.compile(r"\blast hour\b", re.IGNORECASE)


def _scrub_text(text: str) -> str:
    text = LAST_HOUR_RE.sub("last window", text)
    text = PER_HOUR_RE.sub("per window", text)
    return TIME_RANGE_RE.sub("scope medium", text)


def scrub_file(path: Path) -> ScrubResult:
    content = path.read_text(encoding="utf-8")
    updated_content = _scrub_text(content)
    if updated_content != content:
        path.write_text(updated_content, encoding="utf-8")
        return ScrubResult(path=path, updated=True)
    return ScrubResult(path=path, updated=False)


def scrub_workspace(root: Path) -> list[ScrubResult]:
    workstreams_dir = root / "tools" / "hw_checker" / "docs" / "workstreams"
    results: list[ScrubResult] = []
    for path in sorted(workstreams_dir.rglob("*.md")):
        results.append(scrub_file(path))
    return results


def main() -> None:
    repo_root = Path(__file__).resolve().parents[2]
    results = scrub_workspace(repo_root)
    updated = [result for result in results if result.updated]
    print(f"Updated {len(updated)} files to remove time estimates.")
    for result in updated:
        print(f" - {result.path}")


if __name__ == "__main__":
    main()
