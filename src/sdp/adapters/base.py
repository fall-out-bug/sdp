"""Platform adapter base interface.

This module provides the abstract PlatformAdapter interface for multi-platform support
across different AI code assistants (Claude Code, Codex, OpenCode).

Example:
    Implement a concrete adapter:

    >>> from sdp.adapters.base import PlatformAdapter
    >>> class ClaudeCodeAdapter(PlatformAdapter):
    ...     def install(self, target_dir: Path) -> None:
    ...         # Create .claude directory structure
    ...         pass
    ...
    ...     def configure_hooks(self, hooks: list[str]) -> None:
    ...         # Update settings.json with hooks
    ...         pass
    ...
    ...     def load_skill(self, skill_name: str) -> dict[str, Any]:
    ...         # Read skill from .claude/skills/
    ...         pass
    ...
    ...     def get_settings(self) -> dict[str, Any]:
    ...         # Parse .claude/settings.json
    ...         pass

    Detect platform automatically:

    >>> from sdp.adapters.base import detect_platform, PlatformType
    >>> platform = detect_platform()
    >>> if platform == PlatformType.CLAUDE_CODE:
    ...     print("Running in Claude Code environment")
"""

from abc import ABC, abstractmethod
from enum import Enum
from pathlib import Path
from typing import Any


class PlatformType(Enum):
    """Supported AI coding platform types.

    Attributes:
        CLAUDE_CODE: Claude Code platform (.claude/)
        CODEX: Codex platform (.codex/)
        OPENCODE: OpenCode platform (.opencode/)
    """

    CLAUDE_CODE = "claude_code"
    CODEX = "codex"
    OPENCODE = "opencode"


class PlatformAdapter(ABC):
    """Abstract interface for platform-specific operations.

    All platform adapters must implement these methods to provide consistent
    behavior across different AI coding assistants.

    The adapter pattern allows SDP to work seamlessly across:
    - Claude Code (Anthropic)
    - Codex (OpenAI)
    - OpenCode (Open source alternative)

    Abstract Methods:
        install: Set up platform directory structure
        configure_hooks: Register git hooks or platform hooks
        load_skill: Read skill configuration from platform
        get_settings: Parse platform settings file
    """

    @abstractmethod
    def install(self, target_dir: Path) -> None:
        """Install platform adapter to target directory.

        Creates the platform-specific directory structure and configuration files.

        Args:
            target_dir: Directory where platform should be installed

        Raises:
            PermissionError: If target directory is not writable
            FileExistsError: If platform already installed
        """
        pass

    @abstractmethod
    def configure_hooks(self, hooks: list[str]) -> None:
        """Configure platform hooks.

        Registers git hooks or platform-specific hooks for automation.

        Args:
            hooks: List of hook names to configure (e.g., ["pre-commit", "post-build"])

        Raises:
            ValueError: If hook name is invalid
            FileNotFoundError: If hook script does not exist
        """
        pass

    @abstractmethod
    def load_skill(self, skill_name: str) -> dict[str, Any]:
        """Load skill configuration from platform.

        Reads and parses skill definition from platform-specific location.

        Args:
            skill_name: Name of skill to load (e.g., "idea", "build")

        Returns:
            Skill configuration as dictionary with at least:
                - name: Skill name
                - prompt: Skill prompt text
                - description: Skill description

        Raises:
            FileNotFoundError: If skill does not exist
            ValueError: If skill configuration is invalid
        """
        pass

    @abstractmethod
    def get_settings(self) -> dict[str, Any]:
        """Get platform settings.

        Reads and parses platform configuration file.

        Returns:
            Platform settings dictionary

        Raises:
            FileNotFoundError: If settings file does not exist
            ValueError: If settings format is invalid
        """
        pass


def detect_platform(search_path: Path | None = None) -> PlatformType | None:
    """Detect which AI coding platform is active.

    Searches for platform-specific directories and configuration files:
    - .claude/settings.json (Claude Code)
    - .codex/config.yaml or .codex/INSTALL.md (Codex)
    - .opencode/opencode.json (OpenCode)

    The search starts from `search_path` (or current directory) and walks up
    the directory tree until a platform is found or .git directory is reached.

    Priority order (if multiple platforms detected):
    1. Claude Code
    2. Codex
    3. OpenCode

    Args:
        search_path: Starting directory for search. Defaults to current directory.

    Returns:
        Detected platform type or None if no platform found

    Example:
        >>> platform = detect_platform()
        >>> if platform == PlatformType.CLAUDE_CODE:
        ...     print("Using Claude Code")
        >>> elif platform is None:
        ...     print("No platform detected")
    """
    start_path = search_path if search_path else Path.cwd()

    # Walk up directory tree until platform found or .git reached
    current = start_path.resolve()
    while current != current.parent:
        # Check for Claude Code
        if (current / ".claude" / "settings.json").exists():
            return PlatformType.CLAUDE_CODE

        # Check for Codex
        if (current / ".codex" / "config.yaml").exists() or (
            current / ".codex" / "INSTALL.md"
        ).exists():
            return PlatformType.CODEX

        # Check for OpenCode
        if (current / ".opencode" / "opencode.json").exists():
            return PlatformType.OPENCODE

        # Stop at git root (after checking current directory)
        if (current / ".git").exists():
            break

        current = current.parent

    return None
