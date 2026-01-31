"""Hooks package - Git and workflow hooks."""

__all__ = [
    "PostWSCompleteHook",
    "CheckResult",
]

from sdp.hooks.common import CheckResult
from sdp.hooks.ws_complete import PostWSCompleteHook
