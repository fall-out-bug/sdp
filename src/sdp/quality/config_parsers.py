"""Quality gate configuration section parsers."""

from typing import Any, cast

from sdp.quality.models import (
    ArchitectureConfig,
    ComplexityConfig,
    CoverageConfig,
    DocumentationConfig,
    ErrorHandlingConfig,
    FileSizeConfig,
    LayerPattern,
    NamingConfig,
    PerformanceConfig,
    SecurityConfig,
    TestingConfig,
    TypeHintsConfig,
)


class ConfigSectionParsers:
    """Parser methods for each configuration section."""

    def __init__(self, raw_config: dict[str, Any]):
        self._raw_config = raw_config

    def parse_coverage(self) -> CoverageConfig:
        """Parse coverage section."""
        data = self._raw_config.get("coverage", {})
        return CoverageConfig(
            enabled=cast(bool, data.get("enabled", True)),
            minimum=cast(int, data.get("minimum", 80)),
            fail_under=cast(int, data.get("fail_under", 80)),
            exclude_patterns=cast(list[str], data.get("exclude_patterns", [])),
        )

    def parse_complexity(self) -> ComplexityConfig:
        """Parse complexity section."""
        data = self._raw_config.get("complexity", {})
        return ComplexityConfig(
            enabled=cast(bool, data.get("enabled", True)),
            max_cc=cast(int, data.get("max_cc", 10)),
            max_average_cc=cast(int, data.get("max_average_cc", 5)),
        )

    def parse_file_size(self) -> FileSizeConfig:
        """Parse file_size section."""
        data = self._raw_config.get("file_size", {})
        return FileSizeConfig(
            enabled=cast(bool, data.get("enabled", True)),
            max_lines=cast(int, data.get("max_lines", 200)),
            max_imports=cast(int, data.get("max_imports", 20)),
            max_functions=cast(int, data.get("max_functions", 15)),
        )

    def parse_type_hints(self) -> TypeHintsConfig:
        """Parse type_hints section."""
        data = self._raw_config.get("type_hints", {})
        return TypeHintsConfig(
            enabled=cast(bool, data.get("enabled", True)),
            require_return_types=cast(bool, data.get("require_return_types", True)),
            require_param_types=cast(bool, data.get("require_param_types", True)),
            strict_mode=cast(bool, data.get("strict_mode", True)),
            allow_implicit_any=cast(bool, data.get("allow_implicit_any", False)),
        )

    def parse_error_handling(self) -> ErrorHandlingConfig:
        """Parse error_handling section."""
        data = self._raw_config.get("error_handling", {})
        return ErrorHandlingConfig(
            enabled=cast(bool, data.get("enabled", True)),
            forbid_bare_except=cast(bool, data.get("forbid_bare_except", True)),
            forbid_bare_raise=cast(bool, data.get("forbid_bare_raise", True)),
            forbid_pass_with_except=cast(bool, data.get("forbid_pass_with_except", True)),
            require_explicit_errors=cast(bool, data.get("require_explicit_errors", True)),
        )

    def parse_architecture(self) -> ArchitectureConfig:
        """Parse architecture section."""
        data = self._raw_config.get("architecture", {})

        # Parse layer patterns if provided
        layer_patterns = None
        layers_data = data.get("layers")
        if layers_data:
            layer_patterns = []
            for layer_name, layer_config in cast(dict[str, Any], layers_data).items():
                layer_patterns.append(
                    LayerPattern(
                        name=layer_name,
                        path_regex=cast(str, layer_config.get("path_regex", "")),
                        module_regex=cast(str | None, layer_config.get("module_regex")),
                    )
                )

        # Only set forbid_violations if explicitly in config (otherwise None triggers defaults)
        forbid_violations = None
        if "forbid_violations" in data:
            forbid_violations = cast(list[str], data["forbid_violations"])

        # Same for allowed_layer_imports
        allowed_layer_imports = None
        if "allowed_layer_imports" in data:
            allowed_layer_imports = cast(list[str], data["allowed_layer_imports"])

        return ArchitectureConfig(
            enabled=cast(bool, data.get("enabled", True)),
            enforce_layer_separation=cast(bool, data.get("enforce_layer_separation", True)),
            allowed_layer_imports=allowed_layer_imports,
            forbid_violations=forbid_violations,
            layer_patterns=layer_patterns,
        )

    def parse_documentation(self) -> DocumentationConfig | None:
        """Parse documentation section (optional)."""
        data = self._raw_config.get("documentation")
        if not data:
            return None

        return DocumentationConfig(
            enabled=cast(bool, data.get("enabled", True)),
            require_docstrings=cast(bool, data.get("require_docstrings", False)),
            min_docstring_coverage=cast(float, data.get("min_docstring_coverage", 0.5)),
            require_module_docstrings=cast(bool, data.get("require_module_docstrings", True)),
            require_class_docstrings=cast(bool, data.get("require_class_docstrings", False)),
            require_function_docstrings=cast(bool, data.get("require_function_docstrings", False)),
        )

    def parse_testing(self) -> TestingConfig | None:
        """Parse testing section (optional)."""
        data = self._raw_config.get("testing")
        if not data:
            return None

        return TestingConfig(
            enabled=cast(bool, data.get("enabled", True)),
            require_test_for_new_code=cast(bool, data.get("require_test_for_new_code", True)),
            min_test_to_code_ratio=cast(float, data.get("min_test_to_code_ratio", 0.8)),
            require_fast_marker=cast(bool, data.get("require_fast_marker", True)),
            forbid_print_statements=cast(bool, data.get("forbid_print_statements", True)),
        )

    def parse_naming(self) -> NamingConfig | None:
        """Parse naming section (optional)."""
        data = self._raw_config.get("naming")
        if not data:
            return None

        return NamingConfig(
            enabled=cast(bool, data.get("enabled", True)),
            enforce_pep8=cast(bool, data.get("enforce_pep8", True)),
            allow_single_letter=cast(bool, data.get("allow_single_letter", False)),
            min_variable_name_length=cast(int, data.get("min_variable_name_length", 3)),
            max_variable_name_length=cast(int, data.get("max_variable_name_length", 50)),
        )

    def parse_security(self) -> SecurityConfig | None:
        """Parse security section (optional)."""
        data = self._raw_config.get("security")
        if not data:
            return None

        return SecurityConfig(
            enabled=cast(bool, data.get("enabled", True)),
            forbid_hardcoded_secrets=cast(
                bool,
                data.get("forbid_hardcoded_secrets", True),
            ),
            forbid_sql_injection_patterns=cast(
                bool,
                data.get("forbid_sql_injection_patterns", True),
            ),
            forbid_eval_usage=cast(bool, data.get("forbid_eval_usage", True)),
            require_https_urls=cast(bool, data.get("require_https_urls", True)),
        )

    def parse_performance(self) -> PerformanceConfig | None:
        """Parse performance section (optional)."""
        data = self._raw_config.get("performance")
        if not data:
            return None

        return PerformanceConfig(
            enabled=cast(bool, data.get("enabled", True)),
            forbid_sql_queries_in_loops=cast(
                bool,
                data.get("forbid_sql_queries_in_loops", True),
            ),
            max_nesting_depth=cast(int, data.get("max_nesting_depth", 5)),
            warn_large_string_concatenation=cast(
                bool,
                data.get("warn_large_string_concatenation", True),
            ),
        )
