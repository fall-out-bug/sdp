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
]
