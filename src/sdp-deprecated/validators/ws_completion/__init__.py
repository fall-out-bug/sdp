"""Workstream completion verification package.

This module provides backward compatibility by re-exporting all public APIs.
"""

from sdp.validators.ws_completion.models import CheckResult, VerificationResult
from sdp.validators.ws_completion.verifier import WSCompletionVerifier

__all__ = [
    "CheckResult",
    "VerificationResult",
    "WSCompletionVerifier",
]
