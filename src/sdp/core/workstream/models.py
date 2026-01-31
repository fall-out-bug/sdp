"""Workstream data models and enums.

DEPRECATED: This module is deprecated. Import from sdp.domain.workstream instead.
Kept for backward compatibility only.
"""

import warnings

warnings.warn(
    "sdp.core.workstream.models is deprecated. "
    "Use 'from sdp.domain.workstream import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

# Re-export from domain for backward compatibility
from sdp.domain.workstream import (
    AcceptanceCriterion,
    Workstream,
    WorkstreamID,
    WorkstreamSize,
    WorkstreamStatus,
)

__all__ = [
    "WorkstreamStatus",
    "WorkstreamSize",
    "WorkstreamID",
    "AcceptanceCriterion",
    "Workstream",
]
