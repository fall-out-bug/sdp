"""SDP error framework.

Provides structured error classes with categories, remediation steps, and context.

Usage:
    from sdp.errors import (
        SDPError,
        ErrorCategory,
        BeadsNotFoundError,
        CoverageTooLowError,
        # ... etc
    )

    raise SDPError(
        message="Something went wrong",
        category=ErrorCategory.RUNTIME,
        remediation="Try fixing X",
    )
"""

# Base classes
from .base import ErrorCategory, SDPError, format_error_for_json, format_error_for_terminal

# Build errors
from .build import ArtifactValidationError, BuildValidationError, TestFailureError

# Configuration errors
from .config import ConfigurationError, DependencyNotFoundError, HookExecutionError

# Quality-related errors
from .quality import (
    BeadsNotFoundError,
    CoverageTooLowError,
    QualityGateViolationError,
    WorkstreamValidationError,
)

__all__ = [
    # Base
    "ErrorCategory",
    "SDPError",
    "format_error_for_terminal",
    "format_error_for_json",
    # Quality
    "BeadsNotFoundError",
    "CoverageTooLowError",
    "QualityGateViolationError",
    "WorkstreamValidationError",
    # Config
    "ConfigurationError",
    "DependencyNotFoundError",
    "HookExecutionError",
    # Build
    "TestFailureError",
    "BuildValidationError",
    "ArtifactValidationError",
]
