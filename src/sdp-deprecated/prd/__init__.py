"""PRD (Product Requirements Document) module for SDP.

This module provides tools for creating, validating, and maintaining
PRD documents with diagram generation.
"""

from .annotations import (
    Flow,
    FlowStep,
)
from .decorators import (
    get_flow_info,
    prd_flow,
    prd_step,
)
from .detector import detect_project_type
from .generator import (
    generate_diagrams,
    generate_flow_from_steps,
)
from .generator_mermaid import (
    generate_mermaid_component,
    generate_mermaid_deployment,
    generate_mermaid_sequence,
)
from .generator_plantuml import (
    generate_plantuml_component,
    generate_plantuml_deployment,
    generate_plantuml_sequence,
)
from .hash import (
    calculate_diagrams_hash,
    get_stored_hash,
    update_stored_hash,
    validate_diagrams_freshness,
)
from .parser import (
    get_frontmatter,
    parse_prd_sections,
    update_frontmatter,
)
from .parser_bash import (
    parse_bash_annotations,
    parse_directory_bash,
    parse_yaml_annotations,
)
from .parser_python import (
    parse_directory,
    parse_python_annotations,
    parse_python_annotations_ast,
)
from .profiles import (
    PROFILES,
    PRDProfile,
    PRDSection,
    ProjectType,
    get_profile,
    get_section_limit,
)
from .scaffold import (
    create_prd_file,
    generate_prd_template,
    update_prd_frontmatter,
)
from .validator import (
    Severity,
    ValidationIssue,
    format_validation_issues,
    has_critical_issues,
    validate_prd,
    validate_prd_file,
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
    # Parser
    "parse_prd_sections",
    "get_frontmatter",
    "update_frontmatter",
    # Validator
    "Severity",
    "ValidationIssue",
    "validate_prd",
    "validate_prd_file",
    "format_validation_issues",
    "has_critical_issues",
    # Generators
    "generate_diagrams",
    "generate_flow_from_steps",
    "generate_mermaid_sequence",
    "generate_mermaid_component",
    "generate_mermaid_deployment",
    "generate_plantuml_sequence",
    "generate_plantuml_component",
    "generate_plantuml_deployment",
    # Hash
    "calculate_diagrams_hash",
    "get_stored_hash",
    "update_stored_hash",
    "validate_diagrams_freshness",
]
