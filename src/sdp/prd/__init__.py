"""PRD (Product Requirements Document) module for SDP.

This module provides tools for creating, validating, and maintaining
PRD documents with diagram generation.
"""

from .profiles import (
    ProjectType,
    PRDSection,
    PRDProfile,
    PROFILES,
    get_profile,
    get_section_limit,
)

from .detector import detect_project_type

from .scaffold import (
    generate_prd_template,
    create_prd_file,
    update_prd_frontmatter,
)

from .annotations import (
    FlowStep,
    Flow,
)

from .decorators import (
    prd_flow,
    prd_step,
    get_flow_info,
)

from .parser_python import (
    parse_python_annotations,
    parse_python_annotations_ast,
    parse_directory,
)

from .parser_bash import (
    parse_bash_annotations,
    parse_directory_bash,
    parse_yaml_annotations,
)

__all__ = [
    # Profiles
    "ProjectType",
    "PRDSection",
    "PRDProfile",
    "PROFILES",
    "get_profile",
    "get_section_limit",
    # Detector
    "detect_project_type",
    # Scaffold
    "generate_prd_template",
    "create_prd_file",
    "update_prd_frontmatter",
    # Annotations
    "FlowStep",
    "Flow",
    # Decorators
    "prd_flow",
    "prd_step",
    "get_flow_info",
    # Parsers
    "parse_python_annotations",
    "parse_python_annotations_ast",
    "parse_directory",
    "parse_bash_annotations",
    "parse_directory_bash",
    "parse_yaml_annotations",
]
