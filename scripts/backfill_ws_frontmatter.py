"""Backfill missing WS frontmatter fields for legacy files."""

from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
import re


@dataclass
class BackfillResult:
    """Result for a single file."""

    path: Path
    updated: bool


FRONTMATTER_RE = re.compile(r"^---\n(.*?)\n---", re.MULTILINE | re.DOTALL)
WS_ID_RE = re.compile(r"\bWS-\d{2,3}(?:\.\d+)?(?:-\d{2})?(?:-[\w-]+)?\b")
FEATURE_RE = re.compile(r"\bF\d{2,3}\b")
SIZE_RE = re.compile(r"\b(SMALL|MEDIUM|LARGE)\b", re.IGNORECASE)


def _detect_feature(text: str, ws_id: str) -> str:
    match = FEATURE_RE.search(text)
    if match:
        return match.group(0)
    ws_digits = re.search(r"WS-(\d{2,3})", ws_id)
    if ws_digits:
        return f"F{ws_digits.group(1)}"
    return "F00"


def _detect_size(text: str) -> str:
    match = SIZE_RE.search(text)
    if match:
        return match.group(1).upper()
    return "MEDIUM"


def _detect_status(path: Path) -> str:
    if "completed" in path.parts:
        return "completed"
    if "active" in path.parts:
        return "active"
    return "backlog"


def _parse_frontmatter(content: str) -> dict[str, str]:
    match = FRONTMATTER_RE.search(content)
    if not match:
        return {}
    lines = match.group(1).splitlines()
    data: dict[str, str] = {}
    for line in lines:
        if ":" in line:
            key, value = line.split(":", 1)
            data[key.strip()] = value.strip()
    return data


def _render_frontmatter(data: dict[str, str]) -> str:
    lines = ["---"]
    for key, value in data.items():
        lines.append(f"{key}: {value}")
    lines.append("---")
    return "\n".join(lines) + "\n"


def _ensure_frontmatter(path: Path) -> BackfillResult:
    content = path.read_text(encoding="utf-8")
    ws_id_match = WS_ID_RE.search(content)
    if not ws_id_match:
        return BackfillResult(path=path, updated=False)
    ws_id = ws_id_match.group(0)
    feature = _detect_feature(content, ws_id)
    size = _detect_size(content)
    status = _detect_status(path)

    frontmatter = _parse_frontmatter(content)
    updated = False

    if not frontmatter:
        frontmatter = {
            "ws_id": ws_id,
            "feature": feature,
            "status": status,
            "size": size,
            "github_issue": "null",
            "assignee": "null",
            "started": "null",
            "completed": "null",
            "blocked_reason": "null",
        }
        content = _render_frontmatter(frontmatter) + content
        updated = True
    else:
        if "ws_id" not in frontmatter:
            frontmatter["ws_id"] = ws_id
            updated = True
        if "feature" not in frontmatter:
            frontmatter["feature"] = feature
            updated = True
        if "status" not in frontmatter:
            frontmatter["status"] = status
            updated = True
        if "size" not in frontmatter:
            frontmatter["size"] = size
            updated = True
        if updated:
            rendered = _render_frontmatter(frontmatter)
            content = FRONTMATTER_RE.sub(rendered.strip(), content, count=1)

    if updated:
        path.write_text(content, encoding="utf-8")
    return BackfillResult(path=path, updated=updated)


def backfill_workspace(root: Path) -> list[BackfillResult]:
    workstreams_dir = root / "tools" / "hw_checker" / "docs" / "workstreams"
    results: list[BackfillResult] = []
    for path in sorted(workstreams_dir.rglob("WS-*.md")):
        results.append(_ensure_frontmatter(path))
    return results


def main() -> None:
    repo_root = Path(__file__).resolve().parents[2]
    results = backfill_workspace(repo_root)
    updated = [result for result in results if result.updated]
    print(f"Updated {len(updated)} files with frontmatter.")
    for result in updated:
        print(f" - {result.path}")


if __name__ == "__main__":
    main()
