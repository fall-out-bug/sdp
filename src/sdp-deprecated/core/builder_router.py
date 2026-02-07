"""DEPRECATED: Use sdp.core.builder submodule instead.

This module provides backward compatibility by re-exporting from the builder package.
"""

import warnings

from sdp.core.builder import (
    DEFAULT_WEIGHTS,
    BuildResult,
    BuilderRouter,
    HumanEscalationError,
    RetryPolicy,
    select_model_for_tier,
    select_model_weighted,
)

warnings.warn(
    "sdp.core.builder_router module is deprecated. "
    "Use 'from sdp.core.builder import ...' instead.",
    DeprecationWarning,
    stacklevel=2,
)

__all__ = [
    "DEFAULT_WEIGHTS",
    "select_model_weighted",
    "select_model_for_tier",
    "RetryPolicy",
    "BuildResult",
    "HumanEscalationError",
    "BuilderRouter",
]
