"""DEPRECATED: Use sdp.core.project_map submodule instead.

This module provides backward compatibility by re-exporting from the project_map package.
"""

import warnings

from sdp.core.project_map import parse_project_map

warnings.warn(
    "sdp.core.project_map_parser module is deprecated. "
    "Use 'from sdp.core.project_map import parse_project_map' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "parse_project_map",
]
