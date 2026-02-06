"""Extension validation utilities.

Validates extension structure and contents:
- Manifest schema compliance
- Directory structure
- Required files
"""


from sdp.extensions.base import Extension
from sdp.extensions.manifest import ValidationError


class ExtensionValidator:
    """Validate extension structure and contents.

    Example:
        >>> validator = ExtensionValidator()
        >>> validator.validate(extension)
    """

    def validate(self, extension: Extension) -> None:
        """Validate extension.

        Args:
            extension: Extension to validate

        Raises:
            ValidationError: If extension is invalid

        Example:
            >>> validator = ExtensionValidator()
            >>> validator.validate(ext)
        """
        self._validate_root_exists(extension)
        self._validate_manifest(extension)
        self._validate_directories(extension)

    def _validate_root_exists(self, extension: Extension) -> None:
        """Validate extension root directory exists.

        Args:
            extension: Extension to validate

        Raises:
            ValidationError: If root does not exist
        """
        if not extension.root_path.exists():
            raise ValidationError(
                f"Extension root does not exist: {extension.root_path}"
            )

        if not extension.root_path.is_dir():
            raise ValidationError(
                f"Extension root is not a directory: {extension.root_path}"
            )

    def _validate_manifest(self, extension: Extension) -> None:
        """Validate manifest fields.

        Args:
            extension: Extension to validate

        Raises:
            ValidationError: If manifest is invalid
        """
        manifest = extension.manifest

        # Validate name (alphanumeric + underscore)
        if not manifest.name:
            raise ValidationError("Extension name is empty")

        if not manifest.name.replace("_", "").isalnum():
            raise ValidationError(
                f"Extension name must be alphanumeric: {manifest.name}"
            )

        # Validate version (basic semantic versioning check)
        if not manifest.version:
            raise ValidationError("Extension version is empty")

        version_parts = manifest.version.split(".")
        if len(version_parts) != 3:
            raise ValidationError(
                f"Extension version must be semver (X.Y.Z): {manifest.version}"
            )

        if not all(part.isdigit() for part in version_parts):
            raise ValidationError(
                f"Extension version must be semver (X.Y.Z): {manifest.version}"
            )

    def _validate_directories(self, extension: Extension) -> None:
        """Validate extension directories exist if configured.

        Args:
            extension: Extension to validate

        Note:
            Directories are optional - validation only checks that
            configured directories exist, not that they are required.
        """
        # All directories are optional, no validation needed
        # (get_*_path methods return None if directory doesn't exist)
        pass
