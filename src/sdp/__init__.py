"""SDP (Spec-Driven Protocol) - Workstream automation tools.

This package provides development tooling for:
- GitHub Issues integration (create, sync, track workstreams)
- Workstream file parsing and validation
- Feature decomposition and dependency management
- Project map parsing and querying
- Git workflow automation
"""

__version__ = "0.5.2"

# Public API exports
from sdp.core import (
    AcceptanceCriterion,
    CircularDependencyError,
    Constraint,
    Decision,
    Feature,
    MissingDependencyError,
    ProjectMap,
    ProjectMapParseError,
    TechStackItem,
    Workstream,
    WorkstreamParseError,
    WorkstreamSize,
    WorkstreamStatus,
    create_project_map_template,
    get_constraint,
    get_decision,
    load_feature_from_directory,
    load_feature_from_spec,
    parse_project_map,
    parse_workstream,
)

# Error framework
from sdp.errors import (
    ArtifactValidationError,
    BeadsNotFoundError,
    BuildValidationError,
    ConfigurationError,
    CoverageTooLowError,
    DependencyNotFoundError,
    ErrorCategory,
    HookExecutionError,
    QualityGateViolationError,
    SDPError,
    TestFailureError,
    WorkstreamValidationError,
    format_error_for_json,
    format_error_for_terminal,
)

__all__ = [
    "__version__",
    # Core types
    "AcceptanceCriterion",
    "Workstream",
    "WorkstreamStatus",
    "WorkstreamSize",
    "WorkstreamParseError",
    "Feature",
    "CircularDependencyError",
    "MissingDependencyError",
    "ProjectMap",
    "Decision",
    "Constraint",
    "TechStackItem",
    "ProjectMapParseError",
    # Core functions
    "parse_workstream",
    "load_feature_from_spec",
    "load_feature_from_directory",
    "parse_project_map",
    "get_decision",
    "get_constraint",
    "create_project_map_template",
    # Error framework
    "SDPError",
    "ErrorCategory",
    "BeadsNotFoundError",
    "CoverageTooLowError",
    "QualityGateViolationError",
    "WorkstreamValidationError",
    "ConfigurationError",
    "DependencyNotFoundError",
    "HookExecutionError",
    "TestFailureError",
    "BuildValidationError",
    "ArtifactValidationError",
    "format_error_for_terminal",
    "format_error_for_json",
]

