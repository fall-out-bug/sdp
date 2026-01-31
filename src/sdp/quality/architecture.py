"""Portable architecture checking for clean architecture enforcement."""

import ast
import re
from pathlib import Path

from sdp.quality.models import ArchitectureConfig
from sdp.quality.validator_models import QualityGateViolation


class ArchitectureChecker:
    """Check clean architecture layer boundaries using configurable rules."""

    def __init__(
        self,
        config: ArchitectureConfig,
        violations: list[QualityGateViolation],
    ) -> None:
        """Initialize architecture checker.

        Args:
            config: Architecture configuration from quality-gate.toml.
            violations: List to append violations to.
        """
        self._config = config
        self._violations = violations

        # Default layer patterns if not configured
        if self._config.layer_patterns is None:
            from sdp.quality.models import LayerPattern

            self._config.layer_patterns = [
                LayerPattern(
                    name="domain",
                    path_regex=r"(^|/)domain/",
                    module_regex=r"(^|\.)domain(\.|$)",
                ),
                LayerPattern(
                    name="application",
                    path_regex=r"(^|/)application/",
                    module_regex=r"(^|\.)application(\.|$)",
                ),
                LayerPattern(
                    name="infrastructure",
                    path_regex=r"(^|/)infrastructure/",
                    module_regex=r"(^|\.)infrastructure(\.|$)",
                ),
                LayerPattern(
                    name="presentation",
                    path_regex=r"(^|/)presentation/",
                    module_regex=r"(^|\.)presentation(\.|$)",
                ),
            ]

        # Build layer detection patterns from configuration
        self._layer_patterns = self._build_layer_patterns()

    def check_architecture(self, path: Path, tree: ast.AST) -> None:
        """Check file for architecture violations.

        Args:
            path: Path to file being checked.
            tree: AST of the file.
        """
        if not self._config.enabled:
            return

        # Detect which layer this file belongs to
        file_layer = self._detect_layer(path)

        if file_layer is None:
            # Not in a tracked layer
            return

        # Check imports for violations
        for node in ast.walk(tree):
            if isinstance(node, ast.Import):
                self._check_import_node(path, node.names[0].name, file_layer, node.lineno)
            elif isinstance(node, ast.ImportFrom):
                if node.module:
                    self._check_import_node(path, node.module, file_layer, node.lineno)

    def _detect_layer(self, path: Path) -> str | None:
        """Detect which architectural layer a file belongs to.

        Uses regex patterns derived from quality-gate.toml configuration.
        Common patterns: domain, application, infrastructure, presentation.

        Args:
            path: Path to file.

        Returns:
            Layer name or None if not detected.
        """
        path_str = str(path)

        # Check against layer patterns
        for layer, (path_pattern, _module_pattern) in self._layer_patterns.items():
            if re.search(path_pattern, path_str):
                return layer

        return None

    def _build_layer_patterns(self) -> dict[str, tuple[str, str | None]]:
        """Build regex patterns for layer detection from config.

        Returns:
            Dict mapping layer names to (path_regex, module_regex) tuples.
        """
        patterns = {}
        if self._config.layer_patterns:
            for layer in self._config.layer_patterns:
                patterns[layer.name] = (layer.path_regex, layer.module_regex)
        return patterns

    def _check_import_node(
        self,
        file_path: Path,
        import_module: str,
        source_layer: str,
        lineno: int | None,
    ) -> None:
        """Check if an import violates architecture rules.

        Args:
            file_path: Path to file with the import.
            import_module: Module being imported.
            source_layer: Layer of the source file.
            lineno: Line number of import.
        """
        # Detect which layer the imported module belongs to
        # This is a simplified check - in real code, would resolve module paths
        import_layer = self._detect_import_layer(import_module)

        if import_layer is None:
            # Not importing from a tracked layer
            return

        # Check against forbidden patterns
        violation_pattern = f"{source_layer} -> {import_layer}"

        forbid_violations = self._config.forbid_violations
        if forbid_violations is not None and violation_pattern in forbid_violations:
            self._violations.append(
                QualityGateViolation(
                    category="architecture",
                    file_path=str(file_path),
                    line_number=lineno,
                    message=(
                        f"Architecture violation: {source_layer} cannot "
                        f"import from {import_layer}"
                    ),
                    severity="error",
                )
            )

    def _detect_import_layer(self, import_module: str) -> str | None:
        """Detect which layer an imported module belongs to.

        This is a heuristic - assumes module names follow conventions like:
        - myapp.domain.* -> domain
        - myapp.application.* -> application
        - etc.

        Args:
            import_module: Module name from import statement.

        Returns:
            Layer name or None if not detected.
        """
        for layer, (_path_pattern, module_pattern) in self._layer_patterns.items():
            if module_pattern and re.search(module_pattern, import_module):
                return layer
            # Fallback: convert path pattern to module pattern
            if not module_pattern:
                module_pattern_fallback = _path_pattern.replace("/", ".").replace(
                    "(^|)", "(^|.)"
                )
                if re.search(module_pattern_fallback, import_module):
                    return layer

        return None


class LayerViolation:
    """Represents a layer architecture violation."""

    def __init__(
        self,
        source_layer: str,
        target_layer: str,
        file_path: str,
        line: int | None,
        import_statement: str,
    ) -> None:
        """Initialize violation.

        Args:
            source_layer: Layer doing the importing.
            target_layer: Layer being imported.
            file_path: Path to file with violation.
            line: Line number of violation.
            import_statement: The import statement.
        """
        self.source_layer = source_layer
        self.target_layer = target_layer
        self.file_path = file_path
        self.line = line
        self.import_statement = import_statement

    def __str__(self) -> str:
        """String representation of violation."""
        return (
            f"Architecture violation in {self.file_path}:{self.line or '?'}\n"
            f"  {self.source_layer} -> {self.target_layer}\n"
            f"  Import: {self.import_statement}"
        )
