"""Workstream package - parsing and validation for SDP markdown files.

This module provides backward compatibility by re-exporting all public APIs.
Domain entities are now in sdp.domain package.
"""

from sdp.domain.workstream import (
    AcceptanceCriterion,
    Workstream,
    WorkstreamID,
    WorkstreamSize,
    WorkstreamStatus,
)
from sdp.core.workstream.parser import WorkstreamParseError, parse_workstream

__all__ = [
    "WorkstreamStatus",
    "WorkstreamSize",
    "WorkstreamID",
    "AcceptanceCriterion",
    "Workstream",
    "WorkstreamParseError",
    "parse_workstream",
]
