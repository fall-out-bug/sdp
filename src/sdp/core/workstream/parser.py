"""Workstream parsing functionality."""

from pathlib import Path
from typing import Optional

from sdp.core.workstream.markdown_helpers import (
    extract_acceptance_criteria,
    extract_code_blocks,
    extract_dependencies,
    extract_frontmatter,
    extract_section,
    extract_steps,
    extract_title,
    strip_frontmatter,
)
from sdp.core.workstream.models import Workstream, WorkstreamID, WorkstreamSize, WorkstreamStatus
from sdp.errors import ErrorCategory, SDPError


class WorkstreamParseError(SDPError):
    """Workstream parsing error with actionable guidance."""

    def __init__(
        self,
        message: str,
        file_path: Optional[Path] = None,
        parse_error: Optional[str] = None,
    ) -> None:
        super().__init__(
            category=ErrorCategory.VALIDATION,
            message=message,
            remediation=(
                "1. Check WS ID format: PP-FFF-SS (e.g., 00-500-01)\n"
                "   - PP: Project ID (00-99)\n"
                "   - FFF: Feature ID (001-999)\n"
                "   - SS: Sequence number (01-99)\n"
                "2. Ensure file starts with --- frontmatter\n"
                "3. Validate YAML syntax\n"
                "4. See docs/workstreams/TEMPLATE.md for template"
            ),
            docs_url="https://sdp.dev/docs/workstreams#format",
            context={
                "file_path": str(file_path) if file_path else None,
                "parse_error": parse_error,
            }
            if file_path or parse_error
            else None,
        )


def parse_workstream(file_path: Path) -> Workstream:
    """Parse workstream markdown file.

    Args:
        file_path: Path to WS markdown file

    Returns:
        Parsed Workstream instance

    Raises:
        WorkstreamParseError: If file has no frontmatter or required fields missing
    """
    content = file_path.read_text(encoding="utf-8")

    # Extract and validate frontmatter
    try:
        frontmatter = extract_frontmatter(content, str(file_path))
    except ValueError as e:
        raise WorkstreamParseError(
            message=str(e), file_path=file_path, parse_error=str(e)
        ) from e

    # Parse and validate ws_id
    ws_id_raw: str = str(frontmatter["ws_id"])
    try:
        parsed_ws_id = WorkstreamID.parse(ws_id_raw)
        ws_id: str = str(parsed_ws_id)  # Use normalized format
    except ValueError as e:
        raise WorkstreamParseError(
            message=f"Invalid ws_id format: {e}",
            file_path=file_path,
            parse_error=str(e),
        ) from e

    # Parse enum fields
    feature: str = str(frontmatter["feature"])
    status_str: str = str(frontmatter["status"])
    size_str: str = str(frontmatter["size"])

    try:
        status = WorkstreamStatus(status_str)
    except ValueError as e:
        raise WorkstreamParseError(
            message=f"Invalid status: {status_str}",
            file_path=file_path,
            parse_error=str(e),
        ) from e

    try:
        size = WorkstreamSize(size_str)
    except ValueError as e:
        raise WorkstreamParseError(
            message=f"Invalid size: {size_str}",
            file_path=file_path,
            parse_error=str(e),
        ) from e

    # Optional fields
    github_issue_val = frontmatter.get("github_issue")
    github_issue: Optional[int] = int(github_issue_val) if github_issue_val is not None else None

    assignee_val = frontmatter.get("assignee")
    assignee: Optional[str] = str(assignee_val) if assignee_val is not None else None

    # Parse markdown body
    body = strip_frontmatter(content)
    title = extract_title(body)
    goal = extract_section(body, "Goal")
    context = extract_section(body, "Context")
    criteria = extract_acceptance_criteria(body)
    deps = extract_dependencies(body)

    # Merge with frontmatter depends_on (PP-FFF-SS or ws_id format)
    fm_deps = frontmatter.get("depends_on")
    if fm_deps:
        for d in fm_deps if isinstance(fm_deps, list) else [fm_deps]:
            dep_id = str(d).strip()
            if dep_id and dep_id not in deps:
                deps.append(dep_id)

    steps = extract_steps(body)
    code_blocks = extract_code_blocks(body)

    return Workstream(
        ws_id=ws_id,
        feature=feature,
        status=status,
        size=size,
        github_issue=github_issue,
        assignee=assignee,
        title=title,
        goal=goal,
        acceptance_criteria=criteria,
        context=context,
        dependencies=deps,
        steps=steps,
        code_blocks=code_blocks,
        file_path=file_path,
    )
