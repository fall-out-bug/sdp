"""Platform adapters for multi-IDE support.

This package provides abstraction over different IDE platforms:
- Claude Code (.claude/)
- Codex (.codex/)
- OpenCode (.opencode/)

Example:
    >>> from sdp.adapters.base import detect_platform, PlatformType
    >>> platform = detect_platform()
    >>> if platform == PlatformType.CLAUDE_CODE:
    ...     print("Running in Claude Code")
"""

from sdp.adapters.base import PlatformAdapter, PlatformType, detect_platform
from sdp.adapters.claude_code import ClaudeCodeAdapter
from sdp.adapters.codex import CodexAdapter
from sdp.adapters.opencode import OpenCodeAdapter

__all__ = [
    "PlatformAdapter",
    "PlatformType",
    "detect_platform",
    "ClaudeCodeAdapter",
    "CodexAdapter",
    "OpenCodeAdapter",
]
