"""Builder routing package.

This module provides backward compatibility by re-exporting all public APIs.
"""

from sdp.core.builder.model_selector import (
    DEFAULT_WEIGHTS,
    select_model_for_tier,
    select_model_weighted,
)
from sdp.core.builder.policies import BuildResult, HumanEscalationError, RetryPolicy
from sdp.core.builder.router import BuilderRouter

__all__ = [
    "DEFAULT_WEIGHTS",
    "select_model_weighted",
    "select_model_for_tier",
    "RetryPolicy",
    "BuildResult",
    "HumanEscalationError",
    "BuilderRouter",
]
