"""SDP (Spec-Driven Protocol) - Workstream automation tools.

This package provides development tooling for:
- GitHub Issues integration (create, sync, track workstreams)
- Workstream file parsing and validation
- Feature decomposition and dependency management
- Project map parsing and querying
- Git workflow automation
"""

__version__ = "0.3.0"

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
]

