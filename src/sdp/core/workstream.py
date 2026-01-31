"""DEPRECATED: Use sdp.core.workstream submodule instead.

This module provides backward compatibility by re-exporting from the workstream package.
"""

import warnings

from sdp.core.workstream import (
    AcceptanceCriterion,
    Workstream,
    WorkstreamID,
    WorkstreamParseError,
    WorkstreamSize,
    WorkstreamStatus,
    parse_workstream,
)

warnings.warn(
    "sdp.core.workstream module is deprecated. "
    "Use 'from sdp.core.workstream import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "WorkstreamStatus",
    "WorkstreamSize",
    "WorkstreamID",
    "AcceptanceCriterion",
    "Workstream",
    "WorkstreamParseError",
    "parse_workstream",
]
