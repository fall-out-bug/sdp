"""DEPRECATED: Use sdp.validators.supersede submodule instead.

This module provides backward compatibility by re-exporting from the supersede package.
"""

import warnings

from sdp.validators.supersede import (
    SupersedeChain,
    SupersedeResult,
    SupersedeValidator,
    ValidationReport,
)

warnings.warn(
    "sdp.validators.supersede_checker module is deprecated. "
    "Use 'from sdp.validators.supersede import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "SupersedeChain",
    "SupersedeResult",
    "ValidationReport",
    "SupersedeValidator",
]
