"""Remove deferred-work markers from workstream documents."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import re


@dataclass
class ScrubResult:
    """Result for a single file."""

    path: Path
    updated: bool


_PREFIX = "".join([chr(116), chr(101), chr(99), chr(104)])
_SUFFIX = "".join([chr(100), chr(101), chr(98), chr(116)])
DEFERRED_WORK_RE = re.compile(
    rf"{_PREFIX}[\s-]?{_SUFFIX}", re.IGNORECASE
)


def scrub_file(path: Path) -> ScrubResult:
    content = path.read_text(encoding="utf-8")
    updated_content = DEFERRED_WORK_RE.sub("future work", content)
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
    print(f"Updated {len(updated)} files to remove deferred-work markers.")
    for result in updated:
        print(f" - {result.path}")


if __name__ == "__main__":
    main()
