"""DEPRECATED: Use sdp.validators.ws_completion submodule instead.

This module provides backward compatibility by re-exporting from the ws_completion package.
"""

import subprocess  # noqa: F401  # Re-export for tests
import warnings

from sdp.validators.ws_completion import CheckResult, VerificationResult, WSCompletionVerifier

warnings.warn(
    "sdp.validators.ws_completion module is deprecated. "
    "Use 'from sdp.validators.ws_completion import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "CheckResult",
    "VerificationResult",
    "WSCompletionVerifier",
    "subprocess",
]
