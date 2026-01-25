"""Diagram hash calculation for PRD freshness validation.

This module calculates SHA256 hashes of PRD annotations to detect
when diagrams need to be regenerated.
"""

import hashlib
from pathlib import Path

from .parser_bash import parse_directory_bash


def calculate_diagrams_hash(project_path: Path) -> str:
    """Calculate SHA256 hash of all @prd annotations in project.

    This function scans all Python and bash/yaml files in the project
    for @prd annotations and creates a hash of the normalized content.

    Args:
        project_path: Path to the project root

    Returns:
        SHA256 hash (first 12 characters)
    """
    all_steps = []

    # Collect Python annotations
    for py_file in project_path.rglob("*.py"):
        # Skip common non-source directories
        if any(skip in str(py_file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules", ".git", "build", "dist"]):
            continue

        from .parser_python import parse_python_annotations
        all_steps.extend(parse_python_annotations(py_file))

    # Collect bash/yaml annotations
    for ext in ["*.sh", "*.bash", "*.yml", "*.yaml"]:
        for file in project_path.rglob(ext):
            # Skip common non-source directories
            if any(skip in str(file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules", ".git", "build", "dist"]):
                continue

            all_steps.extend(parse_directory_bash(file.parent) if file.parent == project_path else [])

    # Sort for deterministic hash
    sorted_steps = sorted(all_steps, key=lambda s: (s.flow_name, s.step_number, str(s.source_file), s.line_number))

    # Create hash from normalized content
    content = "\n".join(
        f"{s.flow_name}|{s.step_number}|{s.description}|{s.source_file}|{s.line_number}"
        for s in sorted_steps
    )

    return hashlib.sha256(content.encode()).hexdigest()[:12]


def get_stored_hash(project_map_path: Path) -> str | None:
    """Extract diagrams_hash from PROJECT_MAP.md frontmatter.

    Args:
        project_map_path: Path to PROJECT_MAP.md

    Returns:
        Hash string or None if not found
    """
    import re

    if not project_map_path.exists():
        return None

    try:
        content = project_map_path.read_text()
        match = re.search(r'^diagrams_hash:\s*(\w+)', content, re.MULTILINE)
        return match.group(1) if match else None
    except Exception:
        return None


def update_stored_hash(project_map_path: Path, new_hash: str) -> None:
    """Update diagrams_hash in PROJECT_MAP.md frontmatter.

    Args:
        project_map_path: Path to PROJECT_MAP.md
        new_hash: New hash value to store
    """
    import re

    if not project_map_path.exists():
        return

    try:
        content = project_map_path.read_text()

        if "diagrams_hash:" in content:
            # Update existing hash
            content = re.sub(
                r'^diagrams_hash:\s*\w*',
                f'diagrams_hash: {new_hash}',
                content,
                flags=re.MULTILINE
            )
        else:
            # Add hash after prd_version line
            content = re.sub(
                r'^(prd_version:\s*".+")',
                f'\\1\ndiagrams_hash: {new_hash}',
                content,
                flags=re.MULTILINE
            )

        project_map_path.write_text(content)
    except Exception:
        pass


def validate_diagrams_freshness(project_path: Path, project_map_name: str = "PROJECT_MAP.md") -> tuple[bool, str | None]:
    """Check if diagrams are up-to-date.

    Args:
        project_path: Path to the project root
        project_map_name: Name of the PROJECT_MAP file

    Returns:
        Tuple of (is_fresh, stored_hash) where is_fresh is True if hashes match
    """
    docs_path = project_path / "docs" / project_map_name
    if not docs_path.exists():
        # No PROJECT_MAP.md, consider fresh (nothing to check)
        return True, None

    stored_hash = get_stored_hash(docs_path)
    if stored_hash is None:
        # No hash stored, can't validate
        return True, None

    current_hash = calculate_diagrams_hash(project_path)
    return current_hash == stored_hash, stored_hash
