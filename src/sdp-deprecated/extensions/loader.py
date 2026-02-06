"""Extension discovery and loading.

Extensions are discovered in the following locations:
1. Project-local: {project_root}/sdp.local/
2. User-global: ~/.sdp/extensions/{name}/

Each extension directory must contain extension.yaml manifest.
"""

from pathlib import Path
from typing import Iterable

from sdp.extensions.base import BaseExtension, Extension
from sdp.extensions.manifest import ManifestParser, ValidationError


class ExtensionLoader:
    """Discover and load SDP extensions.

    Attributes:
        search_paths: List of directories to search for extensions
        manifest_parser: Parser for extension.yaml files

    Example:
        >>> loader = ExtensionLoader()
        >>> extensions = loader.discover_extensions()
        >>> for ext in extensions:
        ...     print(f"Loaded: {ext.manifest.name} v{ext.manifest.version}")
    """

    MANIFEST_FILENAME: str = "extension.yaml"

    def __init__(
        self,
        search_paths: list[Path] | None = None,
        manifest_parser: ManifestParser | None = None,
    ):
        """Initialize extension loader.

        Args:
            search_paths: Custom search paths (default: auto-detect)
            manifest_parser: Custom manifest parser (default: ManifestParser)

        Example:
            >>> loader = ExtensionLoader(search_paths=[Path("/custom/path")])
        """
        self.search_paths = search_paths or self._default_search_paths()
        self.manifest_parser = manifest_parser or ManifestParser()

    def discover_extensions(self) -> list[Extension]:
        """Discover all extensions in search paths.

        Returns:
            List of loaded extensions

        Example:
            >>> loader = ExtensionLoader()
            >>> extensions = loader.discover_extensions()
            >>> print(len(extensions))
            2
        """
        extensions: list[Extension] = []

        for search_path in self.search_paths:
            if not search_path.exists():
                continue

            for ext_dir in self._scan_extension_directories(search_path):
                try:
                    extension = self.load_extension(ext_dir)
                    extensions.append(extension)
                except ValidationError:
                    # Skip invalid extensions silently
                    # (allows gradual extension adoption)
                    continue

        return extensions

    def load_extension(self, ext_dir: Path) -> Extension:
        """Load extension from directory.

        Args:
            ext_dir: Extension root directory

        Returns:
            Loaded extension

        Raises:
            ValidationError: If extension is invalid

        Example:
            >>> loader = ExtensionLoader()
            >>> ext = loader.load_extension(Path("~/.sdp/extensions/hw_checker"))
        """
        manifest_path = ext_dir / self.MANIFEST_FILENAME
        manifest = self.manifest_parser.parse_file(manifest_path)

        return BaseExtension(manifest=manifest, root_path=ext_dir)

    def _default_search_paths(self) -> list[Path]:
        """Get default extension search paths.

        Returns:
            List of search paths (project-local, user-global)

        Example:
            >>> loader = ExtensionLoader()
            >>> paths = loader._default_search_paths()
            >>> print(paths)
            [PosixPath('/project/sdp.local'), PosixPath('/home/user/.sdp/extensions')]
        """
        paths: list[Path] = []

        # Project-local extensions
        project_root = self._find_project_root()
        if project_root:
            local_ext = project_root / "sdp.local"
            if local_ext.exists():
                paths.append(local_ext)

        # User-global extensions
        user_ext = Path.home() / ".sdp" / "extensions"
        if user_ext.exists():
            paths.append(user_ext)

        return paths

    def _find_project_root(self) -> Path | None:
        """Find project root (directory containing .git or sdp/).

        Returns:
            Project root path or None

        Example:
            >>> loader = ExtensionLoader()
            >>> root = loader._find_project_root()
            >>> print(root)
            PosixPath('/home/user/msu_ai_masters')
        """
        current = Path.cwd()

        for parent in [current, *current.parents]:
            if (parent / ".git").exists() or (parent / "sdp").exists():
                return parent

        return None

    def _scan_extension_directories(self, search_path: Path) -> Iterable[Path]:
        """Scan for extension directories in search path.

        Args:
            search_path: Directory to scan

        Yields:
            Extension directory paths

        Example:
            >>> loader = ExtensionLoader()
            >>> for ext_dir in loader._scan_extension_directories(Path("~/.sdp/extensions")):
            ...     print(ext_dir.name)
        """
        if not search_path.is_dir():
            return

        # If search_path itself contains manifest, treat as extension
        if (search_path / self.MANIFEST_FILENAME).exists():
            yield search_path
            return

        # Otherwise scan subdirectories
        for item in search_path.iterdir():
            if item.is_dir() and (item / self.MANIFEST_FILENAME).exists():
                yield item
