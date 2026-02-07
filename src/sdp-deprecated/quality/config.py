"""Quality gate configuration parser for TOML files."""

from pathlib import Path
from typing import Any, cast

try:
    import tomllib  # Python 3.11+
except ImportError:
    import tomli as tomllib  # Python 3.10 backport

from sdp.quality.config_parsers import ConfigSectionParsers
from sdp.quality.exceptions import ConfigValidationError
from sdp.quality.models import QualityGateConfig as QualityGateConfigModel


class QualityGateConfigLoader:
    """Parser and validator for quality-gate.toml files."""

    DEFAULT_PATHS = [
        "quality-gate.toml",
        ".quality-gate.toml",
        "config/quality-gate.toml",
    ]

    def __init__(self, config_path: str | Path | None = None) -> None:
        """Initialize quality gate configuration.

        Args:
            config_path: Path to quality-gate.toml file. If None, searches default locations.
        """
        self._config_path = self._resolve_path(config_path)
        self._raw_config = self._load_config()
        self._parsers = ConfigSectionParsers(self._raw_config)
        self._config = self._parse_config()

    @property
    def config(self) -> QualityGateConfigModel:
        """Get parsed configuration object."""
        return self._config

    @property
    def raw(self) -> dict[str, Any]:
        """Get raw TOML dictionary."""
        return self._raw_config

    def _resolve_path(self, config_path: str | Path | None) -> Path | None:
        """Resolve configuration file path."""
        if config_path:
            path = Path(config_path)
            if path.exists():
                return path
            raise ConfigValidationError([f"Config file not found: {config_path}"])

        # Search default locations
        for default_path in self.DEFAULT_PATHS:
            path = Path(default_path)
            if path.exists():
                return path

        # Return None (will use defaults)
        return None

    def _load_config(self) -> dict[str, Any]:
        """Load TOML configuration from file or get defaults."""
        if self._config_path is None:
            # Return default configuration
            return self._get_default_config()

        try:
            with open(self._config_path, "rb") as f:
                return cast(dict[str, Any], tomllib.load(f))
        except tomllib.TOMLDecodeError as e:
            raise ConfigValidationError([f"Invalid TOML in {self._config_path}: {e}"])
        except OSError as e:
            raise ConfigValidationError([f"Error reading {self._config_path}: {e}"])

    def _parse_config(self) -> QualityGateConfigModel:
        """Parse raw TOML into typed configuration objects."""
        errors: list[str] = []

        try:
            coverage = self._parsers.parse_coverage()
            complexity = self._parsers.parse_complexity()
            file_size = self._parsers.parse_file_size()
            type_hints = self._parsers.parse_type_hints()
            error_handling = self._parsers.parse_error_handling()
            architecture = self._parsers.parse_architecture()
            documentation = self._parsers.parse_documentation()
            testing = self._parsers.parse_testing()
            naming = self._parsers.parse_naming()
            security = self._parsers.parse_security()
            performance = self._parsers.parse_performance()

            return QualityGateConfigModel(
                coverage=coverage,
                complexity=complexity,
                file_size=file_size,
                type_hints=type_hints,
                error_handling=error_handling,
                architecture=architecture,
                documentation=documentation,
                testing=testing,
                naming=naming,
                security=security,
                performance=performance,
            )

        except (ValueError, TypeError) as e:
            errors.append(f"Configuration error: {e}")
            raise ConfigValidationError(errors)

    def _get_default_config(self) -> dict[str, Any]:
        """Get default configuration when file doesn't exist."""
        return {
            "coverage": {
                "enabled": True,
                "minimum": 80,
                "fail_under": 80,
                "exclude_patterns": ["*/tests/*", "*/test_*.py"],
            },
            "complexity": {"enabled": True, "max_cc": 10, "max_average_cc": 5},
            "file_size": {"enabled": True, "max_lines": 200},
            "type_hints": {"enabled": True, "require_return_types": True},
            "error_handling": {"enabled": True, "forbid_bare_except": True},
            "architecture": {"enabled": True, "enforce_layer_separation": True},
        }

    def validate(self) -> list[str]:
        """Validate configuration and return list of errors (empty if valid)."""
        errors: list[str] = []

        # Validate coverage thresholds
        if self._config.coverage.minimum < 0 or self._config.coverage.minimum > 100:
            errors.append("coverage.minimum must be between 0 and 100")

        if self._config.coverage.fail_under < 0 or self._config.coverage.fail_under > 100:
            errors.append("coverage.fail_under must be between 0 and 100")

        # Validate complexity thresholds
        if self._config.complexity.max_cc < 1:
            errors.append("complexity.max_cc must be at least 1")

        if self._config.complexity.max_average_cc < 1:
            errors.append("complexity.max_average_cc must be at least 1")

        # Validate file size thresholds
        if self._config.file_size.max_lines < 10:
            errors.append("file_size.max_lines must be at least 10")

        if self._config.file_size.max_imports < 1:
            errors.append("file_size.max_imports must be at least 1")

        return errors
