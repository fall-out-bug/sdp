"""Bash/YAML annotation parser.

This module extracts PRD flow annotations from bash and YAML files.
"""

import re
from pathlib import Path

from .annotations import FlowStep


def parse_bash_annotations(path: Path) -> list[FlowStep]:
    """Parse # @prd: comments from bash/yaml files.

    This function looks for comments in the format:
    # @prd: flow=name, step=N, desc=description

    Args:
        path: Path to the bash/yaml file

    Returns:
        List of FlowStep objects found in the file
    """
    try:
        content = path.read_text()
    except Exception:
        return []

    steps = []

    # Pattern: # @prd: flow=name, step=N, desc=description
    # Also supports variations in spacing and quoting
    pattern = re.compile(
        r'^#\s*@prd:\s*flow=([^,\s]+),\s*step=(\d+)(?:,\s*desc=(.+))?$',
        re.MULTILINE
    )

    for match in pattern.finditer(content):
        flow_name = match.group(1).strip()
        step_num = int(match.group(2))
        description = match.group(3).strip() if match.group(3) else ""

        # Clean up description (remove quotes if present)
        if description.startswith('"') and description.endswith('"'):
            description = description[1:-1]
        elif description.startswith("'") and description.endswith("'"):
            description = description[1:-1]

        # Calculate line number
        line_number = content[:match.start()].count('\n') + 1

        steps.append(FlowStep(
            flow_name=flow_name,
            step_number=step_num,
            description=description,
            source_file=path,
            line_number=line_number
        ))

    return steps


def parse_directory_bash(directory: Path) -> list[FlowStep]:
    """Parse all bash/yaml files in directory.

    Args:
        directory: Root directory to search

    Returns:
        List of FlowStep objects from all matching files
    """
    all_steps = []
    extensions = ["*.sh", "*.bash", "*.yml", "*.yaml"]

    try:
        for ext in extensions:
            for file in directory.rglob(ext):
                # Skip common non-source directories
                if any(skip in str(file) for skip in ["venv", ".venv", "__pycache__", ".tox", "node_modules", ".git"]):  # noqa: E501
                    continue

                if file.is_file():
                    all_steps.extend(parse_bash_annotations(file))
    except Exception:
        pass

    return all_steps


def parse_yaml_annotations(path: Path) -> list[FlowStep]:
    """Parse @prd comments from YAML files.

    This is an alias for parse_bash_annotations since the format is the same.

    Args:
        path: Path to the YAML file

    Returns:
        List of FlowStep objects found in the file
    """
    return parse_bash_annotations(path)
