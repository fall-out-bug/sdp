"""Base classes for SDP extension system.

Extensions provide project-specific customization:
- Hooks: validation scripts (pre-build, post-build)
- Patterns: domain-specific coding patterns
- Skills: custom commands (slash commands for Claude)
- Integrations: third-party services (GitHub, Telegram)
"""

from dataclasses import dataclass
from pathlib import Path
from typing import Protocol


@dataclass(frozen=True)
class ExtensionManifest:
    """Extension manifest metadata.

    Attributes:
        name: Unique extension identifier (e.g., 'hw_checker')
        version: Semantic version (e.g., '1.0.0')
        description: Human-readable description
        author: Extension author
        hooks_dir: Relative path to hooks directory
        patterns_dir: Relative path to patterns directory
        skills_dir: Relative path to skills directory
        integrations_dir: Relative path to integrations directory

    Example:
        >>> manifest = ExtensionManifest(
        ...     name="hw_checker",
        ...     version="1.0.0",
        ...     description="Clean Architecture validation",
        ...     author="SDP Team"
        ... )
    """

    name: str
    version: str
    description: str
    author: str
    hooks_dir: str = "hooks"
    patterns_dir: str = "patterns"
    skills_dir: str = "skills"
    integrations_dir: str = "integrations"


class Extension(Protocol):
    """Extension protocol.

    Extensions must implement this protocol to be loadable.
    Use structural subtyping (Protocol) for flexibility.

    Attributes:
        manifest: Extension metadata
        root_path: Absolute path to extension directory

    Example:
        >>> class MyExtension:
        ...     def __init__(self, root_path: Path):
        ...         self.manifest = ExtensionManifest(...)
        ...         self.root_path = root_path
        ...
        ...     def get_hooks_path(self) -> Path | None:
        ...         return self.root_path / self.manifest.hooks_dir
    """

    manifest: ExtensionManifest
    root_path: Path

    def get_hooks_path(self) -> Path | None:
        """Return absolute path to hooks directory.

        Returns:
            Hooks directory path if exists, None otherwise.

        Example:
            >>> ext.get_hooks_path()
            PosixPath('/home/user/.sdp/extensions/hw_checker/hooks')
        """
        ...

    def get_patterns_path(self) -> Path | None:
        """Return absolute path to patterns directory.

        Returns:
            Patterns directory path if exists, None otherwise.

        Example:
            >>> ext.get_patterns_path()
            PosixPath('/home/user/.sdp/extensions/hw_checker/patterns')
        """
        ...

    def get_skills_path(self) -> Path | None:
        """Return absolute path to skills directory.

        Returns:
            Skills directory path if exists, None otherwise.

        Example:
            >>> ext.get_skills_path()
            PosixPath('/home/user/.sdp/extensions/hw_checker/skills')
        """
        ...

    def get_integrations_path(self) -> Path | None:
        """Return absolute path to integrations directory.

        Returns:
            Integrations directory path if exists, None otherwise.

        Example:
            >>> ext.get_integrations_path()
            PosixPath('/home/user/.sdp/extensions/hw_checker/integrations')
        """
        ...


@dataclass
class BaseExtension:
    """Concrete implementation of Extension protocol.

    Attributes:
        manifest: Extension metadata
        root_path: Absolute path to extension directory

    Example:
        >>> manifest = ExtensionManifest(name="hw_checker", version="1.0.0", ...)
        >>> ext = BaseExtension(
        ...     manifest=manifest,
        ...     root_path=Path("/home/user/.sdp/extensions/hw_checker")
        ... )
        >>> hooks = ext.get_hooks_path()
    """

    manifest: ExtensionManifest
    root_path: Path

    def get_hooks_path(self) -> Path | None:
        """Return hooks directory if exists."""
        path = self.root_path / self.manifest.hooks_dir
        return path if path.is_dir() else None

    def get_patterns_path(self) -> Path | None:
        """Return patterns directory if exists."""
        path = self.root_path / self.manifest.patterns_dir
        return path if path.is_dir() else None

    def get_skills_path(self) -> Path | None:
        """Return skills directory if exists."""
        path = self.root_path / self.manifest.skills_dir
        return path if path.is_dir() else None

    def get_integrations_path(self) -> Path | None:
        """Return integrations directory if exists."""
        path = self.root_path / self.manifest.integrations_dir
        return path if path.is_dir() else None
