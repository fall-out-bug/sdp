"""Project map parsing package.

This module provides backward compatibility by re-exporting all public APIs.
"""

from sdp.core.project_map.parser import parse_project_map

__all__ = [
    "parse_project_map",
]
